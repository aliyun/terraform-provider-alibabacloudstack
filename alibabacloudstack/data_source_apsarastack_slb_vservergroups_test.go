package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbServerGroupsDataSource_basic(t *testing.T) {

	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_server_groups.default"
	name := fmt.Sprintf("tf-testaccslbservergroups%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSlbServerGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_server_group.default.load_balancer_id}",
			"name_regex":       "${alibabacloudstack_slb_server_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_server_group.default.load_balancer_id}",
			"name_regex":       "fake_*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_server_group.default.load_balancer_id}",
			"ids":              []string{"${alibabacloudstack_slb_server_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_server_group.default.load_balancer_id}",
			"ids":              []string{"fake-sg-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_server_group.default.load_balancer_id}",
			"ids":              []string{"${alibabacloudstack_slb_server_group.default.id}"},
			"name_regex":       "${alibabacloudstack_slb_server_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_server_group.default.load_balancer_id}",
			"ids":              []string{"fake-sg-id"},
			"name_regex":       "fake_*",
		}),
	}

	var existSlbServerGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"slb_server_groups.#":           "1",
			"slb_server_groups.0.id":        CHECKSET,
			"slb_server_groups.0.name":      name,
			"slb_server_groups.0.servers.#": "1",
			"slb_server_groups.0.servers.0.instance_id": CHECKSET,
			"slb_server_groups.0.servers.0.port":        CHECKSET,
			"slb_server_groups.0.servers.0.weight":      CHECKSET,
		}
	}

	var fakeSlbServerGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "0",
			"names.#":             "0",
			"slb_server_groups.#": "0",
		}
	}

	var slbServerGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbServerGroupsMapFunc,
		fakeMapFunc:  fakeSlbServerGroupsMapFunc,
	}
	slbServerGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceSlbServerGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
%s

%s

%s

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  instance_type = "${local.default_instance_type_id}"
  system_disk_category = "cloud_sperf"
  security_groups = ["${alibabacloudstack_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  name = "${var.name}"
  servers {
      server_ids = ["${alibabacloudstack_instance.default.id}"]
      port = 80
      weight = 100
    }
}

`, name, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages)
}
