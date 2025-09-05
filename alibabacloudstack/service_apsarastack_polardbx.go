package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolardbXService struct {
	client *connectivity.AlibabacloudStackClient
}

type PolardbxAccount struct {
	AccountDescription string `json:"AccountDescription"`
	AccountPrivilege   string `json:"AccountPrivilege,omitempty"`
	DBName             string `json:"DBName,omitempty"`
	AccountType        string `json:"AccountType"`
	AccountName        string `json:"AccountName"`
	DBInstanceName     string `json:"DBInstanceName"`
}

type DoPolardbxDescribeAccountListResponse struct {
	EagleEyeTraceId string            `json:"eagleEyeTraceId"`
	AsapiSuccess    bool              `json:"asapiSuccess"`
	AsapiRequestId  string            `json:"asapiRequestId"`
	Message         string            `json:"Message"`
	RequestID       string            `json:"RequestId"`
	Data            []PolardbxAccount `json:"Data"`
	Success         bool              `json:"Success"`
}

func (s *PolardbXService) DoPolardbxDescribeAccountRequest(id string) (*PolardbxAccount, error) {
	var instanceId, AccountName string
	if parts, err := ParseResourceId(id, 2); err != nil {
		return nil, err
	} else {
		instanceId = parts[0]
		AccountName = parts[1]
	}
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeAccountList", "")
	DoPolardbxDescribeAccountListResponseObj := &DoPolardbxDescribeAccountListResponse{}

	// Call request_params_handler

	request.QueryParams["DBInstanceName"] = instanceId

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceAttribute", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DoPolardbxDescribeAccountListResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(DoPolardbxDescribeAccountListResponseObj.Data) > 0 {
		for _, v := range DoPolardbxDescribeAccountListResponseObj.Data {
			if v.AccountName == AccountName {
				return &v, nil
			}
		}
	}

	return nil, errmsgs.Error(errmsgs.NotFoundMsg, "PolardbxAccount")
}

func (s *PolardbXService) DoPolardbxDescribeSuperAccountRequest(id string) ([]PolardbxAccount, error) {
	var adminAccount, securityAccount, auditAccount *PolardbxAccount
	if response, err := s.client.DoTeaRequest("GET", "polardbx", "2020-02-02", "DescribeAccountList", "", nil, map[string]interface{}{"DBInstanceName": id}, nil); err != nil {
		return nil, err
	} else {
		for _, a := range response["Data"].([]interface{}) {
			data := a.(map[string]interface{})
			if data["AccountType"].(string) == "0" {
				continue
			} else {
				account := PolardbxAccount{
					AccountType:    data["AccountType"].(string),
					AccountName:    data["AccountName"].(string),
					DBInstanceName: data["DBInstanceName"].(string),
				}
				if v, exist := data["AccountDescription"]; exist {
					account.AccountDescription = v.(string)
				}
				if data["AccountType"].(string) == "1" || data["AccountType"].(string) == "2" {
					adminAccount = &account
				} else if data["AccountType"].(string) == "3" {
					securityAccount = &account
				} else if data["AccountType"].(string) == "4" {
					auditAccount = &account
				}
			}
		}
	}
	if securityAccount != nil && auditAccount != nil {
		return []PolardbxAccount{*adminAccount, *securityAccount, *auditAccount}, nil
	} else {
		return []PolardbxAccount{*adminAccount}, nil
	}
}

