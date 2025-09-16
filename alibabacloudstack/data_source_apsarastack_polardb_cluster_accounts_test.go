package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbClusterAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_polardb_cluster_accounts.default"
	name := fmt.Sprintf("tfaccount%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, datasourcePolardbClusterAccountsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_polardb_cluster_account.default.id}"},
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_polardb_cluster_account.default.id}_fake"},
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
	}

	account_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_polardb_cluster_account.default.id}"},
			"account_name":  "${alibabacloudstack_polardb_cluster_account.default.account_name}",
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_polardb_cluster_account.default.id}_fake"},
			"account_name":  "${alibabacloudstack_polardb_cluster_account.default.account_name}_fake",
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
	}

	name_regex_Conf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_polardb_cluster_account.default.id}"},
			"account_name":  "${alibabacloudstack_polardb_cluster_account.default.account_name}",
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_polardb_cluster_account.default.id}_fake"},
			"account_name":  "${alibabacloudstack_polardb_cluster_account.default.account_name}_fake",
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
	}

	AlibabacloudstackPolardbClusterAccountsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, account_nameConf, name_regex_Conf)
}

var existAlibabacloudstackPolardbClusterAccountsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#":                     "1",
		"accounts.0.id":                  CHECKSET,
		"accounts.0.account_description": CHECKSET,
		"accounts.0.account_name":        fmt.Sprintf("tfaccount%d", rand),
		"accounts.0.account_type":        CHECKSET,
	}
}

var fakeAlibabacloudstackPolardbClusterAccountsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#": "0",
	}
}

var AlibabacloudstackPolardbClusterAccountsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_polardb_cluster_accounts.default",
	existMapFunc: existAlibabacloudstackPolardbClusterAccountsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackPolardbClusterAccountsDataMapFunc,
}

func datasourcePolardbClusterAccountsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
	variable "creation" {
		default = "PolarDB"
	}
	%s
	%s
	resource "alibabacloudstack_polardb_cluster_instance" "instance" {
		db_cluster_description 	= "${var.name}"
		db_type            		= "MySQL"
		db_version    			= "5.7"
		instance_name 			= "${var.name}"
		storage_type			= "ESSDPL1"
		storage_space 			= 20
		db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
		db_node_num 			= "1"
		zone_id					= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vpc_id 					= "${alibabacloudstack_vpc_vpc.default.id}"
		vswitch_id 				= "${alibabacloudstack_vpc_vswitch.default.id}"
		sub_category 			= "General"
	}
	resource "alibabacloudstack_polardb_cluster_account" "default" {
		dbcluster_id = "${alibabacloudstack_polardb_cluster.default.id}"
		account_name = "test"
		account_description = "from terraform"
		account_password = "${random_password.password.0.result}"
	}
	`, name, RandomPasswordTestCase(12, 1), VSwitchCommonTestCase)
}
