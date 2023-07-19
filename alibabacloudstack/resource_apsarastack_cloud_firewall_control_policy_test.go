package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackCloudFirewallControlPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudStack_cloud_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name": "ANY",
					"acl_action":       "accept",
					"description":      "放行所有流量",
					"destination_type": "net",
					"destination":      "100.1.1.0/24",
					"direction":        "out",
					"proto":            "ANY",
					"source":           "1.2.3.0/24",
					"source_type":      "net",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": "ANY",
						"acl_action":       "accept",
						"description":      "放行所有流量",
						"destination_type": "net",
						"destination":      "100.1.1.0/24",
						"direction":        "out",
						"proto":            "ANY",
						"source":           "1.2.3.0/24",
						"source_type":      "net",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"release": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"release": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action": "log",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action": "log",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "Any",
					"destination_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "Any",
						"destination_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "www.aliyun.com",
					"destination_type": "domain",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "www.aliyun.com",
						"destination_type": "domain",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source":      "Any",
					"source_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source":      "Any",
						"source_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port_group": "ANY",
					"dest_port_type":  "group",
					"dest_port":       "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port_group": "ANY",
						"dest_port_type":  "group",
						"dest_port":       "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":       "100/200",
					"dest_port_type":  "port",
					"dest_port_group": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":       "100/200",
						"dest_port_type":  "port",
						"dest_port_group": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proto": "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proto": "TCP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "放行TCP流量",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "放行TCP流量",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "100.1.1.0/24",
					"application_name": "ANY",
					"description":      "放行所有流量",
					"source_type":      "net",
					"acl_action":       "accept",
					"destination_type": "net",
					"direction":        "out",
					"source":           "1.2.3.0/24",
					"dest_port_type":   "port",
					"proto":            "ANY",
					"dest_port":        "0/0",
					"release":          "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "100.1.1.0/24",
						"application_name": "ANY",
						"description":      "放行所有流量",
						"source_type":      "net",
						"acl_action":       "accept",
						"destination_type": "net",
						"direction":        "out",
						"source":           "1.2.3.0/24",
						"dest_port_type":   "port",
						"proto":            "ANY",
						"dest_port":        "0/0",
						"release":          "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"ip_version", "source_ip"},
			},
		},
	})
}

func TestAccAlibabacloudStackCloudFirewallControlPolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudStack_cloud_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackCloudFirewallControlPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackCloudFirewallControlPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name": "ANY",
					"acl_action":       "accept",
					"description":      "放行所有流量",
					"destination_type": "net",
					"destination":      "100.1.1.0/24",
					"direction":        "in",
					"proto":            "ANY",
					"source":           "1.2.3.0/24",
					"source_type":      "net",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": "ANY",
						"acl_action":       "accept",
						"description":      "放行所有流量",
						"destination_type": "net",
						"destination":      "100.1.1.0/24",
						"direction":        "in",
						"proto":            "ANY",
						"source":           "1.2.3.0/24",
						"source_type":      "net",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"release": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"release": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_action": "log",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_action": "log",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "Any",
					"destination_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "Any",
						"destination_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source":      "Any",
					"source_type": "group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source":      "Any",
						"source_type": "group",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port":      "100/200",
					"dest_port_type": "port",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dest_port":      "100/200",
						"dest_port_type": "port",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proto": "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proto": "TCP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "放行TCP流量",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "放行TCP流量",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "100.1.1.0/24",
					"application_name": "ANY",
					"description":      "放行所有流量",
					"source_type":      "net",
					"acl_action":       "accept",
					"destination_type": "net",
					"direction":        "out",
					"source":           "1.2.3.0/24",
					"dest_port_type":   "port",
					"proto":            "ANY",
					"dest_port":        "0/0",
					"release":          "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      "100.1.1.0/24",
						"application_name": "ANY",
						"description":      "放行所有流量",
						"source_type":      "net",
						"acl_action":       "accept",
						"destination_type": "net",
						"direction":        "out",
						"source":           "1.2.3.0/24",
						"dest_port_type":   "port",
						"proto":            "ANY",
						"dest_port":        "0/0",
						"release":          "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"ip_version", "source_ip"},
			},
		},
	})
}

var AlibabacloudStackCloudFirewallControlPolicyMap0 = map[string]string{

	"release":          "true",
	"source_ip":        NOSET,
	"application_name": "ANY",
	"description":      "放行所有流量",
	"destination":      "100.1.1.0/24",
	"proto":            "ANY",
	"source":           "1.2.3.0/24",
	"acl_action":       "accept",
	"destination_type": "net",
	"source_type":      "net",
}

func AlibabacloudStackCloudFirewallControlPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
