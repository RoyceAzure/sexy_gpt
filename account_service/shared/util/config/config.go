package config

import (
	"time"

	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	ServiceID               string        `mapstructure:"SERVICEID"`
	Enviornmant             string        `mapstructure:"ENVIRONMENT"`
	DBDriver                string        `mapstructure:"DB_DRIVER"`
	DBSource                string        `mapstructure:"DB_SOURCE"`
	MigrateURL              string        `mapstructure:"MIGRATE_URL"`
	GATEWAY_ACCOUNT_SERVICE string        `mapstructure:"GATEWAY_ACCOUNT_SERVICE"`
	GRPCAccountServiceAddr  string        `mapstructure:"GRPC_ACCOUNT_SERVICE"`
	TokenSymmetricKey       string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration     time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration    time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisQueueAddress       string        `mapstructure:"REDIS_Q_ADDRESS"`
	EmailSenderName         string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress      string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword     string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	AUTH_AUDIENCE           string        `mapstructure:"AUTH_AUDIENCE"`
	AUTH_ISSUER             string        `mapstructure:"AUTH_ISSUER"`
	FRONTEND_KEY            string        `mapstructure:"FRONTEND_KEY"`
}

/*
path: app.env所在的資料夾
*/
func LoadConfig(path string) (config Config, err error) {
	if AppConfig != nil {
		return *AppConfig, nil
	}
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") //JSON XML  這是指extension

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	if err == nil {
		AppConfig = &config
	}
	return
}
