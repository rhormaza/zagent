package main

import (
    "os"
    "zutil"
    "ztcpserver"
)

func main() {
    // Load configuration file
    zutil.LoadConfig(os.Args)

    // Logger always first!
    log := zutil.SetupLogger(zutil.GetConf().LogPath, zutil.GetConf().LogLevel)

    server := ztcpserver.Server {
                Cert    : zutil.GetConf().CertFilePemPath,
                Prvkey  : zutil.GetConf().CertFileKeyPath,
            }
    server.Run(zutil.GetConf().ListenAddressAndPort, false)

    // we need to close the logger to clear the buffers!
    log.Close()
}
