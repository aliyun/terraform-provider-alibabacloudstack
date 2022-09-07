package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlibabacloudStackDnsGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)

	testAccConfig := dataSourceTestAccConfigFunc("data.alibabacloudstack_dns_groups.default", fmt.Sprintf("tf-testacc-%d", rand), dataSourceDnsGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dns_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dns_group.default.name}_fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_dns_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_dns_group.default.id}_fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dns_group.default.name}",
			"ids":        []string{"${alibabacloudstack_dns_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dns_group.default.name}_fake",
			"ids":        []string{"${alibabacloudstack_dns_group.default.id}"},
		}),
	}
	existChangeMap := map[string]string{
		"ids.#":               "1",
		"ids.0":               REMOVEKEY,
		"names.#":             "1",
		"names.0":             "ALL",
		"groups.#":            "1",
		"groups.0.group_id":   "",
		"groups.0.group_name": "ALL",
	}
	nameAllConf := dataSourceTestAccConfig{
		existConfig:   testAccCheckAlibabacloudStackDnsGroupsDataSourceNameRegexAll,
		existChangMap: existChangeMap,
	}

	dnsGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf, nameAllConf)
}

func dataSourceDnsGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_dns_group" "default" {
	name = "%s"
}
`, name)
}

const testAccCheckAlibabacloudStackDnsGroupsDataSourceNameRegexAll = `
data "alibabacloudstack_dns_groups" "default" {
  name_regex = "^ALL"
}`

var existDnsGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":               "1",
		"ids.0":               CHECKSET,
		"names.#":             "1",
		"names.0":             fmt.Sprintf("tf-testacc-%d", rand),
		"groups.0.group_id":   CHECKSET,
		"groups.0.group_name": fmt.Sprintf("tf-testacc-%d", rand),
	}
}

var fakeDnsGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":    "0",
		"names.#":  "0",
		"groups.#": "0",
	}
}

var dnsGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_dns_groups.default",
	existMapFunc: existDnsGroupsMapFunc,
	fakeMapFunc:  fakeDnsGroupsMapFunc,
}
