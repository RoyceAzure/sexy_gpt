package middleware

import (
	"strings"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
)

func CustomOutgoingHeaderMatcher(key string) (string, bool) {
	// 仅转发以"X-Custom-"开头的元数据
	if strings.HasPrefix(key, "X-Custom-") || key == util.DBMSGKey {
		return key, true
	}
	// 不转发其他元数据
	return "", false
}
