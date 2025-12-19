package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackAqsOssScanconfigsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_aqs_oss_scanconfigs.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                    "1",
				"names.#":                  "1",
				"scan_configs.#":           "1",
				"scan_configs.0.enable":    "true",
				"scan_configs.0.scan_mode": "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":          "0",
				"names.#":        "0",
				"scan_configs.#": "0",
			}
		},
		Providers: testYunDunProviders(),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: resourceAqsOssScanconfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_aqs_oss_scanconfig.default.bucket_name}"`,
		}),
		fakeConfig: resourceAqsOssScanconfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_aqs_oss_scanconfig.default.bucket_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: resourceAqsOssScanconfigDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_aqs_oss_scanconfig.default.id}"]`,
		}),
		fakeConfig: resourceAqsOssScanconfigDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_aqs_oss_scanconfig.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: resourceAqsOssScanconfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_aqs_oss_scanconfig.default.bucket_name}"`,
			"ids":        `["${alibabacloudstack_aqs_oss_scanconfig.default.id}"]`,
		}),
		fakeConfig: resourceAqsOssScanconfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_aqs_oss_scanconfig.default.bucket_name}_fake"`,
			"ids":        `["${alibabacloudstack_aqs_oss_scanconfig.default.id}_fake"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

// Dependence template generation method
func resourceAqsOssScanconfigDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
	default = "tf-Acctest%d"
}

resource "alibabacloudstack_aqs_oss_scanconfig" "default" {
	enable                   = true
	start_time               = "00:00:00"
	end_time                 = "03:00:00"
	scan_day_list            = [4, 5]
	scan_mode                = "1"
	bucket_name              = "testtf2"
	decryption               = "OSS"
	key_suffix               = "all"
	key_prefix               = "test"
	last_modified_start_time = "2025-12-12 12:20:27"
}
	
data "alibabacloudstack_aqs_oss_scanconfigs" "default" {
  %s
}

`, rand, strings.Join(pairs, "\n  "))
}
