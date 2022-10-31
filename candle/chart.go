package candle

import (
	"context"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	json "github.com/goccy/go-json"
	"github.com/tk42/victolinux/threadsafe"
)

type Chart struct {
	m       *sync.RWMutex
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
		m:         new(sync.RWMutex),
		CandleNum: candleNum,
		Candles:   make([]Candle, 0, candleNum),
		TimeSet:   mapset.NewSet[int64](),
		in:        in,
		out:       out,
	}
}

func (chart *Chart) Lock() {
	chart.m.Lock()
}

func (chart *Chart) Unlock() {
	chart.m.Unlock()
}

func (chart *Chart) GetLastCandle() (Candle, bool) {
	chart.m.RLock()
	defer chart.m.RUnlock()
	if chart.LastCandle != nil {
		return *chart.LastCandle, true
	} else {
		return Candle{}, false
	}
}

func (chart *Chart) GetCurrentCandle() (Candle, bool) {
	if chart.CurrentCandle != nil {
		return *chart.CurrentCandle, true
	} else {
		return Candle{}, false
	}
}

func (chart *Chart) GetCandles() []Candle {
	chart.m.RLock()
	defer chart.m.RUnlock()
	return chart.Candles
}

func (chart *Chart) GetLastCandleUpdate() chan Candle {
	return chart.out
}

func (chart *Chart) EqualCondition(candle, last Candle) {
	chart.LastCandle = &candle
}

func (chart *Chart) GetCandleClock(ctx context.Context, interval time.Duration) chan Candle {
	ch := make(chan Candle)
	go func(ctx context.Context, interval time.Duration) {
		var last Candle
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c, ok := chart.GetLastCandle()
				if !ok {
					continue
				}
				// MEMO: MarketData is arrived with the delay.
				// ~ So we can't handle the case that the candle is not updated.
				// Even if the interval was appended to Time,
				// the created candle would not be known how much the market data will be delayed.
				// Therefore we create a candle whenerver the market data is updated.

				// TODO: Fix when AMOUNT/VOLUME clock is implemented
				if c.OpenTime.Equal(last.OpenTime) {
					continue
				}
				for _, p := range chart.GetCandles() {
					if p.OpenTime.After(last.OpenTime) && !p.OpenTime.After(c.OpenTime) {
						ch <- p
					}
				}

				last = c.Copy()
			}
		}
	}(ctx, interval)
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

	if candle.OpenTime.Before(chart.StartTime) {
		chart.StartTime = candle.OpenTime
	} else if candle.OpenTime.After(chart.EndTime) {
		chart.EndTime = candle.OpenTime
	}
}

func (c *Chart) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Chart) Unmarshal(b []byte) error {
	return json.Unmarshal(b, c)
}
