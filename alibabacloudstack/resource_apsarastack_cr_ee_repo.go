package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCrEeRepo() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource, resourceAlibabacloudStackCrEeRepoCreate, resourceAlibabacloudStackCrEeRepoRead, resourceAlibabacloudStackCrEeRepoUpdate, resourceAlibabacloudStackCrEeRepoDelete)
	return resource
}

func resourceAlibabacloudStackCrEeRepoCreate(d *schema.ResourceData, meta interface{}) error {
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

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("create ee repo failed, %s", response["asapiErrorMessage"].(string))
	}

	d.SetId(crService.GenResourceId(instanceId, namespace, repoName))

	return nil
}

func resourceAlibabacloudStackCrEeRepoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	
	response, err := crService.DescribeCrEeRepo(d.Id())

	if err != nil {
		return errmsgs.WrapError(err)
	}


	d.Set("instance_id", response["InstanceId"].(string))
	d.Set("namespace", response["RepoNamespaceName"].(string))
	d.Set("name", response["RepoName"].(string))
	d.Set("repo_type", response["RepoType"])
	d.Set("summary", response["Summary"].(string))
	d.Set("repo_id", response["RepoId"].(string))


	return nil
}

func resourceAlibabacloudStackCrEeRepoUpdate(d *schema.ResourceData, meta interface{}) error {
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

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		response := make(map[string]interface{})
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if !response["asapiSuccess"].(bool) {
			return fmt.Errorf("update ee repo failed, %s", response["asapiErrorMessage"].(string))
		}

	}

	return nil
}

func resourceAlibabacloudStackCrEeRepoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	instanceId := d.Get("instance_id").(string)
	repoId := d.Get("repo_id").(string)

	request := client.NewCommonRequest("POST", "cr-ee", "2018-12-01", "DeleteRepository", "")
	request.QueryParams["InstanceId"] = instanceId
	request.QueryParams["RepoId"] = repoId

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := make(map[string]interface{})
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !response["asapiSuccess"].(bool) {
		return fmt.Errorf("delete ee repo failed, %s", response["asapiErrorMessage"].(string))
	}

	return nil
}