import random


JSON_VERSION = '2.0'

def getId(limit=99999999):
    return random.randint(1, limit)


search_log_params = [
    ["ERROR PATTERN","CLEAR PATTERN"],
    [".*stopped.*"  ,".*started.*"],
    [".*hello.*"    ,".*bye.*"],
    [".*hola.*"     ,".*chao.*"],
    [".*foo.*"      ,".*bar.*"]
]

# Multiple valid and invalid json request for test
# - Correct queries have the leading s_
# - Wrong queries have the leading f_
REQ = {
    's_status.info' : { 
        "jsonrpc": JSON_VERSION,
        "method": "status.info",
        "params": {},
        "id": getId()
    },
    'f_status.info' : { 
        "jsonrpc": JSON_VERSION,
        "methods": "status.info",
        "paramss": {},
        "id": getId()
    },
    's_search.log.000' : {
        "jsonrpc": JSON_VERSION,
        "method": "search.log",
        "params": {
            "pattern": search_log_params,
            "filename": "/home/raul/Code/go/src/zagent/util/foo.log.100",
            "hash": "9cd19cc133c3d66f30131cbb05776a05b1853770723b620c5bc9025a9419751c",
            "endpos": 100, 
            "beginpos": 0
            },
        "id": getId()
    },
    's_search.log.100' : {
        "jsonrpc": JSON_VERSION,
        "method": "search.log",
        "params": {
            "pattern": search_log_params, # Look in the top of the file!
            "filename": "/home/raul/Code/go/src/zagent/util/foo.log.100",
            "hash": "9cd19cc133c3d66f30131cbb05776a05b1853770723b620c5bc9025a9419751c",
            "endpos": 300, 
            "beginpos": 100
            },
        "id": getId()
    },
    's_search.log.300' : {
        "jsonrpc": JSON_VERSION,
        "method": "search.log",
        "params": {
            "pattern": search_log_params,
            "filename": "/home/raul/Code/go/src/zagent/util/foo.log.100",
            "hash": "a4142b23f1f5ba5c0e4bf8a0b6e4620e2406bb87e13ce8373dc7c653628e5e96",
            "endpos": 450, 
            "beginpos": 300
            },
        "id": getId()
    },
    'f_search.log.100' : {
        "jsonrpc": JSON_VERSION,
        "method": "search.log",
        "params": {
            "pattern": search_log_params, # Look in the top of the file!
            "filename": "/home/raul/Code/go/src/zagent/util/foo.log.100",
            "hash": "_9cd19cc133c3d66f30131cbb05776a05b1853770723b620c5bc9025a9419751c",
            "endpos": 300, 
            "beginpos": 100
            },
        "id": getId()
    },
    'f_search.log.122' : {
        "jsonrpc": JSON_VERSION,
        "method": "search.log",
        "params": {
            "pattern": search_log_params, # Look in the top of the file!
            "filename": "/home/raul/Code/go/src/zagent/util/foo.log.100",
            "hash": "_9cd19cc133c3d66f30131cbb05776a05b1853770723b620c5bc9025a9419751c",
            "endpos": 122, 
            "beginpos": 0
            },
        "id": getId()
    },
    'f_search.log.300' : {
        "jsonrpc": JSON_VERSION,
        "method": "search.log",
        "params": {
            "pattern": search_log_params, # Look in the top of the file!
            "filename": "/home/raul/Code/go/src/zagent/util/foo.log.100",
            "hash": "a4142b23f1f5ba5c0e4bf8a0b6e4620e2406bb87e13ce8373dc7c653628e5e96",
            "endpos": 450, 
            "beginpos": 0
            },
        "id": getId()
    },
}
