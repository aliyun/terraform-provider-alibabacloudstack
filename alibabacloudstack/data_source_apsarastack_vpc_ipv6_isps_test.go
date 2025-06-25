package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcIpv6IspsDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackVpcIpv6IspsDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_vpc_ipv6_isps.isps"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.service_provider"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.type"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.cidr_block"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.available_count"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.in_use_count"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackVpcIpv6IspsDataSourceBasicConfig = `
data "alibabacloudstack_vpc_ipv6_isps" "isps" {
}
`

func TestAccAlibabacloudStackVpcIpv6IspsServiceProviderDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackVpcIpv6IspsServiceProviderDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_vpc_ipv6_isps.isps"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.service_provider"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.type"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.cidr_block"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.available_count"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.in_use_count"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackVpcIpv6IspsServiceProviderDataSourceBasicConfig = `
data "alibabacloudstack_vpc_ipv6_isps" "isps" {
	service_provider = "BGP"
}
`

func TestAccAlibabacloudStackVpcIpv6IspsLockStatusDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackVpcIpv6IspsLockStatusDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_vpc_ipv6_isps.isps"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.service_provider"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.type"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.cidr_block"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.available_count"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_vpc_ipv6_isps.isps", "ipv6_isps.0.in_use_count"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackVpcIpv6IspsLockStatusDataSourceBasicConfig = `
data "alibabacloudstack_vpc_ipv6_isps" "isps" {
   lock_status = "unlocked"
}
`
