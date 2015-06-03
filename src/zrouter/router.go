package zrouter

import (
    "zjson"
    "zsearch"
    "zstatus"
)

// FIXME: This should be renamed!
type MethodMap map[string] func (*zjson.JsonParams) (interface {}, error)
type RouterMap map[string] MethodMap

var Router RouterMap = RouterMap {
    // _search_ namespace
    "search": MethodMap {
        "log": zsearch.Log, // Method *must* follow defined signature
    },
    "log": MethodMap {
        "extract": zsearch.Extract, // Method *must* follow defined signature
    },

    // _status_ namespace
    "status": MethodMap {
        "info": zstatus.Info, // Method *must* follow defined signature
    },

}



