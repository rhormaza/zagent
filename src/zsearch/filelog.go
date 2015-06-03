package zsearch

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

type LogFileMetaData struct {
    ReadEndPos      int64
    Name            string
    ReadDataLeft    int64
}

var logFileStats *LogFileMetaData

const CHUNK_LIMIT = 4096 // 4K

func calcSha256(b []byte) (digest []byte) {

    hash := sha256.New()
    hash.Write(b)
    digest = hash.Sum(nil)

    return
}

func calcSha256FromChunk(path string, beginPos int64, chunkLength int64) (digestStr string, err error) {

    buf, err := readChunk(path, beginPos, chunkLength)

    digest := calcSha256(buf)
    digestStr = hex.EncodeToString(digest)

    return
}

func calcSha256FromChunkForward(path string, beginPos int64, chunkLength int64) (digestStr string, err error) {
    return calcSha256FromChunk(path, beginPos, chunkLength)
}

func calcSha256FromChunkBackward(path string, beginPos int64, chunkLength int64) (digestStr string, err error) {

    if beginPos == 0 {
        // FIXME(raul): this is a hack, find a better solution!
        return "__zero__", nil
    }

    newBeginPos := beginPos - chunkLength

    if newBeginPos < 0 {
        return calcSha256FromChunk(path, 0, chunkLength)
    } else {
        return calcSha256FromChunk(path, newBeginPos, chunkLength)
    }

}

//
// This will read a chunk of "length" bytes of the file "path" starting from "beginPos"
//
func readChunk(path string, beginPos int64, length int64) (buffer []byte, err error) {

    log.Debug("Reading a %d bytes chunk from %s position in %s",
        length, beginPos, path)

    var file *os.File

    // Opening the file
    if file, err = os.Open(path); err != nil {
        log.Critical("Error reading the file %s", path)
        return
    }
    defer file.Close()

    // Get the file size
    stat, err := file.Stat()
    fileSize := stat.Size()

    if length > CHUNK_LIMIT {
        length = CHUNK_LIMIT
    }

    var bufferSize int64

    if length == 0 || length + beginPos > fileSize {
        bufferSize = fileSize - beginPos

        // record the size of the file.
        logFileStats.ReadEndPos = fileSize
        logFileStats.ReadDataLeft = 0

    } else {

        // record the size of the file.
        bufferSize = length

        logFileStats.ReadEndPos = beginPos + length
        logFileStats.ReadDataLeft = fileSize - logFileStats.ReadEndPos

    }

    // go to right position
    file.Seek(beginPos, os.SEEK_SET)

    buffer = make([]byte, bufferSize)

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


    // We don't care about EOL differences ('\r\n' and '\n')
    // That is due to the fact that we are grepping a line
    // so having an extra character like '\r' is irrelevant.
    eol := []byte{'\n'}

    eolLastIndex := bytes.LastIndex(buffer, eol) + 1
    bufferLen := len(buffer)
    delta := int64(bufferLen - eolLastIndex)
    if delta > 1 {
        // EOL was not at the end of the buffer so 
        // we could be reading in the middle of a line
        // Hence, we discard these chars and modify 
        // counters for position.
        buffer = buffer[:eolLastIndex]
        logFileStats.ReadEndPos -= delta
        logFileStats.ReadDataLeft = logFileStats.ReadDataLeft + delta
        log.Debug("Buffer was returned %s bytes to ensure reading full lines", delta)

    }

    lines := bytes.Split(buffer, eol)
    _beginPos := int(beginPos)
    for i, v := range lines {
        lineNumber := i + _beginPos + 1 // line starts from 1, whereas i starts from 0
        doMatch(int64(lineNumber), string(v), jsonParams, &hitsMap)
    }
    log.Debug("buffer len: %d, buffer hash(sha):%s", len(buffer), hex.EncodeToString(calcSha256(buffer)))

    return
}

