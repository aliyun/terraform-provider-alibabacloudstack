package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

const testAccCheckAlibabacloudStackSlbMasterSlaveServerGroupsDataSourceBasic = ECSInstanceCommonTestCase + `
variable "name" {
  default = "tf-testAccslbmasterslaveservergroupsdatasourcebasic"
}

resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}

resource "alibabacloudstack_ecs_instance" "new" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 40
  system_disk_name     = "test_sys_diskv2"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}


resource "alibabacloudstack_slb_master_slave_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  name = "${var.name}"
  servers {
      server_id = "${alibabacloudstack_ecs_instance.default.id}"
      port = 80
      weight = 100
      server_type = "Master"
  }
  servers {
      server_id = "${alibabacloudstack_ecs_instance.new.id}"
      port = 80
      weight = 100
      server_type = "Slave"
  }
}

data "alibabacloudstack_slb_master_slave_server_groups" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
}`
