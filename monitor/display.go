package monitor

import (
	"github.com/nsf/termbox-go"
	"time"
)

var sortKeys = map[rune]MarketSort{
	'c': ByCurrency,
	's': BySymbol,
	'v': ByVolume,
	't': ByLatest,
	'e': ByClose,
	'a': ByAsk,
	'b': ByBid,
	'h': ByHigh,
	'l': ByLow,
}

func initTermBox() chan *termbox.Event {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	events := make(chan *termbox.Event)
	go func() {
		for {
			ev := termbox.PollEvent()
			events <- &ev
		}
	}()
	return events
}

func Monitor() {
	defer termbox.Close()
	state := &State{
		Log: newLog(5),
	}
	query := &Query{
		Currency: "All",
		Order:    ByVolume,
	}
	state.Update()
	events := initTermBox()
	trades := initFeed()
	ticker := time.NewTicker(1 * time.Minute)
	for {
		Paint(state, query)
		select {
		case <-ticker.C:
			state.Update()
		case t := <-trades:
			state.AddTrade(t)
		case event := <-events:
			switch event.Type {
			case termbox.EventKey:
				switch event.Key {
				case termbox.KeyCtrlZ, termbox.KeyCtrlC:
					return
				case termbox.KeyArrowDown:
					query.Currency = state.Currencies(query).Next(query.Currency)
				case termbox.KeyArrowUp:
					query.Currency = state.Currencies(query).Previous(query.Currency)
				}
				switch event.Ch {
				case 'i':
					query.Ancient = !query.Ancient
				case 'q':
					return
				}
				if sort, ok := sortKeys[event.Ch]; ok {
					if query.Order == sort {
						query.Desc = !query.Desc
					} else {
						query.Order = sort
						query.Desc = false
					}
				}
			}
		}
	}
}
