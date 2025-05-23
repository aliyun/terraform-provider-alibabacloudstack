package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudEssAlarmBasic(t *testing.T) {
	var v ess.Alarm
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand)
	var basicMap = map[string]string{
		// "name":                   name,
		// "description":            "Acc alarm test",
		// "alarm_actions.#":        "1",
		"scaling_group_id": CHECKSET,
		// "metric_type":            "system",
		// "metric_name":            "CpuUtilization",
		// "period":                 "300",
		// "statistics":             "Average",
		// "comparison_operator":    ">=",
		// "evaluation_count":       "2",
		// "threshold":              "200.3",
		"cloud_monitor_group_id": NOSET,
		//"enable":                 "true",
	}
	resourceId := "alibabacloudstack_ess_alarm.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssAlarmConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers: testAccProviders,
		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// "name":                name,
					"description":         "Acc alarm test",
					"alarm_actions":       []string{"${alibabacloudstack_ess_scaling_rule.default.0.ari}"},
					"scaling_group_id":    "${alibabacloudstack_ess_scaling_group.default.id}",
					"metric_type":         "system",
					"metric_name":         "CpuUtilization",
					"period":              "300",
					"statistics":          "Average",
					"threshold":           "200.3",
					"comparison_operator": ">=",
					"evaluation_count":    "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"name": fmt.Sprintf("tf-testAccEssAlarm-%d", rand),
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"name": fmt.Sprintf("tf-testAccEssAlarm-%d", rand),
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"description": "Acc alarm test 123",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"description": "Acc alarm test 123",
			// 		}),
			// 	),
			// },
			{
				Config: testAccConfig(map[string]interface{}{
					"alarm_actions": "${alibabacloudstack_ess_scaling_rule.default.*.ari}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_actions.#": "2",
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"dimensions": map[string]string{
			// 			"device": "eth0",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"dimensions.%": "1",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"metric_name": "PackagesNetIn",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"metric_name": "PackagesNetIn",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"period": "120",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"period": "120",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"statistics": "Minimum",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"statistics": "Minimum",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"threshold": "200.5",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"threshold": "200.5",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"comparison_operator": ">",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"comparison_operator": ">",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"evaluation_count": "3",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"evaluation_count": "3",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"cloud_monitor_group_id": "5390371",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"cloud_monitor_group_id": "5390371",
			// 		}),
			// 	),
			// },
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "true",
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"name":                name,
			// 		"description":         "Acc alarm test",
			// 		"alarm_actions":       []string{"${alibabacloudstack_ess_scaling_rule.default.0.ari}"},
			// 		"scaling_group_id":    "${alibabacloudstack_ess_scaling_group.default.id}",
			// 		"metric_type":         "system",
			// 		"metric_name":         "CpuUtilization",
			// 		"period":              "300",
			// 		"statistics":          "Average",
			// 		"threshold":           "200.3",
			// 		"comparison_operator": ">=",
			// 		"evaluation_count":    "2",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"name":                name,
			// 			"description":         "Acc alarm test",
			// 			"alarm_actions.#":     "1",
			// 			"scaling_group_id":    CHECKSET,
			// 			"metric_type":         "system",
			// 			"metric_name":         "CpuUtilization",
			// 			"period":              "300",
			// 			"statistics":          "Average",
			// 			"comparison_operator": ">=",
			// 			"evaluation_count":    "2",
			// 			"threshold":           "200.3",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"scaling_group_id": "${alibabacloudstack_ess_scaling_group.new.id}",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"scaling_group_id": CHECKSET,
			// 		}),
			// 	),
			// },
		},
	})
}

