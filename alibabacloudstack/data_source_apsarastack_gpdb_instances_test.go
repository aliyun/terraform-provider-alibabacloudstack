package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackGpdbInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackGpdbInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_gpdb_instances.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_gpdb_instances.default", "instances.availability_zone"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_gpdb_instances.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackGpdbInstancesDataSource = `

		data "alibabacloudstack_zones" "default" {

       }
       resource "alibabacloudstack_vpc" "default" {
 			name = "testing"
 			cidr_block = "10.0.0.0/8"
		}
		data "alibabacloudstack_gpdb_instances" "default"{
		}
		resource "alibabacloudstack_vswitch" "default" {
 			vpc_id = alibabacloudstack_vpc.default.id
			cidr_block        = "10.1.0.0/16"
 			name = "apsara_vswitch"
 			availability_zone = data.alibabacloudstack_zones.default.zones.0.id
		}
       resource "alibabacloudstack_gpdb_instance" "default" {
           vswitch_id           = alibabacloudstack_vswitch.default.id
           engine               = "gpdb"
           engine_version       = "4.3"
           instance_class       = "gpdb.group.segsdx2"
           instance_group_count = "2"
           description          = "testing_01"
       }
`
