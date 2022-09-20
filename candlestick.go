package candlestick

import (
	"time"

	candle "github.com/alunir/candlestick/candle"
	amount_chart "github.com/alunir/candlestick/chart/amount"
	count_chart "github.com/alunir/candlestick/chart/count"
	time_chart "github.com/alunir/candlestick/chart/time"
	volume_chart "github.com/alunir/candlestick/chart/volume"
	"github.com/shopspring/decimal"
)

type Candlestick interface {
	GetLastCandleClock() chan candle.Candle
	AddTrade(ti time.Time, value float64, volume float64)
	AddLv2DataCallback(ti time.Time, askPrices []float64, askSizes []float64, bidPrices []float64, bidSizes []float64)
	AddCandle(candle.Candle)
	GetLastCandle() candle.Candle
	GetCurrentCandle() candle.Candle
	GetCandles() []candle.Candle
	GetChartInfo() map[string]interface{}
	Serialized() []byte
	Deserialized([]byte)
	SetLastCandle(candle candle.Candle)
}

func NewCandlestickChart[T time.Duration | float64 | int64](param *ChartParameters[T]) Candlestick {
	chart := candle.NewChart(param.CandleNum)
	switch param.mode {
	case candle.TIME:
		if _, ok := any(param.Resolution).(time.Duration); !ok {
			panic("Resolution must be time.Duration")
		}
		return time_chart.TimeChart{
			Chart:      chart,
			Resolution: time.Duration(param.Resolution),
		}
	case candle.AMOUNT:
		if _, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		}
		return amount_chart.AmountChart{
			Chart:   chart,
			Chunk:   decimal.NewFromFloat(float64(param.Resolution)),
			Buysell: candle.ALL,
		}
	case candle.BUY_AMOUNT:
		if _, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		}
		return amount_chart.AmountChart{
			Chart:   chart,
			Chunk:   decimal.NewFromFloat(float64(param.Resolution)),
			Buysell: candle.BUY,
		}
	case candle.SELL_AMOUNT:
		if _, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		}
		return amount_chart.AmountChart{
			Chart:   chart,
			Chunk:   decimal.NewFromFloat(float64(param.Resolution)),
			Buysell: candle.SELL,
		}
	case candle.VOLUME:
		if _, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		}
		return volume_chart.VolumeChart{
			Chart:   chart,
			Chunk:   decimal.NewFromFloat(float64(param.Resolution)),
			Buysell: candle.ALL,
		}
	case candle.BUY_VOLUME:
		if _, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		}
		return volume_chart.VolumeChart{
			Chart:   chart,
			Chunk:   decimal.NewFromFloat(float64(param.Resolution)),
			Buysell: candle.BUY,
		}
	case candle.SELL_VOLUME:
		if _, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		}
		return volume_chart.VolumeChart{
			Chart:   chart,
			Chunk:   decimal.NewFromFloat(float64(param.Resolution)),
			Buysell: candle.SELL,
		}
	case candle.COUNT:
		if _, ok := any(param.Resolution).(int64); !ok {
			panic("Resolution must be int64")
		}
		return count_chart.CountChart{
			Chart:   chart,
			Chunk:   int64(param.Resolution),
			Buysell: candle.ALL,
		}
	// case candle.RATIO:
	// 	thresholds_str := strings.Split(param.Resolution, ",")
	// 	var thresholds []float64
	// 	for _, s := range thresholds_str {
	// 		threshold, err := strconv.ParseFloat(s, 64)
	// 		if err != nil {
	// 			panic("Failed to parse resolution: " + s)
	// 		}
	// 		thresholds = append(thresholds, threshold)
	// 	}
	// 	if param.ResidueThreshold <= 0.0 {
	// 		panic("ResidueThreshold is zero or negative. This should be positive.")
	// 	}
	// 	return &ratio_chart.RatioChart{
	// 		Chart: candle.Chart{
	// 			CandleNum:  param.CandleNum,
	// 			Candles:    make([]*candle.Candle, 0, param.CandleNum),
	// 			TimeSeries: map[time.Time]*candle.Candle{},
	// 			Clock:      make(chan *candle.Candle),
	// 		},
	// 		Thresholds:       thresholds,
	// 		ResidueThreshold: param.ResidueThreshold,
	// 	}
	case candle.BUY_PRICE:
		panic("not implemented yet")
	case candle.SELL_PRICE:
		panic("not implemented yet")
	default:
		panic("Invalid ChartMode")
	}
}
