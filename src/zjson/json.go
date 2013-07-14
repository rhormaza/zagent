package zjson

import (
    "encoding/json"
    "zutil"
)

// Logger always first!
var log = zutil.SetupLogger("/tmp/zagent.log", 2)

// FIXME: move?
const (
    JSONRPC_VERSION = "2.0"
)


// Deprecated?
func DecodeJson(jsonBlob []byte) (jsonObj interface{}, err error) {
    log.Debug("Decoding data into Json")
    if err = json.Unmarshal(jsonBlob, &jsonObj);  err != nil {
        log.Error("Error %s decoding Json: %s in data", err, jsonBlob)
    }
    return
}

// It should be renamed to DEcodeJson()
// This method will return unmarshalled object from the json blob received 
func DecodeJson2(jsonBlob []byte) (jsonObj JsonObject, err error) {
    log.Debug("Decoding data into Json")
    if err = json.Unmarshal(jsonBlob, &jsonObj);  err != nil {
        log.Error("Error %s decoding Json: %s in data", err, jsonBlob)
    }
    return
}

func EncodeJson(data interface{}) ([]byte, error){
    log.Debug("Encoding Json data: %s", data)
    //jsonBlob, err := json.Marshal(data)
    jsonBlob, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        log.Error("Error %s encoding data: %s in Json", err, data)
        return jsonBlob, err
    }
    return jsonBlob, nil
}

// This will encode a Json Successfully package, that is
// a Json blob with "result" property.
func EncodeJsonSuccess(data interface{}, id int64) ([]byte, error){
    log.Debug("Encoding a success Json reply")

    var reply JsonReplySuccess
    reply.Jsonrpc = JSONRPC_VERSION
    reply.Id = id

    if _result, ok := data.(JsonResult); ok {
        reply.Result = &_result
    } else {
        log.Critical("Error asserting output data")
    }
    //jsonBlob, err := json.Marshal(reply)
    jsonBlob, err := json.MarshalIndent(reply, "", "  ")
    if err != nil {
        log.Error("Error %s encoding data: %s in Json", err, data)
    }
    return jsonBlob, err
}

func EncodeJsonFail(data interface{}, id int64) ([]byte, error){
    log.Debug("Encoding an error Json reply")

    var reply JsonReplyError
    reply.Jsonrpc = JSONRPC_VERSION
    reply.Id = id

    if _error, ok := data.(JsonError); ok {
        reply.Error = &_error
    } else {
        log.Critical("Error asserting output data")
    }
    //jsonBlob, err := json.Marshal(reply)
    jsonBlob, err := json.MarshalIndent(reply, "", "  ")
    if err != nil {
        log.Error("Error %s encoding iFIXME", err)
    }
    return jsonBlob, err
}
