package ratio

import (
	"math"
	"time"

	c "github.com/alunir/candlestick/candle"
)

type RatioChart struct {
	c.Chart
	Thresholds []float64
	LastSign   bool
}

// Only update OHLC, Count and Volume to current candle
func (chart *RatioChart) AddTrade(ti time.Time, price float64, volume float64) {
	candle := c.NewCandleWithBuySell(c.ALL, ti, price, volume, 0.0)
	if chart.CurrentCandle == nil {
		chart.AddCandle(candle)
	} else {
		chart.CurrentCandle.AddCandleWithBuySell(c.ALL, price, volume, 0.0)
	}
}

// Switch to a new candle
func (chart *RatioChart) AddLv2DataCallback(ti time.Time, askPrices []float64, askSizes []float64, bidPrices []float64, bidSizes []float64) {
	askTotalSizes, bidTotalSizes := make([]float64, len(chart.Thresholds)), make([]float64, len(chart.Thresholds))
	for i, s := range askSizes {
		for j, threshold := range chart.Thresholds {
			if askPrices[i] < askPrices[0]*(1.0+threshold) {
				askTotalSizes[j] += s
			}
		}
	}
	for i, s := range bidSizes {
		for j, threshold := range chart.Thresholds {
			if bidPrices[i] > bidPrices[0]*(1.0-threshold) {
				bidTotalSizes[j] += s
			}
		}
	}

	ratio := make([]float64, len(chart.Thresholds))
	for i, as := range askTotalSizes {
		ratio[i] = as/bidTotalSizes[i] - 1.0
	}

	if math.Signbit(ratio[0]) != chart.LastSign {
		chart.SetLastCandle(nil)
		chart.CurrentCandle = nil
		chart.LastSign = math.Signbit(ratio[0])
	}
}
