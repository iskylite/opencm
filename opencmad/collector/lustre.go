package collector

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/pb"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	lustreFlag = 1 << 11
)

var (
	healthCheckPath = kingpin.Flag("collector.lustre.health_check_path",
		"lustre health_check top dir path").Default("/sys/fs/lustre").String()
	targetSizePath = kingpin.Flag("collector.lustre.target_size_path",
		"lustre target size top dir path").Default("/proc/fs/lustre").String()
	llitePath = kingpin.Flag("collector.lustre.llite_path",
		"lustre llite dir path").Default("/sys/kernel/debug/lustre/llite").String()
	lnetPath = kingpin.Flag("collector.lustre.lnet_path",
		"lustre lnet dir path").Default("/sys/kernel/debug/lnet").String()
	collectInterval = kingpin.Flag("collector.lustre.collect_interval",
		"lustre collect interval, unit: ms").Default("500").Int()
)

var errLustreNotAvailable = errors.New("Lustre statistics are not available")

func init() {
	registerCollector("lustre", lustreFlag, NewLustreCollector)
}

type lustreCollector struct {
	subsystem       string
	healthCheckPath string
	targetSizePath  string
	llitePath       string
	lnetPath        string
	logger          log.Logger
}

// NewLustreCollector returns a new Collector exposing Lustre statistics.
func NewLustreCollector(logger log.Logger, subsystem string) (Collector, error) {
	return &lustreCollector{
		subsystem:       subsystem,
		healthCheckPath: *healthCheckPath,
		targetSizePath:  *targetSizePath,
		llitePath:       *llitePath,
		lnetPath:        *lnetPath,
		logger:          logger,
	}, nil
}

func (l *lustreCollector) Update(ch chan<- *pb.CollectData) error {
	if _, err := l.openProcFile(l.healthCheckPath); err != nil {
		if err == errLustreNotAvailable {
			level.Debug(l.logger).Log("err", err)
			return ErrNoData
		}
	}
	// updateLliteStats、updateOSTStats、updateLnetStats 采集0.5s内的数据并做处理，近似瞬时值
	// updateTargetJobStats 暂时不做修改
	updateFuncs := []func(chan<- *pb.CollectData) error{
		l.updateLustreHealth,
		l.updateMGTState,
		l.updateMDTState,
		l.updateOSTState,
		l.updateTargetSize,
		l.updateMDTStats,
		l.updateOSTStats, // test ok
		l.updateOSTBRWStats,
		l.updateTargetJobStats,
		l.updateTargetQuotaSlave,
		l.updateClientState,
		l.updateLliteStats, // test ok
		l.updateLnetNis,
		l.updateLnetStats, // test ok
		l.updateLnetMemUsed,
	}
	var wg sync.WaitGroup
	wg.Add(len(updateFuncs))
	for _, updateFunc := range updateFuncs {
		go func(updateFunc func(chan<- *pb.CollectData) error) {
			if err := updateFunc(ch); err != nil {
				level.Debug(l.logger).Log("msg", err, "subcollector", runtime.FuncForPC(reflect.ValueOf(updateFunc).Pointer()).Name())
			}
			level.Debug(l.logger).Log("msg", "lustre subcollector success", "subcollector", runtime.FuncForPC(reflect.ValueOf(updateFunc).Pointer()).Name())
			wg.Done()
		}(updateFunc)
	}
	wg.Wait()
	return nil
}

func (l *lustreCollector) openProcFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		level.Debug(l.logger).Log("msg", "Cannot open file for reading", "path", path)
		return nil, errLustreNotAvailable
	}
	return file, nil
}

func (l *lustreCollector) updateLustreHealth(ch chan<- *pb.CollectData) error {
	healthCheck, err := utils.ReadAll(filepath.Join(l.healthCheckPath, "health_check"))
	if err != nil {
		return err
	}
	version, err := utils.ReadAll(sysFilePath("fs/lustre/version"))
	if err != nil {
		return err
	}
	tags := map[string]string{
		"health":  strings.TrimSpace(string(healthCheck)),
		"version": strings.TrimSpace(string(version)),
	}
	fields := make(map[string]float64, 3)
	for _, key := range []string{"max_dirty_mb", "memused", "memused_max"} {
		data, err := utils.ReadAll(sysFilePath("fs/lustre/" + key))
		if err != nil {
			fields[key] = 0.0
			level.Debug(l.logger).Log("msg", err)
			continue
		}
		fdata, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		if err != nil {
			fields[key] = 0.0
			level.Debug(l.logger).Log("msg", err)
			continue
		}
		fields[key] = float64(fdata)
	}
	ch <- &pb.CollectData{
		Time:        utils.Now(),
		Measurement: "lustre_health",
		Tags:        tags,
		Fields:      fields,
	}
	return nil
}
