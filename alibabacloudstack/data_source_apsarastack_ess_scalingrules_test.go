package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackEssScalingrulesDataSource(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(map[string]string{
			"ids": `["${alibabacloudstack_ess_scaling_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(map[string]string{
			"ids": `["${alibabacloudstack_ess_scaling_rule.default.id}_fake"]`,
		}),
	}

	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}_fake"`,
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"SimpleScalingRule"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"TargetTrackingScalingRule"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(map[string]string{
			"name_regex": `"${alibabacloudstack_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(map[string]string{
			"name_regex": `"${alibabacloudstack_ess_scaling_rule.default.scaling_rule_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(map[string]string{
			"ids":              `["${alibabacloudstack_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alibabacloudstack_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(map[string]string{
			"ids":              `["${alibabacloudstack_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alibabacloudstack_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
	}

	var existEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                  "1",
			"ids.#":                    "1",
			"names.#":                  "1",
			"rules.0.id":               CHECKSET,
			"rules.0.scaling_group_id": CHECKSET,
			"rules.0.name":             CHECKSET,
			"rules.0.type":             CHECKSET,
			"rules.0.cooldown":         CHECKSET,
			"rules.0.adjustment_type":  "QuantityChangeInCapacity",
			"rules.0.adjustment_value": "1",
			"rules.0.scaling_rule_ari": CHECKSET,
		}
	}

	var fakeEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var EssScalingrulesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_ess_scaling_rules.default",
		existMapFunc: existEssRecordsMapFunc,
		fakeMapFunc:  fakeEssRecordsMapFunc,
	}

	EssScalingrulesRecordsCheckInfo.dataSourceTestCheck(t, -1, idsConf, scalingGroupIdConf, typeConf, nameRegexConf, allConf)
}

func testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDataSourceEssScalingRules"
}

%s

%s

%s

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alibabacloudstack_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}
resource "alibabacloudstack_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
}
resource "alibabacloudstack_ess_scaling_rule" "default"{
	scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
	scaling_rule_name = "${var.name}"
	adjustment_type = "QuantityChangeInCapacity"
	adjustment_value = 1
}
data "alibabacloudstack_ess_scaling_rules" "default" {
  %s
}
`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages, strings.Join(pairs, "\n  "))
	return config
}
