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
					"name":                                "${var.name}",
					"version":                             "1.34.1-aliyun.1",
					"os_type":                             "linux",
					"platform":                            "AliyunLinux",
					"timeout_mins":                        "60",
					"vpc_id":                              "${alibabacloudstack_vpc_vpc.default.id}",
					"master_count":                        "3",
					"master_disk_category":                "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"image_id":                            "${data.alibabacloudstack_images.default.images.0.id}",
					"master_disk_size":                    "40",
					"master_instance_types":               []string{"${local.default_instance_type_id}", "${local.default_instance_type_id}", "${local.default_instance_type_id}"},
					"master_vswitch_ids":                  []string{"${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"},
					"num_of_nodes":                        "1",
					"worker_disk_category":                "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"worker_disk_size":                    "40",
					"worker_instance_types":               []string{"${local.default_instance_type_id}"},
					"worker_vswitch_ids":                  []string{"${alibabacloudstack_vpc_vswitch.default.id}"},
					"enable_ssh":                          "${var.enable_ssh}",
					"password":                            "${random_password.password.0.result}",
					"delete_protection":                   "false",
					"pod_cidr":                            "${var.pod_cidr}",
					"service_cidr":                        "${var.service_cidr}",
					"node_cidr_mask":                      "${var.node_cidr_mask}",
					"is_enterprise_security_group":        "true",
					"new_nat_gateway":                     "false",
					"slb_internet_enabled":                "false",
					"proxy_mode":                          "ipvs",
					"master_storage_set_id":               "${alibabacloudstack_ecs_ebs_storage_set.master.storage_set_id}",
					"master_storage_set_partition_number": "3",
					"master_disk_encrypt_algorithm":       "aes-256",
					"master_disk_kms_key_id":              "${alibabacloudstack_kms_key.default.id}",
					"master_disk_encrypted":               "true",
					"worker_disk_kms_key_id":              "${alibabacloudstack_kms_key.default.id}",
					"worker_disk_encrypted":               "true",
					"worker_disk_encrypt_algorithm":       "aes-256",
					"worker_storage_set_id":               "${alibabacloudstack_ecs_ebs_storage_set.worker.storage_set_id}",
					"worker_storage_set_partition_number": "3",
					"worker_data_disks": []map[string]string{{
						"size":       "40",
						"encrypted":  "true",
						"category":   "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
						"kms_key_id": "${alibabacloudstack_kms_key.default.id}",
					}},
					"tags": map[string]string{
						"Created":                   "TF",
						"For":                       "acceptance test",
						"AAAAAAAAAAAAAAAAAAAAAAAAA": "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
					},
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
						"worker_nodes.#":          "1",
						"tags.%":                  "3",
						"tags.Created":            "TF",
						"tags.For":                "acceptance test",
						"proxy_mode":              "ipvs",
						"node_cidr_mask":          "26",
						"worker_data_disks.#":     "1",
					}),
				),
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
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"enable_ssh", "addons", "cluster_type", "image_id", "is_enterprise_security_group", "master_count", "master_disk_category", "master_disk_size", "master_instance_types", "master_vswitch_ids", "node_port_range", "os_type", "platform",
					"timeout_mins", "new_nat_gateway", "password", "slb_internet_enabled", "master_storage_set_id", "master_storage_set_partition_number", "node_cidr_mask", "proxy_mode", "service_cidr", "worker_data_disks", "worker_storage_set_id", "worker_storage_set_partition_number"},
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
					"version":                      "1.34.1-aliyun.1",
					"os_type":                      "linux",
					"platform":                     "AliyunLinux",
					"timeout_mins":                 "60",
					"vpc_id":                       "${alibabacloudstack_vpc_vpc.default.id}",
					"master_count":                 "3",
					"master_disk_category":         "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"image_id":                     "${data.alibabacloudstack_images.default.images.0.id}",
					"master_disk_size":             "40",
					"master_instance_types":        []string{"${local.default_instance_type_id}", "${local.default_instance_type_id}", "${local.default_instance_type_id}"},
					"master_vswitch_ids":           []string{"${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"},
					"num_of_nodes":                 "1",
					"worker_disk_category":         "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}",
					"worker_disk_size":             "40",
					"worker_instance_types":        []string{"${local.default_instance_type_id}"},
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
						"worker_disk_size":        "40",
						"key_name":                CHECKSET,
						"cpu_policy":              "none",
						"runtime.#":               "1",
						"runtime.0.name":          "containerd",
						"runtime.0.version":       "2.1.5",
						"master_nodes.#":          "3",
						"worker_nodes.#":          "1",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_ssh", "addons", "cluster_type", "image_id", "is_enterprise_security_group", "master_count", "master_disk_category", "master_disk_size", "master_instance_types", "master_vswitch_ids", "node_port_range", "os_type", "platform", "timeout_mins"},
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

// variable "image_id" {
//   default     = "centos_7_9_x64_20G_alibase_20220322.vhd"
// }

%s

# leave it to empty then terraform will create several vswitches

resource "alibabacloudstack_ecs_ebs_storage_set" "master" {
  storage_set_name = "tf-testAcc_storage_set4148"
  maxpartition_number = "3"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_ebs_storage_set" "worker" {
  storage_set_name = "tf-testAcc_storage_set4148"
  maxpartition_number = "3"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

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
	description = "${var.name}"
	protection_level =     "SOFTWARE"
	pending_window_in_days = "7"
}

`, name, SecurityGroupCommonTestCase, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages, RandomPasswordTestCase(12, 1))
}

var CsK8sMap = map[string]string{}
