package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEvpcEvpcsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)

	name := fmt.Sprintf("tf_testevpc_%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(AlibabacloudstackEvpcEvpcsDataCheckInfo.resourceId, name, testAccCheckAlibabacloudstackEvpcEvpcsDataSourceConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_evpc_evpc.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_evpc_evpc.default.id}_fake"},
		}),
	}

	evpc_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${alibabacloudstack_evpc_evpc.default.id}"},
			"evpc_name": "${alibabacloudstack_evpc_evpc.default.evpc_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${alibabacloudstack_evpc_evpc.default.id}_fake"},
			"evpc_name": "${alibabacloudstack_evpc_evpc.default.evpc_name}_fake",
		}),
	}

	name_regexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_evpc_evpc.default.id}"},
			"name_regex": "${alibabacloudstack_evpc_evpc.default.evpc_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_evpc_evpc.default.id}_fake"},
			"name_regex": "${alibabacloudstack_evpc_evpc.default.evpc_name}_fake",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_evpc_evpc.default.id}"},
			"status": "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_evpc_evpc.default.id}"},
			"status": "Pending",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${alibabacloudstack_evpc_evpc.default.id}"},
			"evpc_name": "${alibabacloudstack_evpc_evpc.default.evpc_name}"},
		),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${alibabacloudstack_evpc_evpc.default.id}_fake"},
			"evpc_name": "${alibabacloudstack_evpc_evpc.default.evpc_name}_fake"}),
	}

	AlibabacloudstackEvpcEvpcsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, evpc_nameConf, name_regexConf, statusConf, allConf)
}

var existAlibabacloudstackEvpcEvpcsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"evpcs.#":         "1",
		"evpcs.0.evpc_id": CHECKSET,
	}
}

var fakeAlibabacloudstackEvpcEvpcsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"evpcs.#": "0",
	}
}

var AlibabacloudstackEvpcEvpcsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_evpc_evpcs.default",
	existMapFunc: existAlibabacloudstackEvpcEvpcsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEvpcEvpcsDataMapFunc,
}

func testAccCheckAlibabacloudstackEvpcEvpcsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_evpc_evpc" "default" {
	evpc_name = var.name
	description = var.name
}

`, name)
}
