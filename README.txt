

program (aka zagent) will be designed with the following structure


|-- bin                 // Binaries directory
|   `-- zagent          // Main binary
|-- log                 // Log files directory
|   `-- zagent.log      // Main log file
|-- usr                 // User (zagent) data directory
|   `-- cert            // Certificates directory
|       `-- zagent.pem  // Certificate identifying this zagent



Directory structure

.
|-- DIR_STRUCT.txt
|-- Makefile
|-- README.txt
|-- TODO.txt
|-- bin
|   |-- search
|   `-- zagent
|-- log
|-- pkg
|   `-- linux_amd64
|       |-- log4go.a
|       |-- search.a
|       `-- util.a
|-- src
|   |-- config
|   |   `-- config.go
|   |-- log4go
|   |   |-- LICENSE
|   |   |-- README
|   |   |-- config.go
|   |   |-- examples
|   |   |   |-- ConsoleLogWriter_Manual.go
|   |   |   |-- FileLogWriter_Manual.go
|   |   |   |-- SimpleNetLogServer.go
|   |   |   |-- SocketLogWriter_Manual.go
|   |   |   |-- XMLConfigurationExample.go
|   |   |   `-- example.xml
|   |   |-- filelog.go
|   |   |-- log4go.go
|   |   |-- log4go_test.go
|   |   |-- pattlog.go
|   |   |-- socklog.go
|   |   |-- termlog.go
|   |   `-- wrapper.go
|   |-- main
|   |   `-- zagent.go
|   |-- search
|   |   `-- search.go
|   |-- tcpserver
|   |   |-- archive
|   |   |   `-- archive.go
|   |   |-- handler
|   |   |   `-- agthandler
|   |   |       `-- agthandler.go
|   |   |-- server
|   |   |   `-- server.go
|   |   `-- tcpserver.go
|   `-- util
|       `-- util.go
`-- zagent.log

16 directories, 34 files
