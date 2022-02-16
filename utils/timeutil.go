package utils

import (
	"fmt"
	"time"
)

func GetNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetNowTimeStamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
