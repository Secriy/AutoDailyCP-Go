package utils

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/forgoer/openssl"
	"github.com/google/uuid"
)

type extension struct {
	DeviceID      string  `json:"deviceId"`
	SystemName    string  `json:"systemName"`
	UserID        string  `json:"userId"`
	AppVersion    string  `json:"appVersion"`
	Model         string  `json:"model"`
	Lon           float64 `json:"lon"`
	SystemVersion string  `json:"systemVersion"`
	Lat           float64 `json:"lat"`
}

// 加密
func encrypt(plaintext []byte, key string) []byte {
	hexStr := "0102030405060708"
	data, _ := hex.DecodeString(hexStr)
	text, _ := openssl.DesCBCEncrypt(plaintext, []byte(key), data, openssl.PKCS5_PADDING)
	res := base64.StdEncoding.EncodeToString(text)
	return []byte(res)
}

// GetExtension 获取Extension
func GetExtension(key string) string {
	ext := extension{
		DeviceID:      uuid.New().String(),
		SystemName:    "StudioX",
		UserID:        "23333333",
		AppVersion:    "8.2.22",
		Model:         "X233",
		Lon:           117.0298292824606,
		SystemVersion: "19.1",
		Lat:           32.55562214980632,
	}
	foreEnc, _ := json.Marshal(ext)
	return string(encrypt(foreEnc, key))
}
