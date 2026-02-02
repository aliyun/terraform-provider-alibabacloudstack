package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbxBackupPoliciesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardbx_backup_policies.default"
	name := fmt.Sprintf("tf-testAccPolardbxBackupplicies%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbxBackupPoliciesConfigDependence)

	defaultConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
		// No fake config needed since db_instance_id is required and valid instance always exists in existConfig
	}

	var existPolardbxBackupPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"backup_period":                    CHECKSET,
			"backup_set_retention":             CHECKSET,
			"backup_plan_begin":                CHECKSET,
			"remove_log_retention":             CHECKSET,
			"cold_data_backup_interval":        CHECKSET,
			"local_log_retention_number":       CHECKSET,
			"cold_data_backup_retention":       CHECKSET,
			"force_clean_on_high_space_usage":  CHECKSET,
			"backup_way":                       CHECKSET,
			"local_log_retention":              CHECKSET,
			"backup_type":                      CHECKSET,
			"log_local_retention_space":        CHECKSET,
		}
	}

	// Since this data source requires a valid db_instance_id and always returns values for existing instances,
	// we don't have a realistic "fake" scenario that returns empty results.
	// Thus, we only define existMapFunc and use an empty fakeMapFunc (or same as exist).
	var fakePolardbxBackupPoliciesMapFunc = func(rand int) map[string]string {
		// In practice, querying with invalid db_instance_id would error, not return empty.
		// So we skip fake validation or assume same structure.
		return map[string]string{}
	}

	var polardbxBackupPoliciesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbxBackupPoliciesMapFunc,
		fakeMapFunc:  fakePolardbxBackupPoliciesMapFunc,
	}
	polardbxBackupPoliciesCheckInfo.dataSourceTestCheck(t, rand, defaultConf)
}

func dataSourcePolardbxBackupPoliciesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

data "alibabacloudstack_polardbx_backup_policies" "default" {
  db_instance_id = local.polardbx_instance.id
}
`, name, VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
