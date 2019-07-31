package util

import "time"

func GetDateStr() string {
	t := time.Now()
	return t.Format("20060102150405")
}
