package alibabacloudstack

import (
	"testing"

	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDnsGtmAddressPool_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dns_gtm_addresspool.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDnsGtmAddressPool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDnsLinedependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":         "${var.name}",
					"type":         "A",
					"lba_strategy": "RATIO",
					"addrs": []map[string]interface{}{
						{
							"value":      "192.168.1.1",
							"mode":       "SMART",
							"lba_weight": 20,
						},
						{
							"value":      "127.0.0.1",
							"mode":       "SMART",
							"lba_weight": 80,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"type":               "A",
						"lba_strategy":       "RATIO",
						"addrs.#":            "2",
						"addrs.0.value":      "192.168.1.1",
						"addrs.0.mode":       "SMART",
						"addrs.0.lba_weight": "20",
						"addrs.1.value":      "127.0.0.1",
						"addrs.1.mode":       "SMART",
						"addrs.1.lba_weight": "80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}_update",
					"addrs": []map[string]interface{}{
						{
							"value":      "192.168.2.1",
							"mode":       "ONLINE",
							"lba_weight": 50,
						},
						{
							"value":      "176.16.1.1",
							"mode":       "OFFLINE",
							"lba_weight": 50,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               fmt.Sprintf("%s_update", name),
						"addrs.0.value":      "192.168.2.1",
						"addrs.0.mode":       "ONLINE",
						"addrs.0.lba_weight": "50",
						"addrs.1.value":      "176.16.1.1",
						"addrs.1.mode":       "OFFLINE",
						"addrs.1.lba_weight": "50",
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

func TestAccAlibabacloudStackDnsGtmAddressPool_ipv6(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dns_gtm_addresspool.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDnsGtmAddressPool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDnsLinedependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":         "${var.name}",
					"type":         "AAAA",
					"lba_strategy": "ALL_RR",
					"addrs": []map[string]interface{}{
						{
							"value": "2020:148:2:28::",
							"mode":  "SMART",
						},
						{
							"value": "020:148:3:28::",
							"mode":  "SMART",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":          name,
						"type":          "AAAA",
						"lba_strategy":  "ALL_RR",
						"addrs.#":       "2",
						"addrs.0.value": "2020:148:2:28::",
						"addrs.0.mode":  "SMART",
						"addrs.1.value": "020:148:3:28::",
						"addrs.1.mode":  "SMART",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}_update",
					"addrs": []map[string]interface{}{
						{
							"value": "2020:148:4:28::",
							"mode":  "ONLINE",
						},
						{
							"value": "020:148:5:28::",
							"mode":  "OFFLINE",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":          fmt.Sprintf("%s_update", name),
						"addrs.0.value": "2020:148:4:28::",
						"addrs.0.mode":  "ONLINE",
						"addrs.1.value": "020:148:5:28::",
						"addrs.1.mode":  "OFFLINE",
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

// buildDnsGtmAddressPoolBasicDependency creates the necessary dependent resources for testing alibabacloudstack_dns_gtm_addresspool.
func buildDnsGtmAddressPoolBasicDependency(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
