package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackCSKubernetesClustersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alibabacloudstack_cs_kubernetes_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacckubernetes-%d", rand),
		dataSourceCSKubernetesClustersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			//"enable_details": "true",
			//"ids":            []string{"${alibabacloudstack_cs_kubernetes.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			//"enable_details": "true",
			//"ids":            []string{"${alibabacloudstack_cs_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			//"enable_details": "true",
			"name_regex": "${alibabacloudstack_cs_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			//"enable_details": "true",
			"name_regex": "${alibabacloudstack_cs_kubernetes.default.name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			//"enable_details": "true",
			//"ids":            []string{"${alibabacloudstack_cs_kubernetes.default.id}"},
			"name_regex": "${alibabacloudstack_cs_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			//"enable_details": "true",
			//"ids":            []string{"${alibabacloudstack_cs_kubernetes.default.id}"},
			"name_regex": "${alibabacloudstack_cs_kubernetes.default.name}-fake",
		}),
	}
	var existCSKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			//"ids.#":                                      "1",
			//"ids.0":                                      CHECKSET,
			//"names.#":                                    "1",
			//"names.0":                                    REGEXMATCH + fmt.Sprintf("tf-testacckubernetes-%d", rand),
			//"clusters.#":                                 "1",
			//"clusters.0.id":                              CHECKSET,
			//"clusters.0.name":                            REGEXMATCH + fmt.Sprintf("tf-testacckubernetes-%d", rand),
			//"clusters.0.availability_zone":               CHECKSET,
			//"clusters.0.security_group_id":               CHECKSET,
			//"clusters.0.nat_gateway_id":                  CHECKSET,
			//"clusters.0.vpc_id":                          CHECKSET,
			//"clusters.0.worker_numbers.#":                "1",
			//"clusters.0.worker_numbers.0":                "1",
			//"clusters.0.master_nodes.#":                  "3",
			//"clusters.0.worker_disk_category":            "cloud_efficiency",
			//"clusters.0.master_disk_size":                "50",
			//"clusters.0.master_disk_category":            "cloud_efficiency",
			//"clusters.0.worker_disk_size":                "40",
			//"clusters.0.connections.%":                   "4",
			//"clusters.0.connections.master_public_ip":    CHECKSET,
			//"clusters.0.connections.api_server_internet": CHECKSET,
			//"clusters.0.connections.api_server_intranet": CHECKSET,
			//"clusters.0.connections.service_domain":      CHECKSET,
		}
	}

	var fakeCSKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			//"ids.#":      "0",
			//"names.#":    "0",
			//"clusters.#": "0",
		}
	}

	var csKubernetesClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCSKubernetesClustersMapFunc,
		fakeMapFunc:  fakeCSKubernetesClustersMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
	}
	csKubernetesClustersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceCSKubernetesClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "worker_number" {
  description = "The number of worker nodes in kubernetes cluster."
  default     = 3
}

variable "vpc_id" {
  description = "Existing vpc id used to create several vswitches and other resources."
  default     = "vpc-bs1amyht0w3mdkumc2hf0"
}

variable "vswitch_ids" {
 description = "List of existing vswitch id."
 type        = list(string)
 default     = ["vsw-bs1xvmkkekuy8i7zlhv3j","vsw-bs1xvmkkekuy8i7zlhv3j","vsw-bs1xvmkkekuy8i7zlhv3j"]
}

variable "master_instance_types" {
  description = "The ecs instance types used to launch master nodes."
  default     = ["ecs.e4.large","ecs.e4.large","ecs.e4.large"]
}

variable "worker_instance_types" {
  description = "The ecs instance types used to launch worker nodes."
  default     = ["ecs.e4.large"]
}

variable "node_cidr_mask" {
  description = "The node cidr block to specific how many pods can run on single node."
  default     = 24
}

variable "enable_ssh" {
  description = "Enable login to the node through SSH."
  default     = true
}

variable "password" {
  description = "The password of ECS instance."
  default     = "inputYourCodeHere"
}

variable "service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  default     = "172.21.0.0/20"
}

variable "pod_cidr" {
  description = "The kubernetes pod cidr block. It cannot be equals to vpc's or vswitch's and cannot be in them."
  default     = "172.20.0.0/16"
}

variable "k8s_number" {
  description = "The number of kubernetes cluster."
  default     = 1
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default_m" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Master"
}

data "alibabacloudstack_instance_types" "default_w" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_cs_kubernetes" "default" {
  master_vswitch_ids = "${var.vswitch_ids}"
  new_nat_gateway = "true"
  enable_ssh = "${var.enable_ssh}"
  password = "${var.password}"
  master_disk_category = "cloud_efficiency"
  node_cidr_mask = "${var.node_cidr_mask}"
  vpc_id = "${var.vpc_id}"
  worker_disk_category = "cloud_efficiency"
  worker_instance_types = "${var.worker_instance_types}"
  master_count = "3"
  service_cidr = "${var.service_cidr}"
  os_type = "linux"
  name = "${var.name}"
  master_instance_types = "${var.master_instance_types}"
  platform = "CentOS"
  version = "1.18.8-aliyun.1"
  //worker_data_disks {
  //  size = "100"
  //  encrypted = "false"
  //  category = "cloud_efficiency"
  //}
	//worker_data_disk = "true"
	worker_data_disk_category = "cloud_efficiency"
	worker_data_disk_size = "100"
  proxy_mode = "ipvs"
  master_disk_size = "45"
  worker_vswitch_ids = "${var.vswitch_ids}"
  slb_internet_enabled = "true"
  timeout_mins = "25"
  worker_disk_size = "30"
  num_of_nodes = "${var.worker_number}"
  //count = "${var.k8s_number}"
  pod_cidr = "${var.pod_cidr}"
  delete_protection = "false"
}

//resource "alibabacloudstack_cs_kubernetes" "default" {
//  name = "${var.name}"
//  master_vswitch_ids = "${var.vswitch_ids}"
//  worker_vswitch_ids = "${var.vswitch_ids}"
//  new_nat_gateway = true
//  master_instance_types = "${var.master_instance_types}"
//  worker_instance_types = "${var.worker_instance_types}"
//  num_of_nodes = "${var.worker_number}"
//  vpc_id = "${var.vpc_id}"
//  password = "inputYourCodeHere"
//  pod_cidr = "172.20.0.0/16"
//  service_cidr = "172.21.0.0/20"
//  enable_ssh = true
//  install_cloud_monitor = true
//  worker_disk_category  = "cloud_efficiency"
//  worker_data_disk_category = "cloud_efficiency"
//  worker_data_disk_size =  30
//  master_disk_size = 45
//}
`, name)
}
