package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackEdasSwimmingLaneGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEdasSwimmingLaneGroupsRead,

		Schema: map[string]*schema.Schema{
			"logical_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				//ValidateFunc: validation.StringIsValidRegExp,
				ForceNew: true,
			},
			"swimming_lane_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entry_app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apps": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"logical_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"strategy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEdasSwimmingLaneGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	if v, ok := d.GetOk("logical_region_id"); ok {
		request["logicalRegionId"] = v.(string)
	} else {
		request["logicalRegionId"] = client.RegionId
	}
	response, err := client.DoTeaRequest("GET", "Edas", "2017-08-01", "ListSwimmingLaneGroup", "/pop/v5/trafficmgnt/swimming_lane_groups", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane_group", "ListSwimmingLaneGroup", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.Error("describe edas swimming lane group failed for: " + response["Message"].(string))
	}
	data, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane_group", "ListSwimmingLaneGroup", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var ids []string
	datas := make([]interface{}, 0)
	for _, data := range data.([]interface{}) {
		laneGroup := data.(map[string]interface{})
		id := fmt.Sprintf("%s:%s", laneGroup["NamespaceId"], fmt.Sprint(laneGroup["Id"]))
		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(laneGroup["Name"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[id]; !exist {
				continue
			}
		}
		applist := laneGroup["ApplicationList"].([]interface{})
		apps := make([]string, 0)
		for _, v := range applist {
			app := v.(map[string]interface{})
			apps = append(apps, app["AppId"].(string))
		}
		entryApplication := laneGroup["EntryApplication"].(map[string]interface{})
		i := map[string]interface{}{
			"id": id,

			"name": laneGroup["Name"],

			"logical_region_id": laneGroup["NamespaceId"],

			"group_id": fmt.Sprint(laneGroup["Id"]),

			"entry_app_id": entryApplication["AppId"],

			"apps": apps,
		}

		datas = append(datas, i)

		ids = append(ids, id)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("swimming_lane_groups", datas); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	return nil
}
