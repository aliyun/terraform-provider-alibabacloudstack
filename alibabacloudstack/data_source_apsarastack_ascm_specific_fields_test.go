package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_Service_FieldsDataSource(t *testing.T) { //not completed
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_servicefields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_specific_fields.default"),
					//resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_specific_fields.default", "group_filed"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_servicefields = `

data "alibabacloudstack_ascm_specific_fields" "default" {
  group_filed ="storageType"
  resource_type ="OSS"
}
`
