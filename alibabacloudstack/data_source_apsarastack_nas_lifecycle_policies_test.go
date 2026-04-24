package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNasLifecyclePolicies_DataSource(t *testing.T) {
	rand := getAccTestRandInt(100000, 999999)
	resourceId := "data.alibabacloudstack_nas_lifecycle_policies.default"
	name := fmt.Sprintf("tf-testacc-nas-liecycle-datasource%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccAlibabacloudStackNasLifecyclePoliciesDataSourceConfig)

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_file_system.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "eeeeeeeeee",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nas_lifecycle_policy.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nas_lifecycle_policy.default.id}_fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_file_system.default.id}",
			"ids":            []string{"${alibabacloudstack_nas_lifecycle_policy.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "eeeeeeeeee",
			"ids":            []string{"${alibabacloudstack_nas_lifecycle_policy.default.id}_fake"},
		}),
	}

	LifecyclePoliciesCheckInfo.dataSourceTestCheck(t, rand, descriptionConf, idsConf, allConf)
}

func testAccAlibabacloudStackNasLifecyclePoliciesDataSourceConfig(name string) string {
	clusterFilter := GetOssClusterFilter()
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "${var.name}"
  acl    = "public-read"
  oss_cluster = local.cluster_filter
}

resource "alibabacloudstack_nas_lifecycle_policy" "default" {
	lifecycle_policy_name = "${var.name}"
	file_system_id        = "${alibabacloudstack_nas_file_system.default.id}"
	path                  = "/"
	recursive             = "false"
	lifecycle_rule_name   = "DEFAULT_ATIME_14"
	oss_bucket            = "${alibabacloudstack_oss_bucket.default.id}"

}
data "alibabacloudstack_nas_lifecycle_policies" "default" {
	depends_on = [
		alibabacloudstack_nas_lifecycle_policy.default
	]
}`, name, NasCommonTestCase, clusterFilter)
}

var existLifecyclePoliciesMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"lifecycle_policies.#":                     "1",
		"lifecycle_policies.0.id":                  CHECKSET,
		"lifecycle_policies.0.path":                "/",
		"lifecycle_policies.0.recursive":           "false",
		"lifecycle_policies.0.lifecycle_rule_name": "DEFAULT_ATIME_14",
		"ids.#": "1",
		"ids.0": CHECKSET,
	}
}

var fakeLifecyclePoliciesMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"lifecycle_policies.#": "0",
	}
}

var LifecyclePoliciesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_nas_lifecycle_policies.default",
	existMapFunc: existLifecyclePoliciesMapCheck,
	fakeMapFunc:  fakeLifecyclePoliciesMapCheck,
}
