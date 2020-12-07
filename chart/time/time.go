package time

import (
	"math"
	"time"

	c "github.com/alunir/candlestick/candle"
)

type TimeChart struct {
	c.Chart
	Resolution time.Duration
}

func (chart *TimeChart) AddTrade(ti time.Time, value float64, volume float64) {
	volume = math.Abs(volume)

	var x = ti.Truncate(chart.Resolution)
	var candle = chart.TimeSeries[x]

	if candle != nil {
		candle.AddCandleWithBuySell(c.ALL, value, volume, 0.0) // MEMO: no meaning for stack
		chart.CurrentCandleNew = false
	} else {
		candle = c.NewCandleWithBuySell(c.ALL, x, value, volume, 0.0) // MEMO: no meaning for stack
		chart.SetLastCandle(candle)

		if chart.LastCandle != nil && x.After(chart.LastCandle.Time.Add(chart.Resolution)) {
			chart.backfill(candle.Time, chart.LastCandle.Close)
		}
		chart.AddCandle(candle)
		chart.TimeSeries[candle.Time] = candle
	}
}

func (chart *TimeChart) backfill(x time.Time, value float64) {
	var flatCandle *c.Candle

	for ti := x; !ti.Equal(chart.LastCandle.Time); ti = ti.Add(-chart.Resolution) {
		if chart.TimeSeries[x] == nil {
			flatCandle = c.NewCandle(0, x, value, 0, 0.0)
			if len(chart.Candles) < chart.CandleNum {
				chart.Candles = append(chart.Candles, flatCandle)
			} else {
				chart.Candles = append(chart.Candles[1:chart.CandleNum:chart.CandleNum], flatCandle)
			}

			chart.TimeSeries[x] = flatCandle
		}
	}
}

func (chart *TimeChart) addCandle(candle *c.Candle) {
	chart.CurrentCandle = candle
	chart.CurrentCandleNew = true

	if len(chart.Candles) < chart.CandleNum {
		chart.Candles = append(chart.Candles, candle)
	} else {
		chart.Candles = append(chart.Candles[1:chart.CandleNum:chart.CandleNum], candle)
	}

	if candle.Time.Before(chart.StartTime) {
		chart.StartTime = candle.Time
	} else if candle.Time.After(chart.EndTime) {
		chart.EndTime = candle.Time
	}
}

func (chart *TimeChart) AddLv2DataCallback(ti time.Time, askPrices []float64, askSizes []float64, bidPrices []float64, bidSizes []float64) {
}
