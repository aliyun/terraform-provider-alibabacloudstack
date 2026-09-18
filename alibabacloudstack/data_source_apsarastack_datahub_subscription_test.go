package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDatahubSubscriptionsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_datahub_subscriptions.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_test_scriptions%d", rand),
		dataSourceDatahubSubscriptionsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_subscription.default.project_name}",
			"topic_name":   "${alibabacloudstack_datahub_subscription.default.topic_name}",
			"sub_ids":      []string{"${alibabacloudstack_datahub_subscription.default.sub_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_subscription.default.project_name}",
			"topic_name":   "${alibabacloudstack_datahub_subscription.default.topic_name}",
			"sub_ids":      []string{"${alibabacloudstack_datahub_subscription.default.sub_id}_fake"},
		}),
	}

	topic_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_subscription.default.project_name}",
			"topic_name":   "${alibabacloudstack_datahub_subscription.default.topic_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_subscription.default.project_name}",
			"topic_name":   "${alibabacloudstack_datahub_subscription.default.topic_name}_fake",
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

func dataSourceDatahubSubscriptionsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
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
`, name)
}
