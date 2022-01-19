package main

import (
	_ "net/http/pprof"
	"os"
	"os/user"

	"github.com/iskylite/opencm/opencmad/collector"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"

	"github.com/go-kit/log/level"
	"github.com/prometheus/common/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	// var (
	// 	listenAddress = kingpin.Flag(
	// 		"web.listen-address",
	// 		"Address on which to expose metrics and web interface.",
	// 	).Default(":9100").String()
	// 	metricsPath = kingpin.Flag(
	// 		"web.telemetry-path",
	// 		"Path under which to expose metrics.",
	// 	).Default("/metrics").String()
	// 	disableExporterMetrics = kingpin.Flag(
	// 		"web.disable-exporter-metrics",
	// 		"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
	// 	).Bool()
	// 	maxRequests = kingpin.Flag(
	// 		"web.max-requests",
	// 		"Maximum number of parallel scrape requests. Use 0 to disable.",
	// 	).Default("40").Int()
	// 	configFile = kingpin.Flag(
	// 		"web.config",
	// 		"[EXPERIMENTAL] Path to config yaml file that can enable TLS or authentication.",
	// 	).Default("").String()
	// )

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("node_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting node_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())
	if user, err := user.Current(); err == nil && user.Uid == "0" {
		level.Warn(logger).Log("msg", "Node Exporter is running as root user. This exporter is designed to run as unpriviledged user, root is not required.")
	}
	flags := 1 << 11
	nc, err := collector.NewNodeCollector(logger, uint(flags))
	if err != nil {
		level.Error(logger).Log("err", err)
	}
	datas := nc.Gather()
	for _, data := range datas {
		collector.FormatCollectorData(data, logger)
	}
}
