package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackDatahubSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDatahubSubscriptionsRead,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sub_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDatahubSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	projectName := d.Get("project_name").(string)
	topicName := d.Get("topic_name").(string)

	// 1. Call ListSubscriptions to get all subscriptions for the project and topic
	listReq := client.NewCommonRequest("GET", "datahub", "2019-11-20", "ListSubscriptions", "")
	listReq.QueryParams["ProjectName"] = projectName
	listReq.QueryParams["TopicName"] = topicName
	listReq.QueryParams["PageNumber"] = "1"
	listReq.QueryParams["PageSize"] = "1000"

	listResp, err := client.ProcessCommonRequest(listReq)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	var subscriptions []SubscriptionEntry 
	listResult := &ListSubscriptionResult{}
	err = json.Unmarshal(listResp.GetHttpContentBytes(), listResult)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	subscriptions = listResult.List.Subscriptions

	var filteredSubscriptions []map[string]interface{}
	filterIds := map[string]interface{}{}
	subIds := []string{}

	if v, ok := d.GetOk("sub_ids"); ok {
		for _, vv := range v.([]interface{}) {
			filterIds[vv.(string)] = nil
		}
	}
	for _, subscription := range subscriptions {
		if _, existed := filterIds[subscription.SubId]; len(filterIds) > 0 &&  !existed{
			continue
		}
		subIds = append(subIds, subscription.SubId)
		mapping := map[string]interface{}{
			"id" : fmt.Sprintf("%s:%s:%s", projectName, topicName, subscription.SubId),
			"project_name":     projectName,
			"topic_name":       topicName,
			"sub_id":           subscription.SubId,
			"application_name": subscription.Application,
			"comment":          subscription.Comment,
		}
		filteredSubscriptions = append(filteredSubscriptions, mapping)
	}

	// 4. Set ID and results
	d.Set("sub_ids", subIds)
	d.SetId(dataResourceIdHash(subIds))
	if err := d.Set("subscriptions", filteredSubscriptions); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
