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
			"ids": []string{"${alibabacloudstack_ecs_keypair.default.id}"},
			"finger_print": `"${alibabacloudstack_ecs_keypair.default.FingerPrint}"`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_keypair.default.id}"},
			"finger_print": `"${alibabacloudstack_ecs_keypair.default.FingerPrint}_fake"`,
		}),
	}

	key_pair_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"key_pair_name": `"${alibabacloudstack_ecs_keypair.default.KeyPairName}"`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"key_pair_name": `"${alibabacloudstack_ecs_keypair.default.KeyPairName}_fake"`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": `"${alibabacloudstack_ecs_keypair.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": `"${alibabacloudstack_ecs_keypair.default.ResourceGroupId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_keypair.default.id}"},
			"finger_print":      `"${alibabacloudstack_ecs_keypair.default.FingerPrint}"`,
			"key_pair_name":     `"${alibabacloudstack_ecs_keypair.default.KeyPairName}"`,
			"resource_group_id": `"${alibabacloudstack_ecs_keypair.default.ResourceGroupId}"`}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_keypair.default.id}_fake"},
			"finger_print":      `"${alibabacloudstack_ecs_keypair.default.FingerPrint}_fake"`,
			"key_pair_name":     `"${alibabacloudstack_ecs_keypair.default.KeyPairName}_fake"`,
			"resource_group_id": `"${alibabacloudstack_ecs_keypair.default.ResourceGroupId}_fake"`}),
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
