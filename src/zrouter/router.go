package zrouter

import (
    "zsearchlog"
    "zjson"
    "zstatus"
)
// Deprecated?
type MethodMap map[string] func (map[string]interface{}) (interface {}, error)
// Deprecated?
type RouterMap map[string] MethodMap

// Deprecated?
var ZrouterMap RouterMap = RouterMap {
    "status": MethodMap {
        "info": zstatus.GetAgentInfo,
    },
}

// FIXME: This should be renamed!
type MethodMap2 map[string] func (*zjson.JsonParams) (interface {}, error)
type RouterMap2 map[string] MethodMap2

var Router RouterMap2 = RouterMap2 {
    // _search_ namespace
    "search": MethodMap2 {
        "log": zsearchlog.Process, // Method *must* follow defined signature
    },

    // _status_ namespace
    "status": MethodMap2 {
        "info": zstatus.Process, // Method *must* follow defined signature
    },

}



