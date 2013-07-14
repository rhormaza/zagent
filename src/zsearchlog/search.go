package zsearchlog

import (
    "os"
    "bytes"
    "regexp"
    "strings"
    "crypto/sha256"
    "encoding/hex"
    "zutil"
    "zjson"
)

// Logger always first!
var log = zutil.GetLogger()



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
func doMatch(lineNumber int64, line string, jsonParams *zjson.JsonParams, hitsMap *map[string] []zjson.SearchLogHit) (hit bool) {
    // Let's check our JSON params first!
    if jsonParams == nil {
        return false //FIXME: add a proper error code!
    }

    patternMatrix := (*jsonParams)["pattern"]
    if row, ok := patternMatrix.([]interface{}); ok {
        for _, v := range row {
            
            errClrPattern := make([]string, 2)
            // At this point col has an array like ["ERROR_PATTERN", "CLEAR_PATTERN"]
            if col, ok := v.([]interface{}); ok {
                // Checking for the "Error" pattern first
                if pattern, ok := col[0].(string); ok {
                    errClrPattern[0] = pattern
                } else {
                    log.Error("Json assertion failed")
                    return false
                }

                // Checking for the "Clear" pattern first
                if pattern, ok := col[1].(string); ok {
                    errClrPattern[1] = pattern
                } else {
                    log.Error("Json assertion failed")
                    return false
                }
                
                patternKey := strings.Join(errClrPattern, "<==>")
                // Do matching
                if m, _ := regexp.MatchString(errClrPattern[0], line); m {
                    (*hitsMap)[patternKey] = append((*hitsMap)[patternKey], zjson.SearchLogHit{line, lineNumber, errClrPattern[0], "error"})
                }
                if m, _ := regexp.MatchString(errClrPattern[1], line); m {
                    (*hitsMap)[patternKey] = append((*hitsMap)[patternKey], zjson.SearchLogHit{line, lineNumber, errClrPattern[1], "error"})
                }
                //if m, _ := regexp.MatchString(clearPattern, line); m {
                //    (*hitsMap)[clearPattern] = append((*hitsMap)[clearPattern], zjson.SearchLogHit{line, lineNumber, clearPattern, "clear"})
                //}
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

func ProcessChunk(buffer []byte, beginPos int64, length int64, jsonParams *zjson.JsonParams) (hitsMap map[string] []zjson.SearchLogHit) {
    /*
    TODO: Check for '\r' char in case of a windows box
          Check EOL during the split!
    */

    // It holds the actual result!
    hitsMap = make(map[string] []zjson.SearchLogHit)

    eol := []byte{'\n'}
    _beginPos := int(beginPos)
    for i, v := range bytes.Split(buffer, eol) {
        lineNumber := i + _beginPos + 1 // line starts from 1, whereas i starts from 0
        doMatch(int64(lineNumber), string(v), jsonParams, &hitsMap)
    }
    log.Debug("buffer len: %d, buffer hash(sha):%s", len(buffer), hex.EncodeToString(CalcSha256(buffer)))

    return
}

// Wrapper function that does all the processing. Function
// gets at valid Json Params object to get input parameters.
//
// Main process that walk through file searching for patterns.
//  - A SHA256 hash is calculated for the first 100 bytes
//    of chunk text
//  - If the provided hash matches calculated hash search 
//    is skipped
// 
// TODO: return a valid JSON error if something goes wrong!!!
func Process(jsonParams *zjson.JsonParams) (interface{}, error) {//(zjson.JsonResult, zjson.JsonError) { //(interface {}) {
    log.Debug("Executing Process()")

    var HASH_LIMIT int64 = 100 //only hash first 100 bytes

    filename := (*jsonParams)["filename"]
    filenameStr, _ := filename.(string)

    // Ugly and stupid hack, because the passing 
    // int's couldn't be assert!
    // Hence, the horrible workaround was to pass 
    // begin_pos and end_pos values as strings :-(
    //___endPos := (*jsonParams)["endpos"]
    _endPos, _ := (*jsonParams)["endpos"].(float64)
    endPos := int64(_endPos)

    _beginPos, _ := (*jsonParams)["beginpos"].(float64)
    beginPos := int64(_beginPos)

    hash := (*jsonParams)["hash"]
    hashStr, _ := hash.(string)

    mdStr, _ := CalcSha256FromChunk(filenameStr, beginPos, HASH_LIMIT)
    result := new(zjson.SearchLogPattern) 
    if hashStr == mdStr {
        log.Debug("Hashes match. Given: %s Calculated: %s", hashStr, mdStr)

        // This is the return value if all is successfull
        result.Filename = filenameStr
        result.Hash = mdStr
        result.BeginPos = beginPos
        result.EndPos = endPos
        result.Hits = nil
    } else {
        log.Debug("Hashes differ. Given: %s Calculated: %s", hashStr, mdStr)

        // If hashes are different we need to search the whole file
        buf, _ := ReadChunk(filenameStr, 0, endPos)
        hits := ProcessChunk(buf, beginPos, endPos, jsonParams)

        // This is the return value if all is successfull
        result.Filename = filenameStr
        result.Hash = mdStr
        result.BeginPos = beginPos
        result.EndPos = 4611686018427387900
        result.Hits = hits
    }


    return result, nil // FIXME: return nil as en error is bad!
}

