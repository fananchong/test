package test7

import (
	"go_analysis_mutex_test/test7/d7"
	"time"
)

var tocache *d7.TimeoutCache

func F7() {
	tocache = d7.New(time.Minute*10, time.Minute*10)
}

func init() {
	F7()
}
