package time

import (
	"testing"
	"time"

	candle "github.com/alunir/candlestick/candle"
)

func TestTimeCandles(t *testing.T) {
	candleNum := 5
	interval := time.Minute
	var chart = &TimeChart{
		Chart:      candle.NewChart(candleNum),
		Resolution: interval,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC).Truncate(chart.Resolution)

	chart.AddTrade(start.Add(5*time.Second), 5, 1)
	chart.AddTrade(start.Add(15*time.Second), 25, 1)
	chart.AddTrade(start.Add(25*time.Second), 3, 1)

	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	chart.AddTrade(start.Add(119*time.Second), 13, 2)

	// 120-179, 180-239: No trade.

	// Intentionally empty data series included here, to test flat candles
	chart.AddTrade(start.Add(240*time.Second), 15, 1)
	chart.AddTrade(start.Add(299*time.Second), 15, 5)

	// To make the candle[4], input a trade at 300s.
	chart.AddTrade(start.Add(300*time.Second), 10, 2)

	// for debug
	// &{2009-11-10 23:30:00 +0000 UTC 5 25 3 3 3 33 3 0}
	// &{2009-11-10 23:31:00 +0000 UTC 12 13 12 13 7 86 2 0}
	// &{2009-11-10 23:32:00 +0000 UTC 13 13 13 13 0 0 0 0}
	// &{2009-11-10 23:33:00 +0000 UTC 13 13 13 13 0 0 0 0}
	// &{2009-11-10 23:34:00 +0000 UTC 15 15 15 15 6 90 2 0}
	// fmt.Printf("%v. Got cap: %v len: %v\n", chart.Candles, cap(chart.Candles), len(chart.Candles))
	// for _, c := range chart.Candles {
	// 	fmt.Printf("%v\n", c)
	// }

	// Candles should be like a queue
	if len(chart.Candles) != candleNum {
		t.Logf("Candles are not fulfilled. Size is %v", len(chart.Candles))
		t.Fail()
	}

	if err := chart.Candles[0].AssertOhlcv(t, start, 5, 25, 3, 3, 3, 3); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[1].AssertOhlcv(t, start.Add(chart.Resolution), 12, 13, 12, 13, 7, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[2].AssertOhlcv(t, start.Add(2*chart.Resolution), 13, 13, 13, 13, 0, 0); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[3].AssertOhlcv(t, start.Add(3*chart.Resolution), 13, 13, 13, 13, 0, 0); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.Candles[4].AssertOhlcv(t, start.Add(4*chart.Resolution), 15, 15, 15, 15, 6, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	// This is the Last candle. "Last" means "last traded except current candle"
	if err := chart.LastCandle.AssertOhlcv(t, start.Add(4*chart.Resolution), 15, 15, 15, 15, 6, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	if err := chart.CurrentCandle.AssertOhlcv(t, start.Add(5*chart.Resolution), 10, 10, 10, 10, 2, 1); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}

	// To check, candles are cyclic
	chart.AddTrade(start.Add(360*time.Second), 3, 6)

	// for debug
	// &{2009-11-10 23:31:00 +0000 UTC 12 13 12 13 7 86 2 0}
	// &{2009-11-10 23:32:00 +0000 UTC 13 13 13 13 0 0 0 0}
	// &{2009-11-10 23:33:00 +0000 UTC 13 13 13 13 0 0 0 0}
	// &{2009-11-10 23:34:00 +0000 UTC 15 15 15 15 6 90 2 0}
	// &{2009-11-10 23:35:00 +0000 UTC 10 10 3 3 8 38 2 0}
	// fmt.Printf("%v. Got cap: %v len: %v\n", chart.Candles, cap(chart.Candles), len(chart.Candles))
	// for _, c := range chart.Candles {
	// 	fmt.Printf("%v\n", c)
	// }

	if err := chart.Candles[0].AssertOhlcv(t, start.Add(chart.Resolution), 12, 13, 12, 13, 7, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
}
