package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbBackup_MySQL(t *testing.T) {
	var v *PolardbDescribebackuppolicyResponse
	resourceId := "alibabacloudstack_polardb_backuppolicy.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackPolardbBackupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoPolardbDescribebackuppolicyRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-polardb-backup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackPolardbBackupBasicDependenceMySQL)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":              "${alibabacloudstack_polardb_dbinstance.instance.id}",
					"backup_retention_period":     "9",
					"preferred_backup_time":       "10:00Z-11:00Z", // UTC Time
					"preferred_backup_period":     "Saturday,Sunday",
					"log_backup_retention_period": "7",
					"backup_log":                  "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudStackPolardbBackupMap0 = map[string]string{}

func AlibabacloudStackPolardbBackupBasicDependenceMySQL(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%v"
	}
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
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}
	`, RdsCommonTestCase, name)
}

func TestAccAlibabacloudStackPolardbBackup_PostgreSQL(t *testing.T) {
	var v *PolardbDescribebackuppolicyResponse
	resourceId := "alibabacloudstack_polardb_backuppolicy.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackPolardbBackupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoPolardbDescribebackuppolicyRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-polardb-backup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackPolardbBackupBasicDependencePostgreSQL0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":              "${alibabacloudstack_polardb_dbinstance.instance.id}",
					"backup_retention_period":     "9",
					"preferred_backup_time":       "10:00Z-11:00Z", // UTC Time
					"preferred_backup_period":     "Saturday,Sunday",
					"log_backup_retention_period": "7",
					"backup_log":                  "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func AlibabacloudStackPolardbBackupBasicDependencePostgreSQL0(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%v"
	}
	variable "creation" {
		default = "PolarDB"
	}
	resource "alibabacloudstack_polardb_dbinstance" "instance" {
		engine            = "PolarDB_PG"
		engine_version    = "14"
		instance_name = "${var.name}"
		db_instance_storage_type= "local_ssd"
		db_instance_storage = 10
		db_instance_class = "polardb.x4.medium.2"
		zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}
	`, RdsCommonTestCase, name)
}

