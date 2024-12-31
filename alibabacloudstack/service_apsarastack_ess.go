package alibabacloudstack

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type EssService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *EssService) DescribeEssAlarm(id string) (alarm ess.Alarm, err error) {
	request := ess.CreateDescribeAlarmsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.AlarmTaskId = id
	request.MetricType = "system"
	Alarms, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeAlarms(request)
	})
	AlarmsResponse, ok := Alarms.(*ess.DescribeAlarmsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(Alarms.(*ess.DescribeAlarmsResponse).BaseResponse)
		}
		return alarm, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), Alarms, request.RpcRequest, request)
	systemAlarms := AlarmsResponse.AlarmList.Alarm

	if len(systemAlarms) > 0 {
		return systemAlarms[0], nil
	}

	AlarmsRequest := ess.CreateDescribeAlarmsRequest()
	s.client.InitRpcRequest(*AlarmsRequest.RpcRequest)
	AlarmsRequest.AlarmTaskId = id
	AlarmsRequest.MetricType = "custom"
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeAlarms(AlarmsRequest)
	})
	response, ok := raw.(*ess.DescribeAlarmsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return alarm, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, AlarmsRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(AlarmsRequest.GetActionName(), raw, AlarmsRequest.RpcRequest, AlarmsRequest)
	customAlarms := response.AlarmList.Alarm

	if len(customAlarms) > 0 {
		return customAlarms[0], nil
	}
	return alarm, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EssAlarm", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *EssService) DescribeEssLifecycleHook(id string) (hook ess.LifecycleHook, err error) {
	request := ess.CreateDescribeLifecycleHooksRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.LifecycleHookId = &[]string{id}
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeLifecycleHooks(request)
	})
	response, ok := raw.(*ess.DescribeLifecycleHooksResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return hook, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range response.LifecycleHooks.LifecycleHook {
		if v.LifecycleHookId == id {
			return v, nil
		}
	}
	err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EssLifecycleHook", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	return
}

func (s *EssService) WaitForEssLifecycleHook(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssLifecycleHook(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.LifecycleHookId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.LifecycleHookId, id, errmsgs.ProviderERROR)
		}
	}
}

func (s *EssService) DescribeEssNotification(id string) (notification ess.NotificationConfigurationModel, err error) {
	parts := strings.SplitN(id, ":", 2)
	scalingGroupId, notificationArn := parts[0], parts[1]
	request := ess.CreateDescribeNotificationConfigurationsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ScalingGroupId = scalingGroupId
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeNotificationConfigurations(request)
	})
	response, ok := raw.(*ess.DescribeNotificationConfigurationsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EssNotification", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		return notification, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	for _, v := range response.NotificationConfigurationModels.NotificationConfigurationModel {
		if v.NotificationArn == notificationArn {
			return v, nil
		}
	}
	err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EssNotificationConfiguration", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	return
}

func (s *EssService) WaitForEssNotification(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEssNotification(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		resourceId := fmt.Sprintf("%s:%s", object.ScalingGroupId, object.NotificationArn)
		if resourceId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, resourceId, id, errmsgs.ProviderERROR)
		}
	}
}

func (s *EssService) DescribeEssScalingGroup(id string) (group ess.ScalingGroup, err error) {
	request := ess.CreateDescribeScalingGroupsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	var scids []string
	scids = append(scids, id)
	request.ScalingGroupId = &scids
	raw, e := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingGroups(request)
	})
	response, ok := raw.(*ess.DescribeScalingGroupsResponse)
	if e != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return group, errmsgs.WrapErrorf(e, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range response.ScalingGroups.ScalingGroup {
		if v.ScalingGroupId == id {
			return v, nil
		}
	}
	err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EssScalingGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	return
}

func (s *EssService) DescribeEssScalingConfiguration(id string) (config ess.ScalingConfigurationInDescribeScalingConfigurations, err error) {
	request := ess.CreateDescribeScalingConfigurationsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	var scids []string
	scids = append(scids, id)
	request.ScalingConfigurationId = &scids
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(request)
	})
	response, ok := raw.(*ess.DescribeScalingConfigurationsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return config, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	for _, v := range response.ScalingConfigurations.ScalingConfiguration {
		if v.ScalingConfigurationId == id {
			return v, nil
		}
	}

	err = errmsgs.GetNotFoundErrorFromString(errmsgs.GetNotFoundMessage("Scaling Configuration", id))
	return
}

