package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsNetworkInterface0(t *testing.T) {
	var v ecs.NetworkInterfaceSet

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
					"network_interface_name": name,
					"vswitch_id":             "${alibabacloudstack_vpc_vswitch.default.id}",
					"security_groups":        []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
					"description":            name,
					"private_ip":             "172.16.1.200",
					"private_ips":            []string{"172.16.1.1", "172.16.1.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interface_name": name,
						"private_ips.#":          "2",
						"security_groups.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []string{"${alibabacloudstack_ecs_securitygroup.default.id}", "${alibabacloudstack_ecs_securitygroup.update.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []string{"${alibabacloudstack_ecs_securitygroup.update.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_interface_name": name + "_update",
					"description":            name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interface_name": name + "_update",
						"description":            name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ips": []string{"172.16.1.2", "172.16.1.3", "172.16.1.4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ips.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ips":       REMOVEKEY,
					"private_ips_count": 2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ips.#":     REMOVEKEY, // FIXME:  Strange errors caused by Terraform behavior
						"private_ips_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ips_count": 4,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ips_count": "4",
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
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudTestAccEcsNetworkinterfaceCheckmap = map[string]string{
	"description":            CHECKSET,
	"primary_ip_address":     CHECKSET,
	"network_interface_name": CHECKSET,
	"mac_address":            CHECKSET,
}

func AlibabacloudTestAccEcsNetworkinterfaceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

resource "alibabacloudstack_ecs_securitygroup" "update" {
  name   = "${var.name}_sg1"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

`, name, SecurityGroupCommonTestCase)
}
