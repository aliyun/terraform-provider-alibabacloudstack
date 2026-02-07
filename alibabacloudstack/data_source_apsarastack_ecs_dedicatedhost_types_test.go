package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackEcsDedicatedHostTypesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ecs_dedicated_host_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceDdhTypesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_ecs_dedicated_host_types.anyone.ids.0}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_ecs_dedicated_host_types.anyone.ids.0}_fake"},
		}),
	}

	azConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"availability_zone": "${data.alibabacloudstack_ecs_dedicated_host_types.anyone.ddh_types.0.availability_zones.0}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"availability_zone": "fake_az",
		}),
	}

	var existInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          CHECKSET, // Should contain at least one instance type
			"ddh_types.#":    CHECKSET, // Should contain at least one instance type
			"ddh_types.0.id": CHECKSET,
		}
	}

	var fakeInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"ddh_types.#": "0",
		}
	}

	var instanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existInstanceTypesMapFunc,
		fakeMapFunc:  fakeInstanceTypesMapFunc,
	}
	instanceTypesCheckInfo.dataSourceTestCheck(t, 0, idsConf, azConf)
}

func dataSourceDdhTypesConfigDependence(name string) string {
	return `
data "alibabacloudstack_ecs_dedicated_host_types" "anyone" {
}
`
}
