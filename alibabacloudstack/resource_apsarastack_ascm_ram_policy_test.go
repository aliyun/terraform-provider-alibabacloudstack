package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscmRamPolicyBasic(t *testing.T) {
	var v *RamPolicies

	resourceId := "alibabacloudstack_ascm_ram_policy.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmRamPolicy)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmrampolicy%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccAscm_e_Organization_resource)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_RamPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            "${var.name}",
					"description":     "Testing Policy",
					"policy_document": `{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            "${var.name}Update",
					"description":     "Testing Policy2",
					"policy_document": `{\"Statement\":[{\"Action\":\"vpc:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}`,
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

func testAccCheckAscm_RamPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)

	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ascm_ram_policy" || rs.Type != "alibabacloudstack_ascm_ram_policy" {
			continue
		}
		ascm, err := ascmService.DescribeAscmRamPolicy(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if ascm.AsapiErrorCode != "200" {
			return errmsgs.WrapError(errmsgs.Error("ram role still exist"))
		}
	}

	return nil
}

func testAccAscm_RamPolicy_resource(name string) string {
	return fmt.Sprintf(`
variables "name" {
	default = "%s"
}
`, name)
}

var testAccCheckAscmRamPolicy = map[string]string{
	"name":            CHECKSET,
	"policy_document": CHECKSET,
}
