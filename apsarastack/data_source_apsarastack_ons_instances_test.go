package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackOnsInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceOnsInstancesConfigDependence,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_ons_instances.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ons_instances.default", "instances.instance_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ons_instances.default", "instances.topic_capacity"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ons_instances.default", "instances.tps_receive_max"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ons_instances.default", "instances.tps_send_max"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ons_instances.default", "instances.cluster"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ons_instances.default", "instances.instance_status"),
					resource.TestCheckResourceAttrSet("data.apsarastack_ons_instances.default", "ids.#"),
				),
			},
		},
	})
}

const dataSourceOnsInstancesConfigDependence = `
variable "name" {
  default = "Tf-OnsInstanceDataSource"
}

resource "apsarastack_ons_instance" "default" {
  name = "${var.name}"
  remark = "default-remark"
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
}
data "apsarastack_ons_instances" "default" {
  ids = [apsarastack_ons_instance.default.id]

}
`
