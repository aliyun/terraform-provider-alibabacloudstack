package alibabacloudstack

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
				Optional: true,
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

	request := vpc.CreateAllocateEipAddressProRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.Bandwidth = strconv.Itoa(d.Get("bandwidth").(int))
	request.ClientToken = buildClientToken(request.GetActionName())
	if v, ok := d.GetOk("ip_address"); ok && v != "" {
		request.IpAddress = v.(string)
	}
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.AllocateEipAddressPro(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*vpc.AllocateEipAddressProResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_eip", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.AllocateEipAddressProResponse)
	d.SetId(response.AllocationId)
	err = vpcService.WaitForEip(d.Id(), Available, DefaultTimeoutMedium)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	return resourceAlibabacloudStackEipUpdate(d, meta)
}

func resourceAlibabacloudStackEipRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeEip(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
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
			return errmsgs.WrapError(err)
		}
	}

	update := false
	request := vpc.CreateModifyEipAddressAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
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
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.ModifyEipAddressAttributeResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlibabacloudStackEipRead(d, meta)
}

func resourceAlibabacloudStackEipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateReleaseEipAddressRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.AllocationId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ReleaseEipAddress(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"IncorrectEipStatus"}) {
				return resource.RetryableError(err)
			} else if errmsgs.IsExpectedErrors(err, []string{"InvalidAllocationId.NotFound"}) {
				return nil
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.ReleaseEipAddressResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}
	return errmsgs.WrapError(vpcService.WaitForEip(d.Id(), Deleted, DefaultTimeout))
}
