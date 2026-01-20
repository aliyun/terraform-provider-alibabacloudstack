package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackGpdbConnectionUpdate(t *testing.T) {
	var v *gpdb.DBInstanceNetInfo

	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tftest%d", rand)
	var basicMap = map[string]string{
		"instance_id": CHECKSET,
		"port":        "3306",
	}

	resourceId := "alibabacloudstack_gpdb_connection.default"
	serverFunc := func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeGpdbConnection")
	ra := resourceAttrInit(resourceId, basicMap)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testGpdbConnectionConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${local.gpdb_instance_id}",
					"connection_prefix": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port":              "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix": name + "-update",
					}),
				),
			},
		},
	})
}

func testGpdbConnectionConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
			default = "%s"
		}
	%s
		`, name, GpdbCommonTestCase())
}
