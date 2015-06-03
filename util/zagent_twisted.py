__doc__ = """This module implements classes using Twisted librabries to
query Zagent.
"""
__author__ = "Raul Hormazabal"
__email__ = "rhormaza@gmail.com"

import json
import logging

from OpenSSL import SSL

from twisted.internet.protocol import ClientFactory
from twisted.protocols.basic import LineReceiver
from twisted.internet import ssl, reactor, defer
from twisted.internet.error import ReactorNotRunning

from utils import PersistManager

#logging.basicConfig()
log = logging.getLogger('zen.zagent.twisted')

class ZagentClient(LineReceiver):
    delimiter = b'\r\n\r\n'
    MAX_LENGTH = 65536 # 64*1024
    line = ''

    def __init__(self, req):
        self.request = req

    def connectionMade(self):
        
        out = json.dumps(self.request, indent=2)
        log.debug("Sending query: %s", out)

        self.sendLine(out)

    def connectionLost(self, reason):
        pass

    def lineReceived(self, line):
        log.debug("received: %s", line)
        if line == self.delimiter:
            self.factory.deferred = None
            self.transport.loseConnection()
        else:
            self.line += line

class ZagentClientFactory(ClientFactory):
    protocol = ZagentClient
    line = ''

    def __init__(self, req, stop_reactor=False):
        
        self.request = req
        self.deferred = defer.Deferred()
        self.stop_reactor = stop_reactor

    def buildProtocol(self, addr):
        
        proto = ZagentClient(self.request)
        self.connectedProtocol = proto
        return proto

    def clientConnectionFailed(self, connector, reason):
        
        log.error('connection failed: %s', reason.getErrorMessage())
        
        if self.deferred is not None:
            reactor.callLater(0, self.deferred.callback, reason)

    def clientConnectionLost(self, connector, reason):
        
        log.debug('connection lost: %s', reason.getErrorMessage())

        if self.deferred is not None:
            self.deferred.callback(self.parseData(self.connectedProtocol.line))
        
        if self.stop_reactor:
            reactor.stop()

    def parseData(self, line):
        """
        Process json output and return
        """
        try:
            return json.loads(line)
        except Exception as e:
            log.critical("Error processing json returned output. %s", e)
            raise e

class RequestManager(object):

    def __init__(self, host_or_ip, port, datasource):
        self.host_or_ip = host_or_ip
        self.datasource = datasource

        if type(port) == type(str):
            self.port == int(port)
        else:
            self.port = port

    def on_data_ready(self, result):
        """
        This a callback to do operations like cleaning, caching, etc
        after receiving data successfully from the connection.

        It is not compulsory so it can or cannot be implemented by a subclass.
        """
        pass

    #def request(self, unjson_data, host='127.0.0.1', port='44443',
    #        stop_reactor=True):
    def request(self, unjson_data, stop_reactor=True):
        """
        This sends a JSON query to zagent!
        """

        factory = ZagentClientFactory(unjson_data, stop_reactor)
        # This on_data_ready() callback, see method docs.
        factory.deferred.addCallback(self.on_data_ready)

        ccf = ssl.ClientContextFactory()
        ccf.method = SSL.TLSv1_METHOD

        reactor.connectSSL(self.host_or_ip, self.port, factory, ccf)
        reactor.run()

    def persist(self, host_or_ip, datasrc, unjson_data, meta_data=False):

        try:
            PersistManager.persist_to_fs(
                    host=host_or_ip,
                    datasource=datasrc,
                    data=json.dumps(unjson_data),
                    metadata=meta_data)
        except IOError:
            # Most likely dir doesn't exist.
            # Let's create it and retry.
            PersistManager.create_dir()
            PersistManager.persist_to_fs(
                    host=host_or_ip,
                    datasource=datasrc,
                    data=json.dumps(unjson_data),
                    metadata=meta_data)

    def read_from_cache(self, host_or_ip, datasrc, meta_data=True):
        """
        Read metadata if exist in filesystem, it should be called
        from create_request()
        """
        data = PersistManager.read_from_fs(
                    host=host_or_ip, 
                    datasource=datasrc,
                    metadata=meta_data)
        return json.loads(data)

    def create_request(self, data, method):
        """
        It has to be implement in the subclass
        """
        raise NotImplementedError

    def do_request(self, unjson_params, method):
        self.request(self.create_request(unjson_params, method))
