package alibabacloudstack

import (
	"encoding/json"

	"github.com/PaesslerAG/jsonpath"
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
	querydata := map[string]interface{}{
		"HttpMethod":  "post",
		"ContentType": "application/x-www-form-urlencoded",
		"Uri":         "/api/hsm/gethsmids",
		"Data": map[string]interface{}{
			"vendorCode":  d.Get("vendor_code"),
			"zoneId":      d.Get("zone_id"),
			"productCode": d.Get("product_code"),
			"regionId":    client.RegionId,
		},
	}
	queryString, _ := json.Marshal(querydata)

	request := map[string]interface{}{
		"CryptoServiceId": "testid",
		"Data":            string(queryString),
	}

	response, err := client.DoTeaRequest("POST", "Cspprivate", "2022-02-17", "DescribeProxyCryptoService", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	hsmIds, err := jsonpath.Get("$.data.HsmIds", response)
	if err != nil {
		hsmIds, err = jsonpath.Get("$.data.hsmIds", response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}

	var hsms []map[string]interface{}
	var ids []string
	for _, v := range hsmIds.([]interface{}) {
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
