package alibabacloudstack

import (
	"fmt"
	"os"

	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCsK8s_Basic(t *testing.T) {
	var v *ClusterObject
	resourceId := "alibabacloudstack_cs_kubernetes.k8s"
	ra := resourceAttrInit(resourceId, CsK8sMap)
	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCsKubernetes")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccK8s%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCsK8sConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{

					// "tags": map[string]string{
					// 	"Created": "TF",
					// 	"For":     "acceptance test",
					// },
					"runtime": []map[string]interface{}{
						{"name": "containerd", "version": "2.1.5"},
					},
					"addons": []map[string]interface{}{
						{
							"name": "flannel",
						},
					},
					"name":                 "${var.name}",
					"version":              "1.36.2-aliyun.1",
					"os_type":              "linux",
					"platform":             "AliyunLinux",
					"timeout_mins":         "60",
					"vpc_id":               "${alibabacloudstack_vpc_vpc.default.id}",
					"master_count":         "3",
					"master_disk_category": "${local.disk_category}",
					"image_id":             "${data.alibabacloudstack_images.ack.images.0.id}",
					"work_image_id":        "${data.alibabacloudstack_images.ack.images.0.id}",
					"master_data_disks": []map[string]interface{}{
						{
							"size":       100,
							"category":   "cloud_pperf",
							"encrypted":  true,
							"kms_key_id": "${alibabacloudstack_kms_key.default.id}",
						},
					},
					// "image_type":                          "AliyunLinux4ContainerOptimized",
					"master_disk_size":             "120",
					"master_instance_types":        []string{"${local.ack_instance_type}", "${local.ack_instance_type}", "${local.ack_instance_type}"},
					"master_vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"},
					"num_of_nodes":                 "1",
					"worker_disk_category":         "${local.disk_category}",
					"worker_disk_size":             "120",
					"worker_instance_types":        []string{"${local.ack_instance_type}"},
					"worker_vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"enable_ssh":                   "${var.enable_ssh}",
					"password":                     "${random_password.password.0.result}",
					"delete_protection":            "false",
					"pod_cidr":                     "${var.pod_cidr}",
					"service_cidr":                 "${var.service_cidr}",
					"node_cidr_mask":               "${var.node_cidr_mask}",
					"is_enterprise_security_group": "true",
					"new_nat_gateway":              "false",
					"slb_internet_enabled":         "false",
					"proxy_mode":                   "ipvs",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                           name,
						"image_id":                       CHECKSET,
						"work_image_id":                  CHECKSET,
						"work_image_type":                CHECKSET,
						"nodepool_id":                    CHECKSET,
						"vpc_id":                         CHECKSET,
						"pod_cidr":                       CHECKSET,
						"service_cidr":                   CHECKSET,
						"version":                        CHECKSET,
						"delete_protection":              "false",
						"enable_ssh":                     "true",
						"worker_instance_types.#":        "1",
						"worker_vswitch_ids.#":           "1",
						"master_disk_category":           CHECKSET,
						"master_data_disks.#":            "1",
						"master_data_disks.0.size":       "100",
						"master_data_disks.0.category":   "cloud_pperf",
						"master_data_disks.0.encrypted":  "true",
						"master_data_disks.0.kms_key_id": CHECKSET,
						"worker_disk_category":           CHECKSET,
						"worker_disk_size":               "120",
						"cpu_policy":                     "none",
						"runtime.#":                      "1",
						"runtime.0.name":                 "containerd",
						"runtime.0.version":              "2.1.5",
						"master_nodes.#":                 "3",
						// "worker_nodes.#":          "1",
						"proxy_mode":                   "ipvs",
						"node_cidr_mask":               "24",
						"os_type":                      "linux",
						"platform":                     "AliyunLinux",
						"timeout_mins":                 "60",
						"master_count":                 "3",
						"master_disk_size":             "120",
						"master_instance_types.#":      "3",
						"master_vswitch_ids.#":         "3",
						"num_of_nodes":                 "1",
						"is_enterprise_security_group": "true",
						"new_nat_gateway":              "false",
						"slb_internet_enabled":         "false",
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
					"new_nat_gateway", "slb_internet_enabled", "timeout_mins",
					"master_disk_encrypt_algorithm", "worker_disk_encrypt_algorithm",
					"master_disk_kms_key_id",
					"work_image_id", "work_image_type",
					"master_data_disks",
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
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackCsK8sSecurityGroup(t *testing.T) {
	var v *ClusterObject
	resourceId := "alibabacloudstack_cs_kubernetes.k8s"
	ra := resourceAttrInit(resourceId, CsK8sMap)
	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCsKubernetes")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCsK8sConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCsK8sConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{

					// "tags": map[string]string{
					// 	"Created": "TF",
					// 	"For":     "acceptance test",
					// },
					"runtime": []map[string]interface{}{
						{"name": "containerd", "version": "2.1.5"},
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
						{
							"name":   "logtail-ds",
							"config": "{\\\"IngressDashboardEnabled\\\":\\\"true\\\"}",
						},
						{
							"name":   "ack-node-problem-detector",
							"config": "{\\\"sls_project_name\\\":\\\"\\\"}",
						},
						{
							"name":   "nginx-ingress-controller",
							"config": "{\\\"IngressSlbNetworkType\\\":\\\"intranet\\\"}",
						},
					},
					"name":                         "${var.name}",
					"version":                      "1.36.2-aliyun.1",
					"os_type":                      "linux",
					"platform":                     "AliyunLinux",
					"timeout_mins":                 "60",
					"vpc_id":                       "${alibabacloudstack_vpc_vpc.default.id}",
					"master_count":                 "3",
					"master_disk_category":         "${local.disk_category}",
					"image_id":                     "${data.alibabacloudstack_images.ack.images.0.id}",
					"image_type":                   "AliyunLinux4ContainerOptimized",
					"master_disk_size":             "120",
					"master_instance_types":        []string{"${local.ack_instance_type}", "${local.ack_instance_type}", "${local.ack_instance_type}"},
					"master_vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"},
					"num_of_nodes":                 "1",
					"worker_disk_category":         "${local.disk_category}",
					"worker_disk_size":             "120",
					"worker_instance_types":        []string{"${local.ack_instance_type}"},
					"worker_vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"security_group_id":            "${alibabacloudstack_ecs_securitygroup.default.id}",
					"is_enterprise_security_group": "false",
					"enable_ssh":                   "${var.enable_ssh}",
					"key_name":                     "${alibabacloudstack_ecs_keypair.default.key_name}",
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
						"worker_disk_size":        "120",
						"key_name":                CHECKSET,
						"cpu_policy":              "none",
						"runtime.#":               "1",
						"runtime.0.name":          "containerd",
						"runtime.0.version":       "2.1.5",
						"master_nodes.#":          "3",
						"proxy_mode":              "ipvs",
						"node_cidr_mask":          "24",
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
					"new_nat_gateway", "slb_internet_enabled", "timeout_mins", "security_group_id",
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
					}),
				),
			},
		},
	})
}

