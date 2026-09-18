package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackHsmVendorsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_hsm_vendors.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceHsmVendorsConfigDependence)

	// Basic configuration without filters
	// Note: hsm_vendors data source does not support filtering by ids,
	// it always returns all vendors from the environment
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
	}

	var existHsmVendorsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     CHECKSET,
			"vendors.#":                 CHECKSET,
			"vendors.0.code":            CHECKSET,
			"vendors.0.name":            CHECKSET,
			"vendors.0.products.#":      CHECKSET,
			"vendors.0.products.0.code": CHECKSET,
			"vendors.0.products.0.name": CHECKSET,
		}
	}

	var fakeHsmVendorsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"vendors.#": "0",
		}
	}

	var hsmVendorsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existHsmVendorsMapFunc,
		fakeMapFunc:  fakeHsmVendorsMapFunc,
		Providers:    testYunDunProviders(),
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckYundunCredentials(t)
			testAccPreCheckApiAvailable(t, "GET", "hsm-private", "2018-06-30", "DescribeVendors", "")
		},
	}
	hsmVendorsCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func dataSourceHsmVendorsConfigDependence(name string) string {
	return `
data "alibabacloudstack_hsm_vendors" "anyone" {
}
`
}
