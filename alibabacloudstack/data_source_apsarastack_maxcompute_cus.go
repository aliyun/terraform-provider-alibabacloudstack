package alibabacloudstack

import (
	"encoding/json"
	"regexp"
	"strconv"

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
	request := map[string]interface{}{
		"Region":      client.RegionId,
		"Action":      "ListOdpsCusForAscm",
		"AccessKeyId": client.AccessKey,
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		request["Cluster"] = v.(string)
	}
	response, err := client.DoTeaRequest("GET", "dataworks-private-cloud", "2019-01-17", "ListOdpsCusForAscm", "", nil, request, nil)
	addDebug("ListOdpsCusForAscm", response, request)
	if err != nil {
		err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_cu", "ListOdpsCusForAscm", errmsgs.AlibabacloudStackSdkGoERROR)
		return err
	}
	objects, err := jsonpath.Get("$.Data.data", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "maxcompute_cu", "$.Data.data", response)
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var t []map[string]interface{}
	var ids []string
	for _, object := range objects.([]interface{}) {
		cu_raw := object.(map[string]interface{})
		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(cu_raw["quota_name"].(string)) {
				continue
			}
		}
		var cu_num int
		switch v := cu_raw["max_cu"].(type) {
		case string:
			cu_num, err = strconv.Atoi(v)
			if err != nil {
				return errmsgs.WrapError(errmsgs.Error("illegal max_cu value"))
			}
		case json.Number:
			var floatVal float64
			floatVal, err = v.Float64()
			if err != nil {
				return errmsgs.WrapError(errmsgs.Error("illegal max_cu value"))
			}
			cu_num = int(floatVal)
		case int:
			cu_num = v
		case float64:
			cu_num = int(v)
		default:
			return errmsgs.WrapError(errmsgs.Error("illegal max_cu value type"))
		}
		cu := map[string]interface{}{
			"id":           cu_raw["id"].(string),
			"cu_name":      cu_raw["quota_name"].(string),
			"cu_num":       cu_num,
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
	return nil
}
