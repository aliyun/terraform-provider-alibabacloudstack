package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackEcsStorageSetsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsStorageSetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_storage_sets.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsStorageSetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_storage_sets.default.id}_fake"]`,
		}),
	}

	storage_set_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsStorageSetsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_ecs_storage_sets.default.id}"]`,
			"storage_set_name": `"${alibabacloudstack_ecs_storage_sets.default.StorageSetName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsStorageSetsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_ecs_storage_sets.default.id}_fake"]`,
			"storage_set_name": `"${alibabacloudstack_ecs_storage_sets.default.StorageSetName}_fake"`,
		}),
	}

	zone_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsStorageSetsSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_ecs_storage_sets.default.id}"]`,
			"zone_id": `"${alibabacloudstack_ecs_storage_sets.default.ZoneId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsStorageSetsSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_ecs_storage_sets.default.id}_fake"]`,
			"zone_id": `"${alibabacloudstack_ecs_storage_sets.default.ZoneId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsStorageSetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_storage_sets.default.id}"]`,

			"storage_set_name": `"${alibabacloudstack_ecs_storage_sets.default.StorageSetName}"`,
			"zone_id":          `"${alibabacloudstack_ecs_storage_sets.default.ZoneId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsStorageSetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_storage_sets.default.id}_fake"]`,

			"storage_set_name": `"${alibabacloudstack_ecs_storage_sets.default.StorageSetName}_fake"`,
			"zone_id":          `"${alibabacloudstack_ecs_storage_sets.default.ZoneId}_fake"`}),
	}

	AlibabacloudstackEcsStorageSetsCheckInfo.dataSourceTestCheck(t, rand, idsConf, storage_set_nameConf, zone_idConf, allConf)
}

var existAlibabacloudstackEcsStorageSetsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sets.#":    "1",
		"sets.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsStorageSetsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sets.#": "0",
	}
}

var AlibabacloudstackEcsStorageSetsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_storage_sets.default",
	existMapFunc: existAlibabacloudstackEcsStorageSetsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsStorageSetsMapFunc,
}

func testAccCheckAlibabacloudstackEcsStorageSetsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsStorageSets%d"
}






data "alibabacloudstack_ecs_storage_sets" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
