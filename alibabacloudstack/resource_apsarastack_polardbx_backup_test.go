package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbxBackup_basic(t *testing.T) {
	var v *PolarDbXBackupData
	resourceId := "alibabacloudstack_polardbx_backup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackPolardbxBackupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribePolarDbXBackup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1, 254)
	name := fmt.Sprintf("tf-testAccPolardbxBackupBasic_%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxBackupBasicDependence)

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
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${local.polardbx_instance.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

var AlibabacloudStackPolardbxBackupMap = map[string]string{
	"backup_model":    CHECKSET,
	"backup_set_size": CHECKSET,
	"backup_type":     CHECKSET,
	"status":          "1",
	"backup_set_id":   CHECKSET,
	"end_time":        CHECKSET,
	"begin_time":      CHECKSET,
}

func resourcePolardbxBackupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
%s

%s

 `, name, VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
