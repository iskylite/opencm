package collector

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/transport"
	"github.com/iskylite/procfs"
)

type cpuCollector struct {
	subsystem     string
	fs            procfs.FS
	logger        log.Logger
	cpuStats      []procfs.CPUStat
	cpuStatsMutex sync.Mutex
}

// Idle jump back limit in seconds.
const (
	jumpBackSeconds = 3.0
	cpuFlag         = 1 << 0
)

var (
	jumpBackDebugMessage = fmt.Sprintf("CPU Idle counter jumped backwards more than %f seconds, possible hotplug event, resetting CPU stats", jumpBackSeconds)
)

func init() {
	registerCollector("cpu", cpuFlag, NewCPUCollector)
}

// NewCPUCollector returns a new Collector exposing kernel/system statistics.
func NewCPUCollector(logger log.Logger, subsystem string) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	c := &cpuCollector{
		subsystem: subsystem,
		fs:        fs,
		logger:    logger,
	}
	return c, nil
}

// Update implements Collector and exposes cpu related metrics from /proc/stat and /sys/.../cpu/.
func (c *cpuCollector) Update(ch chan<- *transport.CollectData) error {
	if err := c.updateStat(ch); err != nil {
		return err
	}
	return nil
}

// updateStat reads /proc/stat through procfs and exports CPU-related metrics.
func (c *cpuCollector) updateStat(ch chan<- *transport.CollectData) error {
	rtime := utils.Now()
	stats, err := c.fs.Stat()
	if err != nil {
		return err
	}
	c.updateCPUStats(stats.CPU)

	// Acquire a lock to read the stats.
	c.cpuStatsMutex.Lock()
	defer c.cpuStatsMutex.Unlock()
	ch <- NewCollectorData(rtime, c.subsystem, map[string]string{"cpu": "cpu"}, newCPUDataFields(stats.CPUTotal))
	for cpuID, cpuStat := range c.cpuStats {
		cpuNum := strconv.Itoa(cpuID)
		ch <- NewCollectorData(rtime, c.subsystem, map[string]string{"cpu": "cpu" + cpuNum}, newCPUDataFields(cpuStat))
	}

	return nil
}

// updateCPUStats updates the internal cache of CPU stats.
func (c *cpuCollector) updateCPUStats(newStats []procfs.CPUStat) {

	// Acquire a lock to update the stats.
	c.cpuStatsMutex.Lock()
	defer c.cpuStatsMutex.Unlock()

	// Reset the cache if the list of CPUs has changed.
	if len(c.cpuStats) != len(newStats) {
		c.cpuStats = make([]procfs.CPUStat, len(newStats))
	}

	for i, n := range newStats {
		// If idle jumps backwards by more than X seconds, assume we had a hotplug event and reset the stats for this CPU.
		if (c.cpuStats[i].Idle - n.Idle) >= jumpBackSeconds {
			level.Debug(c.logger).Log("msg", jumpBackDebugMessage, "cpu", i, "old_value", c.cpuStats[i].Idle, "new_value", n.Idle)
			c.cpuStats[i] = procfs.CPUStat{}
		}

		if n.Idle >= c.cpuStats[i].Idle {
			c.cpuStats[i].Idle = n.Idle
		} else {
			level.Debug(c.logger).Log("msg", "CPU Idle counter jumped backwards", "cpu", i, "old_value", c.cpuStats[i].Idle, "new_value", n.Idle)
		}

		if n.User >= c.cpuStats[i].User {
			c.cpuStats[i].User = n.User
		} else {
			level.Debug(c.logger).Log("msg", "CPU User counter jumped backwards", "cpu", i, "old_value", c.cpuStats[i].User, "new_value", n.User)
		}

		if n.Nice >= c.cpuStats[i].Nice {
			c.cpuStats[i].Nice = n.Nice
		} else {
			level.Debug(c.logger).Log("msg", "CPU Nice counter jumped backwards", "cpu", i, "old_value", c.cpuStats[i].Nice, "new_value", n.Nice)
		}

		if n.System >= c.cpuStats[i].System {
			c.cpuStats[i].System = n.System
		} else {
			level.Debug(c.logger).Log("msg", "CPU System counter jumped backwards", "cpu", i, "old_value", c.cpuStats[i].System, "new_value", n.System)
		}

		if n.Iowait >= c.cpuStats[i].Iowait {
			c.cpuStats[i].Iowait = n.Iowait
		} else {
			level.Debug(c.logger).Log("msg", "CPU Iowait counter jumped backwards", "cpu", i, "old_value", c.cpuStats[i].Iowait, "new_value", n.Iowait)
		}

		if n.IRQ >= c.cpuStats[i].IRQ {
			c.cpuStats[i].IRQ = n.IRQ
		} else {
			level.Debug(c.logger).Log("msg", "CPU IRQ counter jumped backwards", "cpu", i, "old_value", c.cpuStats[i].IRQ, "new_value", n.IRQ)
		}

		if n.SoftIRQ >= c.cpuStats[i].SoftIRQ {
			c.cpuStats[i].SoftIRQ = n.SoftIRQ
		} else {
			level.Debug(c.logger).Log("msg", "CPU SoftIRQ counter jumped backwards", "cpu", i, "old_value", c.cpuStats[i].SoftIRQ, "new_value", n.SoftIRQ)
		}

		if n.Steal >= c.cpuStats[i].Steal {
			c.cpuStats[i].Steal = n.Steal
		} else {
			level.Debug(c.logger).Log("msg", "CPU Steal counter jumped backwards", "cpu", i, "old_value", c.cpuStats[i].Steal, "new_value", n.Steal)
		}
	}
}

func newCPUDataFields(stat procfs.CPUStat) map[string]float64 {
	return map[string]float64{
		"user":    stat.User,
		"nice":    stat.Nice,
		"system":  stat.System,
		"idle":    stat.Idle,
		"iowait":  stat.Iowait,
		"irq":     stat.IRQ,
		"softirq": stat.SoftIRQ,
		"steal":   stat.Steal,
	}
}
