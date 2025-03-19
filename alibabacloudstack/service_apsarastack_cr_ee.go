package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"errors"
	"strings"
	"strconv"
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
)

func (c *CrService) ListCrEeInstances(pageNo int, pageSize int) (map[string]interface {}, error) {
	request := c.client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "ListInstance", "")
	
	bresponse, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return nil, fmt.Errorf("read ee repo failed, %s", response["asapiErrorMessage"].(string))
	}
	repoList := response["Instances"].([]interface{})
	if len(repoList) == 0 {
		return nil, errmsgs.WrapError(fmt.Errorf("cr-ee instance not found"))
	}

	return response, nil
}

func (c *CrService) DescribeCrEeInstance(instanceId string) (*cr_ee.GetInstanceResponse, error) {
	request := cr_ee.CreateGetInstanceRequest()
	c.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	resource := instanceId
	action := request.GetActionName()

	raw, err := c.client.WithCrEeClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.GetInstance(request)
	})
	response, ok := raw.(*cr_ee.GetInstanceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, resource, action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(action, raw, request.RpcRequest, request)

	if !response.GetInstanceIsSuccess {
		return response, c.wrapCrServiceError(resource, action, response.Code)
	}
	return response, nil
}

func (c *CrService) GetCrEeInstanceUsage(instanceId string) (map[string]interface{}, error) {
	
	request := c.client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "GetInstanceUsage", "")
	request.QueryParams["InstanceId"] = instanceId
	
	bresponse, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if errmsgs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return nil, fmt.Errorf("read ee repo failed, %s", response["asapiErrorMessage"].(string))
	}

	return response, nil
}

func (c *CrService) ListCrEeInstanceEndpoint(instanceId string) (map[string]interface{}, error) {
	
	request := c.client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "ListInstanceEndpoint", "")
	request.QueryParams["InstanceId"] = instanceId
	
	bresponse, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if errmsgs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return nil, fmt.Errorf("read ee repo failed, %s", response["asapiErrorMessage"].(string))
	}

	return response, nil
}

func (c *CrService) ListCrEeNamespaces(instanceId string, pageNo int, pageSize int) (map[string]interface{}, error) {
	request := c.client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "ListNamespace", "")
	request.QueryParams["InstanceId"] = instanceId
	request.QueryParams["PageNo"] = strconv.Itoa(pageNo)
	request.QueryParams["PageSize"] =  strconv.Itoa(pageSize)

	bresponse, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil,  errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil,  errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return nil, fmt.Errorf("read ee namespace failed, %s", response["errorMessage"].(string))
	}
	return response, nil
}

func (c *CrService) DescribeCrEeNamespace(id string) (map[string]interface{}, error) {
	strRet := c.ParseResourceId(id)
	instanceId := strRet[0]
	namespaceName := strRet[1]
	

	request := c.client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "GetNamespace", "")
	request.QueryParams["InstanceId"] = instanceId
	request.QueryParams["NamespaceName"] = namespaceName

	bresponse, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil,  errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil,  errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		if response["errorMessage"].(string) == "Namespace is not exist."{
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil, fmt.Errorf("read ee namespace failed, %s", response["errorMessage"].(string))
	}
	return response, nil
}

func (c *CrService) ListCrEeRepos(instanceId string, namespace string, pageNo int, pageSize int) (*cr_ee.ListRepositoryResponse, error) {
	response := &cr_ee.ListRepositoryResponse{}
	request := cr_ee.CreateListRepositoryRequest()
	c.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.RepoNamespaceName = namespace
	request.RepoStatus = "ALL"
	request.PageNo = requests.NewInteger(pageNo)
	request.PageSize = requests.NewInteger(pageSize)
	resource := c.GenResourceId(instanceId, namespace)
	action := request.GetActionName()

	raw, err := c.client.WithCrEeClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListRepository(request)
	})
	response, ok := raw.(*cr_ee.ListRepositoryResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, resource, action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListRepositoryResponse)
	if !response.ListRepositoryIsSuccess {
		return response, errmsgs.WrapErrorf(errors.New(response.Code), errmsgs.DataDefaultErrorMsg, resource, action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DescribeCrEeRepo(id string) (map[string]interface{}, error) {
	strRet := c.ParseResourceId(id)
	instanceId := strRet[0]
	namespace := strRet[1]
	repoName := strRet[2]
	
	request := c.client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "GetRepository", "")
	mergeMaps(request.QueryParams, map[string]string{
		"InstanceId":        instanceId,
		"RepoNamespaceName": namespace,
		"RepoName":          repoName,
	})
	
	bresponse, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		if response["errorMessage"].(string) == "Repo is not exist."{
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil, fmt.Errorf("read ee repo failed, %s", response["errorMessage"].(string))
	}

	return response, nil

}

func (c *CrService) ListCrEeRepoTags(instanceId string, repoId string, pageNo int, pageSize int) (*cr_ee.ListRepoTagResponse, error) {
	response := &cr_ee.ListRepoTagResponse{}
	request := cr_ee.CreateListRepoTagRequest()
	c.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.RepoId = repoId
	request.PageNo = requests.NewInteger(pageNo)
	request.PageSize = requests.NewInteger(pageSize)
	resource := c.GenResourceId(instanceId, repoId)
	action := request.GetActionName()

	raw, err := c.client.WithCrEeClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListRepoTag(request)
	})
	response, ok := raw.(*cr_ee.ListRepoTagResponse)
	if err != nil {
		errmsg := ""
		if ok  {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, resource, action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(action, raw, request.RpcRequest, request)

	if !response.ListRepoTagIsSuccess {
		return response, errmsgs.WrapErrorf(errors.New(response.Code), errmsgs.DataDefaultErrorMsg, resource, action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DescribeCrEeSyncRule(id string) (*cr_ee.SyncRulesItem, error) {
	strRet := c.ParseResourceId(id)
	instanceId := strRet[0]
	namespace := strRet[1]
	syncRuleId := strRet[2]

	pageNo := 1
	for {
		response := &cr_ee.ListRepoSyncRuleResponse{}
		request := cr_ee.CreateListRepoSyncRuleRequest()
		c.client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = instanceId
		request.NamespaceName = namespace
		request.PageNo = requests.NewInteger(pageNo)
		request.PageSize = requests.NewInteger(PageSizeLarge)
		raw, err := c.client.WithCrEeClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.ListRepoSyncRule(request)
		})
		response, ok := raw.(*cr_ee.ListRepoSyncRuleResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if !response.ListRepoSyncRuleIsSuccess {
			return nil, c.wrapCrServiceError(id, request.GetActionName(), response.Code)
		}

		for _, rule := range response.SyncRules {
			if rule.SyncRuleId == syncRuleId && rule.LocalInstanceId == instanceId {
				return &rule, nil
			}
		}

		if len(response.SyncRules) < PageSizeLarge {
			return nil, errmsgs.WrapErrorf(errors.New("sync rule not found"), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		pageNo++
	}
}

func (c *CrService) wrapCrServiceError(resource string, action string, code string) error {
	switch code {
	case "INSTANCE_NOT_EXIST", "NAMESPACE_NOT_EXIST", "REPO_NOT_EXIST":
		return errmsgs.WrapErrorf(errors.New(code), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	default:
		return errmsgs.WrapErrorf(errors.New(code), errmsgs.DefaultErrorMsg, resource, action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
}

func (c *CrService) GenResourceId(args ...string) string {
	return strings.Join(args, COLON_SEPARATED)
}

func (c *CrService) ParseResourceId(id string) []string {
	return strings.Split(id, COLON_SEPARATED)
}
