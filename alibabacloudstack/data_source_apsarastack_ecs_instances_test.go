package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackInstancesDataSourceBasic(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_instances.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instances.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackInstancesDataSource = ECSInstanceCommonTestCase + `
variable "name" {
  default = "Tf-EcsInstanceDataSource"
}

data "alibabacloudstack_instances" "default" {
  ids = ["${alibabacloudstack_ecs_instance.default.id}"]
}
`
