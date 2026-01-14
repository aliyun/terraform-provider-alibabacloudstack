package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDatahubKafkaGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDatahubKafkaGroupsRead,

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"kafka_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_modify_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDatahubKafkaGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	projectName := d.Get("project_name").(string)

	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return fmt.Errorf("invalid 'name_regex': %w", err)
		}
		nameRegex = r
	}

	var allGroups []map[string]interface{}
	pageNumber := 1
	pageSize := 50

	for {
		// Prepare request parameters for ListKafkaGroup API
		query := map[string]interface{}{
			"ProjectName": projectName,
			"PageSize":    pageSize,
			"PageNumber":  pageNumber,
		}

		// Call ListKafkaGroup API
		resp, err := client.DoTeaRequest("GET", "datahub", "2019-11-20", "ListKafkaGroup", "", nil, query, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"ResourceNotFound"}) {
				break
			}
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_datahub_kafka_groups", "ListKafkaGroup", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		groupList, ok := resp["List"].([]interface{})
		if !ok {
			groupList = []interface{}{}
		}

		for _, item := range groupList {
			group, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			groupName, _ := group["GroupName"].(string)
			id := fmt.Sprintf("%s:%s", projectName, groupName)

			// Filter by ids
			if len(idsMap) > 0 {
				if _, exists := idsMap[id]; !exists {
					continue
				}
			}

			// Filter by name_regex
			if nameRegex != nil {
				if !nameRegex.MatchString(groupName) {
					continue
				}
			}

			allGroups = append(allGroups, group)
		}

		totalCount, _ := resp["TotalCount"].(float64)
		pageCount, _ := resp["PageCount"].(float64)

		if int64(pageNumber) >= int64(pageCount) || int64(pageSize)*int64(pageNumber) >= int64(totalCount) {
			break
		}
		pageNumber++
	}

	// Prepare the result
	var ids []string
	var groups []map[string]interface{}

	for _, group := range allGroups {
		groupName := group["GroupName"].(string)
		id := fmt.Sprintf("%s:%s", projectName, groupName)

		mapping := map[string]interface{}{
			"group_name":       groupName,
			"comment":          group["Comment"],
			"create_time":      group["CreateTime"],
			"last_modify_time": group["LastModifyTime"],
			"creator":          group["Creator"],
		}

		// Handle topic_list
		topicList := make([]string, 0)
		if topics, ok := group["TopicList"].([]interface{}); ok {
			for _, topic := range topics {
				if topicStr, ok := topic.(string); ok {
					topicList = append(topicList, topicStr)
				}
			}
		}
		mapping["topic_list"] = topicList

		ids = append(ids, id)
		groups = append(groups, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("kafka_groups", groups); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
