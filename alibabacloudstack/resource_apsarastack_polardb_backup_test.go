package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbBackup_basic(t *testing.T) {
	var v *PolardbbackupData
	resourceId := "alibabacloudstack_polardb_backup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackPolardbBackupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoPolardbDescribebackupsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1, 254)
	name := fmt.Sprintf("tf-testAccPolardbBackupBasic_%v", rand)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: resourcePolardbBackupBasicDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_method": "Physical",
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

var AlibabacloudStackPolardbBackupMap = map[string]string{
	"backup_method":     CHECKSET,
	"backup_id":         CHECKSET,
	"backup_mode":       CHECKSET,
	"backup_status":     "Success",
	"backup_size":       CHECKSET,
	"slave_status":      CHECKSET,
	"host_instance_id":  CHECKSET,
	"backup_db_names":   CHECKSET,
	"store_status":      CHECKSET,
	"backup_end_time":   CHECKSET,
	"backup_start_time": CHECKSET,
	"meta_status":       CHECKSET,
	"backup_scale":      CHECKSET,
	"backup_location":   CHECKSET,
}

func resourcePolardbBackupBasicDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_polardb_dbinstance" "default" {
  instance_storage = "5"
  instance_name = "${var.name}"
  storage_type = "local_ssd"
  engine = "MySQL"
  engine_version = "5.7"
  instance_type = "rds.mysql.t1.small"
}

resource "alibabacloudstack_polardb_backup" "default" {
	db_instance_id = "${alibabacloudstack_polardb_dbinstance.default.id}"
	backup_method = "Physical"
}

`, name)
}
