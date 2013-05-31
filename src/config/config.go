package config

type Config struct {
    listenPort int
    bindAddress string
}


func LoadConfig(confFile string) Config {
    return Config{9000, "0.0.0.0"}
}
