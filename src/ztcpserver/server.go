package ztcpserver

import (
  "crypto/tls"
  "net"
  "log"
  "time"
  "zjson"
  "zrouter"
  "strings"
)

const BUFSIZ = 1024
const TIMEOUT = 5

func checkErr(err error) {
  if err != nil {
    log.Fatal(err) //TO-DO: make it more readable
  }
}


type Server struct {
  Cert, Prvkey string
  listener net.Listener
  Router *zrouter.RouterMap
}


func (s *Server) read_data( conn *net.Conn, data *[]byte ) error {
  var (
    buf [BUFSIZ]byte
  )

  // read data
  for {
    n, err := (*conn).Read(buf[0:])
    if err != nil {
      return err 
    }
    *data = append(*data, buf[0:n]...)
    last_four := string((*data)[len(*data)-4:]) // get last four characters
    if last_four == "\r\n\r\n" { //check if it is the end
      break
    }
  }
  return nil
}


func (s *Server) write_data( conn *net.Conn, data *[]byte ) error{
  var (
    size = len(*data) 
    written = 0
  )

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


func (s *Server) Response( conn *net.Conn, msg string ) error {
  msg_b := []byte(msg)
  return s.write_data( conn, &msg_b )
}


func (s *Server) Stop() error{
  return s.listener.Close()
}


type JsonParseError struct {
  msg string
}


func (e *JsonParseError) Error() string {
  return e.msg
}


func getMethodParam(inputObj interface{}) ([]string, map[string]interface{}, error){

  method := ""
  var param map[string]interface{}

  log.Println(inputObj)
  jsonObj, ok := inputObj.(map[string]interface{})
  log.Println(jsonObj)
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


func (s *Server) session(conn *net.Conn) {
    defer (*conn).Close() 
    var input, output []byte

    if err := s.read_data(conn, &input); err != nil {
      log.Println(err)
      s.Response(conn, err.Error()) 
      return
    }

    inputObj, err := zjson.DecodeJson(input)
    if err != nil {
      log.Println(err)
      s.Response(conn, err.Error()) 
      return
    }

    method, param, err := getMethodParam(inputObj)
    if err != nil {
      log.Println(err)
      s.Response(conn, err.Error()) 
      return
    }

    res, err := (*s.Router)[method[0]][method[1]](param)
    if err != nil {
      log.Println(err)
      s.Response(conn, err.Error())
      return
    }

    log.Println(res)
    if output, err = zjson.EncodeJson(res); err != nil {
      log.Println(err)
      s.Response(conn, err.Error())
      return
    }
    
    output = append(output, "\r\n\r\n"...)
    log.Println("before write")

    err = s.write_data(conn, &output)
    log.Println(&output)
    log.Println(output)

    if err != nil {
      log.Println(err)
      return
    }
    log.Println("Done")
}


func (s *Server) Run(laddr string, non_routine bool) {
  cert, err := tls.LoadX509KeyPair(s.Cert, s.Prvkey)
  checkErr(err)

  config := tls.Config{Certificates: []tls.Certificate{cert}}

  s.listener, err = tls.Listen("tcp", laddr, &config)
  checkErr(err)

  for {
    conn, err := s.listener.Accept()
    if err != nil {
      log.Println(err)
      continue
    }
    conn.SetDeadline(time.Now().Add(TIMEOUT * time.Second)) //set timeout
    if non_routine {
      s.session(&conn)
    } else {
      go s.session(&conn)
    }

  }
}
