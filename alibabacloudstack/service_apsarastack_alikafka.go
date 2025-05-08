package alibabacloudstack

import (
	"encoding/json"
	"log"
	"regexp"
	"time"
	"fmt"
	"reflect"
	
	"github.com/PaesslerAG/jsonpath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type TopicListResponse struct {
	EagleEyeTraceId string          `json:"eagleEyeTraceId"`
	AsapiSuccess    bool            `json:"asapiSuccess"`
	RequestId       string          `json:"RequestId"`
	Message         string          `json:"Message"`
	PageSize        int             `json:"PageSize"`
	Code            int             `json:"Code"`
	Success         bool            `json:"Success"`
	ResponseVersion string          `json:"responseVersion"`
	CurrentPage     int             `json:"CurrentPage"`
	Total           int             `json:"Total"`
	TopicList       []AliKafkaTopic `json:"TopicList"`
}

type AliKafkaTopic struct {
	InstanceDo   interface{} `json:"instanceDo"`
	RoleList     []string    `json:"roleList"`
	Tags         []string    `json:"tags"`
	LocalTopic   bool        `json:"localTopic"`
	InstanceId   string      `json:"instanceId"`
	RelationName string      `json:"relationName"`
	HaveAlarm    bool        `json:"haveAlarm"`
	StatusName   string      `json:"statusName"`
	AlarmList    []string    `json:"alarmList"`
	Topic        string      `json:"topic"`
	ChannelName  string      `json:"channelName"`
	AuthType     int         `json:"authType"`
	Status       int         `json:"status"`
}

