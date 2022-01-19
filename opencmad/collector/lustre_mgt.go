package collector

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/transport"
)

func (l *lustreCollector) updateMGTState(ch chan<- *transport.Data) error {
	if _, err := l.openProcFile(procFilePath("fs/lustre/mgs/MGS")); err != nil {
		if err == errLustreNotAvailable {
			level.Debug(l.logger).Log("err", err)
			return ErrNoData
		}
	}
	numExports, err := readAll(procFilePath("fs/lustre/mgs/MGS/num_exports"))
	if err != nil {
		return err
	}
	numExportsInt64, err := strconv.ParseInt(strings.TrimSpace(string(numExports)), 10, 64)
	if err != nil {
		return err
	}
	fps, err := filepath.Glob(procFilePath("fs/lustre/mgs/MGS/live/*"))
	if err != nil {
		return err
	}

	if len(fps) == 0 {
		level.Debug(l.logger).Log("msg", "not found mgs live/* file path")
		return nil
	}
	rtime := utils.Now()
	for _, fspath := range fps {
		if filepath.Base(fspath) == "params" {
			continue
		}
		file, err := os.Open(fspath)
		if err != nil {
			// This file should exist, but there is a race where an exporting pool can remove the files. Ok to ignore.
			level.Debug(l.logger).Log("msg", "Cannot open file for reading", "path", fspath)
			return errLustreNotAvailable
		}

		err = l.parseMGSLiveFile(file, func(fsname string, state string) {
			ch <- &transport.Data{
				Time:        rtime,
				Measurement: "lustre_mgt_state",
				Tags: map[string]string{
					"fsname": fsname,
					"state":  state,
					// "num_exports": strings.TrimSpace(string(numExports)),
				},
				Fields: map[string]float64{
					"num_exports": float64(numExportsInt64),
				},
			}
		})

		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *lustreCollector) parseMGSLiveFile(reader io.Reader, handler func(string, string)) error {
	scanner := bufio.NewScanner(reader)

	var fsname, state string
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case "fsname:":
			fsname = line[1]
		case "state:":
			state = line[1]
			break
		}

	}
	handler(fsname, state)
	return scanner.Err()
}
