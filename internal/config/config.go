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
	MaxAttempts int    `yaml: "max_attempts" env-defalt: "5"`
	HTTPServer  `yaml:http_server`
	PgSQL
}
type HTTPServer struct {
	Address     string        `yaml: "adderss" env-defalt "0.0.0.0:8080"`
	Timeout     time.Duration `yaml: "timeout" env-defalt "5s"`
	IdleTimeout time.Duration `yaml: "idle_timeout" env-defalt "60s"`
}

type PgSQL struct {
	User          string `env: "POSTGRES_USER" env-required:"true"`
	Password      string `env: "POSTGRES_PASSWORD" env-required:"true"`
	Host          string `env: "POSTGRES_HOST" env-default: "postgres"`
	NameDB        string `env: "POSTGRES_DB" env-defalt: "dataBase"`
	Port          string `env: "POSTGRES_PORT" envDefault: "migrations"`
	MigrationPath string `env: "MIGRATIONS" env-required: "true"`
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
func (p *PgSQL) PGLDsetination() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.NameDB,
	)
}
