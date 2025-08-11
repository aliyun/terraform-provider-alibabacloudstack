package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbxBackupPolicy_basic0(t *testing.T) {
	var v *PolarDbXBackupConfig

	resourceId := "alibabacloudstack_polardbx_backup_policy.default"
	backupPolicyCheckset := map[string]string{
		"cold_data_backup_interval":       CHECKSET,
		"local_log_retention_number":      CHECKSET,
		"cold_data_backup_retention":      CHECKSET,
		"force_clean_on_high_space_usage": CHECKSET,
		"backup_way":                      CHECKSET,
		"local_log_retention":             CHECKSET,
		"backup_type":                     CHECKSET,
		"log_local_retention_space":       CHECKSET,
	}
	ra := resourceAttrInit(resourceId, backupPolicyCheckset)
	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDbXBackupConfig")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_pldbx_backup_policy_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxBackupPolicyDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":                  "${local.polardbx_instance.id}",
					"backup_period":                   "Monday,Wednesday,Sunday",
					"backup_set_retention":            "45",
					"backup_plan_begin":               "04:00Z",
					"remove_log_retention":            "45",
					"force_clean_on_high_space_usage": "1",
					"local_log_retention":             "12",
					"log_local_retention_space":       "35",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period":                   "Monday,Wednesday,Sunday",
						"backup_set_retention":            "45",
						"backup_plan_begin":               "04:00Z",
						"remove_log_retention":            "45",
						"force_clean_on_high_space_usage": "1",
						"local_log_retention":             "12",
						"log_local_retention_space":       "35",
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
					"backup_set_retention":            "60",
					"remove_log_retention":            "60",
					"backup_plan_begin":               "06:00Z",
					"backup_period":                   "Monday,Wednesday,Thursday,Sunday",
					"force_clean_on_high_space_usage": "1",
					"local_log_retention":             "7",
					"log_local_retention_space":       "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_set_retention":            "60",
						"remove_log_retention":            "60",
						"backup_plan_begin":               "06:00Z",
						"backup_period":                   "Monday,Wednesday,Thursday,Sunday",
						"force_clean_on_high_space_usage": "1",
						"local_log_retention":             "7",
						"log_local_retention_space":       "30",
					}),
				),
			},
		},
	})
}

func resourcePolardbxBackupPolicyDependence(name string) string {
	return fmt.Sprintf(`
 variable "name" {
   default = "%s"
 }

%s

%s

 `, name, VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
