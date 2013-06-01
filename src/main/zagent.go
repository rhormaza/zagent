package main

import (
    "os"
    "strconv"
    "encoding/json"
////////////
    "util"
    "search"
    "config"
)
// Logger always first!
var log = util.SetupLogger("/tmp/zagent.log", 2)

var jsonBlob = []byte(`{
    "jsonrpc": "2.0",
    "method": "get_error",
    "params": {
        "pattern": [
        [
        "ERROR PATTERN",
        "CLEAR PATTERN"
        ],
        [
        ".*stopped.*",
        ".*started.*"
        ],
        [
        ".*hello.*",
        ".*bye.*"
        ],
        [
        ".*hola.*",
        ".*chao.*"
        ],
        [
        ".*foo.*",
        ".*bar.*"
        ]
        ],
        "filename": "/tmp/foo.txt",
        "hash": "a SHA-256 hash",
        "begin_pos": 0,
        "end_pos": 1125899906842624 
    },
    "id": 2
}`)

//type JsonObject struct {
//    /*Fields start with Uppercase for JSON Marshalling */
//    Jsonrpc string
//    Method  string
//    Params  JsonParams
//    Id      int64
//}
//
//type JsonParams struct {
//    Pattern     [][]string
//    Filename    string
//    Hash        string
//    Begin_pos   int64
//    End_pos     int64
//}
//
//type Hit struct {
//    lineText    string
//    lineNumber  int64
//    lineBegin   int64
//    //more metadata to hold?
//}
//
//type Pattern struct {
//    pattern string
//    hits    []Hit
//}
//
//type Query struct {
//    pattern     string
//    filename    string
//    hash        string
//    beginPos    int64
//    endPos      int64
//}
//
//var hitSlices []Hit
//var PatternHits = make(map[string] []Hit)
type JsonObject search.JsonObject
type JsonParams search.JsonParams
var PatternHits = make(map[string] []search.Hit)

func main() {
    log.Info("Args: %s and config:%s", os.Args, config.LoadConfig("asas"))

    var jsonObj search.JsonObject
    err := json.Unmarshal(jsonBlob, &jsonObj)
    if err != nil {
        log.Critical("error: %s", err)
    }


    last, _ := strconv.Atoi(os.Args[1])
    buf, err := search.ReadChunk(jsonObj.Params.Filename, 0, int64(last))
    search.ProcessChunk(buf, 0, last, &jsonObj)
    mdStr, _ := search.CalcSha256FromChunk("/tmp/foo.txt", 0, int64(last))

    log.Info(">>>>>>>>> %s", mdStr)
    for k, v := range PatternHits {
        log.Info("k: %s and v: %s", k, v)
    }

    //we need to close the logger to clear the buffers!
    log.Close()
}
