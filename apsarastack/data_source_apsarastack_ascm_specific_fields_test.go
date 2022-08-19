package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_Service_FieldsDataSource(t *testing.T) { //not completed
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_servicefields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_specific_fields.default"),
					//resource.TestCheckNoResourceAttr("data.apsarastack_ascm_specific_fields.default", "group_filed"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_servicefields = `

data "apsarastack_ascm_specific_fields" "default" {
  group_filed ="storageType"
  resource_type ="OSS"
}
`
