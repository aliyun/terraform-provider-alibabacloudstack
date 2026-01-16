package alibabacloudstack

import (
	"testing"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbVservergroup0(t *testing.T) {

	var v *slb.DescribeVServerGroupAttributeResponse

	resourceId := "alibabacloudstack_slb_vservergroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbVservergroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribevservergroupattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbv_server_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbVservergroupBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id":   "${alibabacloudstack_slb.default.id}",
					"vserver_group_name": "vserver_group_name",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_ecs_instance.default.id}"},
							"port":       "80",
							"weight":     "100",
							"type":       "ecs",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_group_name":     "vserver_group_name",
						"servers.#":              "1",
						"servers.0.server_ids.#": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_ecs_instance.default.id}"},
							"port":       "80",
							"weight":     "100",
							"type":       "ecs",
						},
						{
							"server_ids": []string{"${alibabacloudstack_ecs_instance.default1.id}"},
							"port":       "100",
							"weight":     "80",
							"type":       "ecs",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#":              "2",
						"servers.0.server_ids.#": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id":   "${alibabacloudstack_slb.default.id}",
					"vserver_group_name": "vserver_group_name",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_ecs_instance.default.id}"},
							"port":       "80",
							"weight":     "100",
							"type":       "ecs",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_group_name":     "vserver_group_name",
						"servers.#":              "1",
						"servers.0.server_ids.#": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id":   "${alibabacloudstack_slb.default.id}",
					"vserver_group_name": "vserver_group_name",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_ecs_instance.default1.id}"},
							"port":       "80",
							"weight":     "100",
							"type":       "ecs",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_group_name":     "vserver_group_name",
						"servers.#":              "1",
						"servers.0.server_ids.#": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id":   "${alibabacloudstack_slb.default.id}",
					"vserver_group_name": "vserver_group_name",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_ecs_instance.default1.id}", "${alibabacloudstack_ecs_instance.default.id}"},
							"port":       "80",
							"weight":     "100",
							"type":       "ecs",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_group_name":     "vserver_group_name",
						"servers.#":              "1",
						"servers.0.server_ids.#": "2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id":   "${alibabacloudstack_slb.default.id}",
					"vserver_group_name": "vserver_group_name",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_ecs_instance.default1.id}"},
							"port":       "80",
							"weight":     "100",
							"type":       "ecs",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_group_name":     "vserver_group_name",
						"servers.#":              "1",
						"servers.0.server_ids.#": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"load_balancer_id":   "${alibabacloudstack_slb.default.id}",
					"vserver_group_name": "vserver_group_name",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_ecs_instance.default1.id}"},
							"port":       "81",
							"weight":     "80",
							"type":       "ecs",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_group_name":     "vserver_group_name",
						"servers.#":              "1",
						"servers.0.server_ids.#": "1",
						"servers.0.port":         "81",
						"servers.0.weight":       "80",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"vserver_group_name": "vserver_group_name_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_group_name": "vserver_group_name_update",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// delete_protection_validation is a local attribute and cannot be loaded from the remote
				// load_balancer_id cannot be read back on the private cloud side temporarily
				ImportStateVerifyIgnore: []string{"delete_protection_validation", "load_balancer_id"},
			},
		},
	})
}

var AlibabacloudTestAccSlbVservergroupCheckmap = map[string]string{

	// "v_server_group_id": CHECKSET,

	// "associated_objects": CHECKSET,

	// "v_server_group_name": CHECKSET,

	// "load_balancer_id": CHECKSET,

	// "backend_servers": CHECKSET,

	// "tags": CHECKSET,
}

