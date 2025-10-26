package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	WriteDSN     string = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	ReadDSN      string = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	MaxOpenConns int    = 10
	MaxIdleConns int    = 10
	HttpHost     string = "0.0.0.0"
	HttpPort     int    = 8080
)

type ValueType interface {
	~string | ~int | ~int64 | ~bool | ~float64
}

type Config struct {
	WriteDSN        string
	ReadDSN         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int64
	HttpHost        string
	HttpPort        int
	Logger          *slog.Logger
}

func NewConfig() *Config {
	godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	slog.Info("loading config from environment variables")

	return &Config{
		WriteDSN:     GetEnv("DATABASE_URL_WRITE", WriteDSN),
		ReadDSN:      GetEnv("DATABASE_URL_READ", ReadDSN),
		MaxOpenConns: GetEnv("DB_MAX_CONN", MaxOpenConns),
		MaxIdleConns: GetEnv("DB_MAX_IDLE_CONN", MaxIdleConns),
		HttpHost:     GetEnv("HTTP_HOST", HttpHost),
		HttpPort:     GetEnv("HTTP_PORT", HttpPort),
		Logger:       logger,
	}
}

func GetEnv[T ValueType](key string, defaultValue T) T {
	v, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	switch any(defaultValue).(type) {
	case string:
		return any(v).(T)
	case int:
		if i, err := strconv.Atoi(v); err == nil {
			return any(i).(T)
		}
	case int64:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return any(i).(T)
		}
	case bool:
		if b, err := strconv.ParseBool(v); err == nil {
			return any(b).(T)
		}
	case float64:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return any(f).(T)
		}
	default:
		return defaultValue
	}
	return defaultValue
}
