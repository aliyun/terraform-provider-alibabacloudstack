package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccCheckCsK8sDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_cs_kubernetes" {
			continue
		}
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
	var v *KubernetesClusterDetail
	resourceId := "alibabacloudstack_cs_kubernetes.k8s"
	ra := resourceAttrInit(resourceId, CsK8sMap)
	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccK8s%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCsK8sConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccCheckCsK8sDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": []map[string]interface{}{
						{"name": "containerd", "version": "2.1.5"},
					},
					"addons": []map[string]interface{}{
						{"name": "flannel"},
						{"name": "csi-plugin"},
						{"name": "csi-provisioner"},
						{"name": "nginx-ingress-controller"},
					},
					"name":                 "${var.name}",
					"version":              "1.34.1-aliyun.1",
					"os_type":              "linux",
					"platform":             "AliyunLinux3",
					"timeout_mins":         "60",
					"vpc_id":               "${alibabacloudstack_vpc_vpc.default.id}",
					"master_count":         "3",
					"master_disk_category": "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					// "image_id":                     "${data.alibabacloudstack_images.aliyunlinux3.images.0.id}",
					"master_disk_size":             "120",
					"master_instance_types":        []string{"${local.ack_instance_type_id}", "${local.ack_instance_type_id}", "${local.ack_instance_type_id}"},
					"master_vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"},
					"num_of_nodes":                 "1",
					"worker_disk_category":         "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"worker_disk_size":             "120",
					"worker_instance_types":        []string{"${local.ack_worker_instance_type_id}"},
					"worker_vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"enable_ssh":                   "${var.enable_ssh}",
					"password":                     "${random_password.password.result}",
					"delete_protection":            "false",
					"pod_cidr":                     "${var.pod_cidr}",
					"service_cidr":                 "${var.service_cidr}",
					"node_cidr_mask":               "${var.node_cidr_mask}",
					"is_enterprise_security_group": "true",
					"new_nat_gateway":              "false",
					"slb_internet_enabled":         "false",
					"proxy_mode":                   "ipvs",
					// "master_storage_set_id":               "${alibabacloudstack_ecs_ebs_storage_set.master.storage_set_id}",
					// "master_storage_set_partition_number": "3",
					// "worker_storage_set_id":               "${alibabacloudstack_ecs_ebs_storage_set.worker.storage_set_id}",
					// "worker_storage_set_partition_number": "3",
					"worker_data_disks": []map[string]string{{
						"size":     "120",
						"category": "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					}},
					// "tags": map[string]string{
					// 	"Created": "TF",
					// 	"For":     "acceptance test",
					// },
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                         name,
						"nodepool_id":                  CHECKSET,
						"vpc_id":                       CHECKSET,
						"pod_cidr":                     CHECKSET,
						"service_cidr":                 CHECKSET,
						"version":                      CHECKSET,
						"delete_protection":            "false",
						"worker_instance_types.#":      "1",
						"worker_vswitch_ids.#":         "1",
						"worker_disk_category":         CHECKSET,
						"worker_disk_size":             "120",
						"cpu_policy":                   "none",
						"runtime.#":                    "1",
						"runtime.0.name":               "containerd",
						"runtime.0.version":            "2.1.5",
						"master_nodes.#":               "3",
						"proxy_mode":                   "ipvs",
						"node_cidr_mask":               "26",
						"worker_data_disks.#":          "1",
						"os_type":                      "linux",
						"platform":                     "AliyunLinux3",
						"timeout_mins":                 "60",
						"master_count":                 "3",
						"master_disk_size":             "120",
						"master_instance_types.#":      "3",
						"master_vswitch_ids.#":         "3",
						"num_of_nodes":                 "1",
						"is_enterprise_security_group": "true",
						"new_nat_gateway":              "false",
						"slb_internet_enabled":         "false",
						// "master_storage_set_partition_number": "3",
						// "worker_storage_set_partition_number": "3",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "addons", "image_type", "is_enterprise_security_group", "os_type",
					"master_storage_set_id", "master_storage_set_partition_number", "master_vswitch_ids",
					"new_nat_gateway", "slb_internet_enabled", "timeout_mins", "worker_storage_set_id",
					"worker_storage_set_partition_number",
					"enable_ssh", "master_disk_category", "master_disk_size", "node_cidr_mask",
					"node_port_range", "proxy_mode", "service_cidr",
				},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"num_of_nodes": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"num_of_nodes": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"num_of_nodes": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"num_of_nodes":            "0",
						"worker_vswitch_ids.#":    REMOVEKEY,
						"worker_disk_size":        REMOVEKEY,
						"nodepool_id":             REMOVEKEY,
						"worker_disk_category":    REMOVEKEY,
						"worker_nodes.#":          REMOVEKEY,
						"worker_instance_types.#": REMOVEKEY,
						"worker_data_disks.#":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackCsK8sSecurityGroup(t *testing.T) {
	var v *KubernetesClusterDetail
	resourceId := "alibabacloudstack_cs_kubernetes.k8s"
	ra := resourceAttrInit(resourceId, CsK8sMap)
	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCsK8sConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCsK8sConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccCheckCsK8sDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"runtime": []map[string]interface{}{
						{"name": "containerd", "version": "2.1.5"},
					},
					"addons": []map[string]interface{}{
						{"name": "flannel"},
						{"name": "csi-plugin"},
						{"name": "csi-provisioner"},
						{"name": "nginx-ingress-controller"},
					},
					"name":                 "${var.name}",
					"version":              "1.34.1-aliyun.1",
					"os_type":              "linux",
					"platform":             "AliyunLinux",
					"timeout_mins":         "60",
					"vpc_id":               "${alibabacloudstack_vpc_vpc.default.id}",
					"master_count":         "3",
					"master_disk_category": "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					// "image_id":                     "${data.alibabacloudstack_images.aliyunlinux3.images.0.id}",
					"image_type":                   "AliyunLinux3",
					"master_disk_size":             "40",
					"master_instance_types":        []string{"${local.ack_instance_type_id}", "${local.ack_instance_type_id}", "${local.ack_instance_type_id}"},
					"master_vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"},
					"num_of_nodes":                 "1",
					"worker_disk_category":         "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"worker_disk_size":             "40",
					"worker_instance_types":        []string{"${local.ack_worker_instance_type_id}"},
					"worker_vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"security_group_id":            "${alibabacloudstack_ecs_securitygroup.default.id}",
					"is_enterprise_security_group": "false",
					"enable_ssh":                   "${var.enable_ssh}",
					"password":                     "${random_password.password.result}",
					"delete_protection":            "false",
					"pod_cidr":                     "${var.pod_cidr}",
					"service_cidr":                 "${var.service_cidr}",
					"node_cidr_mask":               "${var.node_cidr_mask}",
					"new_nat_gateway":              "false",
					"slb_internet_enabled":         "false",
					"proxy_mode":                   "ipvs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                    name,
						"nodepool_id":             CHECKSET,
						"vpc_id":                  CHECKSET,
						"pod_cidr":                CHECKSET,
						"service_cidr":            CHECKSET,
						"version":                 CHECKSET,
						"delete_protection":       "false",
						"worker_instance_types.#": "1",
						"worker_vswitch_ids.#":    "1",
						"worker_disk_category":    CHECKSET,
						"worker_disk_size":        "40",
						"cpu_policy":              "none",
						"runtime.#":               "1",
						"runtime.0.name":          "containerd",
						"runtime.0.version":       "2.1.5",
						"master_nodes.#":          "3",
						"proxy_mode":              "ipvs",
						"node_cidr_mask":          "26",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"num_of_nodes": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"num_of_nodes": "3",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "addons", "image_type", "is_enterprise_security_group", "os_type",
					"master_storage_set_id", "master_storage_set_partition_number", "master_vswitch_ids",
					"new_nat_gateway", "slb_internet_enabled", "timeout_mins", "worker_storage_set_id",
					"worker_storage_set_partition_number", "security_group_id",
				},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"num_of_nodes": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"num_of_nodes":            "0",
						"worker_vswitch_ids.#":    REMOVEKEY,
						"worker_disk_size":        REMOVEKEY,
						"nodepool_id":             REMOVEKEY,
						"worker_disk_category":    REMOVEKEY,
						"worker_nodes.#":          REMOVEKEY,
						"worker_instance_types.#": REMOVEKEY,
						"worker_data_disks.#":     REMOVEKEY,
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

