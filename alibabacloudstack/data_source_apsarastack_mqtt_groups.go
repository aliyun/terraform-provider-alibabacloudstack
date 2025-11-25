package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackMqttGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMqttGroupsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"channel_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMqttGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	instanceId := d.Get("instance_id").(string)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	reqQuery := map[string]interface{}{
		"MqttInstanceId": instanceId,
		"OnsRegionId":    client.RegionId,
		"Platform":       "onsConsole",
		"PreventCache":   time.Now().UnixNano() / 1e6,
		"Dauth_url_hash": "mqtt%2Fconsole%2Finstances%2FinstanceDetail",
		"currentPage":    1,
		"pageSize":       100,
	}

	action := "ConsoleMqttListGroupIdInPage"
	response, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", action, "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "List Mqtt Groups", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if fmt.Sprint(response["Code"]) != "200" {
		return errmsgs.Error(fmt.Sprintf("Mqtt groups list failed!!! Response: %#v", response))
	}

	var data []map[string]interface{}
	if rawData, ok := response["Data"]; ok && rawData != nil {
		dataBytes, _ := json.Marshal(rawData)
		json.Unmarshal(dataBytes, &data)
	}

	var filteredGroups []map[string]interface{}
	var names []string
	var ids []string

	for _, item := range data {
		groupId := item["groupId"].(string)
		resourceId := fmt.Sprintf("%s:%s", instanceId, groupId)

		// Filter by ids
		if len(idsMap) > 0 {
			if _, exists := idsMap[resourceId]; !exists {
				continue
			}
		}

		// Filter by name_regex
		if nameRegex != nil {
			if !nameRegex.MatchString(groupId) {
				continue
			}
		}

		filteredGroups = append(filteredGroups, map[string]interface{}{
			"group_id":           groupId,
			"instance_id":        item["namespaceId"],
			"create_time":        item["createTime"],
			"update_time":        item["updateTime"],
			"independent_naming": item["independentNaming"],
			"channel_name":       item["channelName"],
			"region_name":        item["regionName"],
		})

		names = append(names, groupId)
		ids = append(ids, resourceId)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("groups", filteredGroups); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
