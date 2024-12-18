package alibabacloudstack

import (
	"fmt"

	"github.com/denverdino/aliyungo/cs"

	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccCheckCsK8sDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_cs_kubernetes" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		csService := CsService{client}
		log.Printf("repo ID %s", rs.Primary.ID)
		_, err := csService.DescribeCsKubernetes(rs.Primary.ID)

		if err == nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
	}

	return nil
}

func TestAccAlibabacloudStackCsK8s_Basic(t *testing.T) {
	var v *cs.KubernetesClusterDetail
	resourceId := "alibabacloudstack_cs_kubernetes.k8s"
	ra := resourceAttrInit(resourceId, CsK8sMap)
	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCsK8sConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCsK8sConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCsK8sDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{

					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"runtime": []map[string]interface{}{
						{"name": "docker", "version": "19.03.15"},
					},
					"addons": []map[string]interface{}{
						{
							"name": "flannel",
						},
						{
							"name": "csi-plugin",
						},
						{
							"name": "csi-provisioner",
						},
						{
							"name": "nginx-ingress-controller",
						},
					},
					"name":                  "${var.name}",
					"version":               "1.20.11-aliyun.1",
					"os_type":               "linux",
					"platform":              "AliyunLinux",
					"timeout_mins":          "60",
					"vpc_id":                "${alibabacloudstack_vpc_vpc.default.id}",
					"master_count":          "3",
					"master_disk_category":  "cloud_ssd",
					"image_id":              "${var.image_id}",
					"master_disk_size":      "40",
					"master_instance_types": "${var.master_instance_types}",
					"master_vswitch_ids":    []string{"${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"},
					"num_of_nodes":          "1",
					"worker_disk_category":  "cloud_ssd",
					"worker_disk_size":      "40",
					"worker_instance_types": "${var.worker_instance_types}",
					"worker_vswitch_ids":    []string{"${alibabacloudstack_vpc_vswitch.default.id}"},

					"enable_ssh":                   "${var.enable_ssh}",
					"password":                     "${var.password}",
					"delete_protection":            "false",
					"pod_cidr":                     "${var.pod_cidr}",
					"service_cidr":                 "${var.service_cidr}",
					"node_cidr_mask":               "${var.node_cidr_mask}",
					"is_enterprise_security_group": "true",
					"new_nat_gateway":              "false",
					"slb_internet_enabled":         "false",
					"proxy_mode":                   "ipvs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"detail": "detail",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail": "detail",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"summary": "summary update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"summary": "summary update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_type": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_type": "PRIVATE",
					}),
				),
			},
		},
	})
}

func resourceCsK8sConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "k8s_number" {
  description = "The number of kubernetes cluster."
  default     = 1
}

variable "image_id" {
  default     = "centos_7_9_x64_20G_alibase_20220322.vhd"
}

%s

# leave it to empty then terraform will create several vswitches

variable "runtime" {
 default     = [
		{
			name = "containerd"
  			version = "1.5.13"
		}
	]
}

variable "new_nat_gateway" {
  description = "Whether to create a new nat gateway. In this template, a new nat gateway will create a nat gateway, eip and server snat entries."
  default     = "false"
}

# 3 masters is default settings,so choose three appropriate instance types in the availability zones above.
variable "master_instance_types" {
  description = "The ecs instance types used to launch master nodes."
   default     = ["ecs.n4v2.large","ecs.n4v2.large","ecs.n4v2.large"]
}

variable "worker_instance_types" {
  description = "The ecs instance types used to launch worker nodes."
  default     = ["ecs.n4v2.large"]
}

# options: between 24-28
variable "node_cidr_mask" {
  description = "The node cidr block to specific how many pods can run on single node."
  default     = 26
}

variable "enable_ssh" {
  description = "Enable login to the node through SSH."
  default     = true
}


variable "password" {
  description = "The password of ECS instance."
  default     = "Alibaba@1688"
}

variable "worker_number" {
  description = "The number of worker nodes in kubernetes cluster."
  default     = 3
}

# k8s_pod_cidr is only for flannel network
variable "pod_cidr" {
  description = "The kubernetes pod cidr block. It cannot be equals to vpc's or vswitch's and cannot be in them."
  default     = "172.24.0.0/16"
}

variable "service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  default     = "172.25.0.0/16"
}


`, name, VSwitchCommonTestCase)
}

var CsK8sMap = map[string]string{}
