package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAdbConnection0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_adb_connection.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccAdbConnectionCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoAdbDescribedbclusternetinfoRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadbconnection%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccAdbConnectionBasicdependence)
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

					"connection_string_prefix": "am-bp1j43v9c35ef2cvf80808",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"db_cluster_id": "am-bp1j43v9c35ef2cvf",

						"connection_string_prefix": "am-bp1j43v9c35ef2cvf80808",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"connection_string_prefix": "am-bp1j43v9c35ef2cvf80808",

					"connection_string": "am-bp1j43v9c35ef2cvf907780.ads.aliyuncs.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"connection_string_prefix": "am-bp1j43v9c35ef2cvf80808",

						"connection_string": "am-bp1j43v9c35ef2cvf907780.ads.aliyuncs.com",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccAdbConnectionCheckmap = map[string]string{

	"port": CHECKSET,

	"db_cluster_id": CHECKSET,

	"connection_string": CHECKSET,

	"ip_address": CHECKSET,

	"connection_string_prefix": CHECKSET,
}

func AlibabacloudTestAccAdbConnectionBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
