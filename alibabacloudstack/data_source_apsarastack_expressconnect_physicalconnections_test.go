package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackExpressconnectPhysicalconnectionsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_expressconnect_physical_connections.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		"expressconnect-pc",
		dataSourceExpressconnectPhysicalConnectionsConfigDependence)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"connections.#":    "1",
			"connections.0.id": CHECKSET,
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"connections.#": "0",
		}
	}
	checkInfo := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	checkInfo.dataSourceTestCheck(t, 0,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${data.alibabacloudstack_expressconnect_physical_connections.anyone.ids.0}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${data.alibabacloudstack_expressconnect_physical_connections.anyone.ids.0}_fake"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^${data.alibabacloudstack_expressconnect_physical_connections.anyone.connections.0.physical_connection_name}$",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^tf_fake_fake$",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^${data.alibabacloudstack_expressconnect_physical_connections.anyone.connections.0.physical_connection_name}$",
				"ids":        []string{"${data.alibabacloudstack_expressconnect_physical_connections.anyone.ids.0}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^${data.alibabacloudstack_expressconnect_physical_connections.anyone.connections.0.physical_connection_name}_fake$",
				"ids":        []string{"${data.alibabacloudstack_expressconnect_physical_connections.anyone.ids.0}_fake"},
			}),
		},
	)
}

func dataSourceExpressconnectPhysicalConnectionsConfigDependence(name string) string {
	return `
data "alibabacloudstack_expressconnect_physical_connections" "anyone" {
	name_regex               = ".+"
	include_reservation_data = false
}
`
}
