package alibabacloudstack

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackEdasK8sApplicationScalingRulesDataSource(t *testing.T) {
	rand := getAccTestRandInt(100, 999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEdasScalingRulesDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_edas_k8s_application_scaling_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEdasScalingRulesDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_edas_k8s_application_scaling_rule.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEdasScalingRulesDataSourceName(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_edas_k8s_application_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEdasScalingRulesDataSourceName(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_edas_k8s_application_scaling_rule.default.scaling_rule_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEdasScalingRulesDataSourceName(rand, map[string]string{
			"ids":        `["${alibabacloudstack_edas_k8s_application_scaling_rule.default.id}"]`,
			"name_regex": `"${alibabacloudstack_edas_k8s_application_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEdasScalingRulesDataSourceName(rand, map[string]string{
			"ids":        `["${alibabacloudstack_edas_k8s_application_scaling_rule.default.id}_fake"]`,
			"name_regex": `"${alibabacloudstack_edas_k8s_application_scaling_rule.default.scaling_rule_name}_fake"`,
		}),
	}
	var existAlibabacloudStackEdasScalingRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"scaling_rules.#":                   "1",
			"scaling_rules.0.scaling_rule_name": fmt.Sprintf("testtf%d", rand),
			"scaling_rules.0.scaling_rule_type": CHECKSET,
			"scaling_rules.0.max_replicas":      CHECKSET,
			"scaling_rules.0.min_replicas":      CHECKSET,
			"scaling_rules.0.metrics.#":         "2",
		}
	}
	var fakeAlibabacloudStackEdasScalingRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"scaling_rules.#": "0",
		}
	}
	var alibabacloudstackEdasScalingRulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_edas_k8s_application_scaling_rules.default",
		existMapFunc: existAlibabacloudStackEdasScalingRulesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlibabacloudStackEdasScalingRulesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alibabacloudstackEdasScalingRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlibabacloudStackEdasScalingRulesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {	
	default = "testtf%d"
}

%s


variable "package_version" {	
	default = "2025-10-09 17:17:18"
} 

resource "alibabacloudstack_edas_k8s_application" "default" {
  application_name        	= "${var.name}"
  application_description 	= "This is description of application"
  cluster_id              	= local.edas_cluster_id
  replicas                	= 2
  package_type 				= "FatJar"
  package_url     			= "http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar"
  package_version 			= var.package_version
  jdk             			= "Open JDK 8"
  limit_mem             	= 1024
  requests_mem          	= 1024
  requests_m_cpu        	= 300
  limit_m_cpu           	= 300
}

resource "alibabacloudstack_edas_k8s_application_scaling_rule" "default" {
	app_id = "${alibabacloudstack_edas_k8s_application.default.id}"
	scaling_rule_name = "${var.name}"
	scaling_rule_type = "metric"
	max_replicas = "10"
	min_replicas = "1"
	metrics {
		type = "CPU"
		utilization = "80"
	}
	metrics {
		type = "MEMORY"
		utilization = "80"
	}
}

data "alibabacloudstack_edas_k8s_application_scaling_rules" "default" {	
	app_id = "${alibabacloudstack_edas_k8s_application_scaling_rule.default.app_id}"
	%s
}
`, rand, EdasClusterCommonTestCase(), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"), strings.Join(pairs, " \n "))
	return config
}
