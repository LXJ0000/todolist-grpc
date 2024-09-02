package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	RedisAddr              string `mapstructure:"REDIS_ADDR"`
	RedisExpiration        int    `mapstructure:"REDIS_EXPIRATION"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
	RedisDB                int    `mapstructure:"REDIS_DB"`
	RateLimit              int    `mapstructure:"RATE_LIMIT"`
	RateLimitWindow        int    `mapstructure:"RATE_LIMIT_WINDOW"`
	SnowflakeStartTime     string `mapstructure:"SNOWFLAKE_START_TIME"`
	SnowflakeMachineID     int64  `mapstructure:"SNOWFLAKE_MACHINE_ID"`
	KafkaAddr              string `mapstructure:"KAFKA_ADDR"`
	PrometheusAddress      string `mapstructure:"PROMETHEUS_ADDRESS"`
	MySQLAddress           string `mapstructure:"MYSQL_ADDRESS"`
	LocalStaticPath        string `mapstructure:"LOCAL_STATIC_PATH"`
	UrlStaticPath          string `mapstructure:"URL_STATIC_PATH"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
