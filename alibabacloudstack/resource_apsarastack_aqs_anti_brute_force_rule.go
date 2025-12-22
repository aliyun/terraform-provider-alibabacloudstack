package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAqsAntiBruteForceRule() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"span": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"fail_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"forbidden_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"default_rule": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"enable_smart_rule": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"machine_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackAqsAntiBruteForceRuleCreate, resourceAlibabacloudStackAqsAntiBruteForceRuleRead, resourceAlibabacloudStackAqsAntiBruteForceRuleUpdate, resourceAlibabacloudStackAqsAntiBruteForceRuleDelete)
	return resource
}

func resourceAlibabacloudStackAqsAntiBruteForceRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"From":          "sas",
		"Name":          d.Get("name"),
		"Span":          d.Get("span"),
		"FailCount":     d.Get("fail_count"),
		"ForbiddenTime": d.Get("forbidden_time"),
		"DefaultRule":   d.Get("default_rule"),
	}
	if v, ok := d.GetOk("instance_ids"); ok {
		instance_ids := v.(*schema.Set).List()
		reqQuery["UuidList"] = InstanceIdsHandler(instance_ids, "instanceId", meta)
	}

	resp, err := client.DoTeaRequest("POST", "aegis", "2016-11-11", "CreateAntiBruteForceRule", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	data, ok := resp["data"].(string)
	ruledata := make(map[string]interface{})
	if ok {
		log.Printf("[DEBUG] alibabacloudstack_aqs_operate_common_overall_config CreateAntiBruteForceRule response.data: %s", data)
		err = json.Unmarshal([]byte(data), &ruledata)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}
	ruleId, err := jsonpath.Get("$.CreateAntiBruteForceRule.RuleId", ruledata)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.SetId(fmt.Sprint(ruleId))

	return nil
}

func resourceAlibabacloudStackAqsAntiBruteForceRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	aqsService := AqsService{client}
	targetRule, err := aqsService.DescribeAntiBruteForceRule(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_operate_common_overall_config", "DescribeAntiBruteForceRule", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.Set("name", targetRule["Name"])
	d.Set("span", targetRule["Span"])
	d.Set("fail_count", targetRule["FailCount"])
	d.Set("forbidden_time", targetRule["ForbiddenTime"])
	d.Set("default_rule", targetRule["DefaultRule"])
	instance_ids := InstanceIdsHandler(targetRule["UuidList"].([]interface{}), "uuid", meta)
	d.Set("instance_ids", instance_ids)
	d.Set("enable_smart_rule", targetRule["EnableSmartRule"])
	d.Set("machine_count", targetRule["MachineCount"])
	d.Set("create_timestamp", targetRule["CreateTimestamp"])

	return nil
}

func resourceAlibabacloudStackAqsAntiBruteForceRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("name", "span", "fail_count", "forbidden_time", "default_rule", "instance_ids") {
		request := map[string]interface{}{
			"Id":            d.Id(),
			"From":          "sas",
			"Name":          d.Get("name"),
			"Span":          d.Get("span"),
			"FailCount":     d.Get("fail_count"),
			"ForbiddenTime": d.Get("forbidden_time"),
			"DefaultRule":   d.Get("default_rule"),
		}
		if v, ok := d.GetOk("instance_ids"); ok {
			instance_ids := v.(*schema.Set).List()
			request["UuidList"] = InstanceIdsHandler(instance_ids, "instanceId", meta)
		}

		_, err := client.DoTeaRequest("POST", "aegis", "2016-11-11", "ModifyAntiBruteForceRule", "", nil, request, nil)
		if err != nil {
			return fmt.Errorf("failed to modify anti brute force rule: %v", err)
		}
	}

	return nil
}

func resourceAlibabacloudStackAqsAntiBruteForceRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	reqQuery := map[string]interface{}{
		"From": "sas",
		"Ids":  []string{d.Id()},
	}

	_, err := client.DoTeaRequest("POST", "aegis", "2016-11-11", "DeleteAntiBruteForceRule", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteAntiBruteForceRule", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func InstanceIdsHandler(ids []interface{}, parameterType string, meta interface{}) []string {
	client := meta.(*connectivity.AlibabacloudStackClient)
	aqsService := AqsService{client}
	instances, err := aqsService.DescribeCloudCenterInstances()
	if err != nil {
		return nil
	}
	result := make([]string, 0)
	if parameterType == "uuid" {
		for _, id := range ids {
			for _, v := range instances {
				instance := v.(map[string]interface{})
				if instance["Uuid"].(string) == id.(string) {
					result = append(result, instance["InstanceId"].(string))
				}
			}
		}
	} else if parameterType == "instanceId" {
		for _, id := range ids {
			for _, v := range instances {
				instance := v.(map[string]interface{})
				if instance["InstanceId"].(string) == id.(string) {
					result = append(result, instance["Uuid"].(string))
				}
			}
		}
	}
	return result
}
