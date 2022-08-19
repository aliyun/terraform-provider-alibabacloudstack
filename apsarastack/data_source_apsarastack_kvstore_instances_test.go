package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackKVStoreInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceKVStoreInstancesConfigDependence,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_kvstore_instances.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.charge_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.region_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.create_time"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.vpc_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.instance_class"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.status"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.availability_zone"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.bandwidth"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.user_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.connections"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.vswitch_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_kvstore_instances.default", "instances.connection_domain"),
					resource.TestCheckResourceAttrSet("data.apsarastack_kvstore_instances.default", "ids.#"),
				),
			},
		},
	})
}

const dataSourceKVStoreInstancesConfigDependence = `
variable "name" {
    default = "tf-testAccCheckApsaraStackRKVInstancesDataSource0"
}
data "apsarastack_zones"  "default" {
}
resource "apsarastack_vpc" "default" {
name       = var.name
cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
vpc_id            = apsarastack_vpc.default.id
cidr_block        = "172.16.0.0/24"
availability_zone = data.apsarastack_zones.default.zones[0].id
name              = var.name
}
resource "apsarastack_kvstore_instance" "default" {
instance_class = "redis.master.small.default"
instance_name  = var.name
vswitch_id     = apsarastack_vswitch.default.id
private_ip     = "172.16.0.10"
security_ips   = ["10.0.0.1"]
instance_type  = "Redis"
engine_version = "4.0"
}
data "apsarastack_kvstore_instances" "default" {
  name_regex = apsarastack_kvstore_instance.default.instance_name
}
`
