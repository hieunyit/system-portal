package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	"system-portal/pkg/logger"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Redis    RedisConfig    `mapstructure:"redis"` // NEW: Redis configuration
	Security SecurityConfig `mapstructure:"security"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	Timeout int    `mapstructure:"timeout"`
}

// Database connection settings
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

// DSN builds a PostgreSQL connection string from the database config.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}

type LoggerConfig = logger.LoggerConfig

type JWTConfig struct {
	// Legacy HMAC configuration (deprecated)
	Secret        string `mapstructure:"secret"`
	RefreshSecret string `mapstructure:"refreshSecret"`

	// RSA configuration (recommended)
	UseRSA                     bool          `mapstructure:"useRSA"`
	AccessPrivateKeyPath       string        `mapstructure:"accessPrivateKeyPath"`
	RefreshPrivateKeyPath      string        `mapstructure:"refreshPrivateKeyPath"`
	AccessTokenExpireDuration  time.Duration `mapstructure:"accessTokenExpireDuration"`
	RefreshTokenExpireDuration time.Duration `mapstructure:"refreshTokenExpireDuration"`
}

// Redis configuration
type RedisConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	Database int           `mapstructure:"database"`
	PoolSize int           `mapstructure:"poolSize"`
	TTL      time.Duration `mapstructure:"ttl"`
}

// Security configuration including CORS settings
type SecurityConfig struct {
	EnableSecurityHeaders bool       `mapstructure:"enableSecurityHeaders"`
	CORS                  CORSConfig `mapstructure:"cors"`
	EncryptionKey         string     `mapstructure:"encryptionKey"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowedOrigins"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
}

func Load() (*Config, error) {
	var cfg Config

	// Set config file path based on environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName("config-" + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Set defaults
	setDefaults()

	// Read environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Override port if PORT env variable is set
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}

	return &cfg, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.timeout", 30)

	// Logger defaults
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")

	// JWT defaults
	viper.SetDefault("jwt.useRSA", true)
	viper.SetDefault("jwt.accessPrivateKeyPath", "")
	viper.SetDefault("jwt.refreshPrivateKeyPath", "")
	viper.SetDefault("jwt.accessTokenExpireDuration", time.Hour)
	viper.SetDefault("jwt.refreshTokenExpireDuration", 24*time.Hour)

	// Legacy HMAC defaults (for backward compatibility)
	viper.SetDefault("jwt.secret", "default-hmac-secret-change-in-production")
	viper.SetDefault("jwt.refreshSecret", "default-refresh-secret-change-in-production")

	// Redis defaults
	viper.SetDefault("redis.enabled", false)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.database", 0)
	viper.SetDefault("redis.poolSize", 10)
	viper.SetDefault("redis.ttl", 10*time.Minute)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "user")
	viper.SetDefault("database.password", "pass")
	viper.SetDefault("database.name", "system_portal")
	viper.SetDefault("database.sslmode", "disable")

	// Security defaults
	viper.SetDefault("security.enableSecurityHeaders", true)
	viper.SetDefault("security.cors.allowedOrigins", []string{"*"})
	viper.SetDefault("security.cors.allowedMethods", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	viper.SetDefault("security.cors.allowedHeaders", []string{"Authorization", "Content-Type"})
	viper.SetDefault("security.cors.allowCredentials", true)
	viper.SetDefault("security.encryptionKey", "")
}
