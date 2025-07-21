package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type DrdsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DrdsService) DescribeDrdsInstance(id string) (*drds.DescribeDrdsInstanceResponse, error) {
	request := drds.CreateDescribeDrdsInstanceRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DrdsInstanceId = id
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstance(request)
	})

	response, ok := raw.(*drds.DescribeDrdsInstanceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDrdsInstanceId.NotFound"}) {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if response.Data.Status == "5" {
		return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return response, nil
}

func (s *DrdsService) DrdsInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDrdsInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Data.Status == failState {
				return object, object.Data.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Data.Status))
			}
		}

		return object, object.Data.Status, nil
	}
}

func (s *DrdsService) WaitDrdsInstanceConfigEffect(id string, item map[string]string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		effected := false
		object, err := s.DescribeDrdsInstance(id)

		if err != nil {
			if errmsgs.NotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			return errmsgs.WrapError(err)
		}

		if value, ok := item["description"]; ok {
			if object.Data.Description == value {
				effected = true
			}
		}

		if effected {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Data, item, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

type DrdsDescribedrdsdbResponse struct {
	RequestId string `json:"RequestId"`
	Success   bool   `json:"Success"`

	Data struct {
		DbName     string `json:"DbName"`
		Status     string `json:"Status"`
		CreateTime int64  `json:"CreateTime"`
		Mode       string `json:"Mode"`
		Schema     string `json:"Schema"`
		DbInstType string `json:"DbInstType"`
		InstRole   string `json:"InstRole"`
	} `json:"Data"`
}

func (s *DrdsService) DescribeDrdsDatabase(id string) (*DrdsDescribedrdsdbResponse, error) {

	var instanceId, databaseName string
	if parts, err := ParseResourceId(id, 2); err != nil {
		return nil, err
	} else {
		instanceId = parts[0]
		databaseName = parts[1]
	}

	// api: Drds - 2019-01-23 - DescribeDrdsDB
	request := s.client.NewCommonRequest("GET", "Drds", "2019-01-23", "DescribeDrdsDB", "")
	DrdsDescribedrdsdbResponse := &DrdsDescribedrdsdbResponse{}

	//调用request_params_handler

	request.QueryParams["DbName"] = databaseName

	request.QueryParams["DrdsInstanceId"] = instanceId

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if sdkErr, ok := err.(*errors.ServerError); ok && sdkErr.ErrorCode() == "InvalidDbName.NotFound" {
			return nil, sdkErr
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDrdsDB", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DrdsDescribedrdsdbResponse)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDrdsDB", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return DrdsDescribedrdsdbResponse, nil
}

func (s *DrdsService) DbStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDrdsDatabase(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Data.Status == failState {
				return object, object.Data.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Data.Status))
			}
		}

		return object, object.Data.Status, nil
	}
}

func (s *DrdsService) DescribeDrdsDbTask(id string) (map[string]interface{}, error) {
	var drdsInstanceId, databaseName string
	if parts, err := ParseResourceId(id, 2); err != nil {
		return nil, err
	} else {
		drdsInstanceId = parts[0]
		databaseName = parts[1]
	}

	reqQuery := map[string]interface{}{
		"DbName":         databaseName,
		"DrdsInstanceId": drdsInstanceId,
	}
	repsonse, err := s.client.DoTeaRequest("GET", "Drds", "2019-01-23", "DescribeDrdsDbTasks", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}
	task := repsonse["Tasks"].(map[string]interface{})["Task"].([]interface{})
	if len(task) < 1 {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Not Rearrange Task Found for Drds %s DB %s", drdsInstanceId, databaseName))
	}
	return task[len(task)-1].(map[string]interface{}), nil
}

func (s *DrdsService) DbTaskRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDrdsDbTask(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		taskStatus, err := object["TaskStatus"].(json.Number).Int64()
		if err != nil {
			return nil, "", err
		}
		status := DrdsRearrangeTaskStatus[int(taskStatus)]
		for _, failState := range failStates {
			if status == failState {
				return object, status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, status))
			}
		}

		return object, status, nil
	}
}

type DrdsDescribedrdsdbipwhitelistResponse struct {
	IpWhiteList struct {
		Ip []string `json:"Ip"`
	} `json:"IpWhiteList"`
	RequestId string `json:"RequestId"`
	Success   bool   `json:"Success"`
}

func (s *DrdsService) DoDrdsDescribedrdsdbipwhitelistRequest(id string) (*DrdsDescribedrdsdbipwhitelistResponse, error) {
	var instanceId, databaseName string
	if parts, err := ParseResourceId(id, 2); err != nil {
		return nil, err
	} else {
		instanceId = parts[0]
		databaseName = parts[1]
	}
	// api: Drds - 2019-01-23 - DescribeDrdsDBIpWhiteList
	request := s.client.NewCommonRequest("POST", "Drds", "2019-01-23", "DescribeDrdsDBIpWhiteList", "")
	DrdsDescribedrdsdbipwhitelistResponse := &DrdsDescribedrdsdbipwhitelistResponse{}

	//调用request_params_handler

	request.QueryParams["DbName"] = databaseName

	request.QueryParams["DrdsInstanceId"] = instanceId

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDrdsDBIpWhiteList", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DrdsDescribedrdsdbipwhitelistResponse)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDrdsDBIpWhiteList", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return DrdsDescribedrdsdbipwhitelistResponse, nil
}

