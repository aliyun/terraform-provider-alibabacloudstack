package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCloudfwAddressBook_Basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cloudfw_address_book.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testYundunProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-addrbook-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCloudfwAddressBookConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: nil,
		Providers:    testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_type":   "port",
					"group_name":   "${var.name}",
					"description":  "${var.name}",
					"address_list": []string{"8888", "9999"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_type":     "port",
						"group_name":     name,
						"description":    name,
						"address_list.#": "2",
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
					"description":  "${var.name}_updated",
					"address_list": []string{"8888"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    name + "_updated",
						"address_list.#": "1",
					}),
				),
			},
		},
	})
}

func resourceCloudfwAddressBookConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
