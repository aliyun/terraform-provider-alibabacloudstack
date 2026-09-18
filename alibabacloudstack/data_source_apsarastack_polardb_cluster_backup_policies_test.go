package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbClusterBackupPoliciesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_polardb_cluster_backup_policies.policy"
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbClusterBackupPoliciesConfigDependence)

	// Basic configuration - create a PolarDB cluster and query its backup policies
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.default.id}",
		}),
	}

	var existPolardbBackupPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"db_cluster_id":                               CHECKSET,
			"data_level1_backup_time":                     CHECKSET,
			"data_level1_backup_retention_period":         CHECKSET,
			"data_level2_backup_retention_period":         CHECKSET,
			"backup_retention_policy_on_cluster_deletion": CHECKSET,
			"log_backup_retention_period":                 CHECKSET,
		}
	}

	var fakePolardbBackupPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{}
	}

	var polardbBackupPoliciesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbBackupPoliciesMapFunc,
		fakeMapFunc:  fakePolardbBackupPoliciesMapFunc,
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckApiAvailable(t, "GET", "polardb", "2017-08-01", "DescribeBackupPolicy", "")
		},
	}
	polardbBackupPoliciesCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func dataSourcePolardbClusterBackupPoliciesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
	db_cluster_description 	=  "%s"
	zone_id 				= "${data.alibabacloudstack_zones.default.zones.0.id}"
	db_type 				= "${var.db_type}"
	db_version 				= "${var.db_version}"
	storage_space 			= "20"
	vswitch_id				= "${alibabacloudstack_vpc_vswitch.default.id}"
	db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
	sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
	storage_type 			= "ESSDPL1"
}
`, name, VSwitchCommonTestCase, name)
}
