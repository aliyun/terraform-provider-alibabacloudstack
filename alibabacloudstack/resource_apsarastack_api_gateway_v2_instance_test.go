package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAPIGateWayV2Instance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2Instace")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("testtf-apigw-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, APIGateWayV2InstanceEdasDepDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":      "${var.name}",
					"node_number":        "1",
					"instance_class":     "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}",
					"broker_engine_type": "SCG",
					"deploy_mode":        "edas",
					"edas_app_infos": []map[string]interface{}{
						{
							"edas_namespace": defaultRegionToTest,
							"k8s_cluster_id": "${local.edas_cluster_id}",
							"k8s_namespace":  "default",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"node_number":             "1",
						"instance_class":          "mini",
						"broker_engine_type":      "SCG",
						"deploy_mode":             "edas",
						"edas_app_infos.#":        "1",
						"edas_app_infos.0.app_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": fmt.Sprintf("%s_update", name),
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

func TestAccAlibabacloudStackAPIGateWayV2Instance_custom(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2Instace")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("testtf-apigw-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, APIGateWayV2InstanceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":      "${var.name}",
					"node_number":        "1",
					"instance_class":     "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}",
					"broker_engine_type": "SCG",
					"deploy_mode":        "custom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":                          name,
						"node_number":                            "1",
						"instance_class":                         "mini",
						"broker_engine_type":                     "SCG",
						"deploy_mode":                            "custom",
						"custom_deploy_config.%":                 "4",
						"custom_deploy_config.jarStr":            CHECKSET,
						"custom_deploy_config.serviceYamlStr":    CHECKSET,
						"custom_deploy_config.deploymentYamlStr": CHECKSET,
						"custom_deploy_config.wgetStr":           CHECKSET,
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

func TestAccAlibabacloudStackAPIGateWayV2Instance_HIGRESS(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2Instace")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("testtf-apigw-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AIGateWayV2InstanceK8sDepDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":            "${var.name}",
					"node_number":              "1",
					"instance_class":           "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}",
					"broker_engine_type":       "HIGRESS",
					"deploy_mode":              "k8s",
					"deploy_cluster_code":      "${alibabacloudstack_cs_kubernetes.default.id}",
					"deploy_cluster_namespace": "${var.name}_namespace",
					"ingress_class_name":       "${var.name}_class",
					"sls_enabled":              "true",
					"prometheus_enabled":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":            name,
						"node_number":              "1",
						"instance_class":           "mini",
						"broker_engine_type":       "HIGRESS",
						"deploy_mode":              "k8s",
						"deploy_cluster_code":      CHECKSET,
						"deploy_cluster_namespace": fmt.Sprintf("%s_namespace", name),
						"ingress_class_name":       fmt.Sprintf("%s_class", name),
						"sls_enabled":              "true",
						"prometheus_enabled":       "true",
					}),
				),
			},
			{
				// sls_enabled, prometheus_enabled  not read
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sls_enabled", "prometheus_enabled"},
			},
		},
	})
}

func APIGateWayV2InstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}
`, name)
}

func APIGateWayV2InstanceEdasDepDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}

%s

`, name, EdasClusterCommonTestCase())
}

func AIGateWayV2InstanceK8sDepDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}

%s

%s

resource "alibabacloudstack_cs_kubernetes" "default" {
	count						= local.create_count
	name						= var.name
	version						= "1.30.7-aliyun.1"
	os_type						= "linux"
	platform					= "AliyunLinux"
	num_of_nodes				= "3"
	master_count				= "3"
	master_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"]
	master_instance_types		= ["ecs.n4v2.large","ecs.n4v2.large","ecs.n4v2.large"]
	master_disk_category		= "cloud_ssd"
	vpc_id						= "${alibabacloudstack_vpc_vpc.default.id}"
	worker_instance_types		= ["ecs.n4v2.large"]
	worker_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}"]
	worker_disk_category		= "cloud_ssd"
	password					= random_password.password.0.result
	pod_cidr					= "172.20.0.0/16"
	service_cidr				= "172.21.0.0/20"
	worker_disk_size			= "40"
	master_disk_size			= "40"
	slb_internet_enabled		= "true"
	security_group_id			= alibabacloudstack_ecs_securitygroup.default.id
	runtime	 {
		name	= "containerd"
		version	= "1.6.28"
	}
}

`, name, SecurityGroupCommonTestCase, RandomPasswordTestCase(12, 1))
}
