package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"
)

func TestAccAlibabacloudStackEdasK8sApplicationScalingRulesDataSource(t *testing.T) {
	rand := getAccTestRandInt(100, 999)
	resourceId := "data.alibabacloudstack_edas_k8s_application_scaling_rules.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("testtf%d", rand),
		dataSourceEdasK8sApplicationScalingRulesConfigDependence)

	existMapFunc := func(rand int) map[string]string {
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
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"scaling_rules.#": "0",
		}
	}
	checkInfo := dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existMapFunc,
		fakeMapFunc:       fakeMapFunc,
		ExternalProviders: testAccExternalProviders,
	}

	checkInfo.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"app_id": "${alibabacloudstack_edas_k8s_application_scaling_rule.default.app_id}",
				"ids":    []string{"${alibabacloudstack_edas_k8s_application_scaling_rule.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"app_id": "${alibabacloudstack_edas_k8s_application_scaling_rule.default.app_id}",
				"ids":    []string{"${alibabacloudstack_edas_k8s_application_scaling_rule.default.id}_fake"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"app_id":     "${alibabacloudstack_edas_k8s_application_scaling_rule.default.app_id}",
				"name_regex": "${alibabacloudstack_edas_k8s_application_scaling_rule.default.scaling_rule_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"app_id":     "${alibabacloudstack_edas_k8s_application_scaling_rule.default.app_id}",
				"name_regex": "${alibabacloudstack_edas_k8s_application_scaling_rule.default.scaling_rule_name}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"app_id":     "${alibabacloudstack_edas_k8s_application_scaling_rule.default.app_id}",
				"ids":        []string{"${alibabacloudstack_edas_k8s_application_scaling_rule.default.id}"},
				"name_regex": "${alibabacloudstack_edas_k8s_application_scaling_rule.default.scaling_rule_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"app_id":     "${alibabacloudstack_edas_k8s_application_scaling_rule.default.app_id}",
				"ids":        []string{"${alibabacloudstack_edas_k8s_application_scaling_rule.default.id}_fake"},
				"name_regex": "${alibabacloudstack_edas_k8s_application_scaling_rule.default.scaling_rule_name}_fake",
			}),
		},
	)
}

func dataSourceEdasK8sApplicationScalingRulesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

variable "package_version" {
	default = "2025-10-09 17:17:18"
}

resource "alibabacloudstack_edas_k8s_application" "default" {
	application_name        = "${var.name}"
	application_description = "This is description of application"
	cluster_id              = local.edas_cluster_id
	replicas                = 2
	package_type            = "FatJar"
	package_url             = "http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar"
	package_version         = var.package_version
	jdk                     = "Open JDK 8"
	limit_mem               = 1024
	requests_mem            = 1024
	requests_m_cpu          = 300
	limit_m_cpu             = 300
}

resource "alibabacloudstack_edas_k8s_application_scaling_rule" "default" {
	app_id            = "${alibabacloudstack_edas_k8s_application.default.id}"
	scaling_rule_name = "${var.name}"
	scaling_rule_type = "metric"
	max_replicas      = "10"
	min_replicas      = "1"
	metrics {
		type        = "CPU"
		utilization = "80"
	}
	metrics {
		type        = "MEMORY"
		utilization = "80"
	}
}
`, name, EdasClusterCommonTestCase(), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
