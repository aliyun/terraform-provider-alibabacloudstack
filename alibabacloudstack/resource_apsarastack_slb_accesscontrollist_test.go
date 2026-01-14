package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbAcl0(t *testing.T) {
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
	entry_ipv6 := make([]map[string]string, 0)
	entry_ipv6 = append(entry_ipv6, map[string]string{
		"entry":   "2001:db8::/32",
		"comment": "test_entry_ipv6",
	})

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbAccesscontrollistBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name":           name,
					"address_ip_version": "ipv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_ip_version": "ipv4",
						"acl_name":           name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name": name + "_update",
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
					"entry_list": []map[string]string{{
						"entry":   "192.168.1.0/24",
						"comment": "test_entry1",
					},{
						"entry":   "192.168.2.0/24",
						"comment": "test_entry2",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entry_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entry_list": []map[string]string{{
						"entry":   "192.168.1.0/24",
						"comment": "test_entry1",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entry_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_ip_version": "ipv6",
					"entry_list":         entry_ipv6,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_ip_version": "ipv6",
						"entry_list.#":       "1",
					}),
				),
			},
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
