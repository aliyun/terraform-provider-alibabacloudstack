package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

// Generate dependency resources template for alibabacloudstack_aqs_web_locks
func resourceAqsWebLocksDependenceNew(rand int, attrMap map[string]string) string {
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

resource "alibabacloudstack_aqs_web_lock" "default" {
  instanceid = "${alibabacloudstack_ecs_instance.default.id}"
  status = "on"
  lock_configs {
	dir = "/test/tf/"
	local_backup_dir = "/usr/local/aegis/bak1"
	inclusive_file_type = "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg;gif;png;jspx"
	defence_mode = "block"
	mode = "whitelist"
  }
}

data "alibabacloudstack_aqs_web_locks" "default" {
    %s
}
`, rand, strings.Join(pairs, "\n   "))
}

func TestAccAlibabacloudStackAqsWebLocksDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_aqs_web_locks.default"
	rand := getAccTestRandInt(10000, 20000)
	testCheck := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"weblocks.#":                                    "1",
				"weblocks.0.uuid":                               CHECKSET,
				"weblocks.0.lock_configs.#":                     "1",
				"weblocks.0.lock_configs.0.dir":                 "/test/tf/",
				"weblocks.0.lock_configs.0.inclusive_file_type": "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg;gif;png;jspx",
				"weblocks.0.lock_configs.0.defence_mode":        "block",
				"weblocks.0.lock_configs.0.mode":                "whitelist",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"weblocks.#": "0",
			}
		},
		Providers: testYunDunProviders(),
	}

	testCheck.dataSourceTestCheck(t, rand,

		dataSourceTestAccConfig{
			existConfig: resourceAqsWebLocksDependenceNew(rand, map[string]string{
				"ids": `["${alibabacloudstack_aqs_web_lock.default.id}"]`,
			}),
			fakeConfig: resourceAqsWebLocksDependenceNew(rand, map[string]string{
				"ids": `["${alibabacloudstack_aqs_web_lock.default.id}_fake"]`,
			}),
		},
		// Test with name_regex filter
		dataSourceTestAccConfig{
			existConfig: resourceAqsWebLocksDependenceNew(rand, map[string]string{
				"instanceid": `"${alibabacloudstack_aqs_web_lock.default.instanceid}"`,
			}),
			fakeConfig: resourceAqsWebLocksDependenceNew(rand, map[string]string{
				"instanceid": `"${alibabacloudstack_aqs_web_lock.default.instanceid}_fake"`,
			}),
		},

		dataSourceTestAccConfig{
			existConfig: resourceAqsWebLocksDependenceNew(rand, map[string]string{
				"ids":        `["${alibabacloudstack_aqs_web_lock.default.id}"]`,
				"instanceid": `"${alibabacloudstack_aqs_web_lock.default.instanceid}"`,
			}),
			fakeConfig: resourceAqsWebLocksDependenceNew(rand, map[string]string{
				"name_regex": `"${alibabacloudstack_aqs_web_lock.default.name}_fake"`,
				"instanceid": `"${alibabacloudstack_aqs_web_lock.default.instanceid}_fake"`,
			}),
		},
	)
}
