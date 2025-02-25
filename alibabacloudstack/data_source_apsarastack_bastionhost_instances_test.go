package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackBastionhostInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_bastionhost_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("tf_testAcc%d", rand),
		dataSourceYundunBastionhostInstanceConfigDependency)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_bastionhost_instances.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_bastionhost_instances.default.id}-fake"},
		}),
	}

	// nameRegexConf := dataSourceTestAccConfig{
	// 	existConfig: testAccConfig(map[string]interface{}{
	// 		"description_regex": "${alibabacloudstack_bastionhost_instances.default.description}",
	// 	}),
	// 	fakeConfig: testAccConfig(map[string]interface{}{
	// 		"description_regex": "${alibabacloudstack_bastionhost_instances.default.description}-fake",
	// 	}),
	// }

	// tagsConf := dataSourceTestAccConfig{
	// 	existConfig: testAccConfig(map[string]interface{}{
	// 		"ids": []string{"${alibabacloudstack_bastionhost_instances.default.id}"},
	// 		"tags": map[string]interface{}{
	// 			"Created": "TF",
	// 		},
	// 	}),
	// 	fakeConfig: testAccConfig(map[string]interface{}{
	// 		"ids": []string{"${alibabacloudstack_bastionhost_instances.default.id}-fake"},
	// 		"tags": map[string]interface{}{
	// 			"Created": "TF-fake",
	// 		},
	// 	}),
	// }

	// allConf := dataSourceTestAccConfig{
	// 	existConfig: testAccConfig(map[string]interface{}{
	// 		"description_regex": "${alibabacloudstack_bastionhost_instances.default.description}",
	// 		"ids":               []string{"${alibabacloudstack_bastionhost_instances.default.id}"},
	// 		"tags": map[string]interface{}{
	// 			"For": "acceptance test",
	// 		},
	// 	}),
	// 	fakeConfig: testAccConfig(map[string]interface{}{
	// 		"description_regex": "${alibabacloudstack_bastionhost_instances.default.description}-fake",
	// 		"ids":               []string{"${alibabacloudstack_bastionhost_instances.default.id}-fake"},
	// 		"tags": map[string]interface{}{
	// 			"For": "acceptance test-fake",
	// 		},
	// 	}),
	// }

	var existYundunBastionhostInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "1",
			// "descriptions.#":                    "1",
			// "ids.0":                             CHECKSET,
			// "descriptions.0":                    fmt.Sprintf("tf_testAcc%d", rand),
			// "instances.#":                       "1",
			// "instances.0.description":           fmt.Sprintf("tf_testAcc%d", rand),
			// "instances.0.license_code":          "bhah_ent_50_asset",
			// "instances.0.user_vswitch_id":       CHECKSET,
			// "instances.0.public_network_access": "true",
			// "instances.0.private_domain":        CHECKSET,
			// "instances.0.instance_status":       CHECKSET,
			// "instances.0.security_group_ids.#":  "1",
		}
	}
	var fakeYundunBastionhostInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"descriptions.#": "0",
		}
	}
	var yundunBastionhostInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_bastionhost_instances.default",
		existMapFunc: existYundunBastionhostInstanceMapFunc,
		fakeMapFunc:  fakeYundunBastionhostInstanceMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}

	yundunBastionhostInstanceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)

}

func dataSourceYundunBastionhostInstanceConfigDependency(description string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}				
// data "alibabacloudstack_bastionhost_instances" "default" {
//   }
`, description)
}
