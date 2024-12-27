package alibabacloudstack

import (
	"testing"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsEbsStorageSets_datasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackEcsEbsStorageSet(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ecs_ebs_storage_sets.default", "storages.storage_set_name"),
				),
			},
		},
	})
}

func dataSourceAlibabacloudStackEcsEbsStorageSet() string {
return fmt.Sprintf(`
data "alibabacloudstack_zones"  "default" {
}
resource "alibabacloudstack_ecs_ebs_storage_set" "default" {
  storage_set_name = "tf-testAcc_storage_set%d"
  maxpartition_number = "2"
  zone_id = data.alibabacloudstack_zones.default.zones[0].id
}
data "alibabacloudstack_ecs_ebs_storage_sets" "default"{
 
}

`, getAccTestRandInt(1000, 9999))
}
