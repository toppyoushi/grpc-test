package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type serverConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

//Config 配置
type Config struct {
	Server serverConfig `yaml:"server"`
}

var _config *Config

//GetConfig 获取配置
func GetConfig() *Config {
	return _config
}

//Init 初始化
func Init(configPath string) error {
	f, err := ioutil.ReadFile(configPath)

	if err != nil {
		return err
	}
	_config = &Config{}
	return yaml.Unmarshal(f, _config)
}
