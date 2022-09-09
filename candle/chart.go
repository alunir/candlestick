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
	in               chan *Candle
	out              chan *Candle
	buffer           *RingBuffer[*Candle]
}

func NewChart(candleNum int) *Chart {
	in := make(chan *Candle)
	out := make(chan *Candle, candleNum)
	buffer := NewRingBuffer(in, out)
	go buffer.Run()
	return &Chart{
		CandleNum:  candleNum,
		Candles:    make([]*Candle, 0, candleNum),
		TimeSeries: map[time.Time]*Candle{},
		in:         in,
		out:        out,
		buffer:     buffer,
	}
}

func (chart *Chart) GetLastCandle() *Candle {
	return chart.LastCandle
}

func (chart *Chart) GetCurrentCandle() *Candle {
	return chart.CurrentCandle
}

func (chart *Chart) GetCandles() []*Candle {
	return chart.Candles
}

func (chart *Chart) GetLastCandleClock() chan *Candle {
	return chart.out
}

func (chart *Chart) SetLastCandle(candle *Candle) {
	// (candle, CurrentCandle) -> LastCandle
	// (nil, nil) -> no update
	// (nil, not nil) -> CurrentCandle
	// (not nil, nil) -> candle
	// (not nil, not nil) -> CurrentCandle
	if chart.CurrentCandle != nil {
		chart.LastCandle = chart.CurrentCandle
	} else {
		if candle != nil {
			chart.LastCandle = candle
		} else {
			// no update
			return
		}
	}
	chart.in <- chart.LastCandle
}

func (chart *Chart) AddCandle(candle *Candle) {
	chart.CurrentCandle = candle
	chart.CurrentCandleNew = true

	// Need?
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
