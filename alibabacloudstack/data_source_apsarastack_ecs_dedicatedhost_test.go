package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackEcsDedicatedHostsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_dedicated_hosts.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_dedicated_hosts.default.id}_fake"]`,
		}),
	}

	dedicated_host_cluster_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":                       `["${alibabacloudstack_ecs_dedicated_hosts.default.id}"]`,
			"dedicated_host_cluster_id": `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostClusterId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":                       `["${alibabacloudstack_ecs_dedicated_hosts.default.id}_fake"]`,
			"dedicated_host_cluster_id": `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostClusterId}_fake"`,
		}),
	}

	dedicated_host_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_dedicated_hosts.default.id}"]`,
			"dedicated_host_name": `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_dedicated_hosts.default.id}_fake"]`,
			"dedicated_host_name": `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostName}_fake"`,
		}),
	}

	dedicated_host_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_dedicated_hosts.default.id}"]`,
			"dedicated_host_type": `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_dedicated_hosts.default.id}_fake"]`,
			"dedicated_host_type": `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostType}_fake"`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_dedicated_hosts.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_dedicated_hosts.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_dedicated_hosts.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_dedicated_hosts.default.ResourceGroupId}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_dedicated_hosts.default.id}"]`,
			"status": `"${alibabacloudstack_ecs_dedicated_hosts.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_dedicated_hosts.default.id}_fake"]`,
			"status": `"${alibabacloudstack_ecs_dedicated_hosts.default.Status}_fake"`,
		}),
	}

	zone_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_ecs_dedicated_hosts.default.id}"]`,
			"zone_id": `"${alibabacloudstack_ecs_dedicated_hosts.default.ZoneId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_ecs_dedicated_hosts.default.id}_fake"]`,
			"zone_id": `"${alibabacloudstack_ecs_dedicated_hosts.default.ZoneId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_dedicated_hosts.default.id}"]`,

			"dedicated_host_cluster_id": `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostClusterId}"`,
			"dedicated_host_name":       `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostName}"`,
			"dedicated_host_type":       `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostType}"`,
			"resource_group_id":         `"${alibabacloudstack_ecs_dedicated_hosts.default.ResourceGroupId}"`,
			"status":                    `"${alibabacloudstack_ecs_dedicated_hosts.default.Status}"`,
			"zone_id":                   `"${alibabacloudstack_ecs_dedicated_hosts.default.ZoneId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_dedicated_hosts.default.id}_fake"]`,

			"dedicated_host_cluster_id": `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostClusterId}_fake"`,
			"dedicated_host_name":       `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostName}_fake"`,
			"dedicated_host_type":       `"${alibabacloudstack_ecs_dedicated_hosts.default.DedicatedHostType}_fake"`,
			"resource_group_id":         `"${alibabacloudstack_ecs_dedicated_hosts.default.ResourceGroupId}_fake"`,
			"status":                    `"${alibabacloudstack_ecs_dedicated_hosts.default.Status}_fake"`,
			"zone_id":                   `"${alibabacloudstack_ecs_dedicated_hosts.default.ZoneId}_fake"`}),
	}

	AlibabacloudstackEcsDedicatedHostsCheckInfo.dataSourceTestCheck(t, rand, idsConf, dedicated_host_cluster_idConf, dedicated_host_nameConf, dedicated_host_typeConf, resource_group_idConf, statusConf, zone_idConf, allConf)
}

var existAlibabacloudstackEcsDedicatedHostsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"hosts.#":    "1",
		"hosts.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsDedicatedHostsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"hosts.#": "0",
	}
}

var AlibabacloudstackEcsDedicatedHostsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_dedicated_hosts.default",
	existMapFunc: existAlibabacloudstackEcsDedicatedHostsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsDedicatedHostsMapFunc,
}

func testAccCheckAlibabacloudstackEcsDedicatedHostsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsDedicatedHosts%d"
}






data "alibabacloudstack_ecs_dedicated_hosts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
