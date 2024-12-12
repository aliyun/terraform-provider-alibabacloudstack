package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbAccesscontrollist0(t *testing.T) {
	var v *slb.DescribeAccessControlListAttributeResponse

	resourceId := "alibabacloudstack_slb_accesscontrollist.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbAccesscontrollistCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribeaccesscontrollistattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbaccess_control_list%d", defaultRegionToTest, rand)
	entry := make([]map[string]string, 0)
	entry = append(entry, map[string]string{
		"entry":   "192.168.1.0/24",
		"comment": "test_entry",
	})
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
					"acl_name": "Rdk_test_name01",

					"address_ip_version": "ipv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"address_ip_version": "ipv4",

						"acl_name": "Rdk_test_name01",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"address_ip_version": "ipv6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"address_ip_version": "ipv6",
					}),
				),
			},

			// {
			// 	Config: testAccConfig(map[string]interface{}{

			// 		"entry_list": entry,
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{

			// 			"entry_list.%": "1",
			// 		}),
			// 	),
			// },

			// 	{
			// 		Config: testAccConfig(map[string]interface{}{
			// 			"tags": map[string]string{
			// 				"Created": "TF1",
			// 				"For":     "Test1",
			// 			},
			// 		}),
			// 		Check: resource.ComposeTestCheckFunc(
			// 			testAccCheck(map[string]string{
			// 				"tags.%":       "2",
			// 				"tags.Created": "TF1",
			// 				"tags.For":     "Test1",
			// 			}),
			// 		),
			// 	},
			// 	{
			// 		Config: testAccConfig(map[string]interface{}{
			// 			"tags": map[string]string{
			// 				"Created": "TF-update",
			// 				"For":     "Test-update",
			// 			},
			// 		}),
			// 		Check: resource.ComposeTestCheckFunc(
			// 			testAccCheck(map[string]string{
			// 				"tags.%":       "2",
			// 				"tags.Created": "TF-update",
			// 				"tags.For":     "Test-update",
			// 			}),
			// 		),
			// 	},
			// 	{
			// 		Config: testAccConfig(map[string]interface{}{
			// 			"tags": REMOVEKEY,
			// 		}),
			// 		Check: resource.ComposeTestCheckFunc(
			// 			testAccCheck(map[string]string{
			// 				"tags.%":       "0",
			// 				"tags.Created": REMOVEKEY,
			// 				"tags.For":     REMOVEKEY,
			// 			}),
			// 		),
			// 	},
		},
	})
}

var AlibabacloudTestAccSlbAccesscontrollistCheckmap = map[string]string{

	// "address_ip_version": CHECKSET,

	// "resource_group_id": CHECKSET,

	// "acl_id": CHECKSET,

	// "related_listeners": CHECKSET,

	// "acl_name": CHECKSET,

	// "tags": CHECKSET,
}

func AlibabacloudTestAccSlbAccesscontrollistBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
