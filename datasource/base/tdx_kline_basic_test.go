package base

import (
	"fmt"
	"testing"
)

func TestUpdateAllBasicKLine(t *testing.T) {
	code := "sz300773"
	data := UpdateAllBasicKLine(code)
	fmt.Println(data[0])
	_ = data
}

func Test_totalAdjustmentTimes(t *testing.T) {
	type args struct {
		securityCode string
		startDate    string
		endDate      string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "sz300773",
			args: args{
				securityCode: "sz300773",
				startDate:    "1991-12-19",
				endDate:      "2025-07-08",
			},
			want: 5,
		},
		{
			name: "sz300773-20250603",
			args: args{
				securityCode: "sz300773",
				startDate:    "1991-12-19",
				endDate:      "2025-06-03",
			},
			want: 4,
		},
		{
			name: "sz300773-20250604",
			args: args{
				securityCode: "sz300773",
				startDate:    "1991-12-19",
				endDate:      "2025-06-04",
			},
			want: 4,
		},
		{
			name: "sz300773-20250605",
			args: args{
				securityCode: "sz300773",
				startDate:    "1991-12-19",
				endDate:      "2025-06-05",
			},
			want: 5,
		},
		{
			name: "sz300773-20250606",
			args: args{
				securityCode: "sz300773",
				startDate:    "1991-12-19",
				endDate:      "2025-06-06",
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := totalAdjustmentTimes(tt.args.securityCode, tt.args.startDate, tt.args.endDate); got != tt.want {
				t.Errorf("totalAdjustmentTimes() = %v, want %v", got, tt.want)
			}
		})
	}
}
