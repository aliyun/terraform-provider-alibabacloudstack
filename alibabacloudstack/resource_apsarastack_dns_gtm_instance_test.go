package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDnsGtmInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dns_gtm_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDnsGtmInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDnsGtmInstancedependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":    "${var.name}",
					"prefix":  "${var.name}",
					"zone_id": "${alibabacloudstack_dns_private_domain.default.id}",
					"ttl":     "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":    name,
						"prefix":  name,
						"zone_id": CHECKSET,
						"ttl":     "300",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":   "${var.name}_update",
					"prefix": "${var.name}_update",
					"ttl":    "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   fmt.Sprintf("%s_update", name),
						"prefix": fmt.Sprintf("%s_update", name),
						"ttl":    "500",
					}),
				),
			},
		},
	})
}

func AlibabacloudTestAccDnsGtmInstancedependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
resource "alibabacloudstack_vpc_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	vpc_name   = "${var.name}_vpc0"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS record test"
  vpc_ids = [
    "${alibabacloudstack_vpc_vpc.default.id}"
  ]
}


`, name)
}
