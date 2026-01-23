package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNasNamespaceMountTargets_basic(t *testing.T) {
	resourceId := "data.alibabacloudstack_nas_namespace_mount_targets.default"
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf-testnasfs%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackNasNamespaceMountTarget)
	datasourceData := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"mount_targets.#":                     "1",
				"mount_targets.0.mount_target_domain": CHECKSET,
				"mount_targets.0.network_type":        CHECKSET,
				"mount_targets.0.vpc_id":              CHECKSET,
				"mount_targets.0.vsw_id":              CHECKSET,
				"mount_targets.0.access_group":        CHECKSET,
				"mount_targets.0.status":              CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"mount_targets.#": "0",
			}
		},
	}

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_mount_target.default.nas_namespace_id}",
			"ids":              []string{"${alibabacloudstack_nas_namespace_mount_target.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_mount_target.default.nas_namespace_id}",
			"ids":              []string{"${alibabacloudstack_nas_namespace_mount_target.default.id}_fake"},
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_mount_target.default.nas_namespace_id}",
			"name_regex":       "${alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_mount_target.default.nas_namespace_id}",
			"name_regex":       "${alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain}_fake",
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_mount_target.default.nas_namespace_id}",
			"ids":              []string{"${alibabacloudstack_nas_namespace_mount_target.default.id}"},
			"name_regex":       "${alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_mount_target.default.nas_namespace_id}",
			"name_regex":       "${alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain}_fake",
			"ids":              []string{"${alibabacloudstack_nas_namespace_mount_target.default.id}_fake"},
		}),
	}

	datasourceData.dataSourceTestCheck(t, rand, idsConfig, nameRegexConfig, allConfig)
}

func testAccCheckAlibabacloudStackNasNamespaceMountTarget(name string) string {
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
	nas_namespace_id  = "${alibabacloudstack_nas_namespace.default.id}"
	access_group_name = "${alibabacloudstack_nas_accessgroup.default.access_group_name}"
	vswitch_id 		  = "${alibabacloudstack_vpc_vswitch.default.id}"
	network_type 	  = "Vpc"
}

`, name, VSwitchCommonTestCase)
}
