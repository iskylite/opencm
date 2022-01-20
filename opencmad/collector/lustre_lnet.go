package collector

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/transport"
)

func (l *lustreCollector) updateLnetMemUsed(ch chan<- *transport.CollectData) error {
	fp := filepath.Join(l.lnetPath, "lnet_memused")
	lnetMemUsedByte, err := utils.ReadAll(fp)
	if err != nil {
		return err
	}
	lnetMemUsedInt, err := strconv.Atoi(strings.TrimSpace(string(lnetMemUsedByte)))
	if err != nil {
		return err
	}
	ch <- &transport.CollectData{
		Time:        utils.Now(),
		Measurement: "lustre_lnet_memused",
		Fields: map[string]float64{
			"value": float64(lnetMemUsedInt),
		},
	}
	return nil
}

func (l *lustreCollector) updateLnetNis(ch chan<- *transport.CollectData) error {
	fp := filepath.Join(l.lnetPath, "nis")
	lnetNisFile, err := os.Open(fp)
	if err != nil {
		return err
	}
	err = l.parseLnetNisFile(lnetNisFile, func(nid, state string) {
		ch <- &transport.CollectData{
			Time:        utils.Now(),
			Measurement: "lustre_nis_state",
			Tags: map[string]string{
				"nid":   nid,
				"state": state},
		}
	})
	lnetNisFile.Close()
	return err
}

func (l *lustreCollector) parseLnetNisFile(reader io.Reader, handler func(string, string)) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		switch line[0] {
		case "nid":
		case "0@lo":
		default:
			handler(line[0], line[1])
		}

	}
	return scanner.Err()
}

func (l *lustreCollector) updateLnetStats(ch chan<- *transport.CollectData) error {
	fp := filepath.Join(l.lnetPath, "stats")
	// 第一次采集信息
	lnetStatsFile, err := os.Open(fp)
	if err != nil {
		return err
	}
	fTime := time.Now()
	var fData []int
	err = l.parseLnetStatsFile(lnetStatsFile, func(vals []int) {
		fData = vals
	})
	lnetStatsFile.Close()
	if err != nil {
		return err
	}
	// 等待0.5s
	time.Sleep(time.Duration(*collectInterval) * time.Millisecond)
	lTime := time.Now()
	// 第二次采集信息
	lnetStatsFile, err = os.Open(fp)
	if err != nil {
		return err
	}
	err = l.parseLnetStatsFile(lnetStatsFile, func(lData []int) {
		eTime := lTime.Sub(fTime).Seconds()
		// 结果计算
		fields := map[string]float64{
			"msgs_alloc":  float64(lData[0]),
			"msgs_max":    float64(lData[1]),
			"errors":      float64(lData[2]),
			"send_count":  float64(lData[3]),
			"recv_count":  float64(lData[4]),
			"route_count": float64(lData[5]),
			"drop_count":  float64(lData[6]),
			// *_mb_im 瞬时值
			"send_mb_im":  float64(sub(lData[7], fData[7])) / (eTime * 1024.0 * 1024.0),
			"recv_mb_im":  float64(sub(lData[8], fData[8])) / (eTime * 1024.0 * 1024.0),
			"route_mb_im": float64(sub(lData[9], fData[9])) / (eTime * 1024.0 * 1024.0),
			"drop_mb_im":  float64(sub(lData[10], fData[10])) / (eTime * 1024.0 * 1024.0),
			// *_bytes 累计值
			"send_bytes":  float64(lData[7]),
			"recv_bytes":  float64(lData[8]),
			"route_bytes": float64(lData[9]),
			"drop_bytes":  float64(lData[10]),
		}
		ch <- &transport.CollectData{
			Time:        lTime.Local().Unix(),
			Measurement: "lustre_lnet_stats",
			Tags:        nil,
			Fields:      fields,
		}
	})
	lnetStatsFile.Close()
	return err
}

func sub(a, b int) int {
	if a > b {
		return a - b
	}
	return 0
}

func subf(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return 0
}

func (l *lustreCollector) parseLnetStatsFile(reader io.Reader, handler func([]int)) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		data := make([]int, 0, 11)
		line := strings.Fields(scanner.Text())
		if len(line) != 11 {
			level.Debug(l.logger).Log("msg", "parse lnet stats failed", "data", line)
			continue
		}
		for _, val := range line {
			valStr, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			data = append(data, valStr)
		}
		handler(data)
	}
	return scanner.Err()
}
