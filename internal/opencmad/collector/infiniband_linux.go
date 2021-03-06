package collector

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/internal/pkg/defs"
	"github.com/iskylite/opencm/pb"
	"github.com/iskylite/opencm/pkg/timex"
	"github.com/iskylite/procfs/sysfs"
)

const ibFlag = 1 << 4

type infinibandCollector struct {
	fs        sysfs.FS
	logger    log.Logger
	subsystem string
}

func init() {
	registerCollector("infiniband", ibFlag, NewInfiniBandCollector)
}

// NewInfiniBandCollector returns a new Collector exposing InfiniBand stats.
func NewInfiniBandCollector(logger log.Logger, subsystem string) (Collector, error) {
	var i infinibandCollector
	var err error

	i.fs, err = sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}
	i.logger = logger
	i.subsystem = subsystem
	return &i, nil
}

func (c *infinibandCollector) Update(ch chan<- *pb.CollectData) error {
	rtime := timex.Now()
	devices, err := c.fs.InfiniBandClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "infiniband statistics not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining InfiniBand class info: %w", err)
	}

	for _, device := range devices {
		for _, port := range device.Ports {
			portStr := strconv.FormatUint(uint64(port.Port), 10)
			ch <- &pb.CollectData{
				Time:        rtime,
				Measurement: c.subsystem,
				Tags: map[string]string{
					"device":    device.Name,
					"board_id":  device.BoardID,
					"fw":        device.FirmwareVersion,
					"hca_type":  device.HCAType,
					"port":      portStr,
					"state":     defs.IBState(port.StateID).String(),
					"phy_state": defs.IBPhysicalState(port.PhysStateID).String(),
				},
			}
		}
	}

	return nil
}
