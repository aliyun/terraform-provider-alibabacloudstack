package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackDRDSInstancesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.apsarastack_drds_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sDRDSInstancesDataSource-%d", defaultRegionToTest, rand),
		dataSourceDRDSInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${apsarastack_drds_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${apsarastack_drds_instance.default.description}-fake",
		}),
	}

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${apsarastack_drds_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${apsarastack_drds_instance.default.description}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${apsarastack_drds_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${apsarastack_drds_instance.default.id}-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${apsarastack_drds_instance.default.description}",
			"description_regex": "${apsarastack_drds_instance.default.description}",
			"ids":               []string{"${apsarastack_drds_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${apsarastack_drds_instance.default.description}-fake",
			"description_regex": "${apsarastack_drds_instance.default.description}-fake",
			"ids":               []string{"${apsarastack_drds_instance.default.id}-fake"},
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

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
	}

	drdsInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, descriptionRegexConf, idsConf, allConf)
}

func dataSourceDRDSInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
 	data "apsarastack_zones" "default" {
		available_resource_creation = "VSwitch"
	}
 	variable "name" {
		default = "%s"
	}
	variable "instance_series" {
		default = "drds.sn2.4c16g"
	}
	
	resource "apsarastack_vpc" "default" {
	  name       = "${var.name}"
	  cidr_block = "172.16.0.0/16"
	}
	resource "apsarastack_vswitch" "default" {
	  vpc_id            = "${apsarastack_vpc.default.id}"
	  cidr_block        = "172.16.0.0/24"
	  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	  name              = "${var.name}"
	}	
 	resource "apsarastack_drds_instance" "default" {
  		description = "${var.name}"
  		zone_id = "${apsarastack_vswitch.default.availability_zone}"
  		instance_series = "${var.instance_series}"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${apsarastack_vswitch.default.id}"
  		specification = "drds.sn2.4c16g.8C32G"
}
 `, name)
}
