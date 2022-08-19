package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApsaraStackCrEENamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackCrEENamespaceCreate,
		Read:   resourceApsaraStackCrEENamespaceRead,
		Update: resourceApsaraStackCrEENamespaceUpdate,
		Delete: resourceApsaraStackCrEENamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 30),
			},
			"auto_create": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"default_visibility": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{RepoTypePublic, RepoTypePrivate}, false),
			},
		},
	}
}

func resourceApsaraStackCrEENamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("name").(string)
	autoCreate := d.Get("auto_create").(bool)
	visibility := d.Get("default_visibility").(string)

	response := &cr_ee.CreateNamespaceResponse{}
	request := cr_ee.CreateCreateNamespaceRequest()
	request.RegionId = crService.client.RegionId

	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cr", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.InstanceId = instanceId
	request.NamespaceName = namespace
	request.AutoCreateRepo = requests.NewBoolean(autoCreate)
	request.DefaultRepoType = visibility
	resource := crService.GenResourceId(instanceId, namespace)
	action := request.GetActionName()
	raw, err := crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.CreateNamespace(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, resource, action, ApsaraStackSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.CreateNamespaceResponse)
	if !response.CreateNamespaceIsSuccess {
		return crService.wrapCrServiceError(resource, action, response.Code)
	}

	d.SetId(crService.GenResourceId(instanceId, namespace))

	return resourceApsaraStackCrEENamespaceRead(d, meta)
}

func resourceApsaraStackCrEENamespaceRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	crService := &CrService{client}
	resp, err := crService.DescribeCrEENamespace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", resp.InstanceId)
	d.Set("name", resp.NamespaceName)
	d.Set("auto_create", resp.AutoCreateRepo)
	d.Set("default_visibility", resp.DefaultRepoType)

	return nil
}

func resourceApsaraStackCrEENamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("name").(string)
	if d.HasChanges("auto_create", "default_visibility") {
		autoCreate := d.Get("auto_create").(bool)
		visibility := d.Get("default_visibility").(string)
		response := &cr_ee.UpdateNamespaceResponse{}
		request := cr_ee.CreateUpdateNamespaceRequest()
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "cr", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		request.RegionId = crService.client.RegionId
		request.InstanceId = instanceId
		request.NamespaceName = namespace
		request.AutoCreateRepo = requests.NewBoolean(autoCreate)
		request.DefaultRepoType = visibility
		resource := crService.GenResourceId(instanceId, namespace)
		action := request.GetActionName()
		raw, err := crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.UpdateNamespace(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, resource, action, ApsaraStackSdkGoERROR)
		}
		addDebug(action, raw, request.RpcRequest, request)

		response, _ = raw.(*cr_ee.UpdateNamespaceResponse)
		if !response.UpdateNamespaceIsSuccess {
			return crService.wrapCrServiceError(resource, action, response.Code)
		}
	}

	return resourceApsaraStackCrEENamespaceRead(d, meta)
}

func resourceApsaraStackCrEENamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("name").(string)
	_, err := crService.DeleteCrEENamespace(instanceId, namespace)
	if err != nil {
		if NotFoundError(err) {
			return nil
		} else {
			return WrapError(err)
		}
	}

	return WrapError(crService.WaitForCrEENamespace(instanceId, namespace, Deleted, DefaultTimeout))
}
