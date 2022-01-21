package collector

import (
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/pb"
)

// NewCollectorData 生成Data
func NewCollectorData(rtime int64, measurement string, tags map[string]string, fields map[string]float64) *pb.CollectData {
	return &pb.CollectData{
		Time:        rtime,
		Measurement: measurement,
		Tags:        tags,
		Fields:      fields,
	}
}

// FormatCollectorData 格式化Data，输出字符串
func FormatCollectorData(data *pb.CollectData, logger log.Logger) {
	rtime := strconv.Itoa(int(data.Time))
	var tags, fields []string
	for tk, tv := range data.Tags {
		tags = append(tags, tk+"="+tv)
	}
	for fk, fv := range data.Fields {
		v := strconv.FormatFloat(fv, 'f', 2, 64)
		fields = append(fields, fk+"="+v)
	}
	level.Info(logger).Log("rtime", rtime, "measurement", data.Measurement, "tags", strings.Join(tags, ","), "fields", strings.Join(fields, ","))
}
