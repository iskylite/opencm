package collector

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/transport"
	"github.com/iskylite/procfs"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	defMPILibIncluded    = "libmpi.so"
	defSlurmstepdExe     = "slurmstepd"
	defProcessesExcluded = "^(slurmd)$"
	defMaxRssLimit       = "1048576"
	processesFlag        = 1 << 8
	noJobID              = ""
	unknown              = "unknown"
)

var (
	mpiLibIncluded = kingpin.Flag(
		"collector.processes.mpi-lib-include",
		"Regexp of mpi lib to Include for processes collector.",
	).Default(defMPILibIncluded).String()
	processesExcluded = kingpin.Flag("collector.processes.exclude",
		"exclude some processes for no checking").Default(defProcessesExcluded).String()
	slurmstepdExe = kingpin.Flag("collector.processes.slurmstepd",
		"slurmstepd process name").Default(defSlurmstepdExe).String()
	maxRssLimit = kingpin.Flag("collector.processes.max-rss-limit",
		"max rss limit (unit: kB) for filter all process").Default(defMaxRssLimit).Int()
)

type processCollector struct {
	subsystem                string
	slurmstepdExe            string
	mpiLibIncludePattern     *regexp.Regexp
	processesExcludedPattern *regexp.Regexp
	fs                       procfs.FS
	logger                   log.Logger
}

func init() {
	registerCollector("processes", processesFlag, NewProcessStatCollector)
}

// NewProcessStatCollector returns a new Collector exposing process data read from the proc filesystem.
func NewProcessStatCollector(logger log.Logger, subsystem string) (Collector, error) {
	level.Info(logger).Log("msg", "Parsed flag --collector.processes.mpi-lib-include", "flag", *mpiLibIncluded)
	mpiLibIncludePattern := regexp.MustCompile(*mpiLibIncluded)
	level.Info(logger).Log("msg", "Parsed flag --collector.processes.exclude", "flag", *processesExcluded)
	processesExcludedPattern := regexp.MustCompile(*processesExcluded)
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &processCollector{
		subsystem:                subsystem,
		slurmstepdExe:            *slurmstepdExe,
		mpiLibIncludePattern:     mpiLibIncludePattern,
		processesExcludedPattern: processesExcludedPattern,
		fs:                       fs,
		logger:                   logger,
	}, nil
}
func (c *processCollector) Update(ch chan<- *transport.Data) error {
	slurmstepdProcStats, allProcessesStats, rtime, err := c.getAllProcs()
	if err != nil {
		return err
	}
	err = c.parseProcesses(slurmstepdProcStats, allProcessesStats, func(stat *procStat, jobid string) {
		ch <- &transport.Data{
			Time:        rtime,
			Measurement: c.subsystem,
			Tags:        stat.tags(),
			Fields:      stat.fields(),
		}
	})
	return err
}

type procStat struct {
	jobid      string
	includeMPI bool
	uid        string
	gid        string
	io         procfs.ProcIO
	stat       procfs.ProcStat
	rss        float64
	virtual    float64
}

func (p *procStat) tags() map[string]string {
	return map[string]string{
		"jobid":   p.jobid,
		"pid":     strconv.Itoa(p.stat.PID),
		"ppid":    strconv.Itoa(p.stat.PPID),
		"exe":     p.stat.Comm,
		"state":   p.stat.State,
		"uid":     p.uid,
		"gid":     p.gid,
		"cpu":     strconv.Itoa(int(p.stat.Processor)),
		"threads": strconv.Itoa(int(p.stat.NumThreads)),
	}
}

func (p *procStat) fields() map[string]float64 {
	return map[string]float64{
		"read_char":         float64(p.io.RChar),
		"write_char":        float64(p.io.WChar),
		"read_bytes":        float64(p.io.ReadBytes),
		"write_bytes":       float64(p.io.WriteBytes),
		"rss_mem_bytes":     p.rss,
		"virtual_mem_bytes": p.virtual,
	}
}

