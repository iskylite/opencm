package collector

import (
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/pb"
	"github.com/iskylite/procfs/sysfs"
)

const (
	maxUint64 = ^uint64(0)
	fcFlag    = 1 << 2
)

type fibrechannelCollector struct {
	fs        sysfs.FS
	logger    log.Logger
	subsystem string
}

func init() {
	registerCollector("fibrechannel", fcFlag, NewFibreChannelCollector)
}

// NewFibreChannelCollector returns a new Collector exposing FibreChannel stats.
func NewFibreChannelCollector(logger log.Logger, subsystem string) (Collector, error) {
	var i fibrechannelCollector
	var err error

	i.fs, err = sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}
	i.logger = logger

	// Detailed description for all metrics.
	// descriptions := map[string]string{
	// 	"dumped_frames_total":            "Number of dumped frames",
	// 	"loss_of_signal_total":           "Number of times signal has been lost",
	// 	"loss_of_sync_total":             "Number of failures on either bit or transmission word boundaries",
	// 	"rx_frames_total":                "Number of frames received",
	// 	"error_frames_total":             "Number of errors in frames",
	// 	"invalid_tx_words_total":         "Number of invalid words transmitted by host port",
	// 	"seconds_since_last_reset_total": "Number of seconds since last host port reset",
	// 	"tx_words_total":                 "Number of words transmitted by host port",
	// 	"invalid_crc_total":              "Invalid Cyclic Redundancy Check count",
	// 	"nos_total":                      "Number Not_Operational Primitive Sequence received by host port",
	// 	"fcp_packet_aborts_total":        "Number of aborted packets",
	// 	"rx_words_total":                 "Number of words received by host port",
	// 	"tx_frames_total":                "Number of frames transmitted by host port",
	// 	"link_failure_total":             "Number of times the host port link has failed",
	// 	"name":                           "Name of Fibre Channel HBA",
	// 	"speed":                          "Current operating speed",
	// 	"port_state":                     "Current port state",
	// 	"port_type":                      "Port type, what the port is connected to",
	// 	"symbolic_name":                  "Symoblic Name",
	// 	"node_name":                      "Node Name as hexadecimal string",
	// 	"port_id":                        "Port ID as string",
	// 	"port_name":                      "Port Name as hexadecimal string",
	// 	"fabric_name":                    "Fabric Name; 0 if PTP",
	// 	"dev_loss_tmo":                   "Device Loss Timeout in seconds",
	// 	"supported_classes":              "The FC classes supported",
	// 	"supported_speeds":               "The FC speeds supported",
	// }

	i.subsystem = subsystem

	return &i, nil
}

func (c *fibrechannelCollector) Update(ch chan<- *pb.CollectData) error {
	rtime := utils.Now()
	hosts, err := c.fs.FibreChannelClass()
	if err != nil {
		if os.IsNotExist(err) {
			level.Debug(c.logger).Log("msg", "fibrechannel statistics not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining FibreChannel class info: %s", err)
	}

	for _, host := range hosts {
		tags := map[string]string{
			// fc链接的host，通常是host14这样的形式
			"fc_host": host.Name,
			// fc端口速率
			"port_speed": host.Speed,
			// fc端口状态
			"port_state": host.PortState,
			// fc端口类型
			"port_type": host.PortType,
			// fc端口ID
			"port_id": host.PortID,
			// fc端口名字，对应fc端口启动器的WWPN
			"port_name": host.PortName,
			// fc端口支持速率
			"supported_speeds": host.SupportedSpeeds,
			// fc端口符号名称
			"symbolic_name": host.SymbolicName,
		}

		// 错误统计
		fields := newFibreChannelFields(host)
		ch <- &pb.CollectData{
			Time:        rtime,
			Measurement: c.subsystem,
			Tags:        tags,
			Fields:      fields,
		}

	}

	return nil
}

func newFibreChannelFields(host sysfs.FibreChannelHost) map[string]float64 {
	return map[string]float64{
		"error_frames_total":             float64(host.Counters.ErrorFrames),
		"invalid_crc_total":              float64(host.Counters.InvalidCRCCount),
		"seconds_since_last_reset_total": float64(host.Counters.SecondsSinceLastReset),
		"invalid_tx_words_total":         float64(host.Counters.InvalidTXWordCount),
		"link_failure_total":             float64(host.Counters.LinkFailureCount),
		"loss_of_sync_total":             float64(host.Counters.LossOfSyncCount),
		"loss_of_signal_total":           float64(host.Counters.LossOfSignalCount),
	}
}
