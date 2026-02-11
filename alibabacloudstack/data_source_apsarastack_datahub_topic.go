package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackDatahubTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDatahubTopicsRead,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topics": {
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
						"shard_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"life_cycle": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDatahubTopicsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	projectName := d.Get("project_name").(string)

	// 1. Call ListTopics to get all topic names in the project
	listReq := client.NewCommonRequest("GET", "datahub", "2019-11-20", "ListTopics", "")
	listReq.QueryParams["ProjectName"] = projectName
	listReq.QueryParams["PageNumber"] = "1"
	listReq.QueryParams["PageSize"] = "1000"

	var topicNames []string
	var filteredTopics []map[string]interface{}
	listResp, err := client.ProcessCommonRequest(listReq)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ResourceNotFound") {
			d.Set("names", topicNames)
			d.SetId(dataResourceIdHash(topicNames))
			if err := d.Set("topics", filteredTopics); err != nil {
				return errmsgs.WrapError(err)
			}
			return nil
		}
		return errmsgs.WrapError(err)
	}

	filterNames := map[string]interface{}{}
	nameRegex, nameRegexOk := d.GetOk("name_regex")
	var r *regexp.Regexp
	if nameRegexOk && nameRegex.(string) != "" {
		r, err = regexp.Compile(nameRegex.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if v, ok := d.GetOk("names"); ok {
		for _, vv := range v.([]interface{}) {
			filterNames[vv.(string)] = nil
		}
	}

	listResult := &ListTopicResult{}
	err = json.Unmarshal(listResp.GetHttpContentBytes(), listResult)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	for _, topic := range listResult.List.Topic {
		if nameRegexOk && nameRegex.(string) != "" {
			if !r.MatchString(topic.TopicName) {
				continue
			}
		}
		if _, existed := filterNames[topic.TopicName]; len(filterNames) > 0 && !existed {
			continue
		}
		topicNames = append(topicNames, topic.TopicName)
		filteredTopics = append(filteredTopics, map[string]interface{}{
			"id":          fmt.Sprintf("%s:%s", projectName, topic.TopicName),
			"name":        topic.TopicName,
			"shard_count": topic.ShardCount,
			"life_cycle":  topic.LifeCycle,
			"comment":     topic.Comment,
			"record_type": topic.RecordType,
			"create_time": strconv.FormatInt(topic.CreateTime, 10),
		})
	}
	d.Set("names", topicNames)
	d.SetId(dataResourceIdHash(topicNames))
	if err := d.Set("topics", filteredTopics); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
