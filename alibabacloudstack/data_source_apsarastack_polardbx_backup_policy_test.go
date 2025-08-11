package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbxBackupPolicysDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackPolardbxBackupPolicysDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_polardbx_backup_policys.default"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "backup_period"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "backup_set_retention"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "backup_plan_begin"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "remove_log_retention"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "cold_data_backup_interval"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "local_log_retention_number"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "cold_data_backup_retention"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "force_clean_on_high_space_usage"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "backup_way"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "local_log_retention"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "backup_type"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_polardbx_backup_policys.default", "log_local_retention_space"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackPolardbxBackupPolicysDataSource = `
data "alibabacloudstack_polardbx_backup_policys" "default" {
	db_instance_id = "pxc-unrpi3i87xv25d"
}
`
