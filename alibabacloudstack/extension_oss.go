package alibabacloudstack

import (
	"errors"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type LifecycleRuleStatus string

const (
	ExpirationStatusEnabled  = LifecycleRuleStatus("Enabled")
	ExpirationStatusDisabled = LifecycleRuleStatus("Disabled")
)

func ossNotFoundError(err error) bool {
	var se *oss.ServiceError
	if errors.As(err, &se) &&
		(se.StatusCode == 404 || strings.HasPrefix(se.Code, "NoSuch") || strings.HasPrefix(se.Message, "No Row found")) {
		return true
	}
	return false
}
