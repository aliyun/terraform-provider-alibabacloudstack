package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApsaraStackCrEERepo() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackCrEERepoCreate,
		Read:   resourceApsaraStackCrEERepoRead,
		Update: resourceApsaraStackCrEERepoUpdate,
		Delete: resourceApsaraStackCrEERepoDelete,
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

func resourceApsaraStackCrEERepoCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("namespace").(string)
	repoName := d.Get("name").(string)
	repoType := d.Get("repo_type").(string)
	summary := d.Get("summary").(string)

	response := &cr_ee.CreateRepositoryResponse{}
	request := cr_ee.CreateCreateRepositoryRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cr", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.RegionId = crService.client.RegionId
	request.InstanceId = instanceId
	request.RepoNamespaceName = namespace
	request.RepoName = repoName
	request.RepoType = repoType
	request.Summary = summary
	if detail, ok := d.GetOk("detail"); ok && detail.(string) != "" {
		request.Detail = detail.(string)
	}
	resource := crService.GenResourceId(instanceId, namespace, repoName)
	action := request.GetActionName()

	raw, err := crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.CreateRepository(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, resource, action, ApsaraStackSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.CreateRepositoryResponse)
	if !response.CreateRepositoryIsSuccess {
		return crService.wrapCrServiceError(resource, action, response.Code)
	}

	d.SetId(crService.GenResourceId(instanceId, namespace, repoName))

	return resourceApsaraStackCrEERepoRead(d, meta)
}

func resourceApsaraStackCrEERepoRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	crService := &CrService{client}
	resp, err := crService.DescribeCrEERepo(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", resp.InstanceId)
	d.Set("namespace", resp.RepoNamespaceName)
	d.Set("name", resp.RepoName)
	d.Set("repo_type", resp.RepoType)
	d.Set("summary", resp.Summary)
	d.Set("detail", resp.Detail)
	d.Set("repo_id", resp.RepoId)

	return nil
}

func resourceApsaraStackCrEERepoUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("namespace").(string)
	repo := d.Get("name").(string)
	if d.HasChanges("repo_type", "summary", "detail") {
		response := &cr_ee.UpdateRepositoryResponse{}
		request := cr_ee.CreateUpdateRepositoryRequest()
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cr", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		request.RegionId = crService.client.RegionId
		request.InstanceId = instanceId
		request.RepoId = d.Get("repo_id").(string)
		request.RepoType = d.Get("repo_type").(string)
		request.Summary = d.Get("summary").(string)
		if d.HasChange("detail") {
			request.Detail = d.Get("detail").(string)
		}
		resource := crService.GenResourceId(instanceId, namespace, repo)
		action := request.GetActionName()

		raw, err := crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.UpdateRepository(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, resource, action, ApsaraStackSdkGoERROR)
		}
		addDebug(action, raw, request.RpcRequest, request)

		response, _ = raw.(*cr_ee.UpdateRepositoryResponse)
		if !response.UpdateRepositoryIsSuccess {
			return crService.wrapCrServiceError(resource, action, response.Code)
		}
	}

	return resourceApsaraStackCrEERepoRead(d, meta)
}

func resourceApsaraStackCrEERepoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("namespace").(string)
	repo := d.Get("name").(string)
	repoId := d.Get("repo_id").(string)
	_, err := crService.DeleteCrEERepo(instanceId, namespace, repo, repoId)
	if err != nil {
		if NotFoundError(err) {
			return nil
		} else {
			return WrapError(err)
		}
	}

	return WrapError(crService.WaitForCrEERepo(instanceId, namespace, repo, Deleted, DefaultTimeout))
}
