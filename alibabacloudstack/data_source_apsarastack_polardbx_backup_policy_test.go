package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbxBackupPoliciesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testAccPolardbxBackupplicies%v", rand)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackPolardbxBackupPolicysDataSource(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_polardbx_backup_policies.default"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "backup_period"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "backup_set_retention"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "backup_plan_begin"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "remove_log_retention"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "cold_data_backup_interval"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "local_log_retention_number"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "cold_data_backup_retention"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "force_clean_on_high_space_usage"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "backup_way"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "local_log_retention"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "backup_type"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policies.default", "log_local_retention_space"),
				),
			},
		},
	})
}

func testAccCheckAlibabacloudStackPolardbxBackupPolicysDataSource(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

data "alibabacloudstack_polardbx_backup_policies" "default" {
	db_instance_id = "${local.polardbx_instance.id}"
}
`, name, VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
