package candlestick

import (
	"context"
	"time"

	candle "github.com/alunir/candlestick/candle"
	amount_chart "github.com/alunir/candlestick/chart/amount"
	count_chart "github.com/alunir/candlestick/chart/count"
	time_chart "github.com/alunir/candlestick/chart/time"
	volume_chart "github.com/alunir/candlestick/chart/volume"
)

type Candlestick interface {
	GetLastCandleUpdate() chan candle.Candle
	GetCandleClock(context.Context, time.Duration) chan candle.Candle
	AddTrade(ti time.Time, value float64, volume float64)
	GetLastCandle() candle.Candle
	GetCurrentCandle() candle.Candle
	GetCandles() []candle.Candle
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	AddCandle(candle.Candle)
	SetLastCandle(candle candle.Candle)
}

func NewCandlestickChart[T time.Duration | float64 | int64](param *ChartParameters[T]) Candlestick {
	switch param.mode {
	case candle.TIME:
		if resolution, ok := any(param.Resolution).(time.Duration); !ok {
			panic("Resolution must be time.Duration")
		} else {
			return time_chart.NewTimeChart(resolution, param.CandleNum)
		}
	case candle.AMOUNT:
		if chunk, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		} else {
			return amount_chart.NewAmountChart(chunk, candle.ALL, param.CandleNum)
		}
	case candle.BUY_AMOUNT:
		if chunk, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		} else {
			return amount_chart.NewAmountChart(chunk, candle.BUY, param.CandleNum)
		}
	case candle.SELL_AMOUNT:
		if chunk, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		} else {
			return amount_chart.NewAmountChart(chunk, candle.SELL, param.CandleNum)
		}
	case candle.VOLUME:
		if chunk, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		} else {
			return volume_chart.NewVolumeChart(chunk, candle.ALL, param.CandleNum)
		}
	case candle.BUY_VOLUME:
		if chunk, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		} else {
			return volume_chart.NewVolumeChart(chunk, candle.BUY, param.CandleNum)
		}
	case candle.SELL_VOLUME:
		if chunk, ok := any(param.Resolution).(float64); !ok {
			panic("Resolution must be float64")
		} else {
			return volume_chart.NewVolumeChart(chunk, candle.SELL, param.CandleNum)
		}
	case candle.COUNT:
		if chunk, ok := any(param.Resolution).(int64); !ok {
			panic("Resolution must be float64")
		} else {
			return count_chart.NewCountChart(chunk, candle.ALL, param.CandleNum)
		}
	case candle.BUY_COUNT:
		if chunk, ok := any(param.Resolution).(int64); !ok {
			panic("Resolution must be float64")
		} else {
			return count_chart.NewCountChart(chunk, candle.BUY, param.CandleNum)
		}
	case candle.SELL_COUNT:
		if chunk, ok := any(param.Resolution).(int64); !ok {
			panic("Resolution must be float64")
		} else {
			return count_chart.NewCountChart(chunk, candle.SELL, param.CandleNum)
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
	// case candle.BUY_PRICE:
	// 	panic("not implemented yet")
	// case candle.SELL_PRICE:
	// 	panic("not implemented yet")
	default:
		panic("Invalid ChartMode")
	}
}
