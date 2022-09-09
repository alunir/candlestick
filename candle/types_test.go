package candle

import (
	"testing"
)

func TestChartMode(t *testing.T) {
	m := ChartMode("VOLUME")
	if m != VOLUME {
		t.Logf("Got wrong val: %v", m)
		t.Fail()
	}
}
