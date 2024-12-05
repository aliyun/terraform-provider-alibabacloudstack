package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackMaxcomputeCus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMaxcomputeCusRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cus": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
							ForceNew: true,
						},
						"cu_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"cu_num": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMaxcomputeCusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.Get("name_regex").(string) != "" && d.Get("cluster_name").(string) != "" {
		err := errmsgs.Error("Only one filter condition can be set")
		return err
	}
	if d.Get("name_regex").(string) == "" && d.Get("cluster_name").(string) == "" {
		err := errmsgs.Error("At least one filter condition needs to be set.")
		return err
	}

	action := "ListOdpsCus"

	request := map[string]interface{}{
		"RegionName":      client.RegionId,
		"Product":         "ascm",
		"OrganizationId":  client.Department,
		"ResourceGroupId": client.ResourceGroup,
	}

	filter_query := ""
	if d.Get("name_regex").(string) != "" {
		request["Type"] = "cuName"
		request["CuName"] = d.Get("name_regex").(string)
		filter_query = d.Get("name_regex").(string)
	}

	if d.Get("cluster_name").(string) != "" {
		request["Type"] = "clusterName"
		request["ClusterName"] = d.Get("cluster_name").(string)
		filter_query = d.Get("cluster_name").(string)
	}

	response, err := client.DoTeaRequest("POST", "ASCM", "2019-05-10", action, "", nil, request)

	if err != nil {
		return err
	}

	if errmsgs.IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("MaxcomputeProject", filter_query)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		return err
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		err = errmsgs.Error("ListOdpsCus failed for " + response["asapiErrorMessage"].(string))
		return err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, filter_query, "$", response)
	}
	objects := v.(map[string]interface{})["data"].([]interface{})

	var t []map[string]interface{}
	var ids []string
	for _, object := range objects {
		cu_raw := object.(map[string]interface{})
		max_cu, err := cu_raw["max_cu"].(json.Number).Float64()
		if err != nil {
			return errmsgs.WrapError(errmsgs.Error("illegal max_cu value"))
		}
		cu := map[string]interface{}{
			"id":           cu_raw["id"].(string),
			"cu_name":      cu_raw["quota_name"].(string),
			"cu_num":       int64(max_cu),
			"cluster_name": cu_raw["cluster"].(string),
		}
		t = append(t, cu)
		ids = append(ids, cu["id"].(string))

	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("cus", t); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), t)
	}
	return nil
}
