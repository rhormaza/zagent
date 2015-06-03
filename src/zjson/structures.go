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
//type JsonParams map[string][]interface{}

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

//type JsonResult map[string]interface{}
type JsonResult interface{}
type JsonError interface{}

//type JsonError struct{
//  Errcode int64   `json:"code"`
//  Errmsg  string
//}

type JsonReplySuccess struct {
    Jsonrpc string          `json:"jsonrpc"`    // JSON-RPC version
    Result  *JsonResult     `json:"result"`     // TODO: explain!
    Id      int64           `json:"id"`         // Id of the JSON query
}

type JsonReplyError struct {
    Jsonrpc string          `json:"jsonrpc"`   // JSON-RPC version
    Error   *JsonError      `json:"error"`    // TODO: explain!
    Id      int64           `json:"id"`        // Id of the JSON query
}

type JsonObject struct {
    Jsonrpc string                      // JSON-RPC version
    Id      int64                       // Id of the JSON query
    Method  *string                      // Method to be executed
    Params  *JsonParams                 // Arguments/Parameters of the method provided (It could be empty)
}

// As convention we will use the a Preffix of the 
// package that will use this. 

// This might be merged to JsonError directly
type Err struct {
    Errcode int64   `json:"code"`
    Errmsg  string  `json:"message"`
}

// Struct use to return a searchlog hit!
type SearchLogHit  struct {
    LineText    string  `json:"linetext"`
    LineNum     int64   `json:"linenum"`
    Pattern     string  `json:"pattern"`
    Type        string  `json:"type"` // clear or error
    //LineBegin   int64
    //more metadata to hold?
}

// Struct use to return a searchlog hit!
// FIXME: delme?
type SearchLogPattern struct {
    Filename    string                      `json:"filename"`
    Hash        string                      `json:"hash"`
    BeginPos    int64                       `json:"beginpos"` // This is around 4 exa bytes 10^18
    EndPos      int64                       `json:"endpos"`
    Hits        map[string][]SearchLogHit   `json:"hits"`
}

// Agent status struct
// Note that the name should use NamespaceMethod
type StatusInfo struct { 
    Version string  `json:"version"` 
    Os      string  `json:"os"`
    Port    int64   `json:"port"`
}


type ExtractLog struct {
    BeginPos    int64       `json:"beginpos"`
    EndPos      int64       `json:"endpos"`
    Filename    string      `json:"filename"`
    Lines       []string    `json:"lines"`
}
