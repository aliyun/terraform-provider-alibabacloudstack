package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbListener0(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_listener.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbListenerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribeloadbalancerhttplistenerattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslblistener%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
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

					"load_balancer_id": "alibabacloudstack_slb.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"load_balancer_id": "alibabacloudstack_slb.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",
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

var AlibabacloudTestAccSlbListenerCheckmap = map[string]string{

	"description": CHECKSET,

	"scheduler": CHECKSET,

	"acl_id": CHECKSET,

	"server_group_id": CHECKSET,

	"load_balancer_id": CHECKSET,

	"acl_status": CHECKSET,

	"bandwidth": CHECKSET,

	"acl_type": CHECKSET,
}

func AlibabacloudTestAccSlbListenerBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}




`, name)
}
func TestAccAlibabacloudStackSlbListener1(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_listener.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbListenerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribeloadbalancerhttplistenerattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslblistener%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
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

					"listener_protocol": "https",

					"load_balancer_id": "alibabacloudstack_slb.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"listener_protocol": "https",

						"load_balancer_id": "alibabacloudstack_slb.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"listener_protocol": "https",

					"load_balancer_id": "alibabacloudstack_slb.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"listener_protocol": "https",

						"load_balancer_id": "alibabacloudstack_slb.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
func TestAccAlibabacloudStackSlbListener2(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_listener.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbListenerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribeloadbalancerhttplistenerattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslblistener%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
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

					"listener_protocol": "tcp",

					"load_balancer_id": "alibabacloudstack_slb.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"listener_protocol": "tcp",

						"load_balancer_id": "alibabacloudstack_slb.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"listener_protocol": "tcp",

					"load_balancer_id": "alibabacloudstack_slb.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"listener_protocol": "tcp",

						"load_balancer_id": "alibabacloudstack_slb.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
func TestAccAlibabacloudStackSlbListener3(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_listener.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbListenerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribeloadbalancerhttplistenerattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslblistener%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbListenerBasicdependence)
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

					"load_balancer_id": "${{ref(resource, SLB::LoadBalancer::2.0.0.11.pre::slb.LoadBalancerId)}}",

					"listener_protocol": "http",

					"acl_status": "on",

					"acl_type": "white",

					"description": "testcreate",

					"scheduler": "wrr",

					"acl_id": "${{ref(resource, SLB::AccessControlList::3.0.0::ac.AclId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"load_balancer_id": "${{ref(resource, SLB::LoadBalancer::2.0.0.11.pre::slb.LoadBalancerId)}}",

						"listener_protocol": "http",

						"acl_status": "on",

						"acl_type": "white",

						"description": "testcreate",

						"scheduler": "wrr",

						"acl_id": "${{ref(resource, SLB::AccessControlList::3.0.0::ac.AclId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"listener_protocol": "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"listener_protocol": "http",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"listener_protocol": "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"listener_protocol": "http",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"listener_protocol": "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"listener_protocol": "http",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"listener_protocol": "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"listener_protocol": "http",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "testupdate",

					"scheduler": "rr",

					"listener_protocol": "http",

					"acl_id": "${{ref(resource, SLB::AccessControlList::3.0.0::ac02.AclId)}}",

					"acl_status": "on",

					"acl_type": "black",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "testupdate",

						"scheduler": "rr",

						"listener_protocol": "http",

						"acl_id": "${{ref(resource, SLB::AccessControlList::3.0.0::ac02.AclId)}}",

						"acl_status": "on",

						"acl_type": "black",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"listener_protocol": "http",

					"acl_status": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"listener_protocol": "http",

						"acl_status": "off",
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

