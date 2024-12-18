package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackDBInstancesDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackDBInstanceDataSourceConfig_mysql,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_db_instances.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_db_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_db_instances.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackDBInstanceDataSourceConfig_mysql = RdsCommonTestCase + `

variable "name" {
  default = "tf-testAccDBInstanceConfig"
}

variable "creation" {
		default = "Rds"
}


resource "alibabacloudstack_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  instance_name        = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
  storage_type         = "local_ssd"
}
data "alibabacloudstack_db_instances" "default" {
  name_regex = "${alibabacloudstack_db_instance.default.instance_name}"
  ids        = ["${alibabacloudstack_db_instance.default.id}"]
  status     = "Running"
  tags       = {
    "type" = "database",
    "size" = "tiny"
  }
}
`
