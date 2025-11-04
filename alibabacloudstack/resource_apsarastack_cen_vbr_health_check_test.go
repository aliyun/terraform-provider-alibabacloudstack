package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenVbrHealthCheck_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cen_vbr_health_check.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"health_check_interval":  CHECKSET,
		"health_check_source_ip": "172.16.0.1",
		"health_check_target_ip": "10.0.0.1",
		"healthy_threshold":      CHECKSET,
		"health_check_only":      "true",
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCenVbrHealthCheck")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	randInt := fmt.Sprintf("%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, randInt, CenVbrHealthCheckCommonTestCase)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":                 "${alibabacloudstack_cen_transit_router_vbr_attachment.default.cen_id}",
					"vbr_instance_id":        "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
					"health_check_source_ip": "172.16.0.1",
					"health_check_target_ip": "10.0.0.1",
					"healthy_threshold":      3,
					"health_check_interval":  2,
					"health_check_only":      "true",
					"depends_on":             []string{`alibabacloudstack_cen_transit_router_vbr_attachment.default`},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold":     "3",
						"health_check_interval": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func CenVbrHealthCheckCommonTestCase(randInt string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc%s"
}

resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = "%s"
  virtual_border_router_name = var.name
  vlan_id                    = %s
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
    vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

`, randInt, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), randInt)
}
