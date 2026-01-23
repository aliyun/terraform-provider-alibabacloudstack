package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNasNamespaceFilesystemAttachments_basic(t *testing.T) {
	resourceId := "data.alibabacloudstack_nas_namespace_filesystem_attachments.default"
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf-testnasfs%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackNasNamespaceFilesystemAttachment)
	datasourceData := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"attachments.#":                  "1",
				"attachments.0.mapped_path":      name,
				"attachments.0.attachment_id":    CHECKSET,
				"attachments.0.nas_namespace_id": CHECKSET,
				"attachments.0.file_system_id":   CHECKSET,
				"attachments.0.storage_type":     CHECKSET,
				"attachments.0.file_system_type": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"attachments.#": "0",
			}
		},
	}

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_filesystem_attachment.default.nas_namespace_id}",
			"ids":              []string{"${alibabacloudstack_nas_namespace_filesystem_attachment.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_filesystem_attachment.default.nas_namespace_id}",
			"ids":              []string{"${alibabacloudstack_nas_namespace_filesystem_attachment.default.id}_fake"},
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_filesystem_attachment.default.nas_namespace_id}",
			"name_regex":       "${alibabacloudstack_nas_namespace_filesystem_attachment.default.mapped_path}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_filesystem_attachment.default.nas_namespace_id}",
			"name_regex":       "${alibabacloudstack_nas_namespace_filesystem_attachment.default.mapped_path}_fake",
		}),
	}

	mappedPath := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_filesystem_attachment.default.nas_namespace_id}",
			"mapped_path":      "${alibabacloudstack_nas_namespace_filesystem_attachment.default.mapped_path}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_filesystem_attachment.default.nas_namespace_id}",
			"mapped_path":      "${alibabacloudstack_nas_namespace_filesystem_attachment.default.mapped_path}_fake",
		}),
	}

	fileSystemConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_filesystem_attachment.default.nas_namespace_id}",
			"file_system_id":   "${alibabacloudstack_nas_namespace_filesystem_attachment.default.file_system_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"nas_namespace_id": "${alibabacloudstack_nas_namespace_filesystem_attachment.default.nas_namespace_id}",
			"file_system_id":   "1000000000",
		}),
	}

	datasourceData.dataSourceTestCheck(t, rand, idsConfig, nameRegexConfig, mappedPath, fileSystemConfig)
}

func testAccCheckAlibabacloudStackNasNamespaceFilesystemAttachment(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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

resource "alibabacloudstack_nas_namespace_filesystem_attachment" "default" {
	nas_namespace_id = "${alibabacloudstack_nas_namespace.default.id}"
	file_system_id = "${alibabacloudstack_nas_file_system.default.id}"
	mapped_path = var.name
}

`, name, NasCommonTestCase)
}
