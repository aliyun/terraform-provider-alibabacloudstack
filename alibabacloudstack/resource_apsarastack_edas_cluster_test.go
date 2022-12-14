package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_edas_cluster", &resource.Sweeper{
		Name: "alibabacloudstack_edas_cluster",
		F:    testSweepEdasCluster,
		Dependencies: []string{
			"alibabacloudstack_edas_application",
		},
	})
}

func testSweepEdasCluster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AlibabacloudStack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	clusterListRq := edas.CreateListClusterRequest()
	clusterListRq.RegionId = region

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListCluster(clusterListRq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve edas cluster in service list: %s", err)
		return nil
	}

	listClusterResponse, _ := raw.(*edas.ListClusterResponse)
	if listClusterResponse.Code != 200 {
		log.Printf("[ERROR] Failed to retrieve edas cluster in service list: %s", listClusterResponse)
		return nil
	}

	for _, v := range listClusterResponse.ClusterList.Cluster {
		name := v.ClusterName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}

		if skip {
			log.Printf("[INFO] Skipping edas cluster: %s", name)
			continue
		}
		log.Printf("[INFO] delete edas cluster: %s", name)

		deleteClusterRq := edas.CreateDeleteClusterRequest()
		deleteClusterRq.RegionId = region
		deleteClusterRq.ClusterId = v.ClusterId

		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteCluster(deleteClusterRq)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(deleteClusterRq.GetActionName(), raw, deleteClusterRq.RoaRequest, deleteClusterRq)
			rsp := raw.(*edas.DeleteClusterResponse)
			if rsp.Code == 601 && strings.Contains(rsp.Message, "Operation cannot be processed because there are running instances.") {
				err = Error("Operation cannot be processed because there are running instances.")
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete edas cluster (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAlibabacloudStackEdasCluster_basic(t *testing.T) {
	var v *edas.Cluster
	resourceId := "alibabacloudstack_edas_cluster.name1"
	ra := resourceAttrInit(resourceId, edasClusterBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasclusterbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": "${var.name}",
					"cluster_type": "2",
					"network_mode": "2",
					"vpc_id":       "${alibabacloudstack_vpc.default.id}",
					//"region_id":    "cn-qingdao-apsara-d01",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": name,
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
					"cluster_name": fmt.Sprintf("tf-testacc-edasclusterchange%v", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": fmt.Sprintf("tf-testacc-edasclusterchange%v", rand)}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_type": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_type": "2"}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackEdasCluster_multi(t *testing.T) {
	var v *edas.Cluster
	resourceId := "alibabacloudstack_edas_cluster.default.1"
	ra := resourceAttrInit(resourceId, edasClusterBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasclustermulti%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":        "2",
					"cluster_name": "${var.name}-${count.index}",
					"cluster_type": "2",
					"network_mode": "2",
					"vpc_id":       "${alibabacloudstack_vpc.default.id}",
					//"region_id":    "cn-qingdao-apsara-d01",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var edasClusterBasicMap = map[string]string{
	"cluster_name": CHECKSET,
	"cluster_type": CHECKSET,
	"network_mode": CHECKSET,
}

func testAccCheckEdasClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_edas_cluster" {
			continue
		}

		// Try to find the cluster
		clusterId := rs.Primary.ID
		regionId := client.RegionId

		request := edas.CreateGetClusterRequest()
		request.RegionId = regionId
		request.ClusterId = clusterId
		request.Headers["x-ascm-product-name"] = "Edas"
		request.Headers["x-acs-organizationid"] = client.Department
		request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.GetCluster(request)
		})

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		rsp := raw.(*edas.GetClusterResponse)
		if rsp.Cluster.ClusterId != "" {
			return fmt.Errorf("cluster %s still exist", rsp.Cluster.ClusterId)
		}
	}

	return nil
}

func resourceEdasClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}
	resource "alibabacloudstack_vpc" "default" {
	  cidr_block = "172.16.0.0/16"
	  name = "${var.name}"
	}
		`, name)
}
