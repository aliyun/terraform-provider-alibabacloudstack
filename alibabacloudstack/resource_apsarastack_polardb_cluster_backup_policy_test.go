package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccAlibabacloudStackPolardbClusterBackupPolicy_basic tests basic creation and update of backup policy
func TestAccAlibabacloudStackPolardbClusterBackupPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_polardb_cluster_backup_policy.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"data_level1_backup_frequency":                CHECKSET,
		"backup_retention_policy_on_cluster_deletion": CHECKSET,
	})
	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolardbClusterBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc-polardb-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbClusterBackupPolicyDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":                       "${alibabacloudstack_polardb_cluster_instance.default.id}",
					"data_level1_backup_period":           "Monday,Tuesday,Wednesday,Thursday,Friday",
					"data_level1_backup_time":             "10:00Z-11:00Z",
					"data_level1_backup_retention_period": "7",
					"log_backup_retention_period":         "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_level1_backup_period":           "Monday,Tuesday,Wednesday,Thursday,Friday",
						"data_level1_backup_time":             "10:00Z-11:00Z",
						"data_level1_backup_retention_period": "7",
						"log_backup_retention_period":         "7",
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
					"data_level1_backup_period":           "Monday,Tuesday,Wednesday,Thursday,Friday",
					"data_level1_backup_time":             "12:00Z-13:00Z",
					"data_level1_backup_retention_period": "3",
					"log_backup_retention_period":         "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_level1_backup_period":           "Monday,Tuesday,Wednesday,Thursday,Friday",
						"data_level1_backup_time":             "12:00Z-13:00Z",
						"data_level1_backup_retention_period": "3",
						"log_backup_retention_period":         "5",
					}),
				),
			},
		},
	})
}

func resourcePolardbClusterBackupPolicyDependence(name string) string {
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

	 resource "alibabacloudstack_polardb_cluster_instance" "default" {
	 	db_cluster_description 	= "${var.name}"
	 	db_type            		= "${var.db_type}"
	 	db_version    			= "${var.db_version}"
	 	storage_type			= "ESSDPL1"
	 	storage_space 			= 20
	 	db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
	 	zone_id					= "${data.alibabacloudstack_zones.default.zones.0.id}"
	 	vswitch_id 				= "${alibabacloudstack_vpc_vswitch.default.id}"
	 	sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
	 }
	  `, name, VSwitchCommonTestCase)
}
