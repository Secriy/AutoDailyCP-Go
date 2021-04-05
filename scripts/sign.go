package scripts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"AutoDailyCP-Go/utils"
)

type Sign struct {
	Script
	signInstanceWid string
	signWid         string
}

// DoSign 签到
func (s Sign) DoSign(client *http.Client) bool {
	// List of Sign
	var list = signListJSON{}
	_ = json.Unmarshal(s.getStuSignInfosInOneDay(client), &list)
	sign := list.Datas.UnSignedTasks
	if sign == nil || len(sign) == 0 {
		utils.Log("SignError").Message("AThere is no sign to do.")
		return false
	}
	s.signInstanceWid = sign[0].SignInstanceWid
	s.signWid = sign[0].SignWid
	// Detail of Sign
	detail := s.detailSignInstance(client)
	if detail == "" {
		return false
	}
	// Fill the Sign form
	form := fillSign(detail)
	// Submit Sign
	res := s.submitSign(client, form)
	if res == "" {
		return false
	}
	if strings.Contains(res, "SUCCESS") {
		utils.Log("SignInfo").Message("Success")

		return true
	} else {
		utils.Log("SignError").Message(res)
	}

	return false
}

// getStuSignInfosInOneDay 获取签到列表
func (s Sign) getStuSignInfosInOneDay(client *http.Client) []byte {
	signURL := fmt.Sprintf("https://%v/wec-counselor-sign-apps/stu/sign/getStuSignInfosInOneDay", s.Host)
	req, _ := http.NewRequest("POST", signURL, bytes.NewBuffer([]byte(`{}`)))
	req.Header.Add("User-Agent",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) "+
			"Mobile/15E148 (4789933056)cpdaily/8.2.4  wisedu/8.2.4")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/json")
	res, _ := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		return nil
	}
	ret, _ := ioutil.ReadAll(res.Body)

	return ret
}

// 获取签到详细信息
func (s Sign) detailSignInstance(client *http.Client) string {
	signURL := fmt.Sprintf("https://%v/wec-counselor-sign-apps/stu/sign/detailSignInstance", s.Host)
	res, _ := client.Post(signURL, "application/json", strings.NewReader(
		fmt.Sprintf(`{"signInstanceWid":"%v","signWid":"%v"}`, s.signInstanceWid, s.signWid)))
	if res != nil {
		defer res.Body.Close()
	} else {
		utils.Log("SignError").Message("Error getting detail of the SignInstance.")
		return ""
	}
	ret, _ := ioutil.ReadAll(res.Body)
	if !strings.Contains(string(ret), "学生晨午晚检") {
		utils.Log("SignError").Message("Not sign task.")
		return ""
	}

	return string(ret)
}

// 提交签到信息
func (s Sign) submitSign(client *http.Client, extraFieldItems []signFillForm) string {
	submitURL := fmt.Sprintf("https://%v/wec-counselor-sign-apps/stu/sign/submitSign", s.Host)
	body := signData{
		Position:        s.SignAddr,
		AbnormalReason:  "",
		IsMalposition:   0,
		IsNeedExtra:     1,
		Latitude:        "32.55562214980632",
		Longitude:       "117.0298292824606",
		SignInstanceWid: s.signInstanceWid,
		SignPhotoURL:    "",
		ExtraFieldItems: extraFieldItems,
		UaIsCpadaily:    true,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", submitURL, strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cpdaily-Extension", utils.GetExtension(s.Key))
	res, _ := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		utils.Log("SignError").Message("Error submitting the Sign.")

		return ""
	}
	finalStr, _ := ioutil.ReadAll(res.Body)

	return string(finalStr)
}

// 填写签到信息
func fillSign(details string) []signFillForm {
	var form = signDetailJSON{}
	var ret []signFillForm
	if details == "" {
		return []signFillForm{}
	}
	_ = json.Unmarshal([]byte(details), &form)
	for _, v := range form.Datas.ExtraField {
		for _, t := range v.ExtraFieldItems {
			if t.Content == "腋下温度37.3℃以下" || t.Content == "无" {
				ret = append(ret, signFillForm{
					ExtraFieldItemValue: t.Content,
					ExtraFieldItemWid:   t.Wid,
				})
			}
		}
	}
	return ret
}
