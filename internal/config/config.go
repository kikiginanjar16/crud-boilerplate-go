package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	AppEnv        string
	AppPort       int
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration

	DBHost     string
	DBPort     int
	DBUser     string
	DBPass     string
	DBName     string
	DBSSLMode  string
	DBTimezone string

	JWTSecret       string
	JWTExpireMinute int

	UploadDir string
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func atoi(envKey string, fallback int) int {
	if v := os.Getenv(envKey); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Printf("invalid int for %s=%s, using %d", envKey, v, fallback)
	}
	return fallback
}

func durationFromSeconds(envKey string, fallback int) time.Duration {
	return time.Duration(atoi(envKey, fallback)) * time.Second
}

func ensureDir(path string) {
	if err := os.MkdirAll(path, 0o755); err != nil {
		log.Printf("warn: cannot create dir %s: %v", path, err)
	}
}

func Load() *Config {
	cfg := &Config{
		AppEnv:       getenv("APP_ENV", "development"),
		AppPort:      atoi("APP_PORT", 8080),
		ReadTimeout:  durationFromSeconds("APP_READ_TIMEOUT", 5),
		WriteTimeout: durationFromSeconds("APP_WRITE_TIMEOUT", 10),

		DBHost:     getenv("DB_HOST", "localhost"),
		DBPort:     atoi("DB_PORT", 5432),
		DBUser:     getenv("DB_USER", "postgres"),
		DBPass:     getenv("DB_PASS", "postgres"),
		DBName:     getenv("DB_NAME", "appdb"),
		DBSSLMode:  getenv("DB_SSLMODE", "disable"),
		DBTimezone: getenv("DB_TIMEZONE", "Asia/Jakarta"),

		JWTSecret:       getenv("JWT_SECRET", "supersecretchangeme"),
		JWTExpireMinute: atoi("JWT_EXPIRE_MINUTES", 60),

		UploadDir: getenv("UPLOAD_DIR", "./uploads"),
	}

	// Normalize relative path
	if !filepath.IsAbs(cfg.UploadDir) {
		if wd, err := os.Getwd(); err == nil {
			cfg.UploadDir = filepath.Join(wd, cfg.UploadDir)
		}
	}
	ensureDir(cfg.UploadDir)
	return cfg
}
