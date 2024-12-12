package alibabacloudstack

import (
	"testing"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
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

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbload_balancer%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbLoadbalancerBasicdependence)
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

					"name":          "rdk_test_name",
					"specification": "slb.s1.small",
					"vswitch_id":    "${alibabacloudstack_vswitch.default.id}",
					"address_type":  "intranet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"name":       "rdk_test_name",
						"vswitch_id": CHECKSET,

						"address_type": "intranet",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"specification": "slb.s2.small",

					"name": "Rdk-test-name",
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

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name          = "%s"
  cidr_block        = "192.168.0.0/16"
}


#vsw
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "192.168.0.0/16"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

`,

		name, name)
}

// func TestAccAlibabacloudStackSlbLoadbalancer1(t *testing.T) {

// 	var v map[string]interface{}

// 	resourceId := "alibabacloudstack_slb_loadbalancer.default"
// 	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbLoadbalancerCheckmap)
// 	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
// 		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}, "DoSlbDescribeloadbalancerattributeRequest")
// 	rac := resourceAttrCheckInit(rc, ra)
// 	testAccCheck := rac.resourceAttrMapUpdateSet()

// 	rand := acctest.RandIntRange(10000, 99999)
// 	name := fmt.Sprintf("tf-testacc%sslbload_balancer%d", defaultRegionToTest, rand)

// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbLoadbalancerBasicdependence)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() {

// 			testAccPreCheck(t)
// 		},
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,

// 		CheckDestroy: rac.checkResourceDestroy(),

// 		Steps: []resource.TestStep{

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"internet_charge_type": "PayByBandwidth",

// 					"name": "rdk_test_name",

// 					"address_ip_version": "ipv4",

//

//

// 					"payment_type": "PayOnDemand",

// 					"address_type": "intranet",

// 					"vswitch_id": "alibabacloudstack_vswitch.default.id",

// 					"slave_zone_id": "cn-hangzhou-h",

// 					"vpc_id": "alibabacloudstack_vpc.default.id",

// 					"master_zone_id": "cn-hangzhou-j",

// 					"address": "10.40.0.39",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"internet_charge_type": "PayByBandwidth",

// 						"name": "rdk_test_name",

// 						"address_ip_version": "ipv4",

//

//

// 						"payment_type": "PayOnDemand",

// 						"address_type": "intranet",

// 						"vswitch_id": "alibabacloudstack_vswitch.default.id",

// 						"slave_zone_id": "cn-hangzhou-h",

// 						"vpc_id": "alibabacloudstack_vpc.default.id",

// 						"master_zone_id": "cn-hangzhou-j",

// 						"address": "10.40.0.39",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"status": "active",

// 					"name": "Rdk-test-name",

//

// 					"payment_type": "PrePay",

//
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"status": "active",

// 						"name": "Rdk-test-name",

//

// 						"payment_type": "PrePay",

//
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF",
// 						"For":     "Test",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF",
// 						"tags.For":     "Test",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF-update",
// 						"For":     "Test-update",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF-update",
// 						"tags.For":     "Test-update",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": REMOVEKEY,
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "0",
// 						"tags.Created": REMOVEKEY,
// 						"tags.For":     REMOVEKEY,
// 					}),
// 				),
// 			},
// 		},
// 	})
// }
// func TestAccAlibabacloudStackSlbLoadbalancer2(t *testing.T) {

// 	var v map[string]interface{}

// 	resourceId := "alibabacloudstack_slb_loadbalancer.default"
// 	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbLoadbalancerCheckmap)
// 	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
// 		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}, "DoSlbDescribeloadbalancerattributeRequest")
// 	rac := resourceAttrCheckInit(rc, ra)
// 	testAccCheck := rac.resourceAttrMapUpdateSet()

// 	rand := acctest.RandIntRange(10000, 99999)
// 	name := fmt.Sprintf("tf-testacc%sslbload_balancer%d", defaultRegionToTest, rand)

// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbLoadbalancerBasicdependence)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() {

// 			testAccPreCheck(t)
// 		},
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,

// 		CheckDestroy: rac.checkResourceDestroy(),

// 		Steps: []resource.TestStep{

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"internet_charge_type": "PayByBandwidth",

// 					"name": "rdk_test_name_as",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"internet_charge_type": "PayByBandwidth",

// 						"name": "rdk_test_name_as",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF",
// 						"For":     "Test",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF",
// 						"tags.For":     "Test",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF-update",
// 						"For":     "Test-update",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF-update",
// 						"tags.For":     "Test-update",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": REMOVEKEY,
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "0",
// 						"tags.Created": REMOVEKEY,
// 						"tags.For":     REMOVEKEY,
// 					}),
// 				),
// 			},
// 		},
// 	})
// }
// func TestAccAlibabacloudStackSlbLoadbalancer3(t *testing.T) {

// 	var v map[string]interface{}

// 	resourceId := "alibabacloudstack_slb_loadbalancer.default"
// 	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbLoadbalancerCheckmap)
// 	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
// 		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}, "DoSlbDescribeloadbalancerattributeRequest")
// 	rac := resourceAttrCheckInit(rc, ra)
// 	testAccCheck := rac.resourceAttrMapUpdateSet()

// 	rand := acctest.RandIntRange(10000, 99999)
// 	name := fmt.Sprintf("tf-testacc%sslbload_balancer%d", defaultRegionToTest, rand)

// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbLoadbalancerBasicdependence)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() {

// 			testAccPreCheck(t)
// 		},
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,

// 		CheckDestroy: rac.checkResourceDestroy(),

// 		Steps: []resource.TestStep{

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"internet_charge_type": "PayByBandwidth",

// 					"name": "rdk_test_name_as",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"internet_charge_type": "PayByBandwidth",

// 						"name": "rdk_test_name_as",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"internet_charge_type": "PayByBandwidth",

// 					"name": "rdk_test_name_as",

// 					"status": "inactive",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"internet_charge_type": "PayByBandwidth",

// 						"name": "rdk_test_name_as",

// 						"status": "inactive",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF",
// 						"For":     "Test",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF",
// 						"tags.For":     "Test",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF-update",
// 						"For":     "Test-update",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF-update",
// 						"tags.For":     "Test-update",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": REMOVEKEY,
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "0",
// 						"tags.Created": REMOVEKEY,
// 						"tags.For":     REMOVEKEY,
// 					}),
// 				),
// 			},
// 		},
// 	})
// }
// func TestAccAlibabacloudStackSlbLoadbalancer4(t *testing.T) {

// 	var v map[string]interface{}

// 	resourceId := "alibabacloudstack_slb_loadbalancer.default"
// 	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbLoadbalancerCheckmap)
// 	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
// 		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}, "DoSlbDescribeloadbalancerattributeRequest")
// 	rac := resourceAttrCheckInit(rc, ra)
// 	testAccCheck := rac.resourceAttrMapUpdateSet()

// 	rand := acctest.RandIntRange(10000, 99999)
// 	name := fmt.Sprintf("tf-testacc%sslbload_balancer%d", defaultRegionToTest, rand)

// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbLoadbalancerBasicdependence)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() {

// 			testAccPreCheck(t)
// 		},
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,

// 		CheckDestroy: rac.checkResourceDestroy(),

// 		Steps: []resource.TestStep{

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"internet_charge_type": "PayByBandwidth",

// 					"name": "rdk_test_name_as",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"internet_charge_type": "PayByBandwidth",

// 						"name": "rdk_test_name_as",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"internet_charge_type": "PayByBandwidth",

// 					"name": "rdk_test_name_as",

// 					"status": "inactive",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"internet_charge_type": "PayByBandwidth",

// 						"name": "rdk_test_name_as",

// 						"status": "inactive",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"internet_charge_type": "PayByBandwidth",

// 					"name": "rdk_test_name_as",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"internet_charge_type": "PayByBandwidth",

// 						"name": "rdk_test_name_as",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF",
// 						"For":     "Test",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF",
// 						"tags.For":     "Test",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF-update",
// 						"For":     "Test-update",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF-update",
// 						"tags.For":     "Test-update",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": REMOVEKEY,
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "0",
// 						"tags.Created": REMOVEKEY,
// 						"tags.For":     REMOVEKEY,
// 					}),
// 				),
// 			},
// 		},
// 	})
// }
// func TestAccAlibabacloudStackSlbLoadbalancer5(t *testing.T) {

// 	var v map[string]interface{}

// 	resourceId := "alibabacloudstack_slb_loadbalancer.default"
// 	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbLoadbalancerCheckmap)
// 	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
// 		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}, "DoSlbDescribeloadbalancerattributeRequest")
// 	rac := resourceAttrCheckInit(rc, ra)
// 	testAccCheck := rac.resourceAttrMapUpdateSet()

// 	rand := acctest.RandIntRange(10000, 99999)
// 	name := fmt.Sprintf("tf-testacc%sslbload_balancer%d", defaultRegionToTest, rand)

// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbLoadbalancerBasicdependence)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() {

// 			testAccPreCheck(t)
// 		},
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,

// 		CheckDestroy: rac.checkResourceDestroy(),

// 		Steps: []resource.TestStep{

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"status": "active",

// 					"address_ip_version": "ipv4",

// 					"address": "192.168.10.100",

// 					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::vsw1.VSwitchId)}}",

// 					"slave_zone_id": "cn-hangzhou-h",

//

// 					"name": "tf-createname",

//

// 					"vpc_id": "alibabacloudstack_vpc.default.id",

// 					"payment_type": "PayAsYouGo",

// 					"modification_protection_reason": "test",

// 					"address_type": "intranet",

// 					"master_zone_id": "cn-hangzhou-j",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"status": "active",

// 						"address_ip_version": "ipv4",

// 						"address": "192.168.10.100",

// 						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::vsw1.VSwitchId)}}",

// 						"slave_zone_id": "cn-hangzhou-h",

//

// 						"name": "tf-createname",

//

// 						"vpc_id": "alibabacloudstack_vpc.default.id",

// 						"payment_type": "PayAsYouGo",

// 						"modification_protection_reason": "test",

// 						"address_type": "intranet",

// 						"master_zone_id": "cn-hangzhou-j",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"name": "tf-update",

// 					"delete_protection": "on",

// 					"modification_protection_reason": "test-update",

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"name": "tf-update",

// 						"delete_protection": "on",

// 						"modification_protection_reason": "test-update",

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

//

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

//

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"status": "inactive",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"status": "inactive",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"status": "active",

// 					"address_ip_version": "ipv4",

// 					"address": "192.168.10.100",

// 					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::vsw1.VSwitchId)}}",

// 					"slave_zone_id": "cn-hangzhou-h",

//

// 					"name": "tf-createname",

//

// 					"vpc_id": "alibabacloudstack_vpc.default.id",

// 					"payment_type": "PayAsYouGo",

// 					"modification_protection_reason": "test",

// 					"address_type": "intranet",

// 					"master_zone_id": "cn-hangzhou-j",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"status": "active",

// 						"address_ip_version": "ipv4",

// 						"address": "192.168.10.100",

// 						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::vsw1.VSwitchId)}}",

// 						"slave_zone_id": "cn-hangzhou-h",

//

// 						"name": "tf-createname",

//

// 						"vpc_id": "alibabacloudstack_vpc.default.id",

// 						"payment_type": "PayAsYouGo",

// 						"modification_protection_reason": "test",

// 						"address_type": "intranet",

// 						"master_zone_id": "cn-hangzhou-j",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"name": "tf-update",

// 					"delete_protection": "on",

// 					"modification_protection_reason": "test-update",

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"name": "tf-update",

// 						"delete_protection": "on",

// 						"modification_protection_reason": "test-update",

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"name": "tf-update",

// 					"delete_protection": "on",

// 					"modification_protection_reason": "test-update",

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"name": "tf-update",

// 						"delete_protection": "on",

// 						"modification_protection_reason": "test-update",

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"name": "tf-update",

// 					"delete_protection": "on",

// 					"modification_protection_reason": "test-update",

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"name": "tf-update",

// 						"delete_protection": "on",

// 						"modification_protection_reason": "test-update",

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"name": "tf-update",

// 					"delete_protection": "on",

// 					"modification_protection_reason": "test-update",

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"name": "tf-update",

// 						"delete_protection": "on",

// 						"modification_protection_reason": "test-update",

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"status": "active",

// 					"address_ip_version": "ipv4",

// 					"address": "192.168.10.100",

// 					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::vsw1.VSwitchId)}}",

// 					"slave_zone_id": "cn-hangzhou-h",

//

// 					"name": "tf-createname",

//

// 					"vpc_id": "alibabacloudstack_vpc.default.id",

// 					"payment_type": "PayAsYouGo",

// 					"modification_protection_reason": "test",

// 					"address_type": "intranet",

// 					"master_zone_id": "cn-hangzhou-j",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"status": "active",

// 						"address_ip_version": "ipv4",

// 						"address": "192.168.10.100",

// 						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::vsw1.VSwitchId)}}",

// 						"slave_zone_id": "cn-hangzhou-h",

//

// 						"name": "tf-createname",

//

// 						"vpc_id": "alibabacloudstack_vpc.default.id",

// 						"payment_type": "PayAsYouGo",

// 						"modification_protection_reason": "test",

// 						"address_type": "intranet",

// 						"master_zone_id": "cn-hangzhou-j",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"status": "active",

// 					"address_ip_version": "ipv4",

// 					"address": "192.168.10.100",

// 					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::vsw1.VSwitchId)}}",

// 					"slave_zone_id": "cn-hangzhou-h",

//

// 					"name": "tf-createname",

//

// 					"vpc_id": "alibabacloudstack_vpc.default.id",

// 					"payment_type": "PayAsYouGo",

// 					"modification_protection_reason": "test",

// 					"address_type": "intranet",

// 					"master_zone_id": "cn-hangzhou-j",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"status": "active",

// 						"address_ip_version": "ipv4",

// 						"address": "192.168.10.100",

// 						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.3.2.pre::vsw1.VSwitchId)}}",

// 						"slave_zone_id": "cn-hangzhou-h",

//

// 						"name": "tf-createname",

//

// 						"vpc_id": "alibabacloudstack_vpc.default.id",

// 						"payment_type": "PayAsYouGo",

// 						"modification_protection_reason": "test",

// 						"address_type": "intranet",

// 						"master_zone_id": "cn-hangzhou-j",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF",
// 						"For":     "Test",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF",
// 						"tags.For":     "Test",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF-update",
// 						"For":     "Test-update",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF-update",
// 						"tags.For":     "Test-update",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": REMOVEKEY,
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "0",
// 						"tags.Created": REMOVEKEY,
// 						"tags.For":     REMOVEKEY,
// 					}),
// 				),
// 			},
// 		},
// 	})
// }
// func TestAccAlibabacloudStackSlbLoadbalancer6(t *testing.T) {

// 	var v map[string]interface{}

// 	resourceId := "alibabacloudstack_slb_loadbalancer.default"
// 	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbLoadbalancerCheckmap)
// 	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
// 		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}, "DoSlbDescribeloadbalancerattributeRequest")
// 	rac := resourceAttrCheckInit(rc, ra)
// 	testAccCheck := rac.resourceAttrMapUpdateSet()

// 	rand := acctest.RandIntRange(10000, 99999)
// 	name := fmt.Sprintf("tf-testacc%sslbload_balancer%d", defaultRegionToTest, rand)

// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbLoadbalancerBasicdependence)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() {

// 			testAccPreCheck(t)
// 		},
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,

// 		CheckDestroy: rac.checkResourceDestroy(),

// 		Steps: []resource.TestStep{

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"slave_zone_id": "cn-hangzhou-h",

//

// 					"name": "tf-createname",

//

// 					"payment_type": "PayAsYouGo",

// 					"modification_protection_reason": "test",

// 					"address_type": "internet",

// 					"master_zone_id": "cn-hangzhou-j",

// 					"address_ip_version": "ipv4",

// 					"internet_charge_type": "PayByBandwidth",

//

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"slave_zone_id": "cn-hangzhou-h",

//

// 						"name": "tf-createname",

//

// 						"payment_type": "PayAsYouGo",

// 						"modification_protection_reason": "test",

// 						"address_type": "internet",

// 						"master_zone_id": "cn-hangzhou-j",

// 						"address_ip_version": "ipv4",

// 						"internet_charge_type": "PayByBandwidth",

//

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

// 					"delete_protection": "on",

// 					"modification_protection_reason": "test-update",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

// 						"delete_protection": "on",

// 						"modification_protection_reason": "test-update",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

//

//
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

//

//
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"name": "tf-update",

// 					"internet_charge_type": "paybybandwidth",

//

// 					"payment_type": "PayAsYouGo",

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"name": "tf-update",

// 						"internet_charge_type": "paybybandwidth",

//

// 						"payment_type": "PayAsYouGo",

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF",
// 						"For":     "Test",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF",
// 						"tags.For":     "Test",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF-update",
// 						"For":     "Test-update",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF-update",
// 						"tags.For":     "Test-update",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": REMOVEKEY,
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "0",
// 						"tags.Created": REMOVEKEY,
// 						"tags.For":     REMOVEKEY,
// 					}),
// 				),
// 			},
// 		},
// 	})
// }
// func TestAccAlibabacloudStackSlbLoadbalancer7(t *testing.T) {

// 	var v map[string]interface{}

// 	resourceId := "alibabacloudstack_slb_loadbalancer.default"
// 	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbLoadbalancerCheckmap)
// 	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
// 		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}, "DoSlbDescribeloadbalancerattributeRequest")
// 	rac := resourceAttrCheckInit(rc, ra)
// 	testAccCheck := rac.resourceAttrMapUpdateSet()

// 	rand := acctest.RandIntRange(10000, 99999)
// 	name := fmt.Sprintf("tf-testacc%sslbload_balancer%d", defaultRegionToTest, rand)

// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbLoadbalancerBasicdependence)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() {

// 			testAccPreCheck(t)
// 		},
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,

// 		CheckDestroy: rac.checkResourceDestroy(),

// 		Steps: []resource.TestStep{

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"slave_zone_id": "cn-hangzhou-h",

// 					"name": "tf-createname",

// 					"payment_type": "Subscription",

// 					"address_type": "internet",

// 					"master_zone_id": "cn-hangzhou-j",

// 					"address_ip_version": "ipv4",

// 					"internet_charge_type": "PayByBandwidth",

//

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"slave_zone_id": "cn-hangzhou-h",

// 						"name": "tf-createname",

// 						"payment_type": "Subscription",

// 						"address_type": "internet",

// 						"master_zone_id": "cn-hangzhou-j",

// 						"address_ip_version": "ipv4",

// 						"internet_charge_type": "PayByBandwidth",

//

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"name": "tf-update",

// 					"payment_type": "Subscription",

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"name": "tf-update",

// 						"payment_type": "Subscription",

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"status": "active",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"status": "active",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF",
// 						"For":     "Test",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF",
// 						"tags.For":     "Test",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": map[string]string{
// 						"Created": "TF-update",
// 						"For":     "Test-update",
// 					},
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "2",
// 						"tags.Created": "TF-update",
// 						"tags.For":     "Test-update",
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"tags": REMOVEKEY,
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"tags.%":       "0",
// 						"tags.Created": REMOVEKEY,
// 						"tags.For":     REMOVEKEY,
// 					}),
// 				),
// 			},
// 		},
// 	})
// }
