package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"
)

func TestAccAlibabacloudStackEdasSwimmingLanesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_edas_swimming_lanes.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%s-%d", defaultRegionToTest, rand),
		dataSourceEdasSwimmingLanesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_edas_swimming_lane.default.name}",
			"group_id":   "${alibabacloudstack_edas_swimming_lane.default.group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_edas_swimming_lane.default.name}-fakeTestAcccc",
			"group_id":   "${alibabacloudstack_edas_swimming_lane.default.group_id}",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":      []string{"${alibabacloudstack_edas_swimming_lane.default.id}"},
			"group_id": "${alibabacloudstack_edas_swimming_lane.default.group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":      []string{"${alibabacloudstack_edas_swimming_lane.default.id}-fakeTestAcccc"},
			"group_id": "${alibabacloudstack_edas_swimming_lane.default.group_id}",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_edas_swimming_lane.default.name}",
			"group_id":   "${alibabacloudstack_edas_swimming_lane.default.group_id}",
			"ids":        []string{"${alibabacloudstack_edas_swimming_lane.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_edas_swimming_lane.default.name}-fakeTestAcccc",
			"group_id":   "${alibabacloudstack_edas_swimming_lane.default.group_id}",
			"ids":        []string{"${alibabacloudstack_edas_swimming_lane.default.id}-fakeTestAcccc"},
		}),
	}

	var existEdasSwimmingLanesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"ids.0":                              CHECKSET,
			"swimming_lanes.#":                   "1",
			"swimming_lanes.0.name":              fmt.Sprintf("tf-testAcc%s-%d", defaultRegionToTest, rand),
			"swimming_lanes.0.condition":         CHECKSET,
			"swimming_lanes.0.apps.#":            "2",
			"swimming_lanes.0.logical_region_id": CHECKSET,
			"swimming_lanes.0.group_id":          CHECKSET,
			"swimming_lanes.0.rest_items.#":      "2",
		}
	}

	var fakeEdasSwimmingLanesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"swimming_lanes.#": "0",
		}
	}

	var EdasSwimmingLanesCheckInfo = dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existEdasSwimmingLanesMapFunc,
		fakeMapFunc:       fakeEdasSwimmingLanesMapFunc,
		ExternalProviders: testAccExternalProviders,
	}

	EdasSwimmingLanesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceEdasSwimmingLanesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

variable "package_version" {	
	default = "2025-10-17 17:17:18"
}

%s

resource "alibabacloudstack_edas_k8s_application" "default" {
  count                   	= 2
  application_name        	= "testapp${count.index}"
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

resource "alibabacloudstack_edas_swimming_lane_group" "default" { 
	name 		  = "${var.name}"
	entry_app_id  = "${alibabacloudstack_edas_k8s_application.default[1].id}"
	apps 		  = ["${alibabacloudstack_edas_k8s_application.default[1].id}", "${alibabacloudstack_edas_k8s_application.default[0].id}"]
	strategy_type = "CONTENT"
}

resource "alibabacloudstack_edas_swimming_lane" "default" { 
	name = "${var.name}"
	group_id = "${alibabacloudstack_edas_swimming_lane_group.default.group_id}"
	apps = ["${alibabacloudstack_edas_k8s_application.default[1].id}", "${alibabacloudstack_edas_k8s_application.default[0].id}"]
	priority = "1"
	path = "/root/test"
	condition = "ADD"
	rest_items {
		name =  "test1"
		value = "100"
		type = "cookie"
		cond = "=="
	}
	rest_items{
		name = "test2"
		value = "50"
		type = "header"
		cond = ">="
		operator = "mod"
	}
}
`, name, EdasClusterCommonTestCase(), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
