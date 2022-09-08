package time

import (
	"fmt"
	"testing"
	"time"

	candle "github.com/alunir/candlestick/candle"
)

func TestTimeCandles(t *testing.T) {
	candleNum := 4
	var chart = &TimeChart{
		Chart: candle.Chart{
			CandleNum:  candleNum,
			Candles:    make([]*candle.Candle, 0, candleNum),
			TimeSeries: map[time.Time]*candle.Candle{},
			Clock:      make(chan *candle.Candle),
		},
		Resolution: time.Minute,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC).Truncate(chart.Resolution)

	chart.AddTrade(start, 5, 1)
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	chart.AddTrade(start.Add(25*time.Second), 3, 1)
	var c1 = chart.Candles[0]

	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	chart.AddTrade(start.Add(119*time.Second), 13, 2)
	var c2 = chart.Candles[1]

	// Intentionally empty data series included here, to test flat candles
	chart.AddTrade(start.Add(240*time.Second), 15, 1)
	chart.AddTrade(start.Add(299*time.Second), 15, 5)
	var c3 = chart.Candles[2]
	var c4 = chart.Candles[3]

	if !(c1.Volume == 3 && c1.Open == 5 && c1.Close == 3 &&
		c1.High == 25 && c1.Low == 3) {
		t.Logf("Got wrong val: %v", c1)
		t.Fail()
	}

	if !(c2.Volume == 7 && c2.Open == 12 && c2.Close == 13 &&
		c2.High == 13 && c2.Low == 12) {
		t.Logf("Got wrong val: %v", c2)
		t.Fail()
	}

	if !(c3.Volume == 0 && c3.Open == 13 && c3.Close == 13 &&
		c3.High == 13 && c3.Low == 13) {
		t.Logf("Got wrong val: %v", c3)
		t.Fail()
	}

	if !(c4.Volume == 6 && c4.Open == 15 && c4.Close == 15 &&
		c4.High == 15 && c4.Low == 15) {
		t.Logf("Got wrong val: %v", c4)
		t.Fail()
	}

	if !(chart.LastCandle.Volume == 7 && chart.LastCandle.Open == 12 && chart.LastCandle.Close == 13 &&
		chart.LastCandle.High == 13 && chart.LastCandle.Low == 12) {
		t.Logf("Got wrong chart.LastCandle val: %v", chart.LastCandle)
		t.Fail()
	}

	if !(chart.CurrentCandle.Volume == 6 && chart.CurrentCandle.Open == 15 && chart.CurrentCandle.Close == 15 &&
		chart.CurrentCandle.High == 15 && chart.CurrentCandle.Low == 15) {
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
	if !(c5.Volume == 0 && c5.Open == 13 && c5.Close == 13 &&
		c5.High == 13 && c5.Low == 13) {
		t.Logf("Got wrong val: %v", c5)
		t.Fail()
	}
	fmt.Printf("Got cap: %v len: %v\n", cap(chart.Candles), len(chart.Candles))
	if len(chart.Candles) != 4 {
		t.Logf("Got wrong len: %v", len(chart.Candles))
		t.Fail()
	}
}
