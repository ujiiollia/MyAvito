package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml: "env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	JWTkey      string `yaml: "jwt_key" env-required:"true"`
	HTTPServer  `yaml:http_server`
	PgSQL       `yaml:pgl`
}
type HTTPServer struct {
	Address     string        `yaml: "adderss" env-defalt "0.0.0.0:8080"`
	Timeout     time.Duration `yaml: "timeout" env-defalt "5s"`
	IdleTimeout time.Duration `yaml: "idle_timeout" env-defalt "60s"`
}

type PgSQL struct {
	User          string `yaml: "POSTGRES_USER" env-required:"true"`
	Password      string `yaml: "POSTGRES_PASSWORD" env-required:"true"`
	NameDB        string `yaml: "POSTGRES_DB" env-defalt:"dataBase"`
	Port          int    `yaml: "POSTGRES_PORT" envDefault:"migrations"`
	MigrationPath string `yaml: "MIGRATIONS" env-required:"true"`
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
func (c *PgSQL) PGLDsetination() string {
	return fmt.Sprintf("postgres://%s:%s@postgres:%d/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Port,
		c.NameDB,
	)
}
