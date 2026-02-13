package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDataWorksFileTypesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_dataworks_file_types.default"
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf_filetype%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDataWorksFileTypesConfigDependence)

	// Test with name parameter
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_id": "${alibabacloudstack_data_works_project.default.id}",
			"name":       "${data.alibabacloudstack_dataworks_file_types.anyone.names.0}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_id": "${alibabacloudstack_data_works_project.default.id}",
			"name":       "fake-file-type-name-12345",
		}),
	}

	// Test with ids parameter
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_id": "${alibabacloudstack_data_works_project.default.id}",
			"ids":        []string{"${data.alibabacloudstack_dataworks_file_types.anyone.ids.0}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_id": "${alibabacloudstack_data_works_project.default.id}",
			"ids":        []string{"999999"},
		}),
	}

	var existDataWorksFileTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       CHECKSET,
			"names.#":                     CHECKSET,
			"file_types.#":                CHECKSET,
			"file_types.0.node_type_name": CHECKSET,
			"file_types.0.node_type_id":   CHECKSET,
		}
	}

	var fakeDataWorksFileTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"file_types.#": "0",
		}
	}

	var dataWorksFileTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDataWorksFileTypesMapFunc,
		fakeMapFunc:  fakeDataWorksFileTypesMapFunc,
	}

	// Test basic functionality
	dataWorksFileTypesCheckInfo.dataSourceTestCheck(t, 0, nameConf, idsConf)

}

func dataSourceDataWorksFileTypesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable name {
		default = "%s"
	}
	resource "alibabacloudstack_data_works_project" "default" {
		name =           "${var.name}"
		description =    "${var.name}_desc"
		task_auth_type = "PROJECT"
		}	

data "alibabacloudstack_dataworks_file_types" "anyone" {
	project_id = "${alibabacloudstack_data_works_project.default.id}"
}

`, name)
}
