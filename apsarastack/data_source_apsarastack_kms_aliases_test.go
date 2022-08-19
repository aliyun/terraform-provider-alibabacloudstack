package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackKmsAliasesDataSource(t *testing.T) {
	resourceId := "data.apsarastack_kms_aliases.default"
	rand := acctest.RandIntRange(1000000, 9999999)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("alias/test_kms_ali%d", rand), dataSourceKmsAliasesDependence)

	NameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${apsarastack_kms_alias.this.alias_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${apsarastack_kms_alias.this.alias_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${apsarastack_kms_alias.this.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${apsarastack_kms_alias.this.id}-fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${apsarastack_kms_alias.this.alias_name}",
			"ids":        []string{"${apsarastack_kms_alias.this.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${apsarastack_kms_alias.this.alias_name}-fake",
			"ids":        []string{"${apsarastack_kms_alias.this.id}-fake"},
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
    resource "apsarastack_kms_key" "this" {}

	resource "apsarastack_kms_alias" "this" {
  		alias_name = "%s"
  		key_id = "${apsarastack_kms_key.this.id}"
	}
	`, name)
}
