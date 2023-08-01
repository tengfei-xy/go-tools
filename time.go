package tools

import (
	"math/rand"
	"os"
	"time"
)

func Delay(i int) {
	<-time.After(time.Duration(i) * time.Second)
}
func RemodifyTime(name string, modTime time.Time) {
	if name == "" {
		return
	}
	atime := time.Now()
	os.Chtimes(name, atime, modTime)
}
func RangdomRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())

	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}
func Rangdom(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func Timestamp2String(j int64) string {
	timestamp := int64(j) / 1000
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 15:04:05")
}
func Timestamp2Time(j int64) (time.Duration, time.Duration, time.Duration, time.Duration) {
	duration := time.Duration(j) * time.Second // 转化为 time.Duration 类型的值
	days := duration / (time.Hour * 24)
	hours := (duration % (time.Hour * 24)) / time.Hour
	minutes := (duration % time.Hour) / time.Minute
	seconds := duration % time.Minute

	return days, hours, minutes, seconds
}
func TimeGetTimestamp() int64 {
	return time.Now().UnixNano()
}
