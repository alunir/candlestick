package count

import (
	"math"
	"time"

	c "github.com/alunir/candlestick/candle"
)

type CountChart struct {
	c.Chart
	Chunk   int64
	Buysell c.BuySellType
}

func (chart *CountChart) AddTrade(ti time.Time, value float64, volume float64) {
	if chart.Buysell == c.ALL {
		volume = math.Abs(volume)
	}
	chart.addTradeToCountCandle(ti, value, volume)
}

func (chart *CountChart) addTradeToCountCandle(ti time.Time, value float64, volume float64) {
	if chart.LastCandle != nil {
		if int64(chart.LastCandle.Stack) < chart.Chunk-1 {
			chart.LastCandle.AddCandleWithBuySell(chart.Buysell, value, volume, +1.0)
			chart.CurrentCandleNew = false
		} else {
			candle := c.NewCandleWithBuySell(chart.Buysell, ti, value, volume, 0.0) // reset the counter
			chart.SetLastCandle(candle)
			chart.AddCandle(candle)
		}
	} else {
		candle := c.NewCandleWithBuySell(chart.Buysell, ti, value, volume, 0.0)
		chart.SetLastCandle(candle)
		chart.AddCandle(candle)
	}
}
