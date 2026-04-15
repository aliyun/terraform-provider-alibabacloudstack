package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type MaxcomputeService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *MaxcomputeService) DescribeMaxcomputeProject(id string) (object *MaxComputeProject, err error) {
	client := s.client
	request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", "ListCalcEnginesForAscm", "")
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "ListCalcEnginesForAscm", errmsg)
	}

	response := &MaxComputeProjectDetailResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	for _, v := range response.Data.CalcEngines {
		// EngineStatus == 0 is running   1 is deleted
		if fmt.Sprint(v.EngineId) == id && v.EngineStatus == 0 {
			return &v, nil
		}
	}
	return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("Maxcompute Project", id))
}

func (s *MaxcomputeService) DescribeMaxcomputeProjectEngine(id string) (object *MaxComputeProjectEngineData, err error) {
	client := s.client
	request := client.NewCommonRequest("GET", "dataworks-private-cloud", "2019-01-17", "GetCalcEngineForAscm", "")
	request.QueryParams["EngineId"] = id
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "GetCalcEngineForAscm", errmsg)
	}
	response := MaxComputeProjectEngineDetailResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	return &response.Data, nil
}

func (s *MaxcomputeService) ListOdpsEngineQuotaForAscm(id, name string) (object *OdpsEngineQuotaData, err error) {
	client := s.client
	request := client.NewCommonRequest("GET", "dataworks-private-cloud", "2019-01-17", "ListOdpsEngineQuotaForAscm", "")
	request.QueryParams = map[string]string{
		"Id":          id,
		"ProjectName": name,
	}
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "ListOdpsEngineQuotaForAscm", errmsg)
	}
	response := ListOdpsEngineQuotaForAscmResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	return &response.Data[0], nil
}

func (s *MaxcomputeService) DescribeMaxProjectPropertiesForAscm(id, name string) (map[string]interface{}, error) {
	client := s.client
	request := client.NewCommonRequest("GET", "dataworks-private-cloud", "2019-01-17", "GetOdpsProjectPropertiesForAscm", "")
	request.QueryParams = map[string]string{
		"ProjectName":        name,
		"Properties":         "ENCRYPTION,odps.security.vpc.whitelist",
		"EngineId":           id,
		"DoReplaceTunnelIds": "true",
	}
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "GetOdpsProjectPropertiesForAscm", errmsg)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *MaxcomputeService) DescribeMaxcomputeQuota(cluster_id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"Region":      s.client.RegionId,
		"Action":      "GetOdpsQuotaForAscm",
		"AccessKeyId": s.client.AccessKey,
		"Cluster":     cluster_id,
		"Project":     "odps",
	}

	response, err := s.client.DoTeaRequest("GET", "dataworks-private-cloud", "2019-01-17", "GetOdpsQuotaForAscm", "", nil, request, nil)
	addDebug("ListOdpsCusForAscm", response, request)
	if err != nil {
		err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_project", "GetOdpsQuotaForAscm", errmsgs.AlibabacloudStackSdkGoERROR)
		return
	}

	v, err := jsonpath.Get("$.Data.data", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, cluster_id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MaxcomputeService) DescribeMaxcomputeCu(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"Region":      s.client.RegionId,
		"Action":      "ListOdpsCusForAscm",
		"AccessKeyId": s.client.AccessKey,
		"CuId":        id,
	}

	response, err := s.client.DoTeaRequest("GET", "dataworks-private-cloud", "2019-01-17", "ListOdpsCusForAscm", "", nil, request, nil)
	addDebug("ListOdpsCusForAscm", response, request)
	if err != nil {
		err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_cu", "ListOdpsCusForAscm", errmsgs.AlibabacloudStackSdkGoERROR)
		return
	}
	v, err := jsonpath.Get("$.Data.data", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	datas := v.([]interface{})
	for _, element := range datas {
		if element.(map[string]interface{})["id"].(string) == id {
			return element.(map[string]interface{}), nil
		}
	}
	return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("Maxcompute", id))
}

func (s *MaxcomputeService) DescribeMaxcomputeUsers() ([]interface{}, error) {
	action := "GetOdpsUserList"
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
	request.SetDomain(s.client.Config.Endpoints[connectivity.ASAPICode])

	response := make(map[string]interface{})
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(action, bresponse, request, request.QueryParams)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	datas, err := jsonpath.Get("$.data", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "$.data", response)
	}
	if datas != nil {
		return datas.([]interface{}), nil
	}
	return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
}

func (s *MaxcomputeService) DescribeMaxcomputeUser(id string) (map[string]interface{}, error) {
	action := "GetOdpsUserList"
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
	request.SetDomain(s.client.Config.Endpoints[connectivity.ASAPICode])
	request.QueryParams["Id"] = id

	response := make(map[string]interface{})
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(action, bresponse, request, request.QueryParams)
	if err != nil {
		return response, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return response, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	datas, err := jsonpath.Get("$.data", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.data", response)
	}
	for _, v := range datas.([]interface{}) {
		data := v.(map[string]interface{})
		if fmt.Sprint(data["id"]) == id {
			return data, nil
		}
	}
	return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
}

// DescribeMaxcomputeUserForName 根据用户名查询用户信息
func (s *MaxcomputeService) DescribeMaxcomputeUserForName(username string) (map[string]interface{}, error) {
	users, err := s.DescribeMaxcomputeUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		userMap := user.(map[string]interface{})
		if userMap["userName"].(string) == username {
			return userMap, nil
		}
	}
	return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("Maxcompute User", username))
}

// GetOrCreateMaxcomputeUser 根据用户名获取或创建用户，返回用户ID和AAS PK
func (s *MaxcomputeService) GetOrCreateMaxcomputeUser(username string) (userId string, aasPk string, err error) {
	// 先尝试查询用户是否存在
	userInfo, err := s.DescribeMaxcomputeUserForName(username)
	if err == nil && userInfo != nil {
		// 用户存在，返回 id 和 aasPk
		userId = fmt.Sprintf("%v", userInfo["id"])
		aasPk = fmt.Sprintf("%v", userInfo["aasPk"])
		return userId, aasPk, nil
	}
	// 用户不存在，创建新用户
	action := "CreateOdpsUser"
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
	request.SetDomain(s.client.Config.Endpoints[connectivity.ASAPICode])
	mergeMaps(request.QueryParams, map[string]string{
		"UserName":    username,
		"Description": "Auto created for test",
	})

	response := make(map[string]interface{})
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(action, bresponse, request, request.QueryParams)
	if err != nil {
		return "", "", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return "", "", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// 重新查询用户列表获取新创建的用户信息
	users, err := s.DescribeMaxcomputeUsers()
	if err != nil || len(users) == 0 {
		return "", "", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	for _, user := range users {
		userMap := user.(map[string]interface{})
		if userMap["userName"].(string) == username {
			userId = fmt.Sprintf("%v", userMap["id"])
			aasPk = fmt.Sprintf("%v", userMap["aasPk"])
			return userId, aasPk, nil
		}
	}

	return "", "", errmsgs.Error(errmsgs.GetNotFoundMessage("Maxcompute User after creation", username))
}
