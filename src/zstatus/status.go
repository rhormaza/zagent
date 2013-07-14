package zstatus

import ( 
    "runtime"
    "strconv"
    "zconfig"
    "zjson"
    "zutil"
)

// Always first!
var log = zutil.SetupLogger("/tmp/zagent.log", 2)

// Deprecated?
// It should come from zutil or a config file!
const VERSION = "1.0.2"

// Deprecated?
// This must be define in zjson.structures with *all* structs
// to don't lose track of them
type AgentInfo struct {
  Jsonrpc string            `json:"jsonrpc"`
  ID int                    `json:"id"`
  Result map[string] string `json:"result"`
}


// Deprecated?
//func GetAgentInfo(*zjson.JsonParams) (interface {}, error) {
func GetAgentInfo(map[string]interface{}) (interface {}, error) {
  os := runtime.GOOS // get system type

  //TODO: fix config later
  config := zconfig.LoadConfig("FIXMELATER")
  port := strconv.Itoa(config.ListenPort)
  
  agent_info := AgentInfo{ 
                           Jsonrpc: "2.0",
                           ID: 2,
                           Result: map[string]string{
                             "version": VERSION, 
                             "os": os,
                             "port": port,
                           },
                         }
   //TODO: need return err msg
  return agent_info, nil
}

func Process(jsonParams *zjson.JsonParams) (interface{}, error) {
    // Assert JSON Params, return an error if fails!
    log.Debug("Executing Process()")

    //TODO: fix config later
    config := zconfig.LoadConfig("FIXMELATER")
    port := int64(config.ListenPort)

    // Just for readibility
    result := new(zjson.StatusInfo)
    result.Os = runtime.GOOS
    result.Port = port
    result.Version = VERSION
    return result, nil
}
