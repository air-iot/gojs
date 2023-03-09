package api

import (
	"time"
)

type Lib struct{}

func NewLib() *Lib {
	return &Lib{}
}

func (*Lib) SleepMill(t int64) {
	time.Sleep(time.Millisecond * time.Duration(t))
}
