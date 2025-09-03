package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbAccountDatabaseBinding_basic0(t *testing.T) {
	var v *PolardbDescribeaccountsResponse

	resourceId := "alibabacloudstack_polardb_account_database_binding.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccount")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("accdbbind%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbAccountDatabaseBindingDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id": "${alibabacloudstack_polardb_dbinstance.instance.id}",
					"account_name":   "${alibabacloudstack_polardb_account.default.account_name}",
					"database_privileges": []map[string]interface{}{
						{
							"db_name":   "${alibabacloudstack_polardb_database.default0.data_base_name}",
							"privilege": "ReadOnly",
						},
						{
							"db_name":   "${alibabacloudstack_polardb_database.default1.data_base_name}",
							"privilege": "ReadWrite",
						},
						{
							"db_name":   "${alibabacloudstack_polardb_database.default2.data_base_name}",
							"privilege": "DDLOnly",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"database_privileges.#":           "3",
						"database_privileges.0.db_name":   fmt.Sprintf("%s0", name),
						"database_privileges.0.privilege": "ReadOnly",
						"database_privileges.1.db_name":   fmt.Sprintf("%s1", name),
						"database_privileges.1.privilege": "ReadWrite",
						"database_privileges.2.db_name":   fmt.Sprintf("%s2", name),
						"database_privileges.2.privilege": "DDLOnly",
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
					"database_privileges": []map[string]interface{}{
						{
							"db_name":   "${alibabacloudstack_polardb_database.default0.data_base_name}",
							"privilege": "DDLOnly",
						},
						{
							"db_name":   "${alibabacloudstack_polardb_database.default1.data_base_name}",
							"privilege": "ReadOnly",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"database_privileges.#":           "2",
						"database_privileges.0.db_name":   fmt.Sprintf("%s0", name),
						"database_privileges.0.privilege": "DDLOnly",
						"database_privileges.1.db_name":   fmt.Sprintf("%s1", name),
						"database_privileges.1.privilege": "ReadOnly",
						"database_privileges.2.db_name":   REMOVEKEY,
						"database_privileges.2.privilege": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func resourcePolardbAccountDatabaseBindingDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

variable "creation" {
	default = "PolarDB"
}
resource "alibabacloudstack_polardb_dbinstance" "instance" {
	engine            = "MySQL"
	engine_version    = "5.7"
	instance_name = "${var.name}"
	db_instance_storage_type= "local_ssd"
	db_instance_storage = 5
	db_instance_class = "rds.mysql.t1.small"
	zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}

resource "alibabacloudstack_polardb_database" "default0" {
    data_base_instance_id  = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	data_base_name = "${var.name}0"
	character_set_name = "utf8mb4"
}

resource "alibabacloudstack_polardb_database" "default1" {
    data_base_instance_id  = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	data_base_name = "${var.name}1"
	character_set_name = "utf8mb4"
}

resource "alibabacloudstack_polardb_database" "default2" {
    data_base_instance_id  = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	data_base_name = "${var.name}2"
	character_set_name = "utf8mb4"
}

resource "alibabacloudstack_polardb_account" "default" {
    data_base_instance_id  = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	account_name = "${var.name}"
	account_password     = "Test@1234"
}

 `, name, VSwitchCommonTestCase)
}
