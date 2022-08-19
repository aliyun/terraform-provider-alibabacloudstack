package apsarastack

import (
	"fmt"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"testing"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackEcsCommand_basic(t *testing.T) {
	var v *datahub.EcsDescribeEcsCommandResult
	resourceId := "apsarastack_ecs_command.default"
	ra := resourceAttrInit(resourceId, ApsaraStackEcsCommandMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeEcsCommand")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sApsaraStackEcsCommand%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackEcsCommandBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"command_content": "bHMK",
					"description":     "For Terraform Test",
					"name":            name,
					"type":            "RunShellScript",
					"working_dir":     "/root",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_content": "bHMK",
						"description":     "For Terraform Test",
						"name":            name,
						"type":            "RunShellScript",
						"working_dir":     "/root",
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

var ApsaraStackEcsCommandMap = map[string]string{
	"enable_parameter": "false",
}

func ApsaraStackEcsCommandBasicDependence(name string) string {
	return fmt.Sprintf(`
provider "apsarastack" {
	assume_role {}
}
`)
}
