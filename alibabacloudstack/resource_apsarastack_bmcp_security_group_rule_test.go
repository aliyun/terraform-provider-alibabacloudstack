package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackBmcpSecurityGroupRule_sg_ingress(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_bcmp_security_group_rule.default"
	ra := resourceAttrInit(resourceId, testAccCheckBmcpSecurityGroupRuleBasicMap)
	serviceFunc := func() interface{} {
		return &BcmpService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoEasyAIListSecurityGroupRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc_bmcp_sg_rule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccBmcpSecurityGroupRuleBasic)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckBmcpSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":                     "ingress",
					"ip_protocol":              "tcp",
					"policy":                   "drop",
					"port_range":               "22/22",
					"priority":                 100,
					"security_group_id":        "${alibabacloudstack_bcmp_security_group.default.id}",
					"source_security_group_id": "${alibabacloudstack_bcmp_security_group.new.id}",
					"description":              "abc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":              "abc",
						"type":                     "ingress",
						"policy":                   "drop",
						"port_range":               "22/22",
						"priority":                 "100",
						"security_group_id":        CHECKSET,
						"source_security_group_id": CHECKSET,
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
					"ip_protocol":              "udp",
					"policy":                   "accept",
					"port_range":               "33/33",
					"priority":                 30,
					"description":              "abc_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_security_group_id": REMOVEKEY,
						"cidr_ip":                  "0.0.0.0/0",
						"ip_protocol":              "udp",
						"policy":                   "accept",
						"port_range":               "33/33",
						"priority":                 "30",
						"description":              "abc_update",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackBmcpSecurityGroupRule_sg_egress(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_bcmp_security_group_rule.default"
	ra := resourceAttrInit(resourceId, testAccCheckBmcpSecurityGroupRuleBasicMap)
	serviceFunc := func() interface{} {
		return &BcmpService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoEasyAIListSecurityGroupRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc_bmcp_sg_rule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccBmcpSecurityGroupRuleBasic)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckBmcpSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":                     "egress",
					"ip_protocol":              "tcp",
					"policy":                   "drop",
					"port_range":               "22/22",
					"priority":                 100,
					"security_group_id":        "${alibabacloudstack_bcmp_security_group.default.id}",
					"source_security_group_id": "${alibabacloudstack_bcmp_security_group.new.id}",
					"description":              "abc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":              "abc",
						"type":                     "egress",
						"policy":                   "drop",
						"port_range":               "22/22",
						"priority":                 "100",
						"security_group_id":        CHECKSET,
						"source_security_group_id": CHECKSET,
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
					"ip_protocol":              "udp",
					"policy":                   "accept",
					"port_range":               "33/33",
					"priority":                 30,
					"description":              "abc_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_security_group_id": REMOVEKEY,
						"cidr_ip":                  "0.0.0.0/0",
						"ip_protocol":              "udp",
						"policy":                   "accept",
						"port_range":               "33/33",
						"priority":                 "30",
						"description":              "abc_update",
					}),
				),
			},
		},
	})
}

func testAccBmcpSecurityGroupRuleBasic(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	resource "alibabacloudstack_vpc" "default" {
		name = "${var.name}"
		cidr_block = "172.16.0.0/12"
	}

	resource "alibabacloudstack_bcmp_security_group" "default" {
	  vpc_id = alibabacloudstack_vpc.default.id
	  name = "${var.name}"
	}
	resource "alibabacloudstack_bcmp_security_group" "new" {
	  vpc_id = alibabacloudstack_vpc.default.id
	  name = "${var.name}_new"
	}
	  
	`, name)
}

var testAccCheckBmcpSecurityGroupRuleBasicMap = map[string]string{
	"security_group_id":        CHECKSET,
	"source_security_group_id": CHECKSET,
}

func testAccCheckBmcpSecurityGroupRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	bmcpService := BcmpService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_bcmp_security_group_rule" {
			continue
		}
		_, err := bmcpService.DoEasyAIListSecurityGroupRequest(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil && !errmsgs.NotFoundError(err) {
			return errmsgs.WrapError(err)
		}
	}

	return nil
}
