package indicators

import (
	"fmt"
	"testing"

	"github.com/quant1x/engine/datasource/base"
	"github.com/quant1x/exchange"
	"github.com/quant1x/num/labs"
	"github.com/quant1x/pandas"
)

func TestSar_basic(t *testing.T) {
	code := "300046"
	date := "2024-05-27"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	list := base.CheckoutKLines(code, date)
	rows := len(list)
	high := make([]float64, rows)
	low := make([]float64, rows)
	for i, v := range list {
		high[i] = v.High
		low[i] = v.Low
	}
	data := SAR(high, low)
	df := pandas.LoadStructs(data)
	fmt.Println(df)
	last := data[rows-1]
	// 增量计算SAR 2024-06-14
	latest := last.Incr(18.77, 17.88)
	fmt.Printf("%+v\n", latest)
}

func TestSAR(t *testing.T) {
	code := "600171"
	date := "2024-06-13"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	list := base.CheckoutKLines(code, date)
	rows := len(list)
	highs := make([]float64, rows)
	lows := make([]float64, rows)
	for i, v := range list {
		highs[i] = v.High
		lows[i] = v.Low
	}
	type args struct {
		highs []float64
		lows  []float64
	}
	tests := []struct {
		name string
		args args
		want FeatureSar
	}{
		{
			name: "上海贝岭(sh600171)",
			args: args{highs: highs, lows: lows},
			want: FeatureSar{
				Pos:    6121,
				Bull:   true,
				Af:     0.20,
				Ep:     20.45,
				Sar:    17.30,
				High:   20.45,
				Low:    18.90,
				Period: 12,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SAR(tt.args.highs, tt.args.lows)
			got := result[len(result)-1]
			if !labs.DeepEqual(got, tt.want) {
				t.Errorf("SAR() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFeatureSar_Incr(t *testing.T) {
	type fields struct {
		Pos    int
		Bull   bool
		Af     float64
		Ep     float64
		Sar    float64
		High   float64
		Low    float64
		Period int
	}
	type args struct {
		high float64
		low  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   FeatureSar
	}{
		{
			name: "上海贝岭(sh600171)-2024-06-13",
			fields: fields{
				Pos:    6121,
				Bull:   true,
				Af:     0.20,
				Ep:     20.45,
				Sar:    17.30,
				High:   20.45,
				Low:    18.90,
				Period: 12,
			},
			args: args{high: 18.77, low: 17.88}, // 2024-06-14
			want: FeatureSar{
				Pos:    6122,
				Bull:   true,
				Af:     0.20,
				Ep:     20.45,
				Sar:    17.88,
				High:   18.77,
				Low:    17.88,
				Period: 13,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FeatureSar{
				Pos:    tt.fields.Pos,
				Bull:   tt.fields.Bull,
				Af:     tt.fields.Af,
				Ep:     tt.fields.Ep,
				Sar:    tt.fields.Sar,
				High:   tt.fields.High,
				Low:    tt.fields.Low,
				Period: tt.fields.Period,
			}
			if got := s.Incr(tt.args.high, tt.args.low); !labs.DeepEqual(got, tt.want) {
				t.Errorf("Incr() = %v, want %v", got, tt.want)
			}
		})
	}
}
