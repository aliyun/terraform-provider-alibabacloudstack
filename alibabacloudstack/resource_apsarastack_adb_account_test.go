package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAdbAccount0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_adb_account.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccAdbAccountCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoAdbDescribeaccountsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadbaccount%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccAdbAccountBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"db_cluster_id": "am-bp1j43v9c35ef2cvf",

					"account_type": "Normal",

					"account_name": "nametest123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"db_cluster_id": "am-bp1j43v9c35ef2cvf",

						"account_type": "Normal",

						"account_name": "nametest123",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"db_cluster_id": "am-bp1j43v9c35ef2cvf",

					"account_type": "Normal",

					"account_name": "nametest123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"db_cluster_id": "am-bp1j43v9c35ef2cvf",

						"account_type": "Normal",

						"account_name": "nametest123",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

var AlibabacloudTestAccAdbAccountCheckmap = map[string]string{

	"account_description": CHECKSET,

	"status": CHECKSET,

	"db_cluster_id": CHECKSET,

	"account_type": CHECKSET,

	"account_name": CHECKSET,
}

func AlibabacloudTestAccAdbAccountBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
