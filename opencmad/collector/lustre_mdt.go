package collector

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/pb"
)

func (l *lustreCollector) updateMDTState(ch chan<- *pb.CollectData) error {
	return l.updateTargetState("mdt", ch)
}

func (l *lustreCollector) updateMDTStats(ch chan<- *pb.CollectData) error {
	fps, err := filepath.Glob(procFilePath("fs/lustre/mdt/*/md_stats"))
	if err != nil {
		return err
	}

	if len(fps) == 0 {
		level.Debug(l.logger).Log("msg", "not found mdt md_stats file path")
		return nil
	}
	for _, fp := range fps {
		fpDir := filepath.Dir(fp)
		target := filepath.Base(fpDir)
		fsname := l.parseLustreTarget(target)

		mdtStatsFile, err := os.Open(fp)
		if err != nil {
			return err
		}
		err = l.parseMDTStatsFile(mdtStatsFile, func(rtime int64, fields map[string]float64) {
			ch <- &pb.CollectData{
				Time:        rtime,
				Measurement: "lustre_mdt_stats",
				Tags: map[string]string{
					"fsname": fsname,
					"label":  target,
				},
				Fields: fields,
			}
		})

		mdtStatsFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *lustreCollector) parseMDTStatsFile(reader io.Reader, handler func(int64, map[string]float64)) error {
	scanner := bufio.NewScanner(reader)

	parseFlag := false
	var rtime float64
	var err error
	fields := make(map[string]float64, 14)
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
		reqs, err := strconv.ParseInt(line[1], 10, 64)
		if err != nil {
			return err
		}
		fields[line[0]] = float64(reqs)
	}
	handler(int64(rtime), fields)
	return scanner.Err()
}
