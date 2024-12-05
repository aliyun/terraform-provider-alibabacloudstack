package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:     schema.TypeList,
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
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "cr"
	request.Domain = client.Domain
	request.Version = "2016-06-07"
	request.Scheme = "http"
	request.ApiName = "CreateRepo"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		
		
		"Product":         "cr",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"RegionId":        client.RegionId,
		"Action":          "CreateRepo",
		"Version":         "2016-06-07",
		"X-acs-body":      fmt.Sprintf("{\"%s\":{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}}", "repo", "RepoName", repoName, "RepoNamespace", repoNamespace, "repoType", repoType, "summary", summary, "detail", detail),
	}

	raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cr_repo", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	response := raw.(*responses.CommonResponse)
	log.Printf("repo create response %v", response)
	err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cr_repo", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	log.Printf("repo create unmarshalled response %v", &resp)
	addDebug(request.GetActionName(), raw, request)
	d.SetId(fmt.Sprintf("%s%s%s", repoNamespace, SLASH_SEPARATED, repoName))

	return resourceAlibabacloudStackCRRepoRead(d, meta)
}

type ResponseCr struct {
	Code string `json:"code"`
	Data struct {
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
	if d.HasChange("summary") || d.HasChange("detail") || d.HasChange("repo_type") {
		request := requests.NewCommonRequest()
		request.Method = "POST"
		request.Product = "cr"
		request.Domain = client.Domain
		request.Version = "2016-06-07"
		request.Scheme = "http"
		request.ApiName = "UpdateRepo"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			
			
			"Product":         "cr",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"RegionId":        client.RegionId,
			"Action":          "UpdateRepo",
			"Version":         "2016-06-07",
			"RepoNamespace":   repoNamespace,
			"RepoName":        repoName,
			"X-acs-body":      fmt.Sprintf("{\"%s\":{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}}", "repo", "repoType", repoType, "summary", summary, "detail", detail),
		}
		raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
			return crClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cr_repo", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		response := raw.(*responses.CommonResponse)
		log.Printf("repo create response %v", response)
		err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cr_repo", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		log.Printf("repo create unmarshalled response %v", &resp)
		addDebug(request.GetActionName(), raw, request)
	}
	return resourceAlibabacloudStackCRRepoRead(d, meta)
}

func resourceAlibabacloudStackCRRepoRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := CrService{client}

	object, err := crService.DescribeCrRepo(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("namespace", object.Data.Repo.RepoNamespace)
	d.Set("name", object.Data.Repo.RepoName)
	d.Set("detail", object.Data.Repo.Detail)
	d.Set("summary", object.Data.Repo.Summary)
	d.Set("repo_type", object.Data.Repo.RepoType)

	domainList := make(map[string]string)
	domainList["public"] = object.Data.Repo.RepoDomainList.Public
	domainList["internal"] = object.Data.Repo.RepoDomainList.Internal
	domainList["vpc"] = object.Data.Repo.RepoDomainList.Vpc

	d.Set("domain_list", domainList)

	return nil
}

func resourceAlibabacloudStackCRRepoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	sli := strings.Split(d.Id(), SLASH_SEPARATED)
	repoNamespace := sli[0]
	repoName := sli[1]
	//
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "cr"
	request.Domain = client.Domain
	request.Version = "2016-06-07"
	request.Scheme = "http"
	request.ApiName = "DeleteRepo"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		
		
		"Product":         "cr",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"RegionId":        client.RegionId,
		"Action":          "DeleteRepo",
		"Version":         "2016-06-07",
		"RepoNamespace":   repoNamespace,
		"RepoName":        repoName,
	}
	raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cr_repo", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request)

	return nil
}
