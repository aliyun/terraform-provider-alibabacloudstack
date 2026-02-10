package alibabacloudstack

import (
	"fmt"
	"strings"

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
	page := 1
	request := map[string]interface{}{
		"AclUuid":   parts[0], // XXX: unsupoort 
		"Direction": parts[1],
		"PageSize":  100,
	}
	for {
		request["CurrentPage"] = page
		response, err = s.client.DoTeaRequest("POST", "Cloudfw", "2017-12-07", action, "", nil, nil, request)
		if err != nil {
			return object, err
		}
		v, err := jsonpath.Get("$.Policys", response)
		if err != nil {
			return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Policys", response)
		}
		if len(v.([]interface{})) < 1 {
			break
		}

		for _, vv := range v.([]interface{}) {
			object := vv.(map[string]interface{})
			if object["AclUuid"].(string) != parts[0] {
				continue
			}
			if object["Direction"].(string) != parts[1] {
				continue
			}
			return object, nil
		}
		page += 1
	}

	return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CloudFirewall", id)), errmsgs.NotFoundWithResponse, response)
}

func (s *CloudfwService) DescribeAddressBook(id string) (map[string]interface{}, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, err
	}
	reqQuery := map[string]interface{}{
		"SourceCode":  "yundun",
		"CurrentPage": 1,
		"PageSize":    100,
		"GroupType":   parts[0],
	}

	response, err := s.client.DoTeaRequest("GET", "Cloudfw", "2017-12-07", "DescribeAddressBook", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}
	if response["Acls"] != nil {
		for _, v := range response["Acls"].([]interface{}) {
			acl := v.(map[string]interface{})
			if acl["GroupUuid"] != nil && acl["GroupUuid"].(string) == parts[1] {
				return acl, nil
			}
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Address book with GroupUuid %s not found", id))
}

func (s *CloudfwService) DescribeCloudfwVpcControlPolicy(id string) (map[string]interface{}, error) {
	parst := strings.Split(id, ":")
	reqQuery := map[string]interface{}{
		"SourceCode":    "yundun",
		"CurrentPage":   1,
		"PageSize":      100,
		"AclUuid":       parst[0],
		"VpcFirewallId": "",
	}
	if parst[1] != "" {
		reqQuery["Direction"] = parst[1]
	}

	response, err := s.client.DoTeaRequest("GET", "Cloudfw", "2017-12-07", "DescribeVpcFirewallControlPolicy", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	policys, ok := response["Policys"]
	if !ok {
		return nil, errmsgs.GetNotFoundErrorFromString("Cloudfw control policy not found")
	}

	policyList, ok := policys.([]interface{})
	if !ok || len(policyList) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString("Cloudfw control policy not found")
	}

	for _, policy := range policyList {
		policyMap := policy.(map[string]interface{})
		if policyMap["AclUuid"] == parst[0] {
			return policyMap, nil
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString("Cloudfw control policy not found")
}
