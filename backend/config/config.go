package config

import (
	"bytes"
	_ "embed"
	"github.com/spf13/viper"
	"log"
)

//go:embed config.yaml
var configFile []byte

// Config 用来存储所有的配置信息
type Config struct {
	Database struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`
	Tcos struct {
		BucketName string `mapstructure:"bucketName"`
		Region     string `mapstructure:"region"`
		Host       string `mapstructure:"host"`
	}
	Rds struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
	}
	AliYunAi struct {
		ApiKey string `mapstructure:"apiKey"`
	}
}

var config *Config

func init() {
	viper.SetConfigType("yaml") // 设置配置文件类型为 yaml

	// 使用嵌入的 configFile 加载配置
	if err := viper.ReadConfig(bytes.NewBuffer(configFile)); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	//配置文件映射到结构体
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Uable to decode into struct, %v", err)
	}
}

func LoadConfig() *Config {
	return config
}
