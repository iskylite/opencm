package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/iskylite/opencm/pb"
	"github.com/iskylite/opencm/pkg/timex"
)

const loadFlag = 1 << 5

type loadavgCollector struct {
	logger    log.Logger
	subsystem string
}

func init() {
	registerCollector("loadavg", loadFlag, NewLoadavgCollector)
}

// NewLoadavgCollector returns a new Collector exposing load average stats.
func NewLoadavgCollector(logger log.Logger, subsystem string) (Collector, error) {
	return &loadavgCollector{
		logger:    logger,
		subsystem: subsystem,
	}, nil
}

func (c *loadavgCollector) Update(ch chan<- *pb.CollectData) error {
	rtime := timex.Now()
	loads, err := getLoad()
	if err != nil {
		return fmt.Errorf("couldn't get load: %w", err)
	}
	ch <- &pb.CollectData{
		Time:        rtime,
		Measurement: c.subsystem,
		Tags:        nil,
		Fields: map[string]float64{
			"load1":  loads[0],
			"load5":  loads[1],
			"load15": loads[2],
		},
	}
	return err
}
