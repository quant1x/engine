package config

import (
	"fmt"
	"testing"
)

func TestTimeRange(t *testing.T) {
	text := " 09:30:00 ~ 14:56:30 "
	text = " 14:56:30 - 09:30:00 "
	tr := ParseTimeRange(text)
	fmt.Println(tr)
	text = "09:15:00~09:19:59,09:25:00~11:29:59,13:00:00~14:59:59"
	ts := ParseTradingSession(text)
	fmt.Println(ts)
}

//func TestAddTest(t *testing.T) {
//	type args struct {
//		s1 *intervalset.Set
//		s2 *intervalset.Set
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantRes *intervalset.Set
//	}{
//		{
//			name: "empty + empty = empty",
//			args: args{
//				s1: intervalset.NewSet([]intervalset.Interval{}),
//				s2: intervalset.NewSet([]intervalset.Interval{}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{}),
//		},
//		{
//			name: "empty + [30,111] = [30, 111]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{}),
//				intervalset.NewSet([]intervalset.Interval{&span{30, 111}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{30, 111}}),
//		},
//		{
//			name: "[20, 40] + empty = [20, 40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//		},
//		{
//			name: "[20, 40] + [60,111]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{60, 111}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 40}, &span{60, 111}}),
//		},
//		{
//			name: "[20, 40] + [39,111]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{39, 111}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 111}}),
//		},
//		{
//			name: "[20, 40] + [40,111]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{40, 111}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 111}}),
//		},
//		{
//			name: "[20, 40] + [41,111]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{41, 111}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 111}}),
//		},
//		{
//			name: "[20, 40] + [30,111]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{30, 111}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 111}}),
//		},
//		{
//			name: " [20,40] + [25, 28]  = [20,40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{25, 28}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//		},
//		{
//			name: " [20,40] + [10, 30]  = [10,40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{10, 30}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 40}}),
//		},
//		{
//			name: "[10, 19] + [20,40]  = [10,40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{10, 19}}),
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 40}}),
//		},
//		{
//			name: "[20,40] + [10, 20]  = [10,40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{10, 20}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 40}}),
//		},
//		{
//			name: "[20,40] + [10, 21]  = [10,40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{10, 21}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 40}}),
//		},
//		{
//			name: "[20,40] + [10, 15]  = [10,15] [20,40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{10, 15}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 15}, &span{20, 40}}),
//		},
//		{
//			name: "[20,40] + [10, 100]  = [10,100]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{10, 100}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 100}}),
//		},
//		{
//			name: "[20,40] + [10, 10]  = [10,10] [20 40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{10, 10}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 10}, &span{20, 40}}),
//		},
//		{
//			name: "[20,40] + [19, 19]  = [19 40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{19, 19}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{19, 40}}),
//		},
//		{
//			name: "[20,40] + [20, 20]  = [20 40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{20, 20}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//		},
//		{
//			name: "[20,40] + [21, 21]  = [20 40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{21, 21}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//		},
//		{
//			name: "[20,40] + [25, 25]  = [20 40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{25, 25}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//		},
//		{
//			name: "[20,40] + [39, 39]  = [20 40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{39, 39}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//		},
//		{
//			name: "[20,40] + [40, 40]  = [20 40]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{40, 40}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//		},
//		{
//			name: "[20,40] + [41, 41]  = [20 41]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{41, 41}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 41}}),
//		},
//		{
//			name: "[20,40] + [50, 50]  = [20 40] [50, 50] ",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}}),
//				intervalset.NewSet([]intervalset.Interval{&span{50, 50}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 40}, &span{50, 50}}),
//		},
//		/////////
//		{
//			name: "[20, 40] + [41,50]+ [50,60]+ [61,111] = [20,111]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 40}, &span{41, 50}}),
//				intervalset.NewSet([]intervalset.Interval{&span{50, 60}, &span{61, 111}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{20, 111}}),
//		},
//
//		{
//			name: "[10, 19] +  [20,20] = [10,20]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{10, 19}}),
//				intervalset.NewSet([]intervalset.Interval{&span{20, 20}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 20}}),
//		},
//		{
//			name: " [20,20] + [10, 19]  = [10,20]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 20}}),
//				intervalset.NewSet([]intervalset.Interval{&span{10, 19}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 20}}),
//		},
//		{
//			name: "[10, 19] + [21,30] + [20,20]  = [10,30]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{10, 19}, &span{21, 30}}),
//				intervalset.NewSet([]intervalset.Interval{&span{20, 20}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 30}}),
//		},
//		{
//			name: "[10, 15] +  [20,20]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{10, 15}}),
//				intervalset.NewSet([]intervalset.Interval{&span{20, 20}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 15}, &span{20, 20}}),
//		},
//		{
//			name: " [20,20] + [10, 15]",
//			args: args{
//				intervalset.NewSet([]intervalset.Interval{&span{20, 20}}),
//				intervalset.NewSet([]intervalset.Interval{&span{10, 15}}),
//			},
//			wantRes: intervalset.NewSet([]intervalset.Interval{&span{10, 15}, &span{20, 20}}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if gotRes := AddTest(tt.args.s1, tt.args.s2); !reflect.DeepEqual(gotRes.AllIntervals(), tt.wantRes.AllIntervals()) {
//				t.Errorf("Add() = %v, want %v", gotRes, tt.wantRes)
//			}
//		})
//	}
//}
