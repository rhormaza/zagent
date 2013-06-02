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
            "hash": "a SHA-256 hash",
            "end_pos": "999999999999999999", 
            "begin_pos": "0"
        },
        "id": 2
    }`),
}
