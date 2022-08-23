package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackImageBasic(t *testing.T) {
	var v ecs.Image

	resourceId := "apsarastack_image.default"
	ra := resourceAttrInit(resourceId, testAccImageCheckMap)

	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsImageConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageBasicConfigDependence)
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
					"instance_id": "${apsarastack_instance.default.id}",
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":  name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":   name,
						"description":  fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescriptionChange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescriptionChange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
					}),
				),
			},
			//			{
			//				Config: testAccConfig(map[string]interface{}{
			//					"tags": map[string]string{
			//						"Created": "TF1",
			//						"For":     "acceptance test1232",
			//					},
			//				}),
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{
			//						"tags.%":       "2",
			//						"tags.Created": "TF1",
			//						"tags.For":     "acceptance test1232",
			//					}),
			//				),
			//			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"image_name":  name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"image_name":   name,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test123",
					}),
				),
			},
		},
	})
}

var testAccImageCheckMap = map[string]string{}

func resourceImageBasicConfigDependence(name string) string {
	return fmt.Sprintf(`

%s

%s

%s

variable "name" {
	default = "%s"
}
resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_instance" "default" {
  image_id         = "${data.apsarastack_images.default.ids[0]}"
  instance_type    = "${data.apsarastack_instance_types.default.ids[0]}"
  security_groups  = "${[apsarastack_security_group.default.id]}"
  vswitch_id       = "${apsarastack_vswitch.default.id}"
  instance_name    = "${var.name}"
  system_disk_size = 20
}
`, DataApsarastackVswitchZones, DataApsarastackInstanceTypes, DataApsarastackImages, name)
}
