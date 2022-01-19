package collector

import (
	"regexp"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/transport"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Arch-dependent implementation must define:
// * defMountPointsExcluded
// * defFSTypesExcluded
// * filesystemLabelNames
// * filesystemCollector.GetStats

const (
	filesystemFlag = 1 << 3
)

var (
	mountPointsExcludeSet bool
	mountPointsExclude    = kingpin.Flag(
		"collector.filesystem.mount-points-exclude",
		"Regexp of mount points to exclude for filesystem collector.",
	).Default(defMountPointsExcluded).PreAction(func(c *kingpin.ParseContext) error {
		mountPointsExcludeSet = true
		return nil
	}).String()
	fsTypesIncludeSet bool
	fsTypesInclude    = kingpin.Flag(
		"collector.filesystem.fs-types-Include",
		"Regexp of filesystem types to Include for filesystem collector.",
	).Default(defFSTypesIncluded).PreAction(func(c *kingpin.ParseContext) error {
		fsTypesIncludeSet = true
		return nil
	}).String()
	filesystemLabelNames = []string{"device", "mountpoint", "fstype"}
)

type filesystemCollector struct {
	excludedMountPointsPattern *regexp.Regexp
	includedFSTypesPattern     *regexp.Regexp
	logger                     log.Logger
	subsystem                  string
}

type filesystemLabels struct {
	device, mountPoint, fsType, options string
}

type filesystemStats struct {
	labels            filesystemLabels
	size, free, avail float64
	files, filesFree  float64
	ro, deviceError   float64
}

func init() {
	registerCollector("filesystem", filesystemFlag, NewFilesystemCollector)
}

// NewFilesystemCollector returns a new Collector exposing filesystems stats.
func NewFilesystemCollector(logger log.Logger, subsystem string) (Collector, error) {
	level.Info(logger).Log("msg", "Parsed flag --collector.filesystem.fs-types-Include", "flag", *fsTypesInclude)
	filesystemsTypesPattern := regexp.MustCompile(*fsTypesInclude)
	level.Info(logger).Log("msg", "Parsed flag --collector.filesystem.mount-points-exclude", "flag", *mountPointsExclude)
	mountPointPattern := regexp.MustCompile(*mountPointsExclude)
	return &filesystemCollector{
		excludedMountPointsPattern: mountPointPattern,
		includedFSTypesPattern:     filesystemsTypesPattern,
		logger:                     logger,
		subsystem:                  subsystem,
	}, nil
}

func (c *filesystemCollector) Update(ch chan<- *transport.Data) error {
	rtime := utils.Now()
	stats, err := c.GetStats()
	if err != nil {
		return err
	}
	// Make sure we expose a metric once, even if there are multiple mounts
	seen := map[filesystemLabels]bool{}
	for _, s := range stats {
		if seen[s.labels] {
			continue
		}
		seen[s.labels] = true

		tags := map[string]string{
			"device":  s.labels.device,
			"mp":      s.labels.mountPoint,
			"fs_type": s.labels.fsType,
		}
		fields := map[string]float64{
			"device_error": s.deviceError,
		}
		if s.deviceError == 0 {
			fields["size"] = s.size
			fields["free"] = s.free
			fields["avail"] = s.avail
			fields["size_percent"] = (s.size - s.free) / s.size
			fields["files"] = s.files
			fields["files_free"] = s.filesFree
			fields["files_percent"] = (s.files - s.filesFree) / s.files
			fields["ro"] = s.ro
		}
		ch <- &transport.Data{
			Time:        rtime,
			Measurement: c.subsystem,
			Tags:        tags,
			Fields:      fields,
		}
	}
	return nil
}
