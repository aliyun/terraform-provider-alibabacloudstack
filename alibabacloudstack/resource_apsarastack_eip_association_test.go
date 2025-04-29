package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckEIPAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_eip_association" {
			continue
		}

		if rs.Primary.ID == "" {
			return errmsgs.WrapError(errmsgs.Error("No EIP Association ID is set"))
		}

		// Try to find the EIP
		_, err := vpcService.DescribeEipAssociation(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
	}

	return nil
}

func TestAccAlibabacloudStackEipAssociationBasic(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "alibabacloudstack_eip_association.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipAssociationBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000, 20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
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
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlibabacloudStackEipAssociationMulti(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "alibabacloudstack_eip_association.default.1"
	ra := resourceAttrInit(resourceId, testAccCheckEipAssociationBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000, 20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
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

func TestAccAlibabacloudStackEipAssociationEni(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "alibabacloudstack_eip_association.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipAssociationBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000, 20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
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
variable "name" {
	default = "tf-testAccEipAssociation%d"
}

%s

resource "alibabacloudstack_eip" "default" {
	name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
  allocation_id = "${alibabacloudstack_eip.default.id}"
  instance_id = "${alibabacloudstack_ecs_instance.default.id}"
}
`, rand, ECSInstanceCommonTestCase)
}

func testAccEIPAssociationConfigMulti(rand int) string {
	return fmt.Sprintf(`
%s

%s

%s
variable "name" {
	default = "tf-testAccEipAssociation%d"
}

variable "number" {
		default = "2"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "default" {
  count = "${var.number}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  system_disk_category = "cloud_ssd"
  instance_type = "${local.default_instance_type_id}"

  security_groups = ["${alibabacloudstack_security_group.default.id}"]
  instance_name = "${var.name}"
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "alibabacloudstack_eip" "default" {
	count = "${var.number}"
	name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
  count = "${var.number}"
  allocation_id = "${element(alibabacloudstack_eip.default.*.id,count.index)}"
  instance_id = "${element(alibabacloudstack_instance.default.*.id,count.index)}"
}
`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages, rand)
}

func testAccEIPAssociationConfigEni(rand int) string {
	return fmt.Sprintf(`
%s
variable "name" {
  default = "tf-testAccEipAssociation%d"
}

resource "alibabacloudstack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}


resource "alibabacloudstack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_network_interface" "default" {
	name = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	security_groups = [ "${alibabacloudstack_security_group.default.id}" ]
	private_ip = "192.168.0.2"
}

resource "alibabacloudstack_eip" "default" {
	name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
  allocation_id = "${alibabacloudstack_eip.default.id}"
  instance_id = "${alibabacloudstack_network_interface.default.id}"
  instance_type = "NetworkInterface"
}
`, DataAlibabacloudstackVswitchZones, rand)
}

var testAccCheckEipAssociationBasicMap = map[string]string{
	"allocation_id": CHECKSET,
	"instance_id":   CHECKSET,
}
