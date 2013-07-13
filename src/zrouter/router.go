package zrouter

import (
  //"zsearchlog"
	//"zjson"
  "zstatus"
)

//type MethodMap map[string] func (*zjson.JsonParams) (interface {}, error)
type MethodMap map[string] func (map[string]interface{}) (interface {}, error)
type RouterMap map[string] MethodMap

var ZrouterMap RouterMap = RouterMap {
                "status": MethodMap{
                            "info": zstatus.GetAgentInfo,
                          },
}



