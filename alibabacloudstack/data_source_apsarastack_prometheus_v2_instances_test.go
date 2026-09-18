package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackPrometheusV2InstancesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alibabacloudstack_prometheus_v2_instances.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfacc_prometheus%d", rand),
		dataSourcePrometheusV2InstancesConfigDependence)

	dsa := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                        "1",
				"instances.#":                  "1",
				"instances.0.cluster_name":     fmt.Sprintf("tfacc_prometheus%d", rand),
				"instances.0.status":           CHECKSET,
				"instances.0.http_api":         CHECKSET,
				"instances.0.remote_write_url": CHECKSET,
				"instances.0.push_gateway_url": CHECKSET,
				"instances.0.tag_set.#":        "2",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":       "0",
				"instances.#": "0",
			}
		},
	}

	dsa.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_prometheus_v2_instance.default.cluster_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": fmt.Sprintf("nonexistent-%d", rand),
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_prometheus_v2_instance.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"fake-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": fmt.Sprintf("tfacc_prometheus%d", rand),
				"ids":        []string{"${alibabacloudstack_prometheus_v2_instance.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": fmt.Sprintf("nonexistent-%d", rand),
				"ids":        []string{"fake-id"},
			}),
		},
	)
}

func dataSourcePrometheusV2InstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_prometheus_v2_instance" "default" {
	cluster_name = "${var.name}"
	tags         = ["test1", "test2"]
}
`, name)
}
