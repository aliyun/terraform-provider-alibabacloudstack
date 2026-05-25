package alibabacloudstack

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var SlsClientTimeoutCatcher = Catcher{errmsgs.LogClientTimeout, 15, 5}

type LogService struct {
	client *connectivity.AlibabacloudStackClient
}

type LogProject struct {
	CreateTime         string   `json:"createTime"`
	LastModifyTime     string   `json:"lastModifyTime"`
	Description        string   `json:"description"`
	Owner              string   `json:"owner"`
	ProjectName        string   `json:"projectName"`
	Region             string   `json:"region"`
	Status             string   `json:"status"`
	ResourceGroupId    string   `json:"resourceGroupId"`
	DataRedundancyType string   `json:"dataRedundancyType"`
	ClusterName        string   `json:"clusterName"`
	DataEndpoint       string   `json:"dataEndpoint"`
	RecycleBinEnabled  bool     `json:"recycleBinEnabled"`
	Location           string   `json:"location"`
	Quota              LogQuota `json:"quota"`
}

type LogQuota struct {
	LogStore       int `json:"LogStore"`
	Report         int `json:"Report"`
	SavedSearch    int `json:"SavedSearch"`
	Config         int `json:"Config"`
	MinWriteQuota  int `json:"MinWriteQuota"`
	Dashboard      int `json:"Dashboard"`
	MinReadQuota   int `json:"MinReadQuota"`
	Ingestion      int `json:"Ingestion"`
	MachineGroup   int `json:"MachineGroup"`
	ScheduledSQL   int `json:"ScheduledSQL"`
	Export         int `json:"Export"`
	Alert          int `json:"Alert"`
	ETL            int `json:"ETL"`
	Shard          int `json:"Shard"`
	Chart          int `json:"Chart"`
	MinLengthQuota int `json:"MinLengthQuota"`
}

func (s *LogService) DescribeLogProject(id string) (*LogProject, error) {
	var err error
	var logProject *LogProject

	// Attempt 1: New API (2019-10-23)
	requestBody := map[string]interface{}{"ProjectName": id}
	requestHeaders := map[string]string{"AccessKeyId": s.client.AccessKey}
	response, err := s.client.DoTeaRequest("GET", "Sls", "2020-03-31", "GetProject", "/sls/v1/project/getProject", requestHeaders, requestBody, nil)

	// If new API fails, fallback to old API
	if err != nil && errmsgs.IsExpectedErrors(err, "InvalidVersion") {
		// TODO: Remove old API logic in version 3.20.0
		log.Printf("[WARN] SLS 2019-10-23 GetProject failed: %v, fallback to 2020-03-31 API", err)

		// Attempt 2: Old API (2020-03-31) - will be removed in 3.20.0
		request := s.client.NewCommonRequest("POST", "SLS", "2020-03-31", "GetProject", "")
		request.SetDomain(s.client.Config.Endpoints[connectivity.ASAPICode])
		request.QueryParams["ProjectName"] = id

		bresponse, err := s.client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "ProjectNotExist") {
				return logProject, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
			}
			if bresponse == nil {
				return logProject, errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return logProject, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &logProject)
	} else if err != nil {
		if errcode, ok := response["errorCode"]; ok && errcode.(string) == "ProjectNotExist" {
			return logProject, errmsgs.GetNotFoundErrorFromString("LogProject not found")
		}
		return logProject, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetProject", errmsgs.AlibabacloudStackSdkGoERROR, err)
	} else {
		// Parse response from DoTeaRequest
		if responseBytes, err := json.Marshal(response); err != nil {
			return logProject, err
		} else {
			err = json.Unmarshal(responseBytes, &logProject)
		}
	}

	// if logProject != nil && logProject.ProjectName == "" && len(logProject.LogQuota) > 0 {
	// 	for _, k := range logProject.LogQuota {
	// 		if k.ProjectName == id {
	// 			logProject.ProjectName = k.ProjectName
	// 			logProject.Description = k.Description
	// 			logProject.ClusterName = k.ClusterName
	// 			break
	// 		}
	// 	}
	// }
	if logProject == nil || logProject.ProjectName == "" {
		return logProject, errmsgs.GetNotFoundErrorFromString("LogProject not found")
	}
	return logProject, nil
}

