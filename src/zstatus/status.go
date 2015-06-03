package zstatus

import ( 
    "runtime"
    "strconv"
//    "zconfig"
    "zjson"
    "zutil"
)

// Always first!
//var log = zutil.SetupLogger("/tmp/zagent.log", zutil.GetConf().Log)
var log = zutil.GetLogger()

func Info(jsonParams *zjson.JsonParams) (interface{}, error) {
    // Assert JSON Params, return an error if fails!
    log.Debug("Executing Process()")

    //TODO: fix config later
    //config := zconfig.LoadConfig("FIXMELATER")
   // port := int64(config.ListenPort)
    port, _ := strconv.Atoi(zutil.GetConf().ListenPort)

    // Just for readibility
    result := new(zjson.StatusInfo)
    result.Os = runtime.GOOS
    result.Port = int64(port)
    result.Version = zutil.GetNumericVersion()
    return result, nil
}
