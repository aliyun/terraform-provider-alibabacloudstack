package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasNamespaceMountTarget_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_nas_namespace_mount_target.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNasNamespaceMountTarget)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasNamespaceMountTarget")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-AccNasaccgop%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNasNamespaceMountTargetBasicDependence)
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
					"network_type":      "Vpc",
					"access_group_name": "${alibabacloudstack_nas_accessgroup.default.0.access_group_name}",
					"nas_namespace_id":  "${alibabacloudstack_nas_namespace.default.id}",
					"vswitch_id":        "${alibabacloudstack_vpc_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":        "Vpc",
						"access_group_name":   name + "_0",
						"vpc_id":              CHECKSET,
						"mount_target_domain": CHECKSET,
						"status":              CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alibabacloudstack_nas_accessgroup.default.1.access_group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name + "_1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Active",
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

var AlibabacloudStackNasNamespaceMountTarget = map[string]string{}

func AlibabacloudStackNasNamespaceMountTargetBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_nas_zones" "default" {
}

%s

resource "alibabacloudstack_nas_namespace" "default" {
	zone_id = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
	cluster_id = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id}"
	description = var.name
	storage_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type}"
	protocol_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type}"
	encrypt_type = "0"
}

resource "alibabacloudstack_nas_accessgroup" "default" {
	count             = 2
	access_group_name = "${var.name}_${count.index}"
	access_group_type = "Vpc"
}

`, name, VSwitchCommonTestCase)
}