func (s *LogService) GetSlsDataClient(projectName string) (slsClient sls.ClientInterface, err error) {
	if s.client.Config.Proxy != "" {
		os.Setenv("http_proxy", s.client.Config.Proxy)
		os.Setenv("https_proxy", s.client.Config.Proxy)
	}
	var endpoint string
	if s.client.Config.SlsDataEndpoint != "" {
		endpoint = s.client.Config.SlsDataEndpoint
	} else {
		project, err := s.DescribeLogProject(projectName)
		if err != nil {
			return slsClient, errmsgs.WrapError(err)
		}
		endpoint = project.DataEndpoint
	}
	// Remove protocol prefix from endpoint to avoid signature mismatch
	endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "https://"), "http://")
	slsClient = sls.CreateNormalInterface(endpoint, s.client.Config.AccessKey, s.client.Config.SecretKey, s.client.Config.SecurityToken)
	return slsClient, nil
}

func (s *LogService) WaitForLogProject(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeLogProject(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.ProjectName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.ProjectName, id, errmsgs.ProviderERROR)
		}
	}
}

func (s *LogService) DescribeLogStore(id string) (*sls.LogStore, error) {
	store := &sls.LogStore{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	projectName, name := parts[0], parts[1]

	slsClient, err := s.GetSlsDataClient(projectName)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := slsClient.GetLogStore(projectName, name)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError", errmsgs.LogClientTimeout) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetLogStore", raw, map[string]string{
				"project":  projectName,
				"logstore": name,
			})
		}
		store = raw
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ProjectNotExist", "LogStoreNotExist") {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "GetLogStore", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	if store == nil || store.Name == "" {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("LogStore", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return store, nil
}

func (s *LogService) WaitForLogStore(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogStore(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Name == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Name, name, errmsgs.ProviderERROR)
		}
	}
}

func (s *LogService) DescribeLogStoreIndex(id string) (*sls.Index, error) {
	index := &sls.Index{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return index, errmsgs.WrapError(err)
	}
	projectName, name := parts[0], parts[1]
	client, err := s.GetSlsDataClient(projectName)
	if err != nil {
		return index, errmsgs.WrapError(err)
	}
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.GetIndex(projectName, name)

		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError", errmsgs.LogClientTimeout) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetIndex", raw, map[string]string{
				"project":  projectName,
				"logstore": name,
			})
		}
		index = raw
		return nil
	})

	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ProjectNotExist", "LogStoreNotExist", "IndexConfigNotExist") {
			return index, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return index, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "GetIndex", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	if index == nil || (index.Line == nil && index.Keys == nil) {
		return index, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("LogStoreIndex", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return index, nil
}

func (s *LogService) DescribeLogMachineGroup(id string) (*sls.MachineGroup, error) {
	group := &sls.MachineGroup{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return group, errmsgs.WrapError(err)
	}
	projectName, groupName := parts[0], parts[1]

	slsClient, err := s.GetSlsDataClient(projectName)
	if err != nil {
		return group, errmsgs.WrapError(err)
	}

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := slsClient.GetMachineGroup(projectName, groupName)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError", errmsgs.LogClientTimeout) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetMachineGroup", raw, map[string]string{
				"project":      projectName,
				"machineGroup": groupName,
			})
		}
		group = raw
		return nil
	})

	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ProjectNotExist", "GroupNotExist", "MachineGroupNotExist") {
			return group, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return group, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "GetMachineGroup", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	if group == nil || group.Name == "" {
		return group, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("LogMachineGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return group, nil
}

func (s *LogService) WaitForLogMachineGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogMachineGroup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Name == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Name, name, errmsgs.ProviderERROR)
		}
	}
}

func (s *LogService) DescribeLogtailConfig(id string) (*sls.LogConfig, error) {
	response := &sls.LogConfig{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	projectName, configName := parts[0], parts[2]

	slsClient, err := s.GetSlsDataClient(projectName)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := slsClient.GetConfig(projectName, configName)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetConfig", raw, map[string]string{
				"project": projectName,
				"config":  configName,
			})
		}
		response = raw
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ProjectNotExist", "LogStoreNotExist", "ConfigNotExist") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "GetConfig", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	if response == nil || response.Name == "" {
		return response, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("LogTailConfig", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return response, nil
}

func (s *LogService) WaitForLogtailConfig(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	name := parts[2]
	for {
		object, err := s.DescribeLogtailConfig(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Name == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Name, name, errmsgs.ProviderERROR)
		}
	}
}

