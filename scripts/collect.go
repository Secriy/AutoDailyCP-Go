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
		Rows []struct {
			Wid     string `json:"wid"`
			FormWid string `json:"formWid"`
		} `json:"rows"`
	} `json:"datas"`
}

// 收集详细信息JSON
type collectDetailJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		Collector struct {
			SchoolTaskWid string `json:"schoolTaskWid"`
		} `json:"collector"`
	} `json:"datas"`
}

// 信息表单JSON
type formJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		Rows []rowJSON `json:"rows"`
	} `json:"datas"`
}

// 信息表单RowJSON
type rowJSON struct {
	Wid           string  `json:"wid"`
	FormWid       string  `json:"formWid"`
	FieldType     int     `json:"fieldType"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	MinLength     int     `json:"minLength"`
	MaxLength     int     `json:"maxLength"`
	Sort          string  `json:"sort"`
	IsRequired    int     `json:"isRequired"`
	ImageCount    int     `json:"imageCount"`
	HasOtherItems int     `json:"hasOtherItems"`
	ColName       string  `json:"colName"`
	Value         string  `json:"value"`
	MinValue      float32 `json:"minValue"`
	MaxValue      float32 `json:"maxValue"`
	IsDecimal     bool    `json:"isDecimal"`
	FieldItems    []struct {
		ItemWid       string `json:"itemWid"`
		Content       string `json:"content"`
		IsOtherItems  int    `json:"isOtherItems"`
		ContendExtend int    `json:"contendExtend"`
		IsSelected    int    `json:"isSelected"`
	} `json:"fieldItems"`
}

// 提交收集信息
type collectData struct {
	FormWid       string    `json:"formWid"`
	CollectWid    string    `json:"collectWid"`
	SchoolTaskWid string    `json:"schoolTaskWid"`
	Form          []rowJSON `json:"form"`
	Address       string    `json:"address"`
	UaIsCpadaily  bool      `json:"uaIsCpadaily"`
}

// DoCollect 执行收集
func DoCollect(host string, key string, client *http.Client, collectAddr string) string {
	var info = collectInfoJson{}
	var details = collectDetailJSON{}
	var form = formJSON{}
	_ = json.Unmarshal(queryCollectorProcessingList(host, client), &info)
	if len(info.Datas.Rows) == 0 {
		return "Collect Error: There is no collect to do."
	}
	row := info.Datas.Rows[0]
	_ = json.Unmarshal(detailCollector(host, client, row.Wid), &details)
	collector := details.Datas.Collector
	_ = json.Unmarshal(getFormFields(host, client, row.FormWid, row.Wid), &form)
	retForm := fillFormFields(form.Datas.Rows)

	if collectAddr == "" {
		collectAddr = strings.ReplaceAll(form.Datas.Rows[0].Value, "/", "")
	}

	return submitForm(host, client, key, row.FormWid, row.Wid, collector.SchoolTaskWid, retForm, collectAddr)
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

// detailCollector 查询收集信息
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

// getFormFields 获取历史表单
func getFormFields(host string, client *http.Client, formWid string, collectorWid string) []byte {
	collectURL := fmt.Sprintf("https://%v/wec-counselor-collector-apps/stu/collector/getFormFields", host)
	res, _ := client.Post(collectURL, "application/json", strings.NewReader(fmt.Sprintf(`{"pageSize": 50,"pageNumber": 1,"formWid": %v,"collectorWid": %v}`, formWid, collectorWid)))
	if res != nil {
		defer res.Body.Close()
	} else {
		return nil
	}
	ret, _ := ioutil.ReadAll(res.Body)
	return ret
}

// fillFormFields 填充表单
func fillFormFields(rs []rowJSON) []rowJSON {
	retForm := make([]rowJSON, 0, len(rs))
	for _, rw := range rs {
		for _, fi := range rw.FieldItems {
			if fi.IsSelected == 1 {
				rw.FieldItems = append(rw.FieldItems[0:0:0], fi)
				break
			} else {
				continue
			}
		}
		retForm = append(retForm, rw)
	}

	return retForm
}

// submitForm 提交表单
func submitForm(host string, client *http.Client, key string, formWid string, collectWid string, schoolTaskWid string, form []rowJSON, address string) string {
	collectURL := fmt.Sprintf("https://%v/wec-counselor-collector-apps/stu/collector/submitForm", host)
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
		return "Collect Error: Error submiting the Collector."
	}
	finlStr, _ := ioutil.ReadAll(res.Body)

	return string(finlStr)
}
