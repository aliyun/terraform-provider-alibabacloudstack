package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackNasFileSystem_DataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	storageTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackFileSystemDataSourceConfig(rand, map[string]string{
			"storage_type":      `"${apsarastack_nas_file_system.default.storage_type}"`,
			"description_regex": `"^${apsarastack_nas_file_system.default.description}"`,
		}),
	}
	protocolTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackFileSystemDataSourceConfig(rand, map[string]string{
			"protocol_type":     `"${apsarastack_nas_file_system.default.protocol_type}"`,
			"description_regex": `"^${apsarastack_nas_file_system.default.description}"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${apsarastack_nas_file_system.default.description}"`,
		}),
		fakeConfig: testAccCheckApsaraStackFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${apsarastack_nas_file_system.default.description}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackFileSystemDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_nas_file_system.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackFileSystemDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_nas_file_system.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackFileSystemDataSourceConfig(rand, map[string]string{
			"storage_type":      `"${apsarastack_nas_file_system.default.storage_type}"`,
			"protocol_type":     `"${apsarastack_nas_file_system.default.protocol_type}"`,
			"description_regex": `"^${apsarastack_nas_file_system.default.description}"`,
			"ids":               `["${apsarastack_nas_file_system.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${apsarastack_nas_file_system.default.description}_fake"`,
			"ids":               `["${apsarastack_nas_file_system.default.id}_fake"]`,
		}),
	}

	fileSystemCheckInfo.dataSourceTestCheck(t, rand, storageTypeConf, protocolTypeConf,
		descriptionConf, idsConf, allConf)
}

func testAccCheckApsaraStackFileSystemDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "description" {
  default = "tf-testAccCheckApsaraStackFileSystemsDataSource"
}
variable "storage_type" {
  default = "Capacity"
}
data "apsarastack_nas_protocols" "default" {
        type = "${var.storage_type}"
}
resource "apsarastack_nas_file_system" "default" {
  description = "${var.description}"
  storage_type = "${var.storage_type}"
  protocol_type = "${data.apsarastack_nas_protocols.default.protocols.0}"
}
data "apsarastack_nas_file_systems" "default" {
	%s
}`, strings.Join(pairs, "\n  "))
	return config
}

var existFileSystemMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"systems.0.id":            CHECKSET,
		"systems.0.region_id":     CHECKSET,
		"systems.0.description":   "tf-testAccCheckApsaraStackFileSystemsDataSource",
		"systems.0.protocol_type": CHECKSET,
		"systems.0.storage_type":  "Capacity",
		"systems.0.metered_size":  CHECKSET,
		"systems.0.create_time":   CHECKSET,
		"ids.#":                   "1",
		"ids.0":                   CHECKSET,
		"descriptions.#":          "1",
		"descriptions.0":          "tf-testAccCheckApsaraStackFileSystemsDataSource",
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
	resourceId:   "data.apsarastack_nas_file_systems.default",
	existMapFunc: existFileSystemMapCheck,
	fakeMapFunc:  fakeFileSystemMapCheck,
}
