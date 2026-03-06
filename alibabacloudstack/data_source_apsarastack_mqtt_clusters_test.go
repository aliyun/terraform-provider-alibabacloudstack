package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackMqttClustersDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_mqtt_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceMqttClustersConfigDependence)

	// Basic configuration without filters
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
		// No fake config since this data source queries system-wide Mqtt clusters
		// and doesn't support filtering that results in empty lists
	}

	// Test with name_regex filter
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{`${data.alibabacloudstack_mqtt_clusters.anyone.clusters.0.id}`},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-nonexistent-cluster"},
		}),
	}

	// Test with name_regex filter
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": `^${data.alibabacloudstack_mqtt_clusters.anyone.clusters.0.name}$`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-nonexistent-cluster",
		}),
	}

	var existMqttClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           CHECKSET, // Should contain at least one cluster
			"clusters.#":      CHECKSET, // Should contain at least one cluster
			"clusters.0.id":   CHECKSET,
			"clusters.0.name": CHECKSET,
		}
	}

	// Since this data source queries system-wide Mqtt clusters, fake configurations
	// with non-existent filters should return empty results
	var fakeMqttClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"clusters.#": "0",
		}
	}

	var MqttClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMqttClustersMapFunc,
		fakeMapFunc:  fakeMqttClustersMapFunc,
	}
	MqttClustersCheckInfo.dataSourceTestCheck(t, 0, basicConf, idsConf, nameRegexConf)
}

func dataSourceMqttClustersConfigDependence(name string) string {
	return `
data "alibabacloudstack_mqtt_clusters" "anyone" {
}
`
}
