package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackMaxcomputeUsersDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_maxcompute_users.default"
	name := fmt.Sprintf("tf_testAcck%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceMaxcomputeUsersConfigDependence)

	// Test with name_regex filter (should match the created user)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_maxcompute_user.default.user_name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-nonexistent-user",
		}),
	}

	// Test with ids filter (using the created user's ID)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_maxcompute_user.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-user-id-12345"},
		}),
	}

	var existMaxcomputeUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"users.#":              "1",
			"users.0.id":           CHECKSET,
			"users.0.user_id":      CHECKSET,
			"users.0.user_pk":      CHECKSET,
			"users.0.user_name":    name,
			"users.0.user_type":    CHECKSET,
			"users.0.description":  "TestAccAlibabacloudStackMaxcomputeUser",
		}
	}

	var fakeMaxcomputeUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"users.#":  "0",
		}
	}

	var maxcomputeUsersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMaxcomputeUsersMapFunc,
		fakeMapFunc:  fakeMaxcomputeUsersMapFunc,
	}
	maxcomputeUsersCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf)
}

func dataSourceMaxcomputeUsersConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_maxcompute_user" "default" {
  user_name   = "%s"
  description = "TestAccAlibabacloudStackMaxcomputeUser"
}
`, name)
}
