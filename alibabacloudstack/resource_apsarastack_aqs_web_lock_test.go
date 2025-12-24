package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAqsWebLock_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_aqs_web_lock.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AqsService{testYundunProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeWebLockInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string { return "" })

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: nil,
		Providers:    testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instanceid": "${alibabacloudstack_ecs_instance.default.id}",
					"status":     "on",
					"lock_configs": []map[string]interface{}{
						{
							"dir":                 "/test/tf",
							"local_backup_dir":    "/usr/local/aegis/bak1",
							"inclusive_file_type": "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg;gif;png;jspx",
							"defence_mode":        "block",
							"mode":                "whitelist",
						},
						{
							"dir":                 "/test2/tf",
							"local_backup_dir":    "/usr/local/aegis/bak2",
							"inclusive_file_type": "php;jsp;asp;aspx;js;cgi;",
							"defence_mode":        "block",
							"mode":                "whitelist",
						},
						{
							"dir":                 "/test3/tf",
							"local_backup_dir":    "/usr/local/aegis/bak3",
							"exclusive_file_type": "log;txt;ldb",
							"exclusive_file":      "aaa.txt",
							"exclusive_dir":       "testpath",
							"defence_mode":        "audit",
							"mode":                "blacklist",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":         "on",
						"lock_configs.#": "3",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "off",
					"lock_configs": []map[string]interface{}{
						{
							"dir":                 "/test/tf",
							"local_backup_dir":    "/usr/local/aegis/bak1",
							"inclusive_file_type": "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg;gif;png;jspx",
							"defence_mode":        "block",
							"mode":                "whitelist",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":         "off",
						"lock_configs.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "off",
					"lock_configs": []map[string]interface{}{
						{
							"dir":                 "/test2/tf",
							"local_backup_dir":    "/usr/local/aegis/bak2",
							"inclusive_file_type": "php;jsp;asp;aspx;js;cgi;",
							"defence_mode":        "block",
							"mode":                "whitelist",
						},
						{
							"dir":                 "/test3/tf",
							"local_backup_dir":    "/usr/local/aegis/bak3",
							"exclusive_file_type": "log;txt;ldb",
							"exclusive_file":      "aaa.txt",
							"exclusive_dir":       "testpath",
							"defence_mode":        "audit",
							"mode":                "blacklist",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":         "off",
						"lock_configs.#": "2",
					}),
				),
			},
		},
	})
}

func resourceAqsWebLockDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
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


`, name)
}
