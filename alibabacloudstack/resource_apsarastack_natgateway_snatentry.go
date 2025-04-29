package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSnatEntry() *schema.Resource {
	resource := &schema.Resource{
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
				Computed:      true,
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
				Computed:      true,
				ConflictsWith: strings.Fields("source_vswitch_id"),
			},
			"snat_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, 
		resourceAlibabacloudStackSnatEntryCreate,
		resourceAlibabacloudStackSnatEntryRead,
		resourceAlibabacloudStackSnatEntryUpdate,
		resourceAlibabacloudStackSnatEntryDelete)
	return resource
}

func resourceAlibabacloudStackSnatEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateSnatEntryRequest()
	client.InitRpcRequest(*request.RpcRequest)
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
		bresponse, ok := raw.(*vpc.CreateSnatEntryResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"EIP_NOT_IN_GATEWAY", "OperationUnsupported.EipNatBWPCheck", "OperationUnsupported.EipInBinding"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_snat_entry", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetId(fmt.Sprintf("%s%s%s", request.SnatTableId, COLON_SEPARATED, bresponse.SnatEntryId))
		return nil
	}); err != nil {
		return err
	}

	if err := vpcService.WaitForSnatEntry(d.Id(), Available, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackSnatEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}

	object, err := vpcService.DescribeSnatEntry(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("snat_table_id", object.SnatTableId)
	d.Set("source_cidr", object.SourceCIDR)
	d.Set("source_vswitch_id", object.SourceVSwitchId)
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
		return errmsgs.WrapError(err)
	}

	request := vpc.CreateModifySnatEntryRequest()
	client.InitRpcRequest(*request.RpcRequest)
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
		bresponse, ok := raw.(*vpc.ModifySnatEntryResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err := vpcService.WaitForSnatEntry(d.Id(), Available, DefaultTimeout); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}

func resourceAlibabacloudStackSnatEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := vpc.CreateDeleteSnatEntryRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.SnatTableId = parts[0]
	request.SnatEntryId = parts[1]
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteSnatEntry(request)
		})
		bresponse, ok := raw.(*vpc.DeleteSnatEntryResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"IncorretSnatEntryStatus"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidSnatTableId.NotFound", "InvalidSnatEntryId.NotFound"}) {
			return nil
		}
		return err
	}
	return errmsgs.WrapError(vpcService.WaitForSnatEntry(d.Id(), Deleted, DefaultTimeout))
}
