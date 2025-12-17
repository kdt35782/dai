package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Crypto   CryptoConfig   `mapstructure:"crypto"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int    `mapstructure:"expires_in"`
}

type CryptoConfig struct {
	SM4Key         string `mapstructure:"sm4_key"`
	SM2PrivateKey  string `mapstructure:"sm2_private_key"`
	SM2PublicKey   string `mapstructure:"sm2_public_key"`
}

type UploadConfig struct {
	MaxSize      int64    `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
	UploadPath   string   `mapstructure:"upload_path"`
}

var AppConfig *Config

// LoadConfig 加载配置文件
func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	AppConfig = &config
	return &config
}
