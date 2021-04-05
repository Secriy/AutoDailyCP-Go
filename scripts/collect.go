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

type Collect struct {
	Script
	collectWid    string
	formWid       string
	schoolTaskWid string
}

// DoCollect 执行收集
func (s Collect) DoCollect(client *http.Client) bool {
	var list = collectListJSON{}
	var row rowData
	// List of Collect
	_ = json.Unmarshal(s.queryCollectorProcessingList(client), &list)
	// 判断收集数量
	switch {
	case len(list.Datas.Rows) == 0:
		utils.Log("CollectError").Message("There is no collect to do.")
		return false
	case len(list.Datas.Rows) > 0:
		for k, v := range list.Datas.Rows {
			if strings.Contains(v.Subject, "每日学生健康打卡") {
				row = list.Datas.Rows[k]
				s.collectWid = row.Wid
				s.formWid = row.FormWid
			}
		}
		if row.Subject == "" {
			utils.Log("CollectError").Message("There is no collect to do.")
			return false
		}
	}
	// 查询收集详细信息
	var detail = collectDetailJSON{}
	_ = json.Unmarshal(s.detailCollector(client), &detail)
	s.schoolTaskWid = detail.Datas.Collector.SchoolTaskWid
	// 获取历史表单
	var form = formJSON{}
	_ = json.Unmarshal(s.getFormFields(client), &form)
	retForm := fillFormFields(form.Datas.Rows)
	// 获取历史地址
	if s.CollectAddr == "" {
		s.CollectAddr = strings.ReplaceAll(form.Datas.Rows[0].Value, "/", "")
	}

	ret := s.submitForm(client, retForm)
	if ret == "" {
		return false
	}

	if strings.Contains(ret, "SUCCESS") {
		utils.Log("").Message("Collect Success")

		return true
	} else {
		utils.Log("").Message(ret)
	}

	return false
}

// queryCollectorProcessingList 获取收集列表
func (s Collect) queryCollectorProcessingList(client *http.Client) []byte {
	collectURL := fmt.Sprintf("https://%v/wec-counselor-collector-apps/stu/collector/queryCollectorProcessingList", s.Host)
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
func (s Collect) detailCollector(client *http.Client) []byte {
	collectURL := fmt.Sprintf("https://%v/wec-counselor-collector-apps/stu/collector/detailCollector", s.Host)
	res, _ := client.Post(collectURL, "application/json", strings.NewReader(fmt.Sprintf(`{"collectorWid":"%v"}`, s.collectWid)))
	if res != nil {
		defer res.Body.Close()
	} else {
		return nil
	}
	ret, _ := ioutil.ReadAll(res.Body)

	return ret
}

// getFormFields 获取历史表单
func (s Collect) getFormFields(client *http.Client) []byte {
	collectURL := fmt.Sprintf("https://%v/wec-counselor-collector-apps/stu/collector/getFormFields", s.Host)
	res, _ := client.Post(collectURL, "application/json", strings.NewReader(fmt.Sprintf(`{"pageSize": 50,"pageNumber": 1,"formWid": %v,"collectorWid": %v}`, s.formWid, s.collectWid)))
	if res != nil {
		defer res.Body.Close()
	} else {
		return nil
	}
	ret, _ := ioutil.ReadAll(res.Body)

	return ret
}

// submitForm 提交表单
func (s Collect) submitForm(client *http.Client, form []collectRowJSON) string {
	collectURL := fmt.Sprintf("https://%v/wec-counselor-collector-apps/stu/collector/submitForm", s.Host)
	body := collectData{
		FormWid:       s.formWid,
		CollectWid:    s.collectWid,
		SchoolTaskWid: s.schoolTaskWid,
		Form:          form,
		Address:       s.CollectAddr,
		UaIsCpadaily:  true,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", collectURL, strings.NewReader(string(jsonBody)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cpdaily-Extension", utils.GetExtension(s.Key))
	res, _ := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	} else {
		utils.Log("CollectError").Message("Error submitting the Collector.")
		return ""
	}
	finalStr, _ := ioutil.ReadAll(res.Body)

	return string(finalStr)
}

// fillFormFields 填充表单
func fillFormFields(rs []collectRowJSON) []collectRowJSON {
	retForm := make([]collectRowJSON, 0, len(rs))
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