const (
    HASH_NO_APPLY int64= iota
    HASH_DIFFER
    HASH_MATCH
)

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
func Log(jsonParams *zjson.JsonParams) (interface{}, error) {//(zjson.JsonResult, zjson.JsonError) { //(interface {}) {
    log.Debug("Initiating the execution of Log()")

    var HASH_LIMIT int64 = 100 //only hash first 100 bytes
    logFileStats = new(LogFileMetaData)


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
    recvHashStr, _ := hash.(string)

    searchMethodCase := HASH_NO_APPLY // by default we assume we have to search the whole file
    calcHashStr, _ := calcSha256FromChunkBackward(filenameStr, beginPos, HASH_LIMIT)
    if recvHashStr == calcHashStr {
        searchMethodCase = HASH_MATCH
    } else {
        searchMethodCase = HASH_DIFFER
    }

    result := new(zjson.SearchLogPattern)

    switch searchMethodCase {
    case HASH_MATCH:
        log.Debug("Hashes MATCH!. Given: %s Calculated: %s", recvHashStr, calcHashStr)

        length := endPos - beginPos
        buf, _ := readChunk(filenameStr, beginPos, length)
        if len(buf) == 0 {
            // If we are here, that most likely means we found the end of the file.
            result.Filename = filenameStr
            result.Hash = recvHashStr // send same hash back!
            result.BeginPos = beginPos - HASH_LIMIT // go back
            result.EndPos = beginPos // rewind the beginPos pointer
            // Null, None or nil means NO hits!
            result.Hits = nil
        } else {
            hits := ProcessChunk(buf, beginPos, length, jsonParams)

            // Calc last "HASH_LIMIT" bytes hash!
            retHashStr, _ := calcSha256FromChunkBackward(filenameStr, logFileStats.ReadEndPos, HASH_LIMIT)

            // This is the return value if all is successfull
            result.Filename = filenameStr
            result.Hash = retHashStr
            result.BeginPos = beginPos
            result.EndPos = endPos
            result.Hits = hits
        }

    // HASH_NO_APPLY would probably never be truth, just here for preventing issues.
    case HASH_NO_APPLY, HASH_DIFFER:
        log.Debug("Hashes MISMATCH!. Given: %s Calculated: %s", recvHashStr, calcHashStr)

        length := endPos
        // If hashes are different we need to search the whole file
        buf, _ := readChunk(filenameStr, 0, length)
        // FIXME(raul): len(buf) could be zero...if so most likely file
        // is empty! check and update here accordingly.
        hits := ProcessChunk(buf, 0, length, jsonParams)
        // Calc last "HASH_LIMIT" bytes hash!
        retHashStr, _ := calcSha256FromChunkBackward(filenameStr, logFileStats.ReadEndPos, HASH_LIMIT)

        // This is the return value if all is successfull
        result.Filename = filenameStr
        result.Hash = retHashStr
        result.BeginPos = beginPos
        result.EndPos = logFileStats.ReadEndPos //read from global pointer
        result.Hits = hits
    }

    log.Finest("Result from Log() is:\n%s", result)
    return result, nil // FIXME: return nil as en error is bad!
}



// TODO: implement rounding when getting intervals that will be in the
// middle of a line. We need to search till previou EOL or next EOL


// This should inplement a log extraction method requested by user.
// That is the user request the chunk of a particular log for reading
// TODO(raul): Decide whether lines or positions while be used as params.
func Extract(jsonParams *zjson.JsonParams) (interface{}, error) {
    log.Debug("Initiating the execution of Extract()")

    //FIXME(raul): this is need as ReadChunk method uses this variable
    //this should then be refactored as it does not make much sense.
    //Maybe passing a reference is better or write a simpler
    //ReadChunk() method.
    logFileStats = new(LogFileMetaData)

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

    result := new(zjson.ExtractLog)
    result.Filename = filenameStr
    result.BeginPos = beginPos
    result.EndPos = endPos

    // TODO(raul): put here logic to get chunk!
    buf, _ := readChunk(filenameStr, beginPos, endPos - beginPos)
    log.Finest("Read buffer: %s", buf)
    if len(buf) == 0 {
        result.Lines = []string{}
    } else {
        eol := []byte{'\n'}
        lines := bytes.Split(buf, eol)
        for _, v := range lines {
        //for i, v := range lines {
            //lineNumber := i + _beginPos + 1 // line starts from 1, whereas i starts from 0
            result.Lines = append(result.Lines, string(v))
        }
    }

    log.Finest("Lines found:%s", result.Lines)

    log.Finest("Result from Extract() is:\n%s", result)
    return result, nil // FIXME: return nil as en error is bad!
}
