package zjson

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

type JsonReplySuccess struct {
    Jsonrpc string                      // JSON-RPC version
    Result  *JsonResult                 // TODO: explain!
    Id      int64                       // Id of the JSON query
}

type JsonReplyError struct {
    Jsonrpc string                      // JSON-RPC version
    Error   *JsonError                  // TODO: explain!
    Id      int64                       // Id of the JSON query
}

type JsonReply struct {
    Jsonrpc string                      // JSON-RPC version
    Result  *JsonResult                 // TODO: explain!
    Error   *JsonError                  // TODO: explain!
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

// As convention we will use the a Preffix of the 
// package that will use this. 

// Struct use to return a searchlog hit!
type SearchLogHit  struct {
    LineText    string
    LineNumber  int64
    LineBegin   int64
    //more metadata to hold?
}

// Struct use to return a searchlog hit!
// FIXME: delme?
type SearchLogPattern struct {
    Pattern string
    Hits    []SearchLogHit
}

// Struct use to return a searchlog hit!
// FIXME: delme?
type SearchLogQuery struct {
    pattern     string
    filename    string
    hash        string
    beginPos    int64
    endPos      int64
}
