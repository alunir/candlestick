package time

import (
	"sort"
	"time"

	c "github.com/alunir/candlestick/candle"
	"github.com/shopspring/decimal"
)

type TimeChart struct {
	*c.Chart
	Resolution time.Duration
}

func NewTimeChart(resolution time.Duration, candle_num int) *TimeChart {
	return &TimeChart{
		Chart:      c.NewChart(candle_num),
		Resolution: resolution,
	}
}

func (chart *TimeChart) AddTrade(ti time.Time, val, vol float64) {
	chart.Lock()
	defer chart.Unlock()

	value := decimal.NewFromFloat(val)
	volume := decimal.NewFromFloat(vol).Abs()

	x := ti.Truncate(chart.Resolution)

	if chart.TimeSet.Contains(x.Unix()) {
		chart.CurrentCandle.AddCandleWithBuySell(c.ALL, value, volume, decimal.Zero) // MEMO: no meaning for stack
		chart.CurrentCandleNew = false
	} else {
		candle := c.NewCandleWithBuySell(c.ALL, x, value, volume, decimal.Zero) // MEMO: no meaning for stack
		chart.SetLastCandle(candle)
		if x.After(chart.LastCandle.Time.Add(chart.Resolution)) {
			chart.backfill(x, chart.LastCandle.Close)
		}
		chart.AddCandle(candle)
		chart.TimeSet.Clear()
		chart.TimeSet.Add(x.Unix())
	}
}

func (chart *TimeChart) backfill(x time.Time, value decimal.Decimal) {
	var flatCandle c.Candle
	var tmp []c.Candle

	for ti := x.Add(-chart.Resolution); !ti.Equal(chart.LastCandle.Time); ti = ti.Add(-chart.Resolution) {
		if ok := chart.TimeSet.Contains(ti.Unix()); !ok {
			flatCandle = c.NewCandle(0, ti, value, decimal.Zero, decimal.Zero)
			tmp = append(tmp, flatCandle)
			chart.TimeSet.Add(ti.Unix())
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
