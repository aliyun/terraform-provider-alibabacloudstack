package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_api_gateway_v2_k8s_cluster", &resource.Sweeper{
		Name: "alibabacloudstack_api_gateway_v2_k8s_cluster",
		F:    testSweepApiGatewayV2K8sCluster,
	})
}

func testSweepApiGatewayV2K8sCluster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	apigatewayv2Service := ApiGateWayV2Service{client}

	prefixes := []string{
		"tf_testAccApiGatewayv2cluster",
		"tf-testacc",
	}

	reqBody := map[string]interface{}{
		"current": 1,
		"size":    100,
	}

	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListClusters", "/k8s/listClusters", nil, nil, reqBody)
	if err != nil {
		log.Printf("[ERROR] Failed to list api_gateway_v2_k8s_cluster: %s", err)
		return nil
	}

	records, ok := response["data"].(map[string]interface{})["records"].([]interface{})
	if !ok {
		return nil
	}

	for _, record := range records {
		item, ok := record.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := item["k8sClusterName"].(string)
		k8sClusterCode, _ := item["k8sClusterCode"].(string)

		if name == "" || k8sClusterCode == "" {
			continue
		}

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}

		if skip {
			continue
		}

		log.Printf("[INFO] Deleting api_gateway_v2_k8s_cluster: %s (%s)", name, k8sClusterCode)

		deleteReqBody := map[string]interface{}{
			"k8sClusterCode": k8sClusterCode,
		}

		_, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CancelImportCluster", "/k8s/cancelImportCluster", nil, nil, deleteReqBody)
		if err != nil {
			log.Printf("[ERROR] Failed to delete api_gateway_v2_k8s_cluster (%s): %s", k8sClusterCode, err)
		}
	}

	_ = apigatewayv2Service
	return nil
}

func TestAccAlibabacloudStackApigatewayv2K8sCluster_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_k8s_cluster.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApigatewayv2K8sCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayv2cluster_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, APIGateWayV2InstanceClusterDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cs_cluster_id":    "${local.k8s_cluster_id}",
					"k8s_cluster_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cs_cluster_id":    CHECKSET,
						"k8s_cluster_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_content"},
			},
		},
	})
}

func TestAccAlibabacloudStackApigatewayv2K8sCluster_selfBuild(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_k8s_cluster.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApigatewayv2K8sCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayv2cluster_sb_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, APIGateWayV2InstanceClusterDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"config_content":   "${data.alibabacloudstack_cs_kubernetes_clusters_kubeconfig.k8s_clusters_kubeconfig.kubeconfig}",
					"k8s_cluster_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cs_cluster_id":    NOSET,
						"k8s_cluster_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_content"},
			},
		},
	})
}

func APIGateWayV2InstanceClusterDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

data "alibabacloudstack_cs_kubernetes_clusters_kubeconfig" "k8s_clusters_kubeconfig" {
	cluster_id = local.k8s_cluster_id
  }

`, name, AckK8sCommonTestCase())
}
