package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type OnsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *OnsService) DescribeOnsInstance(instanceid string) (response *OnsInstance, err error) {
	var requestInfo *ons.Client

	request := s.client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleInstanceBaseInfo", "")
	request.QueryParams["OnsRegionId"] = s.client.RegionId
	request.QueryParams["PreventCache"] = ""
	request.QueryParams["InstanceId"] = instanceid

	var resp = &OnsInstance{}
	raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorInstanceNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceid, "ConsoleInstanceBaseInfo", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ConsoleInstanceBaseInfo", response, requestInfo, request)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if bresponse != nil && !resp.Success {
		return resp, errmsgs.WrapError(err)
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
	var requestInfo *ons.Client
	did, err := ParseResourceId(id, 2)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	TopicId := did[0]
	InstanceId := did[1]

	request := s.client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleTopicList", "")
	request.QueryParams["Topic"] = TopicId
	request.QueryParams["OnsRegionId"] = s.client.RegionId
	request.QueryParams["PreventCache"] = ""
	request.QueryParams["InstanceId"] = InstanceId

	var resp = &Topic{}
	raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorTopicNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, did[0], "ConsoleTopicList", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ConsoleTopicList", response, requestInfo, request)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == 200 {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *OnsService) DescribeOnsGroup(id string) (response *OnsGroup, err error) {
	var requestInfo *ons.Client
	did, err := ParseResourceId(id, 2)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	GroupId := did[0]
	InstanceId := did[1]

	request := s.client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleGroupList", "")
	request.QueryParams["GroupId"] = GroupId
	request.QueryParams["OnsRegionId"] = s.client.RegionId
	request.QueryParams["PreventCache"] = ""
	request.QueryParams["InstanceId"] = InstanceId

	var resp = &OnsGroup{}
	raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, did[0], "ConsoleGroupList", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ConsoleGroupList", response, requestInfo, request)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == 200 {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *OnsService) DescribeMqttInstance(id string) (map[string]interface{}, error) {
	reqQuery := map[string]interface{}{
		"Platform":       "onsConsole",
		"OnsRegionId":    s.client.RegionId,
		"MqttInstanceId": id,
		"PreventCache":   time.Now().UnixNano() / 1e6,
	}

	response, err := s.client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttInstanceBaseInfo", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}
	if data, ok := response["Data"]; !ok || data == nil {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("MqttInstance %s not found", id))
	} else {
		return data.(map[string]interface{}), nil
	}
}

func (s *OnsService) DescribeOnsMqttTopic(id string) (map[string]interface{}, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, expected {InstanceId}:{Topic}")
	}
	instanceId := parts[0]
	topic := parts[1]

	reqQuery := map[string]interface{}{
		"InstanceId":  instanceId,
		"CurrentPage": 1,
		"PageSize":    100,
		"isFuzzy":     false,
	}

	response, err := s.client.DoTeaRequest("GET", "Ons-inner", "2018-02-05", "ConsoleTopicListInPage", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}
	if data, ok := response["Data"].([]interface{}); ok {
		for _, item := range data {
			if topicItem, ok := item.(map[string]interface{}); ok {
				if topicItem["topic"].(string) == topic {
					return topicItem, nil
				}
			}
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Mqtt Topic %s not found", id))
}
