package settings

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Настройка заапуска приложения
type Config struct {
	ServiceName string `env:"SERVICE_NAME,required"`
	SecretKey   string `env:"SECRET_KEY,required"`
	ExpTokenDay int    `env:"EXP_TOKEN_DAY,required"`
	// способ заапуска
	Host string `env:"HOST"`
	Port string `env:"PORT"`

	// подключение db
	DB *DbSettings
}

func (ds *DbSettings) UrlPostgres() string {
	return fmt.Sprintf(`postgres://%s:%s@%s:%s/%s`, ds.User, ds.Passw, ds.Host, ds.Port, ds.DbName)
}

func Env(key string) (string, error) {
	ret := os.Getenv(key)
	if len(ret) == 0 {
		err := fmt.Errorf("%s env not find ", key)
		return "", err
	}
	return ret, nil
}

func InitEnv(fnames ...string) (*Config, error) {
	godotenv.Load(fnames...)
	cfg := &Config{
		DB: &DbSettings{},
	}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("unable to parse environment variables: %e", err)
	}
	return cfg, nil
}

// подк. к базам
type DbSettings struct {
	User   string `env:"DB_USER,required"`
	Passw  string `env:"DB_PASSW,required"`
	Host   string `env:"DB_HOST,required"`
	Port   string `env:"DB_PORT,required"`
	Schema string `env:"DB_SCHEMA,required"`
	DbName string `env:"DB_NAME,required"`
}
