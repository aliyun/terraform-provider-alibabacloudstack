package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackKVStoreInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_kvstore_instances.default"
	name := fmt.Sprintf("tf-kvins%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceKVStoreInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${local.kv_instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-nonexistent-instance",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${local.kv_instance_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-instance-id-12345"},
		}),
	}

	instanceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_type": "Redis",
			"name_regex":    "${local.kv_instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_type": "Memcache",
			"name_regex":    "${local.kv_instance_name}",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${local.kv_instance_id}"},
			"name_regex":    "${local.kv_instance_name}",
			"instance_type": "Redis",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"fake-instance-id-12345"},
			"name_regex":    "another-fake-instance",
			"instance_type": "Memcache",
		}),
	}

	var existKVStoreInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "1",
			"names.#":    "1",
			"instances.#": "1",
			"instances.0.name":                CHECKSET,
			"instances.0.id":                  CHECKSET,
			"instances.0.instance_type":       CHECKSET,
			"instances.0.status":              CHECKSET,
			"instances.0.instance_class":      CHECKSET,
			"instances.0.availability_zone":   CHECKSET,
			"instances.0.region_id":           CHECKSET,
			"instances.0.create_time":         CHECKSET,
			"instances.0.connection_domain":    CHECKSET,
			// Numeric fields should be set but exact values may vary
			"instances.0.bandwidth":     CHECKSET,
			"instances.0.connections":   CHECKSET,
			"instances.0.capacity":      CHECKSET,
			"instances.0.port":          CHECKSET,
		}
	}

	var fakeKVStoreInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"instances.#": "0",
		}
	}

	var kvStoreInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKVStoreInstancesMapFunc,
		fakeMapFunc:  fakeKVStoreInstancesMapFunc,
		ExternalProviders: testAccExternalProviders,
	}
	kvStoreInstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, instanceTypeConf, allConf)
}

func dataSourceKVStoreInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

`, name, KVRInstanceCommonTestCase("enterprise", string(KVStoreRedis), ))
}
