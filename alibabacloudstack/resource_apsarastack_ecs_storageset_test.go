package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsStorageset0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_storageset.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsStoragesetCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribestoragesetdetailsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsstorage_set%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsStoragesetBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "wTest",

					"zone_id": "cn-hangzhou-j",

					"region_id": "cn-hangzhou",

					"storage_set_name": "w测试",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "wTest",

						"zone_id": "cn-hangzhou-j",

						"region_id": "cn-hangzhou",

						"storage_set_name": "w测试",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "update",

					"region_id": "cn-hangzhou",

					"storage_set_name": "存储集修改",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "update",

						"region_id": "cn-hangzhou",

						"storage_set_name": "存储集修改",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccEcsStoragesetCheckmap = map[string]string{

	"description": CHECKSET,

	"zone_id": CHECKSET,

	"max_partition_number": CHECKSET,

	"region_id": CHECKSET,

	"storage_set_id": CHECKSET,

	"storage_set_name": CHECKSET,
}

func AlibabacloudTestAccEcsStoragesetBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
