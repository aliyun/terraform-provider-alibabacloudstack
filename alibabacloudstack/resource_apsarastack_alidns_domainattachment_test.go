package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDnsDomainAttachment_basic(t *testing.T) {
	testAccPreCheckWithAPIIsNotSupport(t)
	var v alidns.DescribeInstanceDomainsResponse

	resourceId := "alibabacloudstack_dns_domain_attachment.default"
	ra := resourceAttrInit(resourceId, dnsDomainAttachmnetMap)

	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDnsDomainAttachmentConfigDependence)

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
					"instance_id":  "${alibabacloudstack_dns_instance.default.id}",
					"domain_names": []string{"${alibabacloudstack_dns.default.name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"domain_names.#": "1",
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
					"domain_names": []string{"${alibabacloudstack_dns.default.name}", "${alibabacloudstack_dns.default1.name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_names.#": "2",
					}),
				),
			},
		},
	})
}

func resourceDnsDomainAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	resource "alibabacloudstack_dns_instance" "default" {
 	  dns_security    = "basic"
 	  domain_numbers  = 3
 	  version_code    = "version_personal"
 	  period          = 1
	  renewal_status  = "ManualRenewal"
	}

	resource "alibabacloudstack_dns" "default" {
  	  name = "%s.abc"
	}

	resource "alibabacloudstack_dns" "default1" {
  	  name = "%s1.abc"
	}
`, name, name)
}

var dnsDomainAttachmnetMap = map[string]string{}
