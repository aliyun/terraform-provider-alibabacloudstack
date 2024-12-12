package alibabacloudstack

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSecurityGroupCreate,
		Read:   resourceAlibabacloudStackSecurityGroupRead,
		Update: resourceAlibabacloudStackSecurityGroupUpdate,
		Delete: resourceAlibabacloudStackSecurityGroupDelete,
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

			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"enterprise", "normal"}, false),
				Default:      "normal",
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

func resourceAlibabacloudStackSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ecs.CreateCreateSecurityGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	if v := d.Get("name").(string); v != "" {
		request.SecurityGroupName = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	if v := d.Get("vpc_id").(string); v != "" {
		request.VpcId = v
	}
	request.SecurityGroupType = d.Get("type").(string)
	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSecurityGroup(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.CreateSecurityGroupResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_security_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.CreateSecurityGroupResponse)
	d.SetId(response.SecurityGroupId)
	return resourceAlibabacloudStackSecurityGroupUpdate(d, meta)
}

func resourceAlibabacloudStackSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeSecurityGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	d.Set("name", object.SecurityGroupName)
	d.Set("description", object.Description)
	d.Set("vpc_id", object.VpcId)
	d.Set("inner_access_policy", object.InnerAccessPolicy)

	request := ecs.CreateDescribeSecurityGroupsRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.SecurityGroupId = d.Id()

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroups(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.DescribeSecurityGroupsResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeSecurityGroupsResponse)
	if len(response.SecurityGroups.SecurityGroup) < 1 {
		return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SecurityGroup", d.Id())), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	d.Set("tags", ecsService.tagsToMap(response.SecurityGroups.SecurityGroup[0].Tags.Tag))
	d.Set("type",response.SecurityGroups.SecurityGroup[0].SecurityGroupType)

	return nil
}

func resourceAlibabacloudStackSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	d.Partial(true)

	if err := setTags(client, TagResourceSecurityGroup, d); err != nil {
		return errmsgs.WrapError(err)
	} else {
		//d.SetPartial("tags")
	}

	if d.HasChange("inner_access_policy") {
		policy := GroupInnerAccept
		if v, ok := d.GetOk("inner_access_policy"); ok && v.(string) != "" {
			policy = GroupInnerAccessPolicy(v.(string))
		}
		request := ecs.CreateModifySecurityGroupPolicyRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.SecurityGroupId = d.Id()
		request.InnerAccessPolicy = string(policy)
		request.ClientToken = buildClientToken(request.GetActionName())

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupPolicy(request)
		})
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*ecs.ModifySecurityGroupPolicyResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("inner_access_policy")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackSecurityGroupRead(d, meta)
	}

	update := false
	request := ecs.CreateModifySecurityGroupAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
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
			errmsg := ""
			if response, ok := raw.(*ecs.ModifySecurityGroupAttributeResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("name")
		//d.SetPartial("description")
	}

	d.Partial(false)

	return resourceAlibabacloudStackSecurityGroupRead(d, meta)
}

func resourceAlibabacloudStackSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	request := ecs.CreateDeleteSecurityGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.SecurityGroupId = d.Id()

	err := resource.Retry(6*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSecurityGroup(request)
		})

		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"DependencyViolation"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if response, ok := raw.(*ecs.DeleteSecurityGroupResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultTimeoutMsg, d.Id(), request.GetActionName(), errmsgs.ProviderERROR)
	}
	return errmsgs.WrapError(ecsService.WaitForSecurityGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
