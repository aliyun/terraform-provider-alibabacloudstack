package alibabacloudstack

import (
	"fmt"
	"os"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCSBProject_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_csb_project.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackCSBProjectMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CsbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCsbProjectDetail")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-csbproject%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackCSBProjectBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_TEST_EXISTED_CSB_ID")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"csb_id":          "${var.csb_id}",
					"project_name":    "${var.name}",
					"owner_name":      "${var.name}_user",
					"owner_email":     "${var.name}@aliyun.test",
					"owner_phone_num": "13900000000",
					"description":     "${var.name} desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":    name,
						"owner_name":      name + "_user",
						"owner_email":     name + "@aliyun.test",
						"owner_phone_num": "13900000000",
						"description":     name + " desc",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project_name":    "${var.name}_update",
					"owner_name":      "${var.name}_user1",
					"owner_email":     "${var.name}_update@aliyun.test",
					"owner_phone_num": "15000000000",
					"description":     "${var.name} desc update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":    name + "_update",
						"owner_name":      name + "_user1",
						"owner_email":     name + "_update@aliyun.test",
						"owner_phone_num": "15000000000",
						"description":     name + " desc update",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackCSBProjectMap0 = map[string]string{
	"csb_id":       CHECKSET,
	"project_name": CHECKSET,
}

func AlibabacloudStackCSBProjectBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable name {
		default = "%s"
	}
	
	variable csb_id {
		default = "%s"
	}
	
	`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_CSB_ID"))
}
