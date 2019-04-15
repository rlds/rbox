package def

// TimeCounterResult 时间计数器结果
type TimeCounterResult struct {
	AllCount    uint64 // 总计
	YearCount   uint64 // 年总计数
	MonthCount  uint64 // 月总计数
	DayCount    uint64 // 天总计数
	HourCount   uint64 // 小时总计数
	MinuteCount uint64 // 分总计数
	SecondCount uint64 // 秒总计数
}

// TimeCountRes 时间计数期结果
type TimeCountRes struct {
	Time  string
	Count uint64
}
