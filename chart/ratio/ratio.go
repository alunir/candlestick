package ratio

import (
	"time"

	c "github.com/alunir/candlestick/candle"
)

type RatioChart struct {
	c.Chart
	Threshold float64
}

func (chart *RatioChart) AddTrade(ti time.Time, value float64, volume float64) {
	candle := c.NewCandleWithBuySell(c.ALL, ti, value, volume, isPositive(value))
	if chart.LastCandle == nil {
		chart.SetLastCandle(candle)
		chart.AddCandle(candle)
	} else {
		if (chart.CurrentCandle.Stack > 0.0) && (value > chart.Threshold) {
			// same sign between the latest ratio and the current ratio
			chart.CurrentCandle.AddCandleWithBuySell(c.ALL, value, volume, isPositive(value))
		} else if (chart.CurrentCandle.Stack == 0.0) && (value < -1.0*chart.Threshold) {
			// same sign between the latest ratio and the current ratio
			chart.CurrentCandle.AddCandleWithBuySell(c.ALL, value, volume, isPositive(value))
		} else {
			// different sign between the latest ratio and the current ratio
			chart.SetLastCandle(chart.CurrentCandle)
			chart.AddCandle(candle)
		}
	}
}

func isPositive(r float64) float64 {
	if r > 0.0 {
		return 1.0
	} else {
		return 0.0
	}
}
