package zjson

// Put here all Json queries and results (succesfull and not)A
// This can also be used for do Unit-Testing!
var JsonBlobMap = map[string][]byte { 
    "searchlog_query": []byte(`{
        "jsonrpc": "2.0",
        "method": "get_error",
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
    }`),
}
