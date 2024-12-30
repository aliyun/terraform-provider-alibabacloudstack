package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

func resourceAlibabacloudStackCRRepo() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackCRRepoCreate,
		Read:   resourceAlibabacloudStackCRRepoRead,
		Update: resourceAlibabacloudStackCRRepoUpdate,
		Delete: resourceAlibabacloudStackCRRepoDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 30),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 30),
			},
			"summary": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"repo_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{RepoTypePublic, RepoTypePrivate}, false),
			},
			"detail": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 2000),
			},
			// computed
			"domain_list": {
				Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlibabacloudStackCRRepoCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	resp := ResponseCr{}
	repoNamespace := d.Get("namespace").(string)
	repoName := d.Get("name").(string)
	summary := d.Get("summary").(string)
	repoType := d.Get("repo_type").(string)
	detail := d.Get("detail").(string)

	request := client.NewCommonRequest("PUT", "cr", "2016-06-07", "CreateRepo", "/repos")
	body := map[string]interface{}{
		"repo": map[string]interface{}{
			"RepoName":     repoName,
			"RepoNamespace": repoNamespace,
			"repoType":          repoType,
			"summary":           summary,
			"detail": detail,
		},
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
	}
	request.SetContentType(requests.Json)
	request.SetContent(jsonData)
	response, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), response, request)
	if err != nil {
		if response == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_repo", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	log.Printf("repo create response %v", response)
	err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cr_repo", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	log.Printf("repo create unmarshalled response %v", &resp)
	d.SetId(fmt.Sprintf("%s%s%s", repoNamespace, SLASH_SEPARATED, repoName))

	return resourceAlibabacloudStackCRRepoRead(d, meta)
}

type ResponseCr struct {
	Code               string `json:"code"`
	Data               struct {
		Data struct {
			RepoID int `json:"repoId"`
		} `json:"data"`
	} `json:"data"`
	SuccessResponse bool `json:"successResponse"`
}

func resourceAlibabacloudStackCRRepoUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	resp := ResponseCr{}
	repoNamespace := d.Get("namespace").(string)
	repoName := d.Get("name").(string)
	summary := d.Get("summary").(string)
	repoType := d.Get("repo_type").(string)
	detail := d.Get("detail").(string)

	if d.HasChanges("summary", "detail", "repo_type") {
		request := client.NewCommonRequest("POST", "cr", "2016-06-07", "UpdateRepo", fmt.Sprintf("/repos/%s/%s",repoNamespace,repoName) )
		request.QueryParams["RepoNamespace"] = repoNamespace
		request.QueryParams["RepoName"] = repoName
		body := map[string]interface{}{
			"repo": map[string]interface{}{
				"repoType":          repoType,
				"summary":           summary,
				"detail": detail,
			},
		}
		jsonData, err := json.Marshal(body)
		if err != nil {
			return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
		}
		request.SetContentType(requests.Json)
		request.SetContent(jsonData)
		response, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), response, request)
		if err != nil {
			if response == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_repo", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		log.Printf("repo update response %v", response)
		err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cr_repo", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		log.Printf("repo update unmarshalled response %v", &resp)
	}
	return resourceAlibabacloudStackCRRepoRead(d, meta)
}

func resourceAlibabacloudStackCRRepoRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := CrService{client}

	object, err := crService.DescribeCrRepo(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("namespace", object.Data.Repo.RepoNamespace)
	d.Set("name", object.Data.Repo.RepoName)
	d.Set("detail", object.Data.Repo.Detail)
	d.Set("summary", object.Data.Repo.Summary)
	d.Set("repo_type", object.Data.Repo.RepoType)

	domainList := make([]map[string]string, 1)
	domains := make(map[string]string)
	domains["public"] = object.Data.Repo.RepoDomainList.Public
	domains["internal"] = object.Data.Repo.RepoDomainList.Internal
	domains["vpc"] = object.Data.Repo.RepoDomainList.Vpc
	domainList = append(domainList, domains)

	d.Set("domain_list", domainList)

	return nil
}

func resourceAlibabacloudStackCRRepoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	sli := strings.Split(d.Id(), SLASH_SEPARATED)
	repoNamespace := sli[0]
	repoName := sli[1]

	request := client.NewCommonRequest("DELETE", "cr", "2016-06-07", "DeleteRepo", fmt.Sprintf("/repos/%s/%s",repoNamespace,repoName))
	request.QueryParams["RepoNamespace"] = repoNamespace
	request.QueryParams["RepoName"] = repoName

	response, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), response, request)
	if err != nil {
		if response == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_repo", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	return nil
}
