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
					"node_number":        "2",
					"instance_class":     "mini",
					"broker_engine_type": "SCG",
					"deploy_mode":        "edas",
					"edas_app_infos": []map[string]interface{}{
						{
							"edas_namespace_id": "cn-wulan-env17e-d01",
							"edas_k8s_id":       "a0f1b51f-ca84-4131-8e17-f1439e8e7c36",
							"k8s_namespace":     "default",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":      name,
						"node_number":        "2",
						"instance_class":     "mini",
						"broker_engine_type": "SCG",
						"deploy_mode":        "edas",
						"edas_app_infos.#":   "1",
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
					"instance_name":            "${var.name}",
					"node_number":              "1",
					"instance_class":           "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}",
					"broker_engine_type":       "HIGRESS",
					"deploy_mode":              "k8s",
					"deploy_cluster_code":      "fb63cd07-1272-40e9-a013-787a801eaad4",
					"deploy_cluster_namespace": "tftest-namespace1",
					"ingress_class_name":       "tftest-namespace2",
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
						"deploy_cluster_code":      "fb63cd07-1272-40e9-a013-787a801eaad4",
						"deploy_cluster_namespace": "tftest-namespace1",
						"ingress_class_name":       "tftest-namespace2",
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
