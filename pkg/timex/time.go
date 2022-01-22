package timex

import "time"

// Now 当前本地时间的时间戳
func Now() int64 {
	return time.Now().Local().Unix()
}
