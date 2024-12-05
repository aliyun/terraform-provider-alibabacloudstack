package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCrEESyncRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackCrEESyncRuleCreate,
		Read:   resourceAlibabacloudStackCrEESyncRuleRead,
		Delete: resourceAlibabacloudStackCrEESyncRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 30),
			},
			"repo_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 64),
			},
			"target_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 30),
			},
			"target_repo_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 64),
			},
			"tag_filter": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"sync_direction": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackCrEESyncRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	syncRuleName := d.Get("name").(string)
	instanceId := d.Get("instance_id").(string)
	namespaceName := d.Get("namespace_name").(string)
	targetRegionId := d.Get("target_region_id").(string)
	targetInstanceId := d.Get("target_instance_id").(string)
	targetNamespaceName := d.Get("target_namespace_name").(string)
	tagFilter := d.Get("tag_filter").(string)

	var repoName, targetRepoName string
	if v, ok := d.GetOk("repo_name"); ok {
		repoName = v.(string)
	}
	if v, ok := d.GetOk("target_repo_name"); ok {
		targetRepoName = v.(string)
	}
	if (repoName != "" && targetRepoName == "") || (repoName == "" && targetRepoName != "") {
		return WrapError(Error(DefaultErrorMsg, syncRuleName, "create", "[Params repo_name or target_repo_name is empty]"))
	}

	response := &cr_ee.CreateRepoSyncRuleResponse{}
	request := cr_ee.CreateCreateRepoSyncRuleRequest()
	request.RegionId = crService.client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "cr", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.SyncRuleName = syncRuleName
	request.InstanceId = instanceId
	request.NamespaceName = namespaceName
	request.TargetRegionId = targetRegionId
	request.TargetInstanceId = targetInstanceId
	request.TargetNamespaceName = targetNamespaceName
	request.TagFilter = tagFilter
	request.SyncTrigger = "PASSIVE"
	if repoName != "" && targetRepoName != "" {
		request.SyncScope = "REPO"
		request.RepoName = repoName
		request.TargetRepoName = targetRepoName
	} else {
		request.SyncScope = "NAMESPACE"
	}

	raw, err := crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.CreateRepoSyncRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cr_ee_sync_rule", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.CreateRepoSyncRuleResponse)
	if !response.CreateRepoSyncRuleIsSuccess {
		return crService.wrapCrServiceError("alibabacloudstack_cr_ee_sync_rule", request.GetActionName(), response.Code)
	}

	d.SetId(crService.GenResourceId(instanceId, namespaceName, response.SyncRuleId))

	return resourceAlibabacloudStackCrEESyncRuleRead(d, meta)
}

func resourceAlibabacloudStackCrEESyncRuleRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	resp, err := crService.DescribeCrEESyncRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", resp.SyncRuleName)
	d.Set("rule_id", resp.SyncRuleId)
	d.Set("instance_id", resp.LocalInstanceId)
	d.Set("namespace_name", resp.LocalNamespaceName)
	d.Set("repo_name", resp.LocalRepoName)
	d.Set("target_region_id", resp.TargetRegionId)
	d.Set("target_instance_id", resp.TargetInstanceId)
	d.Set("target_namespace_name", resp.TargetNamespaceName)
	d.Set("target_repo_name", resp.TargetRepoName)
	d.Set("tag_filter", resp.TagFilter)
	d.Set("sync_direction", resp.SyncDirection)
	d.Set("sync_scope", resp.SyncScope)

	return nil
}

func resourceAlibabacloudStackCrEESyncRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	syncRuleId := d.Get("rule_id").(string)
	syncDirection := d.Get("sync_direction").(string)
	if syncDirection != "FROM" {
		return WrapError(Error(DefaultErrorMsg, d.Id(), "delete", "[Please delete sync rule in the source instance]"))
	}

	response := &cr_ee.DeleteRepoSyncRuleResponse{}
	request := cr_ee.CreateDeleteRepoSyncRuleRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "cr", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.RegionId = crService.client.RegionId
	request.InstanceId = instanceId
	request.SyncRuleId = syncRuleId
	raw, err := crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.DeleteRepoSyncRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.DeleteRepoSyncRuleResponse)
	if !response.DeleteRepoSyncRuleIsSuccess {
		return crService.wrapCrServiceError(d.Id(), request.GetActionName(), response.Code)
	}

	return nil
}
