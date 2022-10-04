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
	Time   time.Time       `json:"time"`
	Open   decimal.Decimal `json:"open"`
	High   decimal.Decimal `json:"high"`
	Low    decimal.Decimal `json:"low"`
	Close  decimal.Decimal `json:"close"`
	Volume decimal.Decimal `json:"volume"`
	Amount decimal.Decimal `json:"amount"`
	Count  int             `json:"count"`
	Stack  decimal.Decimal `json:"stack"` // if this touch the resolution, add a new candle.
}

func NewCandle(cnt int, ti time.Time, value decimal.Decimal, volume decimal.Decimal, stack decimal.Decimal) Candle {
	return Candle{
		Count:  cnt,
		Time:   ti,
		High:   value,
		Low:    value,
		Open:   value,
		Close:  value,
		Volume: volume,
		Amount: value.Mul(volume),
		Stack:  stack.Abs(),
	}
}

func NewCandleWithBuySell(buysell BuySellType, ti time.Time, value decimal.Decimal, volume decimal.Decimal, stack decimal.Decimal) Candle {
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

func (candle *Candle) add(value decimal.Decimal, volume decimal.Decimal, stack decimal.Decimal) {
	if value.GreaterThan(candle.High) {
		candle.High = value
	} else if value.LessThan(candle.Low) {
		candle.Low = value
	}

	candle.Volume = candle.Volume.Add(volume)
	candle.Close = value
	candle.Count += 1
	candle.Amount = candle.Amount.Add(value.Mul(volume))
	candle.addStack(stack)
}

func (candle *Candle) addStack(stack decimal.Decimal) {
	candle.Stack = candle.Stack.Add(stack).Abs()
}

func (candle *Candle) AddCandleWithBuySell(buysell BuySellType, value decimal.Decimal, volume decimal.Decimal, stack decimal.Decimal) {
	switch buysell {
	case ALL:
		candle.add(value, volume.Abs(), stack)
	case BUY:
		if volume.GreaterThan(decimal.Zero) {
			candle.add(value, volume, stack)
		} else {
			candle.addStack(stack)
		}
	case SELL:
		if volume.LessThan(decimal.Zero) {
			candle.add(value, volume, stack)
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

func (c Candle) Equal(c2 Candle) bool {
	return c.Time.Equal(c2.Time) &&
		c.Open.Equal(c2.Open) &&
		c.High.Equal(c2.High) &&
		c.Low.Equal(c2.Low) &&
		c.Close.Equal(c2.Close) &&
		c.Volume.Equal(c2.Volume) &&
		c.Amount.Equal(c2.Amount) &&
		c.Count == c2.Count &&
		c.Stack.Equal(c2.Stack)
}

func (c Candle) AssertOhlcv(t *testing.T, ti time.Time, open, high, low, close, volume float64, count int) error {
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
	if c.Time != ti {
		return fmt.Errorf("got wrong Time val: %v but was %v, %v", ti, c.Time, c)
	}
	return nil
}
