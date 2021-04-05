package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"AutoDailyCP-Go/utils"
	"gopkg.in/yaml.v2"
)

type Task struct {
	Account        string `yaml:"Account"`
	Password       string `yaml:"Password"`
	SchoolName     string `yaml:"SchoolName"`
	SignAddress    string `yaml:"SignAddress"`
	CollectAddress string `yaml:"CollectAddress"`
}

// Config 配置文件定义
type Config struct {
	Host  string `yaml:"Host"`
	Key   string `yaml:"Key"`
	Tasks []Task `yaml:"Tasks"`
}

// GetCurrentDirectory 获取程序执行目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		utils.Log("ConfigError").Message("Error get the current directory of execution: " + err.Error())
		os.Exit(2)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// ReadYamlConfig 读取yaml配置
func (c *Config) ReadYamlConfig(path string) {
	path = GetCurrentDirectory() + path
	file, err := os.Open(path)
	if err != nil {
		utils.Log("ConfigError").Message("Error loading configuration file: " + err.Error())
		os.Exit(2)
	}
	content, _ := ioutil.ReadAll(file)
	_ = yaml.Unmarshal(content, &c)
	defer file.Close()
}
