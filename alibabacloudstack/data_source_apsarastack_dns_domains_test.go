package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackDnsDomainDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackDnsDomain,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_dns_domains.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_dns_domains.default", "domains.domain_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_dns_domains.default", "domains.domain_name"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackDnsDomain = `

resource "alibabacloudstack_dns_domain" "default" {
 domain_name = "testdummy."
 remark = "test_dummy"
}
data "alibabacloudstack_dns_domains" "default"{
 domain_name_regex         = alibabacloudstack_dns_domain.default.domain_name
}

`
