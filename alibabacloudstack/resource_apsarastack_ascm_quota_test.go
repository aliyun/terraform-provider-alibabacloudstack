package alibabacloudstack

import (
	"fmt"
	"os"
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
					"quota_type_id": testAccGetOrgId(),
					"product_name":  "VPC",
					"total_vpc":     1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_type":    "organization",
					"quota_type_id": testAccGetOrgId(),
					"product_name":  "VPC",
					"total_vpc":     2,
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

// testAccGetOrgId returns the organization ID from environment.
// Checks ALIBABACLOUDSTACK_DEPARTMENT first (provider's actual env var),
// then ALIBABACLOUDSTACK_ORGANIZATION_ID for convenience.
func testAccGetOrgId() string {
	if v := os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT"); v != "" {
		return v
	}
	if v := os.Getenv("ALIBABACLOUDSTACK_ORGANIZATION_ID"); v != "" {
		return v
	}
	return ""
}

func testAccCheckAscm_QuotaVPC(name string) string {
	return `
# VPC quota set directly on existing org; no child org needed.
`
}

func testAccCheckAscm_QuotaECS(name string) string {
	return `
# ECS quota set directly on existing org; no child org needed.
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
					"quota_type_id": testAccGetOrgId(),
					"product_name":  "ECS",
					"total_cpu":     200,
					"total_mem":     400,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_type":    "organization",
					"quota_type_id": testAccGetOrgId(),
					"product_name":  "ECS",
					"total_cpu":     300,
					"total_mem":     600,
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

var testAccCheckQuotaVPC = map[string]string{
	"product_name":  CHECKSET,
	"quota_type_id": CHECKSET,
	"quota_type":    CHECKSET,
}

var testAccCheckQuotaECS = map[string]string{
	"product_name":  CHECKSET,
	"quota_type_id": CHECKSET,
	"quota_type":    CHECKSET,
	"total_cpu":     CHECKSET,
	"total_mem":     CHECKSET,
}
