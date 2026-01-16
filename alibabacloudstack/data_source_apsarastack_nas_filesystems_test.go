package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNasFileSystem_DataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_nas_file_systems.default"
	name := fmt.Sprintf("tf-testnasfs%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackFileSystemDataSourceConfig)

	storageTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"storage_type":      "${alibabacloudstack_nas_file_system.default.storage_type}",
			"description_regex": "^${alibabacloudstack_nas_file_system.default.description}",
		}),
	}
	protocolTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"protocol_type":     "${alibabacloudstack_nas_file_system.default.protocol_type}",
			"description_regex": "^${alibabacloudstack_nas_file_system.default.description}",
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alibabacloudstack_nas_file_system.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alibabacloudstack_nas_file_system.default.description}_fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nas_file_system.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nas_file_system.default.id}_fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"storage_type":      "${alibabacloudstack_nas_file_system.default.storage_type}",
			"protocol_type":     "${alibabacloudstack_nas_file_system.default.protocol_type}",
			"description_regex": "^${alibabacloudstack_nas_file_system.default.description}",
			"ids":               []string{"${alibabacloudstack_nas_file_system.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alibabacloudstack_nas_file_system.default.description}_fake",
			"ids":               []string{"${alibabacloudstack_nas_file_system.default.id}_fake"},
		}),
	}
	var existFileSystemMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"systems.0.id":            CHECKSET,
			"systems.0.region_id":     CHECKSET,
			"systems.0.description":   name,
			"systems.0.protocol_type": CHECKSET,
			"systems.0.storage_type":  "Capacity",
			"systems.0.metered_size":  CHECKSET,
			"systems.0.create_time":   CHECKSET,
			"ids.#":                   "1",
			"ids.0":                   CHECKSET,
			"descriptions.#":          "1",
			"descriptions.0":          CHECKSET,
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

	fileSystemCheckInfo.dataSourceTestCheck(t, rand, storageTypeConf, protocolTypeConf,
		descriptionConf, idsConf, allConf)
}

func testAccCheckAlibabacloudStackFileSystemDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

`, name, NasCommonTestCase)
}
