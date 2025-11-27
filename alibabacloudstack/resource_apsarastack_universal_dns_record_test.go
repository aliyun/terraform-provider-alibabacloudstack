package alibabacloudstack

import (
	"testing"

	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackUniversalDnsRecord_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_universal_dns_record.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbVservergroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &UniversalDnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeUniversalDNSRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccUniversalDnsRecorddependence)

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
					"zone_id":      "${alibabacloudstack_universal_dns_domain.default.id}",
					"name":         "${var.name}",
					"remark":       "test remark0",
					"type":         "A",
					"ttl":          300,
					"lba_strategy": "RATIO",
					"line_ids":     []string{"default"},
					"rdatas": []map[string]interface{}{
						{
							"value":      "192.168.1.1",
							"lba_weight": "20",
						},
						{
							"value":      "127.0.0.1",
							"lba_weight": "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                name,
						"remark":              "test remark0",
						"type":                "A",
						"ttl":                 "300",
						"line_ids.#":          "1",
						"rdatas.#":            "2",
						"rdatas.0.lba_weight": "20",
						"rdatas.1.lba_weight": "80",
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
					"type":         "CNAME",
					"ttl":          3000,
					"lba_strategy": "RATIO",
					"remark":       "test remark",
					"rdatas": []map[string]interface{}{
						{
							"value":      "test1.",
							"lba_weight": "0",
						},
						{
							"value":      "test2.",
							"lba_weight": "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":                "CNAME",
						"ttl":                 "3000",
						"lba_strategy":        "RATIO",
						"remark":              "test remark",
						"rdatas.#":            "2",
						"rdatas.0.value":      "test1.",
						"rdatas.0.lba_weight": "0",
						"rdatas.1.value":      "test2.",
						"rdatas.1.lba_weight": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lba_strategy": "ALL_RR",
					"rdatas": []map[string]interface{}{
						{
							"value": "test1.",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lba_strategy":        "ALL_RR",
						"rdatas.#":            "1",
						"rdatas.0.value":      "test1.",
						"rdatas.0.lba_weight": REMOVEKEY,
						"rdatas.1.value":      REMOVEKEY,
						"rdatas.1.lba_weight": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func AlibabacloudTestAccUniversalDnsRecorddependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_universal_dns_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS record test"
}

`, name)
}
