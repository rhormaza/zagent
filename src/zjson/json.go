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

// This will encode a Json Successfully package, that is
// a Json blob with "result" property.
func EncodeJsonSuccess(data interface{}) ([]byte){
    log.Debug("Encoding Json data: %s", data)

    var reply JsonReplySuccess
    reply.Jsonrpc = JSONRPC_VERSION
    reply.Id = 1
    // Asserting struct
    if result, ok := data.(map[string][]SearchLogHit); ok {
        log.Info("Attaching result to Json reply. Result is %s", result)
        //
        // Ugly thing to convert actual result "a particular" data structure into a JsonResult Type
        // see: http://golang.org/doc/faq#convert_slice_of_interface
        a := make(JsonResult)
        for i, v := range result {
            a[i] = v
        }
        reply.Result = &a
    } else {
        log.Critical("Error asserting output data")
    }
    // This is the same than above in case of using maps instead of structs to return JsonResult
    //if result, ok := data.(map[string][]map[string]string); ok {
    //    log.Info("Attaching result to Json reply. Result is %s", result)
    //    // Ugly thing to convert actual result "a particular" data structure in a JsonResult Type
    //    a := make(JsonResult)
    //    for i, v := range result {
    //        a[i] = v
    //    }
    //    reply.Result = &a
    //} else {
    //    log.Critical("Error asserting output data")
    //}

    jsonBlob, err := json.Marshal(reply)
    if err != nil {
        log.Error("Error %s encoding data: %s in Json", err, data)
    }
    //return b, err
    return jsonBlob
}

