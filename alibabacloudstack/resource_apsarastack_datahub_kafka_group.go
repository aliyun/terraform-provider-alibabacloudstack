package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDatahubKafkaGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Required: true,
			},
			"topic_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackDatahubKafkaGroupCreate, resourceAlibabacloudStackDatahubKafkaGroupRead, resourceAlibabacloudStackDatahubKafkaGroupUpdate, resourceAlibabacloudStackDatahubKafkaGroupDelete)
	return resource
}

func resourceAlibabacloudStackDatahubKafkaGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	projectName := d.Get("project_name").(string)
	groupName := d.Get("group_name").(string)
	comment := d.Get("comment").(string)

	// Prepare the request parameters for CreateKafkaGroup API
	reqQuery := map[string]interface{}{
		"ProjectName": projectName,
		"GroupName":   groupName,
		"Comment":     comment,
	}

	// Call the CreateKafkaGroup API
	if _, err := client.DoTeaRequest("POST", "datahub", "2019-11-20", "CreateKafkaGroup", "", nil, reqQuery, nil); err != nil {
		return err
	}

	// Generate and set the resource ID: {ProjectName}:{GroupName}
	resourceId := fmt.Sprintf("%s:%s", projectName, groupName)
	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackDatahubKafkaGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	datahubService := DatahubService{client}

	object, err := datahubService.DescribeDatahubKafkaGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_datahub_kafka_group datahubService.DescribeDatahubKafkaGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	// Set resource data from the returned object
	d.Set("group_name", object["GroupName"])
	d.Set("comment", object["Comment"])
	d.Set("create_time", object["CreateTime"])

	// Handle topic_list
	if topicList, ok := object["TopicList"].([]interface{}); ok {
		topics := make([]string, 0, len(topicList))
		for _, topic := range topicList {
			if topicStr, ok := topic.(string); ok {
				topics = append(topics, topicStr)
			}
		}
		d.Set("topic_list", topics)
	}

	// Extract project name from ID and set it
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) == 2 {
		d.Set("project_name", parts[0])
	}

	return nil
}

func resourceAlibabacloudStackDatahubKafkaGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	projectName := parts[0]
	groupName := parts[1]

	if d.HasChange("topic_list") {
		topicList := d.Get("topic_list").(*schema.Set).List()
		kafkaGroupTopicList := make([]string, 0)
		for _, item := range topicList {
			kafkaGroupTopicList = append(kafkaGroupTopicList, item.(string))
		}

		var kafkaGroupTopicListStr string
		if v, err := json.Marshal(kafkaGroupTopicList) ; err == nil {
			kafkaGroupTopicListStr = string(v)
		} else {
			return err
		}
		
		query := map[string]interface{}{
			"RegionId":            client.RegionId,
			"ProjectName":         projectName,
			"GroupName":           groupName,
			"KafkaGroupTopicList": kafkaGroupTopicListStr,
		}

		_, err = client.DoTeaRequest("POST", "datahub", "2019-11-20", "UpdateTopicsForKafkaGroup", "", nil, query, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "UpdateTopicsForKafkaGroup", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	noUpdatesAllowedCheck(d, []string{"comment"})

	return nil
}

func resourceAlibabacloudStackDatahubKafkaGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	projectName := parts[0]
	groupName := parts[1]

	reqQuery := map[string]interface{}{
		"RegionId":    d.Get("region_id"),
		"ProjectName": projectName,
		"GroupName":   groupName,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "datahub", "2019-11-20", "DeleteKafkaGroup", "", nil, reqQuery, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"ResourceNotFound"}) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ResourceNotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteKafkaGroup", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
