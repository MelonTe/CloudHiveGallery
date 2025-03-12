package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 用来存储所有的配置信息
type Config struct {
	Database struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`
}

var config *Config

func init() {
	//初始化config
	viper.AddConfigPath("./config")    //在当前文件夹下寻找
	viper.SetConfigName("config.yaml") //查找指定文件名
	viper.SetConfigType("yaml")        //当没有设置特定的文件后缀名时，必须要指定文件类型

	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
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
