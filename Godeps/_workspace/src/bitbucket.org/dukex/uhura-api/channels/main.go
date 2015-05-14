package channels

import "github.com/op/go-logging"

var log = logging.MustGetLogger("example")

var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)
