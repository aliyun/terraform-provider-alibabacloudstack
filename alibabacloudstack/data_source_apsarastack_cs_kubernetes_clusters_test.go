package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCsK8sClustersDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
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
		resourceId:        resourceId,
		existMapFunc:      existCSKubernetesClustersMapFunc,
		fakeMapFunc:       fakeCSKubernetesClustersMapFunc,
		ExternalProviders: testAccExternalProviders,
	}
	csKubernetesClustersCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceCSKubernetesClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	%s

	%s


	resource "alibabacloudstack_cs_kubernetes" "default" {
		name						= var.name
		version						= "1.30.7-aliyun.1"
		os_type						= "linux"
		platform					= "AliyunLinux"
		num_of_nodes				= "3"
		master_count				= "3"
		master_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"]
		master_instance_types		= ["ecs.n4v2.large","ecs.n4v2.large","ecs.n4v2.large"]
		master_disk_category		= "cloud_ssd"
		vpc_id						= "${alibabacloudstack_vpc_vpc.default.id}"
		worker_instance_types		= ["ecs.n4v2.large"]
		worker_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}"]
		worker_disk_category		= "cloud_ssd"
		password					= random_password.password.0.result
		pod_cidr					= "172.20.0.0/16"
		service_cidr				= "172.21.0.0/20"
		worker_disk_size			= "40"
		master_disk_size			= "40"
		slb_internet_enabled		= "true"
		security_group_id			= alibabacloudstack_ecs_securitygroup.default.id
		runtime	 {
			name	= "containerd"
			version	= "1.6.28"
		}
	}

	locals {
		k8s_cluster_id = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? data.alibabacloudstack_cs_kubernetes_clusters.default.ids.0 : alibabacloudstack_cs_kubernetes.default.0.id
		k8s_cluster_name = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? data.alibabacloudstack_cs_kubernetes_clusters.default.names.0 : alibabacloudstack_cs_kubernetes.default.0.name
	}
	`, name, SecurityGroupCommonTestCase, RandomPasswordTestCase(12, 1))
}
