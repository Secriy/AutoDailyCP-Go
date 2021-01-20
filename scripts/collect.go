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

// 收集列表JSON
type collectInfoJson struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		Wid string `json:"wid"`
	} `json:"datas"`
}

// 收集详细信息JSON
type collectDetailJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		Collector struct {
			FormWid       string `json:"formWid"`
			Wid           string `json:"wid"`
			SchoolTaskWid string `json:"schoolTaskWid"`
		} `json:"collector"`
	} `json:"datas"`
}

// 提交收集选项JSON
type fillCollectForm struct {
	ExtraFieldItemValue string `json:"extraFieldItemValue"`
	ExtraFieldItemWid   int    `json:"extraFieldItemWid"`
}

// 提交收集信息
type collectData struct {
	FormWid       string            `json:"formWid"`
	CollectWid    string            `json:"collectWid"`
	SchoolTaskWid string            `json:"schoolTaskWid"`
	Form          []fillCollectForm `json:"form"`
	Address       string            `json:"address"`
	UaIsCpadaily  bool              `json:"uaIsCpadaily"`
}

func DoCollect(host string, key string, client *http.Client, collectAddr string) string {
	var info = collectInfoJson{}
	var details = collectDetailJSON{}
	var form = collectData{}.Form
	collector := details.Datas.Collector
	_ = json.Unmarshal(queryCollectorProcessingList(host, client), &info)
	if info.Datas.Wid == "" {
		return utils.Message("Collect Error: " + "There is no collect to do.")
	}
	_ = json.Unmarshal(detailCollector(host, client, info.Datas.Wid), &details)
	_ = json.Unmarshal(getFormFields(host, client, collector.FormWid, collector.Wid), &form)

	return utils.Message(submitForm(host, client, key, collector.FormWid, collector.Wid, collector.SchoolTaskWid, form, collectAddr))
}

// queryCollectorProcessingList 获取收集列表
func queryCollectorProcessingList(host string, client *http.Client) []byte {
	collectURL := fmt.Sprintf("https://%v/wec-counselor-collector-apps/stu/collector/queryCollectorProcessingList", host)
	req, _ := http.NewRequest("POST", collectURL, bytes.NewBuffer([]byte(`{"pageSize": 10, "pageNumber": 1}`)))
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

func detailCollector(host string, client *http.Client, collectorWid string) []byte {
	collectURL := fmt.Sprintf("https://%v/wec-counselor-collector-apps/stu/collector/detailCollector", host)
	res, _ := client.Post(collectURL, "application/json", strings.NewReader(fmt.Sprintf(`{"collectorWid":"%v"}`, collectorWid)))
	if res != nil {
		defer res.Body.Close()
	} else {
		return nil
	}
	ret, _ := ioutil.ReadAll(res.Body)
	return ret
}

func getFormFields(host string, client *http.Client, formWid string, collectorWid string) []byte {
	collectURL := fmt.Sprintf("https://%v/wec-counselor-collector-apps/stu/collector/getFormFields", host)
	res, _ := client.Post(collectURL, "application/json", strings.NewReader(fmt.Sprintf(`{"pageSize": 50,
            "pageNumber": 1,
            "formWid": %v,
            "collectorWid": %v}`, formWid, collectorWid)))
	if res != nil {
		defer res.Body.Close()
	} else {
		return nil
	}
	ret, _ := ioutil.ReadAll(res.Body)
	return ret
}

func submitForm(host string, client *http.Client, key string, formWid string, collectWid string, schoolTaskWid string, form []fillCollectForm, address string) string {
	collectURL := fmt.Sprintf("https://%v/wec-counselor-sign-apps/stu/collector/submitForm", host)
	body := collectData{
		FormWid:       formWid,
		CollectWid:    collectWid,
		SchoolTaskWid: schoolTaskWid,
		Form:          form,
		Address:       address,
		UaIsCpadaily:  true,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", collectURL, strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cpdaily-Extension", utils.GetExtension(key))
	res, _ := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		return utils.Message("Collect Error: " + "Error submiting the Collector.")
	}
	finlStr, _ := ioutil.ReadAll(res.Body)
	return string(finlStr)
}
