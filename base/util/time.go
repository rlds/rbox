package util

import "time"

func GetDateStr() string {
	t := time.Now()
	return t.Format("20060102150405")
}

func ToDayDate() string {
	return time.Now().Format("20060102")
}

//获取时间的时分值
func GetTimeHMstr() string {
	return time.Now().Format("15_04")
}
