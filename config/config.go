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
	redisConfig
	mqttConfig
	smtpConfig
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

type redisConfig struct {
	Addr string
	// Password string
}

type mqttConfig struct {
	Addr string
}

type smtpConfig struct {
	SMTP_MAIL string
	SMTP_PASS string
	SMTP_HOST string
	SMTP_ADDR string
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
		redisConfig: redisConfig{
			Addr: getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		},
		mqttConfig: mqttConfig{
			Addr: getEnv("MQTT_ADDR", "127.0.0.1:1883"),
		},
		smtpConfig: smtpConfig{
			SMTP_MAIL: getEnv("SMTP_MAIL", "arthursilva.mailtest@gmail.com"),
			SMTP_PASS: getEnv("SMTP_PASS", "xcyezdmrqithcyuo"),
			SMTP_HOST: getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTP_ADDR: getEnv("SMTP_ADDR", "smtp.gmail.com:587"),
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

func GetRedis() redisConfig {
	return cfg.redisConfig
}

func GetMqtt() mqttConfig {
	return cfg.mqttConfig
}

func GetMailer() smtpConfig {
	return cfg.smtpConfig
}
