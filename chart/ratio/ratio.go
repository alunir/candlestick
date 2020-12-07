package ratio

import (
	"math"
	"time"

	c "github.com/alunir/candlestick/candle"
	"github.com/fatih/structs"
)

type ChartInfo struct {
	AskTotalSizes []float64
	BidTotalSizes []float64
	Ratio         []float64
}

type RatioChart struct {
	c.Chart
	Thresholds []float64
	LastSign   bool
	ChartInfo
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
	if structs.New(chart.ChartInfo).IsZero() {
		chart.ChartInfo = ChartInfo{
			AskTotalSizes: make([]float64, len(chart.Thresholds)),
			BidTotalSizes: make([]float64, len(chart.Thresholds)),
			Ratio:         make([]float64, len(chart.Thresholds)),
		}
	}
	for i, s := range askSizes {
		for j, threshold := range chart.Thresholds {
			if askPrices[i] < askPrices[0]*(1.0+threshold) {
				chart.ChartInfo.AskTotalSizes[j] += s
			}
		}
	}
	for i, s := range bidSizes {
		for j, threshold := range chart.Thresholds {
			if bidPrices[i] > bidPrices[0]*(1.0-threshold) {
				chart.ChartInfo.BidTotalSizes[j] += s
			}
		}
	}
	for i, as := range chart.AskTotalSizes {
		chart.ChartInfo.Ratio[i] = as/chart.ChartInfo.BidTotalSizes[i] - 1.0
	}

	if math.Signbit(chart.ChartInfo.Ratio[0]) != chart.LastSign {
		chart.SetLastCandle(nil)
		chart.CurrentCandle = nil
		chart.LastSign = math.Signbit(chart.ChartInfo.Ratio[0])
	}
}

func (chart *RatioChart) GetChartInfo() map[string]interface{} {
	return structs.New(chart.ChartInfo).Map()
}
