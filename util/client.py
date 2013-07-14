#!/usr/bin/env python
# Usage: python check_agent.py IP 

import socket
import ssl
import sys
import time
import json

from pprint import pprint as pp

SUCCESS = "test.info_success"
ERROR = "test.info_error"

REQ = [
        { 
            "jsonrpc": "2.0",
            "method": "status.info",
            "params": {},
            "id": 2
        },
        {
            "jsonrpc": "2.0",
            "method": "search.log",
            "params": {
                "pattern": [
                    [
                        "ERROR PATTERN",
                        "CLEAR PATTERN"
                        ],
                    [
                        ".*stopped.*",
                        ".*started.*"
                        ],
                    [
                        ".*hello.*",
                        ".*bye.*"
                        ],
                    [
                        ".*hola.*",
                        ".*chao.*"
                        ],
                    [
                        ".*foo.*",
                        ".*bar.*"
                        ]
                    ],
                "filename": "/tmp/foo.txt",
                "hash": "9d07b162e93d76901d07dfa53d755e98806f3d9a8ac765fe8ca47f34d71a4ebe",
                "endpos": 110, 
                "beginpos": 0
                },
            "id": 2
        },
        # bad query
        { 
            "jsonrpc": "2.0",
            "methods": "status.info",
            "paramss": {},
            "id": 2
        },
        ]

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
    return ssl.wrap_socket(s, ssl_version=ssl.PROTOCOL_TLSv1, keyfile="zagent_client.key",
                              certfile="zagent_client.pem")


def request(host, port, method):
    import json
    #REQ["method"] = method
    json_str = json.dumps(REQ[method], indent=4)
    print "request:"
    print json_str
    s = create_socket(host, int(port))
    s.sendall(json_str + "\r\n\r\n")
    response = ""
    while True:
        msg = s.recv(1024)
        response += msg
        if response[-4:] == "\r\n\r\n":
            break

    print "Recevied:"
    print response
    s.close()
  
if __name__ == '__main__':
    try:
        import sys
        if len(sys.argv) <= 1:
            print "IP doesn't specifiy"
            sys.exit(-1)
        ip, reqId = sys.argv[1:3]
        request(ip, 44443, int(reqId))
    except Exception, e:
        sys.stderr.write("Exception: %s" % e)
        sys.exit(-1)
    sys.exit(0)


