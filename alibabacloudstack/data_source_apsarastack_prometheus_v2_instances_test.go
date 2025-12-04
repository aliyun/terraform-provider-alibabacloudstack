package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackPrometheusV2InstancesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alibabacloudstack_prometheus_v2_instances.default"
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
			existConfig: testAccCheckAlibabacloudStackPrometheusV2InstancesDataSourceConfig(rand, map[string]string{
				"name_regex": `"${alibabacloudstack_prometheus_v2_instance.default.cluster_name}"`,
			}),
			fakeConfig: testAccCheckAlibabacloudStackPrometheusV2InstancesDataSourceConfig(rand, map[string]string{
				"name_regex": fmt.Sprintf(`"nonexistent-%d"`, rand),
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccCheckAlibabacloudStackPrometheusV2InstancesDataSourceConfig(rand, map[string]string{
				"ids": `["${alibabacloudstack_prometheus_v2_instance.default.id}"]`,
			}),
			fakeConfig: testAccCheckAlibabacloudStackPrometheusV2InstancesDataSourceConfig(rand, map[string]string{
				"ids": `["fake-id"]`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccCheckAlibabacloudStackPrometheusV2InstancesDataSourceConfig(rand, map[string]string{
				"name_regex": fmt.Sprintf(`"%s"`, fmt.Sprintf("tfacc_prometheus%d", rand)),
				"ids":        `["${alibabacloudstack_prometheus_v2_instance.default.id}"]`,
			}),
			fakeConfig: testAccCheckAlibabacloudStackPrometheusV2InstancesDataSourceConfig(rand, map[string]string{
				"name_regex": fmt.Sprintf(`"%s"`, fmt.Sprintf("nonexistent-%d", rand)),
				"ids":        `["fake-id"]`,
			}),
		},
	)
}

func testAccCheckAlibabacloudStackPrometheusV2InstancesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tfacc_prometheus%d"
}

resource "alibabacloudstack_prometheus_v2_instance" "default" {
  cluster_name = "${var.name}"
  tags = ["test1", "test2"]
}

data "alibabacloudstack_prometheus_v2_instances" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
}
