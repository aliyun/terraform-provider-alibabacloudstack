package alibabacloudstack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackDnsDomainDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackDnsDomain(),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_dns_domains.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_dns_domains.default", "domains.domain_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_dns_domains.default", "domains.domain_name"),
				),
			},
		},
	})
}

func dataSourceAlibabacloudStackDnsDomain () string {
	return  fmt.Sprintf(`

resource "alibabacloudstack_dns_domain" "default" {
 domain_name = "testdummy%d."
}
data "alibabacloudstack_dns_domains" "default"{
 domain_name  = alibabacloudstack_dns_domain.default.domain_name
}`, getAccTestRandInt(10000, 99999))

}
