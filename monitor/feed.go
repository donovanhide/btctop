package monitor

import (
	"net/textproto"
)

type trade struct {
	line []byte
	err  error
}

func initFeed() chan *trade {
	trades := make(chan *trade)
	conn, err := textproto.Dial("tcp", "bitcoincharts.com:27007")
	if err != nil {
		panic(err)
	}
	if _, err = conn.Cmd(string([]byte{255, 246})); err != nil {
		panic(err)
	}
	go func() {
		defer conn.Close()
		for {
			line, err := conn.ReadLineBytes()
			trades <- &trade{line, err}
		}
	}()
	return trades
}
