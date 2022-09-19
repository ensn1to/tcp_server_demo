package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type Config struct {
	HttpAddr string `json:"http-addr" mapstructure:"http-addr"`
	TcpAddr  string `json:"tcp-addr" mapstructure:"tcp-addr"`
}

func ReadFromEnv() {
	viper.AutomaticEnv()
}

// 从配置文件中读取
func ReadFromFile() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetConfigName("tcpserver") // 配置文件名称

	if err := viper.ReadInConfig(); err != nil {
		panic("Fatal error config")
	}

	fmt.Println("Used cfg file is:", viper.ConfigFileUsed())
}

// 监听和重读配置文件
func HotReadCfg() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})
}

func ReadFromConsul(serverMode string) (*Config, error) {
	viper.AddRemoteProvider("consul", "localhost:8500", serverMode)
	viper.SetConfigType("json")
	err := viper.ReadRemoteConfig()
	if err != nil {
		return nil, err
	}

	var c Config

	if viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
