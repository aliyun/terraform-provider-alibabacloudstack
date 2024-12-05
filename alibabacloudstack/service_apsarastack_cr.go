package alibabacloudstack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"log"
	"strings"
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
	//response := &cr.GetNamespaceResponse{}

	var err error
	response := crDescribeNamespaceResponse{}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "cr"
	request.Domain = c.client.Domain
	request.Version = "2016-06-07"
	request.Scheme = "http"
	request.ApiName = "GetNamespace"
	request.Headers = map[string]string{"RegionId": c.client.RegionId}
	request.QueryParams = map[string]string{
		
		
		"Product":         "cr",
		"Department":      c.client.Department,
		"ResourceGroup":   c.client.ResourceGroup,
		"RegionId":        c.client.RegionId,
		"Action":          "GetNamespace",
		"Version":         "2016-06-07",
		"Namespace":       id,
	}
	raw, err := c.client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	resp, _ := raw.(*responses.CommonResponse)
	log.Printf("response for read %v", resp)
	err = json.Unmarshal(resp.GetHttpContentBytes(), &response)
	log.Printf("unmarshal response for read %v", &response)

	if response.Data.Namespace.Namespace != id {
		return nil, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request)

	return &response, nil
}

func (c *CrService) DescribeCrRepo(id string) (GetRepoResponse, error) {
	resp := GetRepoResponse{}
	sli := strings.Split(id, SLASH_SEPARATED)
	repoNamespace := sli[0]
	repoName := sli[1]
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "cr"
	request.Domain = c.client.Domain
	request.Version = "2016-06-07"
	request.Scheme = "http"
	request.ApiName = "GetRepo"
	request.Headers = map[string]string{"RegionId": c.client.RegionId}
	request.QueryParams = map[string]string{
		
		
		"Product":         "cr",
		"Department":      c.client.Department,
		"ResourceGroup":   c.client.ResourceGroup,
		"RegionId":        c.client.RegionId,
		"Action":          "GetRepo",
		"Version":         "2016-06-07",
		"RepoName":        repoName,
		"RepoNamespace":   repoNamespace,
	}
	raw, err := c.client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return resp, err
	}
	response, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
	//if err != nil {
	//	if IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
	//		return response, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
	//	}
	//	return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	//}

	if err != nil {
		if IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request)
	return resp, nil
}
