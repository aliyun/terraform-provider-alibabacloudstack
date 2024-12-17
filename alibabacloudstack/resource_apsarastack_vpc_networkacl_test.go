package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcNetworkacl0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_vpc_networkacl.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcNetworkaclCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribenetworkaclattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcnetwork_acl%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcNetworkaclBasicdependence)
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

					"description": "test",

					"vpc_id": "vpc-bp11lfjeaa57jxr6ovybf",

					"network_acl_name": "rdk-test",

					"region_id": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"vpc_id": "vpc-bp11lfjeaa57jxr6ovybf",

						"network_acl_name": "rdk-test",

						"region_id": "cn-hangzhou",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test-update",

					"network_acl_name": "rdk-test-name",

					"region_id": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test-update",

						"network_acl_name": "rdk-test-name",

						"region_id": "cn-hangzhou",
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

var AlibabacloudTestAccVpcNetworkaclCheckmap = map[string]string{

	"ingress_acl_entries": CHECKSET,

	"status": CHECKSET,

	"description": CHECKSET,

	"network_acl_id": CHECKSET,

	"create_time": CHECKSET,

	"vpc_id": CHECKSET,

	"egress_acl_entries": CHECKSET,

	"network_acl_name": CHECKSET,

	"region_id": CHECKSET,

	"resources": CHECKSET,

	"tags": CHECKSET,
}

func AlibabacloudTestAccVpcNetworkaclBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
