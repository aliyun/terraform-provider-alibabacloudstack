package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCloudmonitorserviceSiteMonitorsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_cloudmonitorservice_site_monitors.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAlibabacloudstackCloudmonitorserviceSiteMonitors%d", rand),
		dataSourceCmsSiteMonitorsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}_fake"},
		}),
	}

	task_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}"},
			"task_id": "${alibabacloudstack_cloudmonitorservice_site_monitors.default.TaskId}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}_fake"},
			"task_id": "${alibabacloudstack_cloudmonitorservice_site_monitors.default.TaskId}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}"},
			"task_id": "${alibabacloudstack_cloudmonitorservice_site_monitors.default.TaskId}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":     []string{"${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}_fake"},
			"task_id": "${alibabacloudstack_cloudmonitorservice_site_monitors.default.TaskId}_fake",
		}),
	}

	AlibabacloudstackCloudmonitorserviceSiteMonitorsCheckInfo.dataSourceTestCheck(t, rand, idsConf, task_idConf, allConf)
}

var existAlibabacloudstackCloudmonitorserviceSiteMonitorsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"monitors.#":    "1",
		"monitors.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackCloudmonitorserviceSiteMonitorsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"monitors.#": "0",
	}
}

var AlibabacloudstackCloudmonitorserviceSiteMonitorsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cloudmonitorservice_site_monitors.default",
	existMapFunc: existAlibabacloudstackCloudmonitorserviceSiteMonitorsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackCloudmonitorserviceSiteMonitorsMapFunc,
}

func dataSourceCmsSiteMonitorsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
