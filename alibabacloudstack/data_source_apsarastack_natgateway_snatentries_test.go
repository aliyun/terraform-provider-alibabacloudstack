package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackSnatEntriesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSnatEntriesBasicConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_snat_entries.default"),
					//resource.TestCheckResourceAttr("data.alibabacloudstack_snat_entries.default", "snat_table_id", "0"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_snat_entries.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackSnatEntriesBasicConfig = `
variable "name" {
	default = "tf-testAccForSnatEntriesDatasource"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alibabacloudstack_eip" "default" {
	name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
	allocation_id = "${alibabacloudstack_eip.default.id}"
	instance_id = "${alibabacloudstack_nat_gateway.default.id}"
}

resource "alibabacloudstack_snat_entry" "default" {
	snat_table_id = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	snat_ip = "${alibabacloudstack_eip.default.ip_address}"
}

data "alibabacloudstack_snat_entries" "default" {
    snat_table_id = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
}`
