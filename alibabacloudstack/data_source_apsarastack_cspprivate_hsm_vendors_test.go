package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackCspprivateVendorsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_cspprivate_hsm_vendors.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceCspprivateHsmVendorsConfigDependence)

	// Basic configuration without filters
	// Note: cspprivate_hsm_vendors data source does not support filtering by ids,
	// it always returns all vendors from the environment
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
	}

	var existCspprivateVendorsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     CHECKSET,
			"vendors.#":                 CHECKSET,
			"vendors.0.id":              CHECKSET,
			"vendors.0.code":            CHECKSET,
			"vendors.0.name":            CHECKSET,
			"vendors.0.products.#":      CHECKSET,
			"vendors.0.products.0.code": CHECKSET,
			"vendors.0.products.0.name": CHECKSET,
		}
	}

	var fakeCspprivateVendorsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"vendors.#": "0",
		}
	}

	var cspprivateVendorsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCspprivateVendorsMapFunc,
		fakeMapFunc:  fakeCspprivateVendorsMapFunc,
		Providers:    testYunDunProviders(),
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckYundunCredentials(t)
			testAccPreCheckApiAvailable(t, "POST", "Cspprivate", "2022-02-17", "DescribeProxyCryptoService", "")
		},
	}
	cspprivateVendorsCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func dataSourceCspprivateHsmVendorsConfigDependence(name string) string {
	return `
data "alibabacloudstack_cspprivate_hsm_vendors" "anyone" {
}
`
}
