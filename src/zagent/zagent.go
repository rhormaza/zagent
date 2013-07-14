package main

import (
    "os"
    //"net"
    "zutil"
    "zconfig"
    //"zjson"
    //"zrouter"
    "ztcpserver"
    "zrouter"
)

// Logger always first!
var log = zutil.SetupLogger("/tmp/zagent.log", 2)

func main() {
    log.Info("Args: %s and config:%s", os.Args, zconfig.LoadConfig("asas"))
    s := ztcpserver.Server { Cert: "certs/zagent_server.pem", Prvkey:"certs/zagent_server.key",
                 Router: &zrouter.ZrouterMap }

    //s.AddHandler("3", agthandler.Test_handler)
    s.Run(":44443", false)

    // we need to close the logger to clear the buffers!
    log.Close()
}
