package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackGpdbBackupPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_gpdb_backup_policy.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tftest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackGpdbBackupPolicyBasicDependence0)
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
					"db_instance_id":          "${local.gpdb_instance_id}",
					"preferred_backup_time":   "02:00Z-03:00Z",
					"preferred_backup_period": "Monday,Wednesday,Friday",
					"backup_retention_period": "7",
					"enable_recovery_point":   "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time":   "02:00Z-03:00Z",
						"preferred_backup_period": "Monday,Wednesday,Friday",
						"backup_retention_period": "7",
						"enable_recovery_point":   "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time":   "04:00Z-05:00Z",
					"preferred_backup_period": "Monday,Wednesday",
					"backup_retention_period": "5",
					"enable_recovery_point":   "true",
					"recovery_point_period":   "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":          CHECKSET,
						"preferred_backup_time":   "04:00Z-05:00Z",
						"preferred_backup_period": "Monday,Wednesday",
						"backup_retention_period": "5",
						"enable_recovery_point":   "true",
						"recovery_point_period":   "2",
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

func AlibabacloudStackGpdbBackupPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
%s
`, name, GpdbCommonTestCase())
}
