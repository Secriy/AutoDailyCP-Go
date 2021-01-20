package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// GetCurrentDirectory 获取程序执行目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Log().Error("获取程序执行目录出错：" + err.Error())
	}

	return strings.Replace(dir, "\\", "/", -1)
}
