package alibabacloudstack

import (
	"testing"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbLoadbalancer0(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse

	resourceId := "alibabacloudstack_slb_loadbalancer.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbLoadbalancerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribeloadbalancerattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbload_balancer%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbLoadbalancerBasicdependence)
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

					"name":          "rdk_test_name",
					"specification": "slb.s1.small",
					"vswitch_id":    "${alibabacloudstack_vpc_vswitch.default.id}",
					"address_type":  "intranet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"name":         "rdk_test_name",
						"vswitch_id":   CHECKSET,
						"address_type": "intranet",
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
					"specification": "slb.s2.small",
					"name":          "Rdk-test-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"specification": "slb.s2.small",

						"name": "Rdk-test-name",
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

func TestAccAlibabacloudStackSlbLoadbalancerClassic(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse

	resourceId := "alibabacloudstack_slb_loadbalancer.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbLoadbalancerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribeloadbalancerattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbload_balancer%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbLoadbalancerBasicdependence)
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

					"name":          "rdk_test_name",
					"specification": "slb.s1.small",
					"network_type":  "classic",
					"address_type":  "intranet",
					"address":       "10.205.44.221",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"name":         "rdk_test_name",
						"address":      "10.205.44.221",
						"address_type": "intranet",
						"network_type": "classic",
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
					"specification": "slb.s2.small",
					"name":          "Rdk-test-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"specification": "slb.s2.small",

						"name": "Rdk-test-name",
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

var AlibabacloudTestAccSlbLoadbalancerCheckmap = map[string]string{

	// "renewal_duration": CHECKSET,

	// "address_ip_version": CHECKSET,

	"address": CHECKSET,

	// "end_time": CHECKSET,

	// "listener_ports_and_protocal": CHECKSET,

	// "resource_group_id": CHECKSET,

	// "listener_ports_and_protocol": CHECKSET,

	// "load_balancer_id": CHECKSET,

	// "backend_servers": CHECKSET,

	// "network_type": CHECKSET,

	// "bandwidth": CHECKSET,

	// "modification_protection_reason": CHECKSET,

	// "payment_type": CHECKSET,

	// "master_zone_id": CHECKSET,

	// "tags": CHECKSET,

	// "status": CHECKSET,

	// "create_time": CHECKSET,

	// "vswitch_id": CHECKSET,

	// "renewal_status": CHECKSET,

	// "renewal_cyc_unit": CHECKSET,

	// "slave_zone_id": CHECKSET,

	// "internet_charge_type": CHECKSET,

	// "region_id_alias": CHECKSET,

	// "name": CHECKSET,

	// "delete_protection": CHECKSET,

	// "vpc_id": CHECKSET,

	// "end_time_stamp": CHECKSET,

	// "region_id": CHECKSET,

	"address_type": CHECKSET,

	// "create_time_stamp": CHECKSET,

	// "auto_release_time": CHECKSET,
}

func AlibabacloudTestAccSlbLoadbalancerBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s 
`,

		name, VSwitchCommonTestCase)
}
