package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDnsPrivateDomain_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dns_private_domain.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"name":             CHECKSET,
		"caller_uid":       CHECKSET,
		"record_count":     "0",
		"create_timestamp": CHECKSET,
		"update_timestamp": CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribePrivateZone")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfacc%d.test.", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudbuildBasicVpcTemplate)

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
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark":  "${var.name}",
					"vpc_ids": []string{"${alibabacloudstack_vpc_vpc.default0.id}", "${alibabacloudstack_vpc_vpc.default1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark":    name,
						"vpc_ids.#": "2",
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

func AlibabacloudbuildBasicVpcTemplate(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_vpc_vpc" "default0" {
	cidr_block = "172.16.0.0/12"
	vpc_name   = "${var.name}_vpc0"
}

resource "alibabacloudstack_vpc_vpc" "default1" {
	cidr_block = "192.168.0.0/16"
	vpc_name   = "${var.name}_vpc1"
}
`, name)
}
