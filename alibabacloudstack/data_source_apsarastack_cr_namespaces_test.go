package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
	"fmt"
)

func TestAccAlibabacloudStackCRNamespacesDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceCRNamespacesConfigDependence(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_cr_namespaces.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cr_namespaces.default", "namespaces.default_visibility"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cr_namespaces.default", "instance_types.auto_create"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cr_namespaces.default", "instance_types.name"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_cr_namespaces.default", "ids.#"),
				),
			},
		},
	})
}

func dataSourceCRNamespacesConfigDependence() string {
	return fmt.Sprintf(`
  resource "alibabacloudstack_cr_namespace" "default" {
  name               = "testacc-cr-namespace%d"
  auto_create        = false
  default_visibility = "PUBLIC"
}

  data "alibabacloudstack_cr_namespaces" "default" {
  name_regex    = alibabacloudstack_cr_namespace.default.name
}
`, getAccTestRandInt(1000000, 9999999))
}
