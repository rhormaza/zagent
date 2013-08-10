#!/usr/bin/env python

# Copyright (c) Twisted Matrix Laboratories.
# See LICENSE for details.

import sys
import time
import json
import random

from OpenSSL import SSL

from twisted.internet.protocol import ClientFactory
from twisted.protocols.basic import LineReceiver
from twisted.internet import ssl, reactor, defer

from jsondata import REQ

class ZagentClient(LineReceiver):
    delimiter = b'\r\n\r\n'
    MAX_LENGTH = 65536 # 64*1024

    def __init__(self, method):
        self.method = method

    def connectionMade(self):
        out = json.dumps(REQ[self.method], indent=2)
        self.sendLine(out)

    def connectionLost(self, reason):
        print 'connection lost (protocol)'

    def lineReceived(self, line):
        print "received:", line
        if line == self.delimiter:
            self.transport.loseConnection()

class ZagentClientFactory(ClientFactory):
    protocol = ZagentClient

    def __init__(self, method):
        self.method = method
        self.deferred = defer.Deferred()

    def buildProtocol(self, addr):
        proto = ZagentClient(self.method)
        self.connectedProtocol = proto
        return proto

    def clientConnectionFailed(self, connector, reason):
        print 'connection failed:', reason.getErrorMessage()
        if self.deferred is not None:        
            reactor.callLater(0, self.deferred.errback, reason)

    def clientConnectionLost(self, connector, reason):
        print 'connection lost:', reason.getErrorMessage()
        if self.deferred is not None:        
            reactor.callLater(0, self.deferred.errback, reason)
        #reactor.stop()

def request(method, host='127.0.0.1', port='44443'):
    """
    This sends a JSON query to zagent!
    """
    factory = ZagentClientFactory(method)

    ccf = ssl.ClientContextFactory()
    ccf.method = SSL.TLSv1_METHOD
    
    reactor.connectSSL(host, int(port), factory, ccf)
    factory1 = ZagentClientFactory(method)

    ccf1 = ssl.ClientContextFactory()
    ccf1.method = SSL.TLSv1_METHOD
    
    reactor.connectSSL(host, int(port), factory1, ccf1)
    reactor.run()

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


