package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDBBackupPolicy_basic(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "alibabacloudstack_db_backup_policy.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-rds-backuppolicy%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, dbBackupPolicyConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":                 "${alibabacloudstack_db_instance.default.id}",
					"preferred_backup_period":     []string{"Sunday", "Thursday", "Tuesday"},
					"preferred_backup_time":       "18:00Z-19:00Z",
					"backup_retention_period":     "34",
					"enable_backup_log":           "true",
					"log_backup_retention_period": "34",
					"local_log_retention_hours":   "24",
					"high_space_usage_protection": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#":   "3",
						"preferred_backup_time":       "18:00Z-19:00Z",
						"backup_retention_period":     "34",
						"enable_backup_log":           "true",
						"log_backup_retention_period": "34",
						"local_log_retention_hours":   "24",
						"high_space_usage_protection": "Enable",
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
					"backup_retention_period": "10",
					"preferred_backup_period": []string{"Monday", "Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period":     "10",
						"preferred_backup_period.#":   "2",
						"preferred_backup_time":       "18:00Z-19:00Z",
						"enable_backup_log":           "true",
						"log_backup_retention_period": "34",
						"local_log_retention_hours":   "24",
						"high_space_usage_protection": "Enable",
					}),
				),
			},
		},
	})
}

func dbBackupPolicyConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

%s
`, name, VSwitchCommonTestCase, RdsMysqlCommonTestCase())
}
