package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type AqsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *AqsService) DescribeAqsOssScanConfig(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"Id":   id,
		"From": "sas",
	}
	response, err := s.client.DoTeaRequest("GET", "aegis", "2016-11-11", "GetOssScanConfig", "", nil, request, nil)
	if err != nil {
		return nil, err
	}
	data, ok := response["Data"].(map[string]interface{})
	if !ok || data == nil {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Aqs OssScanConfig %s not found", id))
	}
	return data, nil
}

func (s *AqsService) DescribeAntiBruteForceRule(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"Id":   id,
		"From": "sas",
	}

	action := "DescribeAntiBruteForceRules"
	response, err := s.client.DoTeaRequest("GET", "aegis", "2016-11-11", action, "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	rules, ok := response["Rules"].([]interface{})
	if !ok || len(rules) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString("Resource not found")
	}
	for _, v := range rules {
		rule := v.(map[string]interface{})
		if fmt.Sprint(rule["Id"]) == id {
			return rule, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString("Resource not found")
}

func (s *AqsService) DescribeCloudCenterInstances() ([]interface{}, error) {
	request := map[string]interface{}{
		"From":         "sas",
		"MachineTypes": "ecs",
	}
	response, err := s.client.DoTeaRequest("GET", "aegis", "2016-11-11", "DescribeCloudCenterInstances", "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	instances, ok := response["Instances"].([]interface{})
	if !ok || len(instances) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString("Resource not found")
	}
	return instances, nil
}

func (s *AqsService) RefreshAssets(assetType string) error {
	request := map[string]interface{}{
		"From":      "sas",
		"AssetType": assetType,
	}
	_, err := s.client.DoTeaRequest("POST", "aegis", "2016-11-11", "RefreshAssets", "", nil, request, nil)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 10)
	return nil
}

func (s *AqsService) DescribeWebLockInstance(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"From":        "sas",
		"CurrentPage": 1,
		"PageSize":    100,
	}

	response, err := s.client.DoTeaRequest("GET", "aegis", "2016-11-11", "DescribeWebLockBindList", "", nil, request, nil)
	if err != nil {
		return nil, err
	}
	bindlist, ok := response["BindList"].([]interface{})
	if !ok && len(bindlist) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString("Resource not found")
	}
	for _, v := range bindlist {
		bind := v.(map[string]interface{})
		if bind["Uuid"] == id {
			return bind, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString("Resource not found")
}

func (s *AqsService) ListWebLockConfigs(id string) ([]interface{}, error) {
	request := map[string]interface{}{
		"From": "sas",
		"Uuid": id,
	}

	response, err := s.client.DoTeaRequest("GET", "aegis", "2016-11-11", "DescribeWebLockConfigList", "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	configList, ok := response["ConfigList"].([]interface{})
	if !ok || len(configList) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString("WebLockConfig not found")
	}

	// result := map[string]interface{}{
	// 	"uuid":                config["Uuid"],
	// 	"dir":                 config["Dir"],
	// 	"inclusive_file_type": config["InclusiveFileType"],
	// 	"exclusive_file":      config["ExclusiveFile"],
	// 	"exclusive_dir":       config["ExclusiveDir"],
	// 	"defence_mode":        config["DefenceMode"],
	// 	"mode":                config["Mode"],
	// 	"local_backup_dir":    config["LocalBackupDir"],
	// 	"exclusive_file_type": config["ExclusiveFileType"],
	// 	"id":                  config["Id"],
	// }

	return configList, nil
}

func (s *AqsService) DeleteWebLockConfig(uuid string, configId int) error {
	deleteConReq := map[string]interface{}{
		"From": "sas",
		"Uuid": uuid,
		"Id":   configId,
	}
	_, err := s.client.DoTeaRequest("GET", "aegis", "2016-11-11", "ModifyWebLockDeleteConfig", "", nil, deleteConReq, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
