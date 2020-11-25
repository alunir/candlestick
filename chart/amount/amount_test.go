package amount

import (
	"fmt"
	"testing"
	"time"

	"github.com/alunir/candlestick/candle"
)

func assertCandle(t *testing.T, c *candle.Candle, count int, volume float64, open float64, close float64, high float64, low float64, time time.Time) {
	expected := candle.Candle{
		Count:  count,
		Volume: volume,
		Open:   open,
		Close:  close,
		High:   high,
		Low:    low,
		Time:   time,
	}
	if c.Count != count {
		t.Logf("Got wrong Count val: %v but was %v, %v", expected.Count, c.Count, c)
		t.Fail()
		return
	}
	if c.Volume != volume {
		t.Logf("Got wrong Volume val: %v but was %v, %v", expected.Volume, c.Volume, c)
		t.Fail()
		return
	}
	if c.Open != open {
		t.Logf("Got wrong Open val: %v but was %v, %v", expected.Open, c.Open, c)
		t.Fail()
		return
	}
	if c.Close != close {
		t.Logf("Got wrong Close val: %v but was %v, %v", expected.Close, c.Close, c)
		t.Fail()
		return
	}
	if c.High != high {
		t.Logf("Got wrong High val: %v but was %v, %v", expected.High, c.High, c)
		t.Fail()
		return
	}
	if c.Low != low {
		t.Logf("Got wrong Low val: %v but was %v, %v", expected.Low, c.Low, c)
		t.Fail()
		return
	}
	if c.Time != time {
		t.Logf("Got wrong Time val: %v but was %v, %v", expected.Time, c.Time, c)
		t.Fail()
		return
	}
}

func TestAmountCandles(t *testing.T) {
	candleNum := 4
	var chart = &AmountChart{
		Chart: candle.Chart{
			CandleNum:  candleNum,
			Candles:    make([]*candle.Candle, 0, candleNum),
			TimeSeries: map[time.Time]*candle.Candle{},
			Clock:      make(chan *candle.Candle),
		},
		Chunk:   30.0,
		Buysell: candle.ALL,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	chart.AddTrade(start, 5, 1)
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	chart.AddTrade(start.Add(25*time.Second), 3, 1)
	var c1 = chart.Candles[0]

	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	chart.AddTrade(start.Add(70*time.Second), 13, 2)
	var c2 = chart.Candles[1]

	// Intentionally empty data series included here, to test flat candles
	chart.AddTrade(start.Add(240*time.Second), 15, 5)
	var c3 = chart.Candles[2]
	var c4 = chart.Candles[3]

	if !(c1.Count == 2 && c1.Open == 5 && c1.Close == 25 &&
		c1.High == 25 && c1.Low == 5 && c1.Volume == 2 && c1.Amount == 30) {
		t.Logf("Got wrong c1 val: %v", c1)
		t.Fail()
	}

	if !(c2.Count == 2 && c2.Open == 3 && c2.Close == 12 &&
		c2.High == 12 && c2.Low == 3 && c2.Volume == 3.25 && c2.Amount == 30) {
		t.Logf("Got wrong c2 val: %v", c2)
		t.Fail()
	}

	if !(c3.Count == 1 && c3.Open == 15 && c3.Close == 15 &&
		c3.High == 15 && c3.Low == 15 && c3.Volume == 2 && c3.Amount == 30) {
		t.Logf("Got wrong c3 val: %v", c3)
		t.Fail()
	}

	if !(c4.Count == 1 && c4.Open == 15 && c4.Close == 15 &&
		c4.High == 15 && c4.Low == 15 && c4.Volume == 0.9333333333333333 && c4.Amount == 14) {
		t.Logf("Got wrong c4 val: %v", c4)
		t.Fail()
	}

	if !(chart.LastCandle.Count == 1 && chart.LastCandle.Open == 15 && chart.LastCandle.Close == 15 &&
		chart.LastCandle.High == 15 && chart.LastCandle.Low == 15 && chart.LastCandle.Volume == 2 && chart.LastCandle.Amount == 30) {
		t.Logf("Got wrong chart.LastCandle val: %v", chart.LastCandle)
		t.Fail()
	}

	if !(chart.CurrentCandle.Count == 1 && chart.CurrentCandle.Open == 15 && chart.CurrentCandle.Close == 15 &&
		chart.CurrentCandle.High == 15 && chart.CurrentCandle.Low == 15 && chart.CurrentCandle.Volume == 0.9333333333333333 && chart.CurrentCandle.Amount == 14) {
		t.Logf("Got wrong chart.CurrentCandle val: %v", chart.CurrentCandle)
		t.Fail()
	}

	if len(chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(chart.Candles))
		t.Fail()
	}

	// Candles should be like a queue
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	chart.AddTrade(start.Add(300*time.Second), 10, 2)
	chart.AddTrade(start.Add(310*time.Second), 3, 6)
	chart.AddTrade(start.Add(370*time.Second), 54, 36)

	var c5 = chart.Candles[0]
	if !(c5.Count == 1 && c5.Open == 54 && c5.Close == 54 &&
		c5.High == 54 && c5.Low == 54 && c5.Volume == 0.5555555555555556) {
		t.Logf("Got wrong c5 val: %v", c5)
		t.Fail()
	}
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	if len(chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(chart.Candles))
		t.Fail()
	}
}

