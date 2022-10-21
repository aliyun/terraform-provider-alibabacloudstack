package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackExpressConnectPhysicalConnectionsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "data.alibabacloudstack_express_connect_physical_connections.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccExpressConnectPhysicalConnectionsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectPhysicalConnectionsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_express_connect_physical_connection.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_express_connect_physical_connection.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alibabacloudstack_express_connect_physical_connection.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
			"ids":        []string{"${alibabacloudstack_express_connect_physical_connection.default.id}-fake"},
		}),
	}
	var existExpressConnectPhysicalConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                          "1",
			"ids.0":                                          CHECKSET,
			"names.#":                                        "1",
			"names.0":                                        name,
			"connections.#":                                  "1",
			"connections.0.id":                               CHECKSET,
			"connections.0.access_point_id":                  CHECKSET,
			"connections.0.ad_location":                      CHECKSET,
			"connections.0.bandwidth":                        CHECKSET,
			"connections.0.business_status":                  CHECKSET,
			"connections.0.circuit_code":                     "",
			"connections.0.create_time":                      CHECKSET,
			"connections.0.description":                      CHECKSET,
			"connections.0.enabled_time":                     "",
			"connections.0.end_time":                         "1970-01-01T00:00:00Z",
			"connections.0.has_reservation_data":             CHECKSET,
			"connections.0.line_operator":                    CHECKSET,
			"connections.0.loa_status":                       "",
			"connections.0.payment_type":                     "AfterPay",
			"connections.0.peer_location":                    CHECKSET,
			"connections.0.physical_connection_id":           CHECKSET,
			"connections.0.physical_connection_name":         CHECKSET,
			"connections.0.port_number":                      CHECKSET,
			"connections.0.port_type":                        CHECKSET,
			"connections.0.redundant_physical_connection_id": "",
			"connections.0.reservation_active_time":          "",
			"connections.0.reservation_internet_charge_type": "",
			"connections.0.reservation_order_type":           "",
			"connections.0.spec":                             CHECKSET,
			"connections.0.status":                           CHECKSET,
			"connections.0.type":                             CHECKSET,
		}
	}

	var fakeExpressConnectPhysicalConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"names.#":       "0",
			"connections.#": "0",
		}
	}

	var ExpressConnectPhysicalConnectionsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existExpressConnectPhysicalConnectionsMapFunc,
		fakeMapFunc:  fakeExpressConnectPhysicalConnectionsMapFunc,
	}

	ExpressConnectPhysicalConnectionsInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf, allConf)
}

func dataSourceExpressConnectPhysicalConnectionsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alibabacloudstack_express_connect_physical_connection" "default" {
  access_point_id          = "ap-cn-qingdao-env17-d01-amtest17"
  device_name			   = "CSW-VM-VPC-G1.AMTEST17"
  line_operator            = "CO"
  peer_location            = var.name
  physical_connection_name = var.name
  type                     = "VPC"
  description              = "my domestic connection"
  port_type                = "10GBase-LR"
  bandwidth                = 100
}`, name)
}
