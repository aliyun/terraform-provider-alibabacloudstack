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
	default = "testacc-agsoss-%d"
}

data "alibabacloudstack_oss_clusters" "default" {
  provider = alibabacloudstack-common
}

resource "alibabacloudstack_oss_bucket" "default" {
  provider = alibabacloudstack-common
  bucket = "${var.name}"
  oss_cluster = "${data.alibabacloudstack_oss_clusters.default.clusters.0.id}"
}

resource "alibabacloudstack_aqs_oss_scanconfig" "default" {
  enable                     = true
  start_time                 = "00:00:00"
  end_time                   = "03:00:00"
  scan_day_list              = [5, 6, 7]
  scan_mode                  = "1"
  bucket_name                = alibabacloudstack_oss_bucket.default.bucket
  decryption                 = "OSS"
  key_prefix                 = "test"
  key_suffix                 = ".py"
  last_modified_start_time   = "2025-12-21 00:00:00"
}

data "alibabacloudstack_aqs_oss_scanconfigs" "default" {
  %s
}

`, rand, strings.Join(pairs, "\n  "))
}
