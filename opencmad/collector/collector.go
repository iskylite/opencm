package collector

import (
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/transport"
)

// Namespace defines the common namespace to be used by all metrics.
const namespace = "node"

const (
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	factories              = make(map[string]func(logger log.Logger, subsystem string) (Collector, error))
	initiatedCollectorsMtx = sync.Mutex{}
	initiatedCollectors    = make(map[string]Collector)
	// CollectorFlag 每个采集模块以及特征值
	CollectorFlag = make(map[uint]string)
)

func registerCollector(collector string, flag uint, factory func(logger log.Logger, subsystem string) (Collector, error)) {
	if _, ok := CollectorFlag[flag]; ok {
		panic(ErrFlagExist)
	}
	CollectorFlag[flag] = collector
	factories[collector] = factory
}

// NodeCollector implements the prometheus.Collector interface.
type NodeCollector struct {
	Collectors map[string]Collector
	logger     log.Logger
}

// NewNodeCollector creates a new NodeCollector.
func NewNodeCollector(logger log.Logger, flags uint) (*NodeCollector, error) {
	collectors := make(map[string]Collector)
	initiatedCollectorsMtx.Lock()
	defer initiatedCollectorsMtx.Unlock()
	for flag, key := range CollectorFlag {
		// 通过二进制判断，如果相与为0，那么这个collector将不会初始化
		if flag&flags == 0 {
			continue
		}
		if collector, ok := initiatedCollectors[key]; ok {
			collectors[key] = collector
		} else {
			collector, err := factories[key](log.With(logger, "collector", key), key)
			if err != nil {
				return nil, err
			}
			collectors[key] = collector
			initiatedCollectors[key] = collector
		}
		level.Debug(logger).Log("msg", "load collector", "collector", key)
	}
	return &NodeCollector{Collectors: collectors, logger: logger}, nil
}

// Collect implements the prometheus.Collector interface.
func (n NodeCollector) Collect(ch chan<- *transport.Data) {
	defer close(ch)
	wg := sync.WaitGroup{}
	wg.Add(len(n.Collectors))
	for name, c := range n.Collectors {
		go func(name string, c Collector) {
			execute(name, c, ch, n.logger)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

func execute(name string, c Collector, ch chan<- *transport.Data, logger log.Logger) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)

	if err != nil {
		if IsNoDataError(err) {
			level.Debug(logger).Log("msg", "collector returned no data", "name", name, "duration_seconds", duration.Seconds(), "err", err)
		} else {
			level.Error(logger).Log("msg", "collector failed", "name", name, "duration_seconds", duration.Seconds(), "err", err)
		}
	} else {
		level.Debug(logger).Log("msg", "collector succeeded", "name", name, "duration_seconds", duration.Seconds())
	}
}

// Gather 获取数据，生成响应
func (n NodeCollector) Gather() []*transport.Data {
	ch := make(chan *transport.Data, runtime.NumCPU())
	datas := make([]*transport.Data, 0, len(n.Collectors))
	go n.Collect(ch)
	for data := range ch {
		datas = append(datas, data)
	}
	return datas
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- *transport.Data) error
}

// ErrNoData indicates the collector found no data to collect, but had no other error.
var ErrNoData = errors.New("collector returned no data")

// ErrFlagExist indicates the collector registerCollector failed, flag repeate
var ErrFlagExist = errors.New("collector flag already exist")

// IsNoDataError tell error is ErrNodata
func IsNoDataError(err error) bool {
	return err == ErrNoData
}
