package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackSlbMasterSlaveServerGroupsDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSlbMasterSlaveServerGroupsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_slb_master_slave_server_groups.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_slb_master_slave_server_groups.default", "load_balancer_id.#", "0"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccCheckAlibabacloudStackSlbMasterSlaveServerGroupsDataSourceBasic = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes + DataAlibabacloudstackImages + `
variable "name" {
  default = "tf-testAccslbmasterslaveservergroupsdatasourcebasic"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
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
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  security_groups = ["${alibabacloudstack_security_group.default.id}"]
  count                      = "2"
  instance_type              = "${local.instance_type_id}"
  image_id                   =  "${data.alibabacloudstack_images.default.images.0.id}"
  instance_name              = "${var.name}"
  vswitch_id                 = alibabacloudstack_vswitch.default.id
  internet_max_bandwidth_out = 10
}

resource "alibabacloudstack_slb_master_slave_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  name = "${var.name}"
  servers {
      server_id = "${alibabacloudstack_instance.default.0.id}"
      port = 80
      weight = 100
      server_type = "Master"
  }
  servers {
      server_id = "${alibabacloudstack_instance.default.1.id}"
      port = 80
      weight = 100
      server_type = "Slave"
  }
}

data "alibabacloudstack_slb_master_slave_server_groups" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
}`
