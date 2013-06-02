package main

import (
    "os"
    "zutil"
    "zconfig"
    "zjson"
    "zrouter"
)

// Logger always first!
var log = zutil.SetupLogger("/tmp/zagent.log", 2)

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
    log.Info("Args: %s and config:%s", os.Args, zconfig.LoadConfig("asas"))

    jsonReply :=  zjson.DecodeJson(zjson.JsonBlobMap["searchlog_query"])

    //m := search.Process(jsonReply.Params)
    m := zrouter.RouterMap["searchlog"](jsonReply.Params)

    log.Debug("%s", m)
    log.Critical("jsonBlob:%s", zjson.EncodeJson(m))

    //we need to close the logger to clear the buffers!
    log.Close()
}
