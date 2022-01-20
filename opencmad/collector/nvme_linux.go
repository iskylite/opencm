package collector

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/transport"
	"github.com/iskylite/procfs/sysfs"
)

const nvmeFlag = 1 << 7

type nvmeCollector struct {
	fs        sysfs.FS
	logger    log.Logger
	subsystem string
}

func init() {
	registerCollector("nvme", nvmeFlag, NewNVMeCollector)
}

// NewNVMeCollector returns a new Collector exposing NVMe stats.
func NewNVMeCollector(logger log.Logger, subsystem string) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &nvmeCollector{
		fs:        fs,
		logger:    logger,
		subsystem: subsystem,
	}, nil
}

func (c *nvmeCollector) Update(ch chan<- *transport.CollectData) error {
	rtime := utils.Now()
	devices, err := c.fs.NVMeClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "nvme statistics not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining NVMe class info: %w", err)
	}

	for _, device := range devices {
		ch <- &transport.CollectData{
			Time:        rtime,
			Measurement: c.subsystem,
			Tags:        map[string]string{"dev": device.Name, "state": device.State, "fw": device.FirmwareRevision, "model": device.Model, "serial": device.Serial},
			Fields:      nil,
		}
	}

	return nil
}
