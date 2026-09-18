package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackHsmsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_hsms.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceHsmsConfigDependence)

	// Basic configuration with required parameters from vendor data source
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vendor_code":  `${data.alibabacloudstack_hsm_vendors.anyone.vendors.0.code}`,
			"product_code": `${data.alibabacloudstack_hsm_vendors.anyone.vendors.0.products.0.code}`,
			"zone_id":      `${data.alibabacloudstack_zones.anyone.zones.0.id}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vendor_code":  `${data.alibabacloudstack_hsm_vendors.anyone.vendors.0.code}`,
			"product_code": `${data.alibabacloudstack_hsm_vendors.anyone.vendors.0.products.0.code}`,
			"zone_id":      "fake-zone-id",
		}),
	}

	var existHsmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     CHECKSET,
			"hsms.#":    CHECKSET,
			"hsms.0.id": CHECKSET,
		}
	}

	var fakeHsmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":  "0",
			"hsms.#": "0",
		}
	}

	var hsmsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existHsmsMapFunc,
		fakeMapFunc:  fakeHsmsMapFunc,
		Providers:    testYunDunProviders(),
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckYundunCredentials(t)
			testAccPreCheckApiAvailable(t, "GET", "hsm-private", "2018-06-30", "GetHsmIds", "")
		},
	}
	hsmsCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func dataSourceHsmsConfigDependence(name string) string {
	return `
data "alibabacloudstack_hsm_vendors" "anyone" {
}

data "alibabacloudstack_zones" "anyone" {
  enable_details = true
}
`
}
