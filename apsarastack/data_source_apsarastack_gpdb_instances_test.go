package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackGpdbInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackGpdbInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_gpdb_instances.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_gpdb_instances.default", "instances.availability_zone"),
					resource.TestCheckResourceAttrSet("data.apsarastack_gpdb_instances.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackGpdbInstancesDataSource = `
		data "apsarastack_zones" "default" {

       }
       resource "apsarastack_vpc" "default" {
 			name = "testing"
 			cidr_block = "10.0.0.0/8"
		}
		data "apsarastack_gpdb_instances" "default"{
		}
		resource "apsarastack_vswitch" "default" {
 			vpc_id = apsarastack_vpc.default.id
			cidr_block        = "10.1.0.0/16"
 			name = "apsara_vswitch"
 			availability_zone = data.apsarastack_zones.default.zones.0.id
		}
       resource "apsarastack_gpdb_instance" "default" {
           vswitch_id           = apsarastack_vswitch.default.id
           engine               = "gpdb"
           engine_version       = "4.3"
           instance_class       = "gpdb.group.segsdx2"
           instance_group_count = "2"
           description          = "testing_01"
       }
`
