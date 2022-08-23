package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackDbsBackupPlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_dbs_backup_plan.default"
	ra := resourceAttrInit(resourceId, ApsaraStackDbsBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DbsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDbsBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdbsbackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackDbsBackupPlanBasicDependence0)
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

var ApsaraStackDbsBackupPlanMap0 = map[string]string{}

func ApsaraStackDbsBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
