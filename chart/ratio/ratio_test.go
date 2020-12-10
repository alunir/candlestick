package ratio

import (
	"testing"
	"time"

	"github.com/alunir/candlestick/candle"
)

func TestRatioCandles(t *testing.T) {
	candleNum := 4
	var chart = &RatioChart{
		Chart: candle.Chart{
			CandleNum:  candleNum,
			Candles:    make([]*candle.Candle, 0, candleNum),
			TimeSeries: map[time.Time]*candle.Candle{},
			Clock:      make(chan *candle.Candle),
		},
		Thresholds:       []float64{0.001, 0.005},
		ResidueThreshold: 1.0,
	}
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	chart.AddLv2DataCallback(start, []float64{100.0, 101.0, 102.0}, []float64{25, 30, 10}, []float64{99.0, 98.0, 97.0}, []float64{21, 16, 10})
	chart.AddTrade(start, 5, 1)
	chart.AddTrade(start.Add(5*time.Second), 25, 1)
	chart.AddTrade(start.Add(25*time.Second), 3, 1)
	chart.AddLv2DataCallback(start.Add(25*time.Second), []float64{100.0, 101.0, 102.0}, []float64{5, 3, 1}, []float64{99.0, 98.0, 97.0}, []float64{21, 16, 10})
	chart.AddTrade(start.Add(60*time.Second), 12, 5)
	chart.AddTrade(start.Add(70*time.Second), 13, 2)
	chart.AddLv2DataCallback(start.Add(70*time.Second), []float64{100.0, 101.0, 102.0}, []float64{15, 13, 11}, []float64{99.0, 98.0, 97.0}, []float64{1, 6, 0})
	chart.AddTrade(start.Add(80*time.Second), 15, 5)
	chart.AddTrade(start.Add(90*time.Second), 16, 2)
	chart.AddLv2DataCallback(start.Add(90*time.Second), []float64{100.0, 101.0, 102.0}, []float64{25, 23, 21}, []float64{99.0, 98.0, 97.0}, []float64{2, 1, 1})
	chart.AddTrade(start.Add(100*time.Second), 11, 6)
	var c1 = chart.Candles[0]

	// Intentionally empty data series included here, to test flat candles
	chart.AddLv2DataCallback(start.Add(110*time.Second), []float64{100.0, 101.0, 102.0}, []float64{5, 3, 1}, []float64{99.0, 98.0, 97.0}, []float64{25, 16, 10})
	chart.AddTrade(start.Add(120*time.Second), 12, 5)
	chart.AddTrade(start.Add(130*time.Second), 13, 2)
	var c2 = chart.Candles[1]

	chart.AddLv2DataCallback(start.Add(140*time.Second), []float64{100.0, 101.0, 102.0}, []float64{125, 123, 211}, []float64{99.0, 98.0, 97.0}, []float64{2, 1, 1})
	chart.AddTrade(start.Add(150*time.Second), 15, 6)
	var c3 = chart.Candles[2]

	if !(c1.Open == 5 && c1.High == 25 && c1.Low == 3 && c1.Close == 16 && c1.Volume == 17 && c1.Amount == 226 && c1.Count == 7) {
		t.Logf("Got wrong c1 val: %v", c1)
		t.Fail()
	}

	if !(c2.Open == 11 && c2.High == 11 && c2.Low == 11 && c2.Close == 11 && c2.Volume == 6 && c2.Amount == 66 && c2.Count == 1) {
		t.Logf("Got wrong c2 val: %v", c2)
		t.Fail()
	}

	if !(c3.Open == 12 && c3.High == 15 && c3.Low == 12 && c3.Close == 15 && c3.Volume == 13 && c3.Amount == 176 && c3.Count == 3) {
		t.Logf("Got wrong c3 val: %v", c3)
		t.Fail()
	}

	if len(chart.Candles) != 3 {
		t.Logf("Got wrong len: %v", len(chart.Candles))
		t.Fail()
	}

	info := chart.GetChartInfo()
	as := info["AskTotalSizes"].([]float64)
	if !(as[0] == 200.0 && as[1] == 200.0) {
		t.Logf("Got wrong AskTotalSizes val: %v", as)
		t.Fail()
	}
	bs := info["BidTotalSizes"].([]float64)
	if !(bs[0] == 72.0 && bs[1] == 72.0) {
		t.Logf("Got wrong BidTotalSizes val: %v", bs)
		t.Fail()
	}
	rs := info["Ratio"].([]float64)
	if !(rs[0] == 1.0216512475319814 && rs[1] == 1.0216512475319814) {
		t.Logf("Got wrong Ratios val: %v", rs)
		t.Fail()
	}

}
