package main

import (
    "os"
    "net"
    "zutil"
    "zconfig"
    //"zjson"
    //"zrouter"
    "ztcpserver"
)

// Logger always first!
var log = zutil.SetupLogger("/tmp/zagent.log", 2)

//
// This function is intended to process to do the steps to process a file.
// The steps are:
// 1.- Launch ztcpserver, this server only handles request from the netwok, 
//     nothing else. It should be simple and lightweight
// 2.- Call zjson.Decode in order to decode a json object. All json decode/encode
//     logic goes here and only here!
// 3.- If decode fails, return the error back to the client 
//     with and known error code
// 4.  If decode is successful, make a call to the zrouter Map with the key
//     extracted from the JSON query. zrouter Map will redirect
//     and call the process.
// 5.- Return result of the method called from the router (zrouter.Map).
// 6.- Done.
//
//func process() {
//    // FIXME: load arguments from a config file!  
//    // This is a example how to load stuff from a config package
//    //log.Info("Args: %s and config:%s", os.Args, zconfig.LoadConfig("asas"))
//
//    // Init our server
//    s := ztcpserver.Server{ Cert: "zagent.pem", Prvkey:"zagent.pem",
//    Handlers: make(map[string] func (*[]byte, *[]byte, *net.Conn) error)}
//
//    // Server now listens from from this port
//    s.Run(":44443", false)
//
//    // We need to get data somehow from the server and put in a byte[]
//    // rawDataQuery :- s.getDataIfReadyToBeRead()
//
//    // pass byte[] to the router.
//    jsonReply, err :=  zjson.DecodeJson(rawDataQuery)
//    if err != nil {
//        // Return valid JSON error reply with error code as well.
//        // In short make a call to zjson.EncodeJsonFail()
//        //resultJson := zjson.EncodeJsonFail(result)
//    } else {
//
//        // Make the call to the router Map with data decoded
//        // Json.Method should have the "key" to the next function call
//        // Json.Params should have paramaters passed within the query
//        result, err := zrouter.RouterMap[jsonReply.method](jsonReply.Params)
//
//        if err != nil {
//            // Return valid JSON error reply with error code as well.
//            // In short make a call to zjson.EncodeJsonFail()
//            //resultJson := zjson.EncodeJsonFail(result)
//        } else {
//            // Encode data returned from requested method (from the router)
//            //resultJson := zjson.EncodeJsonSuccess(result)
//        }
//    }
//    // Write encoded back to the client. 
//    //s.WriteBackToClient(resultJson)
//
//    // Close the logger to clear the buffers!
//    log.Close()
//
//}

func main() {
    log.Info("Args: %s and config:%s", os.Args, zconfig.LoadConfig("asas"))
    s := ztcpserver.Server{ Cert: "zagent.pem", Prvkey:"zagent.pem",
                 Handlers: make(map[string] func (*[]byte, *[]byte, *net.Conn) error)}

    //s.AddHandler("3", agthandler.Test_handler)
    s.Run(":44443", false)

    /*
    jsonReply :=  zjson.DecodeJson(zjson.JsonBlobMap["searchlog_query"])

    //m := search.Process(jsonReply.Params)
    m, _ := zrouter.RouterMap["searchlog"](jsonReply.Params)

    log.Critical("jsonBlob:%s", zjson.EncodeJsonSuccess(m))
    */

    // we need to close the logger to clear the buffers!
    log.Close()
}
