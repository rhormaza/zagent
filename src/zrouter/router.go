package zrouter

import (
  //"zsearchlog"
	"zjson"
  "zstatus"
)


type RouterMap struct {
  methodMap map[string] map[string] func(*zjson.JsonParams) (interface {}, zjson.JsonError)
}


ZrouterMap := RouterMap{
                         {
                           "status":
                                 {
                                   "info": zstatus.GetAgentInfo,
                                 },
                         },
                       }



