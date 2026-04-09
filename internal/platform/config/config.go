package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppName  string
	GinMode  string
	HTTP     HTTPConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

type HTTPConfig struct {
	Port string
}

type DatabaseConfig struct {
	DSN           string
	MigrationsDir string
}

type LoggingConfig struct {
	ServiceName         string
	Output              string
	FilePath            string
	Level               string
	Format              string
	BodyMaxBytes        int
	IncludeRequestBody  bool
	IncludeResponseBody bool
	MaxSizeMB           int
	MaxBackups          int
	MaxAgeDays          int
	Compress            bool
}

type envValue interface {
	string | int | bool
}

func (c HTTPConfig) Addr() string {
	return ":" + c.Port
}

func (c DatabaseConfig) Enabled() bool {
	return c.DSN != ""
}

func Load() Config {
	return Config{
		AppName: getEnv("APP_NAME", "project-example"),
		GinMode: getEnv("GIN_MODE", "debug"),
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			DSN:           getEnv("DATABASE_DSN", ""),
			MigrationsDir: getEnv("DATABASE_MIGRATIONS_DIR", "migrations"),
		},
		Logging: LoggingConfig{
			ServiceName:         getEnv("LOG_SERVICE_NAME", getEnv("APP_NAME", "project-example")),
			Output:              getEnv("LOG_OUTPUT", "both"),
			FilePath:            getEnv("LOG_FILE_PATH", "logs/app.log"),
			Level:               getEnv("LOG_LEVEL", "info"),
			Format:              getEnv("LOG_FORMAT", "json"),
			BodyMaxBytes:        getEnv("LOG_BODY_MAX_BYTES", 4096),
			IncludeRequestBody:  getEnv("LOG_INCLUDE_REQUEST_BODY", false),
			IncludeResponseBody: getEnv("LOG_INCLUDE_RESPONSE_BODY", false),
			MaxSizeMB:           getEnv("LOG_MAX_SIZE_MB", 100),
			MaxBackups:          getEnv("LOG_MAX_BACKUPS", 3),
			MaxAgeDays:          getEnv("LOG_MAX_AGE_DAYS", 7),
			Compress:            getEnv("LOG_COMPRESS", false),
		},
	}
}

func getEnv[T envValue](key string, fallback T) T {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	switch any(fallback).(type) {
	case string:
		return any(value).(T)
	case int:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}

		return any(parsed).(T)
	case bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}

		return any(parsed).(T)
	default:
		return fallback
	}
}
