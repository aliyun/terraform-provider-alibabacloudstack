package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDbsBackupPlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dbs_backup_plan.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDbsBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DbsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDbsBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdbsbackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDbsBackupPlanBasicDependence0)
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
					"backup_method":  "logical",
					"database_type":  "MySQL",
					"instance_class": "large",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_method":  "logical",
						"database_type":  "MySQL",
						"instance_class": "large",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_plan_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_plan_name": name,
					}),
				),
			},
		},
	})
}

var AlibabacloudStackDbsBackupPlanMap0 = map[string]string{}

func AlibabacloudStackDbsBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
