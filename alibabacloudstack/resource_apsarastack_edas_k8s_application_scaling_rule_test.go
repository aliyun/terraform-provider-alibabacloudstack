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

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":            "${alibabacloudstack_edas_k8s_application.default.id}",
					"scaling_rule_name": "testtf",
					"scaling_rule_type": "metric",
					"max_replicas":      "10",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_name":     "testtf",
						"scaling_rule_type":     "metric",
						"max_replicas":          "10",
						"min_replicas":          "1",
						"metrics.#":             "2",
						"metrics.0.type":        "CPU",
						"metrics.0.utilization": "80",
						"metrics.1.type":        "MEMORY",
						"metrics.1.utilization": "80",
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
					"metrics": []map[string]interface{}{
						{
							"type":        "CPU",
							"utilization": "70",
						},
						{
							"type":        "MEMORY",
							"utilization": "70",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_name":     "testtf",
						"scaling_rule_type":     "metric",
						"max_replicas":          "8",
						"min_replicas":          "4",
						"metrics.#":             "2",
						"metrics.0.type":        "CPU",
						"metrics.0.utilization": "70",
						"metrics.1.type":        "MEMORY",
						"metrics.1.utilization": "70",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackEdasK8sApplicationScalingRule_trigger(t *testing.T) {
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

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":            "${alibabacloudstack_edas_k8s_application.default.id}",
					"scaling_rule_name": "${var.name}",
					"scaling_rule_type": "trigger",
					"max_replicas":      "10",
					"min_replicas":      "1",
					"trigger_name":      "${var.name}",
					"trigger_period":    "weekly",
					"trigger_dryrun":    "true",
					"trigger_timer_in_day": []map[string]interface{}{
						{
							"at_time":  "08:00",
							"replicas": "10",
						},
						{
							"at_time":  "20:00",
							"replicas": "4",
						},
					},
					"trigger_timer_in_week": []string{"Sun", "Sat", "Fri", "Thu", "Wed", "Tue", "Mon"},
					"enabled":               "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_rule_name":               name,
						"scaling_rule_type":               "trigger",
						"max_replicas":                    "10",
						"min_replicas":                    "1",
						"trigger_name":                    name,
						"trigger_period":                  "weekly",
						"trigger_dryrun":                  "true",
						"trigger_timer_in_day.#":          "2",
						"trigger_timer_in_day.0.at_time":  "08:00",
						"trigger_timer_in_day.0.replicas": "10",
						"trigger_timer_in_day.1.at_time":  "20:00",
						"trigger_timer_in_day.1.replicas": "4",
						"trigger_timer_in_week.#":         "7",
						"enabled":                         "true",
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
					"max_replicas":   "8",
					"min_replicas":   "4",
					"trigger_dryrun": "true",
					"trigger_timer_in_day": []map[string]interface{}{
						{
							"at_time":  "09:00",
							"replicas": "10",
						},
						{
							"at_time":  "21:00",
							"replicas": "4",
						},
					},
					"trigger_timer_in_week": []string{"Sun", "Sat", "Fri", "Thu", "Wed"},
					"enabled":               "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_replicas":                    "8",
						"min_replicas":                    "4",
						"trigger_dryrun":                  "true",
						"trigger_timer_in_day.#":          "2",
						"trigger_timer_in_day.0.at_time":  "09:00",
						"trigger_timer_in_day.0.replicas": "10",
						"trigger_timer_in_day.1.at_time":  "21:00",
						"trigger_timer_in_day.1.replicas": "4",
						"trigger_timer_in_week.#":         "5",
						"enabled":                         "false",
					}),
				),
			},
		},
	})
}

func resourceEdasK8sApplicationScalingRuleDependence(name string) string {
	edasClusterId := os.Getenv("ALIBABACLOUDSTACK_EDAS_CLUSTER_ID")
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

variable "cluster_id" {	
	default = "%s"
}


variable "package_version" {	
	default = "2025-10-09 17:17:18"
} 

resource "alibabacloudstack_edas_k8s_application" "default" {
  application_name        	= "${var.name}"
  application_description 	= "This is description of application"
  cluster_id              	= var.cluster_id
  replicas                	= 2
  package_type 				= "FatJar"
  package_url     			= "http://fileserver.edas.intra.env212.shuguang.com//prod/demo/SPRING_CLOUD_PROVIDER.jar"
  package_version 			= var.package_version
  jdk             			= "Open JDK 8"
  limit_mem             	= 1024
  requests_mem          	= 1024
  requests_m_cpu        	= 300
  limit_m_cpu           	= 300
}
`, name, edasClusterId)
}
