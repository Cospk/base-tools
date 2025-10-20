package idutil

import (
	"github.com/Cospk/base-tools/utils/encrypt"
	"github.com/Cospk/base-tools/utils/stringutil"
	"github.com/Cospk/base-tools/utils/timeutil"
	"math/rand"
	"strconv"
	"time"
)

// GetMsgIDByMD5 通过 MD5 算法生成消息ID
func GetMsgIDByMD5(sendID string) string {
	t := stringutil.Int64ToString(timeutil.GetCurrentTimestampByNano())
	return encrypt.Md5(t + sendID + stringutil.Int64ToString(rand.Int63n(timeutil.GetCurrentTimestampByNano())))
}

// OperationIDGenerator 生成操作ID
func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}
