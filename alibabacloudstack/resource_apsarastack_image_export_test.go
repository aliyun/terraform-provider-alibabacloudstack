package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackImageExport(t *testing.T) {
	var v ecs.Image
	resourceId := "alibabacloudstack_image_export.default"
	ra := resourceAttrInit(resourceId, testAccExportImageCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}

	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testaccecsimageexportconfigbasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageExportBasicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":   "${alibabacloudstack_image.default.id}",
					"oss_bucket": "${alibabacloudstack_oss_bucket.default.bucket}",
					"oss_prefix": "ecsExport",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_prefix": "ecsExport",
					}),
				),
			},
		},
	})
}

var testAccExportImageCheckMap = map[string]string{
	"image_id":   CHECKSET,
	"oss_bucket": CHECKSET,
}

func resourceImageExportBasicConfigDependence(name string) string {
	return fmt.Sprintf(`

%s

%s

%s

variable "name" {
	default = "%s"
}
resource "alibabacloudstack_vpc" "default" {
 name       = "${var.name}"
 cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
 vpc_id            = "${alibabacloudstack_vpc.default.id}"
 cidr_block        = "172.16.0.0/24"
 availability_zone = data.alibabacloudstack_zones.default.zones[0].id
 name              = "${var.name}"
}
resource "alibabacloudstack_security_group" "default" {
 name   = "${var.name}"
 vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_instance" "default" {
 image_id = "${data.alibabacloudstack_images.default.ids[0]}"
 instance_type = data.alibabacloudstack_instance_types.default.ids[0]
 security_groups = "${[alibabacloudstack_security_group.default.id]}"
 vswitch_id = "${alibabacloudstack_vswitch.default.id}"
 instance_name = "${var.name}"
}
resource "alibabacloudstack_image" "default" {
 instance_id = "${alibabacloudstack_instance.default.id}"
 image_name        = "${var.name}"
}
resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "${var.name}"
}
`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages, name)
}
