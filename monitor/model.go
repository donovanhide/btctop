package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type Query struct {
	Currency string
	Order    MarketSort
	Desc     bool
	Ancient  bool
}

type Trade struct {
	Symbol        string
	ID            int64
	TimeStamp     UnixTime
	Price, Volume float64
}

type Market struct {
	Symbol, Currency                              string
	Open, Close, Ask, Bid, High, Low, Avg, Volume float64
	Latest                                        UnixTime `json:"latest_trade"`
	CurrencyVolume                                float64  `json:"currency_volume"`
	Trades                                        []Trade
}

type Markets []Market

type Currency struct {
	Name    string
	Markets int
}

type Currencies []Currency

type State struct {
	Markets Markets
	Log     *Log
}

type By func(m1, m2 *Market) bool

type MarketSort uint8

const (
	BySymbol MarketSort = iota
	ByVolume
	ByCurrency
	ByLatest
	ByClose
	ByHigh
	ByLow
	ByAsk
	ByBid
)

var sorts = [...]By{
	BySymbol: func(m1, m2 *Market) bool {
		return m1.Symbol < m2.Symbol
	},
	ByVolume: func(m1, m2 *Market) bool {
		return m1.Volume > m2.Volume
	},
	ByCurrency: func(m1, m2 *Market) bool {
		return m1.Currency < m2.Currency
	},
	ByLatest: func(m1, m2 *Market) bool {
		return m2.Latest.Before(m1.Latest)
	},
	ByClose: func(m1, m2 *Market) bool {
		return m1.Close > m2.Close
	},
	ByHigh: func(m1, m2 *Market) bool {
		return m1.High > m2.High
	},
	ByLow: func(m1, m2 *Market) bool {
		return m1.Low > m2.Low
	},
	ByAsk: func(m1, m2 *Market) bool {
		return m1.Ask > m2.Ask
	},
	ByBid: func(m1, m2 *Market) bool {
		return m1.Bid > m2.Bid
	},
}

type marketSorter struct {
	markets Markets
	by      By
}

func (m *marketSorter) Len() int           { return len(m.markets) }
func (m *marketSorter) Swap(i, j int)      { m.markets[i], m.markets[j] = m.markets[j], m.markets[i] }
func (m *marketSorter) Less(i, j int) bool { return m.by(&m.markets[i], &m.markets[j]) }

func (by By) Sort(m Markets, desc bool) {
	ms := &marketSorter{
		markets: m,
		by:      by,
	}
	if desc {
		sort.Sort(sort.Reverse(ms))
	} else {
		sort.Sort(ms)
	}
}

func (c Currencies) Len() int      { return len(c) }
func (c Currencies) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c Currencies) Less(i, j int) bool {
	if c[i].Markets == c[j].Markets {
		return c[i].Name < c[j].Name
	}
	return c[i].Markets > c[j].Markets
}

func (s *State) Currencies(q *Query) (currencies Currencies) {
	counts := make(map[string]int)
	for _, market := range s.Markets {
		if !q.Ancient && market.Latest.Ancient() {
			continue
		}
		counts[market.Currency]++
	}
	currencies = make(Currencies, 0)
	total := 0
	for m, c := range counts {
		total += c
		currencies = append(currencies, Currency{
			Name:    m,
			Markets: c,
		})
	}
	sort.Sort(currencies)
	currencies = append(Currencies{{"All", total}}, currencies...)
	return
}

func (c Currencies) Next(current string) string {
	for i := range c {
		if c[i].Name == current {
			if i < len(c)-1 {
				return c[i+1].Name
			}
		}
	}
	return c[0].Name
}

func (c Currencies) Previous(current string) string {
	for i := range c {
		if c[i].Name == current {
			if i > 0 {
				return c[i-1].Name
			}
		}
	}
	return c[len(c)-1].Name
}

func (m Markets) Query(q *Query) (results Markets) {
	for _, market := range m {
		if !q.Ancient && market.Latest.Ancient() {
			continue
		}
		if q.Currency == "All" || market.Currency == q.Currency {
			results = append(results, market)
		}
	}
	sorts[q.Order].Sort(results, q.Desc)
	return
}

func (s *State) AddTrade(t *trade) {
	if t.err != nil {
		s.Log.Add(t.err.Error())
		return
	}
	var trade Trade
	if err := json.Unmarshal(t.line, &trade); err != nil {
		s.Log.Add(err.Error())
		return
	}
	i := sort.Search(len(s.Markets), func(j int) bool {
		return s.Markets[j].Symbol >= trade.Symbol
	})
	s.Markets[i].Trades = append(s.Markets[i].Trades, trade)
	s.Markets[i].Latest = trade.TimeStamp
	s.Markets[i].Close = trade.Price
	s.Log.Add(fmt.Sprintf("Trade updated: %s %.4f %.4f", trade.Symbol, trade.Price, trade.Volume))
	return
}

func (s *State) Update() {
	resp, err := http.Get("http://api.bitcoincharts.com/v1/markets.json")
	if err != nil {
		s.Log.Add(err.Error())
		return
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&s.Markets); err != nil {
		s.Log.Add(err.Error())
		return
	}
	sorts[BySymbol].Sort(s.Markets, false)
	s.Log.Add("Market updated")
	return
}
