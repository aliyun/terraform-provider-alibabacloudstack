package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"errors"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
)

type CrService struct {
	client *connectivity.AlibabacloudStackClient
}

type crCreateNamespaceRequestPayload struct {
	Namespace struct {
		Namespace string `json:"Namespace"`
	} `json:"Namespace"`
}

type crUpdateNamespaceRequestPayload struct {
	Namespace struct {
		AutoCreate        bool   `json:"AutoCreate"`
		DefaultVisibility string `json:"DefaultVisibility"`
	} `json:"Namespace"`
}

type crListResponse struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data struct {
		Code       string `json:"code"`
		Cost       int    `json:"cost"`
		Message    string `json:"message"`
		Namespaces []struct {
			AuthorizeType     string `json:"authorizeType"`
			Department        int    `json:"Department"`
			NamespaceStatus   string `json:"namespaceStatus"`
			Namespace         string `json:"namespace"`
			DepartmentName    string `json:"DepartmentName"`
			ResourceGroup     int    `json:"ResourceGroup"`
			ResourceGroupName string `json:"ResourceGroupName"`
		} `json:"namespaces"`
		PureListData bool `json:"pureListData"`
		Redirect     bool `json:"redirect"`
		Success      bool `json:"success"`
	} `json:"data"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	Success      bool   `json:"success"`
}

type crDescribeNamespaceResponse struct {
	Code      string `json:"code"`
	RequestId string `json:"requestId"`
	Data      struct {
		Namespace struct {
			Namespace         string `json:"namespace"`
			AuthorizeType     string `json:"authorizeType"`
			DefaultVisibility string `json:"defaultVisibility"`
			AutoCreate        bool   `json:"autoCreate"`
			NamespaceStatus   string `json:"namespaceStatus"`
		} `json:"namespace"`
	} `json:"data"`
}

type crDescribeNamespaceListResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Namespace []struct {
			Namespace       string `json:"namespace"`
			AuthorizeType   string `json:"authorizeType"`
			NamespaceStatus string `json:"namespaceStatus"`
		} `json:"namespaces"`
	} `json:"data"`
}

const (
	RepoTypePublic  = "PUBLIC"
	RepoTypePrivate = "PRIVATE"
)

type crCreateRepoRequestPayload struct {
	Repo struct {
		RepoNamespace string `json:"RepoNamespace"`
		RepoName      string `json:"RepoName"`
		Summary       string `json:"Summary"`
		Detail        string `json:"Detail"`
		RepoType      string `json:"RepoType"`
	} `json:"Repo"`
}

type crUpdateRepoRequestPayload struct {
	Repo struct {
		Summary  string `json:"Summary"`
		Detail   string `json:"Detail"`
		RepoType string `json:"RepoType"`
	} `json:"Repo"`
}

type GetRepoResponse struct {
	Code string `json:"code"`
	Data struct {
		Repo struct {
			Stars          int    `json:"stars"`
			Logo           string `json:"logo"`
			RepoStatus     string `json:"repoStatus"`
			GmtCreate      int64  `json:"gmtCreate"`
			Detail         string `json:"detail"`
			GmtModified    int64  `json:"gmtModified"`
			Summary        string `json:"summary"`
			RepoBuildType  string `json:"repoBuildType"`
			RepoName       string `json:"repoName"`
			RepoNamespace  string `json:"repoNamespace"`
			RepoType       string `json:"repoType"`
			RepoID         int    `json:"repoId"`
			RegionID       string `json:"regionId"`
			RepoOriginType string `json:"repoOriginType"`
			RepoDomainList struct {
				Internal string `json:"internal"`
				Public   string `json:"public"`
				Vpc      string `json:"vpc"`
			} `json:"repoDomainList"`
			RepoAuthorizeType string `json:"repoAuthorizeType"`
			Downloads         int    `json:"downloads"`
		} `json:"repo"`
	} `json:"data"`
}

type crDescribeRepoResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Repo struct {
			Summary        string `json:"summary"`
			Detail         string `json:"detail"`
			RepoNamespace  string `json:"repoNamespace"`
			RepoName       string `json:"repoName"`
			RepoType       string `json:"repoType"`
			RepoDomainList struct {
				Public   string `json:"public"`
				Internal string `json:"internal"`
				Vpc      string `json:"vpc"`
			}
		} `json:"repo"`
	} `json:"data"`
}

type crDescribeReposResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Repos    []crRepo `json:"repos"`
		Total    int      `json:"total"`
		PageSize int      `json:"pageSize"`
		Page     int      `json:"page"`
	} `json:"data"`
}

type crResponseList struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data struct {
		Code         string `json:"code"`
		Cost         int    `json:"cost"`
		Message      string `json:"message"`
		Page         int    `json:"page"`
		PageSize     int    `json:"pageSize"`
		PureListData bool   `json:"pureListData"`
		Redirect     bool   `json:"redirect"`
		Repos        []struct {
			Summary        string `json:"summary"`
			RepoID         int    `json:"repoId"`
			GmtModified    int64  `json:"gmtModified"`
			RepoNamespace  string `json:"repoNamespace"`
			RepoName       string `json:"repoName"`
			RepoOriginType string `json:"repoOriginType"`
			Stars          int    `json:"stars"`
			GmtCreate      int64  `json:"gmtCreate"`
			RepoBuildType  string `json:"repoBuildType"`
			RepoType       string `json:"repoType"`
			RepoDomainList struct {
				Internal string `json:"internal"`
				Public   string `json:"public"`
				Vpc      string `json:"vpc"`
			} `json:"repoDomainList"`
			Downloads         int    `json:"downloads"`
			RegionID          string `json:"regionId"`
			Logo              string `json:"logo"`
			RepoStatus        string `json:"repoStatus"`
			RepoAuthorizeType string `json:"repoAuthorizeType"`
		} `json:"repos"`
		Success bool `json:"success"`
		Total   int  `json:"total"`
	} `json:"data"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	Success      bool   `json:"success"`
}

