package alibabacloudstack

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_edas_k8s_cluster", &resource.Sweeper{
		Name: "alibabacloudstack_edas_k8s_cluster",
		F:    testSweepEdasK8sCluster,
	})
}

func testSweepEdasK8sCluster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapErrorf(err, "error getting AlibabacloudStack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	clusterListRq := edas.CreateListClusterRequest()
	clusterListRq.RegionId = region

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListCluster(clusterListRq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve edas cluster in service list: %s", err)
	}

	listClusterResponse, _ := raw.(*edas.ListClusterResponse)
	if listClusterResponse.Code != 200 {
		log.Printf("[ERROR] Failed to retrieve edas cluster in service list: %s", listClusterResponse.Message)
		return errmsgs.WrapError(errmsgs.Error(listClusterResponse.Message))
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

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteCluster(deleteClusterRq)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(deleteClusterRq.GetActionName(), raw, deleteClusterRq.RoaRequest, deleteClusterRq)
			rsp := raw.(*edas.DeleteClusterResponse)
			if rsp.Code == 601 && strings.Contains(rsp.Message, "Operation cannot be processed because there are running instances.") {
				err = errmsgs.Error("Operation cannot be processed because there are running instances.")
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

func TestAccAlibabacloudStackEdasK8sCluster_basic(t *testing.T) {
	var v *EdasK8sCluster
	resourceId := "alibabacloudstack_edas_k8s_cluster.default"
	ra := resourceAttrInit(resourceId, edasK8sClusterBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	rand := getAccTestRandInt(0, 1000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edask8sclusterbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sClusterConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasK8sClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cs_cluster_id": "${alibabacloudstack_cs_kubernetes.default.id}",
					// "cs_cluster_id": "c89eeac401e7b43d985c6ac2b94ceee66",
					"namespace_id": region,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"cluster_name":  name,
						"cs_cluster_id": CHECKSET,
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

var edasK8sClusterBasicMap = map[string]string{
	"cluster_name":          CHECKSET,
	"cluster_type":          CHECKSET,
	"network_mode":          CHECKSET,
	"cluster_import_status": CHECKSET,
}

func testAccCheckEdasK8sClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_edas_k8s_cluster" {
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
			if errmsgs.NotFoundError(err) {
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

func resourceEdasK8sClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%s"
		}

		variable "cluster_addons" {
			description = "Addon components in kubernetes cluster"

			type = list(object({
				name   = string
				config = string
			}))

			default = [
					{
					"name": "arms-prometheus",
					"config":"",
				},
				{
					"name": "csi-plugin",
					"config":"",
				},
				{
					"name": "csi-provisioner",
					"config":"",
				},
				{
					"name": "logtail-ds",
					"config": "{\"IngressDashboardEnabled\":\"true\"}"
				},
				{
					"name": "ack-node-problem-detector",
					"config": "{\"sls_project_name\":\"\"}"
				},
				{
					"name": "nginx-ingress-controller",
					"config": "{\"IngressSlbNetworkType\":\"intranet\"}"
				}
			]
		}

		%s

		data "alibabacloudstack_images" "default" {
		  name_regex  = "^anolisos_"
		  most_recent = true
		  owners      = "system"
		}

		data "alibabacloudstack_instance_types" "default" {
		  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
		  cpu_core_count       = 2
		  memory_size          = 4
		}

		resource "alibabacloudstack_cs_kubernetes" "default" {
			name 						= var.name
			version 					= "1.30.1-aliyun.1"
			os_type 					= "linux"
			platform 					= "AliyunLinux"
			image_id                    = "${data.alibabacloudstack_images.default.images.0.id}"
			num_of_nodes 				= "3"
			master_count				= "3"
			master_vswitch_ids   		= ["${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"]
			master_instance_types 		= ["${data.alibabacloudstack_instance_types.default.instance_types.0.id}","${data.alibabacloudstack_instance_types.default.instance_types.0.id}","${data.alibabacloudstack_instance_types.default.instance_types.0.id}"]
			master_disk_category 		= "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
			vpc_id 					    = "${alibabacloudstack_vpc_vpc.default.id}"
			worker_instance_types 		= ["${data.alibabacloudstack_instance_types.default.instance_types.0.id}"]
			worker_vswitch_ids 		    = ["${alibabacloudstack_vpc_vswitch.default.id}"]
			worker_disk_category 		= "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
			password 					= "Test@12345"
			pod_cidr 					= "172.20.0.0/16"
			service_cidr 				= "172.21.0.0/20"
			worker_disk_size 			= "40"
			master_disk_size 			= "40"
			slb_internet_enabled 		= "true"
			runtime {
				name    = "containerd"
				version = "1.6.28"
			}
			dynamic "addons" {
				for_each = var.cluster_addons
				content {
					name   = lookup(addons.value, "name", var.cluster_addons)
					config = lookup(addons.value, "config", var.cluster_addons)
				}
			}
		}

		`, name, SecurityGroupCommonTestCase)
}
