package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbAccesscontrollist0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_accesscontrollist.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbAccesscontrollistCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribeaccesscontrollistattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbaccess_control_list%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbAccesscontrollistBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"address_ip_version": "ipv4",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultsRG.ResourceGroupId)}}",

					"acl_name": "Rdk_test_name01",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"address_ip_version": "ipv4",

						"resource_group_id": CHECKSET,

						"acl_name": "Rdk_test_name01",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"acl_name": "Rdk-test-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"acl_name": "Rdk-test-name",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccSlbAccesscontrollistCheckmap = map[string]string{

	"address_ip_version": CHECKSET,

	"resource_group_id": CHECKSET,

	"acl_id": CHECKSET,

	"related_listeners": CHECKSET,

	"acl_name": CHECKSET,

	"tags": CHECKSET,
}

func AlibabacloudTestAccSlbAccesscontrollistBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
