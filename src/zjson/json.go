package zjson

import (
    "encoding/json"
    "zutil"
)
// Logger always first!
var log = zutil.SetupLogger("/tmp/zagent.log", 2)

var JsonBlob = []byte(`{
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
        "end_pos": "999999999999999999", 
        "begin_pos": "0"
    },
    "id": 2
}`)

// This is the main JSON struct, this struct follows JSON-RPC 2.0 standard
//
// {
//     "jsonrpc": "2.0",
//     "method": "METHOD_NAME",
//     "params": [
//          "SOME_PARAMETER",
//          "ANOTHER_PARAMTER",
//          100
//     ],
//     "id": 1
// }
//
// Struct fields start with Uppercase for JSON Marshalling
// All JSON query and procedure should comply with this format.
type JsonParams map[string]interface{}
type JsonQuery struct {
    Jsonrpc string                      // JSON-RPC version
    Method  string                      // Method to be executed
    Params  *JsonParams                 // Arguments/Parameters of the method provided (It could be empty)
    //Params  *map[string]interface{}   // Arguments/Parameters of the method provided (It could be empty)
    Id      int64                       // Id of the JSON query
}


// Succesfull result struct, this struct must follow the format
// {
//     "jsonrpc": "2.0",
//     "result": {
//         "SOME_PROPERTY": "foo",
//         "ANOTHER_PROPERTY": "bar"
//     },
//     "id": 2
// }
//
// An error must follow the format
// {
//     "jsonrpc": "2.0",
//     "error": {
//         "code": -10001,
//         "message": "Method failed, reason: ......"
//     },
//     "id": 2
// }
//
type JsonResult map[string]interface{}
type JsonError map[string]interface{}
type JsonReply struct {
    Jsonrpc string                      // JSON-RPC version
    Result  *JsonResult                 // TODO: explain!
    Error   *JsonError                  // TODO: explain!
    //Result  *map[string]interface{}     // TODO: explain!
    //Error   *map[string]interface{}     // TODO: explain!
    Id      int64                       // Id of the JSON query
}
type JsonObject struct {
    Jsonrpc string                      // JSON-RPC version
    Id      int64                       // Id of the JSON query

    Method  *string                      // Method to be executed
    Params  *JsonParams                 // Arguments/Parameters of the method provided (It could be empty)
    
    Result  *JsonResult                 // TODO: explain!
    Error   *JsonError                  // TODO: explain!
    //Result  *map[string]interface{}     // TODO: explain!
    //Error   *map[string]interface{}     // TODO: explain!
}


func DecodeJson(jsonBlob []byte) (jsonObj JsonObject) {
    log.Debug("Decoding data into Json")
    if err := json.Unmarshal(jsonBlob, &jsonObj);  err != nil {
        log.Error("Error %s decoding Json: %s in data", err, jsonBlob)
    }
    return
}

func EncodeJson(data interface{}) ([]byte){
    log.Debug("Encoding Json data: %s", data)
    jsonBlob, err := json.Marshal(data)
    if err != nil {
        log.Error("Error %s encoding data: %s in Json", err, data)
    }
    //return b, err
    return jsonBlob
}


//
//func main() {
//    log.Info("Args: %s and config:%s", os.Args, config.LoadConfig("asas"))
//
//    var jsonObj JsonQuery
//    err := json.Unmarshal(jsonBlob, &jsonObj)
//    if err != nil {
//        log.Critical("error: %s", err)
//    }
//
//
//    last, _ := strconv.Atoi(os.Args[1])
//    buf, err := search.ReadChunk(jsonObj.Params.Filename, 0, int64(last))
//    search.ProcessChunk(buf, 0, last, &jsonObj)
//    mdStr, _ := search.CalcSha256FromChunk("/tmp/foo.txt", 0, int64(last))
//
//    log.Info(">>>>>>>>> %s", mdStr)
//    for k, v := range PatternHits {
//        log.Info("k: %s and v: %s", k, v)
//    }
//
//    //we need to close the logger to clear the buffers!
//    log.Close()
//}
