package candle

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Chart struct {
	Candles []*Candle
	// Resolution       string // TimeCandles: time.Duration(Nanoseconds, int64), VolumeCandles: float64, AmountCandles: float64, CountCandles: int
	TimeSeries       map[time.Time]*Candle
	LastCandle       *Candle
	CurrentCandle    *Candle
	CurrentCandleNew bool
	StartTime        time.Time
	EndTime          time.Time
	CandleNum        int
	Clock            chan *Candle
}

func (chart *Chart) GetLastCandleClock() chan *Candle {
	return chart.Clock
}

func (chart *Chart) SetLastCandle(candle *Candle) {
	if chart.CurrentCandle == nil {
		chart.LastCandle = candle
	} else {
		chart.LastCandle = chart.CurrentCandle
	}
	select {
	case chart.Clock <- chart.LastCandle:
	default:
	}
}

func (chart *Chart) AddCandle(candle *Candle) {
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

func (c Chart) Serialized() []byte {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(&c)
	if err != nil {
		panic("Failed to Serialized")
	}
	return buf.Bytes()
}

func (c *Chart) Deserialized(b []byte) {
	err := gob.NewDecoder(bytes.NewBuffer(b)).Decode(&c)
	if err != nil {
		panic("Failed to Deserialized. " + err.Error())
	}
}
