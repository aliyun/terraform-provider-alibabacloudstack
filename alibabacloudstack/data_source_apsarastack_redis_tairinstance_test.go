package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackRedisTairInstancesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
		}),
	}

	architecture_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"architecture_type": `"${alibabacloudstack_redis_tair_instances.default.ArchitectureType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"architecture_type": `"${alibabacloudstack_redis_tair_instances.default.ArchitectureType}_fake"`,
		}),
	}

	engine_versionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"engine_version": `"${alibabacloudstack_redis_tair_instances.default.EngineVersion}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"engine_version": `"${alibabacloudstack_redis_tair_instances.default.EngineVersion}_fake"`,
		}),
	}

	instance_classConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"instance_class": `"${alibabacloudstack_redis_tair_instances.default.InstanceClass}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"instance_class": `"${alibabacloudstack_redis_tair_instances.default.InstanceClass}_fake"`,
		}),
	}

	instance_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"instance_type": `"${alibabacloudstack_redis_tair_instances.default.InstanceType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"instance_type": `"${alibabacloudstack_redis_tair_instances.default.InstanceType}_fake"`,
		}),
	}

	network_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"network_type": `"${alibabacloudstack_redis_tair_instances.default.NetworkType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"network_type": `"${alibabacloudstack_redis_tair_instances.default.NetworkType}_fake"`,
		}),
	}

	payment_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"payment_type": `"${alibabacloudstack_redis_tair_instances.default.PaymentType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"payment_type": `"${alibabacloudstack_redis_tair_instances.default.PaymentType}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"status": `"${alibabacloudstack_redis_tair_instances.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"status": `"${alibabacloudstack_redis_tair_instances.default.Status}_fake"`,
		}),
	}

	tair_instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"tair_instance_id": `"${alibabacloudstack_redis_tair_instances.default.TairInstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"tair_instance_id": `"${alibabacloudstack_redis_tair_instances.default.TairInstanceId}_fake"`,
		}),
	}

	tair_instance_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"tair_instance_name": `"${alibabacloudstack_redis_tair_instances.default.TairInstanceName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"tair_instance_name": `"${alibabacloudstack_redis_tair_instances.default.TairInstanceName}_fake"`,
		}),
	}

	vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"vswitch_id": `"${alibabacloudstack_redis_tair_instances.default.VSwitchId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"vswitch_id": `"${alibabacloudstack_redis_tair_instances.default.VSwitchId}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_redis_tair_instances.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_redis_tair_instances.default.VpcId}_fake"`,
		}),
	}

	zone_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_redis_tair_instances.default.id}"]`,
			"zone_id": `"${alibabacloudstack_redis_tair_instances.default.ZoneId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,
			"zone_id": `"${alibabacloudstack_redis_tair_instances.default.ZoneId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_tair_instances.default.id}"]`,

			"architecture_type":  `"${alibabacloudstack_redis_tair_instances.default.ArchitectureType}"`,
			"engine_version":     `"${alibabacloudstack_redis_tair_instances.default.EngineVersion}"`,
			"instance_class":     `"${alibabacloudstack_redis_tair_instances.default.InstanceClass}"`,
			"instance_type":      `"${alibabacloudstack_redis_tair_instances.default.InstanceType}"`,
			"network_type":       `"${alibabacloudstack_redis_tair_instances.default.NetworkType}"`,
			"payment_type":       `"${alibabacloudstack_redis_tair_instances.default.PaymentType}"`,
			"status":             `"${alibabacloudstack_redis_tair_instances.default.Status}"`,
			"tair_instance_id":   `"${alibabacloudstack_redis_tair_instances.default.TairInstanceId}"`,
			"tair_instance_name": `"${alibabacloudstack_redis_tair_instances.default.TairInstanceName}"`,
			"vswitch_id":         `"${alibabacloudstack_redis_tair_instances.default.VSwitchId}"`,
			"vpc_id":             `"${alibabacloudstack_redis_tair_instances.default.VpcId}"`,
			"zone_id":            `"${alibabacloudstack_redis_tair_instances.default.ZoneId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_tair_instances.default.id}_fake"]`,

			"architecture_type":  `"${alibabacloudstack_redis_tair_instances.default.ArchitectureType}_fake"`,
			"engine_version":     `"${alibabacloudstack_redis_tair_instances.default.EngineVersion}_fake"`,
			"instance_class":     `"${alibabacloudstack_redis_tair_instances.default.InstanceClass}_fake"`,
			"instance_type":      `"${alibabacloudstack_redis_tair_instances.default.InstanceType}_fake"`,
			"network_type":       `"${alibabacloudstack_redis_tair_instances.default.NetworkType}_fake"`,
			"payment_type":       `"${alibabacloudstack_redis_tair_instances.default.PaymentType}_fake"`,
			"status":             `"${alibabacloudstack_redis_tair_instances.default.Status}_fake"`,
			"tair_instance_id":   `"${alibabacloudstack_redis_tair_instances.default.TairInstanceId}_fake"`,
			"tair_instance_name": `"${alibabacloudstack_redis_tair_instances.default.TairInstanceName}_fake"`,
			"vswitch_id":         `"${alibabacloudstack_redis_tair_instances.default.VSwitchId}_fake"`,
			"vpc_id":             `"${alibabacloudstack_redis_tair_instances.default.VpcId}_fake"`,
			"zone_id":            `"${alibabacloudstack_redis_tair_instances.default.ZoneId}_fake"`}),
	}

	AlibabacloudstackRedisTairInstancesCheckInfo.dataSourceTestCheck(t, rand, idsConf, architecture_typeConf, engine_versionConf, instance_classConf, instance_typeConf, network_typeConf, payment_typeConf, statusConf, tair_instance_idConf, tair_instance_nameConf, vswitch_idConf, vpc_idConf, zone_idConf, allConf)
}

var existAlibabacloudstackRedisTairInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":    "1",
		"instances.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackRedisTairInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#": "0",
	}
}

var AlibabacloudstackRedisTairInstancesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_redis_tair_instances.default",
	existMapFunc: existAlibabacloudstackRedisTairInstancesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackRedisTairInstancesMapFunc,
}

func testAccCheckAlibabacloudstackRedisTairInstancesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackRedisTairInstances%d"
}






data "alibabacloudstack_redis_tair_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
