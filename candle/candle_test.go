package candle

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestMarshalUnmashal(t *testing.T) {
	c := NewCandleWithBuySell(ALL, time.Now().Truncate(time.Nanosecond), decimal.NewFromFloat(1.0), decimal.NewFromFloat(2.0), decimal.NewFromFloat(3.0))

	b, err := c.Marshal()
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))

	c2 := Candle{}
	err = c2.Unmarshal(b)
	if err != nil {
		panic(err)
	}
	if c2.Time != c.Time {
		t.Errorf("time not equal. %v != %v", c2.Time, c.Time)
		t.Fail()
	}
	if !c2.Open.Equal(c.Open) {
		t.Errorf("open not equal. %v != %v", c2.Open, c.Open)
		t.Fail()
	}
	if !c2.High.Equal(c.High) {
		t.Errorf("high not equal. %v != %v", c2.High, c.High)
		t.Fail()
	}
	if !c2.Low.Equal(c.Low) {
		t.Errorf("low not equal. %v != %v", c2.Low, c.Low)
		t.Fail()
	}
	if !c2.Close.Equal(c.Close) {
		t.Errorf("close not equal. %v != %v", c2.Close, c.Close)
		t.Fail()
	}
	if !c2.Volume.Equal(c.Volume) {
		t.Errorf("volume not equal. %v != %v", c2.Volume, c.Volume)
		t.Fail()
	}
}
