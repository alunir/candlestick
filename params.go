package candlestick

import (
	"fmt"
	"strings"
	"time"

	c "github.com/alunir/candlestick/candle"
	"github.com/fatih/structs"
	"go.uber.org/zap/zapcore"
)

type ChartParameters[T time.Duration | float64 | int64] struct {
	mode             c.ChartMode
	CandleNum        int
	Resolution       T
	ResidueThreshold float64
}

type ParameterOption[T time.Duration | float64 | int64] func(*ChartParameters[T])

func Resolution[T time.Duration | float64 | int64](resolution T) ParameterOption[T] {
	return func(op *ChartParameters[T]) {
		op.Resolution = resolution
	}
}

func ResidueThreshold[T time.Duration | float64 | int64](f float64) ParameterOption[T] {
	return func(op *ChartParameters[T]) {
		op.ResidueThreshold = f
	}
}

func ChartParameter[T time.Duration | float64 | int64](model string, resolution T, candleNum int, ops ...ParameterOption[T]) *ChartParameters[T] {
	params := ChartParameters[T]{
		mode:       c.ChartMode(strings.ToUpper(model)),
		Resolution: resolution,
		CandleNum:  candleNum,
	}
	for _, option := range ops {
		option(&params)
	}
	return &params
}

func (cp ChartParameters[T]) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range structs.New(cp).Map() {
		enc.AddString(k, fmt.Sprint(v))
	}
	return nil
}
