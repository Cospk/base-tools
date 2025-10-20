// Package env 提供环境变量读取工具，支持字符串、整数、浮点数和布尔类型
package env

import (
	"github.com/Cospk/base-tools/errs"
	"os"
	"strconv"
)

// GetString 获取字符串类型的环境变量，如果不存在则返回默认值
func GetString(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	}
	return defaultValue
}

// GetInt 获取整数类型的环境变量，如果不存在则返回默认值
func GetInt(key string, defaultValue int) (int, error) {
	v, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue, errs.WrapMsg(err, "Atoi failed", "value", v)
		}
		return value, nil
	}
	return defaultValue, nil
}

// GetFloat64 获取浮点数类型的环境变量，如果不存在则返回默认值
func GetFloat64(key string, defaultValue float64) (float64, error) {
	v, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return defaultValue, errs.WrapMsg(err, "ParseFloat failed", "value", v)
		}
		return value, nil
	}
	return defaultValue, nil
}

// GetBool 获取布尔类型的环境变量，如果不存在则返回默认值
func GetBool(key string, defaultValue bool) (bool, error) {
	v, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseBool(v)
		if err != nil {
			return defaultValue, errs.WrapMsg(err, "ParseBool failed", "value", v)
		}
		return value, nil
	}
	return defaultValue, nil
}
