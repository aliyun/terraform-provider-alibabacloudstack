package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackDatahubSubscriptionsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_subscriptions.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_subscriptions.default.id}_fake"]`,
		}),
	}

	project_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_datahub_subscriptions.default.id}"]`,
			"project_name": `"${alibabacloudstack_datahub_subscriptions.default.ProjectName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_datahub_subscriptions.default.id}_fake"]`,
			"project_name": `"${alibabacloudstack_datahub_subscriptions.default.ProjectName}_fake"`,
		}),
	}

	topic_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_datahub_subscriptions.default.id}"]`,
			"topic_name": `"${alibabacloudstack_datahub_subscriptions.default.TopicName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_datahub_subscriptions.default.id}_fake"]`,
			"topic_name": `"${alibabacloudstack_datahub_subscriptions.default.TopicName}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_subscriptions.default.id}"]`,

			"project_name": `"${alibabacloudstack_datahub_subscriptions.default.ProjectName}"`,
			"topic_name":   `"${alibabacloudstack_datahub_subscriptions.default.TopicName}"`}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_subscriptions.default.id}_fake"]`,

			"project_name": `"${alibabacloudstack_datahub_subscriptions.default.ProjectName}_fake"`,
			"topic_name":   `"${alibabacloudstack_datahub_subscriptions.default.TopicName}_fake"`}),
	}

	AlibabacloudstackDatahubSubscriptionsCheckInfo.dataSourceTestCheck(t, rand, idsConf, project_nameConf, topic_nameConf, allConf)
}

var existAlibabacloudstackDatahubSubscriptionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"subscriptions.#":    "1",
		"subscriptions.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackDatahubSubscriptionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"subscriptions.#": "0",
	}
}

var AlibabacloudstackDatahubSubscriptionsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_datahub_subscriptions.default",
	existMapFunc: existAlibabacloudstackDatahubSubscriptionsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackDatahubSubscriptionsMapFunc,
}

func testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackDatahubSubscriptions%d"
}






data "alibabacloudstack_datahub_subscriptions" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
