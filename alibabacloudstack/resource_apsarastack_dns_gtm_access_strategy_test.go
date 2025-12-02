package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func buildDnsGtmInstanceTemplate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name = "${var.name}.local."
}

resource "alibabacloudstack_dns_gtm_instance" "default" {
  name     = var.name
  prefix   = var.name
  zone_id  = alibabacloudstack_dns_private_domain.default.id
  ttl      = 300
}

resource "alibabacloudstack_dns_private_line" "default" {
  name     = var.name
  v4_addresses = ["192.168.0.1"]
  v6_addresses = ["2020:148:2:28::", "2020:148:3:28::"]
}

resource "alibabacloudstack_dns_gtm_addresspool" "default" {
  name            = "${var.name}-default"
  type            = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "1.1.1.1"
	mode = "SMART"
  }
}

resource "alibabacloudstack_dns_gtm_addresspool" "failover" {
  name            = "${var.name}-failover"
  type            = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "2.2.2.2"
	mode = "SMART"
  }
}

resource "alibabacloudstack_dns_gtm_addresspool" "update" {
  name            = "${var.name}-update"
  type            = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "3.3.3.3"
	mode = "SMART"
  }
}
`, name)
}

func TestAccAlibabacloudStackDnsGtmAccessStrategy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dns_gtm_access_strategy.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDnsGtmAccessStrategy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, buildDnsGtmInstanceTemplate)

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
					"name":                            "${var.name}",
					"gtm_instance_id":                 "${alibabacloudstack_dns_gtm_instance.default.id}",
					"default_gtm_address_pool_id":     "${alibabacloudstack_dns_gtm_addresspool.default.id}",
					"default_gtm_address_pool_type":   "IPV4",
					"default_min_available_addr_num":  "1",
					"failover_gtm_address_pool_id":    "${alibabacloudstack_dns_gtm_addresspool.failover.id}",
					"failover_gtm_address_pool_type":  "IPV4",
					"failover_min_available_addr_num": "1",
					"switch_mode":                     "BY_PROBE_RESULT",
					"line_ids":                        []string{"${alibabacloudstack_dns_private_line.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                            name,
						"gtm_instance_id":                 CHECKSET,
						"default_gtm_address_pool_id":     CHECKSET,
						"default_gtm_address_pool_type":   "IPV4",
						"default_min_available_addr_num":  "1",
						"failover_gtm_address_pool_id":    CHECKSET,
						"failover_gtm_address_pool_type":  "IPV4",
						"failover_min_available_addr_num": "1",
						"switch_mode":                     "BY_PROBE_RESULT",
						"line_ids.#":                      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"switch_mode":                  "BY_HAND",
					"failover_gtm_address_pool_id": "${alibabacloudstack_dns_gtm_addresspool.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"switch_mode":                  "BY_HAND",
						"failover_gtm_address_pool_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}-updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("%s-updated", name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"failover_min_available_addr_num": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"failover_min_available_addr_num": "2",
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