type PolardbxDescribedbinstanceattributeResponse struct {
	RequestId string `json:"RequestId"`

	DBInstance struct {
		ReadDBInstances []string `json:"ReadDBInstances"`

		LTSVersions []string `json:"LTSVersions"`

		DBNodes []struct {
			Id            string `json:"Id"`
			NodeClass     string `json:"NodeClass"`
			RegionId      string `json:"RegionId"`
			ZoneId        string `json:"ZoneId"`
			ComputeNodeId string `json:"ComputeNodeId"`
			DataNodeId    string `json:"DataNodeId"`
		} `json:"DBNodes"`

		ConnAddrs []struct {
			ConnectionString string `json:"ConnectionString"`
			Port             int    `json:"Port"`
			Type             string `json:"Type"`
			VPCId            string `json:"VPCId"`
			VSwitchId        string `json:"VSwitchId"`
			VpcInstanceId    string `json:"VpcInstanceId"`
		} `json:"ConnAddrs"`

		TagSet []struct {
			Key   string `json:"Key"`
			Value string `json:"Value"`
		} `json:"TagSet"`
		Status                  string `json:"Status"`
		Description             string `json:"Description"`
		ZoneId                  string `json:"ZoneId"`
		VPCId                   string `json:"VPCId"`
		CreateTime              string `json:"CreateTime"`
		Expired                 bool   `json:"Expired"`
		PayType                 string `json:"PayType"`
		DBType                  string `json:"DBType"`
		LockMode                string `json:"LockMode"`
		StorageUsed             int    `json:"StorageUsed"`
		Storage                 int    `json:"Storage"`
		DBVersion               string `json:"DBVersion"`
		Network                 string `json:"Network"`
		RegionId                string `json:"RegionId"`
		Engine                  string `json:"Engine"`
		Id                      string `json:"Id"`
		ConnectionString        string `json:"ConnectionString"`
		Port                    int    `json:"Port"`
		MinorVersion            string `json:"MinorVersion"`
		LatestMinorVersion      string `json:"LatestMinorVersion"`
		DBNodeCount             int    `json:"DBNodeCount"`
		DBInstanceType          string `json:"DBInstanceType"`
		SpecSeries              string `json:"SpecSeries"`
		DNNodeCount             int    `json:"DNNodeCount"`
		DNNodeClass             string `json:"DNNodeClass"`
		CNNodeCount             int    `json:"CNNodeCount"`
		CNNodeClass             string `json:"CNNodeClass"`
		MaintainStartTime       string `json:"MaintainStartTime"`
		MaintainEndTime         string `json:"MaintainEndTime"`
		VSwitchId               string `json:"VSwitchId"`
		CommodityCode           string `json:"CommodityCode"`
		ExpireDate              string `json:"ExpireDate"`
		Type                    string `json:"Type"`
		DBNodeClass             string `json:"DBNodeClass"`
		RightsSeparationStatus  string `json:"RightsSeparationStatus"`
		RightsSeparationEnabled bool   `json:"RightsSeparationEnabled"`
		KindCode                int    `json:"KindCode"`
		ResourceGroupId         string `json:"ResourceGroupId"`
		Series                  string `json:"Series"`
		CpuType                 string `json:"CpuType"`
	} `json:"DBInstance"`
}

func (s *PolardbXService) DoPolardbxDescribedbinstanceattributeRequest(id string) (*PolardbxDescribedbinstanceattributeResponse, error) {
	// api: polardbx - 2020-02-02 - DescribeDBInstanceAttribute
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeDBInstanceAttribute", "")
	PolardbxDescribedbinstanceattributeResponseObj := &PolardbxDescribedbinstanceattributeResponse{}

	// Call request_params_handler

	request.QueryParams["DBInstanceName"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceAttribute", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbxDescribedbinstanceattributeResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbxDescribedbinstanceattributeResponseObj, nil
}

func (s *PolardbXService) PolardbxDescribedbinstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		response, err := s.DoPolardbxDescribedbinstanceattributeRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		object := response.DBInstance
		for _, failState := range failStates {
			if fmt.Sprint(object.Status) == failState {
				return object, fmt.Sprint(object.Status), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object.Status)))
			}
		}
		return object, fmt.Sprint(object.Status), nil
	}
}

type PolardbxDescribedbinstancesResponse struct {
	DBInstances []struct {
		ReadDBInstances []string `json:"ReadDBInstances"`

		Nodes []struct {
			Id        string `json:"Id"`
			ClassCode string `json:"ClassCode"`
			RegionId  string `json:"RegionId"`
			ZoneId    string `json:"ZoneId"`
		} `json:"Nodes"`

		TagSet []struct {
			Key   string `json:"Key"`
			Value string `json:"Value"`
		} `json:"TagSet"`
		Id              string `json:"Id"`
		Description     string `json:"Description"`
		PayType         string `json:"PayType"`
		CreateTime      string `json:"CreateTime"`
		ExpireTime      string `json:"ExpireTime"`
		Expired         bool   `json:"Expired"`
		RegionId        string `json:"RegionId"`
		ZoneId          string `json:"ZoneId"`
		Network         string `json:"Network"`
		VPCId           string `json:"VPCId"`
		Engine          string `json:"Engine"`
		DBType          string `json:"DBType"`
		DBVersion       string `json:"DBVersion"`
		Status          string `json:"Status"`
		LockMode        string `json:"LockMode"`
		LockReason      string `json:"LockReason"`
		NodeCount       int    `json:"NodeCount"`
		NodeClass       string `json:"NodeClass"`
		SpecSeries      string `json:"SpecSeries"`
		DNNodeCount     int    `json:"DNNodeCount"`
		DNNodeClass     string `json:"DNNodeClass"`
		CNNodeCount     int    `json:"CNNodeCount"`
		CNNodeClass     string `json:"CNNodeClass"`
		StorageUsed     int    `json:"StorageUsed"`
		Storage         int    `json:"Storage"`
		CommodityCode   string `json:"CommodityCode"`
		Type            string `json:"Type"`
		MinorVersion    string `json:"MinorVersion"`
		ResourceGroupId string `json:"ResourceGroupId"`
		DBInstanceName  string `json:"DBInstanceName"`
		Series          string `json:"Series"`
		CpuType         string `json:"CpuType"`
	} `json:"DBInstances"`
	RequestId   string `json:"RequestId"`
	PageNumber  int    `json:"PageNumber"`
	PageSize    int    `json:"PageSize"`
	TotalNumber int    `json:"TotalNumber"`
}

