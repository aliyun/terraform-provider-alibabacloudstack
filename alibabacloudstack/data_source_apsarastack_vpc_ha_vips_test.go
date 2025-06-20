package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpcHaVipsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_vpc_ha_vips.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sVpcHaVipsDataSource-%d", defaultRegionToTest, rand),
		dataSourceVpcHaVipsDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_vpc_ha_vip.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_vpc_ha_vip.default.description}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_vpc_ha_vip.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_vpc_ha_vip.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vpc_ha_vip.default.description}",
			"ids":        []string{"${alibabacloudstack_vpc_ha_vip.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vpc_ha_vip.default.description}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_vpc_ha_vip.default.id}-fakeTestAcccc"},
		}),
	}

	var existVpcHaVipsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"ids.0":                              CHECKSET,
			"ha_vips.#":                          "1",
			"ha_vips.0.description":              fmt.Sprintf("tf-testAcc%sVpcHaVipsDataSource-%d", defaultRegionToTest, rand),
			"ha_vips.0.ha_vip_name":              fmt.Sprintf("tf-testAcc%sVpcHaVipsDataSource-%d", defaultRegionToTest, rand),
			"ha_vips.0.ha_vip_id":                CHECKSET,
			"ha_vips.0.ip_address":               CHECKSET,
			"ha_vips.0.associated_instance_type": CHECKSET,
			"ha_vips.0.associated_instances.#":   "2",
		}
	}

	var fakeVpcHaVipsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"ha_vips.#": "0",
		}
	}

	var VpcHaVipsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpcHaVipsMapFunc,
		fakeMapFunc:  fakeVpcHaVipsMapFunc,
	}

	VpcHaVipsCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

func dataSourceVpcHaVipsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

%s

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      instance_type,
	  system_disk_category
    ]
  }
}

resource "alibabacloudstack_vpc_ha_vip" "default" {
  ha_vip_name = "${var.name}"
  description = "${var.name}"
  ip_address  = "172.16.1.88"
  vswitch_id  = "${alibabacloudstack_vpc_vswitch.default.id}"
  vpc_id      = "${alibabacloudstack_vpc_vpc.default.id}"
  associated_instance_type = "EcsInstance"
  associated_instances = ["${alibabacloudstack_ecs_instance.default.id}",]
}

 `, name, DataAlibabacloudstackImages, DataAlibabacloudstackInstanceTypes, SecurityGroupCommonTestCase)
}
