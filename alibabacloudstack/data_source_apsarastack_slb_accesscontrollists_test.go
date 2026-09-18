package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbAclsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_acls.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccSlbAclDataSourceBisic-%d", rand),
		dataSourceSlbAclsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_slb_acl.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_slb_acl.default.name}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_slb_acl.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_slb_acl.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_slb_acl.default.id}"},
			"name_regex": "${alibabacloudstack_slb_acl.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_slb_acl.default.id}_fake"},
			"name_regex": "${alibabacloudstack_slb_acl.default.name}",
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
		resourceId:   resourceId,
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbaclsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceSlbAclsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
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
}
`, name)
}
