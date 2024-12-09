package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackReservedInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackReservedInstanceCreate,
		Read:   resourceAlibabacloudStackReservedInstanceRead,
		Update: resourceAlibabacloudStackReservedInstanceUpdate,
		Delete: resourceAlibabacloudStackReservedInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Region",
				ValidateFunc: validation.StringInSlice([]string{"Region", "Zone"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Windows", "Linux"}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Year",
				ValidateFunc: validation.StringInSlice([]string{"Year"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 3}),
			},
			"offering_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"No Upfront", "Partial Upfront", "All Upfront"}, false),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'reserved_instance_name' instead.",
				ConflictsWith: []string{"reserved_instance_name"},
			},
			"reserved_instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Deprecated:   "Field 'resource_group_id' is deprecated and will be removed in a future release. Please use new field 'reserved_instance_id' instead.",
				ConflictsWith: []string{"reserved_instance_id"},
			},
			"reserved_instance_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"resource_group_id"},
			},
		},
	}
}

func resourceAlibabacloudStackReservedInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	request := ecs.CreatePurchaseReservedInstancesOfferingRequest()
	client.InitRpcRequest(*request.RpcRequest)
	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = v.(string)
	}
	if scope, ok := d.GetOk("scope"); ok {
		request.Scope = scope.(string)
		if v, ok := d.GetOk("zone_id"); ok {
			if scope == "Zone" && v == "" {
				return errmsgs.WrapError(errmsgs.Error("Required when Scope is Zone."))
			}
			request.ZoneId = v.(string)
		}
	}
	if v, ok := d.GetOk("instance_amount"); ok {
		request.InstanceAmount = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("platform"); ok {
		request.Platform = v.(string)
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request.PeriodUnit = v.(string)
	}
	if v, ok := d.GetOk("period"); ok {
		request.Period = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("offering_type"); ok {
		request.OfferingType = v.(string)
	}
	if v, ok := connectivity.GetResourceDataOk(d, "reserved_instance_name", "name"); ok {
		request.ReservedInstanceName = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if v, ok := connectivity.GetResourceDataOk(d, "reserved_instance_id", "resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.PurchaseReservedInstancesOffering(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.PurchaseReservedInstancesOfferingResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_reserved_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.PurchaseReservedInstancesOfferingResponse)
	if len(response.ReservedInstanceIdSets.ReservedInstanceId) != 1 {
		return errmsgs.WrapError(errmsgs.Error("API returned wrong number of collections"))
	}
	d.SetId(response.ReservedInstanceIdSets.ReservedInstanceId[0])

	if err := ecsService.WaitForReservedInstance(d.Id(), Active, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceAlibabacloudStackReservedInstanceRead(d, meta)
}

func resourceAlibabacloudStackReservedInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ecs.CreateModifyReservedInstanceAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ReservedInstanceId = d.Id()
	if d.HasChanges("reserved_instance_name", "description", "name", "description") {
		if v, ok := connectivity.GetResourceDataOk(d, "reserved_instance_name", "name"); ok {
			request.ReservedInstanceName = v.(string)
		}
		if v, ok := d.GetOk("description"); ok {
			request.Description = v.(string)
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyReservedInstanceAttribute(request)
		})
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*ecs.ModifyReservedInstanceAttributeResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlibabacloudStackReservedInstanceRead(d, meta)
}

func resourceAlibabacloudStackReservedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	reservedInstances, err := ecsService.DescribeReservedInstance(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("instance_type", reservedInstances.InstanceType)
	d.Set("scope", reservedInstances.Scope)
	d.Set("zone_id", reservedInstances.ZoneId)
	d.Set("instance_amount", reservedInstances.InstanceAmount)
	d.Set("platform", reservedInstances.Platform)
	d.Set("offering_type", reservedInstances.OfferingType)
	connectivity.SetResourceData(d, reservedInstances.ReservedInstanceName, "reserved_instance_name", "name")
	d.Set("description", reservedInstances.Description)
	connectivity.SetResourceData(d, reservedInstances.ReservedInstanceId, "reserved_instance_id", "resource_group_id")

	return errmsgs.WrapError(err)
}

func resourceAlibabacloudStackReservedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	// PurchaseReservedInstancesOffering can not be release.
	return nil
}
