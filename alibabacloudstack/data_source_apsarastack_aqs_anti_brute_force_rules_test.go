package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

// Generate dependency resources template for alibabacloudstack_aqs_anti_brute_force_rules
func resourceAqsAntiBruteForceRulesDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc%d"
}

data "alibabacloudstack_zones" "default" {
  provider = alibabacloudstack-common
  available_resource_creation = "VSwitch"
  enable_details             = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  provider = alibabacloudstack-common
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  provider = alibabacloudstack-common
  name        = "${var.name}_vsw"
  vpc_id      = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block  = "172.16.0.0/24"
  zone_id     = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  provider = alibabacloudstack-common
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

data "alibabacloudstack_images" "default" {
  provider = alibabacloudstack-common
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  provider = alibabacloudstack-common
  sorted_by        = "Memory"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_ecs_instance" "default" {
  provider = alibabacloudstack-common
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  instance_name        = "${var.name}"
  user_data            = "I_am_user_data"
  security_groups      = ["${alibabacloudstack_ecs_securitygroup.default.id}"]
  vswitch_id          = "${alibabacloudstack_vpc_vswitch.default.id}"
  image_id            = "${data.alibabacloudstack_images.default.images.0.id}"
  security_enhancement_strategy = "Active"
  instance_type       = "${data.alibabacloudstack_instance_types.all.instance_types.0.id}"
  availability_zone   = "${data.alibabacloudstack_zones.default.zones[0].id}"
}

resource "alibabacloudstack_aqs_anti_brute_force_rule" "default" {
  provider = alibabacloudstack-common
  name = "${var.name}_rule"
  span = 10
  fail_count = 80
  forbidden_time = 360
  default_rule = true
  instance_ids = [
    alibabacloudstack_ecs_instance.default.id,
  ]
}

data "alibabacloudstack_aqs_anti_brute_force_rules" "default" {
    %s
}
`, rand, strings.Join(pairs, "\n   "))
}

func TestAccAlibabacloudStackAqsAntiBruteForceRulesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_aqs_anti_brute_force_rules.default"
	rand := getAccTestRandInt(10000, 20000)
	testCheck := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"rules.#":                "1",
				"rules.0.name":           fmt.Sprintf("tf-testacc%d_rule", rand),
				"rules.0.span":           "10",
				"rules.0.fail_count":     "80",
				"rules.0.forbidden_time": "360",
				"rules.0.default_rule":   "true",
				"rules.0.instance_ids.#": "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"rules.#": "0",
			}
		},
		Providers: testYunDunProviders(),
	}

	testCheck.dataSourceTestCheck(t, rand,

		dataSourceTestAccConfig{
			existConfig: resourceAqsAntiBruteForceRulesDependenceNew(rand, map[string]string{
				"ids": `["${alibabacloudstack_aqs_anti_brute_force_rule.default.id}"]`,
			}),
			fakeConfig: resourceAqsAntiBruteForceRulesDependenceNew(rand, map[string]string{
				"ids": `["${alibabacloudstack_aqs_anti_brute_force_rule.default.id}_fake"]`,
			}),
		},
		// Test with name_regex filter
		dataSourceTestAccConfig{
			existConfig: resourceAqsAntiBruteForceRulesDependenceNew(rand, map[string]string{
				"name_regex": `"${alibabacloudstack_aqs_anti_brute_force_rule.default.name}"`,
			}),
			fakeConfig: resourceAqsAntiBruteForceRulesDependenceNew(rand, map[string]string{
				"name_regex": `"${alibabacloudstack_aqs_anti_brute_force_rule.default.name}_fake"`,
			}),
		},

		dataSourceTestAccConfig{
			existConfig: resourceAqsAntiBruteForceRulesDependenceNew(rand, map[string]string{
				"ids":        `["${alibabacloudstack_aqs_anti_brute_force_rule.default.id}"]`,
				"name_regex": `"${alibabacloudstack_aqs_anti_brute_force_rule.default.name}"`,
			}),
			fakeConfig: resourceAqsAntiBruteForceRulesDependenceNew(rand, map[string]string{
				"name_regex": `"${alibabacloudstack_aqs_anti_brute_force_rule.default.name}_fake"`,
				"ids":        `["${alibabacloudstack_aqs_anti_brute_force_rule.default.id}_fake"]`,
			}),
		},
	)
}
