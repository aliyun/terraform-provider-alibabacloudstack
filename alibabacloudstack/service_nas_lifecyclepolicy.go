package alibabacloudstack

import (
	"encoding/json"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

// type
// NasService struct {
// 	client *connectivity.AlibabacloudStackClient
// }

type NasDescribelifecyclepoliciesResponse struct {
	LifecyclePolicies struct {
		LifecyclePolicy []struct {
			FileSystemId        string `json:"FileSystemId"`
			LifecyclePolicyName string `json:"LifecyclePolicyName"`
			Path                string `json:"Path"`
			Recursive           bool   `json:"Recursive"`
			LifecycleRuleName   string `json:"LifecycleRuleName"`
			StorageType         string `json:"StorageType"`
			CreateTime          string `json:"CreateTime"`
			OssBucket           string `json:"OssBucket"`
		} `json:"LifecyclePolicy"`
	} `json:"LifecyclePolicies"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageSize   int    `json:"PageSize"`
	PageNumber int    `json:"PageNumber"`
}

func (s *NasService) DoNasDescribelifecyclepoliciesRequest(id string) (*NasDescribelifecyclepoliciesResponse, error) {
	// api: NAS - 2017-06-26 - DescribeLifecyclePolicies
	request := s.client.NewCommonRequest("GET", "NAS", "2017-06-26", "DescribeLifecyclePolicies", "")
	NasDescribelifecyclepoliciesResponseObj := &NasDescribelifecyclepoliciesResponse{}
	params := strings.Split(id, ":")
	request.QueryParams["FileSystemId"] = params[0]
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeLifecyclePolicies", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &NasDescribelifecyclepoliciesResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeLifecyclePolicies", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return NasDescribelifecyclepoliciesResponseObj, nil
}
