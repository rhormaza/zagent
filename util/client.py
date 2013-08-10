#!/usr/bin/env python2

import socket
import ssl
import sys
import time
import json
import random

from jsondata import REQ

USAGE = "python %s command args1,args2,args3,argsN" % __file__

class ArgumentError(Exception): pass

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

def create_socket(host, port):
    """
    Helper function that creates a SSL connection to "host"
    """
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(300)
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

def request(method, host='127.0.0.1', port='44443'):
    """
    This sends a JSON query to zagent!
    """

    try:
        json_query = json.dumps(REQ[method], indent=4)
        print "============ Request =============="
        print json_query
        s = create_socket(host, int(port))
        s.sendall(json_query + "\r\n\r\n")
    except KeyError as e:
        print "method: <%s> does not exist." % method
        sys.exit(-1)
    except Exception as e:
        raise

    json_response = ""
    while True:
        msg = s.recv(1024)
        json_response += msg
        if json_response[-4:] == "\r\n\r\n":
            break

    print "============= Reply ==============="
    print json_response
    s.close()
  
def create_file(
        filePrefix='foo', 
        number=[10**1, 10**2, 10**3, 10**4,10**5,10**6]
        ):
    lines = [
        'foo%06d\n',
        'bar%06d\n'
    ] 

    for _len in number:
        fileName = '%s.log.%s' % (filePrefix, _len)
        with open(fileName, 'w') as fp:
            for i in range(0, _len):
                fp.write(lines[0] % i)

def create_dict_args(args):
    """
    This converts:
    
        arg_1=val_1,arg_2=val_2,....,arg_N=val_N
    
    in

        { 
            arg_1 : val_1,
            arg_2 : val_2,
            ...
            ...
            arg_N : val_N
        }
    
    """

    exMsg = "You have provided wrong arguments. " + \
            "They should be: " + \
            "arg_1=val_1,arg_2=val_2,....,arg_N=val_N\n"
    _dict = dict()
    try:
        args_list = args.split(',')
        if args_list:
            for a in args_list:
                pair = a.split('=')
                if len(pair) == 2:
                    _dict[pair[0]] = pair[1]
                else:
                    raise
        else:
            raise
    except:
        raise Exception(exMsg)
    else:
        return _dict



if __name__ == '__main__':
    cmdName = sys.argv[1]
    if cmdName == 'create_file':
        if len(sys.argv) != 2:
            raise ArgumentError("wrong arguments!\n")
        else:
            create_file()

    elif cmdName == 'request':
        if len(sys.argv) != 3:
            raise ArgumentError("wrong arguments!\n")
        else:
            request(**create_dict_args(sys.argv[2]))

    else:
        pass
    sys.exit(0)


