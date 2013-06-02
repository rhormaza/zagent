package zsearchlog

import (
    "os"
    "bytes"
    "regexp"
    "strconv"
    "crypto/sha256"
    "encoding/hex"
    "zutil"
    "zjson"
//    "reflect"
)

// Logger always first!
var log = zutil.GetLogger()


// It holds the actual result!
var JsonResult = make(map[string] []zjson.SearchLogHit)

//func readChunk(path string, beginPos int64, length int64) (buf *bytes.Buffer, err error) {
func ReadChunk(path string, beginPos int64, length int64) (buffer []byte, err error) {
    
    log.Debug("Reading file: %s from position: %d until: %d", path, beginPos, beginPos + length)
    var file *os.File

    if file, err = os.Open(path); err != nil {
        log.Critical("Error reading the file %s", path)
        return
    }
    defer file.Close()

    // get the file size
    stat, err := file.Stat()
    
    if length == 0 || length + beginPos > stat.Size() { 
        length = stat.Size() 
    }

    // go to right position
    file.Seek(beginPos, os.SEEK_SET)

    buffer = make([]byte, length - beginPos)

    // read the chunk into buffer
    _, err = file.Read(buffer)
    
    // EOF error will be propagate just because it will be needed
    // Hence, we comment this. It is only left for information.
    //
    //if err == io.EOF {
    //    err = nil
    //}
    return
}

func CalcSha256(b []byte) (digest []byte) {
    
    hash := sha256.New()
    hash.Write(b)
    digest = hash.Sum(nil)

    return 
}

func CalcSha256FromChunk(path string, beginPos int64, length int64) (digestStr string, err error) {
    
    buf, err := ReadChunk(path, beginPos, length)

    digest := CalcSha256(buf)
    digestStr = hex.EncodeToString(digest)
    
    return 
}

//func doMatch(lineNumber int64, line string, jsonParams *map[string]interface{}) (hit bool) {
func doMatch(lineNumber int64, line string, jsonParams *zjson.JsonParams) (hit bool) {
    // Let's check our JSON params first!
    if jsonParams == nil {
        return false //FIXME: add a proper error code!
    }

    patternMatrix := (*jsonParams)["pattern"]
    if row, ok := patternMatrix.([]interface{}); ok {
        for _, v := range row {
            if col, ok := v.([]interface{}); ok {
                if errPattern, ok := col[0].(string); ok {
                    if m, _ := regexp.MatchString(errPattern, line); m {
                        JsonResult[errPattern] = append(JsonResult[errPattern], zjson.SearchLogHit{line, lineNumber, 9999999})
                        // This could also be done with maps, to do so you need to define the a global variable
                        // like : var _JsonResult = make(map[string] []map[string]string)
                        //
                        // And then put something like below in this very place! I prefer use *struct* though
                        //_JsonResult[errPattern] = append(_JsonResult[errPattern], 
                        //    map[string]string {
                        //        "LineText" :  line,
                        //        "LineNumber" : strconv.FormatInt(lineNumber, 10),  // See strconv docs
                        //        "LineBegin" : "999", // FIXME: for now this should a real value
                        //    })
                    }
                } else {
                    log.Error("Json assertion failed")
                    return false
                }
            } else {
                log.Error("Json assertion failed")
                return false
            }

        }
    } else {
        log.Error("Json assertion failed")
        return false
    }
    return true

}

//func ProcessChunk(buffer []byte, beginPos int, length int, jsonParams *map[string]interface{}) {
func ProcessChunk(buffer []byte, beginPos int64, length int64, jsonParams *zjson.JsonParams) {
    /*
    TODO: Check for '\r' char in case of a windows box
          Check EOL during the split!
    */

    eol := []byte{'\n'}
    _beginPos := int(beginPos)
    for i, v := range bytes.Split(buffer, eol) {
        lineNumber := i + _beginPos + 1 // line starts from 1, whereas i starts from 0
        doMatch(int64(lineNumber), string(v), jsonParams)
    }
    log.Debug("buffer len: %d, buffer hash(sha):%s", len(buffer), hex.EncodeToString(CalcSha256(buffer)))
}

// Wrapper function that does all the processing. Function
// gets at valid Json Params object to get input parameters.
//
// TODO: return a valid JSON error if something goes wrong!!!
func Process(jsonParams *zjson.JsonParams) (interface{}, zjson.JsonError) {//(zjson.JsonResult, zjson.JsonError) { //(interface {}) {
    log.Debug("Doing some Process")
 
    filename := (*jsonParams)["filename"]
    filenameStr, _ := filename.(string)

    // Ugly and stupid hack, because the passing 
    // int's couldn't be assert!
    // Hence, the horrible workaround was to pass 
    // begin_pos and end_pos values as strings :-(
    ___endPos := (*jsonParams)["end_pos"]
    __endPos, _ := ___endPos.(string)
    _endPos, _ := strconv.Atoi(__endPos)
    endPos := int64(_endPos)

    ___beginPos := (*jsonParams)["begin_pos"]
    __beginPos, _ := ___beginPos.(string)
    _beginPos, _ := strconv.Atoi(__beginPos)
    beginPos := int64(_beginPos)

    hash := (*jsonParams)["hash"]
    hashStr, _ := hash.(string)

    mdStr, _ := CalcSha256FromChunk(filenameStr, beginPos, endPos)
    log.Debug("Provided hash: %s and calculated hash: %s", hashStr, mdStr)
    log.Debug("begin_pos: %d, end_pos: %d", beginPos, endPos)

    buf, _ := ReadChunk(filenameStr, beginPos, endPos)
    ProcessChunk(buf, beginPos, endPos, jsonParams)

    return JsonResult, nil // FIXME: return nil as en error is bad!
}

