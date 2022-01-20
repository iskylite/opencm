package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/transport"
)

const (
	memFlag = 1 << 6
)

type meminfoCollector struct {
	logger    log.Logger
	subsystem string
}

func init() {
	registerCollector("meminfo", memFlag, NewMeminfoCollector)
}

// NewMeminfoCollector returns a new Collector exposing memory stats.
func NewMeminfoCollector(logger log.Logger, subsystem string) (Collector, error) {
	return &meminfoCollector{logger, subsystem}, nil
}

// Update calls (*meminfoCollector).getMemInfo to get the platform specific
// memory metrics.
func (c *meminfoCollector) Update(ch chan<- *transport.CollectData) error {
	rtime := utils.Now()
	memInfo, err := c.getMemInfo()
	if err != nil {
		return fmt.Errorf("couldn't get meminfo: %w", err)
	}
	level.Debug(c.logger).Log("msg", "Set node_mem", "memInfo", memInfo)

	ch <- &transport.CollectData{
		Time:        rtime,
		Measurement: c.subsystem,
		Tags:        nil,
		Fields: map[string]float64{
			"total":      memInfo["MemTotal_bytes"],
			"free":       memInfo["MemFree_bytes"],
			"avail":      memInfo["MemAvailable_bytes"],
			"cache":      memInfo["Cached_bytes"],
			"buffers":    memInfo["Buffers_bytes"],
			"swap_total": memInfo["SwapTotal_bytes"],
			"swap_free":  memInfo["SwapFree_bytes"],
			"swap_cache": memInfo["SwapCached_bytes"],
		},
	}
	return nil
}
