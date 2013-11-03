package monitor

import (
	"fmt"
	"strconv"
	"time"
)

type UnixTime time.Time

func (t *UnixTime) UnmarshalJSON(b []byte) error {
	if unix, err := strconv.ParseInt(string(b), 10, 64); err != nil {
		return fmt.Errorf("Bad UnixtTime: %s", b)
	} else {
		*t = UnixTime(time.Unix(unix, 0))
	}
	return nil
}

func (t UnixTime) String() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}

func (t UnixTime) Before(o UnixTime) bool {
	return time.Time(t).Before(time.Time(o))
}

func (t UnixTime) Ancient() bool {
	return time.Now().Sub(time.Time(t)) > time.Hour*24*2
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
