package configs

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds the entire application configuration.
type Config struct {
	App         AppConfig         `mapstructure:",squash"`
	Server      ServerConfig      `mapstructure:",squash"`
	Database    DatabaseConfig    `mapstructure:",squash"`
	Redis       RedisConfig       `mapstructure:",squash"`
	RabbitMQ    RabbitMQConfig    `mapstructure:",squash"`
	Kafka       KafkaConfig       `mapstructure:",squash"`
	JWT         JWTConfig         `mapstructure:",squash"`
	SMTP        SMTPConfig        `mapstructure:",squash"`
	Mailgun     MailgunConfig     `mapstructure:",squash"`
	SendGrid    SendGridConfig    `mapstructure:",squash"`
	Storage     StorageConfig     `mapstructure:",squash"`
	AWS         AWSConfig         `mapstructure:",squash"`
	R2          R2Config          `mapstructure:",squash"`
	RateLimit   RateLimitConfig   `mapstructure:",squash"`
	CORS        CORSConfig        `mapstructure:",squash"`
	OTP         OTPConfig         `mapstructure:",squash"`
	Mail        MailConfig        `mapstructure:",squash"`
}

type AppConfig struct {
	Name        string `mapstructure:"APP_NAME"`
	Env         string `mapstructure:"APP_ENV"`
	Debug       bool   `mapstructure:"APP_DEBUG"`
	Version     string `mapstructure:"APP_VERSION"`
	URL         string `mapstructure:"APP_URL"`
	FrontendURL string `mapstructure:"FRONTEND_URL"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"SERVER_HOST"`
	Port            string        `mapstructure:"SERVER_PORT"`
	ReadTimeout     time.Duration `mapstructure:"-"`
	WriteTimeout    time.Duration `mapstructure:"-"`
	IdleTimeout     time.Duration `mapstructure:"-"`
	ShutdownTimeout time.Duration `mapstructure:"-"`
	ReadTimeoutSec  int           `mapstructure:"SERVER_READ_TIMEOUT"`
	WriteTimeoutSec int           `mapstructure:"SERVER_WRITE_TIMEOUT"`
	IdleTimeoutSec  int           `mapstructure:"SERVER_IDLE_TIMEOUT"`
	ShutdownSec     int           `mapstructure:"SERVER_SHUTDOWN_TIMEOUT"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"DB_HOST"`
	Port            string `mapstructure:"DB_PORT"`
	User            string `mapstructure:"DB_USER"`
	Password        string `mapstructure:"DB_PASSWORD"`
	Name            string `mapstructure:"DB_NAME"`
	SSLMode         string `mapstructure:"DB_SSLMODE"`
	Timezone        string `mapstructure:"DB_TIMEZONE"`
	MaxOpenConns    int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime int    `mapstructure:"DB_CONN_MAX_LIFETIME"`
	LogLevel        string `mapstructure:"DB_LOG_LEVEL"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		d.Host, d.User, d.Password, d.Name, d.Port, d.SSLMode, d.Timezone,
	)
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
	PoolSize int    `mapstructure:"REDIS_POOL_SIZE"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

type RabbitMQConfig struct {
	Host          string `mapstructure:"RABBITMQ_HOST"`
	Port          string `mapstructure:"RABBITMQ_PORT"`
	User          string `mapstructure:"RABBITMQ_USER"`
	Password      string `mapstructure:"RABBITMQ_PASSWORD"`
	VHost         string `mapstructure:"RABBITMQ_VHOST"`
	PrefetchCount int    `mapstructure:"RABBITMQ_PREFETCH_COUNT"`
}

func (r *RabbitMQConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s%s", r.User, r.Password, r.Host, r.Port, r.VHost)
}

type KafkaConfig struct {
	Brokers         string `mapstructure:"KAFKA_BROKERS"`
	GroupID         string `mapstructure:"KAFKA_GROUP_ID"`
	AutoOffsetReset string `mapstructure:"KAFKA_AUTO_OFFSET_RESET"`
}

func (k *KafkaConfig) BrokerList() []string {
	return strings.Split(k.Brokers, ",")
}

type JWTConfig struct {
	AccessSecret  string        `mapstructure:"JWT_ACCESS_SECRET"`
	RefreshSecret string        `mapstructure:"JWT_REFRESH_SECRET"`
	AccessExpiry  time.Duration `mapstructure:"-"`
	RefreshExpiry time.Duration `mapstructure:"-"`
	AccessMin     int           `mapstructure:"JWT_ACCESS_EXPIRY"`
	RefreshMin    int           `mapstructure:"JWT_REFRESH_EXPIRY"`
	Issuer        string        `mapstructure:"JWT_ISSUER"`
}

type MailConfig struct {
	Provider string `mapstructure:"MAIL_PROVIDER"`
}

type SMTPConfig struct {
	Host       string `mapstructure:"SMTP_HOST"`
	Port       int    `mapstructure:"SMTP_PORT"`
	User       string `mapstructure:"SMTP_USER"`
	Password   string `mapstructure:"SMTP_PASSWORD"`
	From       string `mapstructure:"SMTP_FROM"`
	FromName   string `mapstructure:"SMTP_FROM_NAME"`
	Encryption string `mapstructure:"SMTP_ENCRYPTION"`
}

type MailgunConfig struct {
	Domain string `mapstructure:"MAILGUN_DOMAIN"`
	APIKey string `mapstructure:"MAILGUN_API_KEY"`
	From   string `mapstructure:"MAILGUN_FROM"`
}

type SendGridConfig struct {
	APIKey   string `mapstructure:"SENDGRID_API_KEY"`
	From     string `mapstructure:"SENDGRID_FROM"`
	FromName string `mapstructure:"SENDGRID_FROM_NAME"`
}

type StorageConfig struct {
	Provider    string `mapstructure:"STORAGE_PROVIDER"`
	LocalPath   string `mapstructure:"STORAGE_LOCAL_PATH"`
	MaxFileSize int64  `mapstructure:"STORAGE_MAX_FILE_SIZE"`
}

type AWSConfig struct {
	Region          string `mapstructure:"AWS_REGION"`
	AccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey  string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	S3Bucket        string `mapstructure:"AWS_S3_BUCKET"`
	S3Endpoint      string `mapstructure:"AWS_S3_ENDPOINT"`
}

type R2Config struct {
	AccountID       string `mapstructure:"R2_ACCOUNT_ID"`
	AccessKeyID     string `mapstructure:"R2_ACCESS_KEY_ID"`
	SecretAccessKey  string `mapstructure:"R2_SECRET_ACCESS_KEY"`
	Bucket          string `mapstructure:"R2_BUCKET"`
	PublicURL       string `mapstructure:"R2_PUBLIC_URL"`
}

type RateLimitConfig struct {
	Requests int `mapstructure:"RATE_LIMIT_REQUESTS"`
	Duration int `mapstructure:"RATE_LIMIT_DURATION"`
}

type CORSConfig struct {
	AllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	AllowedMethods string `mapstructure:"CORS_ALLOWED_METHODS"`
	AllowedHeaders string `mapstructure:"CORS_ALLOWED_HEADERS"`
	MaxAge         int    `mapstructure:"CORS_MAX_AGE"`
}

func (c *CORSConfig) Origins() []string {
	return strings.Split(c.AllowedOrigins, ",")
}

func (c *CORSConfig) Methods() []string {
	return strings.Split(c.AllowedMethods, ",")
}

func (c *CORSConfig) Headers() []string {
	return strings.Split(c.AllowedHeaders, ",")
}

type OTPConfig struct {
	Length      int `mapstructure:"OTP_LENGTH"`
	Expiry      int `mapstructure:"OTP_EXPIRY"`
	MaxAttempts int `mapstructure:"OTP_MAX_ATTEMPTS"`
}

// Load reads configuration from environment files and environment variables.
func Load() *Config {
	v := viper.New()

	v.SetConfigType("env")
	v.AutomaticEnv()

	// Load .env files in order of precedence (last wins)
	envFiles := []string{".env", ".env.development", ".env.production", ".env.local"}
	for _, f := range envFiles {
		v.SetConfigFile(f)
		if err := v.MergeInConfig(); err != nil {
			// It's okay if some env files don't exist
			continue
		}
	}

	cfg := &Config{}

	if err := v.Unmarshal(&cfg.App); err != nil {
		log.Fatalf("Failed to unmarshal App config: %v", err)
	}
	if err := v.Unmarshal(&cfg.Server); err != nil {
		log.Fatalf("Failed to unmarshal Server config: %v", err)
	}
	if err := v.Unmarshal(&cfg.Database); err != nil {
		log.Fatalf("Failed to unmarshal Database config: %v", err)
	}
	if err := v.Unmarshal(&cfg.Redis); err != nil {
		log.Fatalf("Failed to unmarshal Redis config: %v", err)
	}
	if err := v.Unmarshal(&cfg.RabbitMQ); err != nil {
		log.Fatalf("Failed to unmarshal RabbitMQ config: %v", err)
	}
	if err := v.Unmarshal(&cfg.Kafka); err != nil {
		log.Fatalf("Failed to unmarshal Kafka config: %v", err)
	}
	if err := v.Unmarshal(&cfg.JWT); err != nil {
		log.Fatalf("Failed to unmarshal JWT config: %v", err)
	}
	if err := v.Unmarshal(&cfg.Mail); err != nil {
		log.Fatalf("Failed to unmarshal Mail config: %v", err)
	}
	if err := v.Unmarshal(&cfg.SMTP); err != nil {
		log.Fatalf("Failed to unmarshal SMTP config: %v", err)
	}
	if err := v.Unmarshal(&cfg.Mailgun); err != nil {
		log.Fatalf("Failed to unmarshal Mailgun config: %v", err)
	}
	if err := v.Unmarshal(&cfg.SendGrid); err != nil {
		log.Fatalf("Failed to unmarshal SendGrid config: %v", err)
	}
	if err := v.Unmarshal(&cfg.Storage); err != nil {
		log.Fatalf("Failed to unmarshal Storage config: %v", err)
	}
	if err := v.Unmarshal(&cfg.AWS); err != nil {
		log.Fatalf("Failed to unmarshal AWS config: %v", err)
	}
	if err := v.Unmarshal(&cfg.R2); err != nil {
		log.Fatalf("Failed to unmarshal R2 config: %v", err)
	}
	if err := v.Unmarshal(&cfg.RateLimit); err != nil {
		log.Fatalf("Failed to unmarshal RateLimit config: %v", err)
	}
	if err := v.Unmarshal(&cfg.CORS); err != nil {
		log.Fatalf("Failed to unmarshal CORS config: %v", err)
	}
	if err := v.Unmarshal(&cfg.OTP); err != nil {
		log.Fatalf("Failed to unmarshal OTP config: %v", err)
	}

	// Compute durations from seconds/minutes
	cfg.Server.ReadTimeout = time.Duration(cfg.Server.ReadTimeoutSec) * time.Second
	cfg.Server.WriteTimeout = time.Duration(cfg.Server.WriteTimeoutSec) * time.Second
	cfg.Server.IdleTimeout = time.Duration(cfg.Server.IdleTimeoutSec) * time.Second
	cfg.Server.ShutdownTimeout = time.Duration(cfg.Server.ShutdownSec) * time.Second
	cfg.JWT.AccessExpiry = time.Duration(cfg.JWT.AccessMin) * time.Minute
	cfg.JWT.RefreshExpiry = time.Duration(cfg.JWT.RefreshMin) * time.Minute

	return cfg
}

// IsDevelopment returns true if the application is in development mode.
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction returns true if the application is in production mode.
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}
