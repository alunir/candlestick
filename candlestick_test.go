package candlestick

import (
	"testing"
	"time"
)

func TestNewCandlestickChart(t *testing.T) {
	interval := time.Second
	param := ChartParameter("time", interval, 8)
	c := NewCandlestickChart(param)

	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC).Truncate(interval)
	c.AddTrade(start, 5, 2)
	c.AddTrade(start, 15, 1)
	c.AddTrade(start.Add(1*time.Second), 25, 1)

	if err := c.GetLastCandle().AssertOhlcv(t, start, 5, 15, 5, 15, 3, 2); err != nil {
		t.Logf("test failed. %v", err)
		t.Fail()
	}
}
