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

// 签到列表JSON
type signInfoJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		UnSignedTasks []struct {
			SignInstanceWid string `json:"signInstanceWid"`
			SignWid         string `json:"signWid"`
		} `json:"unSignedTasks"`
	} `json:"datas"`
}

// 签到详细信息JSON
type signDetailJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		ExtraField []struct {
			ExtraFieldItems []struct {
				Content string `json:"content"`
				Wid     int    `json:"wid"`
			} `json:"extraFieldItems"`
		} `json:"extraField"`
		IsNeedExtra bool `json:"isNeedExtra"`
	} `json:"datas"`
}

// 提交签到选项JSON
type fillSignForm struct {
	ExtraFieldItemValue string `json:"extraFieldItemValue"`
	ExtraFieldItemWid   int    `json:"extraFieldItemWid"`
}

// 提交签到信息JSON
type signData struct {
	Position        string         `json:"position"`
	AbnormalReason  string         `json:"abnormalReason"`
	IsMalposition   int            `json:"isMalposition"`
	IsNeedExtra     int            `json:"isNeedExtra"`
	Latitude        string         `json:"latitude"`
	Longitude       string         `json:"longitude"`
	SignInstanceWid string         `json:"signInstanceWid"`
	SignPhotoURL    string         `json:"signPhotoUrl"`
	ExtraFieldItems []fillSignForm `json:"extraFieldItems"`
	UaIsCpadaily    bool           `json:"uaIsCpadaily"`
}

// DoSign 签到
func DoSign(host string, key string, client *http.Client, signAddr string) string {
	var info = signInfoJSON{}
	_ = json.Unmarshal(getStuSignInfosInOneDay(host, client), &info)
	if info.Datas.UnSignedTasks == nil || len(info.Datas.UnSignedTasks) == 0 {
		return utils.Message("Sign Error: " + "There is no sign to do.")
	}
	details := detailSignInstance(host, client, info.Datas.UnSignedTasks[0].SignInstanceWid, info.Datas.UnSignedTasks[0].SignWid)
	fill := fillSign(details)
	return utils.Message(submitSign(host, client, key, info.Datas.UnSignedTasks[0].SignInstanceWid, signAddr, fill))
}

// getStuSignInfosInOneDay 获取签到列表
func getStuSignInfosInOneDay(host string, client *http.Client) []byte {
	signURL := fmt.Sprintf("https://%v/wec-counselor-sign-apps/stu/sign/getStuSignInfosInOneDay", host)
	req, _ := http.NewRequest("POST", signURL, bytes.NewBuffer([]byte(`{}`)))
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 (4789933056)cpdaily/8.2.4  wisedu/8.2.4")
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
func detailSignInstance(host string, client *http.Client, signInstanceWid string, signWid string) string {
	signURL := fmt.Sprintf("https://%v/wec-counselor-sign-apps/stu/sign/detailSignInstance", host)
	res, _ := client.Post(signURL, "application/json", strings.NewReader(fmt.Sprintf(`{"signInstanceWid":"%v","signWid":"%v"}`, signInstanceWid, signWid)))
	if res != nil {
		defer res.Body.Close()
	} else {
		return utils.Message("Sign Error: " + "Error getting detail of the SignInstance.")
	}
	ret, _ := ioutil.ReadAll(res.Body)
	return string(ret)
}

// 填写签到信息
func fillSign(details string) []fillSignForm {
	var form = signDetailJSON{}
	var ret []fillSignForm
	if details == "" {
		return []fillSignForm{}
	}
	_ = json.Unmarshal([]byte(details), &form)
	for _, v := range form.Datas.ExtraField {
		for _, t := range v.ExtraFieldItems {
			if t.Content == "腋下温度37.3℃以下" || t.Content == "无" {
				ret = append(ret, fillSignForm{
					ExtraFieldItemValue: t.Content,
					ExtraFieldItemWid:   t.Wid,
				})
			}
		}
	}
	return ret
}

// 提交签到信息
func submitSign(host string, client *http.Client, key string, signInstanceWid string, address string, extraFieldItems []fillSignForm) string {
	submitURL := fmt.Sprintf("https://%v/wec-counselor-sign-apps/stu/sign/submitSign", host)
	body := signData{
		Position:        address,
		AbnormalReason:  "",
		IsMalposition:   0,
		IsNeedExtra:     1,
		Latitude:        "32.55562214980632",
		Longitude:       "117.0298292824606",
		SignInstanceWid: signInstanceWid,
		SignPhotoURL:    "",
		ExtraFieldItems: extraFieldItems,
		UaIsCpadaily:    true,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", submitURL, strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cpdaily-Extension", utils.GetExtension(key))
	res, _ := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		return utils.Message("Sign Error: " + "Error submiting the Sign.")
	}
	finlStr, _ := ioutil.ReadAll(res.Body)
	return string(finlStr)
}
