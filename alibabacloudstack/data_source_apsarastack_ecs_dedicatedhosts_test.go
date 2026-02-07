package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEcsDedicatedHostsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_ecs_dedicated_hosts.default"
	name := fmt.Sprintf("tf_testAccEcsDedicatedHostsDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsDedicatedHostsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_dedicated_host.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_dedicated_host.default.id}-fake"},
		}),
	}

	idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_id": "${alibabacloudstack_ecs_dedicated_host.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_id": "${alibabacloudstack_ecs_dedicated_host.default.id}-fake",
		}),
	}
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_name": "${alibabacloudstack_ecs_dedicated_host.default.dedicated_host_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_name": "${alibabacloudstack_ecs_dedicated_host.default.dedicated_host_name}-fake",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_ecs_dedicated_host.default.id}"},
			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_ecs_dedicated_host.default.id}"},
			"name_regex": name + "fake",
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alibabacloudstack_ecs_dedicated_host.default.id}"},
			"dedicated_host_type": "${alibabacloudstack_ecs_dedicated_host.default.dedicated_host_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alibabacloudstack_ecs_dedicated_host.default.id}"},
			"dedicated_host_type": "${alibabacloudstack_ecs_dedicated_host.default.dedicated_host_type}-fake",
		}),
	}
	zoneConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alibabacloudstack_ecs_dedicated_host.default.zone_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alibabacloudstack_ecs_dedicated_host.default.zone_id}-fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alibabacloudstack_ecs_dedicated_host.default.id}"},
			"dedicated_host_id":   "${alibabacloudstack_ecs_dedicated_host.default.id}",
			"dedicated_host_name": "${alibabacloudstack_ecs_dedicated_host.default.dedicated_host_name}",
			"name_regex":          name,
			"dedicated_host_type": "${alibabacloudstack_ecs_dedicated_host.default.dedicated_host_type}",
			"zone_id":             "${alibabacloudstack_ecs_dedicated_host.default.zone_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alibabacloudstack_ecs_dedicated_host.default.id}-fake"},
			"dedicated_host_id":   "${alibabacloudstack_ecs_dedicated_host.default.id}-fake",
			"dedicated_host_name": "${alibabacloudstack_ecs_dedicated_host.default.dedicated_host_name}-fake",
			"name_regex":          name + "fake",
			"dedicated_host_type": "${alibabacloudstack_ecs_dedicated_host.default.dedicated_host_type}-fake",
			"zone_id":             "${alibabacloudstack_ecs_dedicated_host.default.zone_id}-fake",
		}),
	}

	var existKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"names.#":                     "1",
			"names.0":                     CHECKSET,
			"hosts.0.id":                  CHECKSET,
			"hosts.0.dedicated_host_id":   CHECKSET,
			"hosts.0.dedicated_host_name": CHECKSET,
			"hosts.0.dedicated_host_type": CHECKSET,
			"hosts.0.machine_id":          CHECKSET,
			"hosts.0.physical_gpus":       CHECKSET,
			"hosts.0.sockets":             CHECKSET,
			"hosts.0.supported_instance_types_list.#": CHECKSET,
			"hosts.0.zone_id":                         CHECKSET,
		}
	}

	var fakeKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"hosts.#": "0",
		}
	}

	var ecsDedicatedHostsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsSecretVersionsMapFunc,
		fakeMapFunc:  fakeKmsSecretVersionsMapFunc,
	}

	ecsDedicatedHostsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, typeConf, idConf, nameConf, zoneConf, allConf)
}

func dataSourceEcsDedicatedHostsConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
		resource "alibabacloudstack_ecs_dedicated_host" "default" {
	dedicated_host_name =       var.name
	dedicated_host_type =       local.ddh_type
	dedicated_host_cluster_id = alibabacloudstack_ecs_dedicated_host_cluster.default.id
          tags = {
			Create = "TF"
    		For = "ddh-test",
  			}
		}
	`, AlibabacloudTestAccEcsDedicatedhostBasicdependence(name))
}
