package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcDhcpOptionsSet_basic(t *testing.T) {
	var v *VpcGetdhcpoptionssetResponse
	resourceId := "alibabacloudstack_vpc_dhcp_options_set.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcGetdhcpoptionssetRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccVpcDhcpoptionssetBasic")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVpcDhcpoptionssetBasicDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dhcp_options_set_name":        "${var.name}",
					"dhcp_options_set_description": "${var.name}",
					"domain_name":                  "aliyun.com",
					"domain_name_servers":          "10.82.0.12,10.82.0.13",
					"associate_vpcs": []string{
						"${alibabacloudstack_vpc.default0.id}",
						"${alibabacloudstack_vpc.default1.id}",
						"${alibabacloudstack_vpc.default2.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_name":        name,
						"dhcp_options_set_description": name,
						"domain_name":                  "aliyun.com",
						"domain_name_servers":          "10.82.0.12,10.82.0.13",
						"associate_vpcs.#":             "3",
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
					"dhcp_options_set_name":        "${var.name}_change",
					"dhcp_options_set_description": "${var.name}_change",
					"domain_name":                  "aliyun1.com",
					"domain_name_servers":          "10.82.0.12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_name":        fmt.Sprintf("%s_change", name),
						"dhcp_options_set_description": fmt.Sprintf("%s_change", name),
						"domain_name":                  "aliyun1.com",
						"domain_name_servers":          "10.82.0.12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"associate_vpcs": []string{
						"${alibabacloudstack_vpc.default0.id}",
						"${alibabacloudstack_vpc.default1.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"associate_vpcs.#": "2",
						"associate_vpcs.2": REMOVEKEY,
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF",
			// 			"For":     "Test",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "2",
			// 			"tags.Created": "TF",
			// 			"tags.For":     "Test",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF-update",
			// 			"For":     "Test-update",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "2",
			// 			"tags.Created": "TF-update",
			// 			"tags.For":     "Test-update",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": REMOVEKEY,
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "0",
			// 			"tags.Created": REMOVEKEY,
			// 			"tags.For":     REMOVEKEY,
			// 		}),
			// 	),
			// },
		},
	})
}

func resourceVpcDhcpoptionssetBasicDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

resource "alibabacloudstack_vpc" "default0" {

  name = "${var.name}0"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc" "default1" {

  name = "${var.name}1"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc" "default2" {

  name = "${var.name}2"
  cidr_block = "172.16.0.0/16"
}

`, name)
}