func AlibabacloudTestAccSlbVservergroupBasicdependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alibabacloudstack_slb" "default" {
		name = "${var.name}"
		specification        = "slb.s2.small"
	}
	
	%s

	resource "alibabacloudstack_ecs_instance" "default1" {
		image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type        = "${local.default_instance_type_id}"
		system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
		system_disk_size     = 20
		system_disk_name     = "test_sys_disk"
		security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
		instance_name        = "${var.name}_ecs"
		vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
		zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
		is_outdated          = false
		lifecycle {
			ignore_changes = [
			instance_type
			]
		}
	}

	`, name, ECSInstanceCommonTestCase)
}

// func TestAccAlibabacloudStackSlbVservergroup1(t *testing.T) {

// 	var v map[string]interface{}

// 	resourceId := "alibabacloudstack_slb_vservergroup.default"
// 	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbVservergroupCheckmap)
// 	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
// 		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}, "DoSlbDescribevservergroupattributeRequest")
// 	rac := resourceAttrCheckInit(rc, ra)
// 	testAccCheck := rac.resourceAttrMapUpdateSet()

// 	rand := getAccTestRandInt(10000, 99999)
// 	name := fmt.Sprintf("tf-testacc%sslbv_server_group%d", defaultRegionToTest, rand)

// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbVservergroupBasicdependence)
// 	ResourceTest(t, resource.TestCase{
// 		PreCheck: func() {

// 			testAccPreCheck(t)
// 		},
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,

// 		CheckDestroy: rac.checkResourceDestroy(),

// 		Steps: []resource.TestStep{

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"load_balancer_id": "alibabacloudstack_slb.default.id",

// 					"v_server_group_name": "test-VServerGroupName",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"load_balancer_id": "alibabacloudstack_slb.default.id",

// 						"v_server_group_name": "test-VServerGroupName",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"v_server_group_name": "rdk-test-name",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"v_server_group_name": "rdk-test-name",
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
// func TestAccAlibabacloudStackSlbVservergroup2(t *testing.T) {

// 	var v map[string]interface{}

// 	resourceId := "alibabacloudstack_slb_vservergroup.default"
// 	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbVservergroupCheckmap)
// 	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
// 		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}, "DoSlbDescribevservergroupattributeRequest")
// 	rac := resourceAttrCheckInit(rc, ra)
// 	testAccCheck := rac.resourceAttrMapUpdateSet()

// 	rand := getAccTestRandInt(10000, 99999)
// 	name := fmt.Sprintf("tf-testacc%sslbv_server_group%d", defaultRegionToTest, rand)

// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbVservergroupBasicdependence)
// 	ResourceTest(t, resource.TestCase{
// 		PreCheck: func() {

// 			testAccPreCheck(t)
// 		},
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,

// 		CheckDestroy: rac.checkResourceDestroy(),

// 		Steps: []resource.TestStep{

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"load_balancer_id": "alibabacloudstack_slb.default.id",

// 					"v_server_group_name": "test-VServerGroupName",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"load_balancer_id": "alibabacloudstack_slb.default.id",

// 						"v_server_group_name": "test-VServerGroupName",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"v_server_group_name": "rdk-test-name99",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"v_server_group_name": "rdk-test-name99",
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
// func TestAccAlibabacloudStackSlbVservergroup3(t *testing.T) {

// 	var v map[string]interface{}

// 	resourceId := "alibabacloudstack_slb_vservergroup.default"
// 	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbVservergroupCheckmap)
// 	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
// 		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
// 	}, "DoSlbDescribevservergroupattributeRequest")
// 	rac := resourceAttrCheckInit(rc, ra)
// 	testAccCheck := rac.resourceAttrMapUpdateSet()

// 	rand := getAccTestRandInt(10000, 99999)
// 	name := fmt.Sprintf("tf-testacc%sslbv_server_group%d", defaultRegionToTest, rand)

// 	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbVservergroupBasicdependence)
// 	ResourceTest(t, resource.TestCase{
// 		PreCheck: func() {

// 			testAccPreCheck(t)
// 		},
// 		IDRefreshName: resourceId,
// 		Providers:     testAccProviders,

// 		CheckDestroy: rac.checkResourceDestroy(),

// 		Steps: []resource.TestStep{

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"v_server_group_name": "tfcreate",

// 					"load_balancer_id": "${{ref(resource, SLB::LoadBalancer::2.0.0.11.pre::slb.LoadBalancerId)}}",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"v_server_group_name": "tfcreate",

// 						"load_balancer_id": "${{ref(resource, SLB::LoadBalancer::2.0.0.11.pre::slb.LoadBalancerId)}}",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"v_server_group_name": "tfupdate",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"v_server_group_name": "tfupdate",
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
// 				Config: testAccConfig(map[string]interface{}{}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"v_server_group_name": "tfcreate",

// 					"load_balancer_id": "${{ref(resource, SLB::LoadBalancer::2.0.0.11.pre::slb.LoadBalancerId)}}",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"v_server_group_name": "tfcreate",

// 						"load_balancer_id": "${{ref(resource, SLB::LoadBalancer::2.0.0.11.pre::slb.LoadBalancerId)}}",
// 					}),
// 				),
// 			},

// 			{
// 				Config: testAccConfig(map[string]interface{}{

// 					"v_server_group_name": "tfupdate",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{

// 						"v_server_group_name": "tfupdate",
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
