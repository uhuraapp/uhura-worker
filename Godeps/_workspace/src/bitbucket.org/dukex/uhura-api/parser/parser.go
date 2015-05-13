package parser

import (
	"github.com/op/go-logging"
)

const (
	ITUNES_EXT = "http://www.itunes.com/dtds/podcast-1.0.dtd"
)

var log = logging.MustGetLogger("example")

var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)