func TestAccalibabacloudstackEssAlarmWithExpression(t *testing.T) {
	var v ess.Alarm
	rand := getAccTestRandInt(10000, 999999)
	var basicMap = map[string]string{
		// "name":                   fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand),
		// "description":            "Acc alarm test",
		// "alarm_actions.#":        "1",
		"scaling_group_id": CHECKSET,
		// "metric_type":            "system",
		// "evaluation_count":       "2",
		"cloud_monitor_group_id": NOSET,
		// "enable":                 "true",
		// "expressions.#":          "1",
	}
	resourceId := "alibabacloudstack_ess_alarm.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssAlarmConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// "name":             name,
					"description":      "Acc alarm test",
					"alarm_actions":    []string{"${alibabacloudstack_ess_scaling_rule.default.0.ari}"},
					"scaling_group_id": "${alibabacloudstack_ess_scaling_group.default.id}",
					"metric_type":      "system",
					"evaluation_count": "2",
					"expressions": []map[string]string{{
						"period":              "300",
						"statistics":          "Average",
						"metric_name":         "CpuUtilization",
						"threshold":           "200.3",
						"comparison_operator": ">=",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"name": fmt.Sprintf("tf-testAccEssAlarmExpressions-%d", rand),
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"name": fmt.Sprintf("tf-testAccEssAlarmExpressions-%d", rand),
			// 		}),
			// 	),
			// },

			{
				Config: testAccConfig(map[string]interface{}{
					"expressions": []map[string]string{
						{
							"period":              "120",
							"statistics":          "Average",
							"metric_name":         "CpuUtilization",
							"threshold":           "40.1",
							"comparison_operator": ">=",
						},
						{
							"period":              "120",
							"statistics":          "Minimum",
							"metric_name":         "MemoryUtilization",
							"threshold":           "99.9",
							"comparison_operator": ">",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expressions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expressions_logic_operator": "||",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expressions_logic_operator": "||",
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"name":             name,
			// 		"description":      "Acc alarm test",
			// 		"alarm_actions":    []string{"${alibabacloudstack_ess_scaling_rule.default.0.ari}"},
			// 		"scaling_group_id": "${alibabacloudstack_ess_scaling_group.default.id}",
			// 		"metric_type":      "system",
			// 		"evaluation_count": "2",
			// 		"expressions": []map[string]string{
			// 			{
			// 				"period":              "120",
			// 				"statistics":          "Minimum",
			// 				"metric_name":         "MemoryUtilization",
			// 				"threshold":           "99.9",
			// 				"comparison_operator": ">",
			// 			},
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"name":             name,
			// 			"description":      "Acc alarm test",
			// 			"alarm_actions.#":  "1",
			// 			"scaling_group_id": CHECKSET,
			// 			"metric_type":      "system",
			// 			"evaluation_count": "2",
			// 			"expressions.#":    "1",
			// 		}),
			// 	),
			// },
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccalibabacloudstackEssAlarmWithEffective(t *testing.T) {
	var v ess.Alarm
	rand := getAccTestRandInt(10000, 999999)
	var basicMap = map[string]string{
		// "name":                   fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand),
		// "description":            "Acc alarm test",
		// "alarm_actions.#":        "1",
		"scaling_group_id": CHECKSET,
		// "metric_type":            "system",
		// "evaluation_count":       "2",
		"cloud_monitor_group_id": NOSET,
		// "enable":                 "true",
		// "expressions.#":          "1",
	}
	resourceId := "alibabacloudstack_ess_alarm.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssAlarmConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// "name":             name,
					"description":      "Acc alarm test",
					"alarm_actions":    []string{"${alibabacloudstack_ess_scaling_rule.default.0.ari}"},
					"scaling_group_id": "${alibabacloudstack_ess_scaling_group.default.id}",
					"metric_type":      "system",
					"evaluation_count": "2",
					"effective":        "* * * * * ?",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// "name":             name,
						"description":      "Acc alarm test",
						"scaling_group_id": CHECKSET,
						"metric_type":      "system",
						"evaluation_count": "2",
						"alarm_actions.#":  "1",
						"expressions.#":    "1",
						"effective":        "* * * * * ?",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective": "* * 17-18 * * ?",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective": "* * 17-18 * * ?",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccalibabacloudstackEssAlarmWithEffectiveModify(t *testing.T) {
	var v ess.Alarm
	rand := getAccTestRandInt(10000, 999999)
	var basicMap = map[string]string{
		// "name":                   fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand),
		// "description":            "Acc alarm test",
		// "alarm_actions.#":        "1",
		"scaling_group_id": CHECKSET,
		// "metric_type":            "system",
		// "evaluation_count":       "2",
		"cloud_monitor_group_id": NOSET,
		// "enable":                 "true",
		// "expressions.#":          "1",
	}
	resourceId := "alibabacloudstack_ess_alarm.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssAlarmConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// "name":             name,
					"description":      "Acc alarm test",
					"alarm_actions":    []string{"${alibabacloudstack_ess_scaling_rule.default.0.ari}"},
					"scaling_group_id": "${alibabacloudstack_ess_scaling_group.default.id}",
					"metric_type":      "system",
					"evaluation_count": "2",
					"expressions": []map[string]string{{
						"period":              "300",
						"statistics":          "Average",
						"metric_name":         "CpuUtilization",
						"threshold":           "200.3",
						"comparison_operator": ">=",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// "name":             name,
						"description":      "Acc alarm test",
						"scaling_group_id": CHECKSET,
						"metric_type":      "system",
						"evaluation_count": "2",
						"alarm_actions.#":  "1",
						"expressions.#":    "1",
						"effective":        "* * * * * ?",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective": "* * 17-18 * * ?",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective": "* * 17-18 * * ?",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccalibabacloudstackEssAlarmMulti(t *testing.T) {
	var v ess.Alarm
	rand := getAccTestRandInt(100, 999)
	var basicMap = map[string]string{
		// "name":                fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand),
		// "description":         "Acc alarm test",
		// "alarm_actions.#":     "1",
		"scaling_group_id": CHECKSET,
		// "metric_type":         "system",
		// "metric_name":         "CpuUtilization",
		// "period":              "300",
		// "statistics":          "Average",
		// "comparison_operator": ">=",
		// "evaluation_count":    "2",
		// "threshold":           "200.3",
	}
	resourceId := "alibabacloudstack_ess_alarm.default.9"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssAlarmConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count": "10",
					// "name":                name,
					"description":         "Acc alarm test",
					"alarm_actions":       []string{"${alibabacloudstack_ess_scaling_rule.default.0.ari}"},
					"scaling_group_id":    "${alibabacloudstack_ess_scaling_group.default.id}",
					"metric_type":         "system",
					"metric_name":         "CpuUtilization",
					"period":              "300",
					"statistics":          "Average",
					"threshold":           "200.3",
					"comparison_operator": ">=",
					"evaluation_count":    "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceEssAlarmConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	%s

	%s

	resource "alibabacloudstack_vpc_vswitch" "default" {
		name = "${var.name}_vsw"
		vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
		cidr_block = "192.168.0.0/16"
		zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	  }

// 	resource "alibabacloudstack_vpc_vswitch" "default2" {
// 		vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
// 		cidr_block = "192.168.0.0/16"
// 		availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
// 		name = "${var.name}"
//   }


	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}",]
	}

	// resource "alibabacloudstack_ess_scaling_group" "new" {
	// 	min_size = 1
	// 	max_size = 1
	// 	scaling_group_name = "${var.name}-new"
	// 	removal_policies = ["OldestInstance", "NewestInstance"]
	// 	vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default2.id}"]
	// }

	resource "alibabacloudstack_ess_scaling_rule" "default" {
		count = 2
		scaling_rule_name = "${var.name}-${count.index}"
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}

`, name, DataZoneCommonTestCase, VpcCommonTestCase)
}
