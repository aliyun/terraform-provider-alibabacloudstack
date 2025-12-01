package alibabacloudstack

import (
	"testing"

	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDnsPrivateLine_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dns_private_line.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbVservergroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribePrivateLine")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDnsPrivateLinedependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rc.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{

					"name":         "${var.name}",
					"v4_addresses": []string{"192.168.0.1"},
					"v6_addresses": []string{"2020:148:2:28::", "2020:148:3:28::"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"v4_addresses.#": "1",
						"v6_addresses.#": "2",
						"v4_addresses.0": "192.168.0.1",
						"v6_addresses.0": "2020:148:2:28::",
						"v6_addresses.1": "2020:148:3:28::",
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

					"name":         "${var.name}_update",
					"v4_addresses": []string{"192.168.2.1", "127.0.0.1"},
					"v6_addresses": []string{"2020:148:4:28::"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           fmt.Sprintf("%s_update", name),
						"v4_addresses.#": "2",
						"v6_addresses.#": "1",
						"v4_addresses.0": CHECKSET,
						"v4_addresses.1": CHECKSET,
						"v6_addresses.0": "2020:148:4:28::",
						"v6_addresses.1": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func AlibabacloudTestAccDnsPrivateLinedependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}
