package main

import (
	"net/http"
	"net/http/cookiejar"
	"os"

	"AutoDailyCP-Go/scripts"
	"golang.org/x/net/publicsuffix"
)

var config Config
var script scripts.Script

func init() {
	// Load configuration
	config.ReadYamlConfig("/config.yaml")
}

// Execute
func exec(task Task) {
	script = scripts.Script{
		Host:        config.Host,
		Key:         config.Key,
		SignAddr:    task.SignAddress,
		CollectAddr: task.CollectAddress,
	}
	// Cookie Storage
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client := &http.Client{Jar: cookieJar}
	// Login
	if !task.Login(client) {
		os.Exit(2)
	}
	// Sign
	var sign scripts.Sign
	sign.Script = script
	sign.DoSign(client)
	// Collect
	var collect scripts.Collect
	collect.Script = script
	collect.DoCollect(client)

}

func main() {
	for _, task := range config.Tasks {
		exec(task)
	}
}
