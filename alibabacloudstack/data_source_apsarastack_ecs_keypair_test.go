package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackEcsKeyPairsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_key_pairs.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_key_pairs.default.id}_fake"]`,
		}),
	}

	finger_printConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ecs_key_pairs.default.id}"]`,
			"finger_print": `"${alibabacloudstack_ecs_key_pairs.default.FingerPrint}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ecs_key_pairs.default.id}_fake"]`,
			"finger_print": `"${alibabacloudstack_ecs_key_pairs.default.FingerPrint}_fake"`,
		}),
	}

	key_pair_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ecs_key_pairs.default.id}"]`,
			"key_pair_name": `"${alibabacloudstack_ecs_key_pairs.default.KeyPairName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ecs_key_pairs.default.id}_fake"]`,
			"key_pair_name": `"${alibabacloudstack_ecs_key_pairs.default.KeyPairName}_fake"`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_key_pairs.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_key_pairs.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_key_pairs.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_key_pairs.default.ResourceGroupId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_key_pairs.default.id}"]`,

			"finger_print":      `"${alibabacloudstack_ecs_key_pairs.default.FingerPrint}"`,
			"key_pair_name":     `"${alibabacloudstack_ecs_key_pairs.default.KeyPairName}"`,
			"resource_group_id": `"${alibabacloudstack_ecs_key_pairs.default.ResourceGroupId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_key_pairs.default.id}_fake"]`,

			"finger_print":      `"${alibabacloudstack_ecs_key_pairs.default.FingerPrint}_fake"`,
			"key_pair_name":     `"${alibabacloudstack_ecs_key_pairs.default.KeyPairName}_fake"`,
			"resource_group_id": `"${alibabacloudstack_ecs_key_pairs.default.ResourceGroupId}_fake"`}),
	}

	AlibabacloudstackEcsKeyPairsCheckInfo.dataSourceTestCheck(t, rand, idsConf, finger_printConf, key_pair_nameConf, resource_group_idConf, allConf)
}

var existAlibabacloudstackEcsKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pairs.#":    "1",
		"pairs.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pairs.#": "0",
	}
}

var AlibabacloudstackEcsKeyPairsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_key_pairs.default",
	existMapFunc: existAlibabacloudstackEcsKeyPairsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsKeyPairsMapFunc,
}

func testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsKeyPairs%d"
}






data "alibabacloudstack_ecs_key_pairs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
