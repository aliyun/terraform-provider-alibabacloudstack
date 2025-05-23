package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_arms_dispatch_rule", &resource.Sweeper{
		Name: "alibabacloudstack_arms_dispatch_rule",
		F:    testSweepArmsDispatchRule,
	})
}

func testSweepArmsDispatchRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapErrorf(err, "error getting AlibabacloudStack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	action := "ListDispatchRule"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var response map[string]interface{}
	conn, err := client.NewArmsClient()
	if err != nil {
		return errmsgs.WrapError(err)
	}
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] %s failed: %v", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.DispatchRules", response)
	if err != nil {
		log.Printf("[ERROR] %v", errmsgs.WrapError(err))
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["Name"])
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping dispatch rule: %s ", name)
			continue
		}
		log.Printf("[INFO] delete dispatch rule: %s ", name)

		action = "DeleteDispatchRule"
		request = map[string]interface{}{
			"Id":       fmt.Sprint(item["RuleId"]),
			"RegionId": client.RegionId,
		}
		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s failed: %v", action, err)
		}
	}
	return nil
}

func TestAccAlibabacloudStackArmsDispatchRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_arms_dispatch_rule.default"
	ra := resourceAttrInit(resourceId, ArmsDispatchRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeArmsDispatchRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccArmsDispatchRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ArmsDispatchRuleBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dispatch_rule_name": "${var.name}",
					"group_rules": []map[string]interface{}{
						{
							"group_wait_time": "5",
							"group_interval":  "15",
							"grouping_fields": []string{"alertname"},
							"repeat_interval": "61",
						},
					},
					"dispatch_type": "CREATE_ALERT",
					"label_match_expression_grid": []map[string]interface{}{
						{
							"label_match_expression_groups": []map[string]interface{}{
								{
									"label_match_expressions": []map[string]interface{}{
										{
											"key":      "_aliyun_arms_involvedObject_kind",
											"value":    "app",
											"operator": "eq",
										},
									},
								},
							},
						},
					},
					"notify_rules": []map[string]interface{}{
						{
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_id": "${alibabacloudstack_arms_alert_contact.default.id}",
									"notify_type":      "ARMS_CONTACT",
									"name":             "${var.name}",
								},
							},
							"notify_channels": []string{"dingTalk", "wechat"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dispatch_rule_name":            CHECKSET,
						"group_rules.#":                 "1",
						"dispatch_type":                 "CREATE_ALERT",
						"label_match_expression_grid.#": "1",
						"notify_rules.#":                "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dispatch_rule_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dispatch_rule_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_rules": []map[string]interface{}{
						{
							"group_wait_time": "10",
							"group_interval":  "25",
							"grouping_fields": []string{"alertname2"},
							"repeat_interval": "70",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_rules.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rules": []map[string]interface{}{
						{
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_id": "${alibabacloudstack_arms_alert_contact.default.id}",
									"notify_type":      "ARMS_CONTACT",
									"name":             "${var.name}",
								},
								{
									"notify_object_id": "${alibabacloudstack_arms_alert_contact_group.default.id}",
									"notify_type":      "ARMS_CONTACT_GROUP",
									"name":             "${var.name}",
								},
							},
							"notify_channels": []string{"dingTalk"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_rules.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dispatch_type": "DISCARD_ALERT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dispatch_type": "DISCARD_ALERT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"label_match_expression_grid": []map[string]interface{}{
						{
							"label_match_expression_groups": []map[string]interface{}{
								{
									"label_match_expressions": []map[string]interface{}{
										{
											"key":      "_aliyun_arms_involvedObject_kind",
											"value":    "app",
											"operator": "eq",
										},
										{
											"key":      "_aliyun_arms_alert_name",
											"value":    "tf-testaccapp",
											"operator": "eq",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"label_match_expression_grid.#": "1",
						"notify_rules.#":                "0",
						"group_rules.#":                 "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dispatch_rule_name": "${var.name}",
					"group_rules": []map[string]interface{}{
						{
							"group_wait_time": "5",
							"group_interval":  "15",
							"grouping_fields": []string{"alertname"},
							"repeat_interval": "80",
						},
					},
					"dispatch_type": "CREATE_ALERT",
					"label_match_expression_grid": []map[string]interface{}{
						{
							"label_match_expression_groups": []map[string]interface{}{
								{
									"label_match_expressions": []map[string]interface{}{
										{
											"key":      "_aliyun_arms_involvedObject_kind",
											"value":    "app",
											"operator": "eq",
										},
									},
								},
							},
						},
					},
					"notify_rules": []map[string]interface{}{
						{
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_id": "${alibabacloudstack_arms_alert_contact.default.id}",
									"notify_type":      "ARMS_CONTACT",
									"name":             "${var.name}",
								},
							},
							"notify_channels": []string{"dingTalk", "wechat"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dispatch_rule_name":            CHECKSET,
						"group_rules.#":                 "1",
						"dispatch_type":                 "CREATE_ALERT",
						"label_match_expression_grid.#": "1",
						"notify_rules.#":                "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dispatch_type"},
			},
		},
	})
}

var ArmsDispatchRuleMap = map[string]string{
	"status": CHECKSET,
}

func ArmsDispatchRuleBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alibabacloudstack_arms_alert_contact" "default" {
  alert_contact_name = "${var.name}"
  email = "${var.name}@aaa.com"
ding_robot_webhook_url= "https://oapi.dingtalk.com/robot/send?access_token=91f2f7"
	phone_num=  "14345546451"
}
resource "alibabacloudstack_arms_alert_contact_group" "default" {
  alert_contact_group_name = "${var.name}"
  contact_ids = [alibabacloudstack_arms_alert_contact.default.id]
}
`, name)
}