%s

variable "k8s_number" {
  description = "The number of kubernetes cluster."
  default     = 1
}

%s

%s

// resource "alibabacloudstack_ecs_ebs_storage_set" "master" {
//   storage_set_name = "tf-testAcc_storage_set4148"
//   maxpartition_number = "3"
//   zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
// }

// resource "alibabacloudstack_ecs_ebs_storage_set" "worker" {
//   storage_set_name = "tf-testAcc_storage_set4148"
//   maxpartition_number = "3"
//   zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
// }

variable "runtime" {
 default     = [
		{
			name = "containerd"
  			version = "2.1.5"
		}
	]
}

variable "new_nat_gateway" {
  description = "Whether to create a new nat gateway. In this template, a new nat gateway will create a nat gateway, eip and server snat entries."
  default     = "ture"
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

%s

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

data "alibabacloudstack_instance_types" "ack" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 4
  memory_size       = 8
  sorted_by         = "Memory"
}

data "alibabacloudstack_images" "aliyunlinux3" {
  name_regex  = "^aliyun_3"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "ack_master" {
  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
  kubernetes_node_role = "Master"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "ack_worker" {
  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
  kubernetes_node_role = "Worker"
  sorted_by            = "Memory"
}

locals {
    disk_category = data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0
    ack_instance_type_id        = try(data.alibabacloudstack_instance_types.ack_master.instance_types[0].id, "ecs.g7x-se-x25-c1m4.xlarge")
    ack_worker_instance_type_id = try(data.alibabacloudstack_instance_types.ack_worker.instance_types[0].id, "ecs.s6-k-c1m6.xlarge")
}

`, name, SecurityGroupCommonTestCase, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages, RandomPasswordTestCase(12))
}

var CsK8sMap = map[string]string{}
