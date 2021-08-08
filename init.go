package main

import (
	"vincent.click/pkg/captainslog/v2"
	"vincent.click/pkg/captainslog/v2/format"
)

var log *captainslog.Logger

func init() {
	log = captainslog.NewLogger()
	log.Format = format.Minimal
}
