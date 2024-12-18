package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAccountDataSource_basic(t *testing.T) {
	// 不支持 sts 跳过测试
	t.Skip()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackAccountDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_account.current"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_account.current", "id"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackAccountDataSourceBasic = `
data "alibabacloudstack_account" "current" {
}
`
