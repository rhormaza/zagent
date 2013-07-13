package zstatus

import ( 
  "runtime"
  "strconv"
  "zconfig"
  "zjson"
)

const VERSION = "1.0.2"

type AgentInfo struct {
  jsonrpc string
  id int
  result map[string] string
}




func GetAgentInfo(*zjson.JsonParams) (interface {}, zjson.JsonError) {
  os := runtime.GOOS // get system type

  //TODO: fix config later
  config := zconfig.LoadConfig("FIXMELATER")
  port := strconv.Itoa(config.ListenPort)
  
  agent_info := AgentInfo{ 
                           jsonrpc: "2.0",
                           id: 2,
                           result: map[string]string{
                             "version": VERSION,
                             "os": os,
                             "port": port,
                           },
                         }
   //TODO: need return err msg
  return agent_info, nil
}

