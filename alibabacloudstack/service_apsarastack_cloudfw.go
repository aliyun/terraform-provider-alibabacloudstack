package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type CloudfwService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *CloudfwService) DoCloudfwDescribecontrolpolicyRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeCloudFirewallControlPolicy(id)
}

func (s *CloudfwService) DescribeCloudFirewallControlPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeControlPolicy"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AclUuid":   parts[0],
		"Direction": parts[1],
		"CurrentPage": 1,
		"PageSize":  100,
	}

	response, err = s.client.DoTeaRequest("POST", "Cloudfw", "2017-12-07", action, "", nil, nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.Policys", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Policys", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CloudFirewall", id)), errmsgs.NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}
