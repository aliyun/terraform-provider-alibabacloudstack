package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackNatGateway() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:true,
				ValidateFunc: validation.StringInSlice([]string{"Small", "Middle", "Large"}, false),
				Deprecated:   "Field 'specification' is deprecated and will be removed in a future release. Please use new field 'spec' instead.",
				ConflictsWith: []string{"spec"},
			},
			"spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:true,
				ValidateFunc: validation.StringInSlice([]string{"Small", "Middle", "Large"}, false),
				ConflictsWith: []string{"specification"},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Deprecated: "Field 'name' is deprecated and will be removed in a future release. " +
					"Please use new field 'nat_gateway_name' instead.",
				ConflictsWith: []string{"nat_gateway_name"},
			},
			"nat_gateway_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"bandwidth_package_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"snat_table_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"forward_table_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"bandwidth_packages": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"public_ip_addresses": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				MaxItems: 4,
				Optional: true,
			},
			"tags": tagsSchema(),
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNatGatewayCreate, resourceAlibabacloudStackNatGatewayRead, resourceAlibabacloudStackNatGatewayUpdate, resourceAlibabacloudStackNatGatewayDelete)
	return resource
}

func resourceAlibabacloudStackNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateNatGatewayRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.VpcId = string(d.Get("vpc_id").(string))
	request.Spec = connectivity.GetResourceData(d, "spec", "specification").(string)
	if request.Spec == "" {
		// Default must be nil if computed
		request.Spec = "Small"
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	bandwidthPackages := []vpc.CreateNatGatewayBandwidthPackage{}
	for _, e := range d.Get("bandwidth_packages").([]interface{}) {
		pack := e.(map[string]interface{})
		bandwidthPackage := vpc.CreateNatGatewayBandwidthPackage{
			IpCount:   strconv.Itoa(pack["ip_count"].(int)),
			Bandwidth: strconv.Itoa(pack["bandwidth"].(int)),
		}
		if pack["zone"].(string) != "" {
			bandwidthPackage.Zone = pack["zone"].(string)
		}
		bandwidthPackages = append(bandwidthPackages, bandwidthPackage)
	}

	request.BandwidthPackage = &bandwidthPackages

	if v, ok := connectivity.GetResourceDataOk(d, "nat_gateway_name", "name"); ok {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateNatGateway(&args)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"VswitchStatusError", "TaskConflict"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if bresponse, ok := raw.(*vpc.CreateNatGatewayResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_nat_gateway", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(args.GetActionName(), raw, args.RpcRequest, args)
		response, _ := raw.(*vpc.CreateNatGatewayResponse)
		d.SetId(response.NatGatewayId)
		return nil
	}); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_nat_gateway", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if err := vpcService.WaitForNatGateway(d.Id(), Available, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeNatGateway(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.Name, "nat_gateway_name", "name")
	connectivity.SetResourceData(d, object.Spec, "spec", "specification")

	d.Set("bandwidth_package_ids", strings.Join(object.BandwidthPackageIds.BandwidthPackageId, ","))
	d.Set("snat_table_ids", strings.Join(object.SnatTableIds.SnatTableId, ","))
	d.Set("forward_table_ids", strings.Join(object.ForwardTableIds.ForwardTableId, ","))
	d.Set("description", object.Description)
	d.Set("vpc_id", object.VpcId)
	// bindWidthPackages, err := flattenBandWidthPackages(object.BandwidthPackageIds.BandwidthPackageId, meta, d)
	// if err != nil {
	// 	return errmsgs.WrapError(err)
	// } else {
	// 	d.Set("bandwidth_packages", bindWidthPackages)
	// }
	listTagResourcesObject, err := vpcService.ListTagResources(d.Id(), "NATGATEWAY")
	if err == nil {
		d.Set("tags", tagsToMap(listTagResourcesObject))
	}

	return nil
}

func resourceAlibabacloudStackNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "NATGATEWAY"); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	natGateway, err := vpcService.DescribeNatGateway(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.Partial(true)
	attributeUpdate := false
	modifyNatGatewayAttributeRequest := vpc.CreateModifyNatGatewayAttributeRequest()
	client.InitRpcRequest(*modifyNatGatewayAttributeRequest.RpcRequest)
	modifyNatGatewayAttributeRequest.NatGatewayId = natGateway.NatGatewayId

	if d.HasChanges("name", "nat_gateway_name") {
		var name string
		if v, ok := connectivity.GetResourceDataOk(d, "nat_gateway_name", "name"); ok {
			name = v.(string)
		} else {
			return errmsgs.WrapError(errmsgs.Error("cann't change name to empty string"))
		}
		modifyNatGatewayAttributeRequest.Name = name

		attributeUpdate = true
	}

	if d.HasChange("description") {
		var description string
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		} else {
			return errmsgs.WrapError(errmsgs.Error("can to change description to empty string"))
		}

		modifyNatGatewayAttributeRequest.Description = description

		attributeUpdate = true
	}

	if attributeUpdate {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyNatGatewayAttribute(modifyNatGatewayAttributeRequest)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*vpc.ModifyNatGatewayAttributeResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), modifyNatGatewayAttributeRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(modifyNatGatewayAttributeRequest.GetActionName(), raw, modifyNatGatewayAttributeRequest.RpcRequest, modifyNatGatewayAttributeRequest)
	}

	if d.HasChanges("specification", "spec") {
		modifyNatGatewaySpecRequest := vpc.CreateModifyNatGatewaySpecRequest()
		client.InitRpcRequest(*modifyNatGatewaySpecRequest.RpcRequest)
		modifyNatGatewaySpecRequest.NatGatewayId = natGateway.NatGatewayId
		modifyNatGatewaySpecRequest.Spec = connectivity.GetResourceData(d, "spec", "specification").(string)
		if modifyNatGatewaySpecRequest.Spec == "" {
			modifyNatGatewaySpecRequest.Spec = "Small"
		}

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyNatGatewaySpec(modifyNatGatewaySpecRequest)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*vpc.ModifyNatGatewaySpecResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), modifyNatGatewaySpecRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(modifyNatGatewaySpecRequest.GetActionName(), raw, modifyNatGatewaySpecRequest.RpcRequest, modifyNatGatewaySpecRequest)
	}
	d.Partial(false)
	if err := vpcService.WaitForNatGateway(d.Id(), Available, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateDeleteNatGatewayRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.NatGatewayId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteNatGateway(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"DependencyViolation.BandwidthPackages"}) {
				return resource.RetryableError(err)
			}
			if errmsgs.IsExpectedErrors(err, []string{"InvalidNatGatewayId.NotFound"}) {
				return nil
			}
			errmsg := ""
			if bresponse, ok := raw.(*vpc.DeleteNatGatewayResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return errmsgs.WrapError(vpcService.WaitForNatGateway(d.Id(), Deleted, DefaultTimeoutMedium))
}