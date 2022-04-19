package pubsub

import (
	"time"
)

func toUnixMilliSecond(dt time.Time) int64 {
	return dt.UnixNano() / int64(time.Millisecond)
}

func parseUnixMilliSecond(unixMs int64) time.Time {
	return time.Unix(0, unixMs*int64(time.Millisecond))
}
