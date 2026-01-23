package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNasNamespaces_basic(t *testing.T) {
	resourceId := "data.alibabacloudstack_nas_namespaces.default"
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf-testnasfs%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackNasNamespaceDataSourceConfig)
	datasourceData := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"namespaces.#":                    "1",
				"namespaces.0.description":        name,
				"namespaces.0.protocol_type":      "NFS",
				"namespaces.0.storage_type":       CHECKSET,
				"namespaces.0.file_system_type":   CHECKSET,
				"namespaces.0.zone_id":            CHECKSET,
				"namespaces.0.status":             CHECKSET,
				"namespaces.0.create_time":        CHECKSET,
				"namespaces.0.encrypt_type":       CHECKSET,
				"namespaces.0.mount_target_count": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"namespaces.#": "0",
			}
		},
	}

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nas_namespace.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nas_namespace.default.id}_fake"},
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_nas_namespace.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_nas_namespace.default.description}_fake",
		}),
	}

	zoneIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alibabacloudstack_nas_namespace.default.id}"},
			"zone_id": "${alibabacloudstack_nas_namespace.default.zone_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alibabacloudstack_nas_namespace.default.id}"},
			"zone_id": "${alibabacloudstack_nas_namespace.default.zone_id}_fake",
		}),
	}

	storageTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alibabacloudstack_nas_namespace.default.id}"},
			"storage_type": "${alibabacloudstack_nas_namespace.default.storage_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alibabacloudstack_nas_namespace.default.id}"},
			"storage_type": "${alibabacloudstack_nas_namespace.default.storage_type}_fake",
		}),
	}

	protocolTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_nas_namespace.default.id}"},
			"protocol_type": "${alibabacloudstack_nas_namespace.default.protocol_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_nas_namespace.default.id}"},
			"protocol_type": "SMB",
		}),
	}

	datasourceData.dataSourceTestCheck(t, rand, idsConfig, nameRegexConfig, zoneIdConfig, storageTypeConfig, protocolTypeConfig)
}

func testAccCheckAlibabacloudStackNasNamespaceDataSourceConfig(name string) string {
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
	protocol_type = "NFS"
	encrypt_type = "0"
}

`, name)
}
