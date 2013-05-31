package search

import (
    "os"
    "bytes"
    "regexp"
    "fmt"
    "crypto/sha256"
    "encoding/hex"
    "util"
)

// Logger always first!
var log = util.SetupLogger("zagent.log", 2)

type JsonObject struct {
    /*Fields start with Uppercase for JSON Marshalling */
    Jsonrpc string
    Method  string
    Params  JsonParams
    Id      int64
}

type JsonParams struct {
    Pattern     [][]string
    Filename    string
    Hash        string
    Begin_pos   int64
    End_pos     int64
}

type Hit struct {
    lineText    string
    lineNumber  int64
    lineBegin   int64
    //more metadata to hold?
}

type Pattern struct {
    pattern string
    hits    []Hit
}

type Query struct {
    pattern     string
    filename    string
    hash        string
    beginPos    int64
    endPos      int64
}

var hitSlices []Hit
var PatternHits = make(map[string] []Hit)

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

func doMatch(lineNumber int64, line string, jsonQuery *JsonObject) (hit bool) {
    for _, v := range jsonQuery.Params.Pattern {
        pattern := v[0] //Error pattern
        if m, _ := regexp.MatchString(pattern, line); m {
            PatternHits[pattern] = append(PatternHits[pattern], Hit{line, lineNumber, 9999999})
        }
    }
    return true

}

func ProcessChunk(buffer []byte, beginPos int, length int, jsonInfo *JsonObject) {
    /*
    TODO: Check for '\r' char in case of a windows box
          Check EOL during the split!
    */
    eol := []byte{'\n'}
    for i, v := range bytes.Split(buffer, eol) {
        lineNumber := i + beginPos + 1 // line starts from 1, whereas i starts from 0
        fmt.Println(lineNumber, " -> ", string(v), "len=", len(v))
        doMatch(int64(lineNumber), string(v), jsonInfo)
    }
    fmt.Println("len=", len(buffer), "sha=", hex.EncodeToString(CalcSha256(buffer)))
}


//func main() {
//    //fmt.Println(">>>> ", config.LoadConfig("asas") )
//
//    var jsonObj JsonObject
//    err := json.Unmarshal(jsonBlob, &jsonObj)
//    if err != nil {
//        fmt.Println("error:", err)
//    }
//
//
//    last, _ := strconv.Atoi(os.Args[1])
//    buf, err := ReadChunk(jsonObj.Params.Filename, 0, int64(last))
//    processChunk(buf, 0, last, &jsonObj)
//    fmt.Println("+++++++++++++")
//    mdStr, _ := CalcSha256FromChunk("/tmp/foo.txt", 0, int64(last))
//    fmt.Println(">>>>>>>>> ", mdStr)
//    fmt.Println("+++++++++++++")
//    for k, v := range PatternHits {
//        fmt.Println("k:", k, "v:", v)
//    }
//
//
//    //we need to close the logger to clear the buffers!
//    log.Close()
//}
