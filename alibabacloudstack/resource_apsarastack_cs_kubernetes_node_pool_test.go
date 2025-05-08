package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCSKubernetesNodePool_basic(t *testing.T) {
	var v *NodePoolDetail

	resourceId := "alibabacloudstack_cs_kubernetes_node_pool.default"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"cluster_id":            "c89eeac401e7b43d985c6ac2b94ceee66",
					"vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"instance_types":        []string{"${local.default_instance_type_id}"},
					"node_count":            "1",
					"password":              "1qaz@WSX",
					"system_disk_category":  "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"data_disks":            []map[string]string{{"size": "100", "category": "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"}},
					"tags":                  map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					//"management":            []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "0", "max_unavailable": "0"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"cluster_id":            CHECKSET,
						"vswitch_ids.#":         "1",
						"instance_types.#":      "1",
						"node_count":            "1",
						"key_name":              CHECKSET,
						"system_disk_category":  CHECKSET,
						"system_disk_size":      "40",
						"install_cloud_monitor": "false",
						"data_disks.#":          "1",
						"data_disks.0.size":     "100",
						"data_disks.0.category": CHECKSET,
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.Foo":              "Bar",
						// 						"management.#":                 "1",
						// 						"management.0.auto_repair":     "true",
						// 						"management.0.auto_upgrade":    "true",
						// 						"management.0.surge":           "0",
						// 						"management.0.max_unavailable": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// check: scale out
			{
				Config: testAccConfig(map[string]interface{}{
					"node_count":       "2",
					"system_disk_size": "80",
					"data_disks":       []map[string]string{{"size": "40", "category": "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"}},
					//"management":       []map[string]string{{"auto_repair": "true", "auto_upgrade": "true", "surge": "1", "max_unavailable": "1"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_count":            "2",
						"system_disk_size":      "80",
						"data_disks.#":          "1",
						"data_disks.0.size":     "40",
						"data_disks.0.category": "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
						// 						"management.#":                 "1",
						// 						"management.0.auto_repair":     "true",
						// 						"management.0.auto_upgrade":    "true",
						// 						"management.0.surge":           "1",
						// 						"management.0.max_unavailable": "1",
					}),
				),
			},
			// check: remove nodes
			{
				Config: testAccConfig(map[string]interface{}{
					"node_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_count": "1",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackCSKubernetesNodePool_autoScaling(t *testing.T) {
	var v *NodePoolDetail

	resourceId := "alibabacloudstack_cs_kubernetes_node_pool.autocaling"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"cluster_id":            "c180d1d233d2d47f68f301b129f622665",
					"vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"instance_types":        []string{"${local.default_instance_type_id}"},
					"key_name":              "${alibabacloudstack_key_pair.default.key_name}",
					"system_disk_category":  "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"platform":              "AliyunLinux",
					"scaling_policy":        "release",
					"scaling_config":        []map[string]string{{"min_size": "1", "max_size": "10", "type": "cpu", "is_bond_eip": "true", "eip_internet_charge_type": "PayByBandwidth", "eip_bandwidth": "5"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                         name,
						"cluster_id":                   CHECKSET,
						"vswitch_ids.#":                "1",
						"instance_types.#":             "1",
						"key_name":                     CHECKSET,
						"system_disk_category":         "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
						"system_disk_size":             "40",
						"install_cloud_monitor":        "false",
						"platform":                     "AliyunLinux",
						"scaling_policy":               "release",
						"scaling_config.#":             "1",
						"scaling_config.0.min_size":    "1",
						"scaling_config.0.max_size":    "10",
						"scaling_config.0.type":        "cpu",
						"scaling_config.0.is_bond_eip": "true",
						"scaling_config.0.eip_internet_charge_type": "PayByBandwidth",
						"scaling_config.0.eip_bandwidth":            "5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "node_count"},
			},
			// check: update config
			{
				Config: testAccConfig(map[string]interface{}{
					"platform":       "AliyunLinux",
					"scaling_policy": "release",
					"scaling_config": []map[string]string{{"min_size": "1", "max_size": "20", "type": "cpu", "is_bond_eip": "true", "eip_internet_charge_type": "PayByBandwidth", "eip_bandwidth": "5"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"platform":                                  "AliyunLinux",
						"scaling_policy":                            "release",
						"scaling_config.#":                          "1",
						"scaling_config.0.min_size":                 "1",
						"scaling_config.0.max_size":                 "20",
						"scaling_config.0.type":                     "cpu",
						"scaling_config.0.is_bond_eip":              "true",
						"scaling_config.0.eip_internet_charge_type": "PayByBandwidth",
						"scaling_config.0.eip_bandwidth":            "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_config": []map[string]string{{"min_size": "1", "max_size": "20", "type": "cpu", "is_bond_eip": "false", "eip_internet_charge_type": "PayByBandwidth", "eip_bandwidth": "5"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_config.#":                          "1",
						"scaling_config.0.min_size":                 "1",
						"scaling_config.0.max_size":                 "20",
						"scaling_config.0.type":                     "cpu",
						"scaling_config.0.is_bond_eip":              "false",
						"scaling_config.0.eip_internet_charge_type": "PayByBandwidth",
						"scaling_config.0.eip_bandwidth":            "5",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackCSKubernetesNodePool_PrePaid(t *testing.T) {
	var v *NodePoolDetail

	resourceId := "alibabacloudstack_cs_kubernetes_node_pool.pre_paid_nodepool"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"cluster_id":            "c180d1d233d2d47f68f301b129f622665",
					"vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"password":              "Terraform1234",
					"instance_types":        []string{"${local.default_instance_type_id}"},
					"system_disk_category":  "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"system_disk_size":      "120",
					"install_cloud_monitor": "false",
					"instance_charge_type":  "PrePaid",
					"period":                "1",
					"period_unit":           "Month",
					"auto_renew":            "true",
					"auto_renew_period":     "1",
					"scaling_config":        []map[string]string{{"min_size": "1", "max_size": "10"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                      name,
						"cluster_id":                CHECKSET,
						"password":                  CHECKSET,
						"vswitch_ids.#":             "1",
						"instance_types.#":          "1",
						"system_disk_category":      "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
						"system_disk_size":          "120",
						"instance_charge_type":      "PrePaid",
						"install_cloud_monitor":     "false",
						"period":                    "1",
						"period_unit":               "Month",
						"auto_renew":                "true",
						"auto_renew_period":         "1",
						"scaling_config.#":          "1",
						"scaling_config.0.min_size": "1",
						"scaling_config.0.max_size": "10",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type":  "PrePaid",
					"auto_renew_period":     "2",
					"install_cloud_monitor": "true",
					"scaling_config":        []map[string]string{{"min_size": "2", "max_size": "10"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type":      "PrePaid",
						"auto_renew_period":         "2",
						"install_cloud_monitor":     "true",
						"scaling_config.#":          "1",
						"scaling_config.0.min_size": "2",
						"scaling_config.0.max_size": "10",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackCSKubernetesNodePool_Spot(t *testing.T) {
	var v *NodePoolDetail

	resourceId := "alibabacloudstack_cs_kubernetes_node_pool.spot_nodepool"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePool-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                       name,
					"cluster_id":                 "c180d1d233d2d47f68f301b129f622665",
					"vswitch_ids":                []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"instance_types":             []string{"${local.default_instance_type_id}"},
					"system_disk_category":       "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"system_disk_size":           "120",
					"resource_group_id":          "8",
					"password":                   "Terraform1234",
					"node_count":                 "1",
					"install_cloud_monitor":      "false",
					"internet_charge_type":       "PayByTraffic",
					"internet_max_bandwidth_out": "5",
					"spot_strategy":              "SpotWithPriceLimit",
					"spot_price_limit": []map[string]string{
						{
							"instance_type": "${local.default_instance_type_id}",
							"price_limit":   "0.57",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                             name,
						"cluster_id":                       CHECKSET,
						"vswitch_ids.#":                    "1",
						"instance_types.#":                 "1",
						"system_disk_category":             "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
						"system_disk_size":                 "120",
						"resource_group_id":                CHECKSET,
						"password":                         CHECKSET,
						"node_count":                       "1",
						"install_cloud_monitor":            "false",
						"internet_charge_type":             "PayByTraffic",
						"internet_max_bandwidth_out":       "5",
						"spot_strategy":                    "SpotWithPriceLimit",
						"spot_price_limit.#":               "1",
						"spot_price_limit.0.instance_type": CHECKSET,
						"spot_price_limit.0.price_limit":   "0.57",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_charge_type":       "PayByTraffic",
					"internet_max_bandwidth_out": "10",
					"spot_price_limit": []map[string]string{
						{
							"instance_type": "${local.default_instance_type_id}",
							"price_limit":   "0.60",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type":             "PayByTraffic",
						"internet_max_bandwidth_out":       "10",
						"spot_strategy":                    "SpotWithPriceLimit",
						"spot_price_limit.#":               "1",
						"spot_price_limit.0.instance_type": CHECKSET,
						"spot_price_limit.0.price_limit":   "0.60",
					}),
				),
			},
		},
	})
}

var csdKubernetesNodePoolBasicMap = map[string]string{
	"system_disk_size":     "40",
	"system_disk_category": "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
}

func resourceCSNodePoolConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s




`, name, VSwitchCommonTestCase + DataAlibabacloudstackInstanceTypes)
}
