package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsNetworkinterface0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_networkinterface.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsNetworkinterfaceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribenetworkinterfacesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsnetwork_interface%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsNetworkinterfaceBasicdependence)
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

					"network_interface_name": "eni-test-ecs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"network_interface_name": "eni-test-ecs",
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

var AlibabacloudTestAccEcsNetworkinterfaceCheckmap = map[string]string{

	"description": CHECKSET,

	"resource_group_id": CHECKSET,

	"service_managed": CHECKSET,

	"attachment": CHECKSET,

	"primary_ip_address": CHECKSET,

	"network_interface_id": CHECKSET,

	"ipv6_prefix": CHECKSET,

	"service_id": CHECKSET,

	"ipv6_sets": CHECKSET,

	"associated_public_ip": CHECKSET,

	"ipv4_prefix": CHECKSET,

	"tags": CHECKSET,

	"status": CHECKSET,

	"zone_id": CHECKSET,

	"instance_id": CHECKSET,

	"create_time": CHECKSET,

	"vswitch_id": CHECKSET,

	"network_interface_name": CHECKSET,

	"mac_address": CHECKSET,

	"security_group_ids": CHECKSET,

	"type": CHECKSET,

	"queue_number": CHECKSET,

	"vpc_id": CHECKSET,

	"private_ip_sets": CHECKSET,
}

func AlibabacloudTestAccEcsNetworkinterfaceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
