__doc__ = """This module implements common methods to query Zagent
from cli , it uses the RequestManager class.
"""
__author__ = "Raul Hormazabal"
__email__ = "rhormaza@gmail.com"

import sys
import time
import json
import random
import logging

# Local imports
# FIXME(raul): This is only for a moment, find a better way!
from jsondata import REQ
from utils import PersistManager
from jsondata import JSON_REQ_FMT, JSON_STD_REQ
from zagent_twisted import RequestManager

#log = logging.getLogger('zen.zagent.client')
logging.basicConfig()
log = logging.getLogger('zen.zagent')

def create_dict_args(args):
    """
    This function converts a string such as:
    
        arg_1=val_1,arg_2=val_2,....,arg_N=val_N
    
    in a dict() like:

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

class WrongArgumentError(Exception):
    pass

class SearchLog(RequestManager):

    def _get_metadata(self, result):
        """
        Get data we need for next query, this 
        can raise an KeyError exception that
        we will let propagate here.
        """
        return {
            "beginpos"  : result["result"]["beginpos"],
            "endpos"    : result["result"]["endpos"],
            "hash"      : result["result"]["hash"],
            "filename"  : result["result"]["filename"]
        }

    def create_request(self, data, method='search.log'):
        """
        This creates a JSON request using <search.log> method.
        
        :params data: it holds data needed to create JSON request
        :type data: dict()

        rtype: dict()
        """

        # data must be a dict()
        if type(data) != type(dict()):
            raise WrongArgumentError()

        try:
            # Try to read cached data first
            tmp_meta_data = self.read_from_cache(
                    self.host_or_ip,
                    self.datasource,
                    meta_data=True)
        except Exception as e:
            # If no cached data, then from the beginning
            _hash = "xxx"
            _endpos = 100
            _beginpos = 0
        else:
            _hash = tmp_meta_data.get("hash")
            _beginpos = tmp_meta_data.get("endpos")
            _endpos = _beginpos + 100

        JSON_STD_REQ["method"] = method
        JSON_STD_REQ["params"]["filename"] = data.get("filename")
        JSON_STD_REQ["params"]["pattern"] = data.get("pattern")
        JSON_STD_REQ["params"]["hash"] = _hash 
        JSON_STD_REQ["params"]["endpos"] = _endpos
        JSON_STD_REQ["params"]["beginpos"] = _beginpos

        return JSON_STD_REQ

    def on_data_ready(self, result):
        try:
            unjson_meta_data = self._get_metadata(result)
            self.persist(
                    self.host_or_ip,
                    self.datasource,
                    unjson_meta_data,
                    meta_data=True)
        except KeyError:
            # we don't care...for now :)
            pass 

class LogExtract(RequestManager):
    """
    Simple class that create a <status.info> request to Zagent

    Valid JSON query looks like:

    {
      "params": {}, 
      "jsonrpc": "2.0", 
      "method": "log.extract", 
      "id": 71236624
    }

    Valid JSON reply:

    {
      "jsonrpc": "2.0",
      "result": {
      },
      "id": 71236624
    }
    """

    def create_request(self, data, method='log.extract'):
        """
        This creates a JSON request using <status.info> method.
        
        :params data: it holds data needed to create JSON request
        :type data: dict()

        rtype: dict()
        """

        # data must be a dict()
        if type(data) != type(dict()):
            raise WrongArgumentError()

        JSON_STD_REQ["method"] = method
        JSON_STD_REQ["params"]["filename"] = data.get("filename")
        JSON_STD_REQ["params"]["beginpos"] = data.get("beginpos")
        JSON_STD_REQ["params"]["endpos"] = data.get("endpos")

        return JSON_STD_REQ

class StatusInfo(RequestManager):
    """
    Simple class that create a <status.info> request to Zagent

    Valid JSON query looks like:

    {
      "params": {}, 
      "jsonrpc": "2.0", 
      "method": "status.info", 
      "id": 71236624
    }

    Valid JSON reply:

    {
      "jsonrpc": "2.0",
      "result": {
        "version": "1.0.0",
        "os": "linux",
        "port": 44443
      },
      "id": 71236624
    }
    """

    def create_request(self, data, method='status.info'):
        """
        This creates a JSON request using <status.info> method.
        
        :params data: it holds data needed to create JSON request
        :type data: dict()

        rtype: dict()
        """

        # data must be a dict()
        if type(data) != type(dict()):
            raise WrongArgumentError()

        JSON_STD_REQ["method"] = method

        return JSON_STD_REQ


def process_args():
    args = create_dict_args(sys.argv[2])
    methd = args.get('method', None)
    if methd == 'search.log':
        o = SearchLog("127.0.0.1", 44443, "bar")
        unjson_params = {
            "filename"  : "/home/raul/Code/go/src/zagent/util/foo.log.100",
            "pattern"   : [
                ["ERROR PATTERN", "CLEAR PATTERN"],
                [".*stopped.*"  , ".*started.*"],
                [".*hello.*"    , ".*bye.*"],
                [".*hola.*"     , ".*chao.*"],
                [".*foo.*"      , ".*bar.*"]
            ]
        }
        o.request(o.create_request(unjson_params, "search.log"))
        #o.do_request(unjson_params, "search.log")
    elif methd == 'log.extract':
        o = LogExtract("127.0.0.1", 44443, "bar")
        unjson_params = {
            "filename"  : "/home/raul/Code/go/src/zagent/util/foo.log.100",
            "beginpos" : 5,
            "endpos" : 123
        }
        o.do_request(unjson_params, 'log.extract')
    elif methd == 'status.info':
        o = StatusInfo("127.0.0.1", 44443, "bar")
        unjson_params = {}
        # This uses base class do_request() call, which wraps all
        # processing in one call.
        #o.request(o.create_request(unjson_params, "status.info"))
        o.do_request(unjson_params, "status.info")

    #o.request(**create_dict_args(sys.argv[2]))

def main():
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
            log.setLevel(logging.DEBUG)
            process_args()

    else:
        print "forgot to type command."
        pass
    sys.exit(0)

if __name__ == '__main__':
    main()