func (s *LogService) DescribeLogtailAttachment(id string) (groupName string, err error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return groupName, errmsgs.WrapError(err)
	}
	projectName, configName, name := parts[0], parts[1], parts[2]

	slsClient, err := s.GetSlsDataClient(projectName)
	if err != nil {
		return groupName, errmsgs.WrapError(err)
	}

	var groupNames []string
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := slsClient.GetAppliedMachineGroups(projectName, configName)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetAppliedMachineGroups", raw, map[string]string{
				"project":  projectName,
				"confName": configName,
			})
		}
		groupNames = raw
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ProjectNotExist", "ConfigNotExist", "MachineGroupNotExist") {
			return groupName, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return groupName, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "GetAppliedMachineGroups", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	for _, group_name := range groupNames {
		if group_name == name {
			groupName = group_name
		}
	}
	if groupName == "" {
		return groupName, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("LogtailAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return groupName, nil
}

func (s *LogService) WaitForLogtailAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	name := parts[2]
	for {
		object, err := s.DescribeLogtailAttachment(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object, name, errmsgs.ProviderERROR)
		}
	}
}

func (s *LogService) DescribeLogAlert(id string) (*sls.Alert, error) {
	alert := &sls.Alert{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return alert, errmsgs.WrapError(err)
	}
	projectName, alertName := parts[0], parts[1]

	slsClient, err := s.GetSlsDataClient(projectName)
	if err != nil {
		return alert, errmsgs.WrapError(err)
	}

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := slsClient.GetAlert(projectName, alertName)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError", errmsgs.LogClientTimeout) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetLogstoreAlert", raw, map[string]string{
				"project":    projectName,
				"alert_name": alertName,
			})
		}
		alert = raw
		return nil
	})

	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ProjectNotExist", "JobNotExist") {
			return alert, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return alert, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "GetLogstoreAlert", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	if alert == nil || alert.Name == "" {
		return alert, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("LogstoreAlert", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return alert, nil
}

func (s *LogService) WaitForLogstoreAlert(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogAlert(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Name == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Name, name, errmsgs.ProviderERROR)
		}
	}
}

func (s *LogService) CreateLogDashboard(project, name string) error {
	dashboard := sls.Dashboard{
		DashboardName: name,
		ChartList:     []sls.Chart{},
	}

	slsClient, err := s.GetSlsDataClient(project)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		err := slsClient.CreateDashboard(project, dashboard)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError", errmsgs.LogClientTimeout) {
				return resource.RetryableError(err)
			}
			if err.(*sls.Error).Message == "specified dashboard already exists" {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("CreateLogDashboard", dashboard, map[string]string{
				"project":        project,
				"dashboard_name": name,
			})
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateDashboard(project, name string, client sls.ClientInterface) error {
	dashboard := sls.Dashboard{
		DashboardName: name,
		ChartList:     []sls.Chart{},
	}
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		err := client.CreateDashboard(project, dashboard)
		if err != nil {
			if err.(*sls.Error).Message == "specified dashboard already exists" {
				return nil
			}
			if errmsgs.IsExpectedErrors(err, "InternalServerError", errmsgs.LogClientTimeout) {
				return resource.RetryableError(err)
			}

		}
		if debugOn() {
			addDebug("CreateLogDashboard", dashboard, map[string]string{
				"project":        project,
				"dashboard_name": name,
			})
		}
		return nil
	})
	return err
}

func GetCharTitile(project, dashboard, char string, client *sls.Client) string {
	board, err := client.GetDashboard(project, dashboard)
	// If the query fails to ignore the error, return the original value.
	if err != nil {
		return char
	}
	for _, v := range board.ChartList {
		if v.Display.DisplayName == char {
			return v.Title
		} else {
			return char
		}

	}
	return char
}

func (s *LogService) DescribeLogDashboard(id string) (*sls.Dashboard, error) {
	dashboard := &sls.Dashboard{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return dashboard, errmsgs.WrapError(err)
	}
	projectName, dashboardName := parts[0], parts[1]

	slsClient, err := s.GetSlsDataClient(projectName)
	if err != nil {
		return dashboard, errmsgs.WrapError(err)
	}

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := slsClient.GetDashboard(projectName, dashboardName)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InternalServerError", errmsgs.LogClientTimeout) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetLogstoreDashboard", raw, map[string]string{
				"project":        projectName,
				"dashboard_name": dashboardName,
			})
		}
		dashboard = raw
		return nil
	})

	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ProjectNotExist", "DashboardNotExist") {
			return dashboard, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
		}
		return dashboard, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "GetLogstoreDashboard", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	if dashboard == nil || dashboard.DashboardName == "" {
		return dashboard, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("LogstoreDashboard", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return dashboard, nil
}

func (s *LogService) WaitForLogDashboard(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogDashboard(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.DashboardName == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.DashboardName, name, errmsgs.ProviderERROR)
		}
	}
}
