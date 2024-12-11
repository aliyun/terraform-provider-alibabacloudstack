package alibabacloudstack

import (
	"fmt"
	
	"strings"
	"testing"
)

func TestAccAlibabacloudStackSlbAclsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000,20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_slb_acl.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_slb_acl.default.name}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_slb_acl.default.name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_slb_acl.default.name}"`,
			"tags":       `{Created = "TF1"}`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackSlbAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackSlbAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_acl.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_slb_acl.default.id}"]`,
			"name_regex": `"${alibabacloudstack_slb_acl.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_slb_acl.default.id}_fake"]`,
			"name_regex": `"${alibabacloudstack_slb_acl.default.name}"`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":                     "1",
			"ids.#":                      "1",
			"names.#":                    "1",
			"acls.0.id":                  CHECKSET,
			"acls.0.name":                fmt.Sprintf("tf-testAccSlbAclDataSourceBisic-%d", rand),
			"acls.0.ip_version":          "ipv4",
			"acls.0.entry_list.#":        "2",
			"acls.0.related_listeners.#": "0",
			"acls.0.tags.%":              "2",
			"acls.0.tags.Created":        "TF",
			"acls.0.tags.For":            "acceptance test",
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":  "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var slbaclsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_slb_acls.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbaclsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, tagsConf, idsConf, allConf)
}

func testAccCheckAlibabacloudStackSlbAclsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlbAclDataSourceBisic-%d"
}
variable "ip_version" {
	default = "ipv4"
}

resource "alibabacloudstack_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
    entry = "10.10.10.0/24"
    comment = "first"
  }
  entry_list {
      entry = "168.10.10.0/24"
      comment = "second"
  }
   tags = {
      Created = "TF"
       For     = "acceptance test"
    }
}

data "alibabacloudstack_slb_acls" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
