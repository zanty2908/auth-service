package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env string `mapstructure:"ENV"`

	APIInfo struct {
		Service     string `mapstructure:"service"`
		Version     string `mapstructure:"version"`
		VersionPath string `mapstructure:"version_path"`
		HTTPPort    int    `mapstructure:"http_port"`
		GRPCPort    int    `mapstructure:"grpc_port"`
	} `mapstructure:"api_info"`

	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
	Auth     Auth     `mapstructure:"auth_config"`
}

type Database struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"database"`
	Schema   string `mapstructure:"schema"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`

	MaxIdleConnection    int           `mapstructure:"max_idle_conns"`
	MaxActiveConnection  int           `mapstructure:"max_active_conns"`
	MaxConnectionTimeout time.Duration `mapstructure:"max_conn_timeout"`

	DebugLog bool `mapstructure:"debug_log"`
}

type Redis struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Database     int           `mapstructure:"database"`
	RateLimit    int           `mapstructure:"rate_limit_database"`
	TTL          time.Duration `mapstructure:"ttl"`
	PoolSize     int           `mapstructure:"pool_size"`
	Password     string        `mapstructure:"password"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
}

type Auth struct {
	Token Token `mapstructure:"token"`

	OTP struct {
		ExpiresDuration         time.Duration `mapstructure:"expires_duration"`
		ResendDuration          time.Duration `mapstructure:"resend_duration"`
		EnteredIncorrectlyTimes int           `mapstructure:"entered_incorrectly_times"`
	} `mapstructure:"otp"`
}

type Token struct {
	SecretKey            string        `mapstructure:"secret_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

func (s *Config) DataSourceForSQL() string {
	dbConfig := s.Database
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SSLMode,
	)
}

func (s *Config) DataSourceForPGX() string {
	dbConfig := s.Database
	return fmt.Sprintf(
		"%s://%s:%s@%s:%v/%s?sslmode=%s",
		dbConfig.Driver,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.SSLMode,
	)
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (*Config, error) {
	env := viper.New()
	env.SetConfigName("config")
	env.AddConfigPath(".")          // Look for config in current directory
	env.AddConfigPath("config/")    // Optionally look for config in the working directory.
	env.AddConfigPath("../config/") // Look for config needed for tests.
	env.AddConfigPath("../")        // Look for config needed for tests.
	env.AddConfigPath(path)

	env.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	env.AutomaticEnv()

	err := env.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = env.Unmarshal(config)
	if err != nil { // Handle errors reading the config file
		return nil, err
	}

	return config, err
}
