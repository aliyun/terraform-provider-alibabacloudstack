package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbxAccountDatabaseBinding_basic0(t *testing.T) {
	var v []map[string]string

	resourceId := "alibabacloudstack_polardbx_account_database_binding.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolardbXAccountDBPrivilege")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("accdbbind%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxAccountDatabaseBindingDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alibabacloudstack_polardbx_instance.default.id}",
					"account_name": "${alibabacloudstack_polardbx_account.default.account_name}",
					"db_privileges": []map[string]interface{}{
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default0.database_name}",
							"privilege": "ReadOnly",
						},
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default1.database_name}",
							"privilege": "ReadWrite",
						},
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default2.database_name}",
							"privilege": "DDLOnly",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_privileges.#":           "3",
						"db_privileges.0.db_name":   fmt.Sprintf("%s0", name),
						"db_privileges.0.privilege": "ReadOnly",
						"db_privileges.1.db_name":   fmt.Sprintf("%s1", name),
						"db_privileges.1.privilege": "ReadWrite",
						"db_privileges.2.db_name":   fmt.Sprintf("%s2", name),
						"db_privileges.2.privilege": "DDLOnly",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_privileges": []map[string]interface{}{
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default0.database_name}",
							"privilege": "DDLOnly",
						},
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default1.database_name}",
							"privilege": "ReadOnly",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_privileges.#":           "2",
						"db_privileges.0.db_name":   fmt.Sprintf("%s0", name),
						"db_privileges.0.privilege": "DDLOnly",
						"db_privileges.1.db_name":   fmt.Sprintf("%s1", name),
						"db_privileges.1.privilege": "ReadOnly",
						"db_privileges.2.db_name":   REMOVEKEY,
						"db_privileges.2.privilege": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func resourcePolardbxAccountDatabaseBindingDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "password" {
  default = "%s"
}
%s
resource "alibabacloudstack_polardbx_instance" "default" {
    description = "testtf1111"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage = "50"
	network_type = "vpc"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class = "polarx.x4.medium.2e"
	cn_node_count = "2"
	dn_node_class = "mysql.n4.medium.25"
	dn_node_count = "2"
}

resource "alibabacloudstack_polardbx_database" "default0" {
    instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	database_name = "${var.name}0"
	encode = "utf8mb4"
	mode = "auto"
}

resource "alibabacloudstack_polardbx_database" "default1" {
    instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	database_name = "${var.name}1"
	encode = "utf8mb4"
	mode = "auto"
}

resource "alibabacloudstack_polardbx_database" "default2" {
    instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	database_name = "${var.name}2"
	encode = "utf8mb4"
	mode = "auto"
}

resource "alibabacloudstack_polardbx_account" "default" {
    instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	account_name = "${var.name}"
	account_type = "Normal"
	password     = "${var.password}"
	description  = "Normal user"
}

 `, name, getAccTestPassword(12), VSwitchCommonTestCase)
}
