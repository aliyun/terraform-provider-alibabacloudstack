package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

type AlikafkaService struct {
	client *connectivity.AlibabacloudStackClient
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaInstance(instanceId string) (*alikafka.InstanceVO, error) {
	alikafkaInstance := &alikafka.InstanceVO{}
	instanceListReq := alikafka.CreateGetInstanceListRequest()
	instanceListReq.RegionId = alikafkaService.client.RegionId
        instanceListReq.QueryParams["Product"] = "alikafka"
	wait := incrementalWait(2*time.Second, 1*time.Second)
	var raw interface{}
	var err error
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(client *alikafka.Client) (interface{}, error) {
			return client.GetInstanceList(instanceListReq)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)
		return nil
	})

	if err != nil {
		return alikafkaInstance, WrapErrorf(err, DefaultErrorMsg, instanceId, instanceListReq.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	instanceListResp, _ := raw.(*alikafka.GetInstanceListResponse)
	addDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)

	for _, v := range instanceListResp.InstanceList.InstanceVO {

		// ServiceStatus equals 10 means the instance is released, do not return the instance.
		if v.InstanceId == instanceId && v.ServiceStatus != 10 {
			return &v, nil
		}
	}
	return alikafkaInstance, WrapErrorf(Error(GetNotFoundMessage("AlikafkaInstance", instanceId)), NotFoundMsg, ProviderERROR)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaNodeStatus(instanceId string) (*alikafka.StatusList, error) {
	alikafkaStatusList := &alikafka.StatusList{}
	describeNodeStatusReq := alikafka.CreateDescribeNodeStatusRequest()
	describeNodeStatusReq.RegionId = alikafkaService.client.RegionId
	describeNodeStatusReq.InstanceId = instanceId
        describeNodeStatusReq.QueryParams["Product"] = "alikafka"
	wait := incrementalWait(2*time.Second, 1*time.Second)
	var raw interface{}
	var err error
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(client *alikafka.Client) (interface{}, error) {
			return client.DescribeNodeStatus(describeNodeStatusReq)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(describeNodeStatusReq.GetActionName(), raw, describeNodeStatusReq.RpcRequest, describeNodeStatusReq)
		return nil
	})

	if err != nil {
		return alikafkaStatusList, WrapErrorf(err, DefaultErrorMsg, instanceId, describeNodeStatusReq.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	describeNodeStatusResp, _ := raw.(*alikafka.DescribeNodeStatusResponse)
	addDebug(describeNodeStatusReq.GetActionName(), raw, describeNodeStatusReq.RpcRequest, describeNodeStatusReq)

	return &describeNodeStatusResp.StatusList, nil
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaInstanceByOrderId(orderId string, timeout int) (*alikafka.InstanceVO, error) {
	alikafkaInstance := &alikafka.InstanceVO{}
	instanceListReq := alikafka.CreateGetInstanceListRequest()
	instanceListReq.RegionId = alikafkaService.client.RegionId
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
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)
			return nil
		})

		if err != nil {
			return alikafkaInstance, WrapErrorf(err, DefaultErrorMsg, orderId, instanceListReq.GetActionName(), AlibabacloudStackSdkGoERROR)
		}

		instanceListResp, _ := raw.(*alikafka.GetInstanceListResponse)
		addDebug(instanceListReq.GetActionName(), raw, instanceListReq.RpcRequest, instanceListReq)

		for _, v := range instanceListResp.InstanceList.InstanceVO {
			return &v, nil
		}
		if time.Now().After(deadline) {
			return alikafkaInstance, WrapErrorf(Error(GetNotFoundMessage("AlikafkaInstance", orderId)), NotFoundMsg, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaConsumerGroup(id string) (*alikafka.ConsumerVO, error) {
	alikafkaConsumerGroup := &alikafka.ConsumerVO{}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return alikafkaConsumerGroup, WrapError(err)
	}
	instanceId := parts[0]
	//consumerId := parts[1]

	request := alikafka.CreateGetConsumerListRequest()
	request.InstanceId = instanceId
	request.RegionId = alikafkaService.client.RegionId
        request.QueryParams["Product"] = "alikafka"
	wait := incrementalWait(2*time.Second, 1*time.Second)
	var raw interface{}
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(client *alikafka.Client) (interface{}, error) {
			return client.GetConsumerList(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return alikafkaConsumerGroup, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	//consumerListResp, _ := raw.(*alikafka.GetConsumerListResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	//for _, v := range consumerListResp.ConsumerList {
	//	if v.ConsumerId == consumerId {
	//		return &v, nil
	//	}
	//}
	return alikafkaConsumerGroup, WrapErrorf(Error(GetNotFoundMessage("AlikafkaConsumerGroup", id)), NotFoundMsg, ProviderERROR)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaTopicStatus(id string) (*alikafka.TopicStatus, error) {
	//client := alikafkaService.client
	alikafkaTopicStatus := &alikafka.TopicStatus{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return alikafkaTopicStatus, WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := alikafka.CreateGetTopicStatusRequest()
	request.InstanceId = instanceId
	request.RegionId = alikafkaService.client.RegionId
	request.Topic = topic
	request.QueryParams["Product"] = "alikafka"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.GetTopicStatus(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return alikafkaTopicStatus, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	topicStatusResp, _ := raw.(*alikafka.GetTopicStatusResponse)

	if topicStatusResp.TopicStatus.OffsetTable.OffsetTableItem != nil {
		return &topicStatusResp.TopicStatus, nil
	}

	return alikafkaTopicStatus, WrapErrorf(Error(GetNotFoundMessage("AlikafkaTopicStatus "+ResourceNotfound, id)), ResourceNotfound)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaTopic(id string) (*alikafka.TopicList, error) {

	alikafkaTopic := &alikafka.TopicList{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return alikafkaTopic, WrapError(err)
	}
	instanceId := parts[0]
	topic := parts[1]

	request := alikafka.CreateGetTopicListRequest()
	request.InstanceId = instanceId
	request.RegionId = alikafkaService.client.RegionId
	request.Domain = alikafkaService.client.Config.AlikafkaOpenAPIEndpoint
	

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.GetTopicList(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return alikafkaTopic, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	topicListResp, _ := raw.(*alikafka.GetTopicListResponse)

	for _, v := range topicListResp.TopicList {
		if v.Topic == topic {
			return &v, nil
		}
	}
	return alikafkaTopic, WrapErrorf(Error(GetNotFoundMessage("AlikafkaTopic", id)), NotFoundMsg, ProviderERROR)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaSaslUser(id string) (*alikafka.SaslUserList, error) {
	alikafkaSaslUser := &alikafka.SaslUserList{}

	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return alikafkaSaslUser, WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]

	request := alikafka.CreateDescribeSaslUsersRequest()
	request.InstanceId = instanceId
	request.RegionId = alikafkaService.client.RegionId
	request.Domain = alikafkaService.client.Config.AlikafkaOpenAPIEndpoint
	request.QueryParams["Product"] = "alikafka"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DescribeSaslUsers(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return alikafkaSaslUser, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	userListResp, _ := raw.(*alikafka.DescribeSaslUsersResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range userListResp.SaslUserList {
		if v.Username == username {
			return &v, nil
		}
	}
	return alikafkaSaslUser, WrapErrorf(Error(GetNotFoundMessage("AlikafkaSaslUser", id)), NotFoundMsg, ProviderERROR)
}

func (alikafkaService *AlikafkaService) DescribeAlikafkaSaslAcl(id string) (*alikafka.KafkaAclList, error) {
	alikafkaSaslAcl := &alikafka.KafkaAclList{}

	parts, err := ParseResourceId(id, 6)
	if err != nil {
		return alikafkaSaslAcl, WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]
	aclResourceType := parts[2]
	aclResourceName := parts[3]
	aclResourcePatternType := parts[4]
	aclOperationType := parts[5]

	request := alikafka.CreateDescribeAclsRequest()
	request.InstanceId = instanceId
	request.RegionId = alikafkaService.client.RegionId
	request.Username = username
	request.AclResourceType = aclResourceType
	request.AclResourceName = aclResourceName
	request.AclResourcePatternType = aclResourcePatternType
	request.Domain = alikafkaService.client.Config.AlikafkaOpenAPIEndpoint
	request.QueryParams["Product"] = "alikafka"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DescribeAcls(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"BIZ_SUBSCRIPTION_NOT_FOUND", "BIZ_TOPIC_NOT_FOUND"}) {
			return alikafkaSaslAcl, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return alikafkaSaslAcl, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	aclListResp, _ := raw.(*alikafka.DescribeAclsResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range aclListResp.KafkaAclList {
		if v.AclResourcePatternType == aclResourcePatternType && v.AclOperationType == aclOperationType {
			return &v, nil
		}
	}
	return alikafkaSaslAcl, WrapErrorf(Error(GetNotFoundMessage("AlikafkaSaslAcl", id)), NotFoundMsg, ProviderERROR)
}

func (s *AlikafkaService) WaitForAlikafkaInstanceUpdated(id string, topicQuota int, diskSize int, ioMax int,
	eipMax int, paidType int, specType string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaInstance(id)
		if err != nil {
			return WrapError(err)
		}

		// Wait for all variables be equal.
		if object.InstanceId == id && object.TopicNumLimit == topicQuota && object.DiskSize == diskSize && object.IoMax == ioMax && object.EipMax == eipMax && object.PaidType == paidType && object.SpecType == specType {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaInstance(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		// Process wait for running.
		if object.InstanceId == id && status == Running {

			// ServiceStatus equals 5, means the server is in service.
			if object.ServiceStatus == 5 {
				return nil
			}

		} else if object.InstanceId == id {

			// If target status is not deleted and found a instance, return.
			if status != Deleted {
				return nil
			} else {
				// ServiceStatus equals 10, means the server is in released.
				if object.ServiceStatus == 10 {
					return nil
				}
			}
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAllAlikafkaNodeRelease(id string, status string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaNodeStatus(id)
		if err != nil {
			if NotFoundError(err) {
				return nil
			} else {
				return WrapError(err)
			}
		}

		// Process wait for all node become released.
		allReleased := true
		for _, v := range object.Status {
			if v != status && !strings.HasSuffix(v, status) {
				allReleased = false
			}
		}
		if allReleased {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaConsumerGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaConsumerGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.InstanceId+":"+object.ConsumerId == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceId+":"+object.ConsumerId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) KafkaTopicListRefreshFunc(id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlikafkaTopic(id)
		if err != nil {
			if !IsExpectedErrors(err, []string{ResourceNotfound}) {
				return nil, "", WrapError(err)
			}
		}

		return object, "Creating", nil
	}
}

func (s *AlikafkaService) KafkaTopicStatusRefreshFunc(id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAlikafkaTopicStatus(id)
		if err != nil {
			if !IsExpectedErrors(err, []string{ResourceNotfound}) {
				return nil, "", WrapError(err)
			}
		}

		if object.OffsetTable.OffsetTableItem != nil && len(object.OffsetTable.OffsetTableItem) > 0 {
			return object, "Running", WrapError(err)
		}

		return object, "Creating", nil
	}
}

func (s *AlikafkaService) WaitForAlikafkaTopic(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAlikafkaTopic(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.InstanceId+":"+object.Topic == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceId+":"+object.Topic, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaSaslUser(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	for {
		object, err := s.DescribeAlikafkaSaslUser(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if instanceId+":"+object.Username == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, instanceId+":"+object.Username, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) WaitForAlikafkaSaslAcl(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 6)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	for {
		object, err := s.DescribeAlikafkaSaslAcl(id)
		if err != nil {

			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if instanceId+":"+object.Username+":"+object.AclResourceType+":"+object.AclResourceName+":"+object.AclResourcePatternType+":"+object.AclOperationType == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, instanceId+":"+object.Username, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AlikafkaService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []alikafka.TagResource, err error) {
	request := alikafka.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	request.Domain = s.client.Config.AlikafkaOpenAPIEndpoint
	request.QueryParams["Product"] = "alikafka"
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
			if IsExpectedErrors(err, []string{Throttling, ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	response, _ := raw.(*alikafka.ListTagResourcesResponse)

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
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = string(resourceType)
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			request.QueryParams["Product"] = "alikafka"

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				raw, err := s.client.WithAlikafkaClient(func(client *alikafka.Client) (interface{}, error) {
					return client.UntagResources(request)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
			}
		}

		if len(create) > 0 {
			request := alikafka.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = string(resourceType)
			request.RegionId = s.client.RegionId
			request.QueryParams["Product"] = "alikafka"

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				raw, err := s.client.WithAlikafkaClient(func(client *alikafka.Client) (interface{}, error) {
					return client.TagResources(request)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
			}
		}

		d.SetPartial("tags")
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
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []alikafka.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
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
