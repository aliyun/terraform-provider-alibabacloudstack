package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckSnatEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_snat_entry" {
			continue
		}

		_, err := vpcService.DescribeSnatEntry(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(Error("Snat entry still exist"))
	}

	return nil
}

func TestAccApsaraStackSnatEntryBasic(t *testing.T) {
	var v vpc.SnatTableEntry

	resourceId := "apsarastack_snat_entry.default"
	ra := resourceAttrInit(resourceId, testAccCheckSnatEntryBasicMap)
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

		IDRefreshName: "apsarastack_snat_entry.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSnatEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnatEntryConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSnatEntryConfig_snatIp(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func TestAccApsaraStackSnatEntryMulti(t *testing.T) {
	var v vpc.SnatTableEntry

	resourceId := "apsarastack_snat_entry.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckSnatEntryBasicMap)
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

		IDRefreshName: "apsarastack_snat_entry.default.9",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSnatEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnatEntryConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccSnatEntryConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccSnatEntryConfig%d"
}
data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
	depends_on = [apsarastack_vpc.default]
	vpc_id = "${apsarastack_vswitch.default.vpc_id}"
	specification = "Small"
	name = "${var.name}"
}

resource "apsarastack_eip" "default" {
	name = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
	depends_on = [apsarastack_eip.default, apsarastack_nat_gateway.default]
	allocation_id = "${apsarastack_eip.default.id}"
	instance_id = "${apsarastack_nat_gateway.default.id}"
}

resource "apsarastack_snat_entry" "default"{
	depends_on = [apsarastack_eip_association.default, apsarastack_nat_gateway.default]
	snat_table_id = "${apsarastack_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${apsarastack_vswitch.default.id}"
	snat_ip = "${apsarastack_eip.default.ip_address}"
  
}

resource "apsarastack_snat_entry" "ecs"{
	depends_on = [apsarastack_eip_association.default, apsarastack_nat_gateway.default]
	snat_table_id = "${apsarastack_nat_gateway.default.snat_table_ids}"
	source_cidr = "172.16.0.10/32"
	snat_ip = "${apsarastack_eip.default.ip_address}"
}
`, rand)
}

func testAccSnatEntryConfig_snatIp(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccSnatEntryConfig%d"
}
data "apsarastack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
	depends_on = [apsarastack_vpc.default]
	vpc_id = "${apsarastack_vswitch.default.vpc_id}"
	specification = "Small"
	name = "${var.name}"
}

resource "apsarastack_eip" "default" {
	name = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
	depends_on = [apsarastack_eip.default, apsarastack_nat_gateway.default]
	allocation_id = "${apsarastack_eip.default.id}"
	instance_id = "${apsarastack_nat_gateway.default.id}"
}

resource "apsarastack_snat_entry" "default"{
	depends_on = [apsarastack_eip_association.default, apsarastack_nat_gateway.default]
	snat_table_id = "${apsarastack_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${apsarastack_vswitch.default.id}"
	snat_ip = "${apsarastack_eip.default.ip_address}"

}

resource "apsarastack_snat_entry" "ecs"{
	depends_on = [apsarastack_eip_association.default, apsarastack_nat_gateway.default]
	snat_table_id = "${apsarastack_nat_gateway.default.snat_table_ids}"
	source_cidr = "172.16.0.10/32"
	snat_ip = "${apsarastack_eip.default.ip_address}"
}

`, rand)
}

func testAccSnatEntryConfigMulti(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccSnatEntryMulti%d"
}

data "apsarastack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/16"
}

resource "apsarastack_vswitch" "default" {
    count = 10
    vpc_id            = "${apsarastack_vpc.default.id}"
    cidr_block        = "10.1.${count.index + 1}.0/24"
    availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
    name = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
	depends_on = [apsarastack_vpc.default]
	vpc_id = "${apsarastack_vpc.default.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "apsarastack_eip" "default" {
	name = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
	depends_on = [apsarastack_eip.default, apsarastack_nat_gateway.default]
	allocation_id = "${apsarastack_eip.default.id}"
	instance_id = "${apsarastack_nat_gateway.default.id}"
}

resource "apsarastack_snat_entry" "default"{
	depends_on = [apsarastack_eip_association.default, apsarastack_nat_gateway.default]
	count = "10"
	snat_table_id = "${apsarastack_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${element(apsarastack_vswitch.default.*.id, count.index)}"
	snat_ip = "${apsarastack_eip.default.ip_address}"

}

resource "apsarastack_snat_entry" "ecs"{
	depends_on = [apsarastack_eip_association.default, apsarastack_nat_gateway.default]
	snat_table_id = "${apsarastack_nat_gateway.default.snat_table_ids}"
	source_cidr = "10.1.2.10/32"
	snat_ip = "${apsarastack_eip.default.ip_address}"
}
`, rand)
}

var testAccCheckSnatEntryBasicMap = map[string]string{
	"snat_table_id":     CHECKSET,
	"source_vswitch_id": CHECKSET,
	"snat_ip":           CHECKSET,
	"snat_entry_id":     CHECKSET,
}
