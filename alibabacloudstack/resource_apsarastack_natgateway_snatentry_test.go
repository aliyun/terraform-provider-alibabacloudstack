package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNatgatewaySnatentry0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_natgateway_snatentry.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccNatgatewaySnatentryCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribesnattableentriesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snat_gatewaysnat_entry%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccNatgatewaySnatentryBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"snat_ip": "121.199.28.61",

					"snat_table_id": "stb-bp13yy9elsi8zn4i29m59",

					"snat_entry_name": "rdk_test_name",

					"source_cidr": "10.40.0.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"snat_ip": "121.199.28.61",

						"snat_table_id": "stb-bp13yy9elsi8zn4i29m59",

						"snat_entry_name": "rdk_test_name",

						"source_cidr": "10.40.0.0/24",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"snat_entry_name": "rdk_test_udpate_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"snat_entry_name": "rdk_test_udpate_name",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccNatgatewaySnatentryCheckmap = map[string]string{

	"status": CHECKSET,

	"source_cidr": CHECKSET,

	"snat_ip": CHECKSET,

	"snat_table_id": CHECKSET,

	"source_vswitch_id": CHECKSET,

	"snat_entry_name": CHECKSET,

	"snat_entry_id": CHECKSET,
}

func AlibabacloudTestAccNatgatewaySnatentryBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
