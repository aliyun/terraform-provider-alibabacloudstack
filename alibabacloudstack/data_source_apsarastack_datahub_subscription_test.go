package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackDatahubSubscriptionsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"sub_ids": `["${alibabacloudstack_datahub_subscription.default.sub_id}"]`,
			"project_name": `"${alibabacloudstack_datahub_subscription.default.project_name}"`,
			"topic_name": `"${alibabacloudstack_datahub_subscription.default.topic_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"sub_ids": `["${alibabacloudstack_datahub_subscription.default.sub_id}_fake"]`,
			"project_name": `"${alibabacloudstack_datahub_subscription.default.project_name}"`,
			"topic_name": `"${alibabacloudstack_datahub_subscription.default.topic_name}"`,
		}),
	}

	topic_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"project_name": `"${alibabacloudstack_datahub_subscription.default.project_name}"`,
			"topic_name": `"${alibabacloudstack_datahub_subscription.default.topic_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubSubscriptionsSourceConfig(rand, map[string]string{
			"project_name": `"${alibabacloudstack_datahub_subscription.default.project_name}"`,
			"topic_name": `"${alibabacloudstack_datahub_subscription.default.topic_name}_fake"`,
		}),
	}

	AlibabacloudstackDatahubSubscriptionsCheckInfo.dataSourceTestCheck(t, rand, idsConf, topic_nameConf)
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
	default = "tf_test_scriptions%d"
}

resource "alibabacloudstack_datahub_project" "default" {
  comment = "test"
  name    = var.name
}

resource "alibabacloudstack_datahub_topic" "default" {
  name = var.name
  comment      = "test"
  record_type  = "BLOB"
  project_name = alibabacloudstack_datahub_project.default.name
}
resource "alibabacloudstack_datahub_subscription" "default" {
  comment          = "terraform subscription datasource test"
  application_name = "tf_acc_test"
  project_name     = "${alibabacloudstack_datahub_project.default.name}"
  topic_name       = "${alibabacloudstack_datahub_topic.default.name}"
}

data "alibabacloudstack_datahub_subscriptions" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
