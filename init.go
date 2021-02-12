package main

import (
	"github.com/vincentfiestada/captainslog/v2"
	"github.com/vincentfiestada/captainslog/v2/format"
)

var log *captainslog.Logger

func init() {
	log = captainslog.NewLogger()
	log.Format = format.Minimal
}
