package alibabacloudstack

import (
	"testing"

	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDnsRecursorAcl_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dns_recursor_acl.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDnsRecursorAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDnsRecursorAcldependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{

					"name":     "${var.name}",
					"remark":   "${var.name}",
					"policy":   "ALLOW",
					"line_ids": []string{"${alibabacloudstack_dns_line.ipv4.id}", "${alibabacloudstack_dns_line.ipv6.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":       name,
						"remark":     name,
						"policy":     "ALLOW",
						"line_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"name":     "${var.name}_update",
					"remark":   "${var.name}_update",
					"policy":   "FORBID",
					"line_ids": []string{"${alibabacloudstack_dns_line.ipv4.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":       fmt.Sprintf("%s_update", name),
						"remark":     fmt.Sprintf("%s_update", name),
						"policy":     "FORBID",
						"line_ids.#": "1",
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

func AlibabacloudTestAccDnsRecursorAcldependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_dns_line" "ipv4" {
  name = "${var.name}ipv4"
  v4_addresses = ["192.168.0.1"]
}

resource "alibabacloudstack_dns_line" "ipv6" {
  name = "${var.name}ipv6"
  v6_addresses = ["2020:148:2:28::"]
}

`, name)
}
