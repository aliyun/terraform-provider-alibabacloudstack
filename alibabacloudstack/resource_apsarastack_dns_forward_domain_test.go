package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDnsForwardDomain_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dns_forward_domain.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDnsForwardDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfacc%d.test.", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDnsForwardDomaindependence)
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
					"name":         "${var.name}",
					"remark":       "test1234",
					"forward_mode": "FORWARD_FIRST",
					"forwarders":   []string{"192.168.101.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"remark":       "test1234",
						"forward_mode": "FORWARD_FIRST",
						"forwarders.#": "1",
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
					"forward_mode": "FORWARD_ONLY",
					"forwarders":   []string{"192.168.101.1", "172.16.29.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forward_mode": "FORWARD_ONLY",
						"forwarders.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "111111111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "111111111",
					}),
				),
			},
		},
	})
}

func AlibabacloudTestAccDnsForwardDomaindependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}
