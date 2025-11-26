package alibabacloudstack

import (
	"testing"

	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackUniversalDnsDomain_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_universal_dns_domain.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbVservergroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &UniversalDnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeUniversalZones")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccUniversalDnsDomaindependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{

					"name":   fmt.Sprintf("%s.example.", name),
					"remark": "Created by Terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   fmt.Sprintf("%s.example.", name),
						"remark": "Created by Terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "Updated by Terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "Updated by Terraform",
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

func AlibabacloudTestAccUniversalDnsDomaindependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}
