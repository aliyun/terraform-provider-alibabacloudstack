package alibabacloudstack

import (
	"fmt"
	"testing"

	
)

func TestAccAlibabacloudStackKmsKeysDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_kms_keys.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccKmsKeysDataSource_%d", rand),
		dataSourceKmsKeysConfigDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alibabacloudstack_kms_key.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alibabacloudstack_kms_key.default.description}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_kms_key.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_kms_key.default.id}-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alibabacloudstack_kms_key.default.description}",
			"ids":               []string{"${alibabacloudstack_kms_key.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alibabacloudstack_kms_key.default.description}-fake",
			"ids":               []string{"${alibabacloudstack_kms_key.default.id}"},
		}),
	}

	var existKmsKeysMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"keys.0.description": fmt.Sprintf("tf_testAccKmsKeysDataSource_%d", rand),
		}
	}

	var fakeKmsKeysMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"keys.#": "0",
		}
	}

	var kmsKeysCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsKeysMapFunc,
		fakeMapFunc:  fakeKmsKeysMapFunc,
	}

	kmsKeysCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

func dataSourceKmsKeysConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_kms_key" "default" {
    description = "%s"
    pending_window_in_days = 7
}
`, name)
}
