package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccCheckAlibabacloudStackPolardbInstanceDataSourceConfig_mysql string = RdsCommonTestCase +
	fmt.Sprintf(
		`
variable "name" {
  default = "tf-testAccDBInstanceConfig"
}

variable "creation" {
		default = "PolarDB"
}

resource "alibabacloudstack_polardb_dbinstance" "default" {
	engine            = "MySQL"
	engine_version    = "5.7"
	instance_name = "${var.name}"
	db_instance_storage_type= "local_ssd"
	db_instance_storage = 5
	db_instance_class = "rds.mysql.t1.small"
	zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

data "alibabacloudstack_polardb_dbinstances" "default" {
  db_instance_id        = "${alibabacloudstack_polardb_dbinstance.default.id}"
  db_instance_class = "${alibabacloudstack_polardb_dbinstance.default.db_instance_class}"
  status     = "Running"
  region_id  = "%s"
}`, os.Getenv("ALIBABACLOUDSTACK_REGION"))

func TestAccAlibabacloudStackPolardbInstancesDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackPolardbInstanceDataSourceConfig_mysql,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_polardb_dbinstances.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_polardb_dbinstances.default", "db_instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardb_dbinstances.default", "ids.#"),
				),
			},
		},
	})
}
