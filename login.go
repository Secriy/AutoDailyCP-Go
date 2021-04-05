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

// loginRes JsonStruct of Success
type loginRes struct {
	ResultCode string `json:"resultCode"`
	URL        string `json:"url"`
}

// Login 登录
func (t Task) Login(c *http.Client) bool {
	host := config.Host
	foreURL := fmt.Sprintf("https://%v/iap/login?service=https://mobile.campushoy.com/v6/auth/campus/cas/login", host)
	ltURL := fmt.Sprintf("https://%v/iap/security/lt", host)
	loginURL := fmt.Sprintf("https://%v/iap/doLogin", host)

	// Service
	req, _ := http.NewRequest("GET", foreURL, nil)
	req.Header.Add("CpdailyAuthType", "Login")
	res, _ := c.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		utils.Log("LoginError").Message("Service Response Empty.")
		return false
	}
	var lt = strings.Split(res.Request.URL.String(), "=")[1]
	lt = strings.Replace(lt, "&isCpdaily", "", -1)

	// Security
	req, _ = http.NewRequest("POST", ltURL, strings.NewReader(fmt.Sprintf("lt=%v", lt)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, _ = c.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		utils.Log("LoginError").Message("Security Response Empty.")

		return false
	}
	body, _ := ioutil.ReadAll(res.Body)
	var response = response{}
	_ = json.Unmarshal(body, &response)
	lt = response.Result.Lt

	// DoLogin
	data := url.Values{
		"username":   {t.Account},
		"password":   {t.Password},
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
	res, _ = c.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		utils.Log("LoginError").Message("Security Response Empty.")

		return false
	}
	out, _ := ioutil.ReadAll(res.Body)

	var loginRes = loginRes{}
	_ = json.Unmarshal(out, &loginRes)

	// Redirect
	if loginRes.ResultCode == "REDIRECT" {
		authURL := fmt.Sprintf("https://%v/portal/login", host) + strings.Split(loginRes.URL, "?")[1]
		res, _ = c.Get(authURL)
		utils.Log("LoginInfo").Message("Success")

		return true
	}

	utils.Log("LoginError").Message(loginRes.ResultCode)

	return false
}
