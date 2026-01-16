package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackSecurityGroupRule_sg(t *testing.T) {
	var v ecs.Permission
	resourceId := "alibabacloudstack_security_group_rule.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupRuleBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc_sg_rule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccSecurityGroupRuleBasic)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":                     "ingress",
					"ip_protocol":              "tcp",
					"policy":                   "drop",
					"port_range":               "22/22",
					"priority":                 100,
					"security_group_id":        "${alibabacloudstack_security_group.default.0.id}",
					"source_security_group_id": "${alibabacloudstack_security_group.default.1.id}",
					"description":              "abc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abc",
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
					"source_security_group_id": REMOVEKEY,
					"cidr_ip":                  "0.0.0.0/0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_security_group_id": REMOVEKEY,
						"cidr_ip":                  "0.0.0.0/0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "egress",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "egress",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackSecurityGroupRule_cidr(t *testing.T) {
	var v ecs.Permission
	resourceId := "alibabacloudstack_security_group_rule.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":        "ingress",
		"policy":      "accept",
		"description": "SHDRP-7513",
		"port_range":  "443/443",
		"priority":    "1",
		"cidr_ip":     "182.254.11.243/32",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc_sg_rule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccSecurityGroupRuleBasic)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":              "ingress",
					"ip_protocol":       "tcp",
					"policy":            "accept",
					"port_range":        "443/443",
					"priority":          "1",
					"security_group_id": "${alibabacloudstack_security_group.default.0.id}",
					"cidr_ip":           "182.254.11.243/32",
					"description":       "SHDRP-7513",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "egress",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "egress",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackSecurityGroupRule_cidrv6(t *testing.T) {
	var v ecs.Permission
	resourceId := "alibabacloudstack_security_group_rule.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"policy":      "accept",
		"description": "SHDRP-7513",
		"port_range":  "443/443",
		"priority":    "1",
		"ipv6_cidr_ip":     "2001:0DB8::1428:57ab/128",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc_sg_rule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccSecurityGroupRuleBasic)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":              "egress",
					"ip_protocol":       "tcp",
					"policy":            "accept",
					"port_range":        "443/443",
					"priority":          "1",
					"security_group_id": "${alibabacloudstack_security_group.default.0.id}",
					"ipv6_cidr_ip":      "2001:0DB8::1428:57ab/128",
					"description":       "SHDRP-7513",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "egress",
						"description": "SHDRP-7513",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "ingress",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "ingress",
					}),
				),
			},
		},
	})

}

func testAccSecurityGroupRuleBasic(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	resource "alibabacloudstack_vpc" "default" {
	  name = "${var.name}"
	  cidr_block = "192.168.0.0/16"
	}

	resource "alibabacloudstack_security_group" "default" {
	  count  = 2
	  vpc_id = "${alibabacloudstack_vpc.default.id}"
	  name = "${var.name}"
	}`, name)
}

var testAccCheckSecurityGroupRuleBasicMap = map[string]string{
	"type":                     "ingress",
	"ip_protocol":              "tcp",
	"nic_type":                 "intranet",
	"policy":                   "drop",
	"port_range":               "22/22",
	"priority":                 "100",
	"security_group_id":        CHECKSET,
	"source_security_group_id": CHECKSET,
	"cidr_ip":                  "",
}

func testAccCheckSecurityGroupRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_security_group_rule" {
			continue
		}
		_, err := ecsService.DescribeSecurityGroupRule(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil && !errmsgs.NotFoundError(err) {
			return errmsgs.WrapError(err)
		}
	}

	return nil
}
