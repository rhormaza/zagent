package main

import (
    "os"
    "util"
//    "search"
    "config"
    "jsonrouter"
    "router"
)

// Logger always first!
var log = util.SetupLogger("/tmp/zagent.log", 2)

func see() {
    log.Debug("Calling see")
}

//func main() {
//    mSee := see
//    mSee()
//    see()
//    log.Close()
//}

func main() {
    log.Info("Args: %s and config:%s", os.Args, config.LoadConfig("asas"))

    jsonReply :=  jsonrouter.DecodeJson(jsonrouter.JsonBlob)

    //m := search.Process(jsonReply.Params)
    m := router.RouterMap["searchlog"](jsonReply.Params)
    log.Debug("%s", m)
    log.Critical("jsonBlob:%s", jsonrouter.EncodeJson(m))

    //we need to close the logger to clear the buffers!
    log.Close()
}
