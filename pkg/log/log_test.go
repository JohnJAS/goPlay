package log

import "testing"

func TestLogLevelItoa(t *testing.T){
	if LogLevelItoa(DEBUG) == "" {
		t.Error(`LogLevelAtoi(DEBUG) != DUBUG`)
	}
}