package apsarastack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackAccountDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackAccountDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_account.current"),
					resource.TestCheckResourceAttrSet("data.apsarastack_account.current", "id"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackAccountDataSourceBasic = `
data "apsarastack_account" "current" {
}
`
