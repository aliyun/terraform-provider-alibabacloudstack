package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEcsKeyPairsDataSource(t *testing.T) {

	resourceId := AlibabacloudstackEcsKeyPairsCheckInfo.resourceId
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcsKeyPair%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_keypair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_keypair.default.id}_fake"},
		}),
	}

	finger_printConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alibabacloudstack_ecs_keypair.default.id}"},
			"finger_print": "${alibabacloudstack_ecs_keypair.default.finger_print}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alibabacloudstack_ecs_keypair.default.id}"},
			"finger_print": "${alibabacloudstack_ecs_keypair.default.finger_print}_fake",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ecs_keypair.default.key_pair_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ecs_keypair.default.key_pair_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alibabacloudstack_ecs_keypair.default.id}"},
			"finger_print": "${alibabacloudstack_ecs_keypair.default.finger_print}",
			"name_regex":   "${alibabacloudstack_ecs_keypair.default.key_pair_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alibabacloudstack_ecs_keypair.default.id}_fake"},
			"finger_print": "${alibabacloudstack_ecs_keypair.default.finger_print}_fake",
			"name_regex":   "${alibabacloudstack_ecs_keypair.default.key_pair_name}_fake",
		}),
	}

	AlibabacloudstackEcsKeyPairsCheckInfo.dataSourceTestCheck(t, rand, idsConf, finger_printConf, nameRegexConf, allConf)
}

var existAlibabacloudstackEcsKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"key_pairs.#":    "1",
		"key_pairs.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"key_pairs.#": "0",
	}
}

var AlibabacloudstackEcsKeyPairsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_keypairs.default",
	existMapFunc: existAlibabacloudstackEcsKeyPairsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsKeyPairsMapFunc,
}

func testAccCheckAlibabacloudstackEcsKeyPairsSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_ecs_keypair" "default"{
	key_pair_name = var.name
}


`, name)
}
