package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCmsProjectMetaDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_cms_project_meta.default"
	name := fmt.Sprintf("tf-testacc-cmsprojectmeta%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCmsProjectMetaConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": ".*ECS.*",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace": "acs_ecs_dashboard",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace": "fake_namespace",
		}),
	}

	var existCmsProjectMetaMapFunc = func(rand int) map[string]string {
		// We expect at least one resource matching "acs_ecs" pattern
		// The exact namespace may vary, but we check the structure
		return map[string]string{
			"resources.#":             CHECKSET, // At least one resource expected
			"resources.0.namespace":   CHECKSET,
			"resources.0.description": CHECKSET,
			"resources.0.labels.#":    CHECKSET, // Could be 0 or more, but field exists
		}
	}

	var fakeCmsProjectMetaMapFunc = func(rand int) map[string]string {
		// No resources should match fake pattern
		return map[string]string{
			"resources.#": "0",
		}
	}

	var cmsProjectMetaCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCmsProjectMetaMapFunc,
		fakeMapFunc:  fakeCmsProjectMetaMapFunc,
	}
	cmsProjectMetaCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, namespaceConf)
}

func dataSourceCmsProjectMetaConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

// This data source doesn't require any dependencies as it queries system-wide CMS project metadata
`, name)
}
