package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcIpv6Gateway0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_vpc_ipv6gateway.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcIpv6GatewayCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribeipv6GatewayattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipv6_gateway%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcIpv6GatewayBasicdependence)
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

					"description": "test",

					"ipv6_gateway_name": "rdk-test",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

					"region_id": "cn-beijing",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultRg.ResourceGroupId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"ipv6_gateway_name": "rdk-test",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

						"region_id": "cn-beijing",

						"resource_group_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"region_id": "cn-beijing",

					"description": "test-update",

					"ipv6_gateway_name": "rdk-test-name",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultRg.ResourceGroupId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"region_id": "cn-beijing",

						"description": "test-update",

						"ipv6_gateway_name": "rdk-test-name",

						"resource_group_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"ipv6_gateway_name": "rdk-test",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

					"region_id": "cn-beijing",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultRg.ResourceGroupId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"ipv6_gateway_name": "rdk-test",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

						"region_id": "cn-beijing",

						"resource_group_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"region_id": "cn-beijing",

					"description": "test-update",

					"ipv6_gateway_name": "rdk-test-name",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::defaultRg.ResourceGroupId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"region_id": "cn-beijing",

						"description": "test-update",

						"ipv6_gateway_name": "rdk-test-name",

						"resource_group_id": CHECKSET,
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

var AlibabacloudTestAccVpcIpv6GatewayCheckmap = map[string]string{

	"status": CHECKSET,

	"description": CHECKSET,

	"ipv6_gateway_name": CHECKSET,

	"resource_group_id": CHECKSET,

	"instance_charge_type": CHECKSET,

	"create_time": CHECKSET,

	"ipv6_gateway_id": CHECKSET,

	"business_status": CHECKSET,

	"vpc_id": CHECKSET,

	"expired_time": CHECKSET,

	"region_id": CHECKSET,

	"tags": CHECKSET,
}

func AlibabacloudTestAccVpcIpv6GatewayBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
