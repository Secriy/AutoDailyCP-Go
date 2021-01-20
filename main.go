package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"

	"AutoDailyCP-Go/scripts"
	"golang.org/x/net/publicsuffix"

	"AutoDailyCP-Go/utils"
)

var conf *utils.Config

func init() {
	// Load Configuration
	var err error
	conf, err = utils.ReadYamlConfig(utils.GetCurrentDirectory() + "/config.yaml")
	if err != nil {
		utils.Log().Error("Error loading configuration file: " + err.Error())
		os.Exit(2)
	}
}

func main() {
	for _, auth := range conf.Auth {

		// Cookie Storage
		cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		client := &http.Client{Jar: cookieJar}

		// Login
		loginPrint := Login(conf.Host, client, auth.Account, auth.Password)
		if loginPrint == "" {
			fmt.Println(utils.Message("Login Success"))
		} else {
			fmt.Println(loginPrint)
		}

		// Sign
		scripts.DoSign(conf.Host, conf.Key, client, auth.SignAddress)

		// Collect
		scripts.DoCollect(conf.Host, conf.Key, client, auth.CollectAddress)
	}
}
