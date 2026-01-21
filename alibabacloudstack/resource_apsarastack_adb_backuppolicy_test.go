package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAdbBackupPolicy(t *testing.T) {
	var v *adb.DescribeBackupPolicyResponse
	resourceId := "alibabacloudstack_adb_backup_policy.default"
	serverFunc := func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeAdbBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccAdbBackupPolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAdbBackupPolicyConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// This feature does not have a delete interface, only restore to default
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// Use given id to test with existing instance
					//"db_cluster_id":    "am-3rq9uva152cn34drs",
					"db_cluster_id":           "${local.adb_instance_id}",
					"preferred_backup_period": []string{"Tuesday", "Wednesday"},
					"preferred_backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
						"preferred_backup_time":     "10:00Z-11:00Z",
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
					"preferred_backup_period": []string{"Wednesday", "Monday", "Saturday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "15:00Z-16:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "15:00Z-16:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Tuesday", "Thursday", "Friday", "Sunday"},
					"preferred_backup_time":   "17:00Z-18:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "4",
						"preferred_backup_time":     "17:00Z-18:00Z",
					}),
				),
			}},
	})
}

func resourceAdbBackupPolicyConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	%s
	
`, name, AdbCommonTestCase(false))
}
