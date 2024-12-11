package alibabacloudstack

import (
	"fmt"
	"testing"

	
)

func TestAccAlibabacloudStackEdasApplicationsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_edas_applications.default"
	name := fmt.Sprintf("tf-testacc-edas-applications%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEdasApplicationConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_edas_application.default.application_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_edas_application.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_edas_application.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_edas_application.default.id}"},
			"name_regex": "${alibabacloudstack_edas_application.default.application_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_edas_application.default.id}_fake"},
			"name_regex": "${alibabacloudstack_edas_application.default.application_name}",
		}),
	}

	var existEdasApplicationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"applications.#":                  "1",
			"applications.0.app_name":         name,
			"applications.0.app_id":           CHECKSET,
			"applications.0.application_type": "War",
			"applications.0.build_package_id": CHECKSET,
		}
	}

	var fakeEdasApplicationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"applications.#": "0",
		}
	}

	var edasApplicationCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEdasApplicationsMapFunc,
		fakeMapFunc:  fakeEdasApplicationsMapFunc,
	}

	edasApplicationCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceEdasApplicationConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}

		resource "alibabacloudstack_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = "${var.name}"
		}

		resource "alibabacloudstack_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = "${alibabacloudstack_vpc.default.id}"
         // region_id    = "cn-neimeng-env30-d01"
		}

		resource "alibabacloudstack_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = alibabacloudstack_edas_cluster.default.id
		  package_type = "WAR"
		  build_pack_id = "15"
		}
		`, name)
}
