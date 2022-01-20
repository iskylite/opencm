package collector

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-systemd/dbus"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/transport"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const systemdFlag = 1 << 9

var (
	unitIncludeSet bool
	unitInclude    = kingpin.Flag("collector.systemd.unit-include", "Regexp of systemd units to include. Units must both match include and not match exclude to be included.").Default("ssh|chrony|multipathd|xinetd|zfs-zed").PreAction(func(c *kingpin.ParseContext) error {
		unitIncludeSet = true
		return nil
	}).String()
	unitExcludeSet bool
	unitExclude    = kingpin.Flag("collector.systemd.unit-exclude", "Regexp of systemd units to exclude. Units must both match include and not match exclude to be included.").Default(".+\\.(automount|device|mount|scope|slice)").PreAction(func(c *kingpin.ParseContext) error {
		unitExcludeSet = true
		return nil
	}).String()
)

type systemdCollector struct {
	unitIncludePattern *regexp.Regexp
	unitExcludePattern *regexp.Regexp
	logger             log.Logger
	subsystem          string
}

var unitStatesName = []string{"active", "activating", "deactivating", "inactive", "failed"}

func init() {
	registerCollector("systemd", systemdFlag, NewSystemdCollector)
}

// NewSystemdCollector returns a new Collector exposing systemd statistics.
func NewSystemdCollector(logger log.Logger, subsystem string) (Collector, error) {
	level.Info(logger).Log("msg", "Parsed flag --collector.systemd.unit-include", "flag", *unitInclude)
	unitIncludePattern := regexp.MustCompile(fmt.Sprintf("^(%s).service$", *unitInclude))
	level.Info(logger).Log("msg", "Parsed flag --collector.systemd.unit-exclude", "flag", *unitExclude)
	unitExcludePattern := regexp.MustCompile(fmt.Sprintf("^(?:%s)$", *unitExclude))

	return &systemdCollector{
		unitIncludePattern: unitIncludePattern,
		unitExcludePattern: unitExcludePattern,
		logger:             logger,
		subsystem:          subsystem,
	}, nil
}

// Update gathers metrics from systemd.  Dbus collection is done in parallel
// to reduce wait time for responses.
func (c *systemdCollector) Update(ch chan<- *transport.CollectData) error {
	begin := time.Now()
	rtime := utils.Now()
	conn, err := newSystemdDbusConn()
	if err != nil {
		return fmt.Errorf("couldn't get dbus connection: %w", err)
	}
	defer conn.Close()

	allUnits, err := c.getAllUnits(conn)
	if err != nil {
		return fmt.Errorf("couldn't get units: %w", err)
	}
	level.Debug(c.logger).Log("msg", "getAllUnits took", "duration_seconds", time.Since(begin).Seconds())

	begin = time.Now()
	units := filterUnits(allUnits, c.unitIncludePattern, c.unitExcludePattern, c.logger)
	level.Debug(c.logger).Log("msg", "filterUnits took", "duration_seconds", time.Since(begin).Seconds())

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		begin = time.Now()
		c.collectUnitStatusMetrics(conn, ch, units, rtime)
		level.Debug(c.logger).Log("msg", "collectUnitStatusMetrics took", "duration_seconds", time.Since(begin).Seconds())
	}()
	return err
}

func (c *systemdCollector) collectUnitStatusMetrics(conn *dbus.Conn, ch chan<- *transport.CollectData, units []unit, rtime int64) {
	for _, unit := range units {
		serviceType := ""
		if strings.HasSuffix(unit.Name, ".service") {
			serviceTypeProperty, err := conn.GetUnitTypeProperty(unit.Name, "Service", "Type")
			if err != nil {
				level.Debug(c.logger).Log("msg", "couldn't get unit type", "unit", unit.Name, "err", err)
			} else {
				serviceType = serviceTypeProperty.Value.Value().(string)
			}
		}
		ch <- &transport.CollectData{
			Time:        rtime,
			Measurement: c.subsystem,
			Tags:        map[string]string{"name": unit.Name, "state": unit.ActiveState, "service_type": serviceType},
		}
	}
}

func newSystemdDbusConn() (*dbus.Conn, error) {
	return dbus.New()
}

type unit struct {
	dbus.UnitStatus
}

func (c *systemdCollector) getAllUnits(conn *dbus.Conn) ([]unit, error) {
	allUnits, err := conn.ListUnits()
	if err != nil {
		return nil, err
	}

	result := make([]unit, 0, len(allUnits))
	for _, status := range allUnits {
		unit := unit{
			UnitStatus: status,
		}
		result = append(result, unit)
	}

	return result, nil
}

func filterUnits(units []unit, includePattern, excludePattern *regexp.Regexp, logger log.Logger) []unit {
	filtered := make([]unit, 0, len(units))
	for _, unit := range units {
		if includePattern.MatchString(unit.Name) && !excludePattern.MatchString(unit.Name) && unit.LoadState == "loaded" {
			level.Debug(logger).Log("msg", "Adding unit", "unit", unit.Name)
			filtered = append(filtered, unit)
		} else {
			level.Debug(logger).Log("msg", "Ignoring unit", "unit", unit.Name)
		}
	}

	return filtered
}
