package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckEIPAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_eip_association" {
			continue
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No EIP Association ID is set"))
		}

		// Try to find the EIP
		_, err := vpcService.DescribeEipAssociation(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func TestAccApsaraStackEipAssociationBasic(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "apsarastack_eip_association.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipAssociationBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociationConfigBaisc(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccApsaraStackEipAssociationMulti(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "apsarastack_eip_association.default.1"
	ra := resourceAttrInit(resourceId, testAccCheckEipAssociationBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociationConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccApsaraStackEipAssociationEni(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "apsarastack_eip_association.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipAssociationBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociationConfigEni(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "NetworkInterface",
					}),
				),
			},
		},
	})
}

func testAccEIPAssociationConfigBaisc(rand int) string {
	return fmt.Sprintf(`
%s

%s

%s
provider "apsarastack" {
	assume_role {}
}
variable "name" {
	default = "tf-testAccEipAssociation%d"
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_instance" "default" {
  vswitch_id = "${apsarastack_vswitch.default.id}"
  image_id = "${data.apsarastack_images.default.images.0.id}"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  system_disk_category = "cloud_ssd"
  instance_type = "${local.instance_type_id}"

  security_groups = ["${apsarastack_security_group.default.id}"]
  instance_name = "${var.name}"
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "apsarastack_eip" "default" {
	name = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
  allocation_id = "${apsarastack_eip.default.id}"
  instance_id = "${apsarastack_instance.default.id}"
}
`, DataApsarastackVswitchZones, DataApsarastackInstanceTypes, DataApsarastackImages, rand)
}

func testAccEIPAssociationConfigMulti(rand int) string {
	return fmt.Sprintf(`
%s

%s

%s
provider "apsarastack" {
	assume_role {}
}
variable "name" {
	default = "tf-testAccEipAssociation%d"
}

variable "number" {
		default = "2"
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_instance" "default" {
  count = "${var.number}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
  image_id = "${data.apsarastack_images.default.images.0.id}"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  system_disk_category = "cloud_ssd"
  instance_type = "${local.instance_type_id}"

  security_groups = ["${apsarastack_security_group.default.id}"]
  instance_name = "${var.name}"
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "apsarastack_eip" "default" {
	count = "${var.number}"
	name = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
  count = "${var.number}"
  allocation_id = "${element(apsarastack_eip.default.*.id,count.index)}"
  instance_id = "${element(apsarastack_instance.default.*.id,count.index)}"
}
`, DataApsarastackVswitchZones, DataApsarastackInstanceTypes, DataApsarastackImages, rand)
}

func testAccEIPAssociationConfigEni(rand int) string {
	return fmt.Sprintf(`
%s
provider "apsarastack" {
	assume_role {}
}
variable "name" {
  default = "tf-testAccEipAssociation%d"
}

resource "apsarastack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}


resource "apsarastack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_network_interface" "default" {
	name = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"
	security_groups = [ "${apsarastack_security_group.default.id}" ]
	private_ip = "192.168.0.2"
}

resource "apsarastack_eip" "default" {
	name = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
  allocation_id = "${apsarastack_eip.default.id}"
  instance_id = "${apsarastack_network_interface.default.id}"
  instance_type = "NetworkInterface"
}
`, DataApsarastackVswitchZones, rand)
}

var testAccCheckEipAssociationBasicMap = map[string]string{
	"allocation_id": CHECKSET,
	"instance_id":   CHECKSET,
}
