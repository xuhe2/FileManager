package call

import (
	"time"
)

// GetTimeFromStamp 获取时间戳对应的时间对象
func GetTimeFromStamp(stamp string) (time.Time, error) {
	currentYear := time.Now().Format("2006")

	// 时间格式
	format := "200601021504.05" // CCYYMMDDhhmm.ss

	// 缺秒时补秒
	if len(stamp)%2 == 0 {
		stamp += ".00"
	}

	// 缺世纪年
	if len(stamp) == 11 {
		stamp = currentYear + stamp
	}

	// 缺世纪
	if len(stamp) == 13 {
		stamp = currentYear[:1] + stamp
	}

	// 解析时间字符串
	return time.Parse(format, stamp)
}
