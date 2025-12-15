package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackHsmVendors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackHsmVendorsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vendors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The code of the HSM vendor.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the HSM vendor.",
						},
						"products": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The code of the HSM product.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the HSM product.",
									},
								},
							},
							Description: "A list of HSM products.",
						},
					},
				},
				Description: "A list of HSM vendors.",
			},
		},
	}
}

func dataSourceAlibabacloudStackHsmVendorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})

	response, err := client.DoTeaRequest("GET", "hsm-private", "2018-06-30", "DescribeVendors", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	vendorsRaw, ok := response["Vendors"].([]interface{})
	if !ok {
		return errmsgs.WrapError(errmsgs.GetNotFoundErrorFromString("No hms vendors found"))
	}

	var vendors []map[string]interface{}
	var ids []string
	for _, v := range vendorsRaw {
		vendor := v.(map[string]interface{})
		vendorMapping := map[string]interface{}{
			"code": vendor["Code"],
			"name": vendor["Name"],
		}
		productsque := map[string]interface{}{
			"VendorCode": vendor["Code"],
		}
		response, err := client.DoTeaRequest("GET", "hsm-private", "2018-06-30", "DescribeProducts", "", nil, productsque, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		productsRaw, ok := response["Products"].([]interface{})
		if !ok {
			return errmsgs.WrapError(errmsgs.GetNotFoundErrorFromString("No hms vendor products found"))
		}
		products := []map[string]interface{}{}
		for _, p := range productsRaw {
			product := p.(map[string]interface{})
			productMapping := map[string]interface{}{
				"code": product["Code"],
				"name": product["Name"],
			}
			products = append(products, productMapping)
		}
		vendorMapping["products"] = products
		vendors = append(vendors, vendorMapping)
		ids = append(ids, vendorMapping["code"].(string))
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("vendors", vendors); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
