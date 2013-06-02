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

func main() {
    log.Info("Args: %s and config:%s", os.Args, zconfig.LoadConfig("asas"))

    jsonReply :=  zjson.DecodeJson(zjson.JsonBlobMap["searchlog_query"])

    //m := search.Process(jsonReply.Params)
    m, _ := zrouter.RouterMap["searchlog"](jsonReply.Params)

    log.Critical("jsonBlob:%s", zjson.EncodeJsonSuccess(m))

    // we need to close the logger to clear the buffers!
    log.Close()
}
