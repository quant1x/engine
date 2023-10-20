package base

import "testing"

//func TestBatchRealtime(t *testing.T) {
//	//BatchKLineWideRealtime([]string{"sh000001", "sh000905", "sz399001", "sz399006", "sh600600", "sz002528"})
//	_ = BatchKLineWideRealtime([]string{"sz000966"})
//}

func TestBatchRealtimeBasicKLine(t *testing.T) {
	_ = BatchRealtimeBasicKLine([]string{"sh000001"})
}