func (s *EssService) ActiveEssScalingConfiguration(sgId, id string) error {
	request := ess.CreateModifyScalingGroupRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ScalingGroupId = sgId
	request.ActiveScalingConfigurationId = id
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingGroup(request)
	})
	response, ok := raw.(*ess.ModifyScalingGroupResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return err
}

func (s *EssService) WaitForScalingConfiguration(id string, status Status, timeout int) (err error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssScalingConfiguration(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if object.ScalingConfigurationId == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.ScalingConfigurationId, id, errmsgs.ProviderERROR)
		}
	}
}

// Flattens an array of datadisk into a []map[string]interface{}
func (s *EssService) flattenDataDiskMappings(list []ess.DataDisk) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"size":                 i.Size,
			"category":             i.Category,
			"snapshot_id":          i.SnapshotId,
			"device":               i.Device,
			"delete_with_instance": i.DeleteWithInstance,
			"encrypted":            i.Encrypted,
			//			"kms_key_id":              i.KMSKeyId,
			//			"disk_name":               i.DiskName,
			//			"description":             i.Description,
			//			"auto_snapshot_policy_id": i.AutoSnapshotPolicyId,
		}
		result = append(result, l)
	}
	return result
}

func (s *EssService) flattenVserverGroupList(vServerGroups []ess.VServerGroup) []map[string]interface{} {
	groups := make([]map[string]interface{}, 0, len(vServerGroups))
	for _, v := range vServerGroups {
		vserverGroupAttributes := v.VServerGroupAttributes.VServerGroupAttribute
		attrs := make([]map[string]interface{}, 0, len(vserverGroupAttributes))
		for _, a := range vserverGroupAttributes {
			attr := map[string]interface{}{
				"vserver_group_id": a.VServerGroupId,
				"port":             a.Port,
				"weight":           a.Weight,
			}
			attrs = append(attrs, attr)
		}
		group := map[string]interface{}{
			"loadbalancer_id":    v.LoadBalancerId,
			"vserver_attributes": attrs,
		}
		groups = append(groups, group)
	}
	return groups
}

func (s *EssService) DescribeEssScalingRule(id string) (rule ess.ScalingRule, err error) {
	request := ess.CreateDescribeScalingRulesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	var ruleIds []string
	ruleIds = append(ruleIds, id)
	request.ScalingRuleId = &ruleIds
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingRules(request)
	})
	response, ok := raw.(*ess.DescribeScalingRulesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"InvalidScalingRuleId.NotFound"}) {
			return rule, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return rule, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range response.ScalingRules.ScalingRule {
		if v.ScalingRuleId == id {
			return v, nil
		}
	}

	return rule, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EssScalingRule", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *EssService) WaitForEssScalingRule(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssScalingRule(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errmsgs.WrapError(err)
		}

		if object.ScalingRuleId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.ScalingRuleId, id, errmsgs.ProviderERROR)
		}
	}
}

func (s *EssService) DescribeEssScheduledTask(id string) (task ess.ScheduledTask, err error) {
	request := ess.CreateDescribeScheduledTasksRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	var taskIds []string
	taskIds = append(taskIds, id)
	request.ScheduledTaskId = &taskIds
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScheduledTasks(request)
	})
	response, ok := raw.(*ess.DescribeScheduledTasksResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return task, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, v := range response.ScheduledTasks.ScheduledTask {
		if v.ScheduledTaskId == id {
			task = v
			return
		}
	}
	err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EssSchedule", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	return
}

func (s *EssService) WaitForEssScheduledTask(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEssScheduledTask(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if object.TaskEnabled {
			return nil
		}

		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.ScheduledTaskId, id, errmsgs.ProviderERROR)
		}
	}
}

func (srv *EssService) DoEssDescribescalinggroupsRequest(id string, instanceIds []string) (instances []ess.ScalingInstance, err error) {
	return srv.DescribeEssAttachment(id, instanceIds)
}

