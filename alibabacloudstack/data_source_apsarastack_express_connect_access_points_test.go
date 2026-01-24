package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackExpressConnectAccessPointsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_express_connect_access_points.default"
	name := fmt.Sprintf("tf-testacc-expressConnectAccessPoints%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectAccessPointsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_express_connect_access_points.anyone.points.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake"},
		}),
	}
	name_regex := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${data.alibabacloudstack_express_connect_access_points.anyone.points.0.access_point_name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alibabacloudstack_express_connect_access_points.anyone.points.0.id}"},
			"status": "${data.alibabacloudstack_express_connect_access_points.anyone.points.0.status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"fake"},
			"status": "disabled",
		}),
	}

	var existExpressConnectAccessPointsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"names.#":                  "1",
			"points.#":                 "1",
			"points.0.id":              CHECKSET,
			"points.0.access_point_id": CHECKSET,
			"points.0.host_operator":   CHECKSET,
		}
	}

	var fakeExpressConnectAccessPointsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"points.#": "0",
			"names.#":  "0",
			"ids.#":    "0",
		}
	}

	var ExpressConnectAccessPointsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existExpressConnectAccessPointsMapFunc,
		fakeMapFunc:  fakeExpressConnectAccessPointsMapFunc,
	}

	ExpressConnectAccessPointsCheckInfo.dataSourceTestCheck(t, rand, idsConf, name_regex, statusConf)
}

func dataSourceExpressConnectAccessPointsConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		
		data "alibabacloudstack_express_connect_access_points" "anyone" {
		}
		
		`, name)
}
