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

// 采集客户端的llite下面的stats信息，主要包括当前节点的读写和操作
func (l *lustreCollector) updateClientState(ch chan<- *transport.Data) error {
	// ● /proc/fs/lustre/osc/ytfs-OST0000-osc-ffff8fb1b06d7000/ost_server_uuid
	// ● /proc/fs/lustre/mdc/ytfs-MDT0000-mdc-ffff8fb1b06d7000/mds_server_uuid
	fps, err := filepath.Glob(procFilePath("fs/lustre/*c/*/*_server_uuid"))
	if err != nil {
		return err
	}
	if len(fps) == 0 {
		level.Debug(l.logger).Log("msg", "not found client state file path")
		return nil
	}
	for _, fp := range fps {
		clientStateData, err := readAll(fp)
		if err != nil {
			return err
		}
		targetUUIDAndstate := strings.Fields(string(clientStateData))
		if len(targetUUIDAndstate) != 2 {
			level.Debug(l.logger).Log("msg", "parse client state failed", "path", fp)
			return errLustreNotAvailable
		}
		state := targetUUIDAndstate[1]
		target := strings.TrimRight(targetUUIDAndstate[0], "_UUID")
		targetType := l.parseTargetLabel(target)
		fsname := l.parseLustreTarget(target)
		ch <- &transport.Data{
			Time:        utils.Now(),
			Measurement: "lustre_client_state",
			Tags: map[string]string{
				"fsname": fsname,
				"target": target,
				"type":   targetType,
				"state":  state,
			},
		}
	}
	return nil
}

func (l *lustreCollector) parseLliteToFsname(llite string) string {
	lliteSplit := strings.Split(llite, "-")
	return strings.Join(lliteSplit[0:len(lliteSplit)-1], "-")
}

type lliteStats struct {
	rtime  float64
	fsname string
	fields map[string]float64
}

func (l *lustreCollector) updateLliteStats(ch chan<- *transport.Data) error {
	fps, err := filepath.Glob(filepath.Join(l.llitePath, "*-*/stats"))
	if err != nil {
		return err
	}
	if len(fps) == 0 {
		level.Debug(l.logger).Log("msg", "not found llite stats file path")
		return nil
	}
	// 第一次数据采集
	firstLliteStats := make(map[string]*lliteStats, len(fps))
	for _, fp := range fps {
		fpDir := filepath.Dir(fp)
		llite := filepath.Base(fpDir)
		fsname := l.parseLliteToFsname(llite)

		lliteStatsFile, err := os.Open(fp)
		if err != nil {
			return err
		}
		err = l.parseLliteStatsFile(lliteStatsFile, func(rtime float64, fields map[string]float64) {
			firstLliteStats[fsname] = &lliteStats{
				rtime:  rtime,
				fsname: fsname,
				fields: fields,
			}
		})

		lliteStatsFile.Close()
	}
	time.Sleep(time.Duration(*collectInterval) * time.Millisecond)
	// 第二次数据采集
	for _, fp := range fps {
		fpDir := filepath.Dir(fp)
		llite := filepath.Base(fpDir)
		fsname := l.parseLliteToFsname(llite)

		lliteStatsFile, err := os.Open(fp)
		if err != nil {
			return err
		}
		err = l.parseLliteStatsFile(lliteStatsFile, func(rtime float64, fields map[string]float64) {
			var etime, readBytesNow, writeBytesNow float64
			fdata, ok := firstLliteStats[fsname]
			if !ok {
				etime = float64(*collectInterval / 1000)
				readBytesNow = fields["read_bytes"] / etime / 1024.0 / 1024.0
				writeBytesNow = fields["write_bytes"] / etime / 1024.0 / 1024.0
			} else {
				etime = rtime - fdata.rtime
				readBytesNow = subf(fields["read_bytes"], fdata.fields["read_bytes"]) / etime / 1024.0 / 1024.0
				writeBytesNow = subf(fields["write_bytes"], fdata.fields["write_bytes"]) / etime / 1024.0 / 1024.0
			}
			fields["read_mb_im"] = readBytesNow
			fields["write_mb_im"] = writeBytesNow
			ch <- &transport.Data{
				Time:        int64(rtime),
				Measurement: "lustre_llite_stats",
				Tags: map[string]string{
					"fsname": fsname,
				},
				Fields: fields,
			}
		})

		lliteStatsFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *lustreCollector) parseLliteStatsFile(reader io.Reader, handler func(float64, map[string]float64)) error {
	scanner := bufio.NewScanner(reader)

	parseFlag := false
	var rtime float64
	var err error
	fields := make(map[string]float64, 25)
	for scanner.Scan() {
		if !parseFlag {
			if strings.HasPrefix(scanner.Text(), "snapshot_time") {
				rtime, err = strconv.ParseFloat(strings.Fields(scanner.Text())[1], 64)
				if err != nil {
					return err
				}
				parseFlag = true
			}
			continue
		}
		line := strings.Fields(scanner.Text())
		statsKey := line[0]
		if statsKey == "read_bytes" || statsKey == "write_bytes" {
			bytesCnts, err := strconv.ParseInt(line[6], 10, 64)
			if err != nil {
				return err
			}
			fields[statsKey] = float64(bytesCnts)
			statsKey = strings.Split(statsKey, "_")[0]
		}
		reqs, err := strconv.ParseInt(line[1], 10, 64)
		if err != nil {
			return err
		}
		fields[statsKey] = float64(reqs)
	}
	handler(rtime, fields)
	return scanner.Err()
}
