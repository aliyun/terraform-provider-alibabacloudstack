package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOssEndpointsDataSource_basic(t *testing.T) {

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceOssEndpointsConfigDependence_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_oss_endpoints.default"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_oss_endpoints.default", "endpoints.0.id"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_oss_endpoints.default", "endpoints.0.cluster"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_oss_endpoints.default", "endpoints.0.ha_apsara_stack"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_oss_endpoints.default", "endpoints.0.real_zone"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_oss_endpoints.default", "endpoints.0.oss_ha_enable_single_cluster_access"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_oss_endpoints.default", "endpoints.0.cluster_name"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_oss_endpoints.default", "endpoints.0.is_master_zone"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_oss_endpoints.default", "endpoints.0.location"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_oss_endpoints.default", "endpoints.0.oss_endpoint"),
				),
			},
		},
	})
}

func dataSourceOssEndpointsConfigDependence_basic() string {
	return fmt.Sprintf(`
data "alibabacloudstack_oss_endpoints" "default" {
}
`)
}
