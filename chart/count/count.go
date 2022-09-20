package count

import (
	"math"
	"time"

	c "github.com/alunir/candlestick/candle"
	"github.com/shopspring/decimal"
)

type CountChart struct {
	*c.Chart
	Chunk   int64
	Buysell c.BuySellType
}

func (chart CountChart) AddTrade(ti time.Time, value float64, volume float64) {
	if chart.Buysell == c.ALL {
		volume = math.Abs(volume)
	}
	chart.addTradeToCountCandle(ti, decimal.NewFromFloat(value), decimal.NewFromFloat(volume))
}

func (chart CountChart) addTradeToCountCandle(ti time.Time, value decimal.Decimal, volume decimal.Decimal) {
	if chart.CurrentCandle != nil {
		if chart.CurrentCandle.Stack.IntPart() < chart.Chunk-1 {
			chart.CurrentCandle.AddCandleWithBuySell(chart.Buysell, value, volume, decimal.NewFromInt(1))
			chart.CurrentCandleNew = false
		} else {
			candle := c.NewCandleWithBuySell(chart.Buysell, ti, value, volume, decimal.Zero) // reset the counter
			chart.SetLastCandle(candle)
			chart.AddCandle(candle)
		}
	} else {
		candle := c.NewCandleWithBuySell(chart.Buysell, ti, value, volume, decimal.Zero)
		chart.AddCandle(candle)
	}
}

func (chart CountChart) AddLv2DataCallback(ti time.Time, askPrices []float64, askSizes []float64, bidPrices []float64, bidSizes []float64) {
}

func (chart CountChart) GetChartInfo() map[string]interface{} {
	return make(map[string]interface{})
}
