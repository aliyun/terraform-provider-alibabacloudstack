package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCspprivateHsmInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vendor_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vsm_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"evsm", "gvsm", "svsm"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alias_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCspprivateHsmInstanceCreate, resourceAlibabacloudStackCspprivateHsmInstanceRead, resourceAlibabacloudStackCspprivateHsmInstanceUpdate, resourceAlibabacloudStackCspprivateHsmInstanceDelete)
	return resource
}

func resourceAlibabacloudStackCspprivateHsmInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"ProductCode": d.Get("product_code"),
		"VendorCode":  d.Get("vendor_code"),
		"VsmType":     d.Get("vsm_type"),
		"HsmNum":      1,
		"ZoneId":      d.Get("zone_id"),
		"AliasName":   d.Get("alias_name"),
	}
	if v, ok := d.GetOk("device_id"); ok {
		request["DeviceId"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VswitchId"] = v.(string)
	}

	response, err := client.DoTeaRequest("POST", "Cspprivate", "2022-02-17", "CreateHsm", "", nil, request, nil)
	if err != nil {
		return err
	}

	instanceId, err := jsonpath.Get("$.HsmInfos.HsmInfo[0].HsmId", response)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.SetId(instanceId.(string))

	return nil
}

func resourceAlibabacloudStackCspprivateHsmInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cspprivateService := CspprivateService{client}
	instance, err := cspprivateService.DescribeCspprivateHsmInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if v, ok := instance["ProductCode"]; ok && v.(string) != "" {
		d.Set("product_code", v.(string))
	}
	d.Set("vendor_code", instance["VendorCode"])
	d.Set("vsm_type", instance["VsmType"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("device_id", instance["DeviceId"])
	d.Set("vpc_id", instance["VpcId"])
	d.Set("vswitch_id", instance["VswitchId"])
	d.Set("ip", instance["HsmIp"])
	d.Set("port", instance["HsmPort"])
	d.Set("alias_name", instance["AliasName"])
	d.Set("status", instance["HsmStatus"])
	return nil
}

func resourceAlibabacloudStackCspprivateHsmInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.HasChange("alias_name") {
		reqQuery := map[string]interface{}{
			"HsmId":      d.Id(),
			"InstanceId": d.Id(),
			"AliasName":  d.Get("alias_name"),
		}

		_, err := client.DoTeaRequest("POST", "Cspprivate", "2022-02-17", "UpdateHsmAliasName", "", nil, reqQuery, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_cspprivate_hsm_instance", "UpdateHsmAliasName", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	if d.HasChanges("vpc_id", "vswitch_id", "vpc_cidr_block", "ip") {
		reqQuery := map[string]interface{}{
			"HsmId":        d.Id(),
			"InstanceId":   d.Id(),
			"VpcId":        d.Get("vpc_id"),
			"VpcCidrBlock": d.Get("vpc_cidr_block"),
			"VswitchId":    d.Get("vswitch_id"),
			"Ip":           d.Get("ip"),
		}
		_, err := client.DoTeaRequest("POST", "Cspprivate", "2022-02-17", "ConfigHsmNetwork", "", nil, reqQuery, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_cspprivate_hsm_instance", "ConfigHsmNetwork", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}

func resourceAlibabacloudStackCspprivateHsmInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	reqQuery := map[string]interface{}{
		"HsmId": d.Id(),
	}
	_, err := client.DoTeaRequest("POST", "Cspprivate", "2022-02-17", "DeleteHsm", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
