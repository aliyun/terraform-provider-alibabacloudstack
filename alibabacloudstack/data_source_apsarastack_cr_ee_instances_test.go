package alibabacloudstack

//
//import (
//	"testing"
//)
//
//func TestAccAlibabacloudStackCrEeInstancesDataSource(t *testing.T) {
//	resourceId := "data.alibabacloudstack_cr_ee_instances.default"
//	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "",
//		dataSourceCrEeInstancesConfigDependence)
//
//	nameRegexConf := dataSourceTestAccConfig{
//		existConfig: testAccConfig(map[string]interface{}{}),
//		fakeConfig: testAccConfig(map[string]interface{}{
//			"name_regex": "test-fake.*",
//		}),
//	}
//
//	idsConf := dataSourceTestAccConfig{
//		existConfig: testAccConfig(map[string]interface{}{}),
//		fakeConfig: testAccConfig(map[string]interface{}{
//			"ids": []string{"test-id-fake"},
//		}),
//	}
//
//	allConf := dataSourceTestAccConfig{
//		existConfig: testAccConfig(map[string]interface{}{}),
//		fakeConfig: testAccConfig(map[string]interface{}{
//			"ids":        []string{"test-id-fake"},
//			"name_regex": "test-fake.*",
//		}),
//	}
//
//	var existCrEeInstancesMapFunc = func(rand int) map[string]string {
//		return map[string]string{
//			"names.#":                        CHECKSET,
//			"names.0":                        CHECKSET,
//			"instances.#":                    CHECKSET,
//			"instances.0.id":                 CHECKSET,
//			"instances.0.name":               CHECKSET,
//			"instances.0.region":             CHECKSET,
//			"instances.0.specification":      CHECKSET,
//			"instances.0.namespace_quota":    CHECKSET,
//			"instances.0.namespace_usage":    CHECKSET,
//			"instances.0.repo_quota":         CHECKSET,
//			"instances.0.repo_usage":         CHECKSET,
//			"instances.0.vpc_endpoints.#":    CHECKSET,
//			"instances.0.public_endpoints.#": CHECKSET,
//		}
//	}
//
//	var fakeCrEeInstancesMapFunc = func(rand int) map[string]string {
//		return map[string]string{
//			"ids.#":       "0",
//			"names.#":     "0",
//			"instances.#": "0",
//		}
//	}
//
//	var crEEInstancesCheckInfo = dataSourceAttr{
//		resourceId:   resourceId,
//		existMapFunc: existCrEeInstancesMapFunc,
//		fakeMapFunc:  fakeCrEeInstancesMapFunc,
//	}
//	preCheck := func() {
//		testAccPreCheckWithCrEe(t)
//	}
//	crEEInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, 0, preCheck, nameRegexConf, idsConf, allConf)
//}
//
//func dataSourceCrEeInstancesConfigDependence(name string) string {
//	return ""
//}
