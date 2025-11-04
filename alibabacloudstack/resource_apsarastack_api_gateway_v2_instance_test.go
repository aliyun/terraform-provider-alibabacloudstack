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

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		ExternalProviders: testAccExternalProviders,
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

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		ExternalProviders: testAccExternalProviders,
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

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":            "${var.name}",
					"node_number":              "1",
					"instance_class":           "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}",
					"broker_engine_type":       "HIGRESS",
					"deploy_mode":              "k8s",
					"deploy_cluster_code":      "${local.k8s_cluster_id}",
					"deploy_cluster_namespace": "${var.name}-namespace",
					"ingress_class_name":       "${var.name}-class",
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
						"deploy_cluster_namespace": fmt.Sprintf("%s-namespace", name),
						"ingress_class_name":       fmt.Sprintf("%s-class", name),
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

`, name, AckK8sCommonTestCase())
}
