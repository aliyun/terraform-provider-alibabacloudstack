package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type HologramService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *HologramService) DescribeHologramInstance(id string) (map[string]interface{}, error) {
	pattern := fmt.Sprintf("/api/v1/instances/%s", id)
	response, err := s.client.DoTeaRequest("GET", "Hologram", "2022-06-01", "GetInstance", pattern, nil, nil, nil)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "ListInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	instance, err := jsonpath.Get("$.Instance", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "ListInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if instance != nil {
		return instance.(map[string]interface{}), nil
	}
	return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("hologram_instance", id))
}

func (s *HologramService) HologramInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHologramInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["InstanceStatus"] == failState {
				return object, object["InstanceStatus"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["InstanceStatus"]))
			}
		}
		return object, object["InstanceStatus"].(string), nil
	}
}

func (s *HologramService) CalculateQuota(node int, cluster string) (map[string]interface{}, error) {

	body := map[string]interface{}{
		"node":    node,
		"cluster": cluster,
	}
	request := map[string]interface{}{
		"x-acs-body": body,
	}
	response, err := s.client.DoTeaRequest("POST", "Hologram", "2022-06-01", "CalculateQuota", "/api/v1/instances/calculateQuota", nil, nil, request)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "CalculateQuota", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	quota, err := jsonpath.Get("$.Quota", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "CalculateQuota", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return quota.(map[string]interface{}), nil
}

func (s *HologramService) DescribeHologramInstanceBackup(id string) (map[string]interface{}, error) {
	params := strings.Split(id, ":")
	request := map[string]interface{}{
		"instanceId": params[0],
	}
	response, err := s.client.DoTeaRequest("GET", "Hologram", "2022-06-01", "ListBackupData", "/api/v1/backups", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "ListBackupData", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	backups, err := jsonpath.Get("$.BackupDataList", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "ListBackupData", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	backup := make(map[string]interface{})
	for _, v := range backups.([]interface{}) {
		id := v.(map[string]interface{})["Id"]
		if fmt.Sprint(id) == params[1] {
			backup = v.(map[string]interface{})
		}
	}
	if backup != nil {
		return backup, nil
	} else {
		return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("hologram_instance_backup", id))
	}
}

func (s *HologramService) HologramInstanceBackupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHologramInstanceBackup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"] == failState {
				return object, object["Status"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["Status"]))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *HologramService) DescribeHologramInstanceBackupPolicy(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"instanceId": id,
	}
	response, err := s.client.DoTeaRequest("GET", "Hologram", "2022-06-01", "GetScheduledBackupConfig", "/api/v1/backups/scheduledConfig", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance_backup_policy", "GetScheduledBackupConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	data, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance_backup_policy", "GetScheduledBackupConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	backupConfig := make(map[string]interface{})
	taskParameter := data.(map[string]interface{})["TaskParameter"]
	err = json.Unmarshal([]byte(taskParameter.(string)), &backupConfig)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance_backup_policy", "GetScheduledBackupConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if backupConfig != nil {
		backupConfig["Enabled"] = data.(map[string]interface{})["Enabled"]
		return backupConfig, nil
	} else {
		return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("hologram_instance_backup_policy", id))
	}
}
