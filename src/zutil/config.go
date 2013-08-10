package zutil

import (
    "strings"
    "zconfig"
)

// Configuration struct
type Config struct {
    ConfigFile  string

    CertPath            string
    CertFilePem         string
    CertFileKey         string
    CertFilePemPath     string
    CertFileKeyPath     string

    ListenAddress           string
    ListenPort              string
    ListenAddressAndPort    string

    LogPath     string
    LogLevel    string

    RootPath    string
}

var conf *Config

func processArgs(args []string) (confPath string) {

    binFullPath := strings.Split(args[0], "/")
    lenBinFullPath := len(binFullPath)
    prevDirPos := lenBinFullPath - 2

    if binFullPath[prevDirPos] == "." {
        // Same path
        confPath = "../etc"
    } else if binFullPath[prevDirPos] == "bin" {
        confPath = strings.Join(append(binFullPath[:prevDirPos], "etc"), "/")
    } else {
        confPath = "/opt/zagent/etc"
    }
    return
}

func LoadConfig(args []string) (*Config){
    confFile := processArgs(args) + "/" + "zagent.conf"
    file, _ := zconfig.LoadFile(confFile)

    // Create new sConfig struct to hold our config data
    c := new(Config) 

    c.ConfigFile = confFile

    c.CertPath, _ = file.Get("main", "certificate_path")
    c.CertFilePem, _ = file.Get("main", "certificate_pem")
    c.CertFileKey, _ = file.Get("main", "certificate_key")
    c.CertFilePemPath = c.CertPath + "/" + c.CertFilePem
    c.CertFileKeyPath = c.CertPath + "/" + c.CertFileKey

    c.ListenAddress, _ = file.Get("main", "listen_address")
    c.ListenPort, _ = file.Get("main", "listen_port")
    c.ListenAddressAndPort = c.ListenAddress + ":" + c.ListenPort

    c.LogPath, _ = file.Get("log", "log_path")
    c.LogLevel, _ = file.Get("log", "level")

    c.RootPath, _ = file.Get("main", "root_path")

    // Setting a pointer to our Config data
    conf = c

    return conf
}

func GetConf() (*Config){
    if conf != nil {
        log.Debug("Accessing configuration data structure.")
        return conf
    } else {
        log.Critical("Error configuration has not been loaded yet, run zutil.LoadConfig()")
        panic("Error configuration has not been loaded yet, run zutil.LoadConfig()")
    }
}
