package candle

import (
	"math"
	"time"
)

type Candle struct {
	Count  int
	Time   time.Time
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume float64
	Amount float64
	Stack  float64 // if this touch the resolution, add a new candle.
}

func NewCandle(cnt int, ti time.Time, value float64, volume float64, stack float64) *Candle {
	return &Candle{
		Count:  cnt,
		Time:   ti,
		High:   value,
		Low:    value,
		Open:   value,
		Close:  value,
		Volume: volume,
		Amount: value * volume,
		Stack:  math.Abs(stack),
	}
}

func NewCandleWithBuySell(buysell BuySellType, ti time.Time, value float64, volume float64, stack float64) *Candle {
	switch buysell {
	case ALL:
		return NewCandle(1, ti, value, math.Abs(volume), stack)
	case BUY:
		if volume > 0.0 {
			return NewCandle(1, ti, value, volume, stack)
		} else {
			return NewCandle(1, ti, value, 0.0, stack)
		}
	case SELL:
		if volume < 0.0 {
			return NewCandle(1, ti, value, volume, stack)
		} else {
			return NewCandle(1, ti, value, 0.0, stack)
		}
	default:
		panic("Failed to switch buysell")
	}
}

func (candle *Candle) add(value float64, volume float64, stack float64) {
	if value > candle.High {
		candle.High = value
	} else if value < candle.Low {
		candle.Low = value
	}

	candle.Volume += volume
	candle.Close = value
	candle.Count += 1
	candle.Amount += value * volume
	candle.addStack(stack)
}

func (candle *Candle) addStack(stack float64) {
	candle.Stack += math.Abs(stack)
}

func (candle *Candle) AddCandleWithBuySell(buysell BuySellType, value float64, volume float64, stack float64) {
	switch buysell {
	case ALL:
		candle.add(value, math.Abs(volume), stack)
	case BUY:
		if volume > 0.0 {
			candle.add(value, volume, stack)
		} else {
			candle.addStack(stack)
		}
	case SELL:
		if volume < 0.0 {
			candle.add(value, volume, stack)
		} else {
			candle.addStack(stack)
		}
	}
}
