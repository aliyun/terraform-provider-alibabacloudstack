package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbClusterBackupPoliciesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%v", rand)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourcePolardbClusterBackupPoliciesDataSourceDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_polardb_cluster_backup_policies.policy"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardb_cluster_backup_policies.policy", "db_cluster_id"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardb_cluster_backup_policies.policy", "data_level1_backup_time"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardb_cluster_backup_policies.policy", "data_level1_backup_retention_period"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardb_cluster_backup_policies.policy", "data_level2_backup_retention_period"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardb_cluster_backup_policies.policy", "backup_retention_policy_on_cluster_deletion"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardb_cluster_backup_policies.policy", "log_backup_retention_period"),
				),
			},
		},
	})
}

func datasourcePolardbClusterBackupPoliciesDataSourceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
}

variable "db_type" {
  default = "MySQL"
}

variable "db_version" {
  default = "8.0"
}

%s

data "alibabacloudstack_polardb_cluster_instance_types" "default" {
  db_type = "${var.db_type}"
  db_version = "${var.db_version}"
  sorted_by = "CPU"
  cpu_type = "hygon"
  sub_category = "normal_exclusive"
}

resource "alibabacloudstack_polardb_cluster_instance" "default" {
	db_cluster_description 	=  "${var.name}"
	zone_id 				= "${data.alibabacloudstack_zones.default.zones.0.id}"
	db_type 				= "${var.db_type}"
	db_version 				= "${var.db_version}"
	storage_space 			= "20"
	vswitch_id				= "${alibabacloudstack_vpc_vswitch.default.id}"
	db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
	sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
	storage_type 			= "ESSDPL1"
}

data "alibabacloudstack_polardb_cluster_backup_policies" "policy" {
	db_cluster_id = "${alibabacloudstack_polardb_cluster_instance.default.id}"
}

 `, name, VSwitchCommonTestCase)
}
