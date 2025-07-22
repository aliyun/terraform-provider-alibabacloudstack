package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDrdsInstanceSeriesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_drds_instance_series.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccDrdsInstanceSeriesDataSource_%d", rand),
		dataSourceDrdsInstanceSeriesConfigDependence)

	generationConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"drds.sn2.4c16g"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"drds.sm2.16c128g"},
		}),
	}
	namesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"names": []string{"Basic Edition"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"names": []string{"Basic Puls Pro Max Edition"},
		}),
	}

	var existDrdsInstanceSeriesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         CHECKSET,
			"ids.0":         CHECKSET,
			"series.#":      CHECKSET,
			"series.0.id":   CHECKSET,
			"series.0.name": CHECKSET,
		}
	}

	var fakeDrdsInstanceSeriesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"series.#": "0",
		}
	}

	var DrdsInstanceSeriesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDrdsInstanceSeriesMapFunc,
		fakeMapFunc:  fakeDrdsInstanceSeriesMapFunc,
	}

	DrdsInstanceSeriesCheckInfo.dataSourceTestCheck(t, rand, generationConf, idsConf, namesConf)
}

func dataSourceDrdsInstanceSeriesConfigDependence(name string) string {
	return fmt.Sprintf(`
`)
}
