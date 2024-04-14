package config

import (
	"fmt"
	"time"
)

type Config struct {
	Env              string        `env:"ENV" envDefault:"local"`
	MaxAttempts      int           `env:"max_attempts" envDefault:"5"`
	PostgresUser     string        `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPassword string        `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	PostgresDB       string        `env:"POSTGRES_DB" envDefault:"postgres"`
	PostgresHost     string        `env:"POSTGRES_HOST" envDefault:"0.0.0.0"`
	PostgresPort     int           `env:"POSTGRES_PORT" envDefault:"5432"`
	ServicePort      int           `env:"SERVICE_PORT" envDefault:"8080"`
	Migrations       string        `env:"MIGRATIONS" envDefault:"migrations"`
	HTTPAddress      string        `env:"adderss" envDefault:"0.0.0.0:8080"`
	HTTPTimeout      time.Duration `env:"timeout" envDefault:"5s"`
	HTTPIdleTimeout  time.Duration `env:"idle_timeout" envDefault:"60s"`
}

// func MustLoad() *Config {
// 	path := fetchConfigPath()
// 	if path == "" {
// 		panic("config path is empty")
// 	}
// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		panic("config path is not exist: " + path)
// 	}
// 	var cfg Config
// 	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
// 		panic("failed to read config")
// 	}
// 	return &cfg
// }

//	func fetchConfigPath() string {
//		var res string
//		flag.StringVar(&res, "config", "", "path to config file")
//		flag.Parse()
//		if res == "" {
//			res = os.Getenv("CONFIG_PATH")
//		}
//		return res
//	}
func (c *Config) PGLDsetination() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
	)
}
