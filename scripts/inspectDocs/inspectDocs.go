package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var TerraformProviderPath string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请传入一个参数：terraform项目地址")
		os.Exit(1)
	}
	TerraformProviderPath = os.Args[1]
	mp := make(map[string]string)
	for k, v := range Mp {
		mp[v] = k
	}

	dirPath := TerraformProviderPath + "/alibabacloudstack"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		fullFilePath := filepath.Join(dirPath, file.Name())
		if strings.HasSuffix(fullFilePath, "test.go") {
			continue
		}
		if strings.HasSuffix(fullFilePath, ".go") {
			processSource(fullFilePath)
		}
	}

	//处理docs
	dirPath = TerraformProviderPath + "/website/docs/d"
	files, err = os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		fullFilePath := filepath.Join(dirPath, file.Name())
		processDocs(fullFilePath)
	}

	dirPath = TerraformProviderPath + "/website/docs/r"
	files, err = os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		fullFilePath := filepath.Join(dirPath, file.Name())
		processDocs(fullFilePath)
	}

	// // 打印结果
	// for _, s := range sources {
	// 	fmt.Printf("Source name: %s\nFields: %v\n", s.name, s.field)
	// 	fmt.Println()

	// }

	// for _, s := range docs {
	// 	fmt.Printf("docs name: %s\nFields: %v\n", s.name, s.field)
	// 	fmt.Println()

	// }
	for _, s := range sources {
		filename := strings.TrimPrefix(s.filename, "../../")
		doc, ok := mp[s.name]
		if !ok {
			fmt.Println("[失败] 文件:", filename, " 文档:")
			continue
		}
		field_mp := make(map[string]interface{})
		for i := 0; i < len(s.field); i++ {
			field := s.field[i]
			field_mp[field] = struct{}{}
		}
		var d *source = nil
		for _, de := range docs {

			if doc != de.name {
				continue
			}
			d = de
			for i := 0; i < len(d.field); i++ {
				field := d.field[i]
				delete(field_mp, field)
			}
			break
		}
		if d == nil {
			fmt.Println("[失败] 文件:", filename, " 文档:")
			continue
		}
		if len(field_mp) == 0 {
			continue
		}
		docname := strings.TrimPrefix(d.filename, "../../")
		fmt.Println("[失败] 文件:", filename, " 文档:", docname)
		for k, _ := range field_mp {
			fmt.Println("      参数: ", k, "  缺少")
		}

	}

}

type source struct {
	filename string
	name     string
	field    []string
}

var sources []*source
var docs []*source
var currentSource *source
var readingFields bool

