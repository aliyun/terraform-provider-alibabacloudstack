package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsDedicatedHost_basic0(t *testing.T) {
	time.Sleep(3 * time.Minute)
	var v ecs.DedicatedHost

	resourceId := "alibabacloudstack_ecs_dedicatedhost.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDedicatedhostCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribededicatedhostautorenewRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccddhhost%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDedicatedhostBasicdependence)
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
					"dedicated_host_name":       "${var.name}",
					"dedicated_host_type":       "${local.ddh_type}",
					"dedicated_host_cluster_id": "${alibabacloudstack_ecs_dedicated_host_cluster.default.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"zone_id": "${data.alibabacloudstack_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_name":       name,
						"dedicated_host_cluster_id": CHECKSET,
						"action_on_maintenance":     "Stop",
						"auto_placement":            "on",
						"tags.%":                    "2",
						"tags.Created":              "TF",
						"tags.For":                  "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dedicated_host_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_on_maintenance": "Migrate",
					"auto_placement":        "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_on_maintenance": "Migrate",
						"auto_placement":        "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudTestAccEcsDedicatedhostCheckmap = map[string]string{
	"machine_id":                         CHECKSET,
	"dedicated_host_name":                CHECKSET,
	"status":                             CHECKSET,
	"zone_id":                            CHECKSET,
	"auto_placement":                     CHECKSET,
	"dedicated_host_type":                CHECKSET,
	"supported_instance_type_families.#": CHECKSET,
	"supported_instance_types_list.#":    CHECKSET,
}

func AlibabacloudTestAccEcsDedicatedhostBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ddh_type" {
	default = "%s"
}

%s

data "alibabacloudstack_ecs_dedicated_host_types" "default" {
}

locals {
	ddh_type = var.ddh_type == "" ? data.alibabacloudstack_ecs_dedicated_host_types.default.ids.0 : var.ddh_type
}
	

resource "alibabacloudstack_ecs_dedicated_host_cluster" "default" {
  dedicated_host_cluster_name = "${var.name}_cluster"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_DDH_TYPE"), DataZoneCommonTestCase)
}
