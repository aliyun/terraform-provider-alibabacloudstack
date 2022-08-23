package apsarastack

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackNatGatewayCreate,
		Read:   resourceApsaraStackNatGatewayRead,
		Update: resourceApsaraStackNatGatewayUpdate,
		Delete: resourceApsaraStackNatGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"specification": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Small", "Middle", "Large"}, false),
				Default:      "Small",
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
		},
	}
}

func resourceApsaraStackNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateNatGatewayRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = string(client.Region)
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.VpcId = string(d.Get("vpc_id").(string))
	request.Spec = string(d.Get("specification").(string))
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

	if v, ok := d.GetOk("name"); ok {
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
			if IsExpectedErrors(err, []string{"VswitchStatusError", "TaskConflict"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(args.GetActionName(), raw, args.RpcRequest, args)
		response, _ := raw.(*vpc.CreateNatGatewayResponse)
		d.SetId(response.NatGatewayId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_nat_gateway", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	if err := vpcService.WaitForNatGateway(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceApsaraStackNatGatewayRead(d, meta)
}

func resourceApsaraStackNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeNatGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("specification", object.Spec)
	d.Set("bandwidth_package_ids", strings.Join(object.BandwidthPackageIds.BandwidthPackageId, ","))
	d.Set("snat_table_ids", strings.Join(object.SnatTableIds.SnatTableId, ","))
	d.Set("forward_table_ids", strings.Join(object.ForwardTableIds.ForwardTableId, ","))
	d.Set("description", object.Description)
	d.Set("vpc_id", object.VpcId)
	bindWidthPackages, err := flattenBandWidthPackages(object.BandwidthPackageIds.BandwidthPackageId, meta, d)
	if err != nil {
		return WrapError(err)
	} else {
		d.Set("bandwidth_packages", bindWidthPackages)
	}

	return nil
}

func resourceApsaraStackNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	natGateway, err := vpcService.DescribeNatGateway(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Partial(true)
	attributeUpdate := false
	modifyNatGatewayAttributeRequest := vpc.CreateModifyNatGatewayAttributeRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		modifyNatGatewayAttributeRequest.Scheme = "https"
	} else {
		modifyNatGatewayAttributeRequest.Scheme = "http"
	}
	modifyNatGatewayAttributeRequest.RegionId = natGateway.RegionId
	modifyNatGatewayAttributeRequest.NatGatewayId = natGateway.NatGatewayId
	modifyNatGatewayAttributeRequest.Headers = map[string]string{"RegionId": client.RegionId}
	modifyNatGatewayAttributeRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	if d.HasChange("name") {
		//d.SetPartial("name")
		var name string
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		} else {
			return WrapError(Error("cann't change name to empty string"))
		}
		modifyNatGatewayAttributeRequest.Name = name

		attributeUpdate = true
	}

	if d.HasChange("description") {
		//d.SetPartial("description")
		var description string
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		} else {
			return WrapError(Error("can to change description to empty string"))
		}

		modifyNatGatewayAttributeRequest.Description = description

		attributeUpdate = true
	}

	if attributeUpdate {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyNatGatewayAttribute(modifyNatGatewayAttributeRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyNatGatewayAttributeRequest.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(modifyNatGatewayAttributeRequest.GetActionName(), raw, modifyNatGatewayAttributeRequest.RpcRequest, modifyNatGatewayAttributeRequest)
	}

	if d.HasChange("specification") {
		//d.SetPartial("specification")
		modifyNatGatewaySpecRequest := vpc.CreateModifyNatGatewaySpecRequest()
		modifyNatGatewaySpecRequest.RegionId = natGateway.RegionId
		if strings.ToLower(client.Config.Protocol) == "https" {
			modifyNatGatewaySpecRequest.Scheme = "https"
		} else {
			modifyNatGatewaySpecRequest.Scheme = "http"
		}
		modifyNatGatewaySpecRequest.NatGatewayId = natGateway.NatGatewayId
		modifyNatGatewaySpecRequest.Headers = map[string]string{"RegionId": client.RegionId}
		modifyNatGatewaySpecRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		modifyNatGatewaySpecRequest.Spec = d.Get("specification").(string)

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyNatGatewaySpec(modifyNatGatewaySpecRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyNatGatewaySpecRequest.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(modifyNatGatewaySpecRequest.GetActionName(), raw, modifyNatGatewaySpecRequest.RpcRequest, modifyNatGatewaySpecRequest)
	}
	d.Partial(false)
	if err := vpcService.WaitForNatGateway(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceApsaraStackNatGatewayRead(d, meta)
}

func resourceApsaraStackNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	err := deleteBandwidthPackages(d, meta)
	if err != nil {
		return WrapError(err)
	}
	request := vpc.CreateDeleteNatGatewayRequest()
	request.RegionId = string(client.Region)
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.NatGatewayId = d.Id()
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := vpc.CreateDeleteNatGatewayRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.RegionId = string(client.Region)
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.NatGatewayId = d.Id()
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteNatGateway(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation.BandwidthPackages"}) {
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"InvalidNatGatewayId.NotFound"}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return WrapError(vpcService.WaitForNatGateway(d.Id(), Deleted, DefaultTimeoutMedium))
}

func deleteBandwidthPackages(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	packRequest := vpc.CreateDescribeBandwidthPackagesRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		packRequest.Scheme = "https"
	} else {
		packRequest.Scheme = "http"
	}
	packRequest.RegionId = string(client.Region)
	packRequest.Headers = map[string]string{"RegionId": client.RegionId}
	packRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	packRequest.NatGatewayId = d.Id()
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeBandwidthPackages(packRequest)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug(packRequest.GetActionName(), raw, packRequest.RpcRequest, packRequest)
		response, _ := raw.(*vpc.DescribeBandwidthPackagesResponse)
		retry := false
		if len(response.BandwidthPackages.BandwidthPackage) > 0 {
			for _, pack := range response.BandwidthPackages.BandwidthPackage {
				request := vpc.CreateDeleteBandwidthPackageRequest()
				request.RegionId = string(client.Region)
				if strings.ToLower(client.Config.Protocol) == "https" {
					request.Scheme = "https"
				} else {
					request.Scheme = "http"
				}
				request.BandwidthPackageId = pack.BandwidthPackageId
				request.Headers = map[string]string{"RegionId": client.RegionId}
				request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
				raw, e := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
					return vpcClient.DeleteBandwidthPackage(request)
				})
				if e != nil {
					if IsExpectedErrors(e, []string{"Invalid.RegionId"}) {
						return resource.NonRetryableError(e)
					} else if IsExpectedErrors(e, []string{"INSTANCE_NOT_EXISTS"}) {
						return nil
					}
					err = e
					retry = true
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		if retry {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), packRequest.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return nil
}

func flattenBandWidthPackages(bandWidthPackageIds []string, meta interface{}, d *schema.ResourceData) ([]map[string]interface{}, error) {
	packageLen := len(bandWidthPackageIds)
	result := make([]map[string]interface{}, 0, packageLen)
	for i := packageLen - 1; i >= 0; i-- {
		packageId := bandWidthPackageIds[i]
		bandWidthPackage, err := getPackage(packageId, meta, d)
		if err != nil {
			return result, WrapError(err)
		}
		ipAddress := flattenPackPublicIp(bandWidthPackage.PublicIpAddresses.PublicIpAddresse)
		ipCont, ipContErr := strconv.Atoi(bandWidthPackage.IpCount)
		bandWidth, bandWidthErr := strconv.Atoi(bandWidthPackage.Bandwidth)
		if ipContErr != nil {
			return result, WrapError(ipContErr)
		}
		if bandWidthErr != nil {
			return result, WrapError(bandWidthErr)
		}
		l := map[string]interface{}{
			"ip_count":            ipCont,
			"bandwidth":           bandWidth,
			"zone":                bandWidthPackage.ZoneId,
			"public_ip_addresses": ipAddress,
		}
		result = append(result, l)
	}
	return result, nil
}
func getPackage(packageId string, meta interface{}, d *schema.ResourceData) (pack vpc.BandwidthPackage, err error) {
	client := meta.(*connectivity.ApsaraStackClient)
	request := vpc.CreateDescribeBandwidthPackagesRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId

	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.NatGatewayId = d.Id()
	request.BandwidthPackageId = packageId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeBandwidthPackages(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		packages, _ := raw.(*vpc.DescribeBandwidthPackagesResponse)
		if len(packages.BandwidthPackages.BandwidthPackage) < 1 || packages.BandwidthPackages.BandwidthPackage[0].BandwidthPackageId != packageId {
			return WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		pack = packages.BandwidthPackages.BandwidthPackage[0]
		return nil
	})
	return
}
func flattenPackPublicIp(publicIpAddressList []vpc.PublicIpAddresse) string {
	var result []string
	for _, publicIpAddresses := range publicIpAddressList {
		ipAddress := publicIpAddresses.IpAddress
		result = append(result, ipAddress)
	}
	return strings.Join(result, ",")
}
