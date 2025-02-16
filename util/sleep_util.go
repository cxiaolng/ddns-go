package util

import (
	"log"
	"time"
)

func IsSleepMode(sleepTimeStart string, sleepTimeEnd string) bool {
	var cstZone = time.FixedZone("CST", 8*60*60)
	now := time.Now().In(cstZone)
	currentDate := now.Format("2006-01-02")

	start, err1 := time.ParseInLocation("2006-01-02 15:04", currentDate+" "+sleepTimeStart, cstZone)
	if err1 != nil {
		log.Printf("睡眠免打扰模式判断异常, ERROR: %s\n", err1)
		return false
	}
	end, err2 := time.ParseInLocation("2006-01-02 15:04", currentDate+" "+sleepTimeEnd, cstZone)
	if err2 != nil {
		log.Printf("睡眠免打扰模式判断异常, ERROR: %s\n", err2)
		return false
	}
	if start.After(end) {
		end = end.Add(24 * time.Hour)
	}
	if now.After(start) && now.Before(end) {
		return true
	}
	return false
}
