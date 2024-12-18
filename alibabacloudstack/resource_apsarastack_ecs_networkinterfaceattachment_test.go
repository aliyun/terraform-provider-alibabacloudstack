package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackNetworkInterfaceAttachmentBasic(t *testing.T) {
	var v ecs.NetworkInterfaceSet
	resourceId := "alibabacloudstack_network_interface_attachment.default"
	ra := resourceAttrInit(resourceId, testAccCheckNetworkInterfaceAttachmentCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceAttachmentConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func TestAccAlibabacloudStackNetworkInterfaceAttachmentMulti(t *testing.T) {
	var v ecs.NetworkInterfaceSet
	resourceId := "alibabacloudstack_network_interface_attachment.default.1"
	ra := resourceAttrInit(resourceId, testAccCheckNetworkInterfaceAttachmentCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceAttachmentConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckNetworkInterfaceAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_network_interface_Attachment" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NetworkInterface ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeNetworkInterfaceAttachment(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccNetworkInterfaceAttachmentConfigBasic = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes_Eni2 + DataAlibabacloudstackImages + `

variable "name" {
  default = "tf-testAccNetworkInterfaceAttachment"
}

resource "alibabacloudstack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

resource "alibabacloudstack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${reverse(data.alibabacloudstack_zones.default.zones).0.id}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "default" {
    availability_zone = "${reverse(data.alibabacloudstack_zones.default.zones).0.id}"
    security_groups = ["${alibabacloudstack_security_group.default.id}"]

    instance_type = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
    instance_name        = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    internet_max_bandwidth_out = 10
}

resource "alibabacloudstack_network_interface" "default" {
    name = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    security_groups = [ "${alibabacloudstack_security_group.default.id}" ]
}

resource "alibabacloudstack_network_interface_attachment" "default" {
    instance_id = "${alibabacloudstack_instance.default.id}"
    network_interface_id = "${alibabacloudstack_network_interface.default.id}"
}
`

const testAccNetworkInterfaceAttachmentConfigMulti = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes_Eni2 + DataAlibabacloudstackImages + `

variable "name" {
  default = "tf-testAccNetworkInterfaceAttachment"
}

variable "number" {
		default = "2"
	}

resource "alibabacloudstack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

resource "alibabacloudstack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${reverse(data.alibabacloudstack_zones.default.zones).0.id}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "default" {
	count = "${var.number}"
    availability_zone = "${reverse(data.alibabacloudstack_zones.default.zones).0.id}"
    security_groups = ["${alibabacloudstack_security_group.default.id}"]

    instance_type = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
    instance_name        = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    internet_max_bandwidth_out = 10
}

resource "alibabacloudstack_network_interface" "default" {
    count = "${var.number}"
    name = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    security_groups = [ "${alibabacloudstack_security_group.default.id}" ]
}

resource "alibabacloudstack_network_interface_attachment" "default" {
	count = "${var.number}"
    instance_id = "${element(alibabacloudstack_instance.default.*.id, count.index)}"
    network_interface_id = "${element(alibabacloudstack_network_interface.default.*.id, count.index)}"
}
`

var testAccCheckNetworkInterfaceAttachmentCheckMap = map[string]string{
	"instance_id":          CHECKSET,
	"network_interface_id": CHECKSET,
}
