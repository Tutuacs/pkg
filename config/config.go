package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/Tutuacs/pkg/logs"
)

type config struct {
	apiConfig
	dbConfig
	jwtConfig
}

type apiConfig struct {
	Port string
}

type dbConfig struct {
	Host string
	Port string
	Addr string
	Name string
	User string
	Pass string
}

type jwtConfig struct {
	JWT_EXP    int64
	JWT_SECRET string
}

var cfg *config

func init() {
	logs.MessageLog("Loading configs...")
	cfg = new(config)
	cfg = defaultConfig()
}

func defaultConfig() *config {
	godotenv.Load()
	return &config{
		apiConfig: apiConfig{
			Port: getEnv("API_PORT", ":9000"),
		},
		dbConfig: dbConfig{
			Host: getEnv("DB_HOST", "127.0.0.1"),
			Port: getEnv("DB_PORT", "9999"),
			Addr: getEnv("DB_ADDR", "127.0.0.1:9999"),
			User: getEnv("DB_USER", "user"),
			Pass: getEnv("DB_PASS", "pass"),
			Name: getEnv("DB_NAME", "defaultDb"),
		},
		jwtConfig: jwtConfig{
			JWT_EXP:    getNumberEnv("JWT_EXP", 3600*24*7),
			JWT_SECRET: getEnv("JWT_SECRET", "secret"),
		},
	}
}

func getEnv(key string, defaultInfo string) string {

	info, ok := os.LookupEnv(key)
	if !ok && len(info) != 0 {
		return info
	}

	return defaultInfo
}

func getNumberEnv(key string, defaultInfo int64) int64 {

	if info, ok := os.LookupEnv(key); ok {
		number, err := strconv.ParseInt(info, 10, 64)
		if err != nil || number == 0 {
			return defaultInfo
		}
		return number
	}
	return defaultInfo
}

func GetAPI() apiConfig {
	return cfg.apiConfig
}

func GetDB() dbConfig {
	return cfg.dbConfig
}

func GetJWT() jwtConfig {
	return cfg.jwtConfig
}
