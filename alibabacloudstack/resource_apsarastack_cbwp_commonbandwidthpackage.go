package alibabacloudstack

import (
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackCommonBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackCommonBandwidthPackageCreate,
		Read:   resourceAlibabacloudStackCommonBandwidthPackageRead,
		Update: resourceAlibabacloudStackCommonBandwidthPackageUpdate,
		Delete: resourceAlibabacloudStackCommonBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'bandwidth_package_name' instead.",
				ConflictsWith: []string{"bandwidth_package_name"},
			},
			"bandwidth_package_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PayByTraffic,
				ValidateFunc: validation.StringInSlice([]string{"PayByBandwidth", "PayBy95", "PayByTraffic"}, false),
			},
			"ratio": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      100,
				ValidateFunc: validation.IntBetween(10, 100),
			},
			//"resource_group_id": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	ForceNew: true,
			//	Computed: true,
			//},
		},
	}
}

func resourceAlibabacloudStackCommonBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateCommonBandwidthPackageRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))
	request.Name = connectivity.GetResourceData(d, "bandwidth_package_name", "name").(string)
	request.Description = d.Get("description").(string)
	request.InternetChargeType = d.Get("internet_charge_type").(string)
	request.Ratio = requests.NewInteger(d.Get("ratio").(int))

	wait := incrementalWait(1*time.Second, 1*time.Second)
	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		request.ClientToken = buildClientToken(request.GetActionName())
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateCommonBandwidthPackage(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", errmsgs.Throttling}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if bresponse, ok := raw.(*vpc.CreateCommonBandwidthPackageResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_common_bandwidth_package", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.CreateCommonBandwidthPackageResponse)
		d.SetId(response.BandwidthPackageId)
		return nil
	})
	if err != nil {
		return err
	}

	if err = vpcService.WaitForCommonBandwidthPackage(d.Id(), Available, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceAlibabacloudStackCommonBandwidthPackageRead(d, meta)
}

func resourceAlibabacloudStackCommonBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeCommonBandwidthPackage(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	bandwidth, err := strconv.Atoi(object.Bandwidth)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("bandwidth", bandwidth)
	connectivity.SetResourceData(d, object.Name, "bandwidth_package_name", "name")
	d.Set("description", object.Description)
	d.Set("internet_charge_type", object.InternetChargeType)
	d.Set("ratio", object.Ratio)
	return nil
}

func resourceAlibabacloudStackCommonBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	d.Partial(true)
	update := false
	request := vpc.CreateModifyCommonBandwidthPackageAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.BandwidthPackageId = d.Id()
	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}

	if d.HasChanges("name", "bandwidth_package_name") {
		request.Name = connectivity.GetResourceData(d, "bandwidth_package_name", "name").(string)
		update = true
	}
	log.Printf("111111111111111111111111111111")
	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyCommonBandwidthPackageAttribute(request)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*vpc.ModifyCommonBandwidthPackageAttributeResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	if d.HasChange("bandwidth") {
		request := vpc.CreateModifyCommonBandwidthPackageSpecRequest()
		client.InitRpcRequest(*request.RpcRequest)

		request.BandwidthPackageId = d.Id()
		request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyCommonBandwidthPackageSpec(request)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*vpc.ModifyCommonBandwidthPackageSpecResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return resourceAlibabacloudStackCommonBandwidthPackageRead(d, meta)
}

func resourceAlibabacloudStackCommonBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateDeleteCommonBandwidthPackageRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.BandwidthPackageId = d.Id()
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteCommonBandwidthPackage(request)
	})
	if err != nil {
		errmsg := ""
		if bresponse, ok := raw.(*vpc.DeleteCommonBandwidthPackageResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return errmsgs.WrapError(vpcService.WaitForCommonBandwidthPackage(d.Id(), Deleted, DefaultTimeoutMedium))
}
