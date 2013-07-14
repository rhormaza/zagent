package ztcpserver

import (
    "crypto/tls"
    "net"
    "time"
    "strings"
    "zjson"
    "zrouter"
    "zutil"
)

var log = zutil.SetupLogger("/tmp/zagent.log", 2)

const BUFFSIZE = 1024
const TIMEOUT = 5
const END_OF_JSON = "\r\n\r\n"
const LEN_END_OF_JSON = 4

func checkErr(err error) {
    if err != nil {
        log.Critical(err) //TO-DO: make it more readable
    }
}


type Server struct {
    Cert, Prvkey string
    listener net.Listener
    Router *zrouter.RouterMap
}

func (s *Server) read_data( conn *net.Conn, data *[]byte ) error {
    var (
        buffer [BUFFSIZE]byte
    )

    // Read the data
    for {
        numBytesRecv, err := (*conn).Read(buffer[0:])
        if err != nil {
            return err 
        }
        // numBytes could be less than four which is the minimum data  required.
        if numBytesRecv > LEN_END_OF_JSON {

            *data = append(*data, buffer[0:numBytesRecv]...)
            log.Debug(*data)

            // Get last LEN_END_OF_JSON bytes to confirm end of JSON data.
            last_four := string((*data)[(len(*data) - LEN_END_OF_JSON):])

            if last_four == END_OF_JSON { //check if it is the end
                break
            }
        }
    }
    return nil
}

func (s *Server) write_data( conn *net.Conn, data *[]byte ) error{
    var (
        size = len(*data) 
        written = 0
    )

    log.Debug("About to write data to the network.")
    for {
        n, err := (*conn).Write(*data)
        if err != nil {
            return err 
        }
        written += n
        if written == size {
            break
        }
    }

    return nil
}

func (s *Server) Run(listenAddr string, non_routine bool) {
    cert, err := tls.LoadX509KeyPair(s.Cert, s.Prvkey)
    checkErr(err)

    config := tls.Config{Certificates: []tls.Certificate{cert}}

    s.listener, err = tls.Listen("tcp", listenAddr, &config)
    checkErr(err)

    for {
        conn, err := s.listener.Accept()
        if err != nil {
            log.Info(err)
            continue
        }
        conn.SetDeadline(time.Now().Add(TIMEOUT * time.Second)) //set timeout
        if non_routine {
            s.processQuery(&conn)
        } else {
            go s.processQuery(&conn)
        }

    }
}

func (s *Server) Response( conn *net.Conn, msg string ) error {
    msg_b := []byte(msg)
    return s.write_data( conn, &msg_b )
}


func (s *Server) Stop() error{
    return s.listener.Close()
}



//
// This function is intended to process to do the steps to process a file.
// The steps are:
// 1.- Launch ztcpserver, this server only handles request from the netwok, 
//     nothing else. It should be simple and lightweight
// 2.- Call zjson.Decode in order to decode a json object. All json decode/encode
//     logic goes here and only here!
// 3.- If decode fails, return the error back to the client 
//     with and known error code
// 4.  If decode is successful, make a call to the zrouter Map with the key
//     extracted from the JSON query. zrouter Map will redirect
//     and call the process.
// 5.- Return result of the method called from the router (zrouter.Map).
// 6.- Done.
//
func (s *Server) processQuery(conn *net.Conn) {
    defer func() {
        if x := recover(); x != nil {
            log.Critical("Panic error!")
            err := zjson.Err{-10001, "error!!!"}
            output, _ := zjson.EncodeJsonFail(err, 1)
            output = append(output, END_OF_JSON...)
            s.write_data(conn, &output) 
        }
        (*conn).Close() 
    }()
    var input, output []byte

    if err := s.read_data(conn, &input); err != nil {
        log.Info(err)
        s.Response(conn, err.Error()) 
        return
    }

    jsonObj, err := zjson.DecodeJson2(input)
    if err != nil {
        log.Error("Error decoding JSON data: %s", err)
        o, _ := zjson.EncodeJsonFail(err, 1)
        log.Debug("Sending error reply: %s", o)
        s.write_data(conn, &o) 
        return
    }

    // Obtain the namespace and mathod strings from JSON query
    namespaceAndMethod := strings.Split(*jsonObj.Method, ".")
    namespace := namespaceAndMethod[0]
    method := namespaceAndMethod[1]

    result, err := zrouter.Router[namespace][method](jsonObj.Params)
    if err != nil {
        log.Error("There was an error processing %s.%s(): %s", namespace, method, err)
        //call zjson.EncodeJsonError and return
        s.Response(conn, err.Error())
        return
    }

    log.Info(result)
    // At this point method required, ran sucessfully, hence we return success
    if output, err = zjson.EncodeJsonSuccess(result, (jsonObj.Id)); err != nil {
        log.Info(err)
        s.Response(conn, err.Error())
        return
    }

    output = append(output, END_OF_JSON...)
    log.Debug("About to write data back")

    err = s.write_data(conn, &output)
    log.Debug(output)

    if err != nil {
        log.Error("Error writing JSON data back: %s", err)
        return
    }
    log.Debug("Query processed.")
}

