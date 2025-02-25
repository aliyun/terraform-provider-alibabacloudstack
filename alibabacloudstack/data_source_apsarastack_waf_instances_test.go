package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackWafInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_waf_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("tf_testAcc%d", rand),
		testAccCheckAlicloudWafInstanceDataSourceConfig)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"waf-25d7f7889MuSlV9v2n"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"waf-25d7f7889MuSlV9v2n-fake"},
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"instances.#":                   "1",
			"instances.0.id":                CHECKSET,
			"instances.0.instance_id":       CHECKSET,
			"instances.0.end_date":          CHECKSET,
			"instances.0.in_debt":           CHECKSET,
			"instances.0.remain_day":        CHECKSET,
			"instances.0.status":            "1",
			"instances.0.subscription_type": "Subscription",
			"instances.0.trial":             CHECKSET,
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var wafInstancesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_waf_instances.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	var perCheck = func() {
		testAccPreCheck(t)
	}

	wafInstancesRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, perCheck, idsConf)

}

func testAccCheckAlicloudWafInstanceDataSourceConfig(description string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}				
// data "alibabacloudstack_waf_instances" "default" {
//   }
`, description)
}
