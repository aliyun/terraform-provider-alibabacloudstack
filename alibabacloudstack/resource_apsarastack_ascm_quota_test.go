package alibabacloudstack

import (
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscm_QuotaBasic(t *testing.T) {
	var v *AscmQuota
	resourceId := "alibabacloudstack_ascm_quota.default"
	ra := resourceAttrInit(resourceId, testAccCheckQuota)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckAscm_Quota_Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAscm_Quota,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testAccCheckAscm_Quota_Destroy(s *terraform.State) error { //destroy function
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ascm_quota" || rs.Type != "alibabacloudstack_ascm_quota" {
			continue
		}
		ascm, err := ascmService.DescribeAscmQuota(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if ascm.Message != "" {
			return errmsgs.WrapError(errmsgs.Error("resource  still exist"))
		}
	}

	return nil
}

const testAccCheckAscm_Quota = `
resource "alibabacloudstack_ascm_organization" "default" {
 name = "TestQuota"
 parent_id = "1"
}

resource "alibabacloudstack_ascm_quota" "default" {
  quota_type = "organization"
  quota_type_id = alibabacloudstack_ascm_organization.default.org_id
    product_name = "ECS"
	total_cpu = 100
    total_mem = 100
    total_gpu = 100
    total_disk_cloud_ssd = 100
    total_disk_cloud_efficiency = 100
}
`

var testAccCheckQuota = map[string]string{
	"product_name":  CHECKSET,
	"quota_type_id": CHECKSET,
	"quota_type":    CHECKSET,
}