func TestBuySellAmountCandles(t *testing.T) {
	candleNum := 4
	var b_chart = &AmountChart{
		Chart: candle.Chart{
			CandleNum:  candleNum,
			Candles:    make([]*candle.Candle, 0, candleNum),
			TimeSeries: map[time.Time]*candle.Candle{},
			Clock:      make(chan *candle.Candle),
		},
		Chunk:   30.0,
		Buysell: candle.BUY,
	}
	var s_chart = &AmountChart{
		Chart: candle.Chart{
			CandleNum:  candleNum,
			Candles:    make([]*candle.Candle, 0, candleNum),
			TimeSeries: map[time.Time]*candle.Candle{},
			Clock:      make(chan *candle.Candle),
		},
		Chunk:   30.0,
		Buysell: candle.SELL,
	}

	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	b_chart.AddTrade(start, 5, 1)
	s_chart.AddTrade(start, 5, 1)

	b_chart.AddTrade(start.Add(5*time.Second), 25, 1)
	s_chart.AddTrade(start.Add(5*time.Second), 25, 1)

	b_chart.AddTrade(start.Add(25*time.Second), 3, 1)
	s_chart.AddTrade(start.Add(25*time.Second), 3, 1)

	b_chart.AddTrade(start.Add(60*time.Second), 12, -5)
	s_chart.AddTrade(start.Add(60*time.Second), 12, -5)

	b_chart.AddTrade(start.Add(70*time.Second), 13, 2)
	s_chart.AddTrade(start.Add(70*time.Second), 13, 2)

	b_chart.AddTrade(start.Add(240*time.Second), 15, -5)
	s_chart.AddTrade(start.Add(240*time.Second), 15, -5)

	// assert
	if len(b_chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(b_chart.Candles))
		t.Fail()
	}
	if len(s_chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(s_chart.Candles))
		t.Fail()
	}

	assertCandle(t, b_chart.Candles[0], 3, 2.316666666666667, 12, 15, 15, 12, start.Add(60*time.Second))
	assertCandle(t, s_chart.Candles[0], 1, 0, 12, 12, 12, 12, start.Add(60*time.Second))

	assertCandle(t, b_chart.Candles[1], 1, 2, 15, 15, 15, 15, start.Add(240*time.Second))
	assertCandle(t, s_chart.Candles[1], 1, 0, 15, 15, 15, 15, start.Add(240*time.Second))

	assertCandle(t, b_chart.Candles[2], 1, 2, 15, 15, 15, 15, start.Add(240*time.Second))
	assertCandle(t, s_chart.Candles[2], 1, 0, 15, 15, 15, 15, start.Add(240*time.Second))

	assertCandle(t, b_chart.Candles[3], 1, 0.9333333333333333, 15, 15, 15, 15, start.Add(240*time.Second))
	assertCandle(t, s_chart.Candles[3], 1, 0, 15, 15, 15, 15, start.Add(240*time.Second))
}

