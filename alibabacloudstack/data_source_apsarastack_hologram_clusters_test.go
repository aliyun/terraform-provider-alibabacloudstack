package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHologramClustersDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackHologramClustersDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_hologram_clusters.default"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hologram_clusters.default", "clusters.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hologram_clusters.default", "clusters.0.type"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hologram_clusters.default", "clusters.0.id"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hologram_clusters.default", "clusters.0.cpu"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hologram_clusters.default", "clusters.0.cluster"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackHologramClustersDataSource = `
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}
data "alibabacloudstack_hologram_clusters" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}


`
