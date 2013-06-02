package zconfig

type Config struct {
    ListenPort  int
    BindAddress string
    Log         string
}


func LoadConfig(confFile string) Config {
    return Config{9000, "0.0.0.0", "zagent.log"}
}
