package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"
)

func TestAccAlibabacloudStackEdasSwimmingLaneGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_edas_swimming_lane_groups.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%s-%d", defaultRegionToTest, rand),
		dataSourceEdasSwimmingLaneGroupsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_edas_swimming_lane_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_edas_swimming_lane_group.default.name}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_edas_swimming_lane_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_edas_swimming_lane_group.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_edas_swimming_lane_group.default.name}",
			"ids":        []string{"${alibabacloudstack_edas_swimming_lane_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_edas_swimming_lane_group.default.name}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_edas_swimming_lane_group.default.id}-fakeTestAcccc"},
		}),
	}

	var existEdasSwimmingLaneGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    "1",
			"ids.0":                                    CHECKSET,
			"swimming_lane_groups.#":                   "1",
			"swimming_lane_groups.0.name":              fmt.Sprintf("tf-testAcc%s-%d", defaultRegionToTest, rand),
			"swimming_lane_groups.0.entry_app_id":      CHECKSET,
			"swimming_lane_groups.0.apps.#":            "2",
			"swimming_lane_groups.0.logical_region_id": CHECKSET,
			"swimming_lane_groups.0.group_id":          CHECKSET,
		}
	}

	var fakeEdasSwimmingLaneGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "0",
			"swimming_lane_groups.#": "0",
		}
	}

	var EdasSwimmingLaneGroupsCheckInfo = dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existEdasSwimmingLaneGroupsMapFunc,
		fakeMapFunc:       fakeEdasSwimmingLaneGroupsMapFunc,
		ExternalProviders: testAccExternalProviders,
	}

	EdasSwimmingLaneGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceEdasSwimmingLaneGroupsDependence(name string) string {
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

`, name, EdasClusterCommonTestCase(), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
