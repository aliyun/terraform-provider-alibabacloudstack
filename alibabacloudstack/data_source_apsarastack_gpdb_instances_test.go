package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackGpdbInstancesDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
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
		data "alibabacloudstack_gpdb_instance_types" "default" {
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
           engine_version       = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.engine_version
           instance_class       = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.id
		   db_instance_storage_type         = "local_ssd"
		   db_instance_mode         = "StorageReserver"
           description          = "testing_01"
		   seg_node_num = "2"
		   network_type = "VPC"
       }
`
