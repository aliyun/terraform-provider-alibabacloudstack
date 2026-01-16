package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbMasterSlaveServerGroupsDataSource_basic(t *testing.T) {

	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_master_slave_server_groups.default"
	name := fmt.Sprintf("tf-testaccslbmasterslaveservergroups%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSlbMasterSlaveServerGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_master_slave_server_group.default.load_balancer_id}",
			"name_regex":       "${alibabacloudstack_slb_master_slave_server_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_master_slave_server_group.default.load_balancer_id}",
			"name_regex":       "fake_*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_master_slave_server_group.default.load_balancer_id}",
			"ids":              []string{"${alibabacloudstack_slb_master_slave_server_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_master_slave_server_group.default.load_balancer_id}",
			"ids":              []string{"fake-mssg-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_master_slave_server_group.default.load_balancer_id}",
			"ids":              []string{"${alibabacloudstack_slb_master_slave_server_group.default.id}"},
			"name_regex":       "${alibabacloudstack_slb_master_slave_server_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_master_slave_server_group.default.load_balancer_id}",
			"ids":              []string{"fake-mssg-id"},
			"name_regex":       "fake_*",
		}),
	}

	var existSlbMasterSlaveServerGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"groups.#":                       "1",
			"groups.0.id":                    CHECKSET,
			"groups.0.name":                  name,
			"groups.0.servers.0.instance_id": CHECKSET,
			"groups.0.servers.0.weight":      CHECKSET,
			"groups.0.servers.0.port":        CHECKSET,
			"groups.0.servers.0.server_type": CHECKSET,
		}
	}

	var fakeSlbMasterSlaveServerGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var slbMasterSlaveServerGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbMasterSlaveServerGroupsMapFunc,
		fakeMapFunc:  fakeSlbMasterSlaveServerGroupsMapFunc,
	}
	slbMasterSlaveServerGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceSlbMasterSlaveServerGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
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

resource "alibabacloudstack_instance" "master" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  instance_type = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  security_groups = ["${alibabacloudstack_security_group.default.id}"]
  instance_name = "${var.name}-master"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_instance" "slave" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  instance_type = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  security_groups = ["${alibabacloudstack_security_group.default.id}"]
  instance_name = "${var.name}-slave"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_master_slave_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  name = "${var.name}"
  servers {
      server_id = "${alibabacloudstack_instance.master.id}"
      port = 80
      weight = 100
      server_type = "Master"
  }
  servers {
      server_id = "${alibabacloudstack_instance.slave.id}"
      port = 80
      weight = 100
      server_type = "Slave"
  }
}

`, name, ECSInstanceCommonTestCase)
}
