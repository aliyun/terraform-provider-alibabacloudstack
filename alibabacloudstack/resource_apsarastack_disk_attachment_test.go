package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackDiskAttachment(t *testing.T) {
	var i ecs.Instance
	var v ecs.Disk
	var attachment ecs.Disk
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	diskRc := resourceCheckInit("alibabacloudstack_disk.default", &v, serverFunc)

	instanceRc := resourceCheckInit("alibabacloudstack_instance.default", &i, serverFunc)

	attachmentRc := resourceCheckInit("alibabacloudstack_disk_attachment.default", &attachment, serverFunc)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alibabacloudstack_disk_attachment.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskAttachmentConfig(),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
					resource.TestCheckResourceAttrSet(
						"alibabacloudstack_disk_attachment.default", "device_name"),
				),
			},
			{
				Config: testAccDiskAttachmentConfigResize(),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
					resource.TestCheckResourceAttrSet(
						"alibabacloudstack_disk_attachment.default", "device_name"),
					resource.TestCheckResourceAttr(
						"alibabacloudstack_disk.default", "size", "70"),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackDiskMultiAttachment(t *testing.T) {
	var i ecs.Instance
	var v ecs.Disk
	var attachment ecs.Disk
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	diskRc := resourceCheckInit("alibabacloudstack_disk.default.1", &v, serverFunc)

	instanceRc := resourceCheckInit("alibabacloudstack_instance.default", &i, serverFunc)

	attachmentRc := resourceCheckInit("alibabacloudstack_disk_attachment.default.1", &attachment, serverFunc)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alibabacloudstack_disk_attachment.default.1",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMultiDiskAttachmentConfig(EcsInstanceCommonNoZonesTestCase),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
					resource.TestCheckResourceAttrSet(
						"alibabacloudstack_disk_attachment.default.1", "device_name"),
				),
			},
		},
	})

}

func testAccCheckDiskAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_disk_attachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeDiskAttachment(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func testAccDiskAttachmentConfig() string {
	return fmt.Sprintf(`
	%s
	
	%s
	
    data "alibabacloudstack_images" "default" {
	  # test for windows service
      name_regex  = "^win*"

      most_recent = true
      owners      = "system"
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
   
	variable "name" {
		default = "tf-testAccEcsDiskAttachmentConfig"
	}

	resource "alibabacloudstack_disk" "default" {
	  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
	  size = "50"
	  name = "${var.name}"
	  category = "cloud_efficiency"

	  tags = {
	    Name = "TerraformTest-disk"
	  }
	}

	resource "alibabacloudstack_instance" "default" {
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		availability_zone = data.alibabacloudstack_zones.default.zones[0].id
		system_disk_category = "cloud_ssd"
		system_disk_size = 40
		instance_type = "${local.instance_type_id}"
		security_groups = ["${alibabacloudstack_security_group.default.id}"]
		instance_name = "${var.name}"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}

	resource "alibabacloudstack_disk_attachment" "default" {
	  disk_id = "${alibabacloudstack_disk.default.id}"
	  instance_id = "${alibabacloudstack_instance.default.id}"
	}
	`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes)
}
func testAccDiskAttachmentConfigResize() string {
	return fmt.Sprintf(`
	%s
	
	%s

    data "alibabacloudstack_images" "default" {
	  # test for windows service
      name_regex  = "^win*"

      most_recent = true
      owners      = "system"
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
    
	variable "name" {
		default = "tf-testAccEcsDiskAttachmentConfig"
	}

	resource "alibabacloudstack_disk" "default" {
	  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
	  size = "70"
	  name = "${var.name}"
	  category = "cloud_efficiency"

	  tags = {
	    Name = "TerraformTest-disk"
	  }
	}

	resource "alibabacloudstack_instance" "default" {
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		availability_zone = data.alibabacloudstack_zones.default.zones[0].id
		system_disk_category = "cloud_ssd"
		system_disk_size = 40
		instance_type = "${local.instance_type_id}"
		security_groups = ["${alibabacloudstack_security_group.default.id}"]
		instance_name = "${var.name}"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}

	resource "alibabacloudstack_disk_attachment" "default" {
	  disk_id = "${alibabacloudstack_disk.default.id}"
	  instance_id = "${alibabacloudstack_instance.default.id}"
	}
	`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes)
}
func testAccMultiDiskAttachmentConfig(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEcsDiskAttachmentConfig"
	}

	variable "number" {
		default = "2"
	}

	resource "alibabacloudstack_disk" "default" {
		name = "${var.name}-${count.index}"
		count = "${var.number}"
		availability_zone = data.alibabacloudstack_zones.default.zones[0].id
		size = "50"

		tags = {
			Name = "TerraformTest-disk-${count.index}"
		}
	}

	resource "alibabacloudstack_instance" "default" {
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		availability_zone = data.alibabacloudstack_zones.default.zones[0].id
		system_disk_category = "cloud_ssd"
		system_disk_size = 40
		instance_type = "${local.instance_type_id}"
		security_groups = ["${alibabacloudstack_security_group.default.id}"]
		instance_name = "${var.name}"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}

	resource "alibabacloudstack_disk_attachment" "default" {
		count = "${var.number}"
		disk_id     = "${element(alibabacloudstack_disk.default.*.id, count.index)}"
		instance_id = "${alibabacloudstack_instance.default.id}"
	}
	`, common)
}
