package encoding

import (
	"encoding/base64"
	"github.com/Cospk/base-tools/errs"
)

// Base64Encode 将输入字符串编码为 Base64 字符串。
func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64Decode 解码 Base64 编码的字符串并返回原始字符串。
// 若输入不是合法的 Base64，将返回错误。
func Base64Decode(data string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errs.WrapMsg(err, "DecodeString failed", "data", data)
	}
	return string(decodedBytes), nil
}
