package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackCspprivateVendors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCspprivateVendorsRead,

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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
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

func dataSourceAlibabacloudStackCspprivateVendorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"CryptoServiceId": "testid",
		"Data":            "{\"ContentType\":\"application/x-www-form-urlencoded\",\"Uri\":\"/api/hsm/describevendors\",\"HttpMethod\":\"post\"}",
	}

	response, err := client.DoTeaRequest("POST", "Cspprivate", "2022-02-17", "DescribeProxyCryptoService", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	vendorsRaw, err := jsonpath.Get("$.data.vendors", response)
	if err != nil || len(vendorsRaw.([]interface{})) == 0 {
		return errmsgs.WrapError(err)
	}

	var vendors []map[string]interface{}
	var ids []string
	for _, v := range vendorsRaw.([]interface{}) {
		vendor := v.(map[string]interface{})
		vendorMapping := map[string]interface{}{
			"id":   vendor["code"],
			"code": vendor["code"],
			"name": vendor["name"],
		}
		productsque := map[string]interface{}{
			"CryptoServiceId": "testid",
			"Data":            fmt.Sprintf("{\"HttpMethod\":\"post\",\"ContentType\":\"application/x-www-form-urlencoded\",\"Uri\":\"/api/hsm/describeproducts\",\"Data\":{\"vendorCode\":\"%s\"}}", vendor["code"].(string)),
		}
		response, err := client.DoTeaRequest("POST", "Cspprivate", "2022-02-17", "DescribeProxyCryptoService", "", nil, productsque, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		productsRaw, err := jsonpath.Get("$.data.products", response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		products := []map[string]interface{}{}
		for _, p := range productsRaw.([]interface{}) {
			product := p.(map[string]interface{})
			productMapping := map[string]interface{}{
				"code": product["code"],
				"name": product["name"],
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
