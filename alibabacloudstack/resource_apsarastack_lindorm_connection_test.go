package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackLindormConnection_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_lindorm_connection.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeLindormConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-lindorm-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, LindormConnectionCommonTestCase)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${local.lindorm_instance_id}",
					"engine_type": "lindorm",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_username":                   CHECKSET,
						"engine_password":                   CHECKSET,
						"net_info_list.0.connection_string": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"classic_network": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"classic_network": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"classic_network": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"classic_network": "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func LindormConnectionCommonTestCase(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "lindorm_instance_id" {
  type    = string
  default = "%s"
}

%s

data "alibabacloudstack_lindorm_instances" "default" {
  ids = var.lindorm_instance_id != "" ? [var.lindorm_instance_id] : [" "]
}

data "alibabacloudstack_lindorm_instance_types" "sortbycpu" {
	sorted_by = "CPU"
	engine_type = "lindorm"
}

resource "alibabacloudstack_lindorm_instance" "default" {
	count = length(data.alibabacloudstack_lindorm_instances.default.ids) > 0 ? 0 : 1
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	instance_alias = "${var.name}"
	cpu_brand = "Intel"
	disk_category = "HDD"
	local_disk_size = "200"
	engine_type = "lindorm"
	instance_type = "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[0].name}"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	lindorm_num = 2
}

locals {
  lindorm_instance_id = length(data.alibabacloudstack_lindorm_instances.default.ids) > 0 ? data.alibabacloudstack_lindorm_instances.default.ids.0 : alibabacloudstack_lindorm_instance.default.0.id
}

`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_LINDORM_INSTANCE_ID"), VSwitchCommonTestCase)
}
