package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

// Dependency resource template generation method
func testAccCheckAlibabacloudStackNasDirQuotasDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_nas_dir_quota" "default" {
    file_system_id = "${alibabacloudstack_nas_file_system.default.id}"
	path = "/"
	quotas {
	    quota_type = "Enforcement"
	    user_type  = "Uid"
	    user_id    = "500"
		size_limit = 100
		file_count_limit = 10000
	}
}
`, name, NasCommonTestCase)
}

func TestAccAlibabacloudStackNasDirQuotasDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alibabacloudstack_nas_dir_quotas.default"
	// Common dependency configuration for NAS directory quotas
	name := fmt.Sprintf("tf-testnasdirquotas%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackNasDirQuotasDataSourceConfig)

	testAccCheck := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                  "1",
				"dir_quotas.#":           "1",
				"dir_quotas.0.path":      "/",
				"dir_quotas.0.dir_inode": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":        "0",
				"dir_quotas.#": "0",
			}
		},
	}

	// Define test configurations for different scenarios
	testAccCheck.dataSourceTestCheck(t, rand,

		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"file_system_id": "${alibabacloudstack_nas_dir_quota.default.file_system_id}",
				"path":           "${alibabacloudstack_nas_dir_quota.default.path}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"file_system_id": "${alibabacloudstack_nas_file_system.default.id}",
				"path":           "/fake/path",
			}),
		},
		// Test with ids filter
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"file_system_id": "${alibabacloudstack_nas_dir_quota.default.file_system_id}",
				"ids":            []string{"${alibabacloudstack_nas_dir_quota.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"file_system_id": "${alibabacloudstack_nas_dir_quota.default.file_system_id}",
				"ids":            []string{"${alibabacloudstack_nas_dir_quota.default.id}_fake"},
			}),
		},
		// Test with name_regex filter
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"file_system_id": "${alibabacloudstack_nas_dir_quota.default.file_system_id}",
				"name_regex":     "^/",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"file_system_id": "${alibabacloudstack_nas_dir_quota.default.file_system_id}",
				"name_regex":     "^/fake/path",
			}),
		},
	)
}
