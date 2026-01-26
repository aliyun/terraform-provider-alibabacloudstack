package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasNamespaceGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_nas_namespace_group.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNasNamespaceGroup)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasNamespaceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-AccNasaccgop%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNasNamespaceGroupBasicDependence)
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
					"network_type":        "Vpc",
					"mapped_path":         "${var.name}",
					"nas_namespace_id":    "${alibabacloudstack_nas_namespace.default.id}",
					"mount_target_domain": "${alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":        "Vpc",
						"mapped_path":         name,
						"member_id":           CHECKSET,
						"status":              CHECKSET,
						"mount_target_domain": CHECKSET,
						"create_time":         CHECKSET,
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

func TestAccAlibabacloudStackNasNamespaceGroup_Classic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_nas_namespace_group.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNasNamespaceGroup)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasNamespaceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-AccNasaccgop%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNasNamespaceGroupClassicDependence)
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
					"network_type":        "Classic",
					"mapped_path":         "${var.name}",
					"nas_namespace_id":    "${alibabacloudstack_nas_namespace.default.id}",
					"mount_target_domain": "${alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":        "Classic",
						"mapped_path":         name,
						"member_id":           CHECKSET,
						"status":              CHECKSET,
						"mount_target_domain": CHECKSET,
						"create_time":         CHECKSET,
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

var AlibabacloudStackNasNamespaceGroup = map[string]string{}

func AlibabacloudStackNasNamespaceGroupBasicDependence(name string) string {
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
	access_group_name = "${var.name}"
	access_group_type = "Vpc"
}

resource "alibabacloudstack_nas_namespace_mount_target" "default" {
  	network_type = "Vpc"
	access_group_name = "${alibabacloudstack_nas_accessgroup.default.access_group_name}"
	nas_namespace_id = "${alibabacloudstack_nas_namespace.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}

`, name, VSwitchCommonTestCase)
}

func AlibabacloudStackNasNamespaceGroupClassicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_nas_zones" "default" {
}

resource "alibabacloudstack_nas_namespace" "default" {
	zone_id = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
	cluster_id = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id}"
	description = var.name
	storage_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type}"
	protocol_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type}"
	encrypt_type = "0"
}

resource "alibabacloudstack_nas_accessgroup" "default" {
	access_group_name = "${var.name}"
	access_group_type = "Classic"
}

resource "alibabacloudstack_nas_namespace_mount_target" "default" {
  	network_type = "Classic"
	access_group_name = "${alibabacloudstack_nas_accessgroup.default.access_group_name}"
	nas_namespace_id = "${alibabacloudstack_nas_namespace.default.id}"
}

`, name)
}
