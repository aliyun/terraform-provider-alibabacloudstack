package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackNasFileSystem_DataSource(t *testing.T) {
	rand := getAccTestRandInt(100000, 999999)
	storageTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackFileSystemDataSourceConfig(rand, map[string]string{
			"storage_type":      `"${alibabacloudstack_nas_file_system.default.storage_type}"`,
			"description_regex": `"^${alibabacloudstack_nas_file_system.default.description}"`,
		}),
	}
	protocolTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackFileSystemDataSourceConfig(rand, map[string]string{
			"protocol_type":     `"${alibabacloudstack_nas_file_system.default.protocol_type}"`,
			"description_regex": `"^${alibabacloudstack_nas_file_system.default.description}"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alibabacloudstack_nas_file_system.default.description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alibabacloudstack_nas_file_system.default.description}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackFileSystemDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_file_system.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackFileSystemDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_file_system.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackFileSystemDataSourceConfig(rand, map[string]string{
			"storage_type":      `"${alibabacloudstack_nas_file_system.default.storage_type}"`,
			"protocol_type":     `"${alibabacloudstack_nas_file_system.default.protocol_type}"`,
			"description_regex": `"^${alibabacloudstack_nas_file_system.default.description}"`,
			"ids":               `["${alibabacloudstack_nas_file_system.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alibabacloudstack_nas_file_system.default.description}_fake"`,
			"ids":               `["${alibabacloudstack_nas_file_system.default.id}_fake"]`,
		}),
	}

	fileSystemCheckInfo.dataSourceTestCheck(t, rand, storageTypeConf, protocolTypeConf,
		descriptionConf, idsConf, allConf)
}

func testAccCheckAlibabacloudStackFileSystemDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "description" {
  default = "tf-testAccCheckAlibabacloudStackFileSystemsDataSource"
}
variable "storage_type" {
  default = "Capacity"
}
data "alibabacloudstack_nas_protocols" "default" {
        type = "${var.storage_type}"
}
resource "alibabacloudstack_nas_file_system" "default" {
  description = "${var.description}"
  storage_type = "${var.storage_type}"
  protocol_type = "${data.alibabacloudstack_nas_protocols.default.protocols.0}"
}
data "alibabacloudstack_nas_file_systems" "default" {
	%s
}`, strings.Join(pairs, "\n  "))
	return config
}

var existFileSystemMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"systems.0.id":            CHECKSET,
		"systems.0.region_id":     CHECKSET,
		"systems.0.description":   "tf-testAccCheckAlibabacloudStackFileSystemsDataSource",
		"systems.0.protocol_type": CHECKSET,
		"systems.0.storage_type":  "Capacity",
		"systems.0.metered_size":  CHECKSET,
		"systems.0.create_time":   CHECKSET,
		"ids.#":                   "1",
		"ids.0":                   CHECKSET,
		"descriptions.#":          "1",
		"descriptions.0":          "tf-testAccCheckAlibabacloudStackFileSystemsDataSource",
	}
}

var fakeFileSystemMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"systems.#":      "0",
		"ids.#":          "0",
		"descriptions.#": "0",
	}
}

var fileSystemCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_nas_file_systems.default",
	existMapFunc: existFileSystemMapCheck,
	fakeMapFunc:  fakeFileSystemMapCheck,
}
