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
  router *zrouter.RouterMap
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
  error msg
}

func (e *JsonParseError) Error() string {
  return e.msg
}


func getMethodParam(inputObj interface{}) ([]string, []string, error){

  method := ""
  param := []string{}

  jsonObj, ok = inputObj.(map[string]interface{})
  if !ok { // if failed to do type assertion
     return []string{}, []string{}, JsonParseError{"Request is not a valid json string"}
  }
  for k,v := range jsonObj {
    if k == 'method' {
      method = v
    }
    else if k == 'params' {
      param = v
    }
  }

  // method[0] is namespace, method[1] is mehtodname
  return strings.Split(method, "."), param, nil

}

func (s *Server) session(conn *net.Conn) {
    defer (*conn).Close()
    //var input, output []byte
    var input, output []byte


    if err := s.read_data(conn, &input); err != nil {
      log.Println(err)
      s.Response(conn, err) // we don't care if the msg was sent successfully
      return
    }

    if inputObj, err := json.DecodeJson(input); err != nil {
      log.Println(err)
      s.Response(conn, err) // we don't care if the msg was sent successfully
      return
    }

    if method, param, err := getMethodParam(inputObj); err != nil {
      log.Println(err)
      s.Response(conn, err) // we don't care if the msg was sent successfully
      return
    }

    if res, err := s.router.methodMap[method[0]][method[1]](param); err != nil {
      log.Println(err)
      s.Response(conn, err) // we don't care if the msg was sent successfully
      return
    }

    if output, err := json.EncodeJson(res); err != nil {
      log.Println(err)
      s.Response(conn, err) // we don't care if the msg was sent successfully
      return
    }

    /*
    if _, ok := s.Handlers[code]; !ok {
      log.Println("No such code") // TO-DO: need improve later
      return
    }

    // call handler function
    err = s.Handlers[code](&input, &output, conn)
    if err != nil {
      err_msg := "Failed to call handler function" 
      s.Response(conn, err_msg) // we don't care if the msg was sent successfully
      log.Println(err_msg) //TO-DO: need improve later
    }
    */

    err = s.write_data(conn, &output)
    if err != nil {
      log.Println(err)
      return
    }

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
