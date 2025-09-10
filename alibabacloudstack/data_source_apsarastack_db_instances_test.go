package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDBInstancesDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackDBInstanceDataSourceConfig_mysql(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_db_instances.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_db_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_db_instances.default", "ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckAlibabacloudStackDBInstanceDataSourceConfig_mysql() string {
	return `

variable "name" {
  default = "tf-testAccDBInstanceConfig"
}

variable "creation" {
		default = "Rds"
}
` + VSwitchCommonTestCase + RdsMysqlCommonTestCase() + `

data "alibabacloudstack_db_instances" "default" {
  name_regex = "${alibabacloudstack_db_instance.default.instance_name}"
  ids        = ["${alibabacloudstack_db_instance.default.id}"]
  status     = "Running"
}
`
}
