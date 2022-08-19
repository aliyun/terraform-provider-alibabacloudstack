package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApsaraStackForwardEntryBasic(t *testing.T) {
	var v vpc.ForwardTableEntry
	resourceId := "apsarastack_forward_entry.default"

	rand := acctest.RandInt()
	testAccForwardEntryCheckMap["name"] = fmt.Sprintf("tf-testAccForwardEntryConfig%d", rand)
	ra := resourceAttrInit(resourceId, testAccForwardEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckForwardEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccForwardEntryConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccApsaraStackForwardEntryMulti(t *testing.T) {
	var v vpc.ForwardTableEntry
	resourceId := "apsarastack_forward_entry.default.4"
	rand := acctest.RandInt()
	testAccForwardEntryCheckMap["name"] = fmt.Sprintf("tf-testAccForwardEntryConfig%d", rand)
	ra := resourceAttrInit(resourceId, testAccForwardEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckForwardEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccForwardEntryConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port": "84",
						"internal_port": "8084",
					}),
				),
			},
		},
	})
}

func testAccCheckForwardEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_forward_entry" {
			continue
		}
		if _, err := vpcService.DescribeForwardEntry(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(fmt.Errorf("Forward entry %s still exist", rs.Primary.ID))
	}
	return nil
}

func testAccForwardEntryConfigBasic(rand int) string {
	config := fmt.Sprintf(`
%s

resource "apsarastack_forward_entry" "default"{
	name = "${var.name}"
	forward_table_id = "${apsarastack_nat_gateway.default.forward_table_ids}"
	external_ip = "${apsarastack_eip.default.0.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.4"
	internal_port = "8080"


}
`, testAccForwardEntryConfigCommon(rand))
	return config
}

func testAccForwardEntryConfig_multi(rand int) string {
	config := fmt.Sprintf(`
%s

resource "apsarastack_forward_entry" "default"{
	count = 5
	name = "${var.name}"
	forward_table_id = "${apsarastack_nat_gateway.default.forward_table_ids}"
	external_ip = "${apsarastack_eip.default.0.ip_address}"
	external_port = "${80 + count.index}"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "${8080 + count.index}"
}
`, testAccForwardEntryConfigCommon(rand))
	return config
}

func testAccForwardEntryConfigCommon(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccForwardEntryConfig%d"
}

variable "number" {
	default = "2"
}

data "apsarastack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "apsarastack_vpc" "default" {
	
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	
}

resource "apsarastack_nat_gateway" "default" {
	vpc_id = "${apsarastack_vswitch.default.vpc_id}"
	specification = "Small"
	
}

resource "apsarastack_eip" "default" {
	count = "${var.number}"
	
}

resource "apsarastack_eip_association" "default" {
	count = "${var.number}"
	allocation_id = "${element(apsarastack_eip.default.*.id,count.index)}"
	instance_id = "${apsarastack_nat_gateway.default.id}"
}
`, rand)
}

var testAccForwardEntryCheckMap = map[string]string{
	"forward_table_id": CHECKSET,
	"external_ip":      CHECKSET,
	"external_port":    "80",
	"ip_protocol":      "tcp",
	"internal_port":    "8080",
	"forward_entry_id": CHECKSET,
}
