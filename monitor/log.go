package monitor

import (
	"container/ring"
	"fmt"
	"time"
)

type Log struct {
	ring *ring.Ring
}

func newLog(n int) *Log {
	return &Log{ring: ring.New(5)}
}

func (l *Log) Add(msg string) {
	l.ring.Value = fmt.Sprintf("%s: %s", time.Now().String()[:23], msg)
	l.ring = l.ring.Move(1)
}
func (l *Log) Len() int {
	return l.ring.Len()
}

func (l *Log) Roll() string {
	msg, ok := l.ring.Value.(string)
	if !ok {
		msg = ""
	}
	l.ring = l.ring.Move(1)
	return msg
}
