package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCrEERepo() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackCrEERepoCreate,
		Read:   resourceAlibabacloudStackCrEERepoRead,
		Update: resourceAlibabacloudStackCrEERepoUpdate,
		Delete: resourceAlibabacloudStackCrEERepoDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
				ValidateFunc: validation.StringLenBetween(2, 64),
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

			//Computed
			"repo_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackCrEERepoCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("namespace").(string)
	repoName := d.Get("name").(string)
	repoType := d.Get("repo_type").(string)
	summary := d.Get("summary").(string)
	request := requests.NewCommonRequest()
	request.RegionId = client.RegionId
	request.Product = "cr-ee"
	request.Method = "POST"
	request.Domain = client.Domain
	request.Version = "2018-12-01"
	request.ApiName = "CreateRepository"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret":   client.SecretKey,
		"Product":           "cr-ee",
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"Action":            "CreateRepository",
		"Version":           "2018-12-01",
		"InstanceId":        instanceId,
		"RepoNamespaceName": namespace,
		"RepoName":          repoName,
		"RepoType":          repoType,
		"Summary":           summary,
	}
	if detail, ok := d.GetOk("detail"); ok && detail.(string) != "" {
		request.QueryParams["Detail"] = detail.(string)
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	err = json.Unmarshal(raw.(*responses.CommonResponse).GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("create ee repo failed, %s", response["asapiErrorMessage"].(string))
	}

	d.SetId(crService.GenResourceId(instanceId, namespace, repoName))

	return resourceAlibabacloudStackCrEERepoRead(d, meta)
}

func resourceAlibabacloudStackCrEERepoRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	strRet := crService.ParseResourceId(d.Id())
	instanceId := strRet[0]
	namespace := strRet[1]
	repo := strRet[2]
	request := requests.NewCommonRequest()
	request.RegionId = client.RegionId
	request.Product = "cr-ee"
	request.Method = "POST"
	request.Domain = client.Domain
	request.Version = "2018-12-01"
	request.ApiName = "ListRepository"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret":   client.SecretKey,
		"Product":           "cr-ee",
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"Action":            "ListRepository",
		"Version":           "2018-12-01",
		"InstanceId":        instanceId,
		"RepoNamespaceName": namespace,
		"RepoName":          repo,
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	err = json.Unmarshal(raw.(*responses.CommonResponse).GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("read ee repo failed, %s", response["asapiErrorMessage"].(string))
	}
	repoList := response["Repositories"].([]interface{})
	if len(repoList) == 0 {
		return WrapError(fmt.Errorf("repo %s not found", repoList))
	}
	item := repoList[0].(map[string]interface{})

	d.Set("instance_id", item["InstanceId"].(string))
	d.Set("namespace", item["RepoNamespaceName"].(string))
	d.Set("name", item["RepoName"].(string))
	d.Set("repo_type", item["RepoType"])
	d.Set("summary", item["Summary"].(string))
	d.Set("repo_id", item["RepoId"].(string))

	return nil
}

func resourceAlibabacloudStackCrEERepoUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	instanceId := d.Get("instance_id").(string)
	if d.HasChanges("repo_type", "summary", "detail") {
		request := requests.NewCommonRequest()
		request.RegionId = client.RegionId
		request.Product = "cr-ee"
		request.Method = "POST"
		request.Domain = client.Domain
		request.Version = "2018-12-01"
		request.ApiName = "UpdateRepository"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"Product":         "cr-ee",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Action":          "UpdateRepository",
			"Version":         "2018-12-01",
			"InstanceId":      instanceId,
			"RepoId":          d.Get("repo_id").(string),
			"RepoType":        d.Get("repo_type").(string),
			"Summary":         d.Get("summary").(string),
		}
		if d.HasChange("detail") {
			request.QueryParams["Detail"] = d.Get("detail").(string)
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		response := make(map[string]interface{})
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}

		err = json.Unmarshal(raw.(*responses.CommonResponse).GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if !response["asapiSuccess"].(bool) {
			return fmt.Errorf("update ee repo failed, %s", response["asapiErrorMessage"].(string))
		}

	}

	return resourceAlibabacloudStackCrEERepoRead(d, meta)
}

func resourceAlibabacloudStackCrEERepoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	instanceId := d.Get("instance_id").(string)
	repoId := d.Get("repo_id").(string)
	request := requests.NewCommonRequest()
	request.RegionId = client.RegionId
	request.Product = "cr-ee"
	request.Method = "POST"
	request.Domain = client.Domain
	request.Version = "2018-12-01"
	request.ApiName = "DeleteRepository"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"Product":         "cr-ee",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "DeleteRepository",
		"Version":         "2018-12-01",
		"InstanceId":      instanceId,
		"RepoId":          repoId,
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	err = json.Unmarshal(raw.(*responses.CommonResponse).GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("delete ee repo failed, %s", response["asapiErrorMessage"].(string))
	}

	return nil
}
