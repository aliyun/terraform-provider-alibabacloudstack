package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbBackendServersDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_backend_servers.default"
	name := fmt.Sprintf("tf-testacc-slbbackendservers%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSlbBackendServersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_backend_server.default.load_balancer_id}",
			"ids":              []string{"${alibabacloudstack_ecs_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_backend_server.default.load_balancer_id}",
			"ids":              []string{"fake_id"},
		}),
	}

	var existSlbBackendServersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"backend_servers.#":        "1",
			"backend_servers.0.id":     CHECKSET,
			"backend_servers.0.weight": "100",
			"load_balancer_id":         CHECKSET,
		}
	}

	var fakeSlbBackendServersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"backend_servers.#": "0",
			"load_balancer_id":  CHECKSET,
		}
	}

	var slbBackendServersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbBackendServersMapFunc,
		fakeMapFunc:  fakeSlbBackendServersMapFunc,
	}
	slbBackendServersCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceSlbBackendServersConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	%s

	resource "alibabacloudstack_slb" "default" {
		name = "${var.name}"
		vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	}

	resource "alibabacloudstack_slb_backend_server" "default" {
		load_balancer_id = "${alibabacloudstack_slb.default.id}"

		backend_servers {
			server_id = "${alibabacloudstack_ecs_instance.default.id}"
			weight    = 100
		}
	}
	`, name, ECSInstanceCommonTestCase)
}
