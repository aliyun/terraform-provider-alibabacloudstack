package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackRouterInterfacesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_router_interfaces.default"
	name := fmt.Sprintf("tf-testacc-routerinterfaces%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRouterInterfacesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_router_interface.initiating.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-id"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_router_interface.initiating.name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-regex",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_router_interface.initiating.id}"},
			"status": "${alibabacloudstack_router_interface.initiating.status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_router_interface.initiating.id}"},
			"status": "fake_status",
		}),
	}

	specificationConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_router_interface.initiating.id}"},
			"specification": "Large.2",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_router_interface.initiating.id}"},
			"specification": "Small.1",
		}),
	}

	routerIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_router_interface.initiating.id}"},
			"router_id": "${alibabacloudstack_vpc.default.0.router_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_router_interface.initiating.id}"},
			"router_id": "fake-router-id",
		}),
	}

	roleConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_router_interface.initiating.id}"},
			"role": "InitiatingSide",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_router_interface.initiating.id}"},
			"role": "AcceptingSide",
		}),
	}

	var existRouterInterfacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"interfaces.#":              "1",
			"interfaces.0.id":           CHECKSET,
			"interfaces.0.status":       CHECKSET,
			"interfaces.0.name":         name + "_initiating",
			"interfaces.0.description":  name + "_decription",
			"interfaces.0.role":         "InitiatingSide",
			"interfaces.0.specification": "Large.2",
			"interfaces.0.router_id":    CHECKSET,
			"interfaces.0.router_type":  "VRouter",
			"interfaces.0.vpc_id":       CHECKSET,
			"interfaces.0.creation_time": CHECKSET,
			"interfaces.0.opposite_region_id":        CHECKSET,
			"interfaces.0.opposite_router_type":      "VRouter",
			"interfaces.0.opposite_interface_owner_id": CHECKSET,
		}
	}

	var fakeRouterInterfacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"interfaces.#": "0",
		}
	}

	var routerInterfacesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRouterInterfacesMapFunc,
		fakeMapFunc:  fakeRouterInterfacesMapFunc,
	}
	routerInterfacesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, specificationConf, routerIdConf, roleConf)
}

func dataSourceRouterInterfacesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_vpc" "default" {
  count      = 2
  name       = var.name
  cidr_block = element(["172.16.0.0/12", "192.168.0.0/16"], count.index)
}

data "alibabacloudstack_account" "current"{
}

resource "alibabacloudstack_router_interface" "initiating" {
  opposite_region               = data.alibabacloudstack_account.current.region
  router_type                   = "VRouter"
  router_id                     = alibabacloudstack_vpc.default.0.router_id
  role                          = "InitiatingSide"
  specification                 = "Large.2"
  name                          = "${var.name}_initiating"
  description                   = "${var.name}_decription"
}

resource "alibabacloudstack_router_interface" "opposite" {
  opposite_region               = data.alibabacloudstack_account.current.region
  router_type                   = "VRouter"
  router_id                     = alibabacloudstack_vpc.default.1.router_id
  role                          = "AcceptingSide"
  specification                 = "Large.1"
  name                          = "${var.name}_opposite"
  description                   = "${var.name}_decription"
}

resource "alibabacloudstack_router_interface_connection" "bar" {
  interface_id                  = alibabacloudstack_router_interface.opposite.id
  opposite_interface_id         = alibabacloudstack_router_interface.initiating.id
  opposite_router_id            = alibabacloudstack_vpc.default.1.router_id
  opposite_router_type          = "VRouter"
}

resource "alibabacloudstack_router_interface_connection" "foo" {
  interface_id                  = alibabacloudstack_router_interface.initiating.id
  opposite_interface_id         = alibabacloudstack_router_interface.opposite.id
  opposite_router_id            = alibabacloudstack_vpc.default.0.router_id
  opposite_router_type          = "VRouter"
  depends_on                    = [alibabacloudstack_router_interface_connection.bar]
}

`, name, DataZoneCommonTestCase)
}
