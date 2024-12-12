package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackSlbBackendServersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSlbBackendServersDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_slb_backend_servers.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_slb_backend_servers.default", "load_balancer_id.#", "0"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccCheckAlibabacloudStackSlbBackendServersDataSource = ECSInstanceCommonTestCase + `
variable "name" {
	default = "tf-slbBackendServersdatasourcebasic"
}

resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}

resource "alibabacloudstack_slb_backend_server" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"

  backend_servers {
    server_id = "${alibabacloudstack_ecs_instance.default.id}"
    weight     = 100
  }
}

data "alibabacloudstack_slb_backend_servers" "default" {
 load_balancer_id = "${alibabacloudstack_slb.default.id}"
}
`
