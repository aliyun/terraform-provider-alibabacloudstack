package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alibabacloudstack_ros_stack",
		&resource.Sweeper{
			Name: "alibabacloudstack_ros_stack",
			F:    testSweepRosStack,
		})
}

func testSweepRosStack(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapErrorf(err, "Error getting AlibabacloudStack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   region,
	}
	var response map[string]interface{}
	action := "ListStacks"
	conn, err := client.NewRosClient()
	if err != nil {
		return errmsgs.WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_ros_stack", action, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Stacks", response)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.Stacks", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["StackName"].(string)), strings.ToLower(prefix)) && item["Status"].(string) != "DELETE_COMPLETE" {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ros Stack: %s", item["StackName"].(string))
				continue
			}
			sweeped = true
			action = "DeleteStack"
			request := map[string]interface{}{
				"StackId":  item["StackId"],
				"RegionId": region,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ros Stack (%s): %s", item["StackName"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure Ros Stack have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete Ros Stack success: %s ", item["StackName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlibabacloudStackRosStack_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_ros_stack.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackRosStackMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeRosStack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlibabacloudStackRosStack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackRosStackBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_name":        name,
					"stack_policy_body": `{\"Statement\": [{\"Action\": \"Update:Delete\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}`,
					"template_body":     `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ROS",
					},
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "VpcName",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "InstanceType",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_name":        name,
						"stack_policy_body": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "ROS",
						"parameters.#":      "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_option", "notification_urls", "replacement_option", "retain_all_resources", "retain_resources", "stack_policy_during_update_body", "stack_policy_body", "stack_policy_during_update_url", "stack_policy_url", "template_body", "tags", "template_url", "use_previous_parameters"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_policy_body": `{\"Statement\": [{\"Action\": \"Update:*\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}`,
					"template_body":     `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Description\" : \"模板描述信息，可用于说明模板的适用场景、架构说明等。\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "tf-testacc",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "ECS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_policy_body": CHECKSET,
						"parameters.#":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "ROS Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "ROS Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stack_policy_body": `{\"Statement\": [{\"Action\": \"Update:Delete\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}`,
					"template_body":     `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "VpcName",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "InstanceType",
						},
					},
					"timeout_in_minutes": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_policy_body":  CHECKSET,
						"timeout_in_minutes": "50",
						"parameters.#":       "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ROS",
					},
					"timeout_in_minutes": "60",
					"stack_policy_body":  `{\"Statement\": [{\"Action\": \"Update:*\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}`,
					"template_body":      `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Description\" : \"模板描述信息，可用于说明模板的适用场景、架构说明等。\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "tf-testacc",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "ECS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "ROS",
						"timeout_in_minutes": "60",
						"stack_policy_body":  CHECKSET,
						"parameters.#":       "2",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackRosStackMap = map[string]string{
	"deletion_protection": "Disabled",
	"status":              CHECKSET,
}

func AlibabacloudStackRosStackBasicDependence(name string) string {
	return ""
}
