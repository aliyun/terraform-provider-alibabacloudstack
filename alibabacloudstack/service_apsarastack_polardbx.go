package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolardbXService struct {
	client *connectivity.AlibabacloudStackClient
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
	request := s.client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DescribeDBInstanceAttribute", "")
	PolardbxDescribedbinstanceattributeResponseObj := &PolardbxDescribedbinstanceattributeResponse{}

	//调用request_params_handler

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

func (s *PolardbXService) DoPolardbxDescribedbinstancesRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbxDescribedbinstancesResponse, error) {
	// api: polardbx - 2020-02-02 - DescribeDBInstances
	request := s.client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DescribeDBInstances", "")
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
	request := s.client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DescribeDBInstanceSSL", "")
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
	request := s.client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DescribeDBInstanceTDE", "")
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
	if len(add) > 0 {
		for _, i := range add {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			config[key] = value
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
	request := s.client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DescribeParameters", "")
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
