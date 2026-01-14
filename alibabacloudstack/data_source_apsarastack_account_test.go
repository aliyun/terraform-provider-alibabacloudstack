package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackAccountDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_account.current"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAccountConfigDependence)

	// Basic configuration (no parameters needed for account data source)
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
		// No fake config since this data source always returns the current account
		// and doesn't support filtering that results in empty results
	}

	var existAccountMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"id":                CHECKSET, // Account ID should be set
			"organization_id":   CHECKSET, // Organization ID should be set
		}
	}

	// Since this data source always returns the current account information,
	// there's no realistic scenario where it returns empty results.
	// However, we provide a minimal fake function for framework compatibility.
	var fakeAccountMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"id":              "",
			"organization_id": "",
		}
	}

	var accountCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAccountMapFunc,
		fakeMapFunc:  fakeAccountMapFunc,
	}
	accountCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func dataSourceAccountConfigDependence(name string) string {
	return `
`
}
