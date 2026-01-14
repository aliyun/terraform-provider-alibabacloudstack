package alibabacloudstack

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackMqttTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMqttTopicsRead,

		Schema: map[string]*schema.Schema{
			"store_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"topics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"store_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"order_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"relation": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"relation_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"namespace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unit_flag": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channel_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channel_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMqttTopicsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	instanceId := d.Get("store_instance_id").(string)

	// Build request query parameters
	reqQuery := map[string]interface{}{
		"PreventCache":   time.Now().UnixNano() / 1e6,
		"Dauth_url_hash": "mqtt%2Fconsole%2Finstances%2FinstanceDetail",
		"OnsRegionId":    client.RegionId,
		"InstanceId":     instanceId,
		"CurrentPage":    1,
		"PageSize":       100,
		"isFuzzy":        false,
	}

	// Call API to get topic list
	resp, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleTopicListInPage", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_mqtt_topic", "ConsoleTopicListInPage", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if fmt.Sprint(resp["Code"]) != "200" {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_mqtt_topic", "ConsoleTopicListInPage", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	data, ok := resp["Data"].([]interface{})
	if !ok {
		data = []interface{}{}
	}

	// Process filtering by ids and name_regex
	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var filteredTopics []interface{}
	for _, item := range data {
		topicData, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		topicName := topicData["topic"].(string)
		topicId := topicData["namespaceId"].(string) + ":" + topicName

		// Filter by ids
		if len(idsMap) > 0 {
			if _, exists := idsMap[topicId]; !exists {
				continue
			}
		}

		// Filter by name_regex
		if nameRegex != nil && !nameRegex.MatchString(topicName) {
			continue
		}

		filteredTopics = append(filteredTopics, topicData)
	}

	// Construct topics schema data
	topics := make([]map[string]interface{}, 0, len(filteredTopics))
	ids := make([]string, 0, len(filteredTopics))

	for _, item := range filteredTopics {
		topicData := item.(map[string]interface{})

		mapping := map[string]interface{}{
			"store_instance_id":  topicData["namespaceId"],
			"topic":              topicData["topic"],
			"remark":             topicData["remark"],
			"order_type":         topicData["orderType"],
			"independent_naming": topicData["independentNaming"],
			"update_time":        topicData["updateTime"],
			"relation":           topicData["relation"],
			"relation_name":      topicData["relationName"],
			"create_time":        topicData["createTime"],
			"namespace_id":       topicData["namespaceId"],
			"unit_flag":          topicData["unitFlag"],
			"status_name":        topicData["statusName"],
			"channel_name":       topicData["channelName"],
			"channel_id":         topicData["channelId"],
			"status":             topicData["status"],
		}

		topics = append(topics, mapping)
		ids = append(ids, topicData["namespaceId"].(string)+":"+topicData["topic"].(string))
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("topics", topics); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
