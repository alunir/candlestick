package candlestick

import (
	"strconv"
	"time"

	candle "github.com/alunir/candlestick/candle"
	amount_chart "github.com/alunir/candlestick/chart/amount"
	count_chart "github.com/alunir/candlestick/chart/count"
	ratio_chart "github.com/alunir/candlestick/chart/ratio"
	time_chart "github.com/alunir/candlestick/chart/time"
	volume_chart "github.com/alunir/candlestick/chart/volume"
)

type Candlestick interface {
	GetLastCandleClock() chan *candle.Candle
	AddTrade(ti time.Time, value float64, volume float64)
	AddCandle(*candle.Candle)
	Serialized() []byte
	Deserialized([]byte)
	SetLastCandle(candle *candle.Candle)
}

func NewCandlestickChart(param *ChartParameters) Candlestick {
	switch param.mode {
	case candle.TIME:
		resolution, err := time.ParseDuration(param.Resolution)
		if err != nil {
			panic("Failed to parse resolution")
		}
		return &time_chart.TimeChart{
			Chart: candle.Chart{
				CandleNum:  param.CandleNum,
				Candles:    make([]*candle.Candle, 0, param.CandleNum),
				TimeSeries: map[time.Time]*candle.Candle{},
				Clock:      make(chan *candle.Candle),
			},
			Resolution: resolution,
		}
	case candle.AMOUNT:
		chunk, err := strconv.ParseFloat(param.Resolution, 64)
		if err != nil {
			panic("Failed to parse resolution")
		}
		return &amount_chart.AmountChart{
			Chart: candle.Chart{
				CandleNum:  param.CandleNum,
				Candles:    make([]*candle.Candle, 0, param.CandleNum),
				TimeSeries: map[time.Time]*candle.Candle{},
				Clock:      make(chan *candle.Candle),
			},
			Chunk:   chunk,
			Buysell: candle.ALL,
		}
	case candle.BUY_AMOUNT:
		chunk, err := strconv.ParseFloat(param.Resolution, 64)
		if err != nil {
			panic("Failed to parse resolution")
		}
		return &amount_chart.AmountChart{
			Chart: candle.Chart{
				CandleNum:  param.CandleNum,
				Candles:    make([]*candle.Candle, 0, param.CandleNum),
				TimeSeries: map[time.Time]*candle.Candle{},
				Clock:      make(chan *candle.Candle),
			},
			Chunk:   chunk,
			Buysell: candle.BUY,
		}
	case candle.SELL_AMOUNT:
		chunk, err := strconv.ParseFloat(param.Resolution, 64)
		if err != nil {
			panic("Failed to parse resolution")
		}
		return &amount_chart.AmountChart{
			Chart: candle.Chart{
				CandleNum:  param.CandleNum,
				Candles:    make([]*candle.Candle, 0, param.CandleNum),
				TimeSeries: map[time.Time]*candle.Candle{},
				Clock:      make(chan *candle.Candle),
			},
			Chunk:   chunk,
			Buysell: candle.SELL,
		}
	case candle.VOLUME:
		chunk, err := strconv.ParseFloat(param.Resolution, 64)
		if err != nil {
			panic("Failed to parse resolution")
		}
		return &volume_chart.VolumeChart{
			Chart: candle.Chart{
				CandleNum:  param.CandleNum,
				Candles:    make([]*candle.Candle, 0, param.CandleNum),
				TimeSeries: map[time.Time]*candle.Candle{},
				Clock:      make(chan *candle.Candle),
			},
			Chunk:   chunk,
			Buysell: candle.ALL,
		}
	case candle.BUY_VOLUME:
		chunk, err := strconv.ParseFloat(param.Resolution, 64)
		if err != nil {
			panic("Failed to parse resolution")
		}
		return &volume_chart.VolumeChart{
			Chart: candle.Chart{
				CandleNum:  param.CandleNum,
				Candles:    make([]*candle.Candle, 0, param.CandleNum),
				TimeSeries: map[time.Time]*candle.Candle{},
				Clock:      make(chan *candle.Candle),
			},
			Chunk:   chunk,
			Buysell: candle.BUY,
		}
	case candle.SELL_VOLUME:
		chunk, err := strconv.ParseFloat(param.Resolution, 64)
		if err != nil {
			panic("Failed to parse resolution")
		}
		return &volume_chart.VolumeChart{
			Chart: candle.Chart{
				CandleNum:  param.CandleNum,
				Candles:    make([]*candle.Candle, 0, param.CandleNum),
				TimeSeries: map[time.Time]*candle.Candle{},
				Clock:      make(chan *candle.Candle),
			},
			Chunk:   chunk,
			Buysell: candle.SELL,
		}
	case candle.COUNT:
		chunk, err := strconv.ParseInt(param.Resolution, 10, 64)
		if err != nil {
			panic("Failed to parse resolution")
		}
		return &count_chart.CountChart{
			Chart: candle.Chart{
				CandleNum:  param.CandleNum,
				Candles:    make([]*candle.Candle, 0, param.CandleNum),
				TimeSeries: map[time.Time]*candle.Candle{},
				Clock:      make(chan *candle.Candle),
			},
			Chunk:   chunk,
			Buysell: candle.ALL,
		}
	case candle.RATIO:
		threshold, err := strconv.ParseFloat(param.Resolution, 64)
		if err != nil {
			panic("Failed to parse resolution")
		}
		return &ratio_chart.RatioChart{
			Chart: candle.Chart{
				CandleNum:  param.CandleNum,
				Candles:    make([]*candle.Candle, 0, param.CandleNum),
				TimeSeries: map[time.Time]*candle.Candle{},
				Clock:      make(chan *candle.Candle),
			},
			Threshold: threshold,
		}
	case candle.BUY_PRICE:
		panic("not implemented yet")
	case candle.SELL_PRICE:
		panic("not implemented yet")
	default:
		panic("Invalid ChartMode")
	}
}
