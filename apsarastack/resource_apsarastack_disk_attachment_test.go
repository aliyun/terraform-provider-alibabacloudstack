package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApsaraStackDiskAttachment(t *testing.T) {
	var i ecs.Instance
	var v ecs.Disk
	var attachment ecs.Disk
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	diskRc := resourceCheckInit("apsarastack_disk.default", &v, serverFunc)

	instanceRc := resourceCheckInit("apsarastack_instance.default", &i, serverFunc)

	attachmentRc := resourceCheckInit("apsarastack_disk_attachment.default", &attachment, serverFunc)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "apsarastack_disk_attachment.default",
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
						"apsarastack_disk_attachment.default", "device_name"),
				),
			},
			{
				Config: testAccDiskAttachmentConfigResize(),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
					resource.TestCheckResourceAttrSet(
						"apsarastack_disk_attachment.default", "device_name"),
					resource.TestCheckResourceAttr(
						"apsarastack_disk.default", "size", "70"),
				),
			},
		},
	})

}

func TestAccApsaraStackDiskMultiAttachment(t *testing.T) {
	var i ecs.Instance
	var v ecs.Disk
	var attachment ecs.Disk
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	diskRc := resourceCheckInit("apsarastack_disk.default.1", &v, serverFunc)

	instanceRc := resourceCheckInit("apsarastack_instance.default", &i, serverFunc)

	attachmentRc := resourceCheckInit("apsarastack_disk_attachment.default.1", &attachment, serverFunc)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "apsarastack_disk_attachment.default.1",
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
						"apsarastack_disk_attachment.default.1", "device_name"),
				),
			},
		},
	})

}

func testAccCheckDiskAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_disk_attachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
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
	
    data "apsarastack_images" "default" {
	  # test for windows service
      name_regex  = "^win*"

      most_recent = true
      owners      = "system"
    }
    resource "apsarastack_vpc" "default" {
      name       = "${var.name}"
      cidr_block = "172.16.0.0/16"
    }
    resource "apsarastack_vswitch" "default" {
      vpc_id            = "${apsarastack_vpc.default.id}"
      cidr_block        = "172.16.0.0/24"
      availability_zone = data.apsarastack_zones.default.zones[0].id
      name              = "${var.name}"
    }
    resource "apsarastack_security_group" "default" {
      name   = "${var.name}"
      vpc_id = "${apsarastack_vpc.default.id}"
    }
   
	variable "name" {
		default = "tf-testAccEcsDiskAttachmentConfig"
	}

	resource "apsarastack_disk" "default" {
	  availability_zone = data.apsarastack_zones.default.zones[0].id
	  size = "50"
	  name = "${var.name}"
	  category = "cloud_efficiency"

	  tags = {
	    Name = "TerraformTest-disk"
	  }
	}

	resource "apsarastack_instance" "default" {
		image_id = "${data.apsarastack_images.default.images.0.id}"
		availability_zone = data.apsarastack_zones.default.zones[0].id
		system_disk_category = "cloud_ssd"
		system_disk_size = 40
		instance_type = "${local.instance_type_id}"
		security_groups = ["${apsarastack_security_group.default.id}"]
		instance_name = "${var.name}"
		vswitch_id = "${apsarastack_vswitch.default.id}"
	}

	resource "apsarastack_disk_attachment" "default" {
	  disk_id = "${apsarastack_disk.default.id}"
	  instance_id = "${apsarastack_instance.default.id}"
	}
	`, DataApsarastackVswitchZones, DataApsarastackInstanceTypes)
}
func testAccDiskAttachmentConfigResize() string {
	return fmt.Sprintf(`
	%s
	
	%s

    data "apsarastack_images" "default" {
	  # test for windows service
      name_regex  = "^win*"

      most_recent = true
      owners      = "system"
    }
    resource "apsarastack_vpc" "default" {
      name       = "${var.name}"
      cidr_block = "172.16.0.0/16"
    }
    resource "apsarastack_vswitch" "default" {
      vpc_id            = "${apsarastack_vpc.default.id}"
      cidr_block        = "172.16.0.0/24"
      availability_zone = data.apsarastack_zones.default.zones[0].id
      name              = "${var.name}"
    }
    resource "apsarastack_security_group" "default" {
      name   = "${var.name}"
      vpc_id = "${apsarastack_vpc.default.id}"
    }
    
	variable "name" {
		default = "tf-testAccEcsDiskAttachmentConfig"
	}

	resource "apsarastack_disk" "default" {
	  availability_zone = data.apsarastack_zones.default.zones[0].id
	  size = "70"
	  name = "${var.name}"
	  category = "cloud_efficiency"

	  tags = {
	    Name = "TerraformTest-disk"
	  }
	}

	resource "apsarastack_instance" "default" {
		image_id = "${data.apsarastack_images.default.images.0.id}"
		availability_zone = data.apsarastack_zones.default.zones[0].id
		system_disk_category = "cloud_ssd"
		system_disk_size = 40
		instance_type = "${local.instance_type_id}"
		security_groups = ["${apsarastack_security_group.default.id}"]
		instance_name = "${var.name}"
		vswitch_id = "${apsarastack_vswitch.default.id}"
	}

	resource "apsarastack_disk_attachment" "default" {
	  disk_id = "${apsarastack_disk.default.id}"
	  instance_id = "${apsarastack_instance.default.id}"
	}
	`, DataApsarastackVswitchZones, DataApsarastackInstanceTypes)
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

	resource "apsarastack_disk" "default" {
		name = "${var.name}-${count.index}"
		count = "${var.number}"
		availability_zone = data.apsarastack_zones.default.zones[0].id
		size = "50"

		tags = {
			Name = "TerraformTest-disk-${count.index}"
		}
	}

	resource "apsarastack_instance" "default" {
		image_id = "${data.apsarastack_images.default.images.0.id}"
		availability_zone = data.apsarastack_zones.default.zones[0].id
		system_disk_category = "cloud_ssd"
		system_disk_size = 40
		instance_type = "${local.instance_type_id}"
		security_groups = ["${apsarastack_security_group.default.id}"]
		instance_name = "${var.name}"
		vswitch_id = "${apsarastack_vswitch.default.id}"
	}

	resource "apsarastack_disk_attachment" "default" {
		count = "${var.number}"
		disk_id     = "${element(apsarastack_disk.default.*.id, count.index)}"
		instance_id = "${apsarastack_instance.default.id}"
	}
	`, common)
}
