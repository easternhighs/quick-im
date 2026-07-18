package config

import (
	"log"
	"os"

	"go.yaml.in/yaml/v2"
)

// 总配置内容
type Config struct {
	App   AppConfig   `yaml:"app"`
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
	JWT   JWTConfig   `yaml:"jwt"`
}

type AppConfig struct {
	Env      string `yaml:"env"`
	HTTPAddr string `yaml:"http_addr"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Channel  string `yaml:"channel"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

func ReadConfig(configFileName string) Config {
	//打开配置文件
	configFile, ReadConfigErr := os.OpenFile(configFileName, os.O_RDONLY, 0666)
	if ReadConfigErr != nil {
		log.Fatal("failed to open config file:", ReadConfigErr)
	}
	defer configFile.Close()

	//读取配置文件内容
	config := Config{}
	configFile_byte, ReadConfigErr := os.ReadFile(configFileName)
	if ReadConfigErr != nil {
		log.Fatal("failed to read config file:", ReadConfigErr)
	}
	unmarshalErr := yaml.Unmarshal(configFile_byte, &config)
	if unmarshalErr != nil {
		log.Fatal("failed to unmarshal config file:", unmarshalErr)
	}

	return config
}
