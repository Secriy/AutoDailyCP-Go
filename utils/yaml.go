package utils

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config 配置文件定义
type Config struct {
	Host string `yaml:"Host"`
	Key  string `yaml:"Key"`
	Auth []struct {
		Account        string `yaml:"Account"`
		Password       string `yaml:"Password"`
		SchoolName     string `yaml:"SchoolName"`
		SignAddress    string `yaml:"SignAddress"`
		CollectAddress string `yaml:"CollectAddress"`
	} `yaml:"Auth"`
}

// ReadYamlConfig 读取yaml配置
func ReadYamlConfig(path string) (*Config, error) {
	conf := &Config{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_ = yaml.NewDecoder(f).Decode(conf)

	return conf, nil
}
