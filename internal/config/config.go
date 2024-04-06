package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml: "env" env-default:"local"`
	StoragePath string `yaml: "storage_path" env-reqared "true"`
	HTTPServer  `yaml:http_server`
}
type HTTPServer struct {
	Address     string        `yaml: "adderss" env-defalt "0.0.0.0:8080"`
	Timeout     time.Duration `yaml: "timeout" env-defalt "5s"`
	IdleTimeout time.Duration `yaml: "idle_timeout" env-defalt "60s"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path is not exist: " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config")
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