// first return include slurmstepd, second return include other processes
func (c *processCollector) getAllProcs() (map[int]*procStat, []*procStat, int64, error) {
	p, err := c.fs.AllProcs()
	if err != nil {
		return nil, nil, 0, fmt.Errorf("unable to list all processes: %w", err)
	}

	// 保存slurmstepd进程信息
	slurmstepdProcs := make(map[int]*procStat, 1)
	// 保存除了slurm以外的其他进程信息
	processes := make([]*procStat, 0, 150)

	rtime := utils.Now()

	for _, pid := range p {
		stat, err := pid.Stat()
		if err != nil {
			// PIDs can vanish between getting the list and getting stats.
			if c.isIgnoredError(err) {
				level.Debug(c.logger).Log("msg", "file not found when retrieving stats for pid", "pid", pid.PID, "err", err)
				continue
			}
			level.Debug(c.logger).Log("msg", "error reading stat for pid", "pid", pid.PID, "err", err)
			return nil, nil, 0, fmt.Errorf("error reading stat for pid %d: %w", pid.PID, err)
		}
		// status
		status, err := pid.NewStatus()
		if err != nil {
			if c.isIgnoredError(err) {
				level.Debug(c.logger).Log("msg", "file not found when retrieving status for pid", "pid", pid.PID, "err", err)
				continue
			}
			level.Debug(c.logger).Log("msg", "error reading status for pid", "pid", pid.PID, "err", err)
			return nil, nil, 0, fmt.Errorf("error reading status for pid %d: %w", pid.PID, err)
		}
		// io
		io, err := pid.IO()
		if err != nil {
			if c.isIgnoredError(err) {
				level.Debug(c.logger).Log("msg", "file not found when retrieving io for pid", "pid", pid.PID, "err", err)
				continue
			}
			level.Debug(c.logger).Log("msg", "error reading io for pid", "pid", pid.PID, "err", err)
			return nil, nil, 0, fmt.Errorf("error reading io for pid %d: %w", pid.PID, err)
		}
		// mem
		rss := float64(stat.ResidentMemory())
		virtual := float64(stat.VirtualMemory())
		// 提取slurmstepd进程
		if stat.Comm == c.slurmstepdExe {
			cmds, err := pid.CmdLine()
			if err != nil {
				if c.isIgnoredError(err) {
					level.Debug(c.logger).Log("msg", "file not found when retrieving cmdline for pid", "pid", pid.PID, "err", err)
					continue
				}
				level.Debug(c.logger).Log("msg", "error reading cmdline for pid", "pid", pid.PID, "err", err)
				return nil, nil, 0, fmt.Errorf("error reading cmdline for pid %d: %w", pid.PID, err)
			}
			jobid := parseSlurmstepdCmdline(cmds[0])
			slurmstepdProcs[stat.PID] = &procStat{jobid: jobid, includeMPI: false, stat: stat, uid: status.UIDs[0], gid: status.GIDs[0], io: io, rss: rss, virtual: virtual}
			continue
		}
		// 过滤不需要考虑的进程，比如slurmd进程
		if c.processesExcludedPattern.MatchString(stat.Comm) {
			continue
		}
		// 其他进程
		maps, err := pid.ProcMaps()
		if err != nil {
			// map获取失败，认为不是mpi程序
			if c.isIgnoredError(err) {
				processes = append(processes, &procStat{
					jobid:      noJobID,
					includeMPI: false,
					stat:       stat,
					io:         io,
					uid:        status.UIDs[0],
					gid:        status.GIDs[0],
					rss:        rss,
					virtual:    virtual,
				})
				level.Debug(c.logger).Log("msg", "file not found when retrieving maps for pid", "pid", pid.PID, "err", err)
				continue
			}
			level.Debug(c.logger).Log("msg", "error reading maps for pid", "pid", pid.PID, "err", err)
			return nil, nil, 0, fmt.Errorf("error reading maps for pid %d: %w", pid.PID, err)
		}
		for i := 0; i < len(maps); i++ {
			if c.mpiLibIncludePattern.MatchString(maps[i].Pathname) {
				// 如果map中检查到mpi，则认为是mpi程序
				processes = append(processes, &procStat{
					jobid:      noJobID,
					includeMPI: true,
					stat:       stat,
					io:         io,
					uid:        status.UIDs[0],
					gid:        status.GIDs[0],
					rss:        rss,
					virtual:    virtual,
				})
				break
			}
		}
		// 普通程序
		processes = append(processes, &procStat{
			jobid:      noJobID,
			includeMPI: false,
			stat:       stat,
			io:         io,
			uid:        status.UIDs[0],
			gid:        status.GIDs[0],
			rss:        rss,
			virtual:    virtual,
		})
	}
	return slurmstepdProcs, processes, rtime, nil
}

func (c *processCollector) isIgnoredError(err error) bool {
	if errors.Is(err, os.ErrNotExist) || strings.Contains(err.Error(), syscall.ESRCH.Error()) {
		return true
	}
	return false
}

func (c *processCollector) parseProcesses(slurmstepdProcs map[int]*procStat, allProcesses []*procStat, handler func(stat *procStat, jobid string)) error {
	for _, stat := range allProcesses {
		ppid := stat.stat.PPID
		if slurmstepdStat, ok := slurmstepdProcs[ppid]; ok {
			// slurmstepd衍生的并行作业
			handler(stat, slurmstepdStat.jobid)
			continue
		}
		if stat.includeMPI {
			// 非slurm提交的mpi进程
			handler(stat, stat.jobid)
			continue
		}
		// 实存大于maxRssLimit
		if stat.rss > float64(*maxRssLimit*1024) {
			handler(stat, stat.jobid)
		}
	}
	return nil
}

func parseSlurmstepdCmdline(cmdline string) string {
	cmdlinePattern := regexp.MustCompile("[0-9]+")
	return cmdlinePattern.FindString(cmdline)
}
