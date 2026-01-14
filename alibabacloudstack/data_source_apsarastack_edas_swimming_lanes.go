package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackEdasSwimmingLanes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEdasSwimmingLanesRead,

		Schema: map[string]*schema.Schema{
			"logical_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"swimming_lanes": {
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
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logical_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apps": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"condition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rest_items": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cond": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"lane_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEdasSwimmingLanesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	if v, ok := d.GetOk("logical_region_id"); ok {
		request["logicalRegionId"] = v.(string)
	} else {
		request["logicalRegionId"] = client.RegionId
	}
	request["GroupId"] = d.Get("group_id")
	response, err := client.DoTeaRequest("GET", "Edas", "2017-08-01", "ListSwimmingLane", "/pop/v5/trafficmgnt/swimming_lanes", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane", "ListSwimmingLane", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.Error("describe edas swimming_lane failed for: " + response["Message"].(string))
	}
	data, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_edas_swimming_lane_group", "ListSwimminglane", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	idsMap := getIdsStringFilter(d)

	var ids []string
	datas := make([]interface{}, 0)
	for _, v := range data.([]interface{}) {
		lane := v.(map[string]interface{})
		id := fmt.Sprintf("%s:%s:%s", lane["NamespaceId"], fmt.Sprint(lane["GroupId"]), fmt.Sprint(lane["Id"]))
		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(lane["Name"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[id]; !exist {
				continue
			}
		}
		i := map[string]interface{}{
			"id": id,

			"name": lane["Name"],

			"logical_region_id": lane["NamespaceId"],

			"group_id": fmt.Sprint(lane["GroupId"]),

			"lane_id": fmt.Sprint(lane["Id"]),
		}
		applist := lane["SwimmingLaneAppRelationShipList"].([]interface{})
		apps := make([]string, 0)
		for _, v := range applist {
			app := v.(map[string]interface{})
			apps = append(apps, app["AppId"].(string))
		}
		i["apps"] = apps
		entryRules := make([]map[string]interface{}, 0)
		err = json.Unmarshal([]byte(lane["EntryRule"].(string)), &entryRules)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		restItems := make([]map[string]interface{}, 0)
		if len(entryRules) > 0 {
			entryRule := entryRules[0]
			i["condition"] = entryRule["condition"]
			i["enabled"] = entryRule["enable"]
			i["path"] = entryRule["path"]
			i["priority"] = entryRule["priority"]

			for _, v := range entryRule["restItems"].([]interface{}) {
				item := v.(map[string]interface{})
				restItems = append(restItems, map[string]interface{}{
					"type":     item["type"],
					"name":     item["name"],
					"value":    item["value"],
					"cond":     item["cond"],
					"operator": item["operator"],
				})
			}
			i["rest_items"] = restItems
		}

		datas = append(datas, i)

		ids = append(ids, id)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("swimming_lanes", datas); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	return nil
}
