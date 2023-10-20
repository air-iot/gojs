package log

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	l := NewLogger(SetGroup("testG"), SetModule("testM"), SetKey("testK"), SetProject("testP"))
	l.Debug(1)
	l.Info(1)
	l.Warn(1)
	l.Error(1)
}
