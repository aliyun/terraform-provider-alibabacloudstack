package alibabacloudstack

import (
	"fmt"

	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackImageSharePermission(t *testing.T) {
	var v *ecs.DescribeImageSharePermissionResponse
	resourceId := "alibabacloudstack_image_share_permission.default"
	ra := resourceAttrInit(resourceId, testAccImageSharePermissionCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageShareByImageId")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf-testAccEcsImageShareConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageSharePermissionConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithMultipleAccount(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":   "${alibabacloudstack_image.default.id}",
					"account_id": "123456789",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_id": CHECKSET,
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

var testAccImageSharePermissionCheckMap = map[string]string{
	"image_id": CHECKSET,
}

func resourceImageSharePermissionConfigDependence(name string) string {
	return fmt.Sprintf(`

%s

%s

%s

variable "name" {
	default = "%s"
}

data "alibabacloudstack_instance_types" "new" {
availability_zone = data.alibabacloudstack_zones.default.zones[0].id
cpu_core_count    = 2
memory_size       = 4
}

resource "alibabacloudstack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
    vpc_id = "${alibabacloudstack_vpc.default.id}"
    cidr_block = "172.16.0.0/16"
    availability_zone = data.alibabacloudstack_zones.default.zones.0.id
    name = "${var.name}"
}
resource "alibabacloudstack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "default" {
    image_id = "${data.alibabacloudstack_images.default.images.0.id}"
    instance_type = "${data.alibabacloudstack_instance_types.new.instance_types.0.id}"
    instance_name = "${var.name}"
    security_groups = "${alibabacloudstack_security_group.default.*.id}"
    internet_max_bandwidth_out = "10"
    availability_zone = data.alibabacloudstack_zones.default.zones.0.id
    system_disk_category = "cloud_efficiency"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_image" "default" {
 instance_id = "${alibabacloudstack_instance.default.id}"
 image_name        = "${var.name}"
}
`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages, name)
}
