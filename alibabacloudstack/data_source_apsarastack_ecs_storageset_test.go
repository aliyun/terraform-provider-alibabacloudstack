package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsEbsStorageSets(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackEcsEbsStorageSet,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ecs_ebs_storage_sets.default", "storages.storage_set_name"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackEcsEbsStorageSet = `
data "alibabacloudstack_zones"  "default" {
}
resource "alibabacloudstack_ecs_ebs_storage_set" "default" {
  storage_set_name = "testcc"
  maxpartition_number = "2"
  zone_id = data.alibabacloudstack_zones.default.zones[0].id
}
data "alibabacloudstack_ecs_ebs_storage_sets" "default"{
 
}

`
