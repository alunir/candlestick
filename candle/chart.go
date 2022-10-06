package candle

import (
	"context"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	json "github.com/goccy/go-json"
	"github.com/tk42/victolinux/threadsafe"
)

type Chart struct {
	Candles []Candle
	// mapset.Set[time.Time] can't be compared somehow!
	TimeSet          mapset.Set[int64]
	LastCandle       *Candle
	CurrentCandle    *Candle
	CurrentCandleNew bool
	StartTime        time.Time
	EndTime          time.Time
	CandleNum        int
	in               chan Candle
	out              chan Candle
}

func NewChart(candleNum int) *Chart {
	in := make(chan Candle)
	out := make(chan Candle, candleNum)
	buffer := threadsafe.NewRingBuffer(in, out)
	go buffer.Run(context.Background())
	return &Chart{
		CandleNum: candleNum,
		Candles:   make([]Candle, 0, candleNum),
		TimeSet:   mapset.NewSet[int64](),
		in:        in,
		out:       out,
	}
}

func (chart *Chart) GetLastCandle() Candle {
	return *chart.LastCandle
}

func (chart *Chart) GetCurrentCandle() Candle {
	return *chart.CurrentCandle
}

func (chart *Chart) GetCandles() []Candle {
	return chart.Candles
}

func (chart *Chart) GetLastCandleUpdate() chan Candle {
	return chart.out
}

func (chart *Chart) GetCandleClock(ctx context.Context, interval time.Duration) chan Candle {
	ticker := time.NewTicker(interval)
	ch := make(chan Candle)
	go func(ctx context.Context) {
		var last *Candle
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if chart.LastCandle == nil {
					continue
				}
				c := chart.LastCandle
				if c == last {
					c = &Candle{
						Time:  c.Time.Add(interval),
						Open:  c.Close,
						High:  c.Close,
						Low:   c.Close,
						Close: c.Close,
					}
				}
				ch <- *c
				last = c
			}
		}
	}(ctx)
	return ch
}

func (chart *Chart) SetLastCandle(candle Candle) {
	// (candle, CurrentCandle) -> LastCandle
	// (nil, nil) -> no update
	// (nil, not nil) -> CurrentCandle
	// (not nil, nil) -> candle
	// (not nil, not nil) -> CurrentCandle
	if chart.CurrentCandle != nil {
		chart.LastCandle = chart.CurrentCandle
	} else {
		chart.LastCandle = &candle
		// if candle != nil {
		// 	chart.LastCandle = candle
		// } else {
		// no update
		// panic("SetLastCandle: (candle, CurrentCandle) -> LastCandle")
		// 	return
		// }
	}

	if !(len(chart.Candles) == 0 && chart.CurrentCandle == nil) {
		chart.appendLastCandle()
	}

	chart.in <- *chart.LastCandle
}

func (chart *Chart) appendLastCandle() {
	if len(chart.Candles) < chart.CandleNum {
		chart.Candles = append(chart.Candles, *chart.LastCandle)
	} else {
		chart.Candles = append(chart.Candles[1:chart.CandleNum:chart.CandleNum], *chart.LastCandle)
	}
}

func (chart *Chart) AddCandle(candle Candle) {
	chart.CurrentCandle = &candle
	chart.CurrentCandleNew = true

	// Need? => yes if backfill is not passed
	// if len(chart.Candles) < chart.CandleNum {
	// 	chart.Candles = append(chart.Candles, candle)
	// } else {
	// 	chart.Candles = append(chart.Candles[1:chart.CandleNum:chart.CandleNum], candle)
	// }

	if candle.Time.Before(chart.StartTime) {
		chart.StartTime = candle.Time
	} else if candle.Time.After(chart.EndTime) {
		chart.EndTime = candle.Time
	}
}

func (c Chart) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Chart) Unmarshal(b []byte) error {
	return json.Unmarshal(b, c)
}
