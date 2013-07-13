#!/usr/bin/env python
# Usage: python check_agent.py IP 

import socket
import ssl
import sys
import time
import json

SUCCESS = "test.info_success"
ERROR = "test.info_error"

REQ = { 
        "jsonrpc": "2.0",
        "method": None,
        "params": [],
        "id": 2
      }

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

def create_socket(host, port):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(3)
    while True:
       try:
           s.connect( (host, int(port)) )
       except:
           s.shutdown(socket.SHUT_RDWR)
           s.close()
           time.sleep(1)
       else:
           break
    return ssl.wrap_socket(s, ssl_version=ssl.PROTOCOL_SSLv3 , keyfile="client.pem",
                              certfile="client.pem")


def request(host, port, method):
    import json
    REQ["method"] = method
    json_str = json.dumps(REQ)
    print "request %s" % json_str
    s = create_socket(host, int(port))
    s.sendall(json_str + "\r\n\r\n")
    response = ""
    while True:
        msg = s.recv(1024)
        response += msg
        if response[-4:] == "\r\n\r\n":
            break

    print "Recevied: %s\n" % response
    s.close()
  
if __name__ == '__main__':
    try:
        import sys
        if len(sys.argv) <= 1:
            print "IP doesn't specifiy"
            sys.exit(-1)
        ip = sys.argv[1]
        request(ip, 44443, SUCCESS)
        request(ip, 44443, ERROR)
    except Exception, e:
        sys.stderr.write("Exception: %s" % str(e))
        sys.exit(-1)
    sys.exit(0)