func resourceCsK8sConfigDependence(name string) string {
	ack_instance_type := os.Getenv("ALIBABACLOUDSTACK_ACK_INSTANCE_TYPE")
	var ack_instance_type_expr string
	if ack_instance_type != "" {
		ack_instance_type_expr = fmt.Sprintf(`"%s"`, ack_instance_type)
	} else {
		ack_instance_type_expr = `data.alibabacloudstack_instance_types.ack.instance_types[0].id`
	}
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

variable "k8s_number" {
  description = "The number of kubernetes cluster."
  default     = 1
}

// variable "image_id" {
//   default     = "centos_7_9_x64_20G_alibase_20220322.vhd"
// }

%s

# leave it to empty then terraform will create several vswitches

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
  default     = 24
}

variable "enable_ssh" {
  description = "Enable login to the node through SSH."
  default     = true
}

%s

%s

resource "alibabacloudstack_ecs_keypair" "default" {
  key_name = var.name
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

resource "alibabacloudstack_kms_key" "default" {
	description 			= "${var.name}"
	protection_level 		= "SOFTWARE"
	pending_window_in_days 	= "7"
	rotation_interval 		= "7d"
	automatic_rotation 		= "Enabled"
}

data "alibabacloudstack_instance_types" "ack" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 2
  memory_size       = 4
  sorted_by         = "Memory"
}

data "alibabacloudstack_images" "ack" {
  name_regex  = "^aliyun_4_x86_64_20G_container_optimized.*"
  most_recent = true
  owners      = "system"
}

locals {
    disk_category = "cloud_pperf"
    disk_encrypted = contains(["cloud_pperf", "cloud_sperf"], local.disk_category) ? true : false
	disk_kmskey = local.disk_encrypted? alibabacloudstack_kms_key.default.id:""
	ack_instance_type = %s
}

`, name, SecurityGroupCommonTestCase, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages, RandomPasswordTestCase(12, 1), ack_instance_type_expr)
}

var CsK8sMap = map[string]string{}
