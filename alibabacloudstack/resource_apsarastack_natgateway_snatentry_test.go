package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNatgatewaySnatentry0(t *testing.T) {
	var v vpc.SnatTableEntry

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
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"snat_ip": "${alibabacloudstack_eip.default.ip_address}",

					"snat_table_id": "${alibabacloudstack_nat_gateway.default.snat_table_ids}",

					"source_cidr": "${alibabacloudstack_vpc_vswitch.default.cidr_block}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"snat_ip": CHECKSET,

						"snat_table_id": CHECKSET,

						"source_cidr": CHECKSET,
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

var AlibabacloudTestAccNatgatewaySnatentryCheckmap = map[string]string{

	// "status": CHECKSET,

	// "source_cidr": CHECKSET,

	// "snat_ip": CHECKSET,

	// "snat_table_id": CHECKSET,

	// "source_vswitch_id": CHECKSET,

	// "snat_entry_name": CHECKSET,

	// "snat_entry_id": CHECKSET,
}

func AlibabacloudTestAccNatgatewaySnatentryBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

resource "alibabacloudstack_nat_gateway" "default" {
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alibabacloudstack_eip" "default" {
	name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
	allocation_id = "${alibabacloudstack_eip.default.id}"
	instance_id = "${alibabacloudstack_nat_gateway.default.id}"
}

`, name, VSwitchCommonTestCase)
}
