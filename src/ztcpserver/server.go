package ztcpserver

import (
  "crypto/tls"
  "strconv"
  "net"
  "log"
)

const buf_size = 1024
const end_mark = byte(255) //end marker for data size

func checkErr(err error) {
  if err != nil {
    log.Fatal(err) //TO-DO: make it more readable
  }
}


func try_get_datasize( data *[]byte ) (int, int, error){
  var i int
  size := len(*data)
  for i=0; i<size; i++ {
    if (*data)[i] == end_mark {
      break
    }
  }

  if i != size {
    n, err := strconv.Atoi(string((*data)[1:i]))
    return n, i, err
  }
  return -1, -1, nil
}



type Server struct {
  Cert, Prvkey string
  listener net.Listener
  Handlers map[string] func (*[]byte, *[]byte, *net.Conn) error
}


func (s *Server) read_data( conn *net.Conn, data *[]byte ) error {
  var (
    err error
    n int
    data_size = -1
    begin = -1
    total = 0
    buf [buf_size]byte
  )

  // read data
  for {
    n, err = (*conn).Read(buf[0:])
    if err != nil {
      return err 
    }
    *data = append(*data, buf[0:n]...)
    total += n
    if data_size == -1 {
      data_size, begin, err = try_get_datasize(data)
      if err != nil {
         return err
      }
    }
    if data_size != -1 && begin != -1 && total-begin-1 == data_size {
      break
    }
  }
  *data = (*data)[begin+1:]
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


func (s *Server) session(conn *net.Conn) {
    defer (*conn).Close()
    var input, output []byte

    err := s.read_data(conn, &input)

    // need replace the error handling code with a function
    if err != nil {
      log.Println(err)
      return
    }
    // TO-DO: need check size of input
    code := string(input[0])

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

    err = s.write_data(conn, &output) 
    if err != nil {
      log.Println(err)
      return
    }

}


func (s *Server) AddHandler(code string,
                 handler func (*[]byte, *[]byte, *net.Conn) error) {
  s.Handlers[code] = handler
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

    if non_routine {
      s.session(&conn)
    } else {
      go s.session(&conn)
    }

  }
}
