package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAdbConnection0(t *testing.T) {
	var v *adb.Address

	resourceId := "alibabacloudstack_adb_connection.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccAdbConnectionCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoAdbDescribedbclusternetinfoRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-adbconnection%d", rand)

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
					"db_cluster_id": "${local.adb_instance_id}",
					"connection_prefix": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": name,
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
					"connection_prefix": "${var.name}-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": name+"-update",
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
	"connection_prefix": CHECKSET,
}

func AlibabacloudTestAccAdbConnectionBasicdependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
			default = "%s"
		}

		%s
		
	`, name, AdbCommonTestCase(false))
}
