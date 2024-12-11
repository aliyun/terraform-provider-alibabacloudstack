package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackEcsInstancesDataSource(t *testing.T) {
	// 根据test_meta自动生成的tasecase

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_instances.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_instances.default.id}_fake"]`,
		}),
	}

	instance_network_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_ecs_instances.default.id}"]`,
			"instance_network_type": `"${alibabacloudstack_ecs_instances.default.InstanceNetworkType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_ecs_instances.default.id}_fake"]`,
			"instance_network_type": `"${alibabacloudstack_ecs_instances.default.InstanceNetworkType}_fake"`,
		}),
	}

	payment_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ecs_instances.default.id}"]`,
			"payment_type": `"${alibabacloudstack_ecs_instances.default.PaymentType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ecs_instances.default.id}_fake"]`,
			"payment_type": `"${alibabacloudstack_ecs_instances.default.PaymentType}_fake"`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_instances.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_instances.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_instances.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_instances.default.ResourceGroupId}_fake"`,
		}),
	}

	zone_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_ecs_instances.default.id}"]`,
			"zone_id": `"${alibabacloudstack_ecs_instances.default.ZoneId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_ecs_instances.default.id}_fake"]`,
			"zone_id": `"${alibabacloudstack_ecs_instances.default.ZoneId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_instances.default.id}"]`,

			"instance_network_type": `"${alibabacloudstack_ecs_instances.default.InstanceNetworkType}"`,
			"payment_type":          `"${alibabacloudstack_ecs_instances.default.PaymentType}"`,
			"resource_group_id":     `"${alibabacloudstack_ecs_instances.default.ResourceGroupId}"`,
			"zone_id":               `"${alibabacloudstack_ecs_instances.default.ZoneId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_instances.default.id}_fake"]`,

			"instance_network_type": `"${alibabacloudstack_ecs_instances.default.InstanceNetworkType}_fake"`,
			"payment_type":          `"${alibabacloudstack_ecs_instances.default.PaymentType}_fake"`,
			"resource_group_id":     `"${alibabacloudstack_ecs_instances.default.ResourceGroupId}_fake"`,
			"zone_id":               `"${alibabacloudstack_ecs_instances.default.ZoneId}_fake"`}),
	}

	AlibabacloudstackEcsInstancesDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, instance_network_typeConf, payment_typeConf, resource_group_idConf, zone_idConf, allConf)
}

var existAlibabacloudstackEcsInstancesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":    "1",
		"instances.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsInstancesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#": "0",
	}
}

var AlibabacloudstackEcsInstancesDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_instances.default",
	existMapFunc: existAlibabacloudstackEcsInstancesDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsInstancesDataMapFunc,
}

func testAccCheckAlibabacloudstackEcsInstancesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsInstances%d"
}






data "alibabacloudstack_ecs_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

