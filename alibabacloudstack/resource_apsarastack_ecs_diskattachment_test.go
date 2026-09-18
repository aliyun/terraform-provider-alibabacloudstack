package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackDiskAttachment_basic(t *testing.T) {
	resourceId := "alibabacloudstack_ecs_diskattachment.default"
	var i ecs.Instance
	var v ecs.Disk
	var attachment ecs.Disk
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}

	diskRc := resourceCheckInitWithDescribeMethod("alibabacloudstack_ecs_disk.default", &v, serverFunc, "DescribeDisk")
	instanceRc := resourceCheckInitWithDescribeMethod("alibabacloudstack_ecs_instance.default", &i, serverFunc, "DescribeInstance")
	attachmentRc := resourceCheckInitWithDescribeMethod("alibabacloudstack_ecs_diskattachment.default", &attachment, serverFunc, "DescribeDiskAttachment")

	name := "tf-testAccEcsDiskAttachmentConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEcsDiskAttachmentConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id":     "${alibabacloudstack_ecs_disk.default.id}",
					"instance_id": "${alibabacloudstack_ecs_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
					resource.TestCheckResourceAttrSet(
						"alibabacloudstack_ecs_diskattachment.default", "device_name"),
				),
			},
			{
				Config: resourceEcsDiskAttachmentResizeConfig("tf-testAccEcsDiskAttachmentConfig"),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
					resource.TestCheckResourceAttrSet(
						"alibabacloudstack_ecs_diskattachment.default", "device_name"),
					resource.TestCheckResourceAttr(
						"alibabacloudstack_ecs_disk.default", "size", "30"),
				),
			},
			{
				ResourceName:      "alibabacloudstack_ecs_diskattachment.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlibabacloudStackDiskAttachment_multi(t *testing.T) {
	resourceId := "alibabacloudstack_ecs_diskattachment.default.1"
	var i ecs.Instance
	var v ecs.Disk
	var attachment ecs.Disk
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	diskRc := resourceCheckInitWithDescribeMethod("alibabacloudstack_ecs_disk.default.1", &v, serverFunc, "DescribeDisk")
	instanceRc := resourceCheckInitWithDescribeMethod("alibabacloudstack_ecs_instance.default", &i, serverFunc, "DescribeInstance")
	attachmentRc := resourceCheckInitWithDescribeMethod("alibabacloudstack_ecs_diskattachment.default.1", &attachment, serverFunc, "DescribeDiskAttachment")

	name := "tf-testAccEcsDiskAttachmentConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEcsDiskAttachmentMultiConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "${var.number}",
					"disk_id":     "${element(alibabacloudstack_ecs_disk.default.*.id, count.index)}",
					"instance_id": "${alibabacloudstack_ecs_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
					resource.TestCheckResourceAttrSet(
						"alibabacloudstack_ecs_diskattachment.default.1", "device_name"),
				),
			},
		},
	})
}

func testAccCheckDiskAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_ecs_diskattachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeDiskAttachment(rs.Primary.ID)

		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
	}

	return nil
}

func resourceEcsDiskAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_ecs_disk" "default" {
	availability_zone = data.alibabacloudstack_zones.default.zones[0].id
	size              = "20"
	name              = "${var.name}"
	category          = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	delete_auto_snapshot = true
	tags = {
		Name = "TerraformTest-disk"
	}
}
`, name, ECSInstanceCommonTestCase)
}

func resourceEcsDiskAttachmentResizeConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_ecs_disk" "default" {
	availability_zone = data.alibabacloudstack_zones.default.zones[0].id
	size              = "30"
	name              = "${var.name}"
	category          = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	tags = {
		Name = "TerraformTest-disk"
	}
}

resource "alibabacloudstack_ecs_diskattachment" "default" {
	disk_id     = "${alibabacloudstack_ecs_disk.default.id}"
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
}
`, name, ECSInstanceCommonTestCase)
}

func resourceEcsDiskAttachmentMultiConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "number" {
	default = "2"
}

%s

resource "alibabacloudstack_ecs_disk" "default" {
	name              = "${var.name}-${count.index}"
	count             = "${var.number}"
	availability_zone = data.alibabacloudstack_zones.default.zones[0].id
	category          = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	size              = "20"
	tags = {
		Name = "TerraformTest-disk-${count.index}"
	}
}
`, name, ECSInstanceCommonTestCase)
}
