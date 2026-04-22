package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCSKubernetesNodePool_basic(t *testing.T) {
	var v *NodePoolAlone

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
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 1: 创建基础节点池（使用node_count）
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"cluster_id":            "${local.k8s_cluster_id}",
					"vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"instance_types":        []string{"${local.default_instance_type_id}"},
					"node_count":            "1",
					"password":              "${random_password.password.0.result}",
					"system_disk_category":  "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"data_disks": []map[string]string{
						{
							"size":     "100",
							"category": "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
						},
					},
					"tags": map[string]interface{}{"Created": "TF", "Foo": "Bar"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"cluster_id":            CHECKSET,
						"vswitch_ids.#":         "1",
						"instance_types.#":      "1",
						"node_count":            "1",
						"system_disk_category":  CHECKSET,
						"system_disk_size":      "40",
						"install_cloud_monitor": "false",
						"data_disks.#":          "1",
						"data_disks.0.size":     "100",
						"data_disks.0.category": CHECKSET,
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.Foo":              "Bar",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_types": []string{"${local.default_instance_type_id}", "${local.update_instance_type_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "2",
					}),
				),
			},
			// Step 2: 导入验证
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Step 3: 扩容节点并修改系统盘和数据盘
			{
				Config: testAccConfig(map[string]interface{}{
					"node_count":       "2",
					"system_disk_size": "80",
					"data_disks": []map[string]string{
						{
							"size":     "40",
							"category": "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_count":            "2",
						"system_disk_size":      "80",
						"data_disks.#":          "1",
						"data_disks.0.size":     "40",
						"data_disks.0.category": CHECKSET,
					}),
				),
			},
			// Step 4: 缩容节点
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

func TestAccAlibabacloudStackCSKubernetesNodePool_AutoScaling(t *testing.T) {
	var v *NodePoolAlone

	resourceId := "alibabacloudstack_cs_kubernetes_node_pool.autoscaling"
	ra := resourceAttrInit(resourceId, csdKubernetesNodePoolBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccNodePoolAuto-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSNodePoolConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 1: 创建自动扩缩容节点池（使用scaling_config）
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"cluster_id":            "${local.k8s_cluster_id}",
					"vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"instance_types":        []string{"${local.default_instance_type_id}"},
					"image_id":              "${data.alibabacloudstack_images.default.images.0.id}",
					"key_name":              "${alibabacloudstack_ecs_keypair.default.key_name}",
					"system_disk_category":  "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"system_disk_size":      "40",
					"install_cloud_monitor": "false",
					"platform":              "Custom",
					"scaling_policy":        "release",
					"scaling_config": []map[string]string{
						{
							"min_size":      "1",
							"max_size":      "10",
							"type":          "cpu",
							"is_bond_eip":   "true",
							"eip_bandwidth": "5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"cluster_id":                     CHECKSET,
						"vswitch_ids.#":                  "1",
						"instance_types.#":               "1",
						"key_name":                       CHECKSET,
						"system_disk_category":           CHECKSET,
						"system_disk_size":               "40",
						"install_cloud_monitor":          "false",
						"platform":                       "Custom",
						"scaling_policy":                 "release",
						"scaling_config.#":               "1",
						"scaling_config.0.min_size":      "1",
						"scaling_config.0.max_size":      "10",
						"scaling_config.0.type":          "cpu",
						"scaling_config.0.is_bond_eip":   "true",
						"scaling_config.0.eip_bandwidth": "5",
					}),
				),
			},
			// Step 2: 导入验证
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Step 3: 更新自动扩缩容配置（调整max_size）
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_policy": "release",
					"scaling_config": []map[string]string{
						{
							"min_size":      "1",
							"max_size":      "20",
							"type":          "cpu",
							"is_bond_eip":   "true",
							"eip_bandwidth": "5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_policy":                 "release",
						"scaling_config.#":               "1",
						"scaling_config.0.min_size":      "1",
						"scaling_config.0.max_size":      "20",
						"scaling_config.0.type":          "cpu",
						"scaling_config.0.is_bond_eip":   "true",
						"scaling_config.0.eip_bandwidth": "5",
					}),
				),
			},
			// Step 4: 修改EIP绑定配置
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_config": []map[string]string{
						{
							"min_size":      "1",
							"max_size":      "20",
							"type":          "cpu",
							"is_bond_eip":   "false",
							"eip_bandwidth": "5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_config.#":               "1",
						"scaling_config.0.min_size":      "1",
						"scaling_config.0.max_size":      "20",
						"scaling_config.0.type":          "cpu",
						"scaling_config.0.is_bond_eip":   "false",
						"scaling_config.0.eip_bandwidth": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_types": []string{"${local.default_instance_type_id}", "${local.update_instance_type_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "2",
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

%s

locals {
  update_instance_type_id = coalesce(
    try(local.filtered_default[0].ids[1], null),
    try(local.fallback_all[1], null),
    "no-update-instance-type"
  )
}

resource "alibabacloudstack_ecs_keypair" "default" {
  key_name = var.name
}

`, name, DataAlibabacloudstackImages, AckK8sCommonTestCase())
}