//func process() {
//    // FIXME: load arguments from a config file!  
//    // This is a example how to load stuff from a config package
//    //log.Info("Args: %s and config:%s", os.Args, zconfig.LoadConfig("asas"))
//
//    // Init our server
//    s := ztcpserver.Server{ Cert: "zagent.pem", Prvkey:"zagent.pem",
//    Handlers: make(map[string] func (*[]byte, *[]byte, *net.Conn) error)}
//
//    // Server now listens from from this port
//    s.Run(":44443", false)
//
//    // We need to get data somehow from the server and put in a byte[]
//    // rawDataQuery :- s.getDataIfReadyToBeRead()
//
//    // pass byte[] to the router.
//    jsonReply, err :=  zjson.DecodeJson(rawDataQuery)
//    if err != nil {
//        // Return valid JSON error reply with error code as well.
//        // In short make a call to zjson.EncodeJsonFail()
//        //resultJson := zjson.EncodeJsonFail(result)
//    } else {
//
//        // Make the call to the router Map with data decoded
//        // Json.Method should have the "key" to the next function call
//        // Json.Params should have paramaters passed within the query
//        result, err := zrouter.RouterMap[jsonReply.method](jsonReply.Params)
//
//        if err != nil {
//            // Return valid JSON error reply with error code as well.
//            // In short make a call to zjson.EncodeJsonFail()
//            //resultJson := zjson.EncodeJsonFail(result)
//        } else {
//            // Encode data returned from requested method (from the router)
//            //resultJson := zjson.EncodeJsonSuccess(result)
//        }
//    }
//    // Write encoded back to the client. 
//    //s.WriteBackToClient(resultJson)
//
//    // Close the logger to clear the buffers!
//    log.Close()
//
//}



/* FIXME: really needed?*/

// Deprecated?
type JsonParseError struct {
  msg string
}


// Deprecated?
func (e *JsonParseError) Error() string {
  return e.msg
}


// Deprecated?
func getMethodParam(inputObj interface{}) ([]string, map[string]interface{}, error){

  method := ""
  var param map[string]interface{}

  log.Info(inputObj)
  jsonObj, ok := inputObj.(map[string]interface{})
  log.Info(jsonObj)
  if !ok { // if failed to do type assertion
     return []string{}, map[string]interface{}{}, &JsonParseError{"Request is not a valid json string"}
  }
  for k,v := range jsonObj {
    if k == "method" {
      vv := v.(string)
      method = vv
    } else if k == "params" {
      vv := v.(map[string]interface{})
      param = vv
    }
  }
  // method[0] is namespace, method[1] is mehtodname
  return strings.Split(method, "."), param, nil
}


// Deprecated?
func (s *Server) session(conn *net.Conn) {
    defer (*conn).Close() 
    var input, output []byte

    if err := s.read_data(conn, &input); err != nil {
      log.Info(err)
      s.Response(conn, err.Error()) 
      return
    }

    inputObj, err := zjson.DecodeJson(input)
    if err != nil {
      log.Info(err)
      s.Response(conn, err.Error()) 
      return
    }

    method, param, err := getMethodParam(inputObj)
    if err != nil {
      log.Info(err)
      s.Response(conn, err.Error()) 
      return
    }

    res, err := (*s.Router)[method[0]][method[1]](param)
    if err != nil {
      log.Info(err)
      s.Response(conn, err.Error())
      return
    }

    log.Info(res)
    if output, err = zjson.EncodeJson(res); err != nil {
      log.Info(err)
      s.Response(conn, err.Error())
      return
    }
    
    output = append(output, END_OF_JSON...)
    log.Info("before write")

    err = s.write_data(conn, &output)
    log.Info(&output)
    log.Info(output)

    if err != nil {
      log.Info(err)
      return
    }
    log.Info("Done")
}


