package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDrdsInstanceSpecificationsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_drds_instance_specifications.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccDrdsInstanceSpecificationsDataSource_%d", rand),
		dataSourceDrdsInstanceSpecificationsConfigDependence)

	cpuSortConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"sorted_by": "CPU",
		}),
	}
	memorySortConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"sorted_by": "Memory",
		}),
	}

	seriesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"series": "${data.alibabacloudstack_drds_instance_series.default.series.0.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"series": "xxxxx",
		}),
	}

	testAccConfig = dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccDrdsInstanceSpecificationsDataSource_%d", rand),
		dataSourceDrdsInstanceSpecificationPresetDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_drds_instance_specifications.preset.specifications.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"xxxxx"},
		}),
	}
	namesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"names": []string{"${data.alibabacloudstack_drds_instance_specifications.preset.specifications.0.name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"names": []string{"xxxxx"},
		}),
	}
	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu": "8",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu": "800",
		}),
	}
	memoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"memory": "32",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"memory": "3200",
		}),
	}

	var existDrdsInstanceSpecificationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   CHECKSET,
			"ids.0":                   CHECKSET,
			"specifications.#":        CHECKSET,
			"specifications.0.id":     CHECKSET,
			"specifications.0.name":   CHECKSET,
			"specifications.0.series": CHECKSET,
		}
	}

	var fakeDrdsInstanceSpecificationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"specifications.#": "0",
		}
	}

	var DrdsInstanceSpecificationsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDrdsInstanceSpecificationsMapFunc,
		fakeMapFunc:  fakeDrdsInstanceSpecificationsMapFunc,
	}

	DrdsInstanceSpecificationsCheckInfo.dataSourceTestCheck(t, rand, cpuSortConf, memorySortConf, seriesConf, idsConf, namesConf, cpuConf, memoryConf)
}

func dataSourceDrdsInstanceSpecificationsConfigDependence(name string) string {
	return fmt.Sprintf(`
	data "alibabacloudstack_drds_instance_series" "default" {
	}
`)
}

func dataSourceDrdsInstanceSpecificationPresetDependence(name string) string {
	return fmt.Sprintf(`
	data "alibabacloudstack_drds_instance_series" "default" {
	}
	
	data "alibabacloudstack_drds_instance_specifications" "preset" {
		series = data.alibabacloudstack_drds_instance_series.default.series.0.id
	}
`)
}
