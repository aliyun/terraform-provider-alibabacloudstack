package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackEssScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEssScalingGroupCreate,
		Read:   resourceAlibabacloudStackEssScalingGroupRead,
		Update: resourceAlibabacloudStackEssScalingGroupUpdate,
		Delete: resourceAlibabacloudStackEssScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"scaling_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 40),
			},
			"multi_az_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_cooldown": {
				Type:         schema.TypeInt,
				Default:      300,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 86400),
			},
			"vswitch_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"removal_policies": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				MaxItems: 2,
				MinItems: 1,
			},
			"db_instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MinItems: 0,
			},
			"loadbalancer_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MinItems: 0,
			},
		},
	}
}

func resourceAlibabacloudStackEssScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	request, err := buildAlibabacloudStackEssScalingGroupArgs(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	client.InitRpcRequest(*request.RpcRequest)
	essService := EssService{client}
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateScalingGroup(request)
		})
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ess.CreateScalingGroupResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling, "IncorrectLoadBalancerHealthCheck", "IncorrectLoadBalancerStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_scalinggroup", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		response, _ := raw.(*ess.CreateScalingGroupResponse)
		d.SetId(response.ScalingGroupId)
		return nil
	}); err != nil {
		return err
	}
	if err := essService.WaitForEssScalingGroup(d.Id(), Inactive, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceAlibabacloudStackEssScalingGroupUpdate(d, meta)
}

func resourceAlibabacloudStackEssScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	object, err := essService.DescribeEssScalingGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("min_size", object.MinSize)
	d.Set("max_size", object.MaxSize)
	d.Set("scaling_group_name", object.ScalingGroupName)
	d.Set("default_cooldown", object.DefaultCooldown)
	d.Set("multi_az_policy", object.MultiAZPolicy)
	var polices []string
	if len(object.RemovalPolicies.RemovalPolicy) > 0 {
		for _, v := range object.RemovalPolicies.RemovalPolicy {
			polices = append(polices, v)
		}
	}
	d.Set("removal_policies", polices)
	var dbIds []string
	if len(object.DBInstanceIds.DBInstanceId) > 0 {
		for _, v := range object.DBInstanceIds.DBInstanceId {
			dbIds = append(dbIds, v)
		}
	}
	d.Set("db_instance_ids", dbIds)

	var slbIds []string
	if len(object.LoadBalancerIds.LoadBalancerId) > 0 {
		for _, v := range object.LoadBalancerIds.LoadBalancerId {
			slbIds = append(slbIds, v)
		}
	}
	d.Set("loadbalancer_ids", slbIds)

	var vswitchIds []string
	if len(object.VSwitchIds.VSwitchId) > 0 {
		for _, v := range object.VSwitchIds.VSwitchId {
			vswitchIds = append(vswitchIds, v)
		}
	}
	d.Set("vswitch_ids", vswitchIds)

	return nil
}

func resourceAlibabacloudStackEssScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateModifyScalingGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ScalingGroupId = d.Id()

	d.Partial(true)
	if d.HasChange("scaling_group_name") {
		request.ScalingGroupName = d.Get("scaling_group_name").(string)
	}
	if d.HasChange("multi_az_policy") {
		request.MultiAZPolicy = d.Get("multi_az_policy").(string)
	}
	if d.HasChange("min_size") {
		request.MinSize = requests.NewInteger(d.Get("min_size").(int))
	}

	if d.HasChange("max_size") {
		request.MaxSize = requests.NewInteger(d.Get("max_size").(int))
	}
	if d.HasChange("default_cooldown") {
		request.DefaultCooldown = requests.NewInteger(d.Get("default_cooldown").(int))
	}

	if d.HasChange("vswitch_ids") {
		vSwitchIds := expandStringList(d.Get("vswitch_ids").(*schema.Set).List())
		request.VSwitchIds = &vSwitchIds
	}

	if d.HasChange("removal_policies") {
		policyies := expandStringList(d.Get("removal_policies").([]interface{}))
		s := reflect.ValueOf(request).Elem()
		for i, p := range policyies {
			s.FieldByName(fmt.Sprintf("RemovalPolicy%d", i+1)).Set(reflect.ValueOf(p))
		}
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingGroup(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*ess.ModifyScalingGroupResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.Partial(false)
	return resourceAlibabacloudStackEssScalingGroupRead(d, meta)
}

func resourceAlibabacloudStackEssScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	request := ess.CreateDeleteScalingGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ScalingGroupId = d.Id()
	request.ForceDelete = requests.NewBoolean(true)

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingGroup(request)
	})

	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*ess.DeleteScalingGroupResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		if errmsgs.IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return errmsgs.WrapError(essService.WaitForEssScalingGroup(d.Id(), Deleted, DefaultTimeout))
}

func buildAlibabacloudStackEssScalingGroupArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingGroupRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := SlbService{client}

	request := ess.CreateCreateScalingGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.MinSize = requests.NewInteger(d.Get("min_size").(int))
	request.MaxSize = requests.NewInteger(d.Get("max_size").(int))
	request.DefaultCooldown = requests.NewInteger(d.Get("default_cooldown").(int))

	if v, ok := d.GetOk("scaling_group_name"); ok && v.(string) != "" {
		request.ScalingGroupName = v.(string)
	}
	if v, ok := d.GetOk("multi_az_policy"); ok && v.(string) != "" {
		request.MultiAZPolicy = v.(string)
	}
	if v, ok := d.GetOk("vswitch_ids"); ok {
		ids := expandStringList(v.(*schema.Set).List())
		request.VSwitchIds = &ids
	}

	if dbs, ok := d.GetOk("db_instance_ids"); ok {
		request.DBInstanceIds = convertListToJsonString(dbs.(*schema.Set).List())
	}

	if lbs, ok := d.GetOk("loadbalancer_ids"); ok {
		for _, lb := range lbs.(*schema.Set).List() {
			if err := slbService.WaitForSlb(lb.(string), Active, DefaultTimeout); err != nil {
				return nil, errmsgs.WrapError(err)
			}
		}
		request.LoadBalancerIds = convertListToJsonString(lbs.(*schema.Set).List())
	}

	return request, nil
}
