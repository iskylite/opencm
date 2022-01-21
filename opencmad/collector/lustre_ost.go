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
	"github.com/iskylite/opencm/pb"
)

func (l *lustreCollector) updateOSTState(ch chan<- *pb.CollectData) error {
	return l.updateTargetState("obdfilter", ch)
}

func (l *lustreCollector) updateOSTBRWStats(ch chan<- *pb.CollectData) error {
	fps, err := filepath.Glob(filepath.Join(l.targetSizePath, "osd-*/*OST*/brw_stats"))
	if err != nil {
		return err
	}

	if len(fps) == 0 {
		level.Debug(l.logger).Log("msg", "not found ost brw_stats file path")
		return nil
	}
	for _, fp := range fps {
		fpDir := filepath.Dir(fp)
		target := filepath.Base(fpDir)
		fsname := l.parseLustreTarget(target)

		brwStatsFile, err := os.Open(fp)
		if err != nil {
			return err
		}
		err = l.parseBRWStatsFile(brwStatsFile, func(rtime int64, ioSize string, readIOS, writeIOS float64) {
			ch <- &pb.CollectData{
				Time:        rtime,
				Measurement: "lustre_ost_brw_stats",
				Tags: map[string]string{
					"fsname":  fsname,
					"label":   target,
					"io_size": ioSize,
				},
				Fields: map[string]float64{
					"read_ios":  readIOS,
					"write_ios": writeIOS,
				},
			}
		})
		brwStatsFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *lustreCollector) parseBRWStatsFile(reader io.Reader, handler func(int64, string, float64, float64)) error {
	scanner := bufio.NewScanner(reader)

	parseFlag := false
	var rtime float64
	var err error
	for scanner.Scan() {
		if !parseFlag {
			if strings.HasPrefix(scanner.Text(), "snapshot_time") {
				rtime, err = strconv.ParseFloat(strings.Fields(scanner.Text())[1], 64)
				if err != nil {
					return err
				}
				continue
			}

			if strings.HasPrefix(scanner.Text(), "disk I/O size") {
				parseFlag = true
			}
			continue
		}
		line := strings.Fields(scanner.Text())
		ioSize := strings.TrimRight(line[0], ":")
		readIOS, err := strconv.ParseInt(line[2], 10, 64)
		if err != nil {
			return err
		}
		writeIOS, err := strconv.ParseInt(line[5], 10, 64)
		if err != nil {
			return err
		}
		handler(int64(rtime), ioSize, float64(readIOS), float64(writeIOS))
	}
	return scanner.Err()
}

type ostStats struct {
	rtime  float64
	fsname string
	label  string
	fields map[string]float64
}

func (l *lustreCollector) updateOSTStats(ch chan<- *pb.CollectData) error {
	fps, err := filepath.Glob(procFilePath("fs/lustre/obdfilter/*/stats"))
	if err != nil {
		return err
	}
	// 第一次数据采集
	firstOstStats := make(map[string]*ostStats, len(fps))
	for _, fp := range fps {
		fpDir := filepath.Dir(fp)
		target := filepath.Base(fpDir)
		fsname := l.parseLustreTarget(target)

		ostStatsFile, err := os.Open(fp)
		if err != nil {
			return err
		}
		err = l.parseOSTStatsFile(ostStatsFile, func(rtime float64, fields map[string]float64) {
			firstOstStats[target] = &ostStats{
				rtime:  rtime,
				fsname: fsname,
				label:  target,
				fields: fields,
			}
		})

		ostStatsFile.Close()
		if err != nil {
			return err
		}
	}
	time.Sleep(time.Duration(*collectInterval) * time.Millisecond)
	// 第二次数据采集
	for _, fp := range fps {
		fpDir := filepath.Dir(fp)
		target := filepath.Base(fpDir)
		fsname := l.parseLustreTarget(target)

		ostStatsFile, err := os.Open(fp)
		if err != nil {
			return err
		}
		err = l.parseOSTStatsFile(ostStatsFile, func(rtime float64, fields map[string]float64) {
			var etime, readBytesNow, writeBytesNow float64
			fdata, ok := firstOstStats[target]
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
			ch <- &pb.CollectData{
				Time:        int64(rtime),
				Measurement: "lustre_ost_stats",
				Tags: map[string]string{
					"fsname": fsname,
					"label":  target,
				},
				Fields: fields,
			}
		})

		ostStatsFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *lustreCollector) parseOSTStatsFile(reader io.Reader, handler func(float64, map[string]float64)) error {
	scanner := bufio.NewScanner(reader)

	parseFlag := false
	var rtime float64
	var err error
	fields := make(map[string]float64, 11)
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
