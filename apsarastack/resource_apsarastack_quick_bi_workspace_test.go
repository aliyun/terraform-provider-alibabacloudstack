package apsarastack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlicloudQuickBIWorkspace_basic0(t *testing.T) {
	//t.Skip()
	var v map[string]interface{}
	resourceId := "apsarastack_quick_bi_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudQuickBIWorkspaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuickbiPublicService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeQuickBiWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squickbiWorkspace%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuickBIWorkspaceBasicDependence0)
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
					"workspace_name": name,
					"workspace_desc": "desc-" + name,
					"use_comment":    "false",
					"allow_share":    "false",
					"allow_publish":  "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace_name": name,
						"workspace_desc": "desc-" + name,
						"use_comment":    "false",
						"allow_share":    "false",
						"allow_publish":  "false",
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

var AlicloudQuickBIWorkspaceMap0 = map[string]string{}

func AlicloudQuickBIWorkspaceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
