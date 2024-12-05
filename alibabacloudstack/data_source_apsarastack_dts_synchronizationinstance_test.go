package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackDtsSynchronizationInstancesDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDtsSynchronizationInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_dts_synchronization_instances.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDtsSynchronizationInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_dts_synchronization_instances.default.id}_fake"]`,
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDtsSynchronizationInstancesSourceConfig(rand, map[string]string{
			"ids":  `["${alibabacloudstack_dts_synchronization_instances.default.id}"]`,
			"type": `"${alibabacloudstack_dts_synchronization_instances.default.Type}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDtsSynchronizationInstancesSourceConfig(rand, map[string]string{
			"ids":  `["${alibabacloudstack_dts_synchronization_instances.default.id}_fake"]`,
			"type": `"${alibabacloudstack_dts_synchronization_instances.default.Type}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDtsSynchronizationInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_dts_synchronization_instances.default.id}"]`,

			"type": `"${alibabacloudstack_dts_synchronization_instances.default.Type}"`}),
		fakeConfig: testAccCheckAlibabacloudstackDtsSynchronizationInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_dts_synchronization_instances.default.id}_fake"]`,

			"type": `"${alibabacloudstack_dts_synchronization_instances.default.Type}_fake"`}),
	}

	AlibabacloudstackDtsSynchronizationInstancesCheckInfo.dataSourceTestCheck(t, rand, idsConf, typeConf, allConf)
}

var existAlibabacloudstackDtsSynchronizationInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":    "1",
		"instances.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackDtsSynchronizationInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#": "0",
	}
}

var AlibabacloudstackDtsSynchronizationInstancesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_dts_synchronization_instances.default",
	existMapFunc: existAlibabacloudstackDtsSynchronizationInstancesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackDtsSynchronizationInstancesMapFunc,
}

func testAccCheckAlibabacloudstackDtsSynchronizationInstancesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackDtsSynchronizationInstances%d"
}






data "alibabacloudstack_dts_synchronization_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
