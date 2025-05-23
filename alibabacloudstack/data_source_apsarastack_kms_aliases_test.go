package alibabacloudstack

import (
	"fmt"
	"testing"

	
)

func TestAccAlibabacloudStackKmsAliasesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_kms_aliases.default"
	rand := getAccTestRandInt(1000000, 9999999)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("alias/test_kms_ali%d", rand), dataSourceKmsAliasesDependence)

	NameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_kms_alias.this.alias_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_kms_alias.this.alias_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_kms_alias.this.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_kms_alias.this.id}-fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_kms_alias.this.alias_name}",
			"ids":        []string{"${alibabacloudstack_kms_alias.this.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_kms_alias.this.alias_name}-fake",
			"ids":        []string{"${alibabacloudstack_kms_alias.this.id}-fake"},
		}),
	}
	var existKmsCiphertextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"ids.0":                CHECKSET,
			"names.#":              "1",
			"names.0":              CHECKSET,
			"aliases.#":            "1",
			"aliases.0.id":         CHECKSET,
			"aliases.0.alias_name": CHECKSET,
			"aliases.0.key_id":     CHECKSET,
		}
	}

	var fakeKmsCiphertextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"aliases.#": "0",
		}
	}

	var kmsCipherCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsCiphertextMapFunc,
		fakeMapFunc:  fakeKmsCiphertextMapFunc,
	}

	kmsCipherCheckInfo.dataSourceTestCheck(t, 0, NameRegexConf, idsConf, allConf)
}

func dataSourceKmsAliasesDependence(name string) string {
	return fmt.Sprintf(`
    resource "alibabacloudstack_kms_key" "this" {}

	resource "alibabacloudstack_kms_alias" "this" {
  		alias_name = "%s"
  		key_id = "${alibabacloudstack_kms_key.this.id}"
	}
	`, name)
}
