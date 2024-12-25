package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackForwardEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackForwardEntryCreate,
		Read:   resourceAlibabacloudStackForwardEntryRead,
		Update: resourceAlibabacloudStackForwardEntryUpdate,
		Delete: resourceAlibabacloudStackForwardEntryDelete,

		Schema: map[string]*schema.Schema{
			"forward_table_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateForwardPort,
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "any"}, false),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'forward_entry_name' instead.",
				ConflictsWith: []string{"forward_entry_name"},
			},
			"forward_entry_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateForwardPort,
			},
			"forward_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackForwardEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateForwardEntryRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ForwardTableId = d.Get("forward_table_id").(string)
	request.ExternalIp = d.Get("external_ip").(string)
	request.ExternalPort = d.Get("external_port").(string)
	request.IpProtocol = d.Get("ip_protocol").(string)
	request.InternalIp = d.Get("internal_ip").(string)
	request.InternalPort = d.Get("internal_port").(string)
	if name, ok := connectivity.GetResourceDataOk(d, "forward_entry_name", "name"); ok {
		request.ForwardEntryName = name.(string)
	}
	var raw interface{}
	var err error
	if err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		ar := request
		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateForwardEntry(ar)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidIp.NotInNatgw"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*vpc.CreateForwardEntryResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_forward_entry", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateForwardEntryResponse)

	d.SetId(request.ForwardTableId + COLON_SEPARATED + response.ForwardEntryId)
	if err := vpcService.WaitForForwardEntry(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	return resourceAlibabacloudStackForwardEntryRead(d, meta)
}

func resourceAlibabacloudStackForwardEntryRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}
	forwardEntry, err := vpcService.DescribeForwardEntry(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("forward_table_id", forwardEntry.ForwardTableId)
	d.Set("external_ip", forwardEntry.ExternalIp)
	d.Set("external_port", forwardEntry.ExternalPort)
	d.Set("ip_protocol", forwardEntry.IpProtocol)
	d.Set("internal_ip", forwardEntry.InternalIp)
	d.Set("internal_port", forwardEntry.InternalPort)
	d.Set("forward_entry_id", forwardEntry.ForwardEntryId)
	connectivity.SetResourceData(d, forwardEntry.ForwardEntryName, "forward_entry_name", "name")

	return nil
}

func resourceAlibabacloudStackForwardEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := vpc.CreateModifyForwardEntryRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ForwardEntryId = parts[1]
	request.ForwardTableId = parts[0]

	if d.HasChange("external_ip") {
		request.ExternalIp = d.Get("external_ip").(string)
	}

	if d.HasChange("external_port") {
		request.ExternalPort = d.Get("external_port").(string)
	}

	if d.HasChange("ip_protocol") {
		request.IpProtocol = d.Get("ip_protocol").(string)
	}

	if d.HasChange("internal_ip") {
		request.InternalIp = d.Get("internal_ip").(string)
	}

	if d.HasChange("internal_port") {
		request.InternalPort = d.Get("internal_port").(string)
	}
	if d.HasChange("name") {
		request.ForwardEntryName = connectivity.GetResourceData(d, "forward_entry_name", "name").(string)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyForwardEntry(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*vpc.ModifyForwardEntryResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err := vpcService.WaitForForwardEntry(d.Id(), Available, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	return resourceAlibabacloudStackForwardEntryRead(d, meta)
}

func resourceAlibabacloudStackForwardEntryDelete(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteForwardEntryRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ForwardTableId = parts[0]
	request.ForwardEntryId = parts[1]

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteForwardEntry(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"IncorretForwardEntryStatus"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.DeleteForwardEntryResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidForwardEntryId.NotFound", "InvalidForwardTableId.NotFound"}) {
			return nil
		}
		return err
	}
	return errmsgs.WrapError(vpcService.WaitForForwardEntry(d.Id(), Deleted, DefaultTimeout))
}
