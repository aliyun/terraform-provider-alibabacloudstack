package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackCspprivateHsms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCspprivateHsmsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vendor_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hsms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The code of the HSM vendor.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCspprivateHsmsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"ZoneId":      d.Get("zone_id"),
		"VendorCode":  d.Get("vendor_code"),
		"ProductCode": d.Get("product_code"),
	}

	response, err := client.DoTeaRequest("GET", "hsm-private", "2018-06-30", "GetHsmIds", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	hsmIds, ok := response["HsmIds"].([]interface{})
	if !ok {
		return errmsgs.WrapError(errmsgs.GetNotFoundErrorFromString("No hms ids found"))
	}

	var hsms []map[string]interface{}
	var ids []string
	for _, v := range hsmIds {
		hsmId := v.(string)
		hsms = append(hsms, map[string]interface{}{
			"id": hsmId,
		})
		ids = append(ids, hsmId)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("hsms", hsms); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
