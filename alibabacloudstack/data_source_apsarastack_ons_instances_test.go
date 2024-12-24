package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
	"fmt"
)

func TestAccAlibabacloudStackOnsInstancesDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceOnsInstancesConfigDependence(),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ons_instances.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ons_instances.default", "instances.instance_name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ons_instances.default", "instances.topic_capacity"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ons_instances.default", "instances.tps_receive_max"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ons_instances.default", "instances.tps_send_max"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ons_instances.default", "instances.cluster"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ons_instances.default", "instances.instance_status"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_ons_instances.default", "ids.#"),
				),
			},
		},
	})
}

func dataSourceOnsInstancesConfigDependence() string {
	return fmt.Sprintf(`variable "name" {
  default = "Tf-OnsInstanceDataSource%d"
}

resource "alibabacloudstack_ons_instance" "default" {
  name = "${var.name}"
  remark = "default-remark"
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
}
data "alibabacloudstack_ons_instances" "default" {
  ids = [alibabacloudstack_ons_instance.default.id]

}
`, getAccTestRandInt(10000,20000) )
}
