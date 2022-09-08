package alibabacloudstack

import (
	"fmt"

	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSnatEntryCreate,
		Read:   resourceAlibabacloudStackSnatEntryRead,
		Update: resourceAlibabacloudStackSnatEntryUpdate,
		Delete: resourceAlibabacloudStackSnatEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"snat_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_vswitch_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: strings.Fields("source_cidr"),
			},
			"snat_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_cidr": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: strings.Fields("source_vswitch_id"),
			},
			"snat_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackSnatEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateSnatEntryRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.SnatTableId = d.Get("snat_table_id").(string)
	request.SourceVSwitchId = d.Get("source_vswitch_id").(string)
	request.SnatIp = d.Get("snat_ip").(string)
	if v, ok := d.GetOk("source_cidr"); ok {
		request.SourceCIDR = v.(string)
	}

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ar := request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateSnatEntry(ar)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"EIP_NOT_IN_GATEWAY", "OperationUnsupported.EipNatBWPCheck", "OperationUnsupported.EipInBinding"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.CreateSnatEntryResponse)
		d.SetId(fmt.Sprintf("%s%s%s", request.SnatTableId, COLON_SEPARATED, response.SnatEntryId))
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_snat_entry", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	if err := vpcService.WaitForSnatEntry(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlibabacloudStackSnatEntryRead(d, meta)
}

func resourceAlibabacloudStackSnatEntryRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}

	object, err := vpcService.DescribeSnatEntry(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("snat_table_id", object.SnatTableId)
	if _, ok := d.GetOk("source_cidr"); ok {
		d.Set("source_cidr", object.SourceCIDR)
	} else {
		d.Set("source_vswitch_id", object.SourceVSwitchId)
	}
	d.Set("snat_ip", object.SnatIp)
	d.Set("snat_entry_id", object.SnatEntryId)

	return nil
}

func resourceAlibabacloudStackSnatEntryUpdate(d *schema.ResourceData, meta interface{}) error {

	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := vpc.CreateModifySnatEntryRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.SnatTableId = parts[0]
	request.SnatEntryId = parts[1]
	update := false
	if d.HasChange("snat_ip") {
		update = true
		request.SnatIp = d.Get("snat_ip").(string)
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifySnatEntry(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err := vpcService.WaitForSnatEntry(d.Id(), Available, DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}
	return resourceAlibabacloudStackSnatEntryRead(d, meta)
}

func resourceAlibabacloudStackSnatEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := vpc.CreateDeleteSnatEntryRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.SnatTableId = parts[0]
	request.SnatEntryId = parts[1]
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteSnatEntry(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorretSnatEntryStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSnatTableId.NotFound", "InvalidSnatEntryId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	return WrapError(vpcService.WaitForSnatEntry(d.Id(), Deleted, DefaultTimeout))
}
