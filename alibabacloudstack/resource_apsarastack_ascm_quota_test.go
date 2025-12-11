package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscmQuota_Basic(t *testing.T) {
	var v *AscmQuota
	resourceId := "alibabacloudstack_ascm_quota.default"
	ra := resourceAttrInit(resourceId, testAccCheckQuota)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmorg%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckAscm_Quota)
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
				Config: testAccConfig(map[string]interface{}{
					"quota_type":                  "organization",
					"quota_type_id":               "${alibabacloudstack_ascm_organization.default.id}",
					"product_name":                "VPC",
					"total_vpc":                   10,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_type":                  "organization",
					"quota_type_id":               "${alibabacloudstack_ascm_organization.default.id}",
					"product_name":                "VPC",
					"total_vpc":                   20,
				}),
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

func testAccCheckAscm_Quota(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_ascm_organization" "default" {
 name = "%s"
 parent_id = "1"
}`, name)
}

var testAccCheckQuota = map[string]string{
	"product_name":  CHECKSET,
	"quota_type_id": CHECKSET,
	"quota_type":    CHECKSET,
}
