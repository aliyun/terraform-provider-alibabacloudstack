package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackDBInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackDBInstanceDataSourceConfig_mysql,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_db_instances.default"),
					resource.TestCheckResourceAttr("data.apsarastack_db_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.apsarastack_db_instances.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackDBInstanceDataSourceConfig_mysql = RdsCommonTestCase + `

variable "name" {
  default = "tf-testAccDBInstanceConfig"
}

variable "creation" {
		default = "Rds"
}


resource "apsarastack_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  instance_name        = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
  storage_type         = "local_ssd"
}
data "apsarastack_db_instances" "default" {
  name_regex = "${apsarastack_db_instance.default.instance_name}"
  ids        = ["${apsarastack_db_instance.default.id}"]
  status     = "Running"
  tags       = {
    "type" = "database",
    "size" = "tiny"
  }
}
`
