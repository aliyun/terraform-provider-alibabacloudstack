package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackCspprivateHsmsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_cspprivate_hsms.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceCspprivateHsmsConfigDependence)

	// Basic configuration with required parameters from vendor data source
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vendor_code":  `${data.alibabacloudstack_cspprivate_hsm_vendors.anyone.vendors.0.code}`,
			"product_code": `${data.alibabacloudstack_cspprivate_hsm_vendors.anyone.vendors.0.products.0.code}`,
			"zone_id":      `${data.alibabacloudstack_zones.anyone.zones.0.id}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vendor_code":  `${data.alibabacloudstack_cspprivate_hsm_vendors.anyone.vendors.0.code}`,
			"product_code": `${data.alibabacloudstack_cspprivate_hsm_vendors.anyone.vendors.0.products.0.code}`,
			"zone_id":      "fake-zone-id",
		}),
	}

	var existCspprivateHsmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     CHECKSET,
			"hsms.#":    CHECKSET,
			"hsms.0.id": CHECKSET,
		}
	}

	var fakeCspprivateHsmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":  "0",
			"hsms.#": "0",
		}
	}

	var cspprivateHsmsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCspprivateHsmsMapFunc,
		fakeMapFunc:  fakeCspprivateHsmsMapFunc,
		Providers:    testYunDunProviders(),
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckYundunCredentials(t)
			testAccPreCheckApiAvailable(t, "POST", "Cspprivate", "2022-02-17", "DescribeProxyCryptoService", "")
		},
	}
	cspprivateHsmsCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func dataSourceCspprivateHsmsConfigDependence(name string) string {
	return `
data "alibabacloudstack_cspprivate_hsm_vendors" "anyone" {
}

data "alibabacloudstack_zones" "anyone" {
  enable_details = true
}
`
}
