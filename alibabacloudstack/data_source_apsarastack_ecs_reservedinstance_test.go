package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackEcsReservedInstancesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_reserved_instances.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_reserved_instances.default.id}_fake"]`,
		}),
	}

	instance_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ecs_reserved_instances.default.id}"]`,
			"instance_type": `"${alibabacloudstack_ecs_reserved_instances.default.InstanceType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ecs_reserved_instances.default.id}_fake"]`,
			"instance_type": `"${alibabacloudstack_ecs_reserved_instances.default.InstanceType}_fake"`,
		}),
	}

	offering_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ecs_reserved_instances.default.id}"]`,
			"offering_type": `"${alibabacloudstack_ecs_reserved_instances.default.OfferingType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ecs_reserved_instances.default.id}_fake"]`,
			"offering_type": `"${alibabacloudstack_ecs_reserved_instances.default.OfferingType}_fake"`,
		}),
	}

	reserved_instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":                  `["${alibabacloudstack_ecs_reserved_instances.default.id}"]`,
			"reserved_instance_id": `"${alibabacloudstack_ecs_reserved_instances.default.ReservedInstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":                  `["${alibabacloudstack_ecs_reserved_instances.default.id}_fake"]`,
			"reserved_instance_id": `"${alibabacloudstack_ecs_reserved_instances.default.ReservedInstanceId}_fake"`,
		}),
	}

	reserved_instance_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":                    `["${alibabacloudstack_ecs_reserved_instances.default.id}"]`,
			"reserved_instance_name": `"${alibabacloudstack_ecs_reserved_instances.default.ReservedInstanceName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":                    `["${alibabacloudstack_ecs_reserved_instances.default.id}_fake"]`,
			"reserved_instance_name": `"${alibabacloudstack_ecs_reserved_instances.default.ReservedInstanceName}_fake"`,
		}),
	}

	scopeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":   `["${alibabacloudstack_ecs_reserved_instances.default.id}"]`,
			"scope": `"${alibabacloudstack_ecs_reserved_instances.default.Scope}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":   `["${alibabacloudstack_ecs_reserved_instances.default.id}_fake"]`,
			"scope": `"${alibabacloudstack_ecs_reserved_instances.default.Scope}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_reserved_instances.default.id}"]`,
			"status": `"${alibabacloudstack_ecs_reserved_instances.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_reserved_instances.default.id}_fake"]`,
			"status": `"${alibabacloudstack_ecs_reserved_instances.default.Status}_fake"`,
		}),
	}

	zone_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_ecs_reserved_instances.default.id}"]`,
			"zone_id": `"${alibabacloudstack_ecs_reserved_instances.default.ZoneId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_ecs_reserved_instances.default.id}_fake"]`,
			"zone_id": `"${alibabacloudstack_ecs_reserved_instances.default.ZoneId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_reserved_instances.default.id}"]`,

			"instance_type":          `"${alibabacloudstack_ecs_reserved_instances.default.InstanceType}"`,
			"offering_type":          `"${alibabacloudstack_ecs_reserved_instances.default.OfferingType}"`,
			"reserved_instance_id":   `"${alibabacloudstack_ecs_reserved_instances.default.ReservedInstanceId}"`,
			"reserved_instance_name": `"${alibabacloudstack_ecs_reserved_instances.default.ReservedInstanceName}"`,
			"scope":                  `"${alibabacloudstack_ecs_reserved_instances.default.Scope}"`,
			"status":                 `"${alibabacloudstack_ecs_reserved_instances.default.Status}"`,
			"zone_id":                `"${alibabacloudstack_ecs_reserved_instances.default.ZoneId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_reserved_instances.default.id}_fake"]`,

			"instance_type":          `"${alibabacloudstack_ecs_reserved_instances.default.InstanceType}_fake"`,
			"offering_type":          `"${alibabacloudstack_ecs_reserved_instances.default.OfferingType}_fake"`,
			"reserved_instance_id":   `"${alibabacloudstack_ecs_reserved_instances.default.ReservedInstanceId}_fake"`,
			"reserved_instance_name": `"${alibabacloudstack_ecs_reserved_instances.default.ReservedInstanceName}_fake"`,
			"scope":                  `"${alibabacloudstack_ecs_reserved_instances.default.Scope}_fake"`,
			"status":                 `"${alibabacloudstack_ecs_reserved_instances.default.Status}_fake"`,
			"zone_id":                `"${alibabacloudstack_ecs_reserved_instances.default.ZoneId}_fake"`}),
	}

	AlibabacloudstackEcsReservedInstancesCheckInfo.dataSourceTestCheck(t, rand, idsConf, instance_typeConf, offering_typeConf, reserved_instance_idConf, reserved_instance_nameConf, scopeConf, statusConf, zone_idConf, allConf)
}

var existAlibabacloudstackEcsReservedInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":    "1",
		"instances.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsReservedInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#": "0",
	}
}

var AlibabacloudstackEcsReservedInstancesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_reserved_instances.default",
	existMapFunc: existAlibabacloudstackEcsReservedInstancesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsReservedInstancesMapFunc,
}

func testAccCheckAlibabacloudstackEcsReservedInstancesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsReservedInstances%d"
}






data "alibabacloudstack_ecs_reserved_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
