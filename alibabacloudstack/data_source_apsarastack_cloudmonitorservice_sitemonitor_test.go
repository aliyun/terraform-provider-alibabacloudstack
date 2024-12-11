package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackCloudmonitorserviceSiteMonitorsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCloudmonitorserviceSiteMonitorsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCloudmonitorserviceSiteMonitorsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}_fake"]`,
		}),
	}

	task_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCloudmonitorserviceSiteMonitorsSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}"]`,
			"task_id": `"${alibabacloudstack_cloudmonitorservice_site_monitors.default.TaskId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCloudmonitorserviceSiteMonitorsSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}_fake"]`,
			"task_id": `"${alibabacloudstack_cloudmonitorservice_site_monitors.default.TaskId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCloudmonitorserviceSiteMonitorsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}"]`,

			"task_id": `"${alibabacloudstack_cloudmonitorservice_site_monitors.default.TaskId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackCloudmonitorserviceSiteMonitorsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_site_monitors.default.id}_fake"]`,

			"task_id": `"${alibabacloudstack_cloudmonitorservice_site_monitors.default.TaskId}_fake"`}),
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

func testAccCheckAlibabacloudstackCloudmonitorserviceSiteMonitorsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackCloudmonitorserviceSiteMonitors%d"
}






data "alibabacloudstack_cloudmonitorservice_site_monitors" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
