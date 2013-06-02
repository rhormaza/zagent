package zrouter

import (
    "zsearchlog"
	"zjson"
)

//var RouterMap = map[string] func(*zjson.JsonParams) (interface {}) {
var RouterMap = map[string] func(*zjson.JsonParams) (interface {}, zjson.JsonError) {
    "searchlog":  zsearchlog.Process,
}

