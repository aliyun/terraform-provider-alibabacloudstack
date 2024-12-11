package alibabacloudstack

import (
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSwitchCreate,
		Read:   resourceAlibabacloudStackSwitchRead,
		Update: resourceAlibabacloudStackSwitchUpdate,
		Delete: resourceAlibabacloudStackSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'availability_zone' is deprecated and will be removed in a future release. Please use 'zone_id' instead.",
				ConflictsWith: []string{"zone_id"},
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateSwitchCIDRNetworkAddress,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use 'vswitch_name' instead.",
				ConflictsWith: []string{"vswitch_name"},
			},
			"vswitch_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateVSwitchRequest()
	client.InitRpcRequest(*request.RpcRequest)

	log.Printf("[DEBUG] alibabacloud_vswitch ipv6CidrBlock: %s", d.Get("ipv6_cidr_block").(string))

	if d.Get("enable_ipv6").(bool) {
		request.Ipv6CidrBlock = "0"
	}

	if v, ok := connectivity.GetResourceDataOk(d, "zone_id", "availability_zone"); ok && v.(string) != "" {
		request.ZoneId = v.(string)
	}
	request.VpcId = Trim(d.Get("vpc_id").(string))
	request.CidrBlock = Trim(d.Get("cidr_block").(string))

	if v, ok := connectivity.GetResourceDataOk(d, "vswitch_name", "name"); ok && v.(string) != "" {
		request.VSwitchName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v != "" {
		request.Description = v.(string)
	}

	request.ClientToken = buildClientToken(request.GetActionName())

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVSwitch(&args)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"TaskConflict", "UnknownError", "InvalidStatus.RouteEntry", errmsgs.Throttling, "OperationFailed.IdempotentTokenProcessing"}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if bresponse, ok := raw.(*vpc.CreateVSwitchResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_vswitch", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		response, _ := raw.(*vpc.CreateVSwitchResponse)
		log.Printf("Vswitch response %s", response)
		d.SetId(response.VSwitchId)
		return nil
	}); err != nil {
		return err
	}

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, vpcService.VSwitchStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return resourceAlibabacloudStackSwitchUpdate(d, meta)
}

func resourceAlibabacloudStackSwitchRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	vswitch, err := vpcService.DescribeVSwitch(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, vswitch.ZoneId, "zone_id", "availability_zone")
	d.Set("vpc_id", vswitch.VpcId)
	d.Set("cidr_block", vswitch.CidrBlock)
	d.Set("ipv6_cidr_block", vswitch.Ipv6CidrBlock)
	connectivity.SetResourceData(d, vswitch.VSwitchName, "vswitch_name", "name")
	listTagResourcesObject, err := vpcService.ListTagResources(d.Id(), "VSWITCH")
	if err == nil {
		d.Set("tags", tagsToMap(listTagResourcesObject))
	}
	d.Set("description", vswitch.Description)
	return nil
}

func resourceAlibabacloudStackSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "VSWITCH"); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackSwitchRead(d, meta)
	}

	update := false
	request := vpc.CreateModifyVSwitchAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VSwitchId = d.Id()

	if d.HasChanges("vswitch_name", "name") {
		request.VSwitchName = connectivity.GetResourceData(d, "vswitch_name", "name").(string)
		update = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}

	if d.HasChange("enable_ipv6") && d.Get("enable_ipv6").(bool) {
		request.Ipv6CidrBlock = "0"
		update = true
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVSwitchAttribute(request)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*vpc.ModifyVSwitchAttributeResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlibabacloudStackSwitchRead(d, meta)
}

func resourceAlibabacloudStackSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateDeleteVSwitchRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VSwitchId = d.Id()

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVSwitch(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidRegionId.NotFound"}) {
				return resource.NonRetryableError(err)
			}
			if errmsgs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return nil
			}

			errmsg := ""
			if bresponse, ok := raw.(*vpc.DeleteVSwitchResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, vpcService.VSwitchStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}
