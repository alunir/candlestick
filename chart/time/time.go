package time

import (
	"sort"
	"time"

	c "github.com/alunir/candlestick/candle"
	"github.com/shopspring/decimal"
	"github.com/tk42/victolinux/threadsafe"
)

type TimeChart struct {
	*c.Chart
	Resolution time.Duration
}

func (chart *TimeChart) AddTrade(ti time.Time, val float64, vol float64) {
	value := decimal.NewFromFloat(val)
	volume := decimal.NewFromFloat(vol).Abs()

	var x = ti.Truncate(chart.Resolution)
	candle, ok := chart.TimeSeries.Load(x)

	if ok {
		candle.AddCandleWithBuySell(c.ALL, value, volume, decimal.Zero) // MEMO: no meaning for stack
		chart.CurrentCandleNew = false
	} else {
		candle = c.NewCandleWithBuySell(c.ALL, x, value, volume, decimal.Zero) // MEMO: no meaning for stack
		chart.SetLastCandle(candle)
		if chart.LastCandle != nil && x.After(chart.LastCandle.Time.Add(chart.Resolution)) {
			chart.backfill(x, chart.LastCandle.Close)
		}
		chart.AddCandle(candle)
		chart.TimeSeries = threadsafe.ThreadsafeMap[time.Time, *c.Candle]{}
		chart.TimeSeries.Store(x, candle)
	}
}

func (chart *TimeChart) backfill(x time.Time, value decimal.Decimal) {
	var flatCandle *c.Candle
	var tmp []*c.Candle

	for ti := x.Add(-chart.Resolution); !ti.Equal(chart.LastCandle.Time); ti = ti.Add(-chart.Resolution) {
		if _, ok := chart.TimeSeries.Load(ti); !ok {
			flatCandle = c.NewCandle(0, ti, value, decimal.Zero, decimal.Zero)
			tmp = append(tmp, flatCandle)
			chart.TimeSeries.Store(ti, flatCandle)
		}
	}
	ReverseSlice(tmp)
	chart.Candles = append(chart.Candles, tmp...)[Max(len(chart.Candles)+len(tmp)-chart.CandleNum, 0):]
}

func (chart *TimeChart) AddLv2DataCallback(ti time.Time, askPrices []float64, askSizes []float64, bidPrices []float64, bidSizes []float64) {
}

func (chart *TimeChart) GetChartInfo() map[string]interface{} {
	return make(map[string]interface{})
}

func ReverseSlice[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
