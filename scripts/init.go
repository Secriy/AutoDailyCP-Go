package scripts

type Script struct {
	Host        string
	Key         string
	SignAddr    string
	CollectAddr string
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
	ExtraFieldItems []signFillForm `json:"extraFieldItems"`
	UaIsCpadaily    bool           `json:"uaIsCpadaily"`
	SignVersion     string         `json:"signVersion"`
}

// 签到列表JSON
type signListJSON struct {
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
type signFillForm struct {
	ExtraFieldItemValue string `json:"extraFieldItemValue"`
	ExtraFieldItemWid   int    `json:"extraFieldItemWid"`
}

// 收集数据Data
type rowData struct {
	Wid     string `json:"wid"`
	FormWid string `json:"formWid"`
	Subject string `json:"subject"`
}

// 收集列表JSON
type collectListJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		Rows []rowData `json:"rows"`
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
		Rows []collectRowJSON `json:"rows"`
	} `json:"datas"`
}

// 信息表单collectRowJSON
type collectRowJSON struct {
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
	FormWid       string           `json:"formWid"`
	CollectWid    string           `json:"collectWid"`
	SchoolTaskWid string           `json:"schoolTaskWid"`
	Form          []collectRowJSON `json:"form"`
	Address       string           `json:"address"`
	UaIsCpadaily  bool             `json:"uaIsCpadaily"`
}
