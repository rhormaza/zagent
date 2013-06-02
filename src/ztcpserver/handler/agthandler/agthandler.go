package agthandler

import (
  "net"
)

func Test_handler(input *[]byte, output *[]byte, conn *net.Conn) error{
  *output = append(*output, *input...)
  return nil
}


