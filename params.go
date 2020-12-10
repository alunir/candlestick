package candlestick

import (
	"fmt"
	"strings"

	c "github.com/alunir/candlestick/candle"
	"github.com/fatih/structs"
	"go.uber.org/zap/zapcore"
)

type ChartParameters struct {
	mode             c.ChartMode
	CandleNum        int
	Resolution       string
	ResidueThreshold float64
}

type ParameterOption func(*ChartParameters)

func Resolution(resol string) ParameterOption {
	return func(op *ChartParameters) {
		op.Resolution = strings.Replace(resol, "\"", "", -1)
	}
}

func ResidueThreshold(f float64) ParameterOption {
	return func(op *ChartParameters) {
		op.ResidueThreshold = f
	}
}

func ChartParameter(model string, resolution string, candleNum int, ops ...ParameterOption) *ChartParameters {
	params := ChartParameters{
		mode:       c.ChartMode{Value: strings.ToUpper(model)},
		Resolution: resolution,
		CandleNum:  candleNum,
	}
	for _, option := range ops {
		option(&params)
	}
	return &params
}

func (cp ChartParameters) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range structs.New(cp).Map() {
		enc.AddString(k, fmt.Sprint(v))
	}
	return nil
}
