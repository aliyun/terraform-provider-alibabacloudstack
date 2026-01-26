package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNamespaceGroups_basic(t *testing.T) {
	resourceId := "data.alibabacloudstack_nas_namespace_groups.default"
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf-testnasng%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackNasNamespaceMountTarget)
	datasourceData := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"groups.#":                     "1",
				"groups.0.mount_target_domain": CHECKSET,
				"groups.0.network_type":        CHECKSET,
				"groups.0.vpc_id":              CHECKSET,
				"groups.0.vsw_id":              CHECKSET,
				"groups.0.access_group":        CHECKSET,
				"groups.0.status":              CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"groups.#": "0",
			}
		},
	}

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nas_namespace_groups.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nas_namespace_groups.default.id}_fake"},
		}),
	}

	nasNamespaceConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_groups.default.nas_namespace_id}",
			"ids":              []string{"${alibabacloudstack_nas_namespace_groups.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_groups.default.nas_namespace_id}",
			"ids":              []string{"${alibabacloudstack_nas_namespace_groups.default.id}_fake"},
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_nas_namespace_groups.default.id}"},
			"name_regex": "${alibabacloudstack_nas_namespace_groups.default.mount_target_domain}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_nas_namespace_groups.default.id}"},
			"name_regex": "${alibabacloudstack_nas_namespace_groups.default.mount_target_domain}_fake",
		}),
	}

	mountTargetDomainConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"mount_target_domain": "${alibabacloudstack_nas_namespace_groups.default.mount_target_domain}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"mount_target_domain": "${alibabacloudstack_nas_namespace_groups.default.mount_target_domain}_fake",
		}),
	}

	mappedPathConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"mapped_path": "${alibabacloudstack_nas_namespace_groups.default.mapped_path}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"mapped_path": "${alibabacloudstack_nas_namespace_groups.default.mapped_path}_fake",
		}),
	}

	networkType := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"network_type": "Vpc",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"mapped_path": "Classic",
		}),
	}

	statusType := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "Enabled",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "All",
		}),
	}

	datasourceData.dataSourceTestCheck(t, rand, idsConfig, nasNamespaceConfig, nameRegexConfig, mountTargetDomainConfig, mappedPathConfig, networkType, statusType)
}

func AlibabacloudStackNasNamespaceGroupsBasicDependence(name string) string {
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

resource "alibabacloudstack_nas_namespace_group" "default" {
  	network_type = "Vpc"
	nas_namespace_id = "${alibabacloudstack_nas_namespace.default.id}"
	mount_target_domain = "${alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain}"
	mapped_path = var.name
}

`, name, VSwitchCommonTestCase)
}
