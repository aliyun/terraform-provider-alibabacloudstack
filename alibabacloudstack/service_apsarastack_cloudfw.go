package alibabacloudstack

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type CloudfwService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *CloudfwService) DescribeCloudFirewallControlPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudfwClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeControlPolicy"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AclUuid":     parts[0],
		"Direction":   parts[1],
		"CurrentPage": 1,
		"PageSize":    100,
	}
	request["RegionId"] = s.client.RegionId
	request["Product"] = "Cloudfw"
	request["OrganizationId"] = s.client.Department
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Policys", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Policys", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudFirewall", id)), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}