func processSource(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fieldRegex := regexp.MustCompile(`^\s*\"(\w+)\"`)

	for scanner.Scan() {
		line := scanner.Text()
		heads := strings.Fields(line)
		if len(heads) >= 4 && heads[2] == "*schema.Resource" && heads[3] == "{" {
			readingFields = true
			currentSource = &source{
				name:     heads[1][:len(heads[1])-2],
				filename: filePath,
			}
			sources = append(sources, currentSource)
		} else if readingFields {
			line = strings.TrimSpace(line) // 去除每行的空格
			if len(line) == 0 {            // 忽略空行
				continue
			}
			if strings.HasSuffix(line, "{") {
				matches := fieldRegex.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentSource.field = append(currentSource.field, matches[1])
				}
			}
			if strings.HasSuffix(line, "}") {
				readingFields = false
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

}
func processDocs(filePath string) {
	fieldPattern := regexp.MustCompile(`^\*\s*` + "`" + `([^` + "`" + `]+)` + "`")

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var newSource *source
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# alibabacloudstack") {
			newSource = &source{
				name:     strings.ReplaceAll(line, "\\_", "_"),
				filename: filePath,
			}
			newSource.name = strings.TrimPrefix(newSource.name, "# ")

			docs = append(docs, newSource)
		} else if matches := fieldPattern.FindStringSubmatch(line); matches != nil {
			// 使用正则表达式匹配和提取字段名称
			if len(matches) > 1 {
				// 第一个捕获组是我们感兴趣的字段名
				fieldName := matches[1]
				// 确保docs不为空
				if len(docs) > 0 {
					// 添加字段到最新的source实例
					docs[len(docs)-1].field = append(docs[len(docs)-1].field, fieldName)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}

var Mp = map[string]string{
	"alibabacloudstack_account":                                "dataSourceAlibabacloudStackAccount",
	"alibabacloudstack_adb_clusters":                           "dataSourceAlibabacloudStackAdbDbClusters",
	"alibabacloudstack_adb_zones":                              "dataSourceAlibabacloudStackAdbZones",
	"alibabacloudstack_adb_db_clusters":                        "dataSourceAlibabacloudStackAdbDbClusters",
	"alibabacloudstack_api_gateway_apis":                       "dataSourceAlibabacloudStackApiGatewayApis",
	"alibabacloudstack_api_gateway_apps":                       "dataSourceAlibabacloudStackApiGatewayApps",
	"alibabacloudstack_api_gateway_groups":                     "dataSourceAlibabacloudStackApiGatewayGroups",
	"alibabacloudstack_api_gateway_service":                    "dataSourceAlibabacloudStackApiGatewayService",
	"alibabacloudstack_ascm_resource_groups":                   "dataSourceAlibabacloudStackAscmResourceGroups",
	"alibabacloudstack_ascm_users":                             "dataSourceAlibabacloudStackAscmUsers",
	"alibabacloudstack_ascm_user_groups":                       "dataSourceAlibabacloudStackAscmUserGroups",
	"alibabacloudstack_ascm_logon_policies":                    "dataSourceAlibabacloudStackAscmLogonPolicies",
	"alibabacloudstack_ascm_ram_service_roles":                 "dataSourceAlibabacloudStackAscmRamServiceRoles",
	"alibabacloudstack_ascm_organizations":                     "dataSourceAlibabacloudStackAscmOrganizations",
	"alibabacloudstack_ascm_instance_families":                 "dataSourceAlibabacloudStackInstanceFamilies",
	"alibabacloudstack_ascm_regions_by_product":                "dataSourceAlibabacloudStackRegionsByProduct",
	"alibabacloudstack_ascm_service_cluster_by_product":        "dataSourceAlibabacloudStackServiceClusterByProduct",
	"alibabacloudstack_ascm_ecs_instance_families":             "dataSourceAlibabacloudStackEcsInstanceFamilies",
	"alibabacloudstack_ascm_specific_fields":                   "dataSourceAlibabacloudStackSpecificFields",
	"alibabacloudstack_ascm_environment_services_by_product":   "dataSourceAlibabacloudStackAscmEnvironmentServicesByProduct",
	"alibabacloudstack_ascm_password_policies":                 "dataSourceAlibabacloudStackAscmPasswordPolicies",
	"alibabacloudstack_ascm_quotas":                            "dataSourceAlibabacloudStackQuotas",
	"alibabacloudstack_ascm_metering_query_ecs":                "dataSourceAlibabacloudstackAscmMeteringQueryEcs",
	"alibabacloudstack_ascm_roles":                             "dataSourceAlibabacloudStackAscmRoles",
	"alibabacloudstack_ascm_ram_policies":                      "dataSourceAlibabacloudStackAscmRamPolicies",
	"alibabacloudstack_ascm_ram_policies_for_user":             "dataSourceAlibabacloudStackAscmRamPoliciesForUser",
	"alibabacloudstack_common_bandwidth_packages":              "dataSourceAlibabacloudStackCommonBandwidthPackages",
	"alibabacloudstack_cr_ee_instances":                        "dataSourceAlibabacloudStackCrEEInstances",
	"alibabacloudstack_cr_ee_namespaces":                       "dataSourceAlibabacloudStackCrEENamespaces",
	"alibabacloudstack_cr_ee_repos":                            "dataSourceAlibabacloudStackCrEERepos",
	"alibabacloudstack_cr_ee_sync_rules":                       "dataSourceAlibabacloudStackCrEESyncRules",
	"alibabacloudstack_cr_namespaces":                          "dataSourceAlibabacloudStackCRNamespaces",
	"alibabacloudstack_cr_repos":                               "dataSourceAlibabacloudStackCRRepos",
	"alibabacloudstack_cs_kubernetes_clusters":                 "dataSourceAlibabacloudStackCSKubernetesClusters",
	"alibabacloudstack_cms_alarm_contacts":                     "dataSourceAlibabacloudstackCmsAlarmContacts",
	"alibabacloudstack_cms_alarm_contact_groups":               "dataSourceAlibabacloudstackCmsAlarmContactGroups",
	"alibabacloudstack_cms_project_meta":                       "dataSourceAlibabacloudstackCmsProjectMeta",
	"alibabacloudstack_cms_metric_metalist":                    "dataSourceAlibabacloudstackCmsMetricMetalist",
	"alibabacloudstack_cms_alarms":                             "dataSourceAlibabacloudstackCmsAlarms",
	"alibabacloudstack_datahub_service":                        "dataSourceAlibabacloudStackDatahubService",
	"alibabacloudstack_db_instances":                           "dataSourceAlibabacloudStackDBInstances",
	"alibabacloudstack_db_zones":                               "dataSourceAlibabacloudStackDBZones",
	"alibabacloudstack_disks":                                  "dataSourceAlibabacloudStackDisks",
	"alibabacloudstack_dns_records":                            "dataSourceAlibabacloudStackDnsRecords",
	"alibabacloudstack_dns_groups":                             "dataSourceAlibabacloudStackDnsGroups",
	"alibabacloudstack_drds_instances":                         "dataSourceAlibabacloudStackDRDSInstances",
	"alibabacloudstack_dms_enterprise_instances":               "dataSourceAlibabacloudStackDmsEnterpriseInstances",
	"alibabacloudstack_dms_enterprise_users":                   "dataSourceAlibabacloudStackDmsEnterpriseUsers",
	"alibabacloudstack_ecs_commands":                           "dataSourceAlibabacloudStackEcsCommands",
	"alibabacloudstack_ecs_deployment_sets":                    "dataSourceAlibabacloudStackEcsDeploymentSets",
	"alibabacloudstack_ecs_hpc_clusters":                       "dataSourceAlibabacloudStackEcsHpcClusters",
	"alibabacloudstack_ecs_dedicated_hosts":                    "dataSourceAlibabacloudStackEcsDedicatedHosts",
	"alibabacloudstack_edas_deploy_groups":                     "dataSourceAlibabacloudStackEdasDeployGroups",
	"alibabacloudstack_edas_clusters":                          "dataSourceAlibabacloudStackEdasClusters",
	"alibabacloudstack_edas_applications":                      "dataSourceAlibabacloudStackEdasApplications",
	"alibabacloudstack_eips":                                   "dataSourceAlibabacloudStackEips",
	"alibabacloudstack_ess_scaling_configurations":             "dataSourceAlibabacloudStackEssScalingConfigurations",
	"alibabacloudstack_ess_scaling_groups":                     "dataSourceAlibabacloudStackEssScalingGroups",
	"alibabacloudstack_ess_lifecycle_hooks":                    "dataSourceAlibabacloudStackEssLifecycleHooks",
	"alibabacloudstack_ess_notifications":                      "dataSourceAlibabacloudStackEssNotifications",
	"alibabacloudstack_ess_scaling_rules":                      "dataSourceAlibabacloudStackEssScalingRules",
	"alibabacloudstack_ess_scheduled_tasks":                    "dataSourceAlibabacloudStackEssScheduledTasks",
	"alibabacloudstack_forward_entries":                        "dataSourceAlibabacloudStackForwardEntries",
	"alibabacloudstack_gpdb_accounts":                          "dataSourceAlibabacloudStackGpdbAccounts",
	"alibabacloudstack_gpdb_instances":                         "dataSourceAlibabacloudStackGpdbInstances",
	"alibabacloudstack_hbase_instances":                        "dataSourceAlibabacloudStackHBaseInstances",
	"alibabacloudstack_instances":                              "dataSourceAlibabacloudStackInstances",
	"alibabacloudstack_instance_type_families":                 "dataSourceAlibabacloudStackInstanceTypeFamilies",
	"alibabacloudstack_instance_types":                         "dataSourceAlibabacloudStackInstanceTypes",
	"alibabacloudstack_images":                                 "dataSourceAlibabacloudStackImages",
	"alibabacloudstack_key_pairs":                              "dataSourceAlibabacloudStackKeyPairs",
	"alibabacloudstack_kms_aliases":                            "dataSourceAlibabacloudStackKmsAliases",
	"alibabacloudstack_kms_ciphertext":                         "dataSourceAlibabacloudStackKmsCiphertext",
	"alibabacloudstack_kms_keys":                               "dataSourceAlibabacloudStackKmsKeys",
	"alibabacloudstack_kms_secrets":                            "dataSourceAlibabacloudStackKmsSecrets",
	"alibabacloudstack_kvstore_instances":                      "dataSourceAlibabacloudStackKVStoreInstances",
	"alibabacloudstack_kvstore_zones":                          "dataSourceAlibabacloudStackKVStoreZones",
	"alibabacloudstack_kvstore_instance_classes":               "dataSourceAlibabacloudStackKVStoreInstanceClasses",
	"alibabacloudstack_kvstore_instance_engines":               "dataSourceAlibabacloudStackKVStoreInstanceEngines",
	"alibabacloudstack_mongodb_instances":                      "dataSourceAlibabacloudStackMongoDBInstances",
	"alibabacloudstack_mongodb_zones":                          "dataSourceAlibabacloudStackMongoDBZones",
	"alibabacloudstack_maxcompute_cus":                         "dataSourceAlibabacloudStackMaxcomputeCus",
	"alibabacloudstack_maxcompute_users":                       "dataSourceAlibabacloudStackMaxcomputeUsers",
	"alibabacloudstack_maxcompute_clusters":                    "dataSourceAlibabacloudStackMaxcomputeClusters",
	"alibabacloudstack_maxcompute_cluster_qutaos":              "dataSourceAlibabacloudStackMaxcomputeClusterQutaos",
	"alibabacloudstack_maxcompute_projects":                    "dataSourceAlibabacloudStackMaxcomputeProjects",
	"alibabacloudstack_nas_zones":                              "dataSourceAlibabacloudStackNasZones",
	"alibabacloudstack_nas_protocols":                          "dataSourceAlibabacloudStackNasProtocols",
	"alibabacloudstack_nas_file_systems":                       "dataSourceAlibabacloudStackFileSystems",
	"alibabacloudstack_nas_mount_targets":                      "dataSourceAlibabacloudStackNasMountTargets",
	"alibabacloudstack_nas_access_rules":                       "dataSourceAlibabacloudStackAccessRules",
	"alibabacloudstack_nat_gateways":                           "dataSourceAlibabacloudStackNatGateways",
	"alibabacloudstack_network_acls":                           "dataSourceAlibabacloudStackNetworkAcls",
	"alibabacloudstack_network_interfaces":                     "dataSourceAlibabacloudStackNetworkInterfaces",
	"alibabacloudstack_oss_buckets":                            "dataSourceAlibabacloudStackOssBuckets",
	"alibabacloudstack_oss_bucket_objects":                     "dataSourceAlibabacloudStackOssBucketObjects",
	"alibabacloudstack_ons_instances":                          "dataSourceAlibabacloudStackOnsInstances",
	"alibabacloudstack_ons_topics":                             "dataSourceAlibabacloudStackOnsTopics",
	"alibabacloudstack_ons_groups":                             "dataSourceAlibabacloudStackOnsGroups",
	"alibabacloudstack_ots_tables":                             "dataSourceAlibabacloudStackOtsTables",
	"alibabacloudstack_ots_instances":                          "dataSourceAlibabacloudStackOtsInstances",
	"alibabacloudstack_ots_instances_attachment":               "dataSourceAlibabacloudStackOtsInstanceAttachments",
	"alibabacloudstack_ots_service":                            "dataSourceAlibabacloudStackOtsService",
	"alibabacloudstack_quick_bi_users":                         "dataSourceAlibabacloudStackQuickBiUsers",
	"alibabacloudstack_router_interfaces":                      "dataSourceAlibabacloudStackRouterInterfaces",
	"alibabacloudstack_ram_service_role_products":              "dataSourceAlibabacloudstackRamServiceRoleProducts",
	"alibabacloudstack_route_tables":                           "dataSourceAlibabacloudStackRouteTables",
	"alibabacloudstack_route_entries":                          "dataSourceAlibabacloudStackRouteEntries",
	"alibabacloudstack_ros_stacks":                             "dataSourceAlibabacloudStackRosStacks",
	"alibabacloudstack_ros_templates":                          "dataSourceAlibabacloudStackRosTemplates",
	"alibabacloudstack_security_groups":                        "dataSourceAlibabacloudStackSecurityGroups",
	"alibabacloudstack_security_group_rules":                   "dataSourceAlibabacloudStackSecurityGroupRules",
	"alibabacloudstack_snapshots":                              "dataSourceAlibabacloudStackSnapshots",
	"alibabacloudstack_slb_listeners":                          "dataSourceAlibabacloudStackSlbListeners",
	"alibabacloudstack_slb_server_groups":                      "dataSourceAlibabacloudStackSlbServerGroups",
	"alibabacloudstack_slb_acls":                               "dataSourceAlibabacloudStackSlbAcls",
	"alibabacloudstack_slb_domain_extensions":                  "dataSourceAlibabacloudStackSlbDomainExtensions",
	"alibabacloudstack_slb_rules":                              "dataSourceAlibabacloudStackSlbRules",
	"alibabacloudstack_slb_master_slave_server_groups":         "dataSourceAlibabacloudStackSlbMasterSlaveServerGroups",
	"alibabacloudstack_slbs":                                   "dataSourceAlibabacloudStackSlbs",
	"alibabacloudstack_slb_zones":                              "dataSourceAlibabacloudStackSlbZones",
	"alibabacloudstack_snat_entries":                           "dataSourceAlibabacloudStackSnatEntries",
	"alibabacloudstack_slb_server_certificates":                "dataSourceAlibabacloudStackSlbServerCertificates",
	"alibabacloudstack_slb_ca_certificates":                    "dataSourceAlibabacloudStackSlbCACertificates",
	"alibabacloudstack_slb_backend_servers":                    "dataSourceAlibabacloudStackSlbBackendServers",
	"alibabacloudstack_tsdb_zones":                             "dataSourceAlibabacloudStackTsdbZones",
	"alibabacloudstack_vpn_gateways":                           "dataSourceAlibabacloudStackVpnGateways",
	"alibabacloudstack_vpn_customer_gateways":                  "dataSourceAlibabacloudStackVpnCustomerGateways",
	"alibabacloudstack_vpn_connections":                        "dataSourceAlibabacloudStackVpnConnections",
	"alibabacloudstack_vpc_ipv6_gateways":                      "dataSourceAlibabacloudStackVpcIpv6Gateways",
	"alibabacloudstack_vpc_ipv6_egress_rules":                  "dataSourceAlibabacloudStackVpcIpv6EgressRules",
	"alibabacloudstack_vpc_ipv6_addresses":                     "dataSourceAlibabacloudStackVpcIpv6Addresses",
	"alibabacloudstack_vpc_ipv6_internet_bandwidths":           "dataSourceAlibabacloudStackVpcIpv6InternetBandwidths",
	"alibabacloudstack_vswitches":                              "dataSourceAlibabacloudStackVSwitches",
	"alibabacloudstack_vpcs":                                   "dataSourceAlibabacloudStackVpcs",
	"alibabacloudstack_zones":                                  "dataSourceAlibabacloudStackZones",
	"alibabacloudstack_elasticsearch_instances":                "dataSourceAlibabacloudStackElasticsearch",
	"alibabacloudstack_elasticsearch_zones":                    "dataSourceAlibabacloudStackElaticsearchZones",
	"alibabacloudstack_ehpc_job_templates":                     "dataSourceAlibabacloudStackEhpcJobTemplates",
	"alibabacloudstack_oos_executions":                         "dataSourceAlibabacloudStackOosExecutions",
	"alibabacloudstack_oos_templates":                          "dataSourceAlibabacloudStackOosTemplates",
	"alibabacloudstack_express_connect_physical_connections":   "dataSourceAlibabacloudStackExpressConnectPhysicalConnections",
	"alibabacloudstack_express_connect_access_points":          "dataSourceAlibabacloudStackExpressConnectAccessPoints",
	"alibabacloudstack_express_connect_virtual_border_routers": "dataSourceAlibabacloudStackExpressConnectVirtualBorderRouters",
	"alibabacloudStack_cloud_firewall_control_policies":        "dataSourceAlibabacloudStackCloudFirewallControlPolicies",
	"alibabacloudstack_ecs_ebs_storage_sets":                   "dataSourceAlibabacloudStackEcsEbsStorageSets",
	"alibabacloudstack_ess_scaling_configuration":              "resourceAlibabacloudStackEssScalingConfiguration",
	"alibabacloudstack_adb_account":                            "resourceAlibabacloudStackAdbAccount",
	"alibabacloudstack_adb_backup_policy":                      "resourceAlibabacloudStackAdbBackupPolicy",
	"alibabacloudstack_adb_cluster":                            "resourceAlibabacloudStackAdbDbCluster",
	"alibabacloudstack_adb_connection":                         "resourceAlibabacloudStackAdbConnection",
	"alibabacloudstack_adb_db_cluster":                         "resourceAlibabacloudStackAdbDbCluster",
	"alibabacloudstack_alikafka_sasl_acl":                      "resourceAlibabacloudStackAlikafkaSaslAcl",
	"alibabacloudstack_alikafka_sasl_user":                     "resourceAlibabacloudStackAlikafkaSaslUser",
	"alibabacloudstack_alikafka_topic":                         "resourceAlibabacloudStackAlikafkaTopic",
	"alibabacloudstack_api_gateway_api":                        "resourceAlibabacloudStackApigatewayApi",
	"alibabacloudstack_api_gateway_app":                        "resourceAlibabacloudStackApigatewayApp",
	"alibabacloudstack_api_gateway_app_attachment":             "resourceAliyunApigatewayAppAttachment",
	"alibabacloudstack_api_gateway_group":                      "resourceAlibabacloudStackApigatewayGroup",
	"alibabacloudstack_api_gateway_vpc_access":                 "resourceAlibabacloudStackApigatewayVpc",
	"alibabacloudstack_application_deployment":                 "resourceAlibabacloudStackEdasApplicationPackageAttachment",
	"alibabacloudstack_ascm_custom_role":                       "resourceAlibabacloudStackAscmRole",
	"alibabacloudstack_ascm_logon_policy":                      "resourceAlibabacloudStackLogonPolicy",
	"alibabacloudstack_ascm_organization":                      "resourceAlibabacloudStackAscmOrganization",
	"alibabacloudstack_ascm_password_policy":                   "resourceAlibabacloudStackAscmPasswordPolicy",
	"alibabacloudstack_ascm_quota":                             "resourceAlibabacloudStackAscmQuota",
	"alibabacloudstack_ascm_ram_policy":                        "resourceAlibabacloudStackAscmRamPolicy",
	"alibabacloudstack_ascm_ram_policy_for_role":               "resourceAlibabacloudStackAscmRamPolicyForRole",
	"alibabacloudstack_ascm_ram_role":                          "resourceAlibabacloudStackAscmRamRole",
	"alibabacloudstack_ascm_resource_group":                    "resourceAlibabacloudStackAscmResourceGroup",
	"alibabacloudstack_ascm_user":                              "resourceAlibabacloudStackAscmUser",
	"alibabacloudstack_ascm_user_group":                        "resourceAlibabacloudStackAscmUserGroup",
	"alibabacloudstack_ascm_user_group_resource_set_binding":   "resourceAlibabacloudStackAscmUserGroupResourceSetBinding",
	"alibabacloudstack_ascm_user_group_role_binding":           "resourceAlibabacloudStackAscmUserGroupRoleBinding",
	"alibabacloudstack_ascm_user_role_binding":                 "resourceAlibabacloudStackAscmUserRoleBinding",
	"alibabacloudstack_ascm_usergroup_user":                    "resourceAlibabacloudStackAscmUserGroupUser",
	"alibabacloudstack_cms_alarm":                              "resourceAlibabacloudStackCmsAlarm",
	"alibabacloudstack_cms_alarm_contact":                      "resourceAlibabacloudstackCmsAlarmContact",
	"alibabacloudstack_cms_alarm_contact_group":                "resourceAlibabacloudstackCmsAlarmContactGroup",
	"alibabacloudstack_cms_site_monitor":                       "resourceAlibabacloudStackCmsSiteMonitor",
	"alibabacloudstack_common_bandwidth_package":               "resourceAlibabacloudStackCommonBandwidthPackage",
	"alibabacloudstack_common_bandwidth_package_attachment":    "resourceAlibabacloudStackCommonBandwidthPackageAttachment",
	"alibabacloudstack_cr_ee_namespace":                        "resourceAlibabacloudStackCrEENamespace",
	"alibabacloudstack_cr_ee_repo":                             "resourceAlibabacloudStackCrEERepo",
	"alibabacloudstack_cr_ee_sync_rule":                        "resourceAlibabacloudStackCrEESyncRule",
	"alibabacloudstack_cr_namespace":                           "resourceAlibabacloudStackCRNamespace",
	"alibabacloudstack_cr_repo":                                "resourceAlibabacloudStackCRRepo",
	"alibabacloudstack_cs_kubernetes":                          "resourceAlibabacloudStackCSKubernetes",
	"alibabacloudstack_cs_kubernetes_node_pool":                "resourceAlibabacloudStackCSKubernetesNodePool",
	"alibabacloudstack_datahub_project":                        "resourceAlibabacloudStackDatahubProject",
	"alibabacloudstack_datahub_subscription":                   "resourceAlibabacloudStackDatahubSubscription",
	"alibabacloudstack_datahub_topic":                          "resourceAlibabacloudStackDatahubTopic",
	"alibabacloudstack_db_account":                             "resourceAlibabacloudStackDBAccount",
	"alibabacloudstack_db_account_privilege":                   "resourceAlibabacloudStackDBAccountPrivilege",
	"alibabacloudstack_db_backup_policy":                       "resourceAlibabacloudStackDBBackupPolicy",
	"alibabacloudstack_db_connection":                          "resourceAlibabacloudStackDBConnection",
	"alibabacloudstack_db_database":                            "resourceAlibabacloudStackDBDatabase",
	"alibabacloudstack_db_instance":                            "resourceAlibabacloudStackDBInstance",
	"alibabacloudstack_db_read_write_splitting_connection":     "resourceAlibabacloudStackDBReadWriteSplittingConnection",
	"alibabacloudstack_db_readonly_instance":                   "resourceAlibabacloudStackDBReadonlyInstance",
	"alibabacloudstack_disk":                                   "resourceAlibabacloudStackDisk",
	"alibabacloudstack_disk_attachment":                        "resourceAlibabacloudStackDiskAttachment",
	"alibabacloudstack_dms_enterprise_instance":                "resourceAlibabacloudStackDmsEnterpriseInstance",
	"alibabacloudstack_dms_enterprise_user":                    "resourceAlibabacloudStackDmsEnterpriseUser",
	"alibabacloudstack_dns_domain":                             "resourceAlibabacloudStackDnsDomain",
	"alibabacloudstack_dns_domain_attachment":                  "resourceAlibabacloudStackDnsDomainAttachment",
	"alibabacloudstack_dns_group":                              "resourceAlibabacloudStackDnsGroup",
	"alibabacloudstack_dns_record":                             "resourceAlibabacloudStackDnsRecord",
	"alibabacloudstack_drds_instance":                          "resourceAlibabacloudStackDRDSInstance",
	"alibabacloudstack_dts_subscription_job":                   "resourceAlibabacloudStackDtsSubscriptionJob",
	"alibabacloudstack_dts_synchronization_instance":           "resourceAlibabacloudStackDtsSynchronizationInstance",
	"alibabacloudstack_dts_synchronization_job":                "resourceAlibabacloudStackDtsSynchronizationJob",
	"alibabacloudstack_ecs_command":                            "resourceAlibabacloudStackEcsCommand",
	"alibabacloudstack_ecs_dedicated_host":                     "resourceAlibabacloudStackEcsDedicatedHost",
	"alibabacloudstack_ecs_deployment_set":                     "resourceAlibabacloudStackEcsDeploymentSet",
	"alibabacloudstack_ecs_hpc_cluster":                        "resourceAlibabacloudStackEcsHpcCluster",
	"alibabacloudstack_ecs_ebs_storage_set":                    "resourceAlibabacloudStackEcsEbsStorageSets",
	"alibabacloudstack_edas_application":                       "resourceAlibabacloudStackEdasApplication",
	"alibabacloudstack_edas_application_scale":                 "resourceAlibabacloudStackEdasInstanceApplicationAttachment",
	"alibabacloudstack_edas_cluster":                           "resourceAlibabacloudStackEdasCluster",
	"alibabacloudstack_edas_deploy_group":                      "resourceAlibabacloudStackEdasDeployGroup",
	"alibabacloudstack_edas_instance_cluster_attachment":       "resourceAlibabacloudStackEdasInstanceClusterAttachment",
	"alibabacloudstack_edas_k8s_application":                   "resourceAlibabacloudStackEdasK8sApplication",
	"alibabacloudstack_edas_k8s_cluster":                       "resourceAlibabacloudStackEdasK8sCluster",
	"alibabacloudstack_edas_slb_attachment":                    "resourceAlibabacloudStackEdasSlbAttachment",
	"alibabacloudstack_ehpc_job_template":                      "resourceAlibabacloudStackEhpcJobTemplate",
	"alibabacloudstack_eip":                                    "resourceAlibabacloudStackEip",
	"alibabacloudstack_eip_association":                        "resourceAlibabacloudStackEipAssociation",
	"alibabacloudstack_ess_alarm":                              "resourceAlibabacloudStackEssAlarm",
	"alibabacloudstack_ess_attachment":                         "resourceAlibabacloudstackEssAttachment",
	"alibabacloudstack_ess_lifecycle_hook":                     "resourceAlibabacloudStackEssLifecycleHook",
	"alibabacloudstack_ess_notification":                       "resourceAlibabacloudStackEssNotification",
	"alibabacloudstack_ess_scaling_group":                      "resourceAlibabacloudStackEssScalingGroup",
	"alibabacloudstack_ess_scaling_rule":                       "resourceAlibabacloudStackEssScalingRule",
	"alibabacloudstack_ess_scalinggroup_vserver_groups":        "resourceAlibabacloudStackEssScalingGroupVserverGroups",
	"alibabacloudstack_ess_scheduled_task":                     "resourceAlibabacloudStackEssScheduledTask",
	"alibabacloudstack_forward_entry":                          "resourceAlibabacloudStackForwardEntry",
	"alibabacloudstack_gpdb_account":                           "resourceAlibabacloudStackGpdbAccount",
	"alibabacloudstack_gpdb_connection":                        "resourceAlibabacloudStackGpdbConnection",
	"alibabacloudstack_gpdb_instance":                          "resourceAlibabacloudStackGpdbInstance",
	"alibabacloudstack_hbase_instance":                         "resourceAlibabacloudStackHBaseInstance",
	"alibabacloudstack_image":                                  "resourceAlibabacloudStackImage",
	"alibabacloudstack_image_copy":                             "resourceAlibabacloudStackImageCopy",
	"alibabacloudstack_image_export":                           "resourceAlibabacloudStackImageExport",
	"alibabacloudstack_image_import":                           "resourceAlibabacloudStackImageImport",
	"alibabacloudstack_image_share_permission":                 "resourceAlibabacloudStackImageSharePermission",
	"alibabacloudstack_instance":                               "resourceAlibabacloudStackInstance",
	"alibabacloudstack_key_pair":                               "resourceAlibabacloudStackKeyPair",
	"alibabacloudstack_key_pair_attachment":                    "resourceAlibabacloudStackKeyPairAttachment",
	"alibabacloudstack_kms_alias":                              "resourceAlibabacloudStackKmsAlias",
	// "alibabacloudstack_kms_ciphertext": "resourceAlibabacloudStackKmsCiphertext",
	"alibabacloudstack_kms_key":                               "resourceAlibabacloudStackKmsKey",
	"alibabacloudstack_kms_secret":                            "resourceAlibabacloudStackKmsSecret",
	"alibabacloudstack_kvstore_account":                       "resourceAlibabacloudStackKVstoreAccount",
	"alibabacloudstack_kvstore_backup_policy":                 "resourceAlibabacloudStackKVStoreBackupPolicy",
	"alibabacloudstack_kvstore_connection":                    "resourceAlibabacloudStackKvstoreConnection",
	"alibabacloudstack_kvstore_instance":                      "resourceAlibabacloudStackKVStoreInstance",
	"alibabacloudstack_launch_template":                       "resourceAlibabacloudStackLaunchTemplate",
	"alibabacloudstack_log_machine_group":                     "resourceAlibabacloudStackLogMachineGroup",
	"alibabacloudstack_log_project":                           "resourceAlibabacloudStackLogProject",
	"alibabacloudstack_log_store":                             "resourceAlibabacloudStackLogStore",
	"alibabacloudstack_log_store_index":                       "resourceAlibabacloudStackLogStoreIndex",
	"alibabacloudstack_logtail_attachment":                    "resourceAlibabacloudStackLogtailAttachment",
	"alibabacloudstack_logtail_config":                        "resourceAlibabacloudStackLogtailConfig",
	"alibabacloudstack_maxcompute_project":                    "resourceAlibabacloudStackMaxcomputeProject",
	"alibabacloudstack_maxcompute_user":                       "resourceAlibabacloudStackMaxcomputeUser",
	"alibabacloudstack_maxcompute_cu":                         "resourceAlibabacloudStackMaxcomputeCu",
	"alibabacloudstack_mongodb_instance":                      "resourceAlibabacloudStackMongoDBInstance",
	"alibabacloudstack_mongodb_sharding_instance":             "resourceAlibabacloudStackMongoDBShardingInstance",
	"alibabacloudstack_nas_access_group":                      "resourceAlibabacloudStackNasAccessGroup",
	"alibabacloudstack_nas_access_rule":                       "resourceAlibabacloudStackNasAccessRule",
	"alibabacloudstack_nas_file_system":                       "resourceAlibabacloudStackNasFileSystem",
	"alibabacloudstack_nas_mount_target":                      "resourceAlibabacloudStackNasMountTarget",
	"alibabacloudstack_nat_gateway":                           "resourceAlibabacloudStackNatGateway",
	"alibabacloudstack_network_acl":                           "resourceAlibabacloudStackNetworkAcl",
	"alibabacloudstack_network_acl_attachment":                "resourceAlibabacloudStackNetworkAclAttachment",
	"alibabacloudstack_network_acl_entries":                   "resourceAlibabacloudStackNetworkAclEntries",
	"alibabacloudstack_network_interface":                     "resourceAlibabacloudStackNetworkInterface",
	"alibabacloudstack_network_interface_attachment":          "resourceNetworkInterfaceAttachment",
	"alibabacloudstack_ons_group":                             "resourceAlibabacloudStackOnsGroup",
	"alibabacloudstack_ons_instance":                          "resourceAlibabacloudStackOnsInstance",
	"alibabacloudstack_ons_topic":                             "resourceAlibabacloudStackOnsTopic",
	"alibabacloudstack_oss_bucket":                            "resourceAlibabacloudStackOssBucket",
	"alibabacloudstack_oss_bucket_quota":                      "resourceAlibabacloudStackOssBucketQuota",
	"alibabacloudstack_oss_bucket_kms":                        "resourceAlibabacloudStackOssBucketKms",
	"alibabacloudstack_oss_bucket_object":                     "resourceAlibabacloudStackOssBucketObject",
	"alibabacloudstack_ots_instance":                          "resourceAlibabacloudStackOtsInstance",
	"alibabacloudstack_ots_instance_attachment":               "resourceAlibabacloudStackOtsInstanceAttachment",
	"alibabacloudstack_ots_table":                             "resourceAlibabacloudStackOtsTable",
	"alibabacloudstack_quick_bi_user":                         "resourceAlibabacloudStackQuickBiUser",
	"alibabacloudstack_quick_bi_user_group":                   "resourceAlibabacloudStackQuickBiUserGroup",
	"alibabacloudstack_quick_bi_workspace":                    "resourceAlibabacloudStackQuickBiWorkspace",
	"alibabacloudstack_ram_role_attachment":                   "resourceAlibabacloudStackRamRoleAttachment",
	"alibabacloudstack_reserved_instance":                     "resourceAlibabacloudStackReservedInstance",
	"alibabacloudstack_ros_stack":                             "resourceAlibabacloudStackRosStack",
	"alibabacloudstack_ros_template":                          "resourceAlibabacloudStackRosTemplate",
	"alibabacloudstack_route_entry":                           "resourceAlibabacloudStackRouteEntry",
	"alibabacloudstack_route_table":                           "resourceAlibabacloudStackRouteTable",
	"alibabacloudstack_route_table_attachment":                "resourceAlibabacloudStackRouteTableAttachment",
	"alibabacloudstack_router_interface":                      "resourceAlibabacloudStackRouterInterface",
	"alibabacloudstack_router_interface_connection":           "resourceAlibabacloudStackRouterInterfaceConnection",
	"alibabacloudstack_security_group":                        "resourceAlibabacloudStackSecurityGroup",
	"alibabacloudstack_security_group_rule":                   "resourceAlibabacloudStackSecurityGroupRule",
	"alibabacloudstack_slb":                                   "resourceAlibabacloudStackSlb",
	"alibabacloudstack_slb_acl":                               "resourceAlibabacloudStackSlbAcl",
	"alibabacloudstack_slb_backend_server":                    "resourceAlibabacloudStackSlbBackendServer",
	"alibabacloudstack_slb_ca_certificate":                    "resourceAlibabacloudStackSlbCACertificate",
	"alibabacloudstack_slb_domain_extension":                  "resourceAlibabacloudStackSlbDomainExtension",
	"alibabacloudstack_slb_listener":                          "resourceAlibabacloudStackSlbListener",
	"alibabacloudstack_slb_master_slave_server_group":         "resourceAlibabacloudStackSlbMasterSlaveServerGroup",
	"alibabacloudstack_slb_rule":                              "resourceAlibabacloudStackSlbRule",
	"alibabacloudstack_slb_server_certificate":                "resourceAlibabacloudStackSlbServerCertificate",
	"alibabacloudstack_slb_server_group":                      "resourceAlibabacloudStackSlbServerGroup",
	"alibabacloudstack_snapshot":                              "resourceAlibabacloudStackSnapshot",
	"alibabacloudstack_snapshot_policy":                       "resourceAlibabacloudStackSnapshotPolicy",
	"alibabacloudstack_snat_entry":                            "resourceAlibabacloudStackSnatEntry",
	"alibabacloudstack_vpc":                                   "resourceAlibabacloudStackVpc",
	"alibabacloudstack_vpc_ipv6_egress_rule":                  "resourceAlibabacloudStackVpcIpv6EgressRule",
	"alibabacloudstack_vpc_ipv6_gateway":                      "resourceAlibabacloudStackVpcIpv6Gateway",
	"alibabacloudstack_vpc_ipv6_internet_bandwidth":           "resourceAlibabacloudStackVpcIpv6InternetBandwidth",
	"alibabacloudstack_vpn_connection":                        "resourceAlibabacloudStackVpnConnection",
	"alibabacloudstack_vpn_customer_gateway":                  "resourceAlibabacloudStackVpnCustomerGateway",
	"alibabacloudstack_vpn_gateway":                           "resourceAlibabacloudStackVpnGateway",
	"alibabacloudstack_vpn_route_entry":                       "resourceAlibabacloudStackVpnRouteEntry",
	"alibabacloudstack_vswitch":                               "resourceAlibabacloudStackSwitch",
	"alibabacloudstack_data_works_folder":                     "resourceAlibabacloudStackDataWorksFolder",
	"alibabacloudstack_data_works_connection":                 "resourceAlibabacloudStackDataWorksConnection",
	"alibabacloudstack_data_works_user":                       "resourceAlibabacloudStackDataWorksUser",
	"alibabacloudstack_data_works_project":                    "resourceAlibabacloudStackDataWorksProject",
	"alibabacloudstack_data_works_user_role_binding":          "resourceAlibabacloudStackDataWorksUserRoleBinding",
	"alibabacloudstack_data_works_remind":                     "resourceAlibabacloudStackDataWorksRemind",
	"alibabacloudstack_elasticsearch_instance":                "resourceAlibabacloudStackElasticsearch",
	"alibabacloudstack_dbs_backup_plan":                       "resourceAlibabacloudStackDbsBackupPlan",
	"alibabacloudstack_express_connect_physical_connection":   "resourceAlibabacloudStackExpressConnectPhysicalConnection",
	"alibabacloudstack_express_connect_virtual_border_router": "resourceAlibabacloudStackExpressConnectVirtualBorderRouter",
	"alibabacloudstack_oos_template":                          "resourceAlibabacloudStackOosTemplate",
	"alibabacloudstack_oos_execution":                         "resourceAlibabacloudStackOosExecution",
	"alibabacloudstack_arms_alert_contact":                    "resourceAlibabacloudStackArmsAlertContact",
	"alibabacloudstack_arms_alert_contact_group":              "resourceAlibabacloudStackArmsAlertContactGroup",
	"alibabacloudstack_arms_dispatch_rule":                    "resourceAlibabacloudStackArmsDispatchRule",
	"alibabacloudstack_arms_prometheus_alert_rule":            "resourceAlibabacloudStackArmsPrometheusAlertRule",
	"alibabacloudstack_elasticsearch_k8s_instance":            "resourceAlibabacloudStackElasticsearchOnk8s",
	"alibabacloudstack_cloud_firewall_control_policy":         "resourceAlibabacloudStackCloudFirewallControlPolicy",
	"alibabacloudstack_cloud_firewall_control_policy_order":   "resourceAlibabacloudStackCloudFirewallControlPolicyOrder",
	"alibabacloudstack_csb_project":                           "resourceAlibabacloudStackCsbProject",
	"alibabacloudstack_graph_database_db_instance":            "resourceAlibabacloudStackGraphDatabaseDbInstance",
}
