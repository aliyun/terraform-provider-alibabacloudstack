package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEdasK8sApplicationScalingRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_edas_k8s_application_scaling_rule.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeEdasK8sApplicationScalingRule")
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-app-scalingrule%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sApplicationScalingRuleDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":            "${alibabacloudstack_edas_k8s_application.default.id}",
					"scaling_rule_name": "testtf",
					"scaling_rule_type": "metric",
					"max_replicas":      "4",
					"min_replicas":      "1",
					"metrics": []map[string]interface{}{
						{
							"type":        "CPU",
							"utilization": "80",
						},
						{
							"type":        "MEMORY",
							"utilization": "80",
						},
					},
					"scale_up_stabilization_window_seconds": "2",
					"scale_up_select_policy":                "Max",
					"scale_up_policies": []map[string]interface{}{
						{
							"type":           "Percent",
							"value":          "10",
							"period_seconds": "20",
						},
						{
							"type":           "Pods",
							"value":          "10",
							"period_seconds": "20",
						},
					},
					"scale_down_stabilization_window_seconds": "2",
					"scale_down_select_policy":                "Min",
					"scale_down_policies": []map[string]interface{}{
						{
							"type":           "Percent",
							"value":          "10",
							"period_seconds": "20",
						},
						{
							"type":           "Pods",
							"value":          "10",
							"period_seconds": "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_name":                       "testtf",
						"scaling_rule_type":                       "metric",
						"max_replicas":                            "4",
						"min_replicas":                            "1",
						"metrics.#":                               "2",
						"metrics.0.type":                          "CPU",
						"metrics.0.utilization":                   "80",
						"metrics.1.type":                          "MEMORY",
						"metrics.1.utilization":                   "80",
						"scale_up_stabilization_window_seconds":   "2",
						"scale_up_select_policy":                  "Max",
						"scale_up_policies.#":                     "2",
						"scale_up_policies.0.type":                "Percent",
						"scale_down_stabilization_window_seconds": "2",
						"scale_down_select_policy":                "Min",
						"scale_down_policies.#":                   "2",
						"scale_down_policies.0.type":              "Percent",
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"max_replicas": "5",
					"min_replicas": "3",
					"metrics": []map[string]interface{}{
						{
							"type":        "CPU",
							"utilization": "75",
						},
						{
							"type":        "MEMORY",
							"utilization": "75",
						},
					},
					"scale_up_stabilization_window_seconds":   "10",
					"scale_down_stabilization_window_seconds": "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_name":                       "testtf",
						"scaling_rule_type":                       "metric",
						"max_replicas":                            "5",
						"min_replicas":                            "3",
						"metrics.#":                               "2",
						"metrics.0.type":                          "CPU",
						"metrics.0.utilization":                   "75",
						"metrics.1.type":                          "MEMORY",
						"metrics.1.utilization":                   "75",
						"scale_up_stabilization_window_seconds":   "10",
						"scale_down_stabilization_window_seconds": "8",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackEdasK8sApplicationScalingRule_triggerNew(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_edas_k8s_application_scaling_rule.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeEdasK8sApplicationScalingRule")
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tfappscalingrule%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sApplicationScalingRuleDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":            "${alibabacloudstack_edas_k8s_application.default.id}",
					"scaling_rule_name": "${var.name}",
					"scaling_rule_type": "trigger",
					"max_replicas":      "10",
					"min_replicas":      "1",
					"triggers": []map[string]interface{}{
						{
							"name":   "test1",
							"period": "weekly",
							"timer_in_day": []map[string]interface{}{
								{
									"at_time":  "08:00",
									"replicas": "10",
								},
								{
									"at_time":  "20:00",
									"replicas": "4",
								},
							},
							"timer_in_week": []string{
								"Sun",
								"Sat",
								"Fri",
								"Thu",
								"Wed",
								"Tue",
								"Mon",
							},
						},
						{
							"name":   "test2",
							"period": "daily",
							"timer_in_day": []map[string]interface{}{
								{
									"at_time":  "08:00",
									"replicas": "10",
								},
								{
									"at_time":  "20:00",
									"replicas": "4",
								},
							},
						},
					},
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_name":                 name,
						"scaling_rule_type":                 "trigger",
						"max_replicas":                      "10",
						"min_replicas":                      "1",
						"triggers.#":                        "2",
						"triggers.0.name":                   "test1",
						"triggers.0.period":                 "weekly",
						"triggers.0.timer_in_day.#":         "2",
						"triggers.0.timer_in_day.0.at_time": "08:00",
						"triggers.0.timer_in_week.#":        "7",
						"triggers.1.name":                   "test2",
						"triggers.1.period":                 "daily",
						"triggers.1.timer_in_day.#":         "2",
						"triggers.1.timer_in_day.0.at_time": "08:00",
						"enabled":                           "true",
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"max_replicas": "8",
					"min_replicas": "4",
					"triggers": []map[string]interface{}{
						{
							"name":   "test1",
							"period": "monthly",
							"timer_in_day": []map[string]interface{}{
								{
									"at_time":  "09:00",
									"replicas": "10",
								},
								{
									"at_time":  "21:00",
									"replicas": "4",
								},
							},
							"timer_in_month": []string{
								"1",
								"2",
							},
						},
						{
							"name":   "test2",
							"period": "daily",
							"timer_in_day": []map[string]interface{}{
								{
									"at_time":  "10:00",
									"replicas": "10",
								},
								{
									"at_time":  "22:00",
									"replicas": "4",
								},
							},
						},
					},
					"scale_up_stabilization_window_seconds": "2",
					"scale_up_select_policy":                "Max",
					"scale_up_policies": []map[string]interface{}{
						{
							"type":           "Percent",
							"value":          "10",
							"period_seconds": "20",
						},
						{
							"type":           "Pods",
							"value":          "10",
							"period_seconds": "20",
						},
					},
					"scale_down_stabilization_window_seconds": "2",
					"scale_down_select_policy":                "Min",
					"scale_down_policies": []map[string]interface{}{
						{
							"type":           "Percent",
							"value":          "10",
							"period_seconds": "20",
						},
						{
							"type":           "Pods",
							"value":          "10",
							"period_seconds": "20",
						},
					},
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_replicas":                            "8",
						"min_replicas":                            "4",
						"triggers.#":                              "2",
						"triggers.0.name":                         "test1",
						"triggers.0.period":                       "monthly",
						"triggers.0.timer_in_day.#":               "2",
						"triggers.0.timer_in_day.0.at_time":       "09:00",
						"triggers.0.timer_in_month.#":             "2",
						"triggers.0.timer_in_week.#":              "0",
						"triggers.1.name":                         "test2",
						"triggers.1.period":                       "daily",
						"triggers.1.timer_in_day.#":               "2",
						"triggers.1.timer_in_day.0.at_time":       "10:00",
						"enabled":                                 "false",
						"scale_up_stabilization_window_seconds":   "2",
						"scale_up_select_policy":                  "Max",
						"scale_up_policies.#":                     "2",
						"scale_up_policies.0.type":                "Percent",
						"scale_down_stabilization_window_seconds": "2",
						"scale_down_select_policy":                "Min",
						"scale_down_policies.#":                   "2",
						"scale_down_policies.0.type":              "Percent",
					}),
				),
			},
		},
	})
}

func resourceEdasK8sApplicationScalingRuleDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

variable "package_version" {	
	default = "2025-10-17 17:17:18"
} 

resource "alibabacloudstack_edas_k8s_application" "default" {
  application_name        	= "${var.name}"
  application_description 	= "This is description of application"
  cluster_id              	= local.edas_cluster_id
  replicas                	= 2
  package_type 				= "FatJar"
  package_url     			= "http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar"
  package_version 			= var.package_version
  jdk             			= "Open JDK 8"
  limit_mem             	= 1024
  requests_mem          	= 1024
  requests_m_cpu        	= 300
  limit_m_cpu           	= 300
} 
`, name, EdasClusterCommonTestCase(), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
