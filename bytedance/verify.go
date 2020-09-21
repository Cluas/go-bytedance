package bytedance

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

// Verify 校验加密信息
func Verify(tpToken string, timestamp string, nonce string, encrypt string, msgSignature string) bool {
	values := []string{tpToken, timestamp, nonce, encrypt}
	sort.Strings(values)
	newMsgSignature := Sha1(strings.Join(values, ""))
	return newMsgSignature == msgSignature
}

// Sha1 sha1 加密
func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	encodeStr := fmt.Sprintf("%x", h.Sum(nil))
	return encodeStr
}