func (s *PolardbXService) DoPolardbxDescribedbinstancesRequest() (*PolardbxDescribedbinstancesResponse, error) {
	// api: polardbx - 2020-02-02 - DescribeDBInstances
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeDBInstances", "")
	PolardbxDescribedbinstancesResponseObj := &PolardbxDescribedbinstancesResponse{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbxDescribedbinstancesResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbxDescribedbinstancesResponseObj, nil
}

func (s *PolardbXService) DoPolardbxDescribeDBInstanceSSLRequest(id string) (bool, error) {
	// api: polardbx - 2020-02-02 - DescribeDBInstances
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeDBInstanceSSL", "")
	request.QueryParams["DBInstanceName"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return false, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return false, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeDBInstanceSSL", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

	if err != nil {
		return false, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "DescribeDBInstanceSSL", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	enable_ssl, err := jsonpath.Get("$.Data.SSLEnabled", response)
	if err != nil {
		return false, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "DescribeDBInstanceSSL", "$.Data.SSLEnabled", response)
	}
	return enable_ssl.(bool), nil
}

func (s *PolardbXService) DoPolardbxDescribeDBInstanceTDERequest(id string) (bool, error) {
	// api: polardbx - 2020-02-02 - DescribeDBInstances
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeDBInstanceTDE", "")
	request.QueryParams["DBInstanceName"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return false, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return false, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeDBInstanceTDE", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

	if err != nil {
		return false, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "DescribeDBInstanceTDE", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	tde_status, err := jsonpath.Get("$.Data.TDEStatus", response)
	if err != nil {
		return false, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "DescribeDBInstanceTDE", "$.Data.TDEStatus", response)
	}
	return tde_status.(string) == "1", nil
}

func (s *PolardbXService) ModifyParameters(d *schema.ResourceData, attribute string) error {
	request := s.client.NewCommonRequest("POST", "polardbx", "2020-02-02", "ModifyParameter", "")
	request.QueryParams["DBInstanceId"] = d.Id()
	request.QueryParams["ParamLevel"] = attribute
	config := make(map[string]string)
	o, n := d.GetChange(fmt.Sprintf("%s_parameters", attribute))
	os, ns := o.(*schema.Set), n.(*schema.Set)
	add := ns.Difference(os).List()
	var parameters = make(map[string]string)
	if response, err := s.DoPolardbxDescribeParametersRequest(d.Id(), attribute); err != nil {
		return errmsgs.WrapError(err)
	} else {
		running_parameters, err := jsonpath.Get("$.Data.RunningParameters", response)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "DescribeParameters", "$.Data.RunningParameters", response)
		}
		if running_parameters != nil {
			for _, v := range running_parameters.([]interface{}) {
				parameter := v.(map[string]interface{})
				if vv, ok := parameter["ParameterName"]; ok {
					parameters[vv.(string)] = parameter["ParameterValue"].(string)
				}
			}
		}
	}
	if len(add) > 0 {
		for _, i := range add {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			if parameters[key] == value {
				continue
			}
			config[key] = value
		}
		if len(config) == 0 {
			return nil
		}
		cfg, _ := json.Marshal(config)
		request.QueryParams["Parameters"] = string(cfg)
		// wait instance status is Normal before modifying
		bresponse, err := s.client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		stateConf := BuildStateConf([]string{"CONFIG_SWITCHING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, s.PolardbxDescribedbinstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	//d.SetPartial(attribute)
	return nil
}

func (s *PolardbXService) DoPolardbxDescribeParametersRequest(id, param_level string) (map[string]interface{}, error) {
	// api: polardbx - 2020-02-02 - DescribeDBInstances
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeParameters", "")
	request.QueryParams["DBInstanceId"] = id
	request.QueryParams["ParamLevel"] = param_level

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeParameters", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "DescribeParameters", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return response, nil
}

func (s *PolardbXService) RefreshParameters(d *schema.ResourceData, attribute string) error {
	var param []map[string]interface{}
	documented, ok := d.GetOk(fmt.Sprintf("%s_parameters", attribute))
	if !ok {
		d.Set(fmt.Sprintf("%s_parameters", attribute), param)
		return nil
	}
	response, err := s.DoPolardbxDescribeParametersRequest(d.Id(), attribute)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	running_parameters, err := jsonpath.Get("$.Data.RunningParameters", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "DescribeParameters", "$.Data.RunningParameters", response)
	}
	var parameters = make(map[string]interface{})
	if running_parameters != nil {
		for _, v := range running_parameters.([]interface{}) {
			parameter := v.(map[string]interface{})
			if _, ok := v.(map[string]interface{})["ParameterName"]; ok {
				p := map[string]interface{}{
					"name":  parameter["ParameterName"],
					"value": parameter["ParameterValue"],
				}
				parameters[parameter["ParameterName"].(string)] = p
			}
		}
	}
	for _, parameter := range documented.(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, value := range parameters {
			if value.(map[string]interface{})["name"] == name {
				param = append(param, value.(map[string]interface{}))
				break
			}
		}
	}
	if err := d.Set(fmt.Sprintf("%s_parameters", attribute), param); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

type PolardbxDBAccount struct {
	AccountPrivilege string `json:"AccountPrivilege"`
	AccountName      string `json:"AccountName"`
}

type PolardbxDatabase struct {
	CharacterSetName string              `json:"CharacterSetName"`
	DBDescription    string              `json:"DBDescription"`
	DBName           string              `json:"DBName"`
	Accounts         []PolardbxDBAccount `json:"Accounts"`
	DBInstanceName   string              `json:"DBInstanceName"`
}

type PolardbxDescribeDbListResponse struct {
	EagleEyeTraceId string             `json:"eagleEyeTraceId"`
	AsapiSuccess    bool               `json:"asapiSuccess"`
	AsapiRequestId  string             `json:"asapiRequestId"`
	RequestID       string             `json:"RequestId"`
	Message         string             `json:"Message"`
	Data            []PolardbxDatabase `json:"Data"`
	Success         bool               `json:"Success"`
}

func (s *PolardbXService) DoPolardbxDescribeDbListRequest(id string) (*PolardbxDatabase, error) {
	// api: polardbx - 2020-02-02 - DescribeDBInstances
	var InstanceId, databaseName string
	db := &PolardbxDatabase{}
	if parts, err := ParseResourceId(id, 2); err != nil {
		return db, err
	} else {
		InstanceId = parts[0]
		databaseName = parts[1]
	}
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeDbList", "")
	request.QueryParams["DBInstanceName"] = InstanceId
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeDbList", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	var response PolardbxDescribeDbListResponse
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "DescribeDbList", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	for _, v := range response.Data {
		if v.DBName == databaseName {
			db = &v
			break
		}
	}
	if db == nil {
		return nil, errmsgs.Error(errmsgs.NotFoundMsg, "PolardbxDatabase")
	} else {
		return db, nil
	}
}

func (s *PolardbXService) CreatePolardbxAccount(instanceId, accountName, accountPassword, accountDescription string) (err error) {

	action := "CreateAccount"
	reqQuery := map[string]interface{}{
		"DBInstanceName":     instanceId,
		"AccountName":        accountName,
		"AccountType":        "Normal",
		"AccountPassword":    accountPassword,
		"AccountDescription": accountDescription,
	}
	_, err = s.client.DoTeaRequest("POST", "polardbx", "2020-02-02", action, "", nil, reqQuery, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *PolardbXService) DeletePolardbxAccount(instance_id, account_name string) (err error) {
	request := s.client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DeleteAccount", "")
	request.QueryParams["DBInstanceName"] = instance_id
	request.QueryParams["AccountName"] = account_name

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardbx_account", "DeleteAccount", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}

type PolardbXAccountDBPrivilege struct {
	InstanceId  string
	AccountName string
	DBName      string
	Privilege   string
}

func (s *PolardbXService) DescribePolardbXAccountDBPrivilege(id string) ([]map[string]string, error) {
	db_names := make([]string, 0)
	privileges := make([]string, 0)
	var instanceId, AccountName string
	if parts, err := ParseResourceId(id, 2); err != nil {
		return nil, err
	} else {
		instanceId = parts[0]
		AccountName = parts[1]
	}
	account, err := s.DoPolardbxDescribeAccountRequest(fmt.Sprintf("%s:%s", instanceId, AccountName))
	if err != nil {
		return nil, err
	}
	db_privileges := make([]map[string]string, 0)
	if account.DBName != "" {
		db_names = strings.Split(account.DBName, ",")
	}
	if account.AccountPrivilege != "" {
		privileges = strings.Split(account.AccountPrivilege, ",")
	}
	log.Printf("[DEBUG]  DBName: %s   ====================", account.DBName)
	log.Printf("[DEBUG]  AccountPrivilege: %s   ====================", account.AccountPrivilege)
	for i, db_name := range db_names {
		log.Printf("[DEBUG]  db_name: %s   ====================    privilege: %s", db_name, privileges[i])
		db_privileges = append(db_privileges, map[string]string{
			"db_name":   db_name,
			"privilege": privileges[i],
		})
	}
	return db_privileges, nil
}

func (s *PolardbXService) AccountPrivilegeHaschange(account *PolardbxAccount, db_name string, privilege string) (string, string, bool) {
	if account.DBName == "" {
		return db_name, privilege, true
	} else {
		old_db_names := strings.Split(account.DBName, ",")
		old_privileges := strings.Split(account.AccountPrivilege, ",")
		new_db_names := make([]string, 0)
		new_privileges := make([]string, 0)
		for index, name := range old_db_names {
			if name == db_name {
				if privilege == "none" {
					continue
				}
				new_privileges = append(new_privileges, privilege)
			} else {
				new_privileges = append(new_privileges, old_privileges[index])
			}
			new_db_names = append(new_db_names, name)
		}
		new_db_name := strings.Join(new_db_names, ",")
		new_privilege := strings.Join(new_privileges, ",")
		if new_db_name == account.DBName && new_privilege == account.AccountPrivilege {
			return new_db_name, new_privilege, false
		}
		return new_db_name, new_privilege, true
	}
}

func (s *PolardbXService) PolardbxAccountDatabaseBinding(id string, privileges []interface{}) error {
	var instanceId, accountName string
	if parts, err := ParseResourceId(id, 2); err != nil {
		return err
	} else {
		instanceId = parts[0]
		accountName = parts[1]
	}
	new_dbnames := make([]string, 0)
	new_privileges := make([]string, 0)
	for _, db_privilege := range privileges {
		p := db_privilege.(map[string]interface{})
		new_dbnames = append(new_dbnames, p["db_name"].(string))
		new_privileges = append(new_privileges, p["privilege"].(string))
	}
	new_prprivilege_str := strings.Join(new_privileges, ",")
	new_dbname_str := strings.Join(new_dbnames, ",")
	reqQuery := map[string]interface{}{
		"DBInstanceName":   instanceId,
		"AccountName":      accountName,
		"AccountPrivilege": new_prprivilege_str,
		"DbName":           new_dbname_str,
	}
	_, err := s.client.DoTeaRequest("POST", "polardbx", "2020-02-02", "ModifyAccountPrivilege", "", nil, reqQuery, nil)
	if err != nil {
		return err
	}
	return nil
}

type PolardbxSecurityIpGroupItems struct {
	GroupName string `xml:"GroupName" json:"GroupName"`
	IpLists   string `xml:"SecurityIPList" json:"SecurityIPList"`
}

type DescribeSecurityIpsResponse struct {
	RequestId string `xml:"RequestId" json:"RequestId"`
	Success   bool   `xml:"Success" json:"Success"`
	Message   string `xml:"Message" json:"Message"`
	Data      struct {
		DBInstanceName string                         `xml:"DBInstanceName" json:"DBInstanceName"`
		GroupItems     []PolardbxSecurityIpGroupItems `xml:"GroupItems" json:"GroupItems"`
	} `xml:"Data" json:"Data"`
}

func (s *PolardbXService) DoPolardbxDescribeSecurityIpsRequest(id string) (*[]PolardbxSecurityIpGroupItems, error) {
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeAccountList", "")
	DescribeSecurityIpsResponseObj := &DescribeSecurityIpsResponse{}

	// Call request_params_handler

	request.QueryParams["DBInstanceName"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeAccountList", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DescribeSecurityIpsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeAccountList", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return &DescribeSecurityIpsResponseObj.Data.GroupItems, nil
}

func (s *PolardbXService) ModifySecurityIps(instance_id string, old, new []interface{}) error {
	old_securitygroup := make(map[string]interface{})
	new_securitygroup := make(map[string]interface{})
	requests := make([]map[string]interface{}, 0)
	for _, v := range old {
		old_securitygroup[v.(map[string]interface{})["group_name"].(string)] = v.(map[string]interface{})["ips"].(string)
	}
	for _, v := range new {
		new_securitygroup[v.(map[string]interface{})["group_name"].(string)] = v.(map[string]interface{})["ips"].(string)
	}
	for old_group_name, old_ips := range old_securitygroup {
		new_ips, ok := new_securitygroup[old_group_name]
		if !ok {
			requests = append(requests, map[string]interface{}{
				"modifyMode":            "2",
				"GroupName":             old_group_name,
				"DBInstanceIPArrayName": old_group_name,
				"SecurityIPList":        old_ips,
				"DBInstanceName":        instance_id,
			})
		} else if old_ips != new_ips {
			requests = append(requests, map[string]interface{}{
				"modifyMode":     "0",
				"GroupName":      old_group_name,
				"SecurityIPList": new_ips,
				"DBInstanceName": instance_id,
			})
		}

	}
	for new_group_name, new_ips := range new_securitygroup {
		_, ok := new_securitygroup[new_group_name]
		if !ok {
			requests = append(requests, map[string]interface{}{
				"modifyMode":     "1",
				"DBInstanceName": instance_id,
				"SecurityIPList": new_ips,
				"GroupName":      new_group_name,
			})
		}
	}
	for _, request := range requests {
		_, err := s.client.DoTeaRequest("POST", "polardbx", "2020-02-02", "ModifySecurityIps", "", nil, request, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

type PolardbXDBSecurityIPGroup struct {
	GroupName      string `json:"GroupName"`
	SecurityIPList string `json:"SecurityIPList"`
}

type PolardbXDBSecurityIPGroupResponse struct {
	EagleEyeTraceId string `json:"eagleEyeTraceId"`
	AsapiSuccess    bool   `json:"asapiSuccess"`
	AsapiRequestId  string `json:"asapiRequestId"`
	RequestId       string `json:"RequestId"`
	Message         string `json:"Message"`
	Data            struct {
		GroupItems     []PolardbXDBSecurityIPGroup `json:"GroupItems"`
		DBInstanceName string                      `json:"DBInstanceName"`
	} `json:"Data"`
	Success bool `json:"Success"`
}

func (s *PolardbXService) DescribePolardbXDBSecurityIPGroup(instance_id string) ([]PolardbXDBSecurityIPGroup, error) {
	PolardbXDBSecurityIPGroupResponseObj := PolardbXDBSecurityIPGroupResponse{}
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeSecurityIps", "")

	// Call request_params_handler

	request.QueryParams["DBInstanceName"] = instance_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeSecurityIps", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbXDBSecurityIPGroupResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeSecurityIps", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbXDBSecurityIPGroupResponseObj.Data.GroupItems, nil
}

// BackupResponse Backup response structure
type PolarDbXBackupResponse struct {
	EagleEyeTraceId string               `json:"eagleEyeTraceId"`
	AsapiSuccess    bool                 `json:"asapiSuccess"`
	AsapiRequestId  string               `json:"asapiRequestId"`
	Message         string               `json:"Message"`
	RequestId       string               `json:"RequestId"`
	PageSize        int                  `json:"PageSize"`
	PageNumber      int                  `json:"PageNumber"`
	TotalNumber     int                  `json:"TotalNumber"`
	Data            []PolarDbXBackupData `json:"Data"`
	Success         bool                 `json:"Success"`
}

// BackupData Backup data structure
type PolarDbXBackupData struct {
	BackupModel   int    `json:"BackupModel"`
	Status        int    `json:"Status"`
	EndTime       int    `json:"EndTime,omitempty"`
	BeginTime     int    `json:"BeginTime"`
	BackupType    int    `json:"BackupType"`
	BackupSetId   string `json:"BackupSetId"`
	BackupSetSize int    `json:"BackupSetSize"`
}

func (s *PolardbXService) DescribePolarDbXBackups(instance_id string) ([]PolarDbXBackupData, error) {
	PolarDbXBackupResponseObj := PolarDbXBackupResponse{}
	polardbx_backups := make([]PolarDbXBackupData, 0)
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeBackupSetList", "")
	page := 1
	page_size := 50
	for {
		request.QueryParams["PageNumber"] = fmt.Sprint(page)
		request.QueryParams["PageSize"] = fmt.Sprint(page_size)
		request.QueryParams["DBInstanceName"] = instance_id
		bresponse, err := s.client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardbx_backup", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolarDbXBackupResponseObj)

		if err != nil {
			return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardbx_backup", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		polardbx_backups = append(polardbx_backups, PolarDbXBackupResponseObj.Data...)
		if page*page_size >= PolarDbXBackupResponseObj.TotalNumber {
			break
		}
		page++
	}
	return polardbx_backups, nil
}

func (s *PolardbXService) CheckPolarDbXBackupTaskExists(instance_id string) (bool, error) {
	backup_list, err := s.DescribePolarDbXBackups(instance_id)
	if err != nil {
		return false, err
	}
	for _, backup := range backup_list {
		if backup.Status == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (s *PolardbXService) DescribePolarDbXBackup(backup_id string) (*PolarDbXBackupData, error) {
	var instanceId, backupSetId string
	if parts, err := ParseResourceId(backup_id, 2); err != nil {
		return nil, err
	} else {
		instanceId = parts[0]
		backupSetId = parts[1]
	}
	backup_list, err := s.DescribePolarDbXBackups(instanceId)
	if err != nil {
		return nil, err
	}
	for _, backup := range backup_list {
		if backup.BackupSetId == backupSetId {
			return &backup, nil
		}
	}

	return nil, errmsgs.Error(errmsgs.NotFoundMsg, "PolarDbXBackup")
}

func (s *PolardbXService) DescribePolarDbXBackupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDbXBackup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object.Status) == failState {
				return object, fmt.Sprint(object.Status), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object.Status)))
			}
		}
		return object, fmt.Sprint(object.Status), nil
	}
}

type PolarDbXBackupConfigResponse struct {
	EagleEyeTraceId string               `json:"eagleEyeTraceId"`
	AsapiSuccess    bool                 `json:"asapiSuccess"`
	AsapiRequestId  string               `json:"asapiRequestId"`
	Message         string               `json:"Message"`
	RequestId       string               `json:"RequestId"`
	Data            PolarDbXBackupConfig `json:"Data"`
	Success         bool                 `json:"Success"`
}

type PolarDbXBackupConfig struct {
	BackupPeriod               string `json:"BackupPeriod"`
	IsEnabled                  int    `json:"IsEnabled"`
	BackupSetRetention         int    `json:"BackupSetRetention"`
	BackupPlanBegin            string `json:"BackupPlanBegin"`
	ColdDataBackupInterval     int    `json:"ColdDataBackupInterval"`
	RemoveLogRetention         int    `json:"RemoveLogRetention"`
	LocalLogRetentionNumber    int    `json:"LocalLogRetentionNumber"`
	ColdDataBackupRetention    int    `json:"ColdDataBackupRetention"`
	ForceCleanOnHighSpaceUsage int    `json:"ForceCleanOnHighSpaceUsage"`
	BackupWay                  string `json:"BackupWay"`
	LocalLogRetention          int    `json:"LocalLogRetention"`
	BackupType                 string `json:"BackupType"`
	LogLocalRetentionSpace     int    `json:"LogLocalRetentionSpace"`
	DBInstanceName             string `json:"DBInstanceName"`
}

func (s *PolardbXService) DescribePolarDbXBackupConfig(instanceId string) (*PolarDbXBackupConfig, error) {
	PolarDbXBackupConfigResponseObj := PolarDbXBackupConfigResponse{}
	request := s.client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeBackupPolicy", "")
	request.QueryParams["DBInstanceName"] = instanceId
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardbx_backup", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolarDbXBackupConfigResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardbx_backup", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return &PolarDbXBackupConfigResponseObj.Data, nil
}

type PolarDbXLogEngineInfo struct {
	InstanceName string `json:"InstanceName"`
	GroupName    string `json:"GroupName"`
	Comment      string `json:"Comment"`
	HashLevel    string `json:"HashLevel"`
	ClusterType  string `json:"ClusterType"`
	NodeClass    string `json:"NodeClass"`
	NodeCount    int    `json:"NodeCount"`
}

func (s *PolardbXService) DescribePolardbxLogEngine(instanceId string) ([]PolarDbXLogEngineInfo, error) {
	var result []PolarDbXLogEngineInfo
	if response, err := s.client.DoTeaRequest("GET", "polardbx", "2020-02-02", "DescribeCdcInfo", "", nil, map[string]interface{}{"DBInstanceName": instanceId}, nil); err != nil {
		return result, err
	} else {
		for _, item := range response["Data"].(map[string]interface{})["InstanceTopologyList"].([]interface{}) {
			data := item.(map[string]interface{})
			var info PolarDbXLogEngineInfo
			if _, ok := data["GroupName"]; ok {
				info = PolarDbXLogEngineInfo{
					InstanceName: data["InstanceName"].(string),
					GroupName:    data["GroupName"].(string),
					Comment:      data["Comment"].(string),
					HashLevel:    data["HashLevel"].(string),
					NodeClass:    data["PhysicalNodes"].([]interface{})[0].(map[string]interface{})["NodeClass"].(string),
					NodeCount:    len(data["PhysicalNodes"].([]interface{})),
					ClusterType:  data["ClusterType"].(string),
				}
			} else {
				info = PolarDbXLogEngineInfo{
					InstanceName: data["InstanceName"].(string),
					Comment:      data["Comment"].(string),
					NodeClass:    data["PhysicalNodes"].([]interface{})[0].(map[string]interface{})["NodeClass"].(string),
					NodeCount:    len(data["PhysicalNodes"].([]interface{})),
				}
			}
			result = append(result, info)
		}
	}
	return result, nil
}

func (s *PolardbXService) WaitCdcNodeReady(instanceId string) error {
	time.Sleep(10 * time.Second)
	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		if response, err := s.client.DoTeaRequest("GET", "polardbx", "2020-02-02", "DescribeCdcInfo", "", nil, map[string]interface{}{"DBInstanceName": instanceId}, nil); err != nil {
			if errmsgs.NotFoundError(err) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		} else {
			for _, item := range response["Data"].(map[string]interface{})["InstanceTopologyList"].([]interface{}) {
				nodes := item.(map[string]interface{})
				for _, node := range nodes["PhysicalNodes"].([]interface{}) {
					info := node.(map[string]interface{})
					if info["Status"].(string) != "ACTIVATION" {
						time.Sleep(5 * time.Second)
						return resource.RetryableError(fmt.Errorf("Node for instance %s is not ready", nodes["InstanceName"].(string)))
					}
				}
			}
			return nil
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *PolardbXService) ModifyCdcClass(reqQuery map[string]interface{}) error {
	if _, err := s.client.DoTeaRequest("POST", "polardbx", "2020-02-02", "ModifyCdcClass", "", nil, reqQuery, nil); err != nil {
		if sdkError, ok := err.(*errmsgs.ComplexError); ok {
			err = sdkError.Cause
		}
		if sdkError, ok := err.(*tea.SDKError); !ok || *sdkError.Code != "IncorrectTargetClasscode" {
			return err
		}
	}
	if err := s.WaitCdcNodeReady(reqQuery["DBInstanceName"].(string)); err != nil {
		return err
	}
	return nil
}
func (s *PolardbXService) ModifyAccountPassword(instanceId, accountName, accountPassword string) error {
	reqQuery := map[string]interface{}{
		"AccountName":     accountName,
		"DBInstanceName":  instanceId,
		"AccountPassword": accountPassword,
	}
	if _, err := s.client.DoTeaRequest("POST", "polardbx", "2020-02-02", "ResetAccountPassword", "", nil, reqQuery, nil); err != nil {
		return err
	}
	return nil
}

func (s *PolardbXService) ModifyAccountDescription(instanceId, accountName, description string) error {
	request := s.client.NewCommonRequest("POST", "polardbx", "2020-02-02", "ModifyAccountDescription", "")

	request.QueryParams["AccountDescription"] = description
	request.QueryParams["AccountName"] = accountName
	request.QueryParams["DBInstanceName"] = instanceId

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_polardbx_account", "ModifyAccountDescription", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}

func (s *PolardbXService) DescribePolardbxReadWriteSplittingConfig(id string) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"DBInstanceName": id,
		"ConfigName":     "htap",
	}

	response, err := s.client.DoTeaRequest("GET", "polardbx", "2020-02-02", "DescribeDBInstanceConfig", "", nil, query, nil)
	if err != nil {
		return nil, err
	}

	if success, ok := response["asapiSuccess"].(bool); !ok || !success {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("PolarDBX read write splitting config not found for instance %s", id))
	}

	data, ok := response["Data"].(map[string]interface{})
	if !ok || data == nil {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("PolarDBX read write splitting config not found for instance %s", id))
	}

	return data, nil
}
