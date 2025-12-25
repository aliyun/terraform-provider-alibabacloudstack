package alibabacloudstack

import (
	"encoding/json"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackRegionsByProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackRegionsByProductRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
			},
			"region_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackRegionsByProductRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "GetRegionsByProduct", "/ascm/basic/cmdb/region/product")
	response := RegionsByProduct{}

	for {
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_regions_by_product", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Code == 200 || len(response.Body.RegionList) < 1 {
			break
		}

	}

	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Body.RegionList {
		mapping := map[string]interface{}{
			"region_id":   rg.RegionID,
			"region_type": rg.RegionType,
		}
		ids = append(ids, rg.RegionID)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("region_list", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
