package zjson

import (
    "encoding/json"
    "zutil"
)

// Logger always first!
var log = zutil.SetupLogger("/tmp/zagent.log", 2)

const (
    JSONRPC_VERSION = "2.0"
)


func DecodeJson(jsonBlob []byte) (jsonObj interface{}, err error) {
    log.Debug("Decoding data into Json")
    if err = json.Unmarshal(jsonBlob, &jsonObj);  err != nil {
        log.Error("Error %s decoding Json: %s in data", err, jsonBlob)
    }
    return
}


func EncodeJson(data interface{}) ([]byte, error){
    log.Debug("Encoding Json data: %s", data)
    jsonBlob, err := json.Marshal(data)
    if err != nil {
        log.Error("Error %s encoding data: %s in Json", err, data)
        return jsonBlob, err
    }
    return jsonBlob, nil
}

// This will encode a Json Successfully package, that is
// a Json blob with "result" property.
func EncodeJsonSuccess(data interface{}, id int64) ([]byte){
    log.Debug("Encoding a success Json data")

    var reply JsonReplySuccess
    reply.Jsonrpc = JSONRPC_VERSION
    reply.Id = id

    if _result, ok := data.(JsonResult); ok {
        reply.Result = &_result
    } else {
        log.Critical("Error asserting output data")
    }
    jsonBlob, err := json.Marshal(reply)
    if err != nil {
        log.Error("Error %s encoding data: %s in Json", err, data)
    }
    //return b, err
    return jsonBlob
}
