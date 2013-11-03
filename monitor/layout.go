package monitor

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"reflect"
)

type Column struct {
	Header   string
	Width    int
	Field    string
	Format   string
	Sort     MarketSort
	Shortcut int
}

type Columns []Column

var columns = Columns{
	{"Symbol", 15, "Symbol", "%s", BySymbol, 0},
	{"Latest", 20, "Latest", "%s", ByLatest, 2},
	{"Close", 12, "Close", "%11.5f", ByClose, 4},
	{"Volume", 12, "Volume", "%11.5f", ByVolume, 0},
	{"High", 12, "High", "%11.5f", ByHigh, 0},
	{"Low", 12, "Low", "%11.5f", ByLow, 0},
	{"Ask", 12, "Ask", "%11.5f", ByAsk, 0},
	{"Bid", 12, "Bid", "%11.5f", ByBid, 0},
}

func Paint(s *State, query *Query) {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	filtered := s.Markets.Query(query)
	currencies := s.Currencies(query)
	_, height := termbox.Size()
	logY := height - s.Log.Len()
	filtered.Draw(0, 0, 100, logY-1, query)
	currencies.Draw(120, 0, 10, logY-1, query)
	inactive := "visible"
	if !query.Ancient {
		inactive = "hidden"
	}
	drawString(120, logY-1, "Inactive "+inactive, 0, termbox.ColorYellow, termbox.ColorBlack)
	drawString(0, logY-1, "Data courtesy of http://bitcoincharts.com/", -1, termbox.ColorYellow, termbox.ColorBlack)
	s.Log.Draw(0, logY)
	termbox.HideCursor()
	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}

func (q *Query) Sorted(str string, sort MarketSort) string {
	switch {
	case q.Order == sort && q.Desc:
		return str + " △"
	case q.Order == sort && !q.Desc:
		return str + " ▽"
	default:
		return str
	}
}

func (c Currencies) Draw(x, y, width, height int, query *Query) {
	drawString(x, y, query.Sorted("Currencies", ByCurrency), 0, termbox.ColorYellow|termbox.AttrBold, termbox.ColorBlack)
	for i, cur := range c[:min(height-1, len(c))] {
		bg := termbox.ColorBlack
		if cur.Name == query.Currency {
			bg = termbox.ColorGreen
		}
		drawString(x, y+i+1, fmt.Sprintf("%s (%d)", cur.Name, cur.Markets), -1, termbox.ColorWhite, bg)
	}
}

func (markets Markets) Draw(x, y, width, height int, query *Query) {
	for _, col := range columns {
		drawString(x, y, query.Sorted(col.Header, col.Sort), col.Shortcut, termbox.ColorYellow|termbox.AttrBold, termbox.ColorBlack)
		for y, m := range markets[:min(height-1, len(markets))] {
			v := reflect.ValueOf(m).FieldByName(col.Field).Interface()
			drawString(x, y+1, fmt.Sprintf(col.Format, v), -1, termbox.ColorWhite, termbox.ColorBlack)
		}
		x += col.Width
	}
}

func (l *Log) Draw(x, y int) {
	for i := 0; i < l.Len(); i++ {
		drawString(x, y+i, l.Roll(), -1, termbox.ColorWhite, termbox.ColorBlack)
	}
}

func drawString(x, y int, str string, shortcut int, fg, bg termbox.Attribute) int {
	for i, s := range str {
		if i == shortcut {
			termbox.SetCell(x, y, s, fg|termbox.AttrUnderline, bg)
		} else {
			termbox.SetCell(x, y, s, fg, bg)
		}
		x++
	}
	return x
}