type DrdsDescribedrdsdbsResponse struct {
	Data struct {
		Db []struct {
			DbName     string `json:"DbName"`
			Status     string `json:"Status"`
			CreateTime int64  `json:"CreateTime"`
			Mode       string `json:"Mode"`
			Schema     string `json:"Schema"`
			DbInstType string `json:"DbInstType"`
		} `json:"Db"`
	} `json:"Data"`
	RequestId  string `json:"RequestId"`
	Success    bool   `json:"Success"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
	Total      int    `json:"Total"`
}

func (s *DrdsService) DoDrdsDescribedrdsdbsRequest(id string) (*DrdsDescribedrdsdbsResponse, error) {
	// api: Drds - 2019-01-23 - DescribeDrdsDBs
	request := s.client.NewCommonRequest("POST", "Drds", "2019-01-23", "DescribeDrdsDBs", "")
	DrdsDescribedrdsdbsResponse := &DrdsDescribedrdsdbsResponse{}

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDrdsDBs", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DrdsDescribedrdsdbsResponse)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDrdsDBs", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return DrdsDescribedrdsdbsResponse, nil
}

type DrdsDescribeinstanceAccount struct {
	DbPrivileges struct {
					DbPrivilege []struct {
						DbName    string `json:"DbName"`
						Privilege string `json:"Privilege"`
					} `json:"DbPrivilege"`
				} `json:"DbPrivileges"`
				AccountName string `json:"AccountName"`
				Host        string `json:"Host"`
				AccountType int    `json:"AccountType"`
				Description string `json:"Description"`
}

type DrdsDescribeinstanceaccountsResponse struct {
	InstanceAccounts struct {
		InstanceAccount []DrdsDescribeinstanceAccount `json:"InstanceAccount"`
	} `json:"InstanceAccounts"`
	RequestId string `json:"RequestId"`
	Success   bool   `json:"Success"`
}

func (s *DrdsService) DescribeDrdsAccount(id string) (*DrdsDescribeinstanceAccount, error) {
	var instanceId, drdsAccountName string
	if parts, err := ParseResourceId(id, 2); err != nil {
		return nil, err
	} else {
		instanceId = parts[0]
		drdsAccountName = parts[1]
	}
	// api: Drds - 2019-01-23 - DescribeInstanceAccounts
	request := s.client.NewCommonRequest("POST", "Drds", "2019-01-23", "DescribeInstanceAccounts", "")
	drdsDescribeinstanceaccountsResponse := &DrdsDescribeinstanceaccountsResponse{}

	//调用request_params_handler
	request.QueryParams["DrdsInstanceId"] = instanceId

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeInstanceAccounts", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &drdsDescribeinstanceaccountsResponse)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeInstanceAccounts", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	for _, account := range drdsDescribeinstanceaccountsResponse.InstanceAccounts.InstanceAccount{
		if account.AccountName == drdsAccountName {
			return &account, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Can't Found Account %s in Drds %s", drdsAccountName, instanceId))
}

func (s *DrdsService) DescribeDrdsRdsInstance(id string) (map[string]interface{}, error) {

	var drdsInstanceId, rdsInstanceId string
	if parts, err := ParseResourceId(id, 2); err != nil {
		return nil, err
	} else {
		drdsInstanceId = parts[0]
		rdsInstanceId = parts[1]
	}

	reqQuery := map[string]interface{}{
		"DrdsInstanceId": drdsInstanceId,
	}

	response, err := s.client.DoTeaRequest("GET", "Drds", "2019-01-23", "DescribeDrdsRdsInstances", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	dbInstances := response["DbInstances"].(map[string]interface{})["DbInstance"].([]interface{})

	if len(dbInstances) < 1 {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("No private rds for drds Instance %s", drdsInstanceId))
	}

	for _, d := range dbInstances {
		dbInstance := d.(map[string]interface{})
		if dbInstance["DBInstanceId"].(string) != rdsInstanceId {
			continue
		}
		return dbInstance, nil
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Private Rds %s for Drds %s Not Found", rdsInstanceId, drdsInstanceId))

}

// WaitForInstance waits for instance to given status
func (s *DrdsService) PrivateRdsStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDrdsRdsInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		if object == nil {
			return nil, "", nil
		}

		var status string
		if v, err := object["DBInstanceStatus"].(json.Number).Int64(); err != nil {
			return nil, "", err
		} else {
			if v < 0 {
				return nil, "", nil
			}
			status = DrdsPrivateRdsDbInstanceStatus[int(v)]
		}

		for _, failState := range failStates {
			if status == failState {
				return object, status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, status))
			}
		}

		return object, status, nil
	}
}
