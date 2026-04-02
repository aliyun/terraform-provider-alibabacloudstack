package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackBcmpKeyPairsDataSource(t *testing.T) {

	resourceId := AlibabacloudstackBcmpKeyPairsCheckInfo.resourceId
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccBcmpKeyPair%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudstackBcmpKeyPairsSourceConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_bcmp_keypair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_bcmp_keypair.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_bcmp_keypair.default.key_pair_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_bcmp_keypair.default.key_pair_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":      []string{"${alibabacloudstack_bcmp_keypair.default.id}"},
			"name_regex": "${alibabacloudstack_bcmp_keypair.default.key_pair_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":      []string{"${alibabacloudstack_bcmp_keypair.default.id}_fake"},
			"name_regex": "${alibabacloudstack_bcmp_keypair.default.key_pair_name}_fake",
		}),
	}

	AlibabacloudstackBcmpKeyPairsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

var existAlibabacloudstackBcmpKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"key_pairs.#":    "1",
		"key_pairs.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackBcmpKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"key_pairs.#": "0",
	}
}

var AlibabacloudstackBcmpKeyPairsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_bcmp_keypairs.default",
	existMapFunc: existAlibabacloudstackBcmpKeyPairsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackBcmpKeyPairsMapFunc,
}

func testAccCheckAlibabacloudstackBcmpKeyPairsSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_bcmp_keypair" "default"{
	key_pair_name = var.name
}

`, name)
}
