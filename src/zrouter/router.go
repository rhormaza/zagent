package zrouter

import (
    "zsearchlog"
	"zjson"
)
var RouterMap = map[string] func(*zjson.JsonParams) (interface {}) {
    "searchlog":  zsearchlog.Process,
}

