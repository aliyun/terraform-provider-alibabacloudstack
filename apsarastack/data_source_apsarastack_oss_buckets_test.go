package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackOssBucketsDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceOssBucketsConfigDependence_basic,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_oss_buckets.default"),
					resource.TestCheckResourceAttr("data.apsarastack_oss_buckets.default", "buckets.#", "0"),
					resource.TestCheckNoResourceAttr("data.apsarastack_oss_buckets.default", "buckets.0.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_oss_buckets.default", "buckets.0.acl"),
					resource.TestCheckNoResourceAttr("data.apsarastack_oss_buckets.default", "buckets.0.extranet_endpoint"),
					resource.TestCheckNoResourceAttr("data.apsarastack_oss_buckets.default", "buckets.0.intranet_endpoint"),
					resource.TestCheckNoResourceAttr("data.apsarastack_oss_buckets.default", "buckets.0.location"),
					resource.TestCheckNoResourceAttr("data.apsarastack_oss_buckets.default", "buckets.0.owner"),
					resource.TestCheckNoResourceAttr("data.apsarastack_oss_buckets.default", "buckets.0.storage_class"),
					resource.TestCheckNoResourceAttr("data.apsarastack_oss_buckets.default", "buckets.0.creation_date"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

const dataSourceOssBucketsConfigDependence_basic = `

//resource "apsarastack_oss_bucket" "demo" {
//  bucket = "your-buset-nalfme"
//  acl    = "public-read"
//}

data "apsarastack_oss_buckets" "default" {
  name_regex = "your-buset-nalfme"
}
`
