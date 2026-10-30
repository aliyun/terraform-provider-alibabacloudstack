package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccAlibabacloudStackPolardbClusterBackup_basic tests the manual backup of a PolarDB cluster edition instance
func TestAccAlibabacloudStackPolardbClusterBackup_basic(t *testing.T) {
	var v *PolardbClusterBackupData
	resourceId := "alibabacloudstack_polardb_cluster_backup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackPolardbClusterBackupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribePolardbClusterBackup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1, 254)
	name := fmt.Sprintf("tf-testAccPolardbClusterBackupBasic_%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbClusterBackupBasicDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckApiAvailable(t, "POST", "polardb", "2017-08-01", "CreateBackup", "")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": "${local.polardb_cluster_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_job_id"},
			},
		},
	})
}

var AlibabacloudStackPolardbClusterBackupMap = map[string]string{
	"db_cluster_id":     CHECKSET,
	"backup_id":         CHECKSET,
	"backup_method":     CHECKSET,
	"backup_mode":       CHECKSET,
	"backup_type":       CHECKSET,
	"backup_status":     "Success",
	"backups_level":     CHECKSET,
	"store_status":      CHECKSET,
	"backup_set_size":   CHECKSET,
	"backup_start_time": CHECKSET,
	"backup_end_time":   CHECKSET,
}

func resourcePolardbClusterBackupBasicDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

`, name, VSwitchCommonTestCase+PolardbClusterCommonTestCase("PostgreSQL", "16", true))
}