func (srv *EssService) DescribeEssAttachment(id string, instanceIds []string) (instances []ess.ScalingInstance, err error) {
	request := ess.CreateDescribeScalingInstancesRequest()
	srv.client.InitRpcRequest(*request.RpcRequest)
	request.ScalingGroupId = id
	s := reflect.ValueOf(request).Elem()

	if len(instanceIds) > 0 {
		for i, id := range instanceIds {
			s.FieldByName(fmt.Sprintf("InstanceId%d", i+1)).Set(reflect.ValueOf(id))
		}
	}

	raw, err := srv.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingInstances(request)
	})
	response, ok := raw.(*ess.DescribeScalingInstancesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
			err = errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		} else {
			return instances, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(response.ScalingInstances.ScalingInstance) < 1 {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EssAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		return
	}
	return response.ScalingInstances.ScalingInstance, nil
}

func (s *EssService) DescribeEssScalingConfifurations(id string) (configs []ess.ScalingConfigurationInDescribeScalingConfigurations, err error) {
	request := ess.CreateDescribeScalingConfigurationsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ScalingGroupId = id
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	for {
		raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingConfigurations(request)
		})
		response, ok := raw.(*ess.DescribeScalingConfigurationsResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return configs, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
			break
		}
		configs = append(configs, response.ScalingConfigurations.ScalingConfiguration...)
		if len(response.ScalingConfigurations.ScalingConfiguration) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return configs, errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	if len(configs) < 1 {
		return configs, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EssScalingConfifuration", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return
}

func (srv *EssService) EssRemoveInstances(id string, instanceIds []string) error {

	if len(instanceIds) < 1 {
		return nil
	}
	group, err := srv.DescribeEssScalingGroup(id)

	if err != nil {
		return errmsgs.WrapError(err)
	}

	if group.LifecycleState == string(Inactive) {
		return errmsgs.WrapError(errmsgs.Error("Scaling group current status is %s, please active it before attaching or removing ECS instances.", group.LifecycleState))
	} else {
		if err := srv.WaitForEssScalingGroup(group.ScalingGroupId, Active, DefaultTimeout); err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return errmsgs.WrapError(err)
		}
	}

	removed := instanceIds
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := ess.CreateRemoveInstancesRequest()
		srv.client.InitRpcRequest(*request.RpcRequest)
		if len(removed) > 0 {
			request.InstanceId = &removed
		} else {
			return nil
		}
		raw, err := srv.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.RemoveInstances(request)
		})
		if err != nil {
			errmsg := ""
			if _, ok := raw.(*ess.RemoveInstancesResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(raw.(*ess.RemoveInstancesResponse).BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
				return nil
			}
			if errmsgs.IsExpectedErrors(err, []string{"IncorrectCapacity.MinSize"}) {
				instances, err := srv.DescribeEssAttachment(id, instanceIds)
				if len(instances) > 0 {
					if group.MinSize == 0 {
						return resource.RetryableError(errmsgs.WrapError(err))
					}
					return resource.NonRetryableError(errmsgs.WrapError(err))
				}
			}
			if errmsgs.IsExpectedErrors(err, []string{"ScalingActivityInProgress", "IncorrectScalingGroupStatus"}) {
				time.Sleep(5)
				return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		time.Sleep(3 * time.Second)
		instances, err := srv.DescribeEssAttachment(id, instanceIds)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(errmsgs.WrapError(err))
		}
		if len(instances) > 0 {
			removed = make([]string, 0)
			for _, inst := range instances {
				removed = append(removed, inst.InstanceId)
			}
			return resource.RetryableError(errmsgs.WrapError(errmsgs.Error("There are still ECS instances in the scaling group.")))
		}

		return nil
	}); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

// ess dimensions to map
func (s *EssService) flattenDimensionsToMap(dimensions []ess.Dimension) map[string]string {
	result := make(map[string]string)
	for _, dimension := range dimensions {
		if dimension.DimensionKey == UserId || dimension.DimensionKey == ScalingGroup {
			continue
		}
		result[dimension.DimensionKey] = dimension.DimensionValue
	}
	return result
}

func (s *EssService) WaitForEssScalingGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssScalingGroup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.LifecycleState == string(status) {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.LifecycleState, string(status), errmsgs.ProviderERROR)
		}
	}
}

func (s *EssService) WaitForEssAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssAttachment(id, make([]string, 0))
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if len(object) > 0 && status != Deleted {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), errmsgs.ProviderERROR)
		}
	}
}

func (s *EssService) WaitForEssAlarm(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEssAlarm(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.AlarmTaskId == id && status != Deleted {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.AlarmTaskId, id, errmsgs.ProviderERROR)
		}
	}
}
