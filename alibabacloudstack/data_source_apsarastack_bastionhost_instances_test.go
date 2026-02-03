package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackBastionhostInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_bastionhost_instances.default"
	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                        "1",
				"instances.#":                  "1",
				"instances.0.license_code":     CHECKSET,
				"instances.0.user_vswitch_id":  CHECKSET,
				"instances.0.highavailability": CHECKSET,
				"instances.0.vpc_id":           CHECKSET,
				"instances.0.disasterrecovery": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":          "0",
				"descriptions.#": "0",
			}
		},
		Providers: testYunDunProviders(),
	}

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("tf_testAcc%d", rand),
		dataSourceYundunBastionhostInstanceConfigDependency)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_bastionhost_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_bastionhost_instance.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_bastionhost_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_bastionhost_instance.default.description}-fake",
		}),
	}

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

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_bastionhost_instance.default.description}",
			"ids":               []string{"${alibabacloudstack_bastionhost_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_bastionhost_instance.default.description}-fake",
			"ids":               []string{"${alibabacloudstack_bastionhost_instance.default.id}-fake"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

// Dependence template generation method
func dataSourceYundunBastionhostInstanceConfigDependency(description string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}	
  
data "alibabacloudstack_zones" "default" {
provider = alibabacloudstack-common
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "vpc" {	
provider = alibabacloudstack-common
	vpc_name = var.name
	cidr_block = "192.168.0.0/16" # VPC CIDR block
}
resource "alibabacloudstack_vswitch" "vsw" {
provider = alibabacloudstack-common
	vpc_id = alibabacloudstack_vpc.vpc.id
	cidr_block = "192.168.0.0/16" # Subnet CIDR block
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id # Availability zone
}

resource "alibabacloudstack_bastionhost_instance" "default" {
description = var.name
  license_code = "bastionhostah_small_lic"
  vpc_id = "${alibabacloudstack_vpc.vpc.id}"
  vswitch_id = "${alibabacloudstack_vswitch.vsw.id}"
	highavailability= "false"
	disasterrecovery= "false"
	asset=            "50"
}
`, description)
}
