package collector

import (
	"errors"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/transport"
)

const zfsFlag = 1 << 10

var errZFSNotAvailable = errors.New("ZFS / ZFS statistics are not available")

type zfsSysctl string

func init() {
	registerCollector("zfs", zfsFlag, NewZFSCollector)
}

type zfsCollector struct {
	linuxProcpathBase    string
	linuxZpoolIoPath     string
	linuxZpoolObjsetPath string
	linuxZpoolStatePath  string
	subsystem            string
	linuxPathMap         map[string]string
	logger               log.Logger
}

// NewZFSCollector returns a new Collector exposing ZFS statistics.
func NewZFSCollector(logger log.Logger, subsystem string) (Collector, error) {
	return &zfsCollector{
		linuxProcpathBase:    "spl/kstat/zfs",
		linuxZpoolIoPath:     "/*/io",
		linuxZpoolObjsetPath: "/*/objset-*",
		linuxZpoolStatePath:  "/*/state",
		subsystem:            subsystem,
		linuxPathMap: map[string]string{
			"zfs_abd":         "abdstats",
			"zfs_arc":         "arcstats",
			"zfs_dbuf":        "dbuf_stats",
			"zfs_dmu_tx":      "dmu_tx",
			"zfs_dnode":       "dnodestats",
			"zfs_fm":          "fm",
			"zfs_vdev_cache":  "vdev_cache_stats", // vdev_cache is deprecated
			"zfs_vdev_mirror": "vdev_mirror_stats",
			"zfs_xuio":        "xuio_stats", // no known consumers of the XUIO interface on Linux exist
			"zfs_zfetch":      "zfetchstats",
			"zfs_zil":         "zil",
		},
		logger: logger,
	}, nil
}

func (c *zfsCollector) Update(ch chan<- *transport.Data) error {

	if _, err := c.openProcFile(c.linuxProcpathBase); err != nil {
		if err == errZFSNotAvailable {
			level.Debug(c.logger).Log("err", err)
			return ErrNoData
		}
	}

	// for subsystem := range c.linuxPathMap {
	// 	if err := c.updateZfsStats(subsystem, ch); err != nil {
	// 		if err == errZFSNotAvailable {
	// 			level.Debug(c.logger).Log("err", err)
	// 			// ZFS /proc files are added as new features to ZFS arrive, it is ok to continue
	// 			continue
	// 		}
	// 		return err
	// 	}
	// }

	// Pool stats
	return c.updatePoolStats(ch)
}

// func (c *zfsCollector) constSysctlMetric(subsystem string, sysctl zfsSysctl, value uint64) *transport.Data {
// 	metricName := sysctl.metricName()

// 	return prometheus.MustNewConstMetric(
// 		prometheus.NewDesc(
// 			prometheus.BuildFQName(namespace, subsystem, metricName),
// 			string(sysctl),
// 			nil,
// 			nil,
// 		),
// 		prometheus.UntypedValue,
// 		float64(value),
// 	)
// }

func (c *zfsCollector) constPoolMetric(rtime int64, poolName string, io map[string]float64) *transport.Data {
	return &transport.Data{
		Time:        rtime,
		Measurement: c.subsystem + "_io",
		Tags: map[string]string{
			"pool": poolName,
		},
		Fields: io,
	}
}

// func (c *zfsCollector) constPoolObjsetMetric(poolName string, datasetName string, sysctl zfsSysctl, value uint64) *transport.Data {
// 	metricName := sysctl.metricName()

// 	return prometheus.MustNewConstMetric(
// 		prometheus.NewDesc(
// 			prometheus.BuildFQName(namespace, "zfs_zpool_dataset", metricName),
// 			string(sysctl),
// 			[]string{"zpool", "dataset"},
// 			nil,
// 		),
// 		prometheus.UntypedValue,
// 		float64(value),
// 		poolName,
// 		datasetName,
// 	)
// }

func (c *zfsCollector) constPoolStateMetric(rtime int64, poolName string, stateName string) *transport.Data {
	return &transport.Data{
		Time:        rtime,
		Measurement: c.subsystem + "_state",
		Tags: map[string]string{
			"pool":  poolName,
			"state": stateName,
		},
	}
}
