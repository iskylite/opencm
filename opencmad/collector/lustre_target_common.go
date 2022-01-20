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
	"gopkg.in/yaml.v2"
)

// UnknownTargetType 未知的lustre存储目标类型，通常包括MGS、MDT、OST。
const UnknownTargetType = "UNKNOWN"

// updateTargetState targetType即为lustre存储目标类型，这里用于指定实际目录，即mdt、obdfilter
func (l *lustreCollector) updateTargetState(targetType string, ch chan<- *transport.CollectData) error {
	fps, err := filepath.Glob(procFilePath(filepath.Join("fs/lustre", targetType, "*")))
	if err != nil {
		return err
	}

	if len(fps) == 0 {
		level.Debug(l.logger).Log("msg", "not found target state file path", "target_type", targetType)
		return nil
	}

	for _, fp := range fps {
		numExports, err := utils.ReadAll(filepath.Join(fp, "num_exports"))
		if err != nil {
			return err
		}
		numExportsInt64, err := strconv.ParseInt(strings.TrimSpace(string(numExports)), 10, 64)
		if err != nil {
			return err
		}
		recoveryStatusFile, err := os.Open(filepath.Join(fp, "recovery_status"))
		if err != nil {
			return err
		}
		err = l.parseRecoveryStatus(recoveryStatusFile, func(tags map[string]string) {
			tags["fsname"] = l.parseLustreTarget(filepath.Base(fp))
			tags["type"] = l.parseTargetType(targetType)
			tags["label"] = filepath.Base(fp)
			// tags["num_exports"] = strings.TrimSpace(string(numExports))

			ch <- &transport.CollectData{
				Time:        utils.Now(),
				Measurement: "lustre_target_state",
				Tags:        tags,
				Fields: map[string]float64{
					"num_exports": float64(numExportsInt64),
				},
			}
		})
		recoveryStatusFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// mdt、obdfilter ==> MDT、OST， otherwise UnknownTargetType
func (l *lustreCollector) parseTargetType(targetType string) string {
	switch strings.ToLower(targetType) {
	case "mdt":
		return "MDT"
	case "obdfilter":
		return "OST"
	default:
		return UnknownTargetType
	}
}

func (l *lustreCollector) parseLustreTarget(fp string) string {
	if fp == "MGS" {
		return fp
	}
	return fp[0 : len(fp)-8]
}

func (l *lustreCollector) parseRecoveryStatus(reader io.Reader, handler func(map[string]string)) error {
	scanner := bufio.NewScanner(reader)

	tags := make(map[string]string, 12)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) != 2 {
			continue
		}
		key := strings.TrimRight(line[0], ":")
		val := line[1]
		tags[key] = val
	}
	handler(tags)
	return scanner.Err()
}

// lustre 存储目标卷数据采集
func (l *lustreCollector) updateTargetSize(ch chan<- *transport.CollectData) error {
	fps, err := filepath.Glob(filepath.Join(l.targetSizePath, "osd-*/*/mntdev"))
	if err != nil {
		return err
	}
	if len(fps) == 0 {
		level.Debug(l.logger).Log("msg", "not found target size file path")
		return nil
	}
	for _, fp := range fps {
		fpDir := filepath.Dir(fp)
		target := filepath.Base(fpDir)
		fsname := l.parseLustreTarget(target)
		tags := map[string]string{
			"fsname": fsname,
			"label":  target,
			"type":   l.parseTargetLabel(target),
		}
		// 存储卷类型，设备名和存储硬盘类型
		for _, tKey := range []string{"fstype", "mntdev", "nonrotational"} {
			tval, err := utils.ReadAll(filepath.Join(fpDir, tKey))
			if err != nil {
				return err
			}
			if tKey == "nonrotational" {
				if tval[0] == '0' {
					tags[tKey] = "HDD"
				} else {
					tags[tKey] = "SSD"
				}
			} else {
				tags[tKey] = strings.TrimSpace(string(tval))
			}
		}
		// 存储卷容量
		fields := make(map[string]float64, 5)
		for _, fKey := range []string{"filesfree", "filestotal", "kbytesavail", "kbytestotal", "kbytesfree"} {
			fdata, err := utils.ReadAll(filepath.Join(fpDir, fKey))
			if err != nil {
				return err
			}
			fdataint64, err := strconv.ParseInt(strings.TrimSpace(string(fdata)), 10, 64)
			if err != nil {
				return err
			}
			fields[fKey] = float64(fdataint64)
		}
		ch <- &transport.CollectData{
			Time:        utils.Now(),
			Measurement: "lustre_target_size",
			Tags:        tags,
			Fields:      fields,
		}
	}
	return nil
}

func (l *lustreCollector) parseTargetLabel(label string) string {
	if label == "MGS" {
		return label
	}
	return label[len(label)-7 : len(label)-4]
}

