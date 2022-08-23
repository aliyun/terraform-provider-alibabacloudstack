package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"

	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackGpdbConnectionUpdate(t *testing.T) {
	var v *gpdb.DBInstanceNetInfo

	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	var basicMap = map[string]string{
		"instance_id": CHECKSET,
		"port":        "3306",
	}

	resourceId := "apsarastack_gpdb_connection.default"
	serverFunc := func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeGpdbConnection")
	ra := resourceAttrInit(resourceId, basicMap)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", testGpdbConnectionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${apsarastack_gpdb_instance.default.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%s", rand),
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
					"instance_id":       "${apsarastack_gpdb_instance.default.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%s", rand),
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
        data "apsarastack_zones" "default" {
            available_resource_creation = "Gpdb"
        }
        variable "name" {
            default = "tf-testAccGpdbInstance"
        }
		resource "apsarastack_vpc" "default" {
  			name = "testing"
  			cidr_block = "10.0.0.0/8"
		}
		resource "apsarastack_vswitch" "default" {
  			vpc_id = apsarastack_vpc.default.id
			cidr_block        = "10.1.0.0/16"
  			name = "apsara_vswitch"
  			availability_zone = data.apsarastack_zones.default.zones.0.id
		}
        resource "apsarastack_gpdb_instance" "default" {
            vswitch_id           = apsarastack_vswitch.default.id
            engine               = "gpdb"
            engine_version       = "4.3"
            instance_class       = "gpdb.group.segsdx2"
            instance_group_count = "2"
            description          = "${var.name}"
	}`)
}
