package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackDnsDomainDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackDnsDomain,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_dns_domains.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_dns_domains.default", "domains.domain_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_dns_domains.default", "domains.domain_name"),
				),
			},
		},
	})
}

const dataSourceApsaraStackDnsDomain = `

resource "apsarastack_dns_domain" "default" {
 domain_name = "testdummy."
 remark = "test_dummy"
}
data "apsarastack_dns_domains" "default"{
 domain_name_regex         = apsarastack_dns_domain.default.domain_name
}

`
