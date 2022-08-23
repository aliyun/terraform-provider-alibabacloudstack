package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
	"time"
)

func resourceApsaraStackForwardEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackForwardEntryCreate,
		Read:   resourceApsaraStackForwardEntryRead,
		Update: resourceApsaraStackForwardEntryUpdate,
		Delete: resourceApsaraStackForwardEntryDelete,

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
				Type:     schema.TypeString,
				Optional: true,
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

func resourceApsaraStackForwardEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateForwardEntryRequest()
	request.RegionId = string(client.Region)
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.ForwardTableId = d.Get("forward_table_id").(string)
	request.ExternalIp = d.Get("external_ip").(string)
	request.ExternalPort = d.Get("external_port").(string)
	request.IpProtocol = d.Get("ip_protocol").(string)
	request.InternalIp = d.Get("internal_ip").(string)
	request.InternalPort = d.Get("internal_port").(string)
	if name, ok := d.GetOk("name"); ok {
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
			if IsExpectedErrors(err, []string{"InvalidIp.NotInNatgw"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_forward_entry", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateForwardEntryResponse)

	d.SetId(request.ForwardTableId + COLON_SEPARATED + response.ForwardEntryId)
	if err := vpcService.WaitForForwardEntry(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return resourceApsaraStackForwardEntryRead(d, meta)
}

func resourceApsaraStackForwardEntryRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}
	forwardEntry, err := vpcService.DescribeForwardEntry(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("forward_table_id", forwardEntry.ForwardTableId)
	d.Set("external_ip", forwardEntry.ExternalIp)
	d.Set("external_port", forwardEntry.ExternalPort)
	d.Set("ip_protocol", forwardEntry.IpProtocol)
	d.Set("internal_ip", forwardEntry.InternalIp)
	d.Set("internal_port", forwardEntry.InternalPort)
	d.Set("forward_entry_id", forwardEntry.ForwardEntryId)
	d.Set("name", forwardEntry.ForwardEntryName)

	return nil
}

func resourceApsaraStackForwardEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := vpc.CreateModifyForwardEntryRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = string(client.Region)
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
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
		request.ForwardEntryName = d.Get("name").(string)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyForwardEntry(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err := vpcService.WaitForForwardEntry(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceApsaraStackForwardEntryRead(d, meta)
}

func resourceApsaraStackForwardEntryDelete(d *schema.ResourceData, meta interface{}) error {
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteForwardEntryRequest()
	request.RegionId = string(client.Region)
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.ForwardTableId = parts[0]
	request.ForwardEntryId = parts[1]

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteForwardEntry(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"UnknownError"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidForwardEntryId.NotFound", "InvalidForwardTableId.NotFound"}) {
			return nil
		}
		WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return WrapError(vpcService.WaitForForwardEntry(d.Id(), Deleted, DefaultTimeout))
}
