package apsarastack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccApsaraStackKeyPairsDataSourceBasic(t *testing.T) {
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKeyPairsDataSourceConfig(map[string]string{
			"name_regex": `"${apsarastack_key_pair.default.key_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackKeyPairsDataSourceConfig(map[string]string{
			"name_regex": `"${apsarastack_key_pair.default.key_name}_fake"`,
		}),
	}
	/*tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKeyPairsDataSourceConfig(map[string]string{
			"name_regex": `"${apsarastack_key_pair.default.key_name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckApsaraStackKeyPairsDataSourceConfig(map[string]string{
			"name_regex": `"${apsarastack_key_pair.default.key_name}"`,
			"tags":       `{Created = "TF1"}`,
		}),
	}*/
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKeyPairsDataSourceConfig(map[string]string{
			"ids": `["${apsarastack_key_pair.default.key_name}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackKeyPairsDataSourceConfig(map[string]string{
			"ids": `["${apsarastack_key_pair.default.key_name}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackKeyPairsDataSourceConfig(map[string]string{
			"name_regex": `"${apsarastack_key_pair.default.key_name}"`,
			"ids":        `["${apsarastack_key_pair.default.key_name}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackKeyPairsDataSourceConfig(map[string]string{
			"name_regex": `"${apsarastack_key_pair.default.key_name}"`,
			"ids":        `["${apsarastack_key_pair.default.key_name}_fake"]`,
		}),
	}
	keyPairsCheckInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, allConf)
}

func testAccCheckApsaraStackKeyPairsDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
resource "apsarastack_key_pair" "default" {
	key_name = "tf-key-test"
}
data "apsarastack_key_pairs" "default" {
	%s
}`, strings.Join(pairs, "\n  "))
	return config
}

var existKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":                 "1",
		"ids.#":                   "1",
		"key_pairs.#":             "1",
		"key_pairs.0.id":          CHECKSET,
		"key_pairs.0.key_name":    "tf-key-test",
		"key_pairs.0.instances.#": "0",
		/*"key_pairs.0.tags.%":       "2",
		"key_pairs.0.tags.Created": "TF",
		"key_pairs.0.tags.For":     "acceptance test",*/
	}
}

var fakeKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":     "0",
		"ids.#":       "0",
		"key_pairs.#": "0",
		//"key_pairs.0.tags.%": "0",
	}
}

var keyPairsCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_key_pairs.default",
	existMapFunc: existKeyPairsMapFunc,
	fakeMapFunc:  fakeKeyPairsMapFunc,
}
