package candle

import (
	"fmt"
	"testing"
	"time"

	json "github.com/goccy/go-json"

	"github.com/shopspring/decimal"
)

// Time, OHLCV,A,C
type Candle struct {
	Open      decimal.Decimal
	OpenTime  time.Time
	High      decimal.Decimal
	HighTime  time.Time
	Low       decimal.Decimal
	LowTime   time.Time
	Close     decimal.Decimal
	CloseTime time.Time
	Volume    decimal.Decimal
	Amount    decimal.Decimal
	Count     int
	Stack     decimal.Decimal // if this touch the resolution, add a new candle.
}

func NewCandle(cnt int, ti time.Time, value, volume, stack decimal.Decimal) Candle {
	return Candle{
		Open:      value,
		OpenTime:  ti,
		High:      value,
		HighTime:  ti,
		Low:       value,
		LowTime:   ti,
		Close:     value,
		CloseTime: ti,
		Volume:    volume,
		Amount:    value.Mul(volume),
		Count:     cnt,
		Stack:     stack.Abs(),
	}
}

func NewCandleWithBuySell(buysell BuySellType, ti time.Time, value, volume, stack decimal.Decimal) Candle {
	switch buysell {
	case ALL:
		return NewCandle(1, ti, value, volume.Abs(), stack)
	case BUY:
		if volume.IsPositive() {
			return NewCandle(1, ti, value, volume, stack)
		} else {
			return NewCandle(1, ti, value, decimal.Zero, stack)
		}
	case SELL:
		if volume.IsNegative() {
			return NewCandle(1, ti, value, volume, stack)
		} else {
			return NewCandle(1, ti, value, decimal.Zero, stack)
		}
	default:
		panic("Failed to switch buysell")
	}
}

func (candle *Candle) add(ti time.Time, value, volume, stack decimal.Decimal) {
	if value.GreaterThan(candle.High) {
		candle.HighTime = ti
		candle.High = value
	} else if value.LessThan(candle.Low) {
		candle.LowTime = ti
		candle.Low = value
	}

	candle.Volume = candle.Volume.Add(volume)
	candle.Close = value
	candle.CloseTime = ti
	candle.Count += 1
	candle.Amount = candle.Amount.Add(value.Mul(volume))
	candle.addStack(stack)
}

func (candle *Candle) addStack(stack decimal.Decimal) {
	candle.Stack = candle.Stack.Add(stack).Abs()
}

func (candle *Candle) AddCandleWithBuySell(buysell BuySellType, ti time.Time, value, volume, stack decimal.Decimal) {
	switch buysell {
	case ALL:
		candle.add(ti, value, volume.Abs(), stack)
	case BUY:
		if volume.GreaterThan(decimal.Zero) {
			candle.add(ti, value, volume, stack)
		} else {
			candle.addStack(stack)
		}
	case SELL:
		if volume.LessThan(decimal.Zero) {
			candle.add(ti, value, volume, stack)
		} else {
			candle.addStack(stack)
		}
	}
}

func (c Candle) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Candle) Unmarshal(b []byte) error {
	return json.Unmarshal(b, c)
}

func (c Candle) IsZero() bool {
	return c.Count == 0 // satisfied?
}

func (c Candle) Copy() Candle {
	return Candle{
		Open:      c.Open.Copy(),
		OpenTime:  c.OpenTime,
		High:      c.High.Copy(),
		HighTime:  c.HighTime,
		Low:       c.Low.Copy(),
		LowTime:   c.LowTime,
		Close:     c.Close.Copy(),
		CloseTime: c.CloseTime,
		Volume:    c.Volume.Copy(),
		Amount:    c.Amount.Copy(),
		Count:     c.Count,
		Stack:     c.Stack.Copy(),
	}
}

func (c Candle) Equal(c2 Candle) bool {
	return c.OpenTime.Equal(c2.OpenTime) &&
		c.Open.Equal(c2.Open) &&
		c.High.Equal(c2.High) &&
		c.HighTime.Equal(c2.HighTime) &&
		c.Low.Equal(c2.Low) &&
		c.LowTime.Equal(c2.LowTime) &&
		c.Close.Equal(c2.Close) &&
		c.CloseTime.Equal(c2.CloseTime) &&
		c.Volume.Equal(c2.Volume) &&
		c.Amount.Equal(c2.Amount) &&
		c.Count == c2.Count &&
		c.Stack.Equal(c2.Stack)
}

func (c Candle) AssertOhlcv(t *testing.T, open_time, high_time, low_time, close_time time.Time, open, high, low, close, volume float64, count int) error {
	if c.Count != count {
		return fmt.Errorf("got wrong Count val: %v but was %v, %v", count, c.Count, c)
	}
	if !c.Volume.Equal(decimal.NewFromFloat(volume)) {
		return fmt.Errorf("got wrong Volume val: %v but was %v, %v", volume, c.Volume, c)
	}
	if !c.Open.Equal(decimal.NewFromFloat(open)) {
		return fmt.Errorf("got wrong Open val: %v but was %v, %v", open, c.Open, c)
	}
	if !c.Close.Equal(decimal.NewFromFloat(close)) {
		return fmt.Errorf("got wrong Close val: %v but was %v, %v", close, c.Close, c)
	}
	if !c.High.Equal(decimal.NewFromFloat(high)) {
		return fmt.Errorf("got wrong High val: %v but was %v, %v", high, c.High, c)
	}
	if !c.Low.Equal(decimal.NewFromFloat(low)) {
		return fmt.Errorf("got wrong Low val: %v but was %v, %v", low, c.Low, c)
	}
	if c.OpenTime != open_time {
		return fmt.Errorf("got wrong OpenTime val: %v but was %v, %v", open_time, c.OpenTime, c)
	}
	if c.HighTime != high_time {
		return fmt.Errorf("got wrong HighTime val: %v but was %v, %v", high_time, c.HighTime, c)
	}
	if c.LowTime != low_time {
		return fmt.Errorf("got wrong LowTime val: %v but was %v, %v", low_time, c.LowTime, c)
	}
	if c.CloseTime != close_time {
		return fmt.Errorf("got wrong CloseTime val: %v but was %v, %v", close_time, c.CloseTime, c)
	}
	return nil
}
