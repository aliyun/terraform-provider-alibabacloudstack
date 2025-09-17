package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbClusterAccountDatabaseBinding_basic0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_polardb_cluster_account_database_binding.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolardbClusterAccount")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("accdbbind%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbClusterAccountDatabaseBindingDependence)

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
					"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
					"account_name":  "${alibabacloudstack_polardb_cluster_account.default.account_name}",
					"database_privileges": []map[string]interface{}{
						{
							"db_name":   "${alibabacloudstack_polardb_cluster_database.default0.db_name}",
							"privilege": "ReadOnly",
						},
						{
							"db_name":   "${alibabacloudstack_polardb_cluster_database.default1.db_name}",
							"privilege": "ReadWrite",
						},
						{
							"db_name":   "${alibabacloudstack_polardb_cluster_database.default2.db_name}",
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
							"db_name":   "${alibabacloudstack_polardb_cluster_database.default0.db_name}",
							"privilege": "DDLOnly",
						},
						{
							"db_name":   "${alibabacloudstack_polardb_cluster_database.default1.db_name}",
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

func resourcePolardbClusterAccountDatabaseBindingDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
variable "creation" {
	default = "PolarDB"
}
	
variable "db_type" {
	default = "MySQL"
}

variable "db_version" {
	default = "5.7"
}

data "alibabacloudstack_polardb_cluster_instance_types" "default" {
	db_type = "${var.db_type}"
	db_version = "${var.db_version}"
	sorted_by = "CPU"
	sub_category = "normal_exclusive"
}

%s
%s

resource "alibabacloudstack_polardb_cluster_instance" "instance" {
	db_cluster_description 	= "${var.name}"
	db_type            		= "${var.db_type}"
	db_version    			= "${var.db_version}"
	storage_type			= "ESSDPL1"
	storage_space 			= 20
	db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
	zone_id					= "${data.alibabacloudstack_zones.default.zones.0.id}"
	vpc_id 					= "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id 				= "${alibabacloudstack_vpc_vswitch.default.id}"
	sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
}

resource "alibabacloudstack_polardb_cluster_database" "default0" {
    db_cluster_id  = "${alibabacloudstack_polardb_cluster_instance.instance.id}"
	db_name = "${var.name}0"
	character_set_name = "utf8mb4"
}

resource "alibabacloudstack_polardb_cluster_database" "default1" {
    db_cluster_id  = "${alibabacloudstack_polardb_cluster_instance.instance.id}"
	db_name = "${var.name}1"
	character_set_name = "utf8mb4"
}

resource "alibabacloudstack_polardb_cluster_database" "default2" {
    db_cluster_id  = "${alibabacloudstack_polardb_cluster_instance.instance.id}"
	db_name = "${var.name}2"
	character_set_name = "utf8mb4"
}

resource "alibabacloudstack_polardb_cluster_account" "default" {
    db_cluster_id  = "${alibabacloudstack_polardb_cluster_instance.instance.id}"
	account_name = "${var.name}"
	account_password     = "${random_password.password.0.result}"
}

 `, name, VSwitchCommonTestCase, RandomPasswordTestCase(10, 1))
}
