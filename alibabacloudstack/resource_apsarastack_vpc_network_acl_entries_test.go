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
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVpcCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvpc%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpckNetworkAclEntriesdependence)
	resource.Test(t, resource.TestCase{
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
							"port":           "-1/-1",
							"source_cidr_ip": "0.0.0.0/32",
							"name":           "${var.name}",
							"entry_type":     "custom",
							"policy":         "accept",
							"description":    "${var.name}",
						},
					},
					"egress": []map[string]interface{}{
						{
							"protocol":            "all",
							"port":                "-1/-1",
							"destination_cidr_ip": "0.0.0.0/32",
							"name":                "${var.name}",
							"entry_type":          "custom",
							"policy":              "accept",
							"description":         "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_acl_id":           CHECKSET,
						"ingress.#":                "1",
						"ingress.0.protocol":       "all",
						"ingress.0.port":           "-1/-1",
						"ingress.0.source_cidr_ip": "0.0.0.0/32",
						"ingress.0.name":           name,
						"ingress.0.entry_type":     "custom",
						"ingress.0.policy":         "accept",
						"ingress.0.description":    name,
						"egress.#":                 "1",
						"egress.0.protocol":        "all",
						"egress.0.port":            "-1/-1",
						"egress.0.source_cidr_ip":  "0.0.0.0/32",
						"egress.0.name":            name,
						"egress.0.entry_type":      "custom",
						"egress.0.policy":          "accept",
						"egress.0.description":     name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ingress": []map[string]interface{}{
						{
							"protocol":       "TCP",
							"port":           "-1/-1",
							"source_cidr_ip": "192.168.0.0/32",
							"name":           "${var.name}_updated",
							"entry_type":     "system",
							"policy":         "drop",
							"description":    "${var.name}_updated",
						},
					},
					"egress": []map[string]interface{}{
						{
							"protocol":            "TCP",
							"port":                "30/30",
							"destination_cidr_ip": "192.168.10.0/32",
							"name":                "${var.name}_updated",
							"entry_type":          "system",
							"policy":              "drop",
							"description":         "${var.name}_updated",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ingress.#":                "1",
						"ingress.0.protocol":       "TCP",
						"ingress.0.port":           "30/30",
						"ingress.0.source_cidr_ip": "192.168.0.0/32",
						"ingress.0.name":           name + "_updated",
						"ingress.0.entry_type":     "system",
						"ingress.0.policy":         "drop",
						"ingress.0.description":    name + "_updated",
						"egress.#":                 "1",
						"egress.0.protocol":        "TCP",
						"egress.0.port":            "30/30",
						"egress.0.source_cidr_ip":  "192.168.10.0/32",
						"egress.0.name":            name + "_updated",
						"egress.0.entry_type":      "system",
						"egress.0.policy":          "drop",
						"egress.0.description":     name + "_updated",
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

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
  name              = var.name
}

resource "alibabacloudstack_network_acl_attachment" "default" {
  network_acl_id = alibabacloudstack_network_acl.default.id
  resources {
    resource_id   = alibabacloudstack_vswitch.default.id
    resource_type = "VSwitch"
  }
}
`, name)
}
