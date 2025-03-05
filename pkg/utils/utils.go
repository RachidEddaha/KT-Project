package utils

import (
	"time"
)

func TimeNowInUTC() time.Time {
	return time.Now().In(time.UTC)
}