func TestBuyAmountOnlyCandles(t *testing.T) {
	candleNum := 4
	var b_chart = &AmountChart{
		Chart: candle.Chart{
			CandleNum:  candleNum,
			Candles:    make([]*candle.Candle, 0, candleNum),
			TimeSeries: map[time.Time]*candle.Candle{},
			Clock:      make(chan *candle.Candle),
		},
		Chunk:   30.0,
		Buysell: candle.BUY,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	b_chart.AddTrade(start, 5, 1)
	b_chart.AddTrade(start.Add(5*time.Second), 25, 1)
	b_chart.AddTrade(start.Add(25*time.Second), 3, 1)
	b_chart.AddTrade(start.Add(60*time.Second), 12, -8) // 8 * 12 > 30
	b_chart.AddTrade(start.Add(70*time.Second), 13, 3)  // 3 * 13 = 39
	b_chart.AddTrade(start.Add(240*time.Second), 15, -5)

	// assert
	if len(b_chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(b_chart.Candles))
		t.Fail()
	}

	// count, volume, open, close, high, low
	assertCandle(t, b_chart.Candles[0], 2, 2.184615384615385, 13, 15, 15, 13, start.Add(70*time.Second))
	assertCandle(t, b_chart.Candles[1], 1, 2, 15, 15, 15, 15, start.Add(240*time.Second))
	assertCandle(t, b_chart.Candles[2], 1, 2, 15, 15, 15, 15, start.Add(240*time.Second))
	assertCandle(t, b_chart.Candles[3], 1, 0.2, 15, 15, 15, 15, start.Add(240*time.Second))
}

func TestAMOUNTGetLastCandleClock(t *testing.T) {
	candleNum := 4
	var chart = &AmountChart{
		Chart: candle.Chart{
			CandleNum:  candleNum,
			Candles:    make([]*candle.Candle, 0, candleNum),
			TimeSeries: map[time.Time]*candle.Candle{},
			Clock:      make(chan *candle.Candle),
		},
		Chunk:   30.0,
		Buysell: candle.ALL,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	go func() {
		for {
			select {
			case lastCandle := <-chart.GetLastCandleClock():
				fmt.Printf("lastCandle is updated: %v\n", lastCandle)
			}
		}
	}()

	chart.AddTrade(start, 5, 1)
	time.Sleep(time.Second)
	fmt.Println("[TestAMOUNTGetLastCandleClock] 1 lastCandle should be updated next.")
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	chart.AddTrade(start.Add(25*time.Second), 3, 1)
	time.Sleep(time.Second)
	fmt.Println("[TestAMOUNTGetLastCandleClock] 2 lastCandle should be updated next.")
	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	chart.AddTrade(start.Add(70*time.Second), 13, 2)
	time.Sleep(time.Second)
}

func TestBUYAMOUNTGetLastCandleClock(t *testing.T) {
	candleNum := 4
	var chart = &AmountChart{
		Chart: candle.Chart{
			CandleNum:  candleNum,
			Candles:    make([]*candle.Candle, 0, candleNum),
			TimeSeries: map[time.Time]*candle.Candle{},
			Clock:      make(chan *candle.Candle),
		},
		Chunk:   30.0,
		Buysell: candle.BUY,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	go func() {
		for {
			select {
			case lastCandle := <-chart.GetLastCandleClock():
				fmt.Printf("lastCandle is updated: %v\n", lastCandle)
			}
		}
	}()

	chart.AddTrade(start, 5, 1)
	time.Sleep(time.Second)
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	fmt.Println("[TestBUYAMOUNTGetLastCandleClock] 1 lastCandle should be updated next.")
	chart.AddTrade(start.Add(25*time.Second), 3, 1)
	time.Sleep(time.Second)
	fmt.Println("[TestBUYAMOUNTGetLastCandleClock] 2 lastCandle should be updated next.")
	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	chart.AddTrade(start.Add(70*time.Second), 13, 2)
	time.Sleep(time.Second)
}