type crRepo struct {
	Summary        string `json:"summary"`
	RepoNamespace  string `json:"repoNamespace"`
	RepoName       string `json:"repoName"`
	RepoType       string `json:"repoType"`
	RegionId       string `json:"regionId"`
	RepoDomainList struct {
		Public   string `json:"public"`
		Internal string `json:"internal"`
		Vpc      string `json:"vpc"`
	} `json:"repoDomainList"`
}

type crDescribeRepoTagsResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Tags     []crTag `json:"tags"`
		Total    int     `json:"total"`
		PageSize int     `json:"pageSize"`
		Page     int     `json:"page"`
	} `json:"data"`
}

type crTag struct {
	ImageId     string `json:"imageId"`
	Digest      string `json:"digest"`
	Tag         string `json:"tag"`
	Status      string `json:"status"`
	ImageUpdate int    `json:"imageUpdate"`
	ImageCreate int    `json:"imageCreate"`
	ImageSize   int    `json:"imageSize"`
}

type crResponse struct {
	Code string `json:"code"`
	Data struct {
		Data struct {
			NamespaceID int `json:"namespaceId"`
		} `json:"data"`
	} `json:"data"`
	SuccessResponse bool `json:"successResponse"`
}

func (c *CrService) DescribeCrNamespace(id string) (*crDescribeNamespaceResponse, error) {
	response := crDescribeNamespaceResponse{}
	request := c.client.NewCommonRequest("GET", "cr", "2016-06-07", "GetNamespace", "/namespace/"+id)
	request.QueryParams["Namespace"] = id
	resp, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		errmsg := ""
		if resp == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	log.Printf("response for read %v", resp)
	err = json.Unmarshal(resp.GetHttpContentBytes(), &response)
	log.Printf("unmarshal response for read %v", &response)

	if response.Data.Namespace.Namespace != id {
		return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), resp, request)

	return &response, nil
}

func (c *CrService) DescribeCrRepo(id string) (GetRepoResponse, error) {
	resp := GetRepoResponse{}
	sli := strings.Split(id, SLASH_SEPARATED)
	repoNamespace := sli[0]
	repoName := sli[1]
	request := c.client.NewCommonRequest("GET", "cr", "2016-06-07", "GetRepo", fmt.Sprintf("/repos/%s/%s", repoNamespace, repoName))
	request.QueryParams["RepoName"] = repoName
	request.QueryParams["RepoNamespace"] = repoNamespace
	response, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		if response == nil {
			return resp, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), response, request)
	return resp, nil
}

func (c *CrService) ListCrEeInstances(pageNo int, pageSize int) (map[string]interface{}, error) {
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
	if code, ok := response["Code"].(string); ok && code != "success" {
		return nil, fmt.Errorf("read ee repo failed, %s", response)
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
	request.QueryParams["PageSize"] = strconv.Itoa(pageSize)

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
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	success, ok := response["asapiSuccess"].(bool)
	if !ok || !success {
		errMsg, ok := response["errorMessage"].(string)
		if !ok {
			return nil, fmt.Errorf("read ee namespace failed, unknown error")
		}
		errMesg := strings.ToLower(errMsg)
		if strings.Contains(errMesg, "namespace is not exist") || strings.Contains(errMesg, "namespace does not exist") {
			return nil, errmsgs.GetNotFoundErrorFromString(response["errorMessage"].(string))
		}
		return nil, fmt.Errorf("read ee namespace failed, %s", response["errorMessage"].(string))
	}
	return response, nil
}

func (c *CrService) ListCrEeRepos(instanceId string, namespace string, pageNo int, pageSize int) (map[string]interface{}, error) {
	request := c.client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "ListRepository", "")
	request.QueryParams["InstanceId"] = instanceId
	request.QueryParams["RepoNamespaceName"] = namespace
	request.QueryParams["RepoStatus"] = "ALL"
	request.QueryParams["PageNo"] = strconv.Itoa(pageNo)
	request.QueryParams["PageSize"] = strconv.Itoa(pageSize)
	resource := c.GenResourceId(instanceId, namespace)

	bresponse, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, resource, "ListRepository", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return nil, fmt.Errorf("read ee repo failed, %s", response["errorMessage"].(string))
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
		if response["errorMessage"].(string) == "Repo is not exist." {
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
		if ok {
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
func (c *CrService) DescribeCrEEArtifactLifecycleRule(id string) (map[string]interface{}, error) {
	strRet := c.ParseResourceId(id)
	instanceId := strRet[0]
	ruleId := strRet[1]
	reqQuery := map[string]interface{}{
		"PageNo":          1,
		"PageSize":        30,
		"InstanceId":      instanceId,
		"RuleId":          ruleId,
		"EnableDeleteTag": true,
	}

	if response, err := c.client.DoTeaRequest("GET", "cr-ee", "2018-12-01", "ListArtifactLifecycleRule", "", nil, reqQuery, nil); err != nil {
		return nil, err
	} else {
		if rules, ok := response["Rules"].([]interface{}); ok {
			for _, rule := range rules {
				ruleMap := rule.(map[string]interface{})
				if rule_id, ok := ruleMap["RuleId"].(string); ok && rule_id == ruleId {
					return ruleMap, nil
				}
			}
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("CrArtifactLifecycleRule %s not found", id))
}
