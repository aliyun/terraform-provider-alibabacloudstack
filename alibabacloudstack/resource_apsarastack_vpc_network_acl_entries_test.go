package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNetworkAclEntries_basic(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_vpc_network_acl_entries.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvpc%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpckNetworkAclEntriesdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"network_acl_id": "${alibabacloudstack_network_acl.default.id}",
					"ingress": []map[string]interface{}{
						{
							"protocol":       "all",
							"port":           "100/100",
							"source_cidr_ip": "0.10.0.0/32",
							"name":           "${var.name}1",
							"entry_type":     "custom",
							"policy":         "accept",
							"description":    "${var.name}1",
						},
						{
							"protocol":       "tcp",
							"port":           "200/200",
							"source_cidr_ip": "0.20.0.0/32",
							"name":           "${var.name}2",
							"entry_type":     "custom",
							"policy":         "drop",
							"description":    "${var.name}2",
						},
						{
							"protocol":       "udp",
							"port":           "300/300",
							"source_cidr_ip": "0.30.0.0/32",
							"name":           "${var.name}3",
							"entry_type":     "custom",
							"policy":         "drop",
							"description":    "${var.name}3",
						},
					},
					"egress": []map[string]interface{}{
						{
							"protocol":            "all",
							"port":                "111/111",
							"destination_cidr_ip": "0.0.10.0/32",
							"name":                "${var.name}1",
							"entry_type":          "custom",
							"policy":              "drop",
							"description":         "${var.name}1",
						},
						{
							"protocol":            "tcp",
							"port":                "222/222",
							"destination_cidr_ip": "0.0.20.0/32",
							"name":                "${var.name}2",
							"entry_type":          "custom",
							"policy":              "accept",
							"description":         "${var.name}2",
						},
						{
							"protocol":            "udp",
							"port":                "333/333",
							"destination_cidr_ip": "0.0.30.0/32",
							"name":                "${var.name}3",
							"entry_type":          "custom",
							"policy":              "accept",
							"description":         "${var.name}3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_acl_id":               CHECKSET,
						"ingress.#":                    "3",
						"ingress.0.protocol":           "all",
						"ingress.0.port":               "100/100",
						"ingress.0.source_cidr_ip":     "0.10.0.0/32",
						"ingress.0.name":               name + "1",
						"ingress.0.entry_type":         "custom",
						"ingress.0.policy":             "accept",
						"ingress.0.description":        name + "1",
						"ingress.1.protocol":           "tcp",
						"ingress.1.port":               "200/200",
						"ingress.1.source_cidr_ip":     "0.20.0.0/32",
						"ingress.1.name":               name + "2",
						"ingress.1.entry_type":         "custom",
						"ingress.1.policy":             "drop",
						"ingress.1.description":        name + "2",
						"ingress.2.name":               name + "3",
						"ingress.2.description":        name + "3",
						"egress.#":                     "3",
						"egress.0.protocol":            "all",
						"egress.0.port":                "111/111",
						"egress.0.destination_cidr_ip": "0.0.10.0/32",
						"egress.0.name":                name + "1",
						"egress.0.entry_type":          "custom",
						"egress.0.policy":              "drop",
						"egress.0.description":         name + "1",
						"egress.1.name":                name + "2",
						"egress.1.description":         name + "2",
						"egress.2.name":                name + "3",
						"egress.2.description":         name + "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ingress": []map[string]interface{}{
						{
							"protocol":       "udp",
							"port":           "300/300",
							"source_cidr_ip": "0.30.0.0/32",
							"name":           "${var.name}3",
							"entry_type":     "custom",
							"policy":         "drop",
							"description":    "${var.name}3",
						},
						{
							"protocol":       "all",
							"port":           "100/100",
							"source_cidr_ip": "0.10.0.0/32",
							"name":           "${var.name}1",
							"entry_type":     "custom",
							"policy":         "accept",
							"description":    "${var.name}1",
						},
						{
							"protocol":       "tcp",
							"port":           "200/200",
							"source_cidr_ip": "0.20.0.0/32",
							"name":           "${var.name}2",
							"entry_type":     "custom",
							"policy":         "drop",
							"description":    "${var.name}2",
						},
					},
					"egress": []map[string]interface{}{
						{
							"protocol":            "udp",
							"port":                "333/333",
							"destination_cidr_ip": "0.0.30.0/32",
							"name":                "${var.name}3",
							"entry_type":          "custom",
							"policy":              "accept",
							"description":         "${var.name}3",
						},
						{
							"protocol":            "tcp",
							"port":                "222/222",
							"destination_cidr_ip": "0.0.20.0/32",
							"name":                "${var.name}2",
							"entry_type":          "custom",
							"policy":              "accept",
							"description":         "${var.name}2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ingress.#":                    "3",
						"ingress.0.protocol":           "udp",
						"ingress.0.port":               "300/300",
						"ingress.0.source_cidr_ip":     "0.30.0.0/32",
						"ingress.0.name":               name + "3",
						"ingress.0.entry_type":         "custom",
						"ingress.0.policy":             "drop",
						"ingress.0.description":        name + "3",
						"ingress.1.protocol":           "all",
						"ingress.1.port":               "100/100",
						"ingress.1.source_cidr_ip":     "0.10.0.0/32",
						"ingress.1.name":               name + "1",
						"ingress.1.entry_type":         "custom",
						"ingress.1.policy":             "accept",
						"ingress.1.description":        name + "1",
						"ingress.2.name":               name + "2",
						"ingress.2.description":        name + "2",
						"egress.#":                     "2",
						"egress.0.protocol":            "udp",
						"egress.0.port":                "333/333",
						"egress.0.destination_cidr_ip": "0.0.30.0/32",
						"egress.0.name":                name + "3",
						"egress.0.entry_type":          "custom",
						"egress.0.policy":              "accept",
						"egress.0.description":         name + "3",
						"egress.1.name":                name + "2",
						"egress.1.description":         name + "2",
						"egress.2.name":                REMOVEKEY,
						"egress.2.description":         REMOVEKEY,
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

func AlibabacloudTestAccVpckNetworkAclEntriesdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_network_acl" "default" {
  vpc_id = alibabacloudstack_vpc.default.id
  name   = var.name
}

`, name)
}
