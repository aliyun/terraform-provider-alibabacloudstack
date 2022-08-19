package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackSnatEntriesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSnatEntriesBasicConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_snat_entries.default"),
					//resource.TestCheckResourceAttr("data.apsarastack_snat_entries.default", "snat_table_id", "0"),
					resource.TestCheckResourceAttrSet("data.apsarastack_snat_entries.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackSnatEntriesBasicConfig = `
variable "name" {
	default = "tf-testAccForSnatEntriesDatasource"
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
	vpc_id = "${apsarastack_vpc.default.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "apsarastack_eip" "default" {
	name = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
	allocation_id = "${apsarastack_eip.default.id}"
	instance_id = "${apsarastack_nat_gateway.default.id}"
}

resource "apsarastack_snat_entry" "default" {
	snat_table_id = "${apsarastack_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${apsarastack_vswitch.default.id}"
	snat_ip = "${apsarastack_eip.default.ip_address}"
}

data "apsarastack_snat_entries" "default" {
    snat_table_id = "${apsarastack_nat_gateway.default.snat_table_ids}"
}`
