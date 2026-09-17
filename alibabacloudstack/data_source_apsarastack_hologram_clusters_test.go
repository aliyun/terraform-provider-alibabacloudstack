package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHologramClustersDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_hologram_clusters.default"
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackHologramClustersDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID(resourceId),
					resource.TestCheckResourceAttrSet(resourceId, "clusters.#"),
					resource.TestCheckResourceAttrSet(resourceId, "clusters.0.type"),
					resource.TestCheckResourceAttrSet(resourceId, "clusters.0.id"),
					resource.TestCheckResourceAttrSet(resourceId, "clusters.0.cpu"),
					resource.TestCheckResourceAttrSet(resourceId, "clusters.0.cluster"),
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
  name_regex = "Hologres.*"
}


`
