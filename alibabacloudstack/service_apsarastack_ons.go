package alibabacloudstack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"strings"
)

type OnsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *OnsService) DescribeOnsInstance(instanceid string) (response *OnsInstance, err error) {
	var requestInfo *ecs.Client

	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "Ons-inner",
		"Action":          "ConsoleInstanceBaseInfo",
		"Version":         "2018-02-05",
		"OnsRegionId":     s.client.RegionId,
		"PreventCache":    "",
		"InstanceId":      instanceid,
	}
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.ServiceCode = "Ons-inner"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ConsoleInstanceBaseInfo"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &OnsInstance{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorInstanceNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, instanceid, "ConsoleInstanceBaseInfo", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ConsoleInstanceBaseInfo", response, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if bresponse != nil || resp.Success != true {
		return resp, WrapError(err)
	}

	return resp, nil
}

type TopicStruct struct {
	Data      string `json:"Data"`
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
	Success   bool   `json:"Success"`
	Code      int    `json:"Code"`
}

func (s *OnsService) DescribeOnsTopic(id string) (response *Topic, err error) {
	var requestInfo *ecs.Client
	did, err := ParseResourceId(id, 2)
	if err != nil {
		return response, WrapError(err)
	}
	TopicId := did[0]
	InstanceId := did[1]
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "Ons-inner",
		"Action":          "ConsoleTopicList",
		"Version":         "2018-02-05",
		"topic":           TopicId,
		"OnsRegionId":     s.client.RegionId,
		"PreventCache":    "",
		"InstanceId":      InstanceId,
	}
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.ServiceCode = "Ons-inner"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ConsoleTopicList"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &Topic{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorTopicNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, did[0], "ConsoleTopicList", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ConsoleTopicList", response, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == 200 {
		return resp, WrapError(err)
	}

	return resp, nil
}
func (s *OnsService) DescribeOnsGroup(id string) (response *OnsGroup, err error) {
	var requestInfo *ecs.Client
	did, err := ParseResourceId(id, 2)
	if err != nil {
		return response, WrapError(err)
	}
	GroupId := did[0]
	InstanceId := did[1]
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "Ons-inner",
		"Action":          "ConsoleGroupList",
		"Version":         "2018-02-05",
		"GroupId":         GroupId,
		"OnsRegionId":     s.client.RegionId,
		"PreventCache":    "",
		"InstanceId":      InstanceId,
	}
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.ServiceCode = "Ons-inner"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ConsoleGroupList"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &OnsGroup{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, did[0], "ConsoleGroupList", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ConsoleGroupList", response, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == 200 {
		return resp, WrapError(err)
	}

	return resp, nil
}
