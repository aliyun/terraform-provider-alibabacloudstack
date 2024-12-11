package alibabacloudstack

import (
	"fmt"
	"testing"

	
)

func TestAccAlibabacloudStackDRDSInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_drds_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sDRDSInstancesDataSource-%d", defaultRegionToTest, rand),
		dataSourceDRDSInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_drds_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_drds_instance.default.description}-fake",
		}),
	}

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_drds_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_drds_instance.default.description}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_drds_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_drds_instance.default.id}-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_drds_instance.default.description}",
			"description_regex": "${alibabacloudstack_drds_instance.default.description}",
			"ids":               []string{"${alibabacloudstack_drds_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_drds_instance.default.description}-fake",
			"description_regex": "${alibabacloudstack_drds_instance.default.description}-fake",
			"ids":               []string{"${alibabacloudstack_drds_instance.default.id}-fake"},
		}),
	}

	var existDRDSInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"descriptions.#":           "1",
			"ids.0":                    CHECKSET,
			"descriptions.0":           fmt.Sprintf("tf-testAcc%sDRDSInstancesDataSource-%d", defaultRegionToTest, rand),
			"instances.#":              "1",
			"instances.0.description":  fmt.Sprintf("tf-testAcc%sDRDSInstancesDataSource-%d", defaultRegionToTest, rand),
			"instances.0.type":         "PRIVATE",
			"instances.0.zone_id":      CHECKSET,
			"instances.0.id":           CHECKSET,
			"instances.0.network_type": "VPC",
			"instances.0.create_time":  CHECKSET,
		}
	}

	var fakeDRDSInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"descriptions.#": "0",
			"instances.#":    "0",
		}
	}

	var drdsInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDRDSInstancesMapFunc,
		fakeMapFunc:  fakeDRDSInstancesMapFunc,
	}

	drdsInstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, descriptionRegexConf, idsConf, allConf)
}

func dataSourceDRDSInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
 	data "alibabacloudstack_zones" "default" {
		available_resource_creation = "VSwitch"
	}
 	variable "name" {
		default = "%s"
	}
	variable "instance_series" {
		default = "drds.sn2.4c16g"
	}
	
	resource "alibabacloudstack_vpc" "default" {
	  name       = "${var.name}"
	  cidr_block = "172.16.0.0/16"
	}
	resource "alibabacloudstack_vswitch" "default" {
	  vpc_id            = "${alibabacloudstack_vpc.default.id}"
	  cidr_block        = "172.16.0.0/24"
	  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	  name              = "${var.name}"
	}	
 	resource "alibabacloudstack_drds_instance" "default" {
  		description = "${var.name}"
  		zone_id = "${alibabacloudstack_vswitch.default.availability_zone}"
  		instance_series = "${var.instance_series}"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
  		specification = "drds.sn2.4c16g.8C32G"
}
 `, name)
}
