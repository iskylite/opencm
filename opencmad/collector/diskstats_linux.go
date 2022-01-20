package collector

import (
	"fmt"
	"regexp"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/transport"
	"github.com/iskylite/procfs/blockdevice"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	secondsPerTick = 1.0 / 1000.0
	diskstatsFlag  = 1 << 1
)

var (
	ignoredDevices = kingpin.Flag("collector.diskstats.ignored-devices", "Regexp of devices to ignore for diskstats.").Default("^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+|sr\\d+$").String()
)

type diskstatsCollector struct {
	subsystem             string
	ignoredDevicesPattern *regexp.Regexp
	fs                    blockdevice.FS
	logger                log.Logger
}

func init() {
	registerCollector("diskstats", diskstatsFlag, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
// Docs from https://www.kernel.org/doc/Documentation/iostats.txt
func NewDiskstatsCollector(logger log.Logger, subsystem string) (Collector, error) {
	fs, err := blockdevice.NewFS(*procPath, *sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &diskstatsCollector{
		subsystem:             subsystem,
		ignoredDevicesPattern: regexp.MustCompile(*ignoredDevices),
		fs:                    fs,
		logger:                logger,
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- *transport.CollectData) error {
	rtime := utils.Now()
	diskStats, err := c.fs.ProcDiskstats()
	if err != nil {
		return fmt.Errorf("couldn't get diskstats: %w", err)
	}
	for _, stats := range diskStats {
		dev := stats.DeviceName
		if c.ignoredDevicesPattern.MatchString(dev) {
			level.Debug(c.logger).Log("msg", "Ignoring device", "device", dev, "pattern", c.ignoredDevicesPattern)
			continue
		}

		diskSectorSize := 512.0
		blockQueue, err := c.fs.SysBlockDeviceQueueStats(dev)
		if err != nil {
			level.Debug(c.logger).Log("msg", "Error getting queue stats", "device", dev, "err", err)
		} else {
			diskSectorSize = float64(blockQueue.LogicalBlockSize)
		}

		ch <- NewCollectorData(rtime, c.subsystem, map[string]string{"dev": dev}, newDiskStatsFields(stats, diskSectorSize))

	}
	return nil
}

func newDiskStatsFields(stats blockdevice.Diskstats, diskSectorSize float64) map[string]float64 {
	return map[string]float64{
		"read_complete_total_cnts":          float64(stats.ReadIOs),
		"read_merged_total_cnts":            float64(stats.ReadMerges),
		"read_total_bytes":                  float64(stats.ReadSectors) * diskSectorSize,
		"read_time_total_seconds":           float64(stats.ReadTicks) * secondsPerTick,
		"write_complete_total_cnts":         float64(stats.WriteIOs),
		"write_merged_total_cnts":           float64(stats.WriteMerges),
		"write_total_bytes":                 float64(stats.WriteSectors) * diskSectorSize,
		"write_time_total_seconds":          float64(stats.WriteTicks) * secondsPerTick,
		"io_now_procs_cnts":                 float64(stats.IOsInProgress),
		"io_time_total_seconds":             float64(stats.IOsTotalTicks) * secondsPerTick,
		"io_time_weighted_total_seconds":    float64(stats.WeightedIOTicks) * secondsPerTick,
		"discard_complete_total_cnts":       float64(stats.DiscardIOs),
		"discard_merged_total_cnts":         float64(stats.DiscardMerges),
		"discard_total_sectors":             float64(stats.DiscardSectors),
		"discard_time_total_seconds":        float64(stats.DiscardTicks) * secondsPerTick,
		"flush_requests_total_cnts":         float64(stats.FlushRequestsCompleted),
		"flush_requests_time_total_seconds": float64(stats.TimeSpentFlushing) * secondsPerTick,
	}
}