// lustre job_stats采集
func (l *lustreCollector) updateTargetJobStats(ch chan<- *transport.CollectData) error {
	fps, err := filepath.Glob(procFilePath("fs/lustre/*/*/job_stats"))
	if err != nil {
		return err
	}

	if len(fps) == 0 {
		level.Debug(l.logger).Log("msg", "not found target job_stats file path")
		return nil
	}

	for _, fp := range fps {
		fpDir := filepath.Dir(fp)
		target := filepath.Base(fpDir)
		fsname := l.parseLustreTarget(target)
		targetType := l.parseTargetLabel(target)

		jobStatsDataBytes, err := utils.ReadAll(fp)
		if err != nil {
			return err
		}

		jobStatsData := make(map[string][]map[string]interface{})
		err = yaml.Unmarshal(jobStatsDataBytes, &jobStatsData)
		if err != nil {
			return err
		}
		err = l.parseJobStatsData(jobStatsData["job_stats"], func(rtime int64, jobID string, fields map[string]float64) {
			ch <- &transport.CollectData{
				Time:        rtime,
				Measurement: "lustre_" + strings.ToLower(targetType) + "_job_stats",
				Tags: map[string]string{
					"fsname": fsname,
					"type":   targetType,
					"label":  target,
					"job_id": jobID,
				},
				Fields: fields,
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// 重新解析由yaml转换job_stats为map的数据，并使用handler函数处理
func (l *lustreCollector) parseJobStatsData(jobStatsData []map[string]interface{}, handler func(int64, string, map[string]float64)) error {
	for _, jobStats := range jobStatsData {
		fields := make(map[string]float64, 21)
		var rtime int
		var jobID int
		var ok bool
		for jKey, jVal := range jobStats {
			switch jKey {
			case "job_id":
				jobID, ok = jVal.(int)
				if !ok {
					level.Debug(l.logger).Log("mgs", "assert failed", "key", jKey, "val", jVal)
					return errLustreNotAvailable
				}
			case "snapshot_time":
				rtime, ok = jVal.(int)
				if !ok {
					level.Debug(l.logger).Log("mgs", "assert failed", "key", jKey, "val", jVal)
					return errLustreNotAvailable
				}
			default:
				jVal, ok := jVal.(map[interface{}]interface{})
				if !ok {
					level.Debug(l.logger).Log("mgs", "assert failed", "key", jKey, "val", jVal)
					return errLustreNotAvailable
				}
				if strings.HasSuffix(jKey, "_bytes") {
					jValSum, ok := jVal["sum"].(int)
					if !ok {
						level.Debug(l.logger).Log("mgs", "assert sum failed", "key", jKey, "val", jobStats[jKey])
						return errLustreNotAvailable
					}
					fields[jKey] = float64(jValSum)
					jKey = strings.TrimRight(jKey, "_bytes")
				}
				jValReqs, ok := jVal["samples"].(int)
				if !ok {
					level.Debug(l.logger).Log("mgs", "assert samples failed", "key", jKey, "val", jobStats[jKey])
					return errLustreNotAvailable
				}
				fields[jKey] = float64(jValReqs)
			}
		}
		handler(int64(rtime), strconv.Itoa(jobID), fields)
	}
	return nil
}

// lustre配额信息采集
func (l *lustreCollector) updateTargetQuotaSlave(ch chan<- *transport.CollectData) error {
	fps, err := filepath.Glob(procFilePath("fs/lustre/*/*/quota_slave/acct_*"))
	if err != nil {
		return err
	}

	if len(fps) == 0 {
		level.Debug(l.logger).Log("msg", "not found target quota_slave file path")
		return nil
	}
	for _, fp := range fps {
		fpList := strings.Split(fp, string(os.PathSeparator))
		target := fpList[len(fpList)-3]
		fsname := l.parseLustreTarget(target)
		targetType := l.parseTargetLabel(target)
		acctFileName := filepath.Base(fp)
		var yamlGroup string
		switch acctFileName {
		case "acct_user":
			yamlGroup = "usr_accounting"
		case "acct_group":
			yamlGroup = "grp_accounting"
		case "acct_project":
			yamlGroup = "prj_accounting"
		}
		rtime := utils.Now()

		quotaSlaveDataBytes, err := utils.ReadAll(fp)
		if err != nil {
			return err
		}
		// 如果没有配置配额，那么acct_project内容将是not supported，故skip
		if len(quotaSlaveDataBytes) <= 14 {
			continue
		}
		quotaSlaveData := make(map[string][]map[string]interface{})
		err = yaml.Unmarshal(quotaSlaveDataBytes, &quotaSlaveData)
		if err != nil {
			return err
		}
		err = l.parseQuotaSlaveData(quotaSlaveData[yamlGroup], func(id string, fields map[string]float64) {
			ch <- &transport.CollectData{
				Time:        rtime,
				Measurement: "lustre_" + acctFileName,
				Tags: map[string]string{
					"fsname": fsname,
					"type":   targetType,
					"label":  target,
					"id":     id,
				},
				Fields: fields,
			}
		})
	}
	return nil
}

func (l *lustreCollector) parseQuotaSlaveData(quotaSlaveData []map[string]interface{}, handler func(string, map[string]float64)) error {
	for _, quotaSlave := range quotaSlaveData {
		id, ok := quotaSlave["id"].(int)
		if !ok {
			level.Debug(l.logger).Log("mgs", "assert failed", "key", "id", "val", quotaSlave["id"])
			return errLustreNotAvailable
		}
		usage, ok := quotaSlave["usage"].(map[interface{}]interface{})
		if !ok {
			level.Debug(l.logger).Log("mgs", "assert failed", "key", "usage", "val", quotaSlave["usage"])
			return errLustreNotAvailable
		}
		inodes, ok := usage["inodes"].(int)
		if !ok {
			level.Debug(l.logger).Log("mgs", "assert failed", "key", "usage.inodes", "val", usage["inodes"])
			return errLustreNotAvailable
		}
		kbytes, ok := usage["kbytes"].(int)
		if !ok {
			level.Debug(l.logger).Log("mgs", "assert failed", "key", "usage.kbytes", "val", usage["kbytes"])
			return errLustreNotAvailable
		}
		handler(strconv.Itoa(id), map[string]float64{
			"inodes": float64(inodes),
			"kbytes": float64(kbytes),
		})
	}
	return nil
}
