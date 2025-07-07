package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackExpressconnectVbrHa_basic0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_expressconnect_vbr_ha.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackExpressconnectVbrHaCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressconnectService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribeVbrHaRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testaccexpressconnect-vbrha%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackExpressconnectVbrHaDependence0(rand))
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID")
			testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID_SECOND")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":     name,
					"vbr_id":      "${alibabacloudstack_express_connect_virtual_border_router.default1.id}",
					"peer_vbr_id": "${alibabacloudstack_express_connect_virtual_border_router.default2.id}",
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     name,
						"description": name,
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

var AlibabacloudStackExpressconnectVbrHaCheckMap = map[string]string{
	"name":     CHECKSET,
	"vbr_id":      CHECKSET,
	"peer_vbr_id": CHECKSET,
}

func AlibabacloudStackExpressconnectVbrHaDependence0(vlanId int) func(string) string {
	return func(name string) string {
		return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default1" {
	physical_connection_id =     "%s"
	vlan_id =                    %d
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}_1"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default2" {
	physical_connection_id =     "%s"
	vlan_id =                    %d
	local_gateway_ip =           "10.0.1.1"
	peer_gateway_ip =            "10.0.1.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}_2"
}


`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlanId, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID_SECOND"), vlanId+3)
	}
}