type AlikafkaService struct {
	client *connectivity.AlibabacloudStackClient
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaInstance(instanceId string) (*alikafka.InstanceVO, error) {
	alikafkaInstance := &alikafka.InstanceVO{}
	instanceListReq := alikafka.CreateGetInstanceListRequest()
	alikafkaService.client.InitRpcRequest(*instanceListReq.RpcRequest)
	instanceListReq.QueryParams["Product"] = "alikafka"
	wait := incrementalWait(2*time.Second, 1*time.Second)
	var raw interface{}
	var err error
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(client *alikafka.Client) (interface{}, error) {
			return client.GetInstanceList(instanceListReq)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)
		return nil
	})

	instanceListResp, ok := raw.(*alikafka.GetInstanceListResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(instanceListResp.BaseResponse)
		}
		return alikafkaInstance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceId, instanceListReq.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)

	for _, v := range instanceListResp.InstanceList.InstanceVO {
		if v.InstanceId == instanceId && v.ServiceStatus != 10 {
			return &v, nil
		}
	}
	return alikafkaInstance, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AlikafkaInstance", instanceId)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaInstanceByOrderId(orderId string, timeout int) (*alikafka.InstanceVO, error) {
	alikafkaInstance := &alikafka.InstanceVO{}
	instanceListReq := alikafka.CreateGetInstanceListRequest()
	alikafkaService.client.InitRpcRequest(*instanceListReq.RpcRequest)
	instanceListReq.OrderId = orderId
	instanceListReq.QueryParams["Product"] = "alikafka"
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		wait := incrementalWait(2*time.Second, 1*time.Second)
		var raw interface{}
		var err error
		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			raw, err = alikafkaService.client.WithAlikafkaClient(func(client *alikafka.Client) (interface{}, error) {
				return client.GetInstanceList(instanceListReq)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)
			return nil
		})

		instanceListResp, ok := raw.(*alikafka.GetInstanceListResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(instanceListResp.BaseResponse)
			}
			return alikafkaInstance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, orderId, instanceListReq.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)

		for _, v := range instanceListResp.InstanceList.InstanceVO {
			return &v, nil
		}
		if time.Now().After(deadline) {
			return alikafkaInstance, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AlikafkaInstance", orderId)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaConsumerGroup(id string) (*alikafka.ConsumerVO, error) {
	alikafkaConsumerGroup := &alikafka.ConsumerVO{}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return alikafkaConsumerGroup, errmsgs.WrapError(err)
	}
	instanceId := parts[0]

	request := alikafka.CreateGetConsumerListRequest()
	alikafkaService.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.QueryParams["Product"] = "alikafka"
	wait := incrementalWait(2*time.Second, 1*time.Second)
	var raw interface{}
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(client *alikafka.Client) (interface{}, error) {
			return client.GetConsumerList(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	consumerListResp, ok := raw.(*alikafka.GetConsumerListResponse)
	if err != nil {
		errmsg := ""

		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(consumerListResp.BaseResponse)
		}
		return alikafkaConsumerGroup, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	//for _, v := range consumerListResp.ConsumerList {
	//	if v.ConsumerId == consumerId {
	//		return &v, nil
	//	}
	//}
	return alikafkaConsumerGroup, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AlikafkaConsumerGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaTopicStatus(id string) (*alikafka.TopicStatus, error) {
	alikafkaTopicStatus := &alikafka.TopicStatus{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return alikafkaTopicStatus, errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := alikafka.CreateGetTopicStatusRequest()
	alikafkaService.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.Topic = topic
	request.QueryParams["Product"] = "alikafka"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.GetTopicStatus(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	topicStatusResp, ok := raw.(*alikafka.GetTopicStatusResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(topicStatusResp.BaseResponse)
		}
		return alikafkaTopicStatus, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if topicStatusResp.TopicStatus.OffsetTable.OffsetTableItem != nil {
		return &topicStatusResp.TopicStatus, nil
	}

	return alikafkaTopicStatus, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AlikafkaTopicStatus "+errmsgs.ResourceNotfound, id)), errmsgs.ResourceNotfound)
}
func (alikafkaService *AlikafkaService) DoAlikafkaGettopiclistRequest(id string) (*AliKafkaTopic, error) {
	return alikafkaService.DescribeAlikafkaTopic(id)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaTopic(id string) (*AliKafkaTopic, error) {
	alikafkaTopic := &AliKafkaTopic{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return alikafkaTopic, errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	// request := alikafka.CreateGetTopicListRequest()
	request := alikafkaService.client.NewCommonRequest("POST", "alikafka", "2019-09-16", "GetTopicList", "")
	request.QueryParams["InstanceId"] = instanceId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.ProcessCommonRequest(request)
		})
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	topicListResp := TopicListResponse{}
	bresponse, ok := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &topicListResp)
	if err != nil && !ok {
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return alikafkaTopic, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	for _, v := range topicListResp.TopicList {
		if v.Topic == topic {
			return &v, nil
		}
	}
	return alikafkaTopic, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AlikafkaTopic", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaSaslUser(id string) (*alikafka.SaslUserVO, error) {
	alikafkaSaslUser := &alikafka.SaslUserVO{}

	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return alikafkaSaslUser, errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]

	request := alikafka.CreateDescribeSaslUsersRequest()
	alikafkaService.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DescribeSaslUsers(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	userListResp, ok := raw.(*alikafka.DescribeSaslUsersResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(userListResp.BaseResponse)
		}
		return alikafkaSaslUser, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range userListResp.SaslUserList.SaslUserVO {
		if v.Username == username {
			return &v, nil
		}
	}
	return alikafkaSaslUser, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AlikafkaSaslUser", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaSaslAcl(id string) (*alikafka.KafkaAclVO, error) {
	alikafkaSaslAcl := &alikafka.KafkaAclVO{}

	parts, err := ParseResourceId(id, 6)
	if err != nil {
		return alikafkaSaslAcl, errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]
	aclResourceType := parts[2]
	aclResourceName := parts[3]
	aclResourcePatternType := parts[4]
	aclOperationType := parts[5]

	request := alikafka.CreateDescribeAclsRequest()
	alikafkaService.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.Username = username
	request.AclResourceType = aclResourceType
	request.AclResourceName = aclResourceName
	request.AclResourcePatternType = aclResourcePatternType
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DescribeAcls(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	aclListResp, ok := raw.(*alikafka.DescribeAclsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(aclListResp.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"BIZ_SUBSCRIPTION_NOT_FOUND", "BIZ_TOPIC_NOT_FOUND"}) {
			return alikafkaSaslAcl, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return alikafkaSaslAcl, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range aclListResp.KafkaAclList.KafkaAclVO {
		if v.AclResourcePatternType == aclResourcePatternType && v.AclOperationType == aclOperationType {
			return &v, nil
		}
	}
	return alikafkaSaslAcl, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AlikafkaSaslAcl", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *AlikafkaService) WaitForAlikafkaInstanceUpdated(id string, topicQuota int, diskSize int, ioMax int, eipMax int, paidType int, specType string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaInstance(id)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		if object.InstanceId == id && object.TopicNumLimit == topicQuota && object.DiskSize == diskSize && object.IoMax == ioMax && object.EipMax == eipMax && object.PaidType == paidType && object.SpecType == specType {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceId, id, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if object.InstanceId == id && status == Running {
			if object.ServiceStatus == 5 {
				return nil
			}
		} else if object.InstanceId == id {
			if status != Deleted {
				return nil
			} else if object.ServiceStatus == 10 {
				return nil
			}
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceId, id, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaConsumerGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaConsumerGroup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if object.InstanceId+":"+object.ConsumerId == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceId+":"+object.ConsumerId, id, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) KafkaTopicListRefreshFunc(id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlikafkaTopic(id)
		if err != nil {
			if !errmsgs.IsExpectedErrors(err, []string{errmsgs.ResourceNotfound}) {
				return nil, "", errmsgs.WrapError(err)
			}
		}

		return object, "Creating", nil
	}
}

func (s *AlikafkaService) KafkaTopicStatusRefreshFunc(id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlikafkaTopicStatus(id)
		if err != nil {
			if !errmsgs.IsExpectedErrors(err, []string{errmsgs.ResourceNotfound}) {
				return nil, "", errmsgs.WrapError(err)
			}
		}

		if object.OffsetTable.OffsetTableItem != nil && len(object.OffsetTable.OffsetTableItem) > 0 {
			return object, "Running", errmsgs.WrapError(err)
		}

		return object, "Creating", nil
	}
}

func (s *AlikafkaService) WaitForAlikafkaTopic(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaTopic(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if object.InstanceId+":"+object.Topic == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceId+":"+object.Topic, id, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaSaslUser(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	for {
		object, err := s.DescribeAlikafkaSaslUser(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if instanceId+":"+object.Username == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, instanceId+":"+object.Username, id, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaSaslAcl(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 6)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	for {
		object, err := s.DescribeAlikafkaSaslAcl(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if instanceId+":"+object.Username+":"+object.AclResourceType+":"+object.AclResourceName+":"+object.AclResourcePatternType+":"+object.AclOperationType == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, instanceId+":"+object.Username, id, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []alikafka.TagResource, err error) {
	request := alikafka.CreateListTagResourcesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []alikafka.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, alikafka.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.ListTagResources(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling, errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	response, ok := raw.(*alikafka.ListTagResourcesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, resourceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}

	return response.TagResources.TagResource, nil
}

func (s *AlikafkaService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := alikafka.CreateUntagResourcesRequest()
			s.client.InitRpcRequest(*request.RpcRequest)
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = string(resourceType)
			request.TagKey = &tagKey

			wait := incrementalWait(2*time.Second, 1*time.Second)
			var raw interface{}
			var err error
			err = resource.Retry(10*time.Minute, func() *resource.RetryError {
				raw, err = s.client.WithAlikafkaClient(func(client *alikafka.Client) (interface{}, error) {
					return client.UntagResources(request)
				})
				if err != nil {
					if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			response, ok := raw.(*alikafka.UntagResourcesResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
		}

		if len(create) > 0 {
			request := alikafka.CreateTagResourcesRequest()
			s.client.InitRpcRequest(*request.RpcRequest)
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = string(resourceType)

			wait := incrementalWait(2*time.Second, 1*time.Second)
			var raw interface{}
			var err error
			err = resource.Retry(10*time.Minute, func() *resource.RetryError {
				raw, err = s.client.WithAlikafkaClient(func(client *alikafka.Client) (interface{}, error) {
					return client.TagResources(request)
				})
				if err != nil {
					if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			response, ok := raw.(*alikafka.TagResourcesResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
		}
	}

	return nil
}

func (s *AlikafkaService) tagsToMap(tags []alikafka.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *AlikafkaService) ignoreTag(t alikafka.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *AlikafkaService) diffTags(oldTags, newTags []alikafka.TagResourcesTag) ([]alikafka.TagResourcesTag, []alikafka.TagResourcesTag) {
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	var remove []alikafka.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *AlikafkaService) tagsFromMap(m map[string]interface{}) []alikafka.TagResourcesTag {
	result := make([]alikafka.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, alikafka.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *AlikafkaService) AliKafkaInstanceStateRefreshFunc(id, attribute string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlikafkaInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		
		rv := reflect.ValueOf(object)
		// 查找字段
		field := rv.FieldByName(attribute)
		if !field.IsValid() || !field.CanInterface(){
			return nil, "", nil
		}
	
		state := field.Interface().(string)
		

		for _, failState := range failStates {

			if fmt.Sprint(state) == failState {
				return object, fmt.Sprint(state), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(state)))
			}
		}
		return object, fmt.Sprint(state), nil
	}
}

func (s *AlikafkaService) GetQuotaTip(instanceId string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetQuotaTip"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": instanceId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, request)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, errmsgs.WrapError(err)
	}
	v, err := jsonpath.Get("$.QuotaData", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, instanceId, "$.QuotaData", response)
	}
	return v.(map[string]interface{}), nil
}
