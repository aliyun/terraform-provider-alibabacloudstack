package alibabacloudstack

import (
	"testing"
)

// The data source answers for the account the provider is authenticated with, so the test reads
// the RAM role behind ASCM role id 1 - the built-in operations administration role, which every
// deployment grants to the environment administrator.
//
// No prerequisite resource is created on purpose: the organization and the user are resolved by
// the backend from the caller's identity, so an organization or a user created by the test could
// never be the subject of the query. Creating a fresh ASCM custom role would not help either,
// because it only becomes readable once it is granted to the administrator itself, which would
// change the default role of a user shared by the whole environment.
func TestAccAlibabacloudStackAccountAscmCustomRoleArnDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_account_ascm_custom_role_arn.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAccountAscmCustomRoleArnConfigDependence)

	arnConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ascm_role_id": 1,
		}),
	}

	// The ASCM role id is taken from the configuration, so the role id carried by the materialized
	// RAM role name is asserted exactly. The organization id and the user id are assigned by the
	// backend at read time, so only their position in the name is pinned.
	var existAccountAscmCustomRoleArnMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ascm_role_id":                "1",
			"id":                          REGEXMATCH + "^acs:ram::[0-9]+:role/ascm-role-[0-9]+-1-[0-9]+$",
			"arn":                         REGEXMATCH + "^acs:ram::[0-9]+:role/ascm-role-[0-9]+-1-[0-9]+$",
			"role_name":                   REGEXMATCH + "^ASCM-ROLE-[0-9]+-1-[0-9]+$",
			"ram_role_id":                 REGEXMATCH + "^[0-9]+$",
			"assume_role_policy_document": REGEXMATCH + "sts:AssumeRole",
		}
	}

	var fakeAccountAscmCustomRoleArnMapFunc = func(rand int) map[string]string {
		// No fake step is built below: a role the authenticated account does not hold makes the
		// data source return an error instead of an empty result, so there is no empty state to assert.
		return map[string]string{
			"id":  NOSET,
			"arn": NOSET,
		}
	}

	var accountAscmCustomRoleArnCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAccountAscmCustomRoleArnMapFunc,
		fakeMapFunc:  fakeAccountAscmCustomRoleArnMapFunc,
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckApiAvailable(t, "POST", "ascm", "2019-05-10", "GetRamRoleAndRamPolicy", "/ascm/auth/role/getRamRoleAndRamPolicy")
		},
	}
	accountAscmCustomRoleArnCheckInfo.dataSourceTestCheck(t, rand, arnConf)
}

func dataSourceAccountAscmCustomRoleArnConfigDependence(name string) string {
	return ""
}
