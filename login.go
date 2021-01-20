package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"AutoDailyCP-Go/utils"
)

// Response JsonStruct of Responses
type response struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Result  result `json:"result"`
}

// Result JsonStruct of Result
type result struct {
	EncryptSalt  string `json:"_encryptSalt"`
	Lt           string `json:"_lt"`
	ForgetPwdURL string `json:"forgetPwdUrl"`
	NeedCaptcha  string `json:"needCaptcha"`
}

// Success JsonStruct of Success
type success struct {
	ResultCode string `json:"resultCode"`
	URL        string `json:"url"`
}

// Login 登录
func Login(host string, client *http.Client, username string, password string) string {
	foreURL := fmt.Sprintf("https://%v/iap/login?service=https://mobile.campushoy.com/v6/auth/campus/cas/login", host)
	ltURL := fmt.Sprintf("https://%v/iap/security/lt", host)
	loginURL := fmt.Sprintf("https://%v/iap/doLogin", host)

	// Service
	req, _ := http.NewRequest("GET", foreURL, nil)
	req.Header.Add("CpdailyAuthType", "Login")
	res, _ := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		return utils.Message("Login Error: " + "Service Response Empty.")
	}
	var lt = strings.Split(res.Request.URL.String(), "=")[1]
	lt = strings.Replace(lt, "&isCpdaily", "", -1)

	// Security
	req, _ = http.NewRequest("POST", ltURL, strings.NewReader(fmt.Sprintf("lt=%v", lt)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, _ = client.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		return utils.Message("Login Error: " + "Security Response Empty.")
	}
	body, _ := ioutil.ReadAll(res.Body)
	var response = response{}
	_ = json.Unmarshal(body, &response)
	lt = response.Result.Lt

	// DoLogin
	data := url.Values{
		"username":   {username},
		"password":   {password},
		"lt":         {lt},
		"dllt":       {"cpdaily"},
		"mobile":     {""},
		"captcha":    {""},
		"rememberMe": {"false"},
	}
	req, _ = http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 (4789933056)cpdaily/8.2.4  wisedu/8.2.4")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, _ = client.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		return utils.Message("Login Error: " + "Security Response Empty.")
	}
	out, _ := ioutil.ReadAll(res.Body)

	var success = success{}
	_ = json.Unmarshal(out, &success)

	// Redirect
	if success.ResultCode == "REDIRECT" {
		authURL := fmt.Sprintf("https://%v/portal/login", host) + strings.Split(success.URL, "?")[1]
		res, _ = client.Get(authURL)
		return ""
	}

	return utils.Message("Login Error: " + success.ResultCode)
}
