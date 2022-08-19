package apsarastack

import (
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackCommonBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackCommonBandwidthPackageCreate,
		Read:   resourceApsaraStackCommonBandwidthPackageRead,
		Update: resourceApsaraStackCommonBandwidthPackageUpdate,
		Delete: resourceApsaraStackCommonBandwidthPackageDelete,
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
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

func resourceApsaraStackCommonBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateCommonBandwidthPackageRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))
	request.Name = d.Get("name").(string)
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
			if IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", Throttling}) {
				wait()
				return resource.RetryableError(err)

			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.CreateCommonBandwidthPackageResponse)
		d.SetId(response.BandwidthPackageId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_common_bandwidth_package", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	if err = vpcService.WaitForCommonBandwidthPackage(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceApsaraStackCommonBandwidthPackageRead(d, meta)
}

func resourceApsaraStackCommonBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeCommonBandwidthPackage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	bandwidth, err := strconv.Atoi(object.Bandwidth)
	if err != nil {
		return WrapError(err)
	}
	d.Set("bandwidth", bandwidth)
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("internet_charge_type", object.InternetChargeType)
	d.Set("ratio", object.Ratio)
	return nil
}

func resourceApsaraStackCommonBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	d.Partial(true)
	update := false
	request := vpc.CreateModifyCommonBandwidthPackageAttributeRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.BandwidthPackageId = d.Id()
	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyCommonBandwidthPackageAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("description")
		//d.SetPartial("name")
	}

	if d.HasChange("bandwidth") {
		request := vpc.CreateModifyCommonBandwidthPackageSpecRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}

		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.BandwidthPackageId = d.Id()
		request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyCommonBandwidthPackageSpec(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("bandwidth")
	}

	d.Partial(false)
	return resourceApsaraStackCommonBandwidthPackageRead(d, meta)
}

func resourceApsaraStackCommonBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateDeleteCommonBandwidthPackageRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.BandwidthPackageId = d.Id()
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteCommonBandwidthPackage(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(vpcService.WaitForCommonBandwidthPackage(d.Id(), Deleted, DefaultTimeoutMedium))
}
