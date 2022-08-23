package apsarastack

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackSecurityGroupCreate,
		Read:   resourceApsaraStackSecurityGroupRead,
		Update: resourceApsaraStackSecurityGroupUpdate,
		Delete: resourceApsaraStackSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"inner_access_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Accept", "Drop"}, false),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceApsaraStackSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := ecs.CreateCreateSecurityGroupRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if v := d.Get("name").(string); v != "" {
		request.SecurityGroupName = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	if v := d.Get("vpc_id").(string); v != "" {
		request.VpcId = v
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSecurityGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_security_group", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.CreateSecurityGroupResponse)
	d.SetId(response.SecurityGroupId)
	return resourceApsaraStackSecurityGroupUpdate(d, meta)
}

func resourceApsaraStackSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeSecurityGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	d.Set("name", object.SecurityGroupName)
	d.Set("description", object.Description)
	d.Set("vpc_id", object.VpcId)
	d.Set("inner_access_policy", object.InnerAccessPolicy)

	request := ecs.CreateDescribeSecurityGroupsRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.SecurityGroupId = d.Id()

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroups(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeSecurityGroupsResponse)
	if len(response.SecurityGroups.SecurityGroup) < 1 {
		return WrapErrorf(Error(GetNotFoundMessage("SecurityGroup", d.Id())), NotFoundMsg, ProviderERROR)
	}
	d.Set("tags", ecsService.tagsToMap(response.SecurityGroups.SecurityGroup[0].Tags.Tag))

	return nil
}

func resourceApsaraStackSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	d.Partial(true)

	if err := setTags(client, TagResourceSecurityGroup, d); err != nil {
		return WrapError(err)
	} else {
		//d.SetPartial("tags")
	}

	if d.HasChange("inner_access_policy") {
		policy := GroupInnerAccept
		if v, ok := d.GetOk("inner_access_policy"); ok && v.(string) != "" {
			policy = GroupInnerAccessPolicy(v.(string))
		}
		request := ecs.CreateModifySecurityGroupPolicyRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.SecurityGroupId = d.Id()
		request.InnerAccessPolicy = string(policy)
		request.ClientToken = buildClientToken(request.GetActionName())

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupPolicy(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("inner_access_policy")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceApsaraStackSecurityGroupRead(d, meta)
	}

	update := false
	request := ecs.CreateModifySecurityGroupAttributeRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.SecurityGroupId = d.Id()
	if d.HasChange("name") {
		request.SecurityGroupName = d.Get("name").(string)
		update = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("name")
		//d.SetPartial("description")
	}

	d.Partial(false)

	return resourceApsaraStackSecurityGroupRead(d, meta)
}

func resourceApsaraStackSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	request := ecs.CreateDeleteSecurityGroupRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.SecurityGroupId = d.Id()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	err := resource.Retry(6*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSecurityGroup(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR)
	}
	return WrapError(ecsService.WaitForSecurityGroup(d.Id(), Deleted, DefaultTimeoutMedium))

}
