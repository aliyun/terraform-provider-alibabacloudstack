package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlicloudAlikafkaInstancesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_alikafka_instances.default"
	name := fmt.Sprintf("tf-testacc-alikafkainstance%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlikafkaInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alibabacloudstack_alikafka_instance.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alibabacloudstack_alikafka_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_alikafka_instance.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alibabacloudstack_alikafka_instance.default.id}"},
			"name_regex":     "${alibabacloudstack_alikafka_instance.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_alikafka_instance.default.id}_fake"},
			"name_regex": "${alibabacloudstack_alikafka_instance.default.name}_fake",
		}),
	}

	var existAlikafkaInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"instances.#":                            "1",
			"instances.0.name":                       name,
			"instances.0.id":                         CHECKSET,
			"instances.0.replicas":                   CHECKSET,
			"instances.0.disk_num":                   CHECKSET,
			"instances.0.sasl":                       CHECKSET,
			"instances.0.plaintext":                  CHECKSET,
			"instances.0.message_max_bytes":          "10000000",
			"instances.0.num_partitions":             "3",
			"instances.0.auto_create_topics_enable":  "false",
			"instances.0.num_io_threads":             "16",
			"instances.0.queued_max_requests":        "80",
			"instances.0.replica_fetch_wait_max_ms":  "500",
			"instances.0.replica_lag_time_max_ms":    "30000",
			"instances.0.num_network_threads":        "3",
			"instances.0.log_retention_bytes":        "-1",
			"instances.0.replica_fetch_max_bytes":    "10000000",
			"instances.0.num_replica_fetchers":       "4",
			"instances.0.default_replication_factor": "3",
			"instances.0.offsets_retention_minutes":  "10080",
			"instances.0.background_threads":         "10",
		}
	}

	var fakeAlikafkaInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var alikafkaInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAlikafkaInstancesMapFunc,
		fakeMapFunc:  fakeAlikafkaInstancesMapFunc,
	}
	alikafkaInstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)

}

func dataSourceAlikafkaInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
%s

resource "alibabacloudstack_alikafka_instance" "default" {
	name = "${var.name}"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	sasl =      true
	plaintext = true
	spec =      "Broker4C16G"
}

`, name, DataZoneCommonTestCase)
}
