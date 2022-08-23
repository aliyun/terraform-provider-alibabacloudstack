package apsarastack

import (
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccApsaraStackAscm_RamPolicyBasic(t *testing.T) {
	var v *RamPolicies

	resourceId := "apsarastack_ascm_ram_policy.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmRamPolicy)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_RamPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAscm_RamPolicy_resource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscm_RamPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)

	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ascm_ram_policy" || rs.Type != "apsarastack_ascm_ram_policy" {
			continue
		}
		ascm, err := ascmService.DescribeAscmRamPolicy(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if ascm.AsapiErrorCode != "200" {
			return WrapError(Error("ram role still exist"))
		}
	}

	return nil
}

const testAccAscm_RamPolicy_resource = `
resource "apsarastack_ascm_ram_policy" "default" {
  name = "TestingRamPolicy"
  description = "Testing Policy"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"
}
`

var testAccCheckAscmRamPolicy = map[string]string{
	"name":            CHECKSET,
	"policy_document": CHECKSET,
}
