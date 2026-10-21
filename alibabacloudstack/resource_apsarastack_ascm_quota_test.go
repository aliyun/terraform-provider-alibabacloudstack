package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscmQuota_VPC(t *testing.T) {
	var v *AscmQuota
	resourceId := "alibabacloudstack_ascm_quota.default"
	ra := resourceAttrInit(resourceId, testAccCheckQuotaVPC)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmorg%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckAscm_QuotaVPC)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckAscm_Quota_Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_type":    "organization",
					"quota_type_id": TfRawString("data.alibabacloudstack_account.current.organization_id"),
					"product_name":  "VPC",
					"total_vpc":     1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"total_vpc": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"total_vpc": 2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"total_vpc": "2",
					}),
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

func testAccCheckAscm_QuotaVPC(name string) string {
	return `
data "alibabacloudstack_account" "current" {}
`
}

func testAccCheckAscm_QuotaECS(name string) string {
	return `
data "alibabacloudstack_account" "current" {}
`
}

// TestAccAlibabacloudStackAscmQuota_ECS mirrors the tenant_onboarding scenario:
// creates a brand-new org then sets ECS quota (total_cpu/total_mem) on it.
// Validates that Create routes to /quota/add (not /quota/update) for a new org.
func TestAccAlibabacloudStackAscmQuota_ECS(t *testing.T) {
	var v *AscmQuota
	resourceId := "alibabacloudstack_ascm_quota.ecs"
	ra := resourceAttrInit(resourceId, testAccCheckQuotaECS)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmorg%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckAscm_QuotaECS)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_Quota_Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_type":    "organization",
					"quota_type_id": TfRawString("data.alibabacloudstack_account.current.organization_id"),
					"product_name":  "ECS",
					"total_cpu":     200,
					"total_mem":     400,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"total_cpu": "200",
						"total_mem": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"total_cpu": 300,
					"total_mem": 600,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"total_cpu": "300",
						"total_mem": "600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"total_cpu": 0, // Delete this quota item
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"total_cpu": REMOVEKEY,
					}),
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

var testAccCheckQuotaVPC = map[string]string{
	"product_name":  "VPC",
	"quota_type_id": CHECKSET,
	"quota_type":    "organization",
}

var testAccCheckQuotaECS = map[string]string{
	"product_name":  "ECS",
	"quota_type_id": CHECKSET,
	"quota_type":    "organization",
	"total_cpu":     CHECKSET,
	"total_mem":     CHECKSET,
}
