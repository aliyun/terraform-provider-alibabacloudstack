package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackInstancesDataSourceBasic(t *testing.T) {
	resourceId := "data.alibabacloudstack_instances.default"
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID(resourceId),
					resource.TestCheckResourceAttr(resourceId, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "ids.#"),
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
