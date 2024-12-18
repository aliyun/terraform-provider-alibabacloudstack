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

	rand := getAccTestRandInt(10000,20000)
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", testGpdbConnectionConfigDependence)

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
					"instance_id":       "${alibabacloudstack_gpdb_instance.default.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alibabacloudstack_gpdb_instance.default.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%d", rand),
					"port":              "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
		},
	})
}

func testGpdbConnectionConfigDependence(name string) string {
	return fmt.Sprintf(`
        data "alibabacloudstack_zones" "default" {
            available_resource_creation = "Gpdb"
        }
        variable "name" {
            default = "tf-testAccGpdbInstance"
        }
		resource "alibabacloudstack_vpc" "default" {
  			name = "testing"
  			cidr_block = "10.0.0.0/8"
		}
		resource "alibabacloudstack_vswitch" "default" {
  			vpc_id = alibabacloudstack_vpc.default.id
			cidr_block        = "10.1.0.0/16"
  			name = "apsara_vswitch"
  			availability_zone = data.alibabacloudstack_zones.default.zones.0.id
		}
        resource "alibabacloudstack_gpdb_instance" "default" {
            vswitch_id           = alibabacloudstack_vswitch.default.id
            engine               = "gpdb"
            engine_version       = "4.3"
            instance_class       = "gpdb.group.segsdx2"
            instance_group_count = "2"
            description          = "${var.name}"
	}`)
}
