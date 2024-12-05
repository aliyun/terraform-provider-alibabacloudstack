package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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

	request := client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "CreateRepository", "")
	mergeMaps(request.QueryParams, map[string]string{
		"InstanceId":        instanceId,
		"RepoNamespaceName": namespace,
		"RepoName":          repoName,
		"RepoType":          repoType,
		"Summary":           summary,
	})
	if detail, ok := d.GetOk("detail"); ok && detail.(string) != "" {
		request.QueryParams["Detail"] = detail.(string)
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
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

	request := client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "ListRepository", "")
	mergeMaps(request.QueryParams, map[string]string{
		"InstanceId":        instanceId,
		"RepoNamespaceName": namespace,
		"RepoName":          repo,
	})

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("read ee repo failed, %s", response["asapiErrorMessage"].(string))
	}
	repoList := response["Repositories"].([]interface{})
	if len(repoList) == 0 {
		return errmsgs.WrapError(fmt.Errorf("repo %s not found", repoList))
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

		request := client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "UpdateRepository", "")
		mergeMaps(request.QueryParams, map[string]string{
			"InstanceId": instanceId,
			"RepoId":     d.Get("repo_id").(string),
			"RepoType":   d.Get("repo_type").(string),
			"Summary":    d.Get("summary").(string),
		})
		if d.HasChange("detail") {
			request.QueryParams["Detail"] = d.Get("detail").(string)
		}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		response := make(map[string]interface{})
		addDebug(request.GetActionName(), raw, request, request.QueryParams)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
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

	request := client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "DeleteRepository", "")
	request.QueryParams["InstanceId"] = instanceId
	request.QueryParams["RepoId"] = repoId

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), raw, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("delete ee repo failed, %s", response["asapiErrorMessage"].(string))
	}

	return nil
}
