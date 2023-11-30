package alibabacloudstack

import (
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEipCreate,
		Read:   resourceAlibabacloudStackEipRead,
		Update: resourceAlibabacloudStackEipUpdate,
		Delete: resourceAlibabacloudStackEipDelete,
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
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackEipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateAllocateEipAddressRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = string(client.Region)
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.AllocateEipAddress(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_eip", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.AllocateEipAddressResponse)
	d.SetId(response.AllocationId)
	err = vpcService.WaitForEip(d.Id(), Available, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlibabacloudStackEipUpdate(d, meta)
}

func resourceAlibabacloudStackEipRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeEip(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	// 314 版本API字段名称 Descritpion
	//d.Set("description", object.Descritpion)
	// 316版本API字段名称 Description
	d.Set("description", object.Description)
	bandwidth, _ := strconv.Atoi(object.Bandwidth)
	d.Set("bandwidth", bandwidth)
	d.Set("ip_address", object.IpAddress)
	d.Set("status", object.Status)
	if tag := object.Tags.Tag; tag != nil {
		d.Set("tags", vpcService.tagToMap(tag))
	}
	return nil
}

func resourceAlibabacloudStackEipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "EIP"); err != nil {
			return WrapError(err)
		}
	}

	update := false
	request := vpc.CreateModifyEipAddressAttributeRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.AllocationId = d.Id()

	if d.HasChange("bandwidth") && !d.IsNewResource() {
		update = true
		request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
	}
	if d.HasChange("name") {
		update = true
		request.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		update = true
		request.Description = d.Get("description").(string)
	}
	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyEipAddressAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlibabacloudStackEipRead(d, meta)
}

func resourceAlibabacloudStackEipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateReleaseEipAddressRequest()
	request.AllocationId = d.Id()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ReleaseEipAddress(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectEipStatus"}) {
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"InvalidAllocationId.NotFound"}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	return WrapError(vpcService.WaitForEip(d.Id(), Deleted, DefaultTimeout))
}
