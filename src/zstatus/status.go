package zstatus

import ( 
  "runtime"
  "strconv"
  "zconfig"
  //"zjson"
)

const VERSION = "1.0.2"

type AgentInfo struct {
  Jsonrpc string            `json:"jsonrpc"`
  ID int                    `json:"id"`
  Result map[string] string `json:"result"`
}


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

