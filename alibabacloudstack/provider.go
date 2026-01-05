package alibabacloudstack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/go-homedir"
)

const DefaultProviderName = "registry.terraform.io/aliyun/alibabacloudstack"

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ACCESS_KEY", nil),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SECRET_KEY", nil),
				Description: descriptions["secret_key"],
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_REGION", nil),
				Description: descriptions["region"],
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				//DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_REGION", nil),
				Description: descriptions["region_id"],
				Deprecated:  "Field region_id is Deprecated, Please use parameter region replace it.",
			},
			"role_arn": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   descriptions["assume_role_role_arn"],
				DefaultFunc:   schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ASSUME_ROLE_ARN", nil),
				ConflictsWith: []string{"security_token"},
			},
			"security_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SECURITY_TOKEN", nil),
				Description:   descriptions["security_token"],
				ConflictsWith: []string{"role_arn"},
			},
			"ecs_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ECS_ROLE_NAME", nil),
				Description: descriptions["ecs_role_name"],
			},
			"skip_region_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: descriptions["skip_region_validation"],
				Deprecated:  "always skip to valiate region in apsarastack",
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_PROFILE", nil),
			},
			"endpoints": endpointsSchema(),
			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_credentials_file"],
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SHARED_CREDENTIALS_FILE", nil),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_INSECURE", false),
				Description: descriptions["insecure"],
			},
			"assume_role": assumeRoleSchema(),
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ALIBABACLOUDSTACK_PROTOCOL", "HTTP"),
				Description:  descriptions["protocol"],
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"client_read_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLIENT_READ_TIMEOUT", 60000),
				Description: descriptions["client_read_timeout"],
			},
			"client_connect_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLIENT_CONNECT_TIMEOUT", 60000),
				Description: descriptions["client_connect_timeout"],
			},
			"source_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SOURCE_IP", nil),
				Description: descriptions["source_ip"],
			},
			"security_transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SECURITY_TRANSPORT", nil),
				//Deprecated:  "It has been deprecated from version 1.136.0 and using new field secure_transport instead.",
			},
			"secure_transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SECURE_TRANSPORT", nil),
				Description: descriptions["secure_transport"],
			},
			"configuration_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  descriptions["configuration_source"],
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_PROXY", nil),
				Description: descriptions["proxy"],
			},
			"is_center_region": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_CENTER_REGION", "true"),
				Description: descriptions["is_center_region"],
			},
			"popgw_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_POPGW_DOMAIN", nil),
				Description: descriptions["popgw_domain"],
			},
			"ossservice_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["ossservice_domain"],
				Deprecated:  "This parameter will no longer be valid and will be removed in 3.19.0.",
			},
			"organization_accesskey": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ORGANIZATION_ACCESSKEY", nil),
				Description: descriptions["organization_accesskey"],
				Deprecated:  "Use access_key replace organization_accesskey.",
			},
			"organization_secretkey": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ORGANIZATION_SECRETKEY", nil),
				Description: descriptions["organization_secretkey"],
				Deprecated:  "Use secret_key replace organization_secretkey.",
			},
			"sls_openapi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SLS_OPENAPI_ENDPOINT", nil),
				Description: descriptions["sls_openapi_endpoint"],
				Deprecated:  "Use schema endpoints replace sls_openapi_endpoint.",
			},
			"sts_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_STS_ENDPOINT", nil),
				Description: descriptions["sts_endpoint"],
				Deprecated:  "Use schema endpoints replace sts_endpoint.",
			},
			"department": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_DEPARTMENT", nil),
				Description: descriptions["department"],
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_RESOURCE_GROUP", nil),
				Description: descriptions["resource_group"],
			},
			"resource_group_set_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_RESOURCE_GROUP_SET", nil),
				Description: descriptions["resource_group_set_name"],
			},
			"kms_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_KMS_ENDPOINT", nil),
				Description: descriptions["kms_endpoint"],
				Deprecated:  "Use schema endpoints replace kms_endpoint.",
			},
			"asapi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ASAPI_ENDPOINT", nil),
				Description: descriptions["asapi_endpoint"],
				Deprecated:  "Use schema endpoints replace asapi_endpoint.",
			},
			"sls_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SLS_ENDPOINT", nil),
				Description: descriptions["sls_endpoint"],
				Deprecated:  "Use schema endpoints replace sls_endpoint.",
			},
			"max_retry_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MAX_RETRY_TIMEOUT", 0),
				Description: descriptions["max_retry_timeout"],
			},
		},
		DataSourcesMap: getDataSourcesMap(),
		ResourcesMap:   getResourcesMap(),
		ConfigureFunc:  providerConfigure,
	}
}

var providerConfig map[string]interface{}

func stringToBool(value string) (bool, error) {
	// Convert string to lowercase for comparison
	value = strings.ToLower(value)

	// Check common boolean representations
	switch value {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value for environment variable: %s", value)
	}
}

func getDataSourcesMap() map[string]*schema.Resource {
	maps := map[string]*schema.Resource{
		"alibabacloudstack_account":                                          dataSourceAlibabacloudStackAccount(),
		"alibabacloudstack_adb_clusters":                                     dataSourceAlibabacloudStackAdbDbClusters(),
		"alibabacloudstack_adb_zones":                                        dataSourceAlibabacloudStackAdbZones(),
		"alibabacloudstack_adb_db_clusters":                                  dataSourceAlibabacloudStackAdbDbClusters(),
		"alibabacloudstack_adb_dbclusters":                                   dataSourceAlibabacloudStackAdbDbClusters(),
		"alibabacloudstack_alikafka_instances":                               dataSourceAlicloudAlikafkaInstances(),
		"alibabacloudstack_api_gateway_apis":                                 dataSourceAlibabacloudStackApiGatewayApis(),
		"alibabacloudstack_apigateway_apis":                                  dataSourceAlibabacloudStackApiGatewayApis(),
		"alibabacloudstack_api_gateway_apps":                                 dataSourceAlibabacloudStackApiGatewayApps(),
		"alibabacloudstack_apigateway_apps":                                  dataSourceAlibabacloudStackApiGatewayApps(),
		"alibabacloudstack_api_gateway_groups":                               dataSourceAlibabacloudStackApiGatewayGroups(),
		"alibabacloudstack_apigateway_apigroups":                             dataSourceAlibabacloudStackApiGatewayGroups(),
		"alibabacloudstack_api_gateway_service":                              dataSourceAlibabacloudStackApiGatewayService(),
		"alibabacloudstack_ascm_resource_groups":                             dataSourceAlibabacloudStackAscmResourceGroups(),
		"alibabacloudstack_ascm_users":                                       dataSourceAlibabacloudStackAscmUsers(),
		"alibabacloudstack_ascm_user_groups":                                 dataSourceAlibabacloudStackAscmUserGroups(),
		"alibabacloudstack_ascm_logon_policies":                              dataSourceAlibabacloudStackAscmLogonPolicies(),
		"alibabacloudstack_ascm_ram_service_roles":                           dataSourceAlibabacloudStackAscmRamServiceRoles(),
		"alibabacloudstack_ascm_organizations":                               dataSourceAlibabacloudStackAscmOrganizations(),
		"alibabacloudstack_ascm_instance_families":                           dataSourceAlibabacloudStackInstanceFamilies(),
		"alibabacloudstack_ascm_regions_by_product":                          dataSourceAlibabacloudStackRegionsByProduct(),
		"alibabacloudstack_ascm_ecs_instance_families":                       dataSourceAlibabacloudStackEcsInstanceFamilies(),
		"alibabacloudstack_ascm_specific_fields":                             dataSourceAlibabacloudStackSpecificFields(),
		"alibabacloudstack_ascm_password_policies":                           dataSourceAlibabacloudStackAscmPasswordPolicies(),
		"alibabacloudstack_ascm_quotas":                                      dataSourceAlibabacloudStackQuotas(),
		"alibabacloudstack_ascm_roles":                                       dataSourceAlibabacloudStackAscmRoles(),
		"alibabacloudstack_ascm_ram_roles":                                   dataSourceAlibabacloudStackAscmRoles(),
		"alibabacloudstack_ascm_ram_policies":                                dataSourceAlibabacloudStackAscmRamPolicies(),
		"alibabacloudstack_ascm_ram_policies_for_user":                       dataSourceAlibabacloudStackAscmRamPoliciesForUser(),
		"alibabacloudstack_common_bandwidth_packages":                        dataSourceAlibabacloudStackCommonBandwidthPackages(),
		"alibabacloudstack_cbwp_commonbandwidthpackages":                     dataSourceAlibabacloudStackCommonBandwidthPackages(),
		"alibabacloudstack_cr_ee_instances":                                  dataSourceAlibabacloudStackCrEeInstances(),
		"alibabacloudstack_cr_ee_namespaces":                                 dataSourceAlibabacloudStackCrEeNamespaces(),
		"alibabacloudstack_cr_ee_repos":                                      dataSourceAlibabacloudStackCrEeRepos(),
		"alibabacloudstack_cr_ee_sync_rules":                                 dataSourceAlibabacloudStackCrEeSyncRules(),
		"alibabacloudstack_cr_namespaces":                                    dataSourceAlibabacloudStackCRNamespaces(),
		"alibabacloudstack_cr_repositories":                                  dataSourceAlibabacloudStackCRRepos(),
		"alibabacloudstack_cr_repos":                                         dataSourceAlibabacloudStackCRRepos(),
		"alibabacloudstack_cs_kubernetes_clusters":                           dataSourceAlibabacloudStackCSKubernetesClusters(),
		"alibabacloudstack_ack_clusters":                                     dataSourceAlibabacloudStackCSKubernetesClusters(),
		"alibabacloudstack_ack_templates":                                    dataSourceAlibabacloudStackAckTemplates(),
		"alibabacloudstack_cs_kubernetes_clusters_kubeconfig":                dataSourceAlibabacloudStackCSKubernetesClustersKubeConfig(),
		"alibabacloudstack_cms_metric_rule_templates":                        dataSourceAlibabacloudStackCmsMetricRuleTemplates(),
		"alibabacloudstack_cloudmonitorservice_metricruletemplates":          dataSourceAlibabacloudStackCmsMetricRuleTemplates(),
		"alibabacloudstack_cms_alarm_contacts":                               dataSourceAlibabacloudStackCmsAlarmContacts(),
		"alibabacloudstack_cloudmonitorservice_alarmcontacts":                dataSourceAlibabacloudStackCmsAlarmContacts(),
		"alibabacloudstack_cms_alarm_contact_groups":                         dataSourceAlibabacloudStackCmsAlarmContactGroups(),
		"alibabacloudstack_cloudmonitorservice_alarmcontactgroups":           dataSourceAlibabacloudStackCmsAlarmContactGroups(),
		"alibabacloudstack_cms_project_meta":                                 dataSourceAlibabacloudStackCmsProjectMeta(),
		"alibabacloudstack_cms_metric_metalist":                              dataSourceAlibabacloudStackCmsMetricMetalist(),
		"alibabacloudstack_cms_alarms":                                       dataSourceAlibabacloudStackCmsAlarms(),
		"alibabacloudstack_cloudmonitorservice_metricalarmrules":             dataSourceAlibabacloudStackCmsAlarms(),
		"alibabacloudstack_datahub_projects":                                 dataSourceAlibabacloudStackDatahubProjects(),
		"alibabacloudstack_datahub_topics":                                   dataSourceAlibabacloudStackDatahubTopics(),
		"alibabacloudstack_datahub_subscriptions":                            dataSourceAlibabacloudStackDatahubSubscriptions(),
		"alibabacloudstack_datahub_kafka_groups":                             dataSourceAlibabacloudStackDatahubKafkaGroups(),
		"alibabacloudstack_rds_instance_types":                               dataSourceAlibabacloudStackRdsInstanceTypes(),
		"alibabacloudstack_db_instances":                                     dataSourceAlibabacloudStackDBInstances(),
		"alibabacloudstack_rds_dbinstances":                                  dataSourceAlibabacloudStackDBInstances(),
		"alibabacloudstack_db_zones":                                         dataSourceAlibabacloudStackDBZones(),
		"alibabacloudstack_disks":                                            dataSourceAlibabacloudStackDisks(),
		"alibabacloudstack_ecs_disks":                                        dataSourceAlibabacloudStackDisks(),
		"alibabacloudstack_dns_records":                                      dataSourceAlibabacloudStackDnsRecords(),
		"alibabacloudstack_dns_domains":                                      dataSourceAlibabacloudStackDnsDomains(),
		"alibabacloudstack_drds_instances":                                   dataSourceAlibabacloudStackDRDSInstances(),
		"alibabacloudstack_drds_rds_instances":                               dataSourceAlibabacloudStackDrdsRdsInstances(),
		"alibabacloudstack_drds_instance_series":                             dataSourceAlibabacloudStackDrdsInstanceSeries(),
		"alibabacloudstack_drds_instance_specifications":                     dataSourceAlibabacloudStackDrdsInstanceSpecifications(),
		"alibabacloudstack_drds_databases":                                   dataSourceAlibabacloudStackDrdsDatabases(),
		"alibabacloudstack_drds_accounts":                                    dataSourceAlibabacloudStackDrdsAccounts(),
		"alibabacloudstack_dms_enterprise_instances":                         dataSourceAlibabacloudStackDmsEnterpriseInstances(),
		"alibabacloudstack_dmsenterprise_instances":                          dataSourceAlibabacloudStackDmsEnterpriseInstances(),
		"alibabacloudstack_dms_enterprise_users":                             dataSourceAlibabacloudStackDmsEnterpriseUsers(),
		"alibabacloudstack_dmsenterprise_users":                              dataSourceAlibabacloudStackDmsEnterpriseUsers(),
		"alibabacloudstack_ecs_commands":                                     dataSourceAlibabacloudStackEcsCommands(),
		"alibabacloudstack_ecs_deployment_sets":                              dataSourceAlibabacloudStackEcsDeploymentSets(),
		"alibabacloudstack_ecs_deploymentsets":                               dataSourceAlibabacloudStackEcsDeploymentSets(),
		"alibabacloudstack_ecs_hpc_clusters":                                 dataSourceAlibabacloudStackEcsHpcClusters(),
		"alibabacloudstack_ecs_hpcclusters":                                  dataSourceAlibabacloudStackEcsHpcClusters(),
		"alibabacloudstack_ecs_dedicated_hosts":                              dataSourceAlibabacloudStackEcsDedicatedHosts(),
		"alibabacloudstack_ecs_dedicatedhosts":                               dataSourceAlibabacloudStackEcsDedicatedHosts(),
		"alibabacloudstack_ecs_dedicated_host_cluster":                       dataSourceAlibabacloudStackEcsDedicatedHostClusters(),
		"alibabacloudstack_edas_deploy_groups":                               dataSourceAlibabacloudStackEdasDeployGroups(),
		"alibabacloudstack_edas_deploygroups":                                dataSourceAlibabacloudStackEdasDeployGroups(),
		"alibabacloudstack_edas_clusters":                                    dataSourceAlibabacloudStackEdasClusters(),
		"alibabacloudstack_edas_applications":                                dataSourceAlibabacloudStackEdasApplications(),
		"alibabacloudstack_edas_slbattachments":                              dataSourceAlibabacloudStackEdasApplications(),
		"alibabacloudstack_edas_namespaces":                                  dataSourceAlibabacloudStackEdasNamespaces(),
		"alibabacloudstack_eips":                                             dataSourceAlibabacloudStackEips(),
		"alibabacloudstack_eip_addresses":                                    dataSourceAlibabacloudStackEips(),
		"alibabacloudstack_autoscaling_scalingconfigurations":                dataSourceAlibabacloudStackEssScalingConfigurations(),
		"alibabacloudstack_ess_scaling_configurations":                       dataSourceAlibabacloudStackEssScalingConfigurations(),
		"alibabacloudstack_autoscaling_scalinggroups":                        dataSourceAlibabacloudStackEssScalingGroups(),
		"alibabacloudstack_ess_scaling_groups":                               dataSourceAlibabacloudStackEssScalingGroups(),
		"alibabacloudstack_ess_lifecycle_hooks":                              dataSourceAlibabacloudStackEssLifecycleHooks(),
		"alibabacloudstack_autoscaling_lifecyclehooks":                       dataSourceAlibabacloudStackEssLifecycleHooks(),
		"alibabacloudstack_ess_notifications":                                dataSourceAlibabacloudStackEssNotifications(),
		"alibabacloudstack_autoscaling_notifications":                        dataSourceAlibabacloudStackEssNotifications(),
		"alibabacloudstack_autoscaling_scalingrules":                         dataSourceAlibabacloudStackEssScalingRules(),
		"alibabacloudstack_ess_scaling_rules":                                dataSourceAlibabacloudStackEssScalingRules(),
		"alibabacloudstack_ess_scheduled_tasks":                              dataSourceAlibabacloudStackEssScheduledTasks(),
		"alibabacloudstack_autoscaling_scheduledtasks":                       dataSourceAlibabacloudStackEssScheduledTasks(),
		"alibabacloudstack_forward_entries":                                  dataSourceAlibabacloudStackForwardEntries(),
		"alibabacloudstack_natgateway_forwardentries":                        dataSourceAlibabacloudStackForwardEntries(),
		"alibabacloudstack_gpdb_accounts":                                    dataSourceAlibabacloudStackGpdbAccounts(),
		"alibabacloudstack_gpdb_instances":                                   dataSourceAlibabacloudStackGpdbInstances(),
		"alibabacloudstack_gpdb_dbinstances":                                 dataSourceAlibabacloudStackGpdbInstances(),
		"alibabacloudstack_hbase_instances":                                  dataSourceAlibabacloudStackHBaseInstances(),
		"alibabacloudstack_hbase_clusters":                                   dataSourceAlibabacloudStackHBaseInstances(),
		"alibabacloudstack_instances":                                        dataSourceAlibabacloudStackInstances(),
		"alibabacloudstack_ecs_instances":                                    dataSourceAlibabacloudStackInstances(),
		"alibabacloudstack_instance_type_families":                           dataSourceAlibabacloudStackInstanceTypeFamilies(),
		"alibabacloudstack_instance_types":                                   dataSourceAlibabacloudStackInstanceTypes(),
		"alibabacloudstack_images":                                           dataSourceAlibabacloudStackImages(),
		"alibabacloudstack_ecs_images":                                       dataSourceAlibabacloudStackImages(),
		"alibabacloudstack_key_pairs":                                        dataSourceAlibabacloudStackKeyPairs(),
		"alibabacloudstack_ecs_keypairs":                                     dataSourceAlibabacloudStackKeyPairs(),
		"alibabacloudstack_kms_aliases":                                      dataSourceAlibabacloudStackKmsAliases(),
		"alibabacloudstack_kms_ciphertexts":                                  dataSourceAlibabacloudStackKmsCiphertext(),
		"alibabacloudstack_kms_keys":                                         dataSourceAlibabacloudStackKmsKeys(),
		"alibabacloudstack_kms_secrets":                                      dataSourceAlibabacloudStackKmsSecrets(),
		"alibabacloudstack_kvstore_instances":                                dataSourceAlibabacloudStackKVStoreInstances(),
		"alibabacloudstack_redis_tairinstances":                              dataSourceAlibabacloudStackKVStoreInstances(),
		"alibabacloudstack_rds_backups":                                      dataSourceAlibabacloudStackRdsBackups(),
		"alibabacloudstack_kvstore_zones":                                    dataSourceAlibabacloudStackKVStoreZones(),
		"alibabacloudstack_kvstore_instance_classes":                         dataSourceAlibabacloudStackKVStoreInstanceClasses(),
		"alibabacloudstack_kvstore_instance_engines":                         dataSourceAlibabacloudStackKVStoreInstanceEngines(),
		"alibabacloudstack_mongodb_instance_types":                           dataSourceAlibabacloudStackMongoDBInstanceTypes(),
		"alibabacloudstack_mongodb_instances":                                dataSourceAlibabacloudStackMongoDBInstances(),
		"alibabacloudstack_mongodb_accounts":                                 dataSourceAlibabacloudStackMongodbAccounts(),
		"alibabacloudstack_mongodb_zones":                                    dataSourceAlibabacloudStackMongoDBZones(),
		"alibabacloudstack_mongodb_backups":                                  dataSourceAlibabacloudStackMongodbBackups(),
		"alibabacloudstack_maxcompute_cus":                                   dataSourceAlibabacloudStackMaxcomputeCus(),
		"alibabacloudstack_maxcompute_users":                                 dataSourceAlibabacloudStackMaxcomputeUsers(),
		"alibabacloudstack_maxcompute_clusters":                              dataSourceAlibabacloudStackMaxcomputeClusters(),
		"alibabacloudstack_maxcompute_cluster_quotas":                        dataSourceAlibabacloudStackMaxcomputeClusterQuotas(),
		"alibabacloudstack_maxcompute_projects":                              dataSourceAlibabacloudStackMaxcomputeProjects(),
		"alibabacloudstack_nas_zones":                                        dataSourceAlibabacloudStackNasZones(),
		"alibabacloudstack_nas_file_systems":                                 dataSourceAlibabacloudStackFileSystems(),
		"alibabacloudstack_nas_filesystems":                                  dataSourceAlibabacloudStackFileSystems(),
		"alibabacloudstack_nas_mount_targets":                                dataSourceAlibabacloudStackNasMountTargets(),
		"alibabacloudstack_nas_mounttargets":                                 dataSourceAlibabacloudStackNasMountTargets(),
		"alibabacloudstack_nas_access_rules":                                 dataSourceAlibabacloudStackAccessRules(),
		"alibabacloudstack_nas_accessrules":                                  dataSourceAlibabacloudStackAccessRules(),
		"alibabacloudstack_nat_gateways":                                     dataSourceAlibabacloudStackNatGateways(),
		"alibabacloudstack_nas_lifecycle_policies":                           dataSourceAlibabacloudStackNasLifecyclePolicies(),
		"alibabacloudstack_natgateway_natgateways":                           dataSourceAlibabacloudStackNatGateways(),
		"alibabacloudstack_natgateway_bandwidth_packages":                    dataSourceAlibabacloudStackNatgatewayBandwidthPackages(),
		"alibabacloudstack_network_acls":                                     dataSourceAlibabacloudStackNetworkAcls(),
		"alibabacloudstack_vpc_networkacls":                                  dataSourceAlibabacloudStackNetworkAcls(),
		"alibabacloudstack_network_interfaces":                               dataSourceAlibabacloudStackNetworkInterfaces(),
		"alibabacloudstack_ecs_networkinterfaces":                            dataSourceAlibabacloudStackNetworkInterfaces(),
		"alibabacloudstack_oss_buckets":                                      dataSourceAlibabacloudStackOssBuckets(),
		"alibabacloudstack_oss_bucket_objects":                               dataSourceAlibabacloudStackOssBucketObjects(),
		"alibabacloudstack_oss_clusters":                                     dataSourceAlibabacloudStackOssClusters(),
		"alibabacloudstack_ons_instances":                                    dataSourceAlibabacloudStackOnsInstances(),
		"alibabacloudstack_ons_topics":                                       dataSourceAlibabacloudStackOnsTopics(),
		"alibabacloudstack_ons_groups":                                       dataSourceAlibabacloudStackOnsGroups(),
		"alibabacloudstack_ots_tables":                                       dataSourceAlibabacloudStackOtsTables(),
		"alibabacloudstack_ots_instances":                                    dataSourceAlibabacloudStackOtsInstances(),
		"alibabacloudstack_ots_instance_attachments":                         dataSourceAlibabacloudStackOtsInstanceAttachments(),
		"alibabacloudstack_ots_service":                                      dataSourceAlibabacloudStackOtsService(),
		"alibabacloudstack_quick_bi_users":                                   dataSourceAlibabacloudStackQuickBiUsers(),
		"alibabacloudstack_router_interfaces":                                dataSourceAlibabacloudStackRouterInterfaces(),
		"alibabacloudstack_expressconnect_routerinterfaces":                  dataSourceAlibabacloudStackRouterInterfaces(),
		"alibabacloudstack_route_tables":                                     dataSourceAlibabacloudStackRouteTables(),
		"alibabacloudstack_vpc_routetables":                                  dataSourceAlibabacloudStackRouteTables(),
		"alibabacloudstack_route_entries":                                    dataSourceAlibabacloudStackRouteEntries(),
		"alibabacloudstack_ros_stacks":                                       dataSourceAlibabacloudStackRosStacks(),
		"alibabacloudstack_ros_templates":                                    dataSourceAlibabacloudStackRosTemplates(),
		"alibabacloudstack_security_groups":                                  dataSourceAlibabacloudStackSecurityGroups(),
		"alibabacloudstack_ecs_securitygroups":                               dataSourceAlibabacloudStackSecurityGroups(),
		"alibabacloudstack_security_group_rules":                             dataSourceAlibabacloudStackSecurityGroupRules(),
		"alibabacloudstack_snapshots":                                        dataSourceAlibabacloudStackSnapshots(),
		"alibabacloudstack_ecs_snapshots":                                    dataSourceAlibabacloudStackSnapshots(),
		"alibabacloudstack_slb_listeners":                                    dataSourceAlibabacloudStackSlbListeners(),
		"alibabacloudstack_slb_server_groups":                                dataSourceAlibabacloudStackSlbServerGroups(),
		"alibabacloudstack_slb_vservergroups":                                dataSourceAlibabacloudStackSlbServerGroups(),
		"alibabacloudstack_slb_acls":                                         dataSourceAlibabacloudStackSlbAcls(),
		"alibabacloudstack_slb_accesscontrollists":                           dataSourceAlibabacloudStackSlbAcls(),
		"alibabacloudstack_slb_domain_extensions":                            dataSourceAlibabacloudStackSlbDomainExtensions(),
		"alibabacloudstack_slb_domainextensions":                             dataSourceAlibabacloudStackSlbDomainExtensions(),
		"alibabacloudstack_slb_rules":                                        dataSourceAlibabacloudStackSlbRules(),
		"alibabacloudstack_slb_master_slave_server_groups":                   dataSourceAlibabacloudStackSlbMasterSlaveServerGroups(),
		"alibabacloudstack_slb_masterslaveservergroups":                      dataSourceAlibabacloudStackSlbMasterSlaveServerGroups(),
		"alibabacloudstack_slbs":                                             dataSourceAlibabacloudStackSlbs(),
		"alibabacloudstack_slb_loadbalancers":                                dataSourceAlibabacloudStackSlbs(),
		"alibabacloudstack_slb_zones":                                        dataSourceAlibabacloudStackSlbZones(),
		"alibabacloudstack_snat_entries":                                     dataSourceAlibabacloudStackSnatEntries(),
		"alibabacloudstack_natgateway_snatentries":                           dataSourceAlibabacloudStackSnatEntries(),
		"alibabacloudstack_slb_server_certificates":                          dataSourceAlibabacloudStackSlbServerCertificates(),
		"alibabacloudstack_slb_servercertificates":                           dataSourceAlibabacloudStackSlbServerCertificates(),
		"alibabacloudstack_slb_ca_certificates":                              dataSourceAlibabacloudStackSlbCACertificates(),
		"alibabacloudstack_slb_cacertificates":                               dataSourceAlibabacloudStackSlbCACertificates(),
		"alibabacloudstack_slb_backend_servers":                              dataSourceAlibabacloudStackSlbBackendServers(),
		"alibabacloudstack_slb_backendservers":                               dataSourceAlibabacloudStackSlbBackendServers(),
		"alibabacloudstack_tsdb_zones":                                       dataSourceAlibabacloudStackTsdbZones(),
		"alibabacloudstack_vpn_gateways":                                     dataSourceAlibabacloudStackVpnGateways(),
		"alibabacloudstack_vpngateway_vpngateways":                           dataSourceAlibabacloudStackVpnGateways(),
		"alibabacloudstack_vpn_customer_gateways":                            dataSourceAlibabacloudStackVpnCustomerGateways(),
		"alibabacloudstack_vpngateway_customergateways":                      dataSourceAlibabacloudStackVpnCustomerGateways(),
		"alibabacloudstack_vpn_connections":                                  dataSourceAlibabacloudStackVpnConnections(),
		"alibabacloudstack_vpngateway_vpnconnections":                        dataSourceAlibabacloudStackVpnConnections(),
		"alibabacloudstack_vpngateway_ssl_vpnservers":                        dataSourceAlibabacloudStackVpngatewaySslVpnServers(),
		"alibabacloudstack_vpngateway_ssl_vpn_client_certs":                  dataSourceAlibabacloudStackVpngatewaySslVpnClientCerts(),
		"alibabacloudstack_vpngateway_sslvpnclientcerts":                     dataSourceAlibabacloudStackVpngatewaySslVpnClientCerts(),
		"alibabacloudstack_vpc_ipv6_gateways":                                dataSourceAlibabacloudStackVpcIpv6Gateways(),
		"alibabacloudstack_vpc_ipv6_egress_rules":                            dataSourceAlibabacloudStackVpcIpv6EgressRules(),
		"alibabacloudstack_vpc_ipv6_egressrules":                             dataSourceAlibabacloudStackVpcIpv6EgressRules(),
		"alibabacloudstack_vpc_ipv6_addresses":                               dataSourceAlibabacloudStackVpcIpv6Addresses(),
		"alibabacloudstack_vpc_ipv6_internet_bandwidths":                     dataSourceAlibabacloudStackVpcIpv6InternetBandwidths(),
		"alibabacloudstack_vpc_ipv6_internetbandwidths":                      dataSourceAlibabacloudStackVpcIpv6InternetBandwidths(),
		"alibabacloudstack_vswitches":                                        dataSourceAlibabacloudStackVSwitches(),
		"alibabacloudstack_vpc_vswitches":                                    dataSourceAlibabacloudStackVSwitches(),
		"alibabacloudstack_vpcs":                                             dataSourceAlibabacloudStackVpcs(),
		"alibabacloudstack_vpc_vpcs":                                         dataSourceAlibabacloudStackVpcs(),
		"alibabacloudstack_zones":                                            dataSourceAlibabacloudStackZones(),
		"alibabacloudstack_elasticsearch_instances":                          dataSourceAlibabacloudStackElasticsearch(),
		"alibabacloudstack_ehpc_job_templates":                               dataSourceAlibabacloudStackEhpcJobTemplates(),
		"alibabacloudstack_oos_executions":                                   dataSourceAlibabacloudStackOosExecutions(),
		"alibabacloudstack_oos_templates":                                    dataSourceAlibabacloudStackOosTemplates(),
		"alibabacloudstack_express_connect_physical_connections":             dataSourceAlibabacloudStackExpressConnectPhysicalConnections(),
		"alibabacloudstack_expressconnect_physicalconnections":               dataSourceAlibabacloudStackExpressConnectPhysicalConnections(),
		"alibabacloudstack_express_connect_access_points":                    dataSourceAlibabacloudStackExpressConnectAccessPoints(),
		"alibabacloudstack_expressconnect_accesspoints":                      dataSourceAlibabacloudStackExpressConnectAccessPoints(),
		"alibabacloudstack_expressconnect_vbr_pconn_associations":            dataSourceAlibabacloudStackExpressconnectVbrPconnAssociations(),
		"alibabacloudstack_express_connect_virtual_border_routers":           dataSourceAlibabacloudStackExpressConnectVirtualBorderRouters(),
		"alibabacloudstack_expressconnect_virtualborderrouters":              dataSourceAlibabacloudStackExpressConnectVirtualBorderRouters(),
		"alibabacloudstack_expressconnect_bgp_networks":                      dataSourceAlibabacloudStackExpressconnectBgpNetworks(),
		"alibabacloudstack_cloud_firewall_control_policies":                  dataSourceAlibabacloudStackCloudFirewallControlPolicies(),
		"alibabacloudstack_ecs_ebs_storage_sets":                             dataSourceAlibabacloudStackEcsEbsStorageSets(),
		"alibabacloudstack_polardb_zones":                                    dataSourceAlibabacloudStackPolardbZones(),
		"alibabacloudstack_polardb_instance_types":                           dataSourceAlibabacloudStackPolardbInstanceTypes(),
		"alibabacloudstack_polardb_databases":                                dataSourceAlibabacloudStackPolardbDatabases(),
		"alibabacloudstack_polardb_dbinstances":                              dataSourceAlibabacloudStackPolardbDbInstances(),
		"alibabacloudstack_polardb_instances":                                dataSourceAlibabacloudStackPolardbDbInstances(),
		"alibabacloudstack_polardb_accounts":                                 dataSourceAlibabacloudStackPolardbAccounts(),
		"alibabacloudstack_polardb_backups":                                  dataSourceAlibabacloudStackPolardbBackups(),
		"alibabacloudstack_polardb_proxies":                                  dataSourceAlibabacloudStackPolardbProxies(),
		"alibabacloudstack_polardb_cluster_instance_types":                   dataSourceAlibabacloudStackPolardbClusterInstanceTypes(),
		"alibabacloudstack_polardb_cluster_proxy_types":                      dataSourceAlibabacloudStackPolardbClusterProxyTypes(),
		"alibabacloudstack_polardb_cluster_instances":                        dataSourceAlibabacloudStackPolardbClusterInstances(),
		"alibabacloudstack_polardb_cluster_accounts":                         dataSourceAlibabacloudStackPolardbClusterAccounts(),
		"alibabacloudstack_polardb_cluster_databases":                        dataSourceAlibabacloudStackPolardbClusterDatabases(),
		"alibabacloudstack_polardb_cluster_proxies":                          dataSourceAlibabacloudStackPolardbClusterProxies(),
		"alibabacloudstack_polardb_cluster_backup_policies":                  dataSourceAlibabacloudStackPolardbClusterBackupPolicies(),
		"alibabacloudstack_acm_configurations":                               dataSourceAlibabacloudStackAcmConfigurations(),
		"alibabacloudstack_bastionhost_instances":                            dataSourceAlibabacloudStackBastionhostInstances(),
		"alibabacloudstack_waf_instances":                                    dataSourceAlibabacloudStackWafInstances(),
		"alibabacloudstack_flink_namespaces":                                 dataSourceAlibabacloudStackFlinkNamespaces(),
		"alibabacloudstack_ebs_diskreplicapairs":                             dataSourceAlibabacloudStackEbsDiskReplicaPairs(),
		"alibabacloudstack_ebs_diskreplicagroups":                            dataSourceAlibabacloudStackEbsDiskReplicaGroups(),
		"alibabacloudstack_ecs_snapshot_groups":                              dataSourceAlibabacloudStackEcsSnapshotGroups(),
		"alibabacloudstack_ecs_invocations":                                  dataSourceAlibabacloudStackEcsInvocations(),
		"alibabacloudstack_vpc_dhcp_options_sets":                            dataSourceAlibabacloudStackVpcDhcpOptionsSets(),
		"alibabacloudstack_slb_access_logs":                                  dataSourceAlibabacloudStackSlbAccessLogs(),
		"alibabacloudstack_vpc_ipv6_isps":                                    dataSourceAlibabacloudStackVpcIpv6Isps(),
		"alibabacloudstack_vpc_ha_vips":                                      dataSourceAlibabacloudStackVpcHaVips(),
		"alibabacloudstack_vpngateway_vpn_pbr_route_entries":                 dataSourceAlibabacloudStackVpngatewayVpnPbrRouteEntries(),
		"alibabacloudstack_expressconnect_bgp_groups":                        dataSourceAlibabacloudStackExpressconnectBgpGroups(),
		"alibabacloudstack_expressconnect_bgp_peers":                         dataSourceAlibabacloudStackExpressconnectBgpPeers(),
		"alibabacloudstack_db_proxies":                                       dataSourceAlibabacloudStackRdsDbProxies(),
		"alibabacloudstack_polardbx_instances":                               dataSourceAlibabacloudStackPolardbxInstances(),
		"alibabacloudstack_polardbx_accounts":                                dataSourceAlibabacloudStackPolardbxAccounts(),
		"alibabacloudstack_polardbx_databases":                               dataSourceAlibabacloudStackPolardbxDatabases(),
		"alibabacloudstack_polardbx_instance_types":                          dataSourceAlibabacloudStackPolardbxInstanceTypes(),
		"alibabacloudstack_polardbx_backups":                                 dataSourceAlibabacloudStackPolardbxBackups(),
		"alibabacloudstack_polardbx_backup_policies":                         dataSourceAlibabacloudStackPolardbxBackupPolicys(),
		"alibabacloudstack_polardbx_cdc_classes":                             dataSourceAlibabacloudStackPolardbxCdcClasses(),
		"alibabacloudstack_cen_instances":                                    dataSourceAlibabacloudStackCenCenInstances(),
		"alibabacloudstack_cen_transit_router_route_tables":                  dataSourceAlibabacloudStackCenTransitRouterRouterTables(),
		"alibabacloudstack_cen_transit_router_route_entries":                 dataSourceAlibabacloudStackCenTransitRouterRouteEntries(),
		"alibabacloudstack_cen_transit_router_vpc_attachments":               dataSourceAlibabacloudStackCenTransitRouterVpcAttachments(),
		"alibabacloudstack_cen_transit_router_vbr_attachments":               dataSourceAlibabacloudStackCenTransitRouterVbrAttachments(),
		"alibabacloudstack_cen_transit_router_connect_attachments":           dataSourceAlibabacloudStackCenTransitRouterConnectAttachments(),
		"alibabacloudstack_cen_transit_router_route_table_associations":      dataSourceAlibabacloudStackCenTransitRouterRouterTableAssociations(),
		"alibabacloudstack_cen_transit_router_route_table_propagations":      dataSourceAlibabacloudStackCenTransitRouterRouterTablePropagations(),
		"alibabacloudstack_cen_route_maps":                                   dataSourceAlibabacloudStackCenTransitRouterRouteMaps(),
		"alibabacloudstack_cen_transit_router_multicast_domains":             dataSourceAlibabacloudStackCenTransitRouterMulticastDomains(),
		"alibabacloudstack_cen_transit_router_multicast_domain_associations": dataSourceAlibabacloudStackCenTransitRouterMulticastDomainAssociations(),
		"alibabacloudstack_cen_transit_router_multicast_domain_sources":      dataSourceAlibabacloudStackCenTransitRouterMulticastDomainSources(),
		"alibabacloudstack_cen_transit_router_multicast_domain_members":      dataSourceAlibabacloudStackCenTransitRouterMulticastDomainMembers(),
		"alibabacloudstack_cen_vbr_health_checks":                            dataSourceAlibabacloudStackCenVbrHealthChecks(),
		"alibabacloudstack_cen_transit_router_connect_peers":                 dataSourceAlibabacloudStackCenTransitRouterConnectPeers(),
		"alibabacloudstack_edas_k8s_application_scaling_rules":               dataSourceAlibabacloudStackEdasScalingRules(),
		"alibabacloudstack_edas_swimming_lane_groups":                        dataSourceAlibabacloudStackEdasSwimmingLaneGroups(),
		"alibabacloudstack_edas_swimming_lanes":                              dataSourceAlibabacloudStackEdasSwimmingLanes(),
		"alibabacloudstack_hologram_clusters":                                dataSourceAlibabacloudStackHologramClusters(),
		"alibabacloudstack_hologram_instances":                               dataSourceAlibabacloudStackHologramInstances(),
		"alibabacloudstack_lindorm_instances":                                dataSourceAlibabacloudStackLindormInstances(),
		"alibabacloudstack_lindorm_instance_types":                           dataSourceAlibabacloudStackLindormInstanceTypes(),
		"alibabacloudstack_api_gateway_v2_k8s_clusters":                      dataSourceAlibabacloudStackAPIGatewayV2K8sClusters(),
		"alibabacloudstack_api_gateway_v2_instance_types":                    dataSourceAlibabacloudStackAPIGatewayV2InstanceTypes(),
		"alibabacloudstack_api_gateway_v2_instances":                         dataSourceAlibabacloudStackAPIGatewayV2Instances(),
		"alibabacloudstack_api_gateway_v2_certificates":                      dataSourceAlibabacloudStackAPIGatewayV2Certificates(),
		"alibabacloudstack_api_gateway_v2_domains":                           dataSourceAlibabacloudStackAPIGatewayV2Domains(),
		"alibabacloudstack_api_gateway_v2_signatures":                        dataSourceAlibabacloudStackAPIGatewayV2Signatures(),
		"alibabacloudstack_api_gateway_v2_consumers":                         dataSourceAlibabacloudStackAPIGatewayV2Consumers(),
		"alibabacloudstack_api_gateway_v2_service_sources":                   dataSourceAlibabacloudStackAPIGatewayV2ServiceSources(),
		"alibabacloudstack_api_gateway_v2_route_groups":                      dataSourceAlibabacloudStackAPIGatewayV2RouteGroups(),
		"alibabacloudstack_api_gateway_v2_cascade_instances":                 dataSourceAlibabacloudStackAPIGatewayV2CascadeInstances(),
		"alibabacloudstack_api_gateway_v2_cascade_links":                     dataSourceAlibabacloudStackApiGatewayV2CascadeLinks(),
		"alibabacloudstack_api_gateway_v2_services":                          dataSourceAlibabacloudStackAPIGatewayV2Services(),
		"alibabacloudstack_api_gateway_v2_routes":                            dataSourceAlibabacloudStackAPIGatewayV2Routes(),
		"alibabacloudstack_api_gateway_v2_mcpservers":                        dataSourceAlibabacloudStackAPIGatewayV2McpServers(),
		"alibabacloudstack_mqtt_instances":                                   dataSourceAlibabacloudStackMqttInstances(),
		"alibabacloudstack_mqtt_topics":                                      dataSourceAlibabacloudStackMqttTopics(),
		"alibabacloudstack_mqtt_groups":                                      dataSourceAlibabacloudStackMqttGroups(),
		"alibabacloudstack_universal_dns_domains":                            dataSourceAlibabacloudStackUniversalDnsDomains(),
		"alibabacloudstack_universal_dns_records":                            dataSourceAlibabacloudStackUniversalDnsRecords(),
		"alibabacloudstack_universal_dns_lines":                              dataSourceAlibabacloudStackUniversalDnsLines(),
		"alibabacloudstack_dns_forward_domains":                              dataSourceAlibabacloudStackDnsForwardDomains(),
		"alibabacloudstack_dns_lines":                                        dataSourceAlibabacloudStackDnsLines(),
		"alibabacloudstack_dns_recursor_acls":                                dataSourceAlibabacloudStackDnsRecursorAcls(),
		"alibabacloudstack_dns_private_domains":                              dataSourceAlibabacloudStackDnsPrivateDomains(),
		"alibabacloudstack_dns_private_records":                              dataSourceAlibabacloudStackDnsPrivateRecords(),
		"alibabacloudstack_dns_private_lines":                                dataSourceAlibabacloudStackDnsPrivateLines(),
		"alibabacloudstack_dns_gtm_instances":                                dataSourceAlibabacloudStackDnsGtmInstances(),
		"alibabacloudstack_dns_gtm_addresspools":                             dataSourceAlibabacloudStackDnsGtmAddressPools(),
		"alibabacloudstack_dns_gtm_access_strategies":                        dataSourceAlibabacloudStackDnsGtmAccessStrategies(),
		"alibabacloudstack_prometheus_v2_instances":                          dataSourceAlibabacloudStackPrometheusV2Instances(),
		"alibabacloudstack_prometheus_v2_notify_groups":                      dataSourceAlibabacloudStackPrometheusV2NotifyGroups(),
		"alibabacloudstack_prometheus_v2_contacts":                           dataSourceAlibabacloudStackPrometheusV2Contacts(),
		"alibabacloudstack_prometheus_v2_alerts":                             dataSourceAlibabacloudStackPrometheusV2Alerts(),
		"alibabacloudstack_schedulerx2_app_groups":                           dataSourceAlibabacloudStackSchedulerx2AppGroups(),
		"alibabacloudstack_schedulerx2_jobs":                                 dataSourceAlibabacloudStackSchedulerx2Jobs(),
		"alibabacloudstack_schedulerx2_workflows":                            dataSourceAlibabacloudStackSchedulerx2Workflows(),
		"alibabacloudstack_apfs_file_systems":                                dataSourceAlibabacloudStackApfsFileSystems(),
		"alibabacloudstack_apfs_zones":                                       dataSourceAlibabacloudStackApfsZones(),
		"alibabacloudstack_hsm_instances":                                    dataSourceAlibabacloudStackHsmInstances(),
		"alibabacloudstack_hsm_vendors":                                      dataSourceAlibabacloudStackHsmVendors(),
		"alibabacloudstack_hsms":                                             dataSourceAlibabacloudStackHsms(),
		"alibabacloudstack_hsm_clusters":                                     dataSourceAlibabacloudStackHsmClusters(),
		"alibabacloudstack_cspprivate_hsm_instances":                         dataSourceAlibabacloudStackCspprivateHsmInstances(),
		"alibabacloudstack_cspprivate_hsm_vendors":                           dataSourceAlibabacloudStackCspprivateVendors(),
		"alibabacloudstack_cspprivate_hsms":                                  dataSourceAlibabacloudStackCspprivateHsms(),
		"alibabacloudstack_cspprivate_hsm_groups":                            dataSourceAlibabacloudStackCspprivateHsmGroups(),
		"alibabacloudstack_bms_keypairs":                                     dataSourceAlibabacloudStackBmsKeypairs(),
		"alibabacloudstack_aqs_oss_scanconfigs":                              dataSourceAlibabacloudStackAqsOssScanconfigs(),
		"alibabacloudstack_aqs_anti_brute_force_rules":                       dataSourceAlibabacloudStackAqsAntiBruteForceRules(),
		"alibabacloudstack_aqs_web_locks":                                    dataSourceAlibabacloudStackAqsWebLocks(),
		"alibabacloudstack_cr_ee_attestor_lifecycle_rules":                   dataSourceAlibabacloudStackCrEEAttestorLifecycleRules(),
		"alibabacloudstack_cloudfw_address_books":                            dataSourceAlibabacloudStackCloudfwAddressBooks(),
		"alibabacloudstack_cloudfw_vpc_control_policies":                     dataSourceAlibabacloudStackCloudfwVpcControlPolicies(),
	}
	if v, err := stringToBool(os.Getenv("APSARASTACK_IN_ALIBABACLOUDSTACK")); err == nil && !v {
		return maps
	}
	new_map := map[string]*schema.Resource{}
	for key, value := range maps {
		new_map[key] = value
		if strings.HasPrefix(key, "alibabacloudstack_") {
			new_key := strings.Replace(key, "alibabacloudstack_", "apsarastack_", 1)
			new_map[new_key] = value
		}
	}
	return new_map
}

func getResourcesMap() map[string]*schema.Resource {
	maps := map[string]*schema.Resource{
		"alibabacloudstack_ess_scaling_configuration":              resourceAlibabacloudStackEssScalingConfiguration(),
		"alibabacloudstack_adb_account":                            resourceAlibabacloudStackAdbAccount(),
		"alibabacloudstack_adb_backup_policy":                      resourceAlibabacloudStackAdbBackupPolicy(),
		"alibabacloudstack_adb_backuppolicy":                       resourceAlibabacloudStackAdbBackupPolicy(),
		"alibabacloudstack_adb_cluster":                            resourceAlibabacloudStackAdbDbCluster(),
		"alibabacloudstack_adb_connection":                         resourceAlibabacloudStackAdbConnection(),
		"alibabacloudstack_adb_db_cluster":                         resourceAlibabacloudStackAdbDbCluster(),
		"alibabacloudstack_adb_dbcluster":                          resourceAlibabacloudStackAdbDbCluster(),
		"alibabacloudstack_alikafka_instance":                      resourceAlibabacloudStackAlikafkaInstance(),
		"alibabacloudstack_alikafka_sasl_acl":                      resourceAlibabacloudStackAlikafkaSaslAcl(),
		"alibabacloudstack_alikafka_sasl_user":                     resourceAlibabacloudStackAlikafkaSaslUser(),
		"alibabacloudstack_alikafka_topic":                         resourceAlibabacloudStackAlikafkaTopic(),
		"alibabacloudstack_api_gateway_api":                        resourceAlibabacloudStackApigatewayApi(),
		"alibabacloudstack_apigateway_api":                         resourceAlibabacloudStackApigatewayApi(),
		"alibabacloudstack_api_gateway_app":                        resourceAlibabacloudStackApigatewayApp(),
		"alibabacloudstack_api_gateway_app_attachment":             resourceAlibabacloudStackApigatewayAppAttachment(),
		"alibabacloudstack_apigateway_app":                         resourceAlibabacloudStackApigatewayAppAttachment(),
		"alibabacloudstack_api_gateway_group":                      resourceAlibabacloudStackApigatewayGroup(),
		"alibabacloudstack_apigateway_apigroup":                    resourceAlibabacloudStackApigatewayGroup(),
		"alibabacloudstack_api_gateway_vpc_access":                 resourceAlibabacloudStackApigatewayVpc(),
		"alibabacloudstack_apigateway_vpc":                         resourceAlibabacloudStackApigatewayVpc(),
		"alibabacloudstack_application_deployment":                 resourceAlibabacloudStackEdasApplicationPackageAttachment(),
		"alibabacloudstack_edas_applicationpackageattachment":      resourceAlibabacloudStackEdasApplicationPackageAttachment(),
		"alibabacloudstack_ascm_custom_role":                       resourceAlibabacloudStackAscmRole(),
		"alibabacloudstack_ascm_logon_policy":                      resourceAlibabacloudStackLogonPolicy(),
		"alibabacloudstack_ascm_organization":                      resourceAlibabacloudStackAscmOrganization(),
		"alibabacloudstack_ascm_password_policy":                   resourceAlibabacloudStackAscmPasswordPolicy(),
		"alibabacloudstack_ascm_quota":                             resourceAlibabacloudStackAscmQuota(),
		"alibabacloudstack_ascm_ram_policy":                        resourceAlibabacloudStackAscmRamPolicy(),
		"alibabacloudstack_ascm_ram_policy_for_role":               resourceAlibabacloudStackAscmRamPolicyForRole(),
		"alibabacloudstack_ascm_ram_role":                          resourceAlibabacloudStackAscmRamRole(),
		"alibabacloudstack_ascm_resource_group":                    resourceAlibabacloudStackAscmResourceGroup(),
		"alibabacloudstack_ascm_user":                              resourceAlibabacloudStackAscmUser(),
		"alibabacloudstack_ascm_user_group":                        resourceAlibabacloudStackAscmUserGroup(),
		"alibabacloudstack_ascm_user_group_resource_set_binding":   resourceAlibabacloudStackAscmUserGroupResourceSetBinding(),
		"alibabacloudstack_ascm_user_group_role_binding":           resourceAlibabacloudStackAscmUserGroupRoleBinding(),
		"alibabacloudstack_ascm_user_role_binding":                 resourceAlibabacloudStackAscmUserRoleBinding(),
		"alibabacloudstack_ascm_usergroup_user":                    resourceAlibabacloudStackAscmUserGroupUser(),
		"alibabacloudstack_ascm_resource_group_user_attachment":    resourceAlibabacloudStackAscmResourceGroupUserAttachment(),
		"alibabacloudstack_ascm_service_ram_role":                  resourceAlibabacloudStackAscmServiceRamRole(),
		"alibabacloudstack_ascm_ram_service_role":                  resourceAlibabacloudStackAscmServiceRamRole(),
		"alibabacloudstack_cms_alarm":                              resourceAlibabacloudStackCmsAlarm(),
		"alibabacloudstack_cloudmonitorservice_metricalarmrule":    resourceAlibabacloudStackCmsAlarm(),
		"alibabacloudstack_cms_alarm_contact":                      resourceAlibabacloudStackCmsAlarmContact(),
		"alibabacloudstack_cloudmonitorservice_alarmcontact":       resourceAlibabacloudStackCmsAlarmContact(),
		"alibabacloudstack_cms_alarm_contact_group":                resourceAlibabacloudStackCmsAlarmContactGroup(),
		"alibabacloudstack_cloudmonitorservice_alarmcontactgroup":  resourceAlibabacloudStackCmsAlarmContactGroup(),
		"alibabacloudstack_cms_metric_rule_template":               resourceAlibabacloudStackCmsMetricRuleTemplate(),
		"alibabacloudstack_cloudmonitorservice_metricruletemplate": resourceAlibabacloudStackCmsMetricRuleTemplate(),
		"alibabacloudstack_cms_site_monitor":                       resourceAlibabacloudStackCmsSiteMonitor(),
		"alibabacloudstack_cloudmonitorservice_sitemonitor":        resourceAlibabacloudStackCmsSiteMonitor(),
		"alibabacloudstack_common_bandwidth_package":               resourceAlibabacloudStackCommonBandwidthPackage(),
		"alibabacloudstack_cbwp_commonbandwidthpackage":            resourceAlibabacloudStackCommonBandwidthPackage(),
		"alibabacloudstack_common_bandwidth_package_attachment":    resourceAlibabacloudStackCommonBandwidthPackageAttachment(),
		"alibabacloudstack_cbwp_commonbandwidthpackageattachment":  resourceAlibabacloudStackCommonBandwidthPackageAttachment(),
		"alibabacloudstack_cr_ee_namespace":                        resourceAlibabacloudStackCrEeNamespace(),
		"alibabacloudstack_cr_ee_repo":                             resourceAlibabacloudStackCrEeRepo(),
		"alibabacloudstack_cr_ee_sync_rule":                        resourceAlibabacloudStackCrEeSyncRule(),
		"alibabacloudstack_cr_namespace":                           resourceAlibabacloudStackCRNamespace(),
		"alibabacloudstack_cr_repo":                                resourceAlibabacloudStackCRRepo(),
		"alibabacloudstack_cr_repository":                          resourceAlibabacloudStackCRRepo(),
		"alibabacloudstack_cs_kubernetes":                          resourceAlibabacloudStackCSKubernetes(),
		"alibabacloudstack_ack_cluster":                            resourceAlibabacloudStackCSKubernetes(),
		"alibabacloudstack_ack_template":                           resourceAlibabacloudStackAckTemplate(),
		"alibabacloudstack_cs_kubernetes_node_pool":                resourceAlibabacloudStackCSKubernetesNodePool(),
		"alibabacloudstack_datahub_project":                        resourceAlibabacloudStackDatahubProject(),
		"alibabacloudstack_datahub_subscription":                   resourceAlibabacloudStackDatahubSubscription(),
		"alibabacloudstack_datahub_topic":                          resourceAlibabacloudStackDatahubTopic(),
		"alibabacloudstack_datahub_kafka_group":                    resourceAlibabacloudStackDatahubKafkaGroup(),
		"alibabacloudstack_db_account":                             resourceAlibabacloudStackDBAccount(),
		"alibabacloudstack_rds_account":                            resourceAlibabacloudStackDBAccount(),
		"alibabacloudstack_db_account_privilege":                   resourceAlibabacloudStackDBAccountPrivilege(),
		"alibabacloudstack_db_backup_policy":                       resourceAlibabacloudStackDBBackupPolicy(),
		"alibabacloudstack_rds_backuppolicy":                       resourceAlibabacloudStackDBBackupPolicy(),
		"alibabacloudstack_db_connection":                          resourceAlibabacloudStackDBConnection(),
		"alibabacloudstack_rds_dbinstance":                         resourceAlibabacloudStackDBInstance(),
		"alibabacloudstack_db_database":                            resourceAlibabacloudStackDBDatabase(),
		"alibabacloudstack_rds_database":                           resourceAlibabacloudStackDBDatabase(),
		"alibabacloudstack_db_instance":                            resourceAlibabacloudStackDBInstance(),
		"alibabacloudstack_db_read_write_splitting_connection":     resourceAlibabacloudStackDBReadWriteSplittingConnection(),
		"alibabacloudstack_db_readonly_instance":                   resourceAlibabacloudStackDBReadonlyInstance(),
		"alibabacloudstack_rds_backup":                             resourceAlibabacloudStackRdsBackup(),
		"alibabacloudstack_disk":                                   resourceAlibabacloudStackDisk(),
		"alibabacloudstack_ecs_disk":                               resourceAlibabacloudStackDisk(),
		"alibabacloudstack_disk_attachment":                        resourceAlibabacloudStackDiskAttachment(),
		"alibabacloudstack_ecs_diskattachment":                     resourceAlibabacloudStackDiskAttachment(),
		"alibabacloudstack_dms_enterprise_instance":                resourceAlibabacloudStackDmsEnterpriseInstance(),
		"alibabacloudstack_dmsenterprise_instance":                 resourceAlibabacloudStackDmsEnterpriseInstance(),
		"alibabacloudstack_dms_enterprise_user":                    resourceAlibabacloudStackDmsEnterpriseUser(),
		"alibabacloudstack_dmsenterprise_user":                     resourceAlibabacloudStackDmsEnterpriseUser(),
		"alibabacloudstack_dns_domain":                             resourceAlibabacloudStackDnsDomain(),
		"alibabacloudstack_dns_domain_attachment":                  resourceAlibabacloudStackDnsDomainAttachment(),
		"alibabacloudstack_alidns_domainattachment":                resourceAlibabacloudStackDnsDomainAttachment(),
		"alibabacloudstack_dns_record":                             resourceAlibabacloudStackDnsRecord(),
		"alibabacloudstack_drds_instance":                          resourceAlibabacloudStackDrdsInstance(),
		"alibabacloudstack_drds_readonly_instance":                 resourceAlibabacloudStackDrdsReadonlyInstance(),
		"alibabacloudstack_drds_database":                          resourceAlibabacloudStackDrdsDatabase(),
		"alibabacloudstack_drds_account":                           resourceAlibabacloudStackDrdsAccount(),
		"alibabacloudstack_drds_rds_instance":                      resourceAlibabacloudStackDrdsRdsInstnace(),
		"alibabacloudstack_dts_subscription_job":                   resourceAlibabacloudStackDtsSubscriptionJob(),
		"alibabacloudstack_dts_subscriptionjob":                    resourceAlibabacloudStackDtsSubscriptionJob(),
		"alibabacloudstack_dts_synchronization_instance":           resourceAlibabacloudStackDtsSynchronizationInstance(),
		"alibabacloudstack_dts_synchronizationinstance":            resourceAlibabacloudStackDtsSynchronizationInstance(),
		"alibabacloudstack_dts_synchronization_job":                resourceAlibabacloudStackDtsSynchronizationJob(),
		"alibabacloudstack_ecs_command":                            resourceAlibabacloudStackEcsCommand(),
		"alibabacloudstack_ecs_dedicated_host":                     resourceAlibabacloudStackEcsDedicatedHost(),
		"alibabacloudstack_ecs_dedicatedhost":                      resourceAlibabacloudStackEcsDedicatedHost(),
		"alibabacloudstack_ecs_dedicated_host_cluster":             resourceAlibabacloudStackEcsDedicatedhostcluster(),
		"alibabacloudstack_ecs_deployment_set":                     resourceAlibabacloudStackEcsDeploymentSet(),
		"alibabacloudstack_ecs_deploymentset":                      resourceAlibabacloudStackEcsDeploymentSet(),
		"alibabacloudstack_ecs_hpc_cluster":                        resourceAlibabacloudStackEcsHpcCluster(),
		"alibabacloudstack_ecs_hpccluster":                         resourceAlibabacloudStackEcsHpcCluster(),
		"alibabacloudstack_ecs_ebs_storage_set":                    resourceAlibabacloudStackEcsEbsStorageSets(),
		"alibabacloudstack_ecs_storageset":                         resourceAlibabacloudStackEcsEbsStorageSets(),
		"alibabacloudstack_edas_application":                       resourceAlibabacloudStackEdasApplication(),
		"alibabacloudstack_edas_k8s_service":                       resourceAlibabacloudStackEdasK8sService(),
		"alibabacloudstack_edas_slbattachment":                     resourceAlibabacloudStackEdasSlbAttachment(),
		"alibabacloudstack_edas_application_scale":                 resourceAlibabacloudStackEdasInstanceApplicationAttachment(),
		"alibabacloudstack_edas_cluster":                           resourceAlibabacloudStackEdasCluster(),
		"alibabacloudstack_edas_deploy_group":                      resourceAlibabacloudStackEdasDeployGroup(),
		"alibabacloudstack_edas_deploygroup":                       resourceAlibabacloudStackEdasDeployGroup(),
		"alibabacloudstack_edas_instance_cluster_attachment":       resourceAlibabacloudStackEdasInstanceClusterAttachment(),
		"alibabacloudstack_edas_instanceclusterattachment":         resourceAlibabacloudStackEdasInstanceClusterAttachment(),
		"alibabacloudstack_edas_k8s_application":                   resourceAlibabacloudStackEdasK8sApplication(),
		"alibabacloudstack_edas_k8s_cluster":                       resourceAlibabacloudStackEdasK8sCluster(),
		"alibabacloudstack_edas_namespace":                         resourceAlibabacloudStackEdasNamespace(),
		"alibabacloudstack_edas_slb_attachment":                    resourceAlibabacloudStackEdasSlbAttachment(),
		"alibabacloudstack_ehpc_job_template":                      resourceAlibabacloudStackEhpcJobTemplate(),
		"alibabacloudstack_eip":                                    resourceAlibabacloudStackEip(),
		"alibabacloudstack_eip_address":                            resourceAlibabacloudStackEip(),
		"alibabacloudstack_eip_association":                        resourceAlibabacloudStackEipAssociation(),
		"alibabacloudstack_ess_alarm":                              resourceAlibabacloudStackEssAlarm(),
		"alibabacloudstack_autoscaling_alarmtask":                  resourceAlibabacloudStackEssAlarm(),
		"alibabacloudstack_ess_attachment":                         resourceAlibabacloudStackEssAttachment(),
		"alibabacloudstack_ess_lifecycle_hook":                     resourceAlibabacloudStackEssLifecycleHook(),
		"alibabacloudstack_autoscaling_lifecyclehook":              resourceAlibabacloudStackEssLifecycleHook(),
		"alibabacloudstack_ess_notification":                       resourceAlibabacloudStackEssNotification(),
		"alibabacloudstack_autoscaling_notification":               resourceAlibabacloudStackEssNotification(),
		"alibabacloudstack_ess_scaling_group":                      resourceAlibabacloudStackEssScalingGroup(),
		"alibabacloudstack_autoscaling_scalingrule":                resourceAlibabacloudStackEssScalingRule(),
		"alibabacloudstack_ess_scaling_rule":                       resourceAlibabacloudStackEssScalingRule(),
		"alibabacloudstack_ess_scalinggroup_vserver_groups":        resourceAlibabacloudStackEssScalingGroupVserverGroups(),
		"alibabacloudstack_ess_scheduled_task":                     resourceAlibabacloudStackEssScheduledTask(),
		"alibabacloudstack_autoscaling_scheduled_task":             resourceAlibabacloudStackEssScheduledTask(),
		"alibabacloudstack_autoscaling_scheduledtask":              resourceAlibabacloudStackEssScheduledTask(),
		"alibabacloudstack_forward_entry":                          resourceAlibabacloudStackForwardEntry(),
		"alibabacloudstack_natgateway_forwardentry":                resourceAlibabacloudStackForwardEntry(),
		"alibabacloudstack_natgateway_bandwidth_package":           resourceAlibabacloudStackNatgatewayBandwidthPackage(),
		"alibabacloudstack_gpdb_account":                           resourceAlibabacloudStackGpdbAccount(),
		"alibabacloudstack_gpdb_connection":                        resourceAlibabacloudStackGpdbConnection(),
		"alibabacloudstack_gpdb_publicconnection":                  resourceAlibabacloudStackGpdbConnection(),
		"alibabacloudstack_gpdb_instance":                          resourceAlibabacloudStackGpdbInstance(),
		"alibabacloudstack_gpdb_dbinstance":                        resourceAlibabacloudStackGpdbInstance(),
		"alibabacloudstack_hbase_instance":                         resourceAlibabacloudStackHBaseInstance(),
		"alibabacloudstack_hbase_cluster":                          resourceAlibabacloudStackHBaseInstance(),
		"alibabacloudstack_image":                                  resourceAlibabacloudStackImage(),
		"alibabacloudstack_ecs_image":                              resourceAlibabacloudStackImage(),
		"alibabacloudstack_image_copy":                             resourceAlibabacloudStackImageCopy(),
		"alibabacloudstack_image_export":                           resourceAlibabacloudStackImageExport(),
		"alibabacloudstack_image_import":                           resourceAlibabacloudStackImageImport(),
		"alibabacloudstack_image_share_permission":                 resourceAlibabacloudStackImageSharePermission(),
		"alibabacloudstack_instance":                               resourceAlibabacloudStackInstance(),
		"alibabacloudstack_ecs_instance":                           resourceAlibabacloudStackInstance(),
		"alibabacloudstack_key_pair":                               resourceAlibabacloudStackKeyPair(),
		"alibabacloudstack_ecs_keypair":                            resourceAlibabacloudStackKeyPair(),
		"alibabacloudstack_key_pair_attachment":                    resourceAlibabacloudStackKeyPairAttachment(),
		"alibabacloudstack_ecs_keypairattachment":                  resourceAlibabacloudStackKeyPairAttachment(),
		"alibabacloudstack_kms_alias":                              resourceAlibabacloudStackKmsAlias(),
		"alibabacloudstack_kms_ciphertext":                         resourceAlibabacloudStackKmsCiphertext(),
		"alibabacloudstack_kms_key":                                resourceAlibabacloudStackKmsKey(),
		"alibabacloudstack_kms_secret":                             resourceAlibabacloudStackKmsSecret(),
		"alibabacloudstack_kvstore_account":                        resourceAlibabacloudStackKVstoreAccount(),
		"alibabacloudstack_redis_account":                          resourceAlibabacloudStackKVstoreAccount(),
		"alibabacloudstack_kvstore_backup_policy":                  resourceAlibabacloudStackKVStoreBackupPolicy(),
		"alibabacloudstack_kvstore_connection":                     resourceAlibabacloudStackKvstoreConnection(),
		"alibabacloudstack_redis_connection":                       resourceAlibabacloudStackKvstoreConnection(),
		"alibabacloudstack_kvstore_instance":                       resourceAlibabacloudStackKVStoreInstance(),
		"alibabacloudstack_redis_tairinstance":                     resourceAlibabacloudStackKVStoreInstance(),
		// This resource is not yet supported by the private cloud frontend
		// "alibabacloudstack_launch_template":                        resourceAlibabacloudStackLaunchTemplate(),
		// "alibabacloudstack_ecs_launchtemplate":                     resourceAlibabacloudStackLaunchTemplate(),
		"alibabacloudstack_log_alert":                                       resourceAlibabacloudStackLogAlert(),
		"alibabacloudstack_log_machine_group":                               resourceAlibabacloudStackLogMachineGroup(),
		"alibabacloudstack_log_project":                                     resourceAlibabacloudStackLogProject(),
		"alibabacloudstack_log_store":                                       resourceAlibabacloudStackLogStore(),
		"alibabacloudstack_log_store_index":                                 resourceAlibabacloudStackLogStoreIndex(),
		"alibabacloudstack_logtail_attachment":                              resourceAlibabacloudStackLogtailAttachment(),
		"alibabacloudstack_logtail_config":                                  resourceAlibabacloudStackLogtailConfig(),
		"alibabacloudstack_maxcompute_project":                              resourceAlibabacloudStackMaxcomputeProject(),
		"alibabacloudstack_maxcompute_user":                                 resourceAlibabacloudStackMaxcomputeUser(),
		"alibabacloudstack_maxcompute_cu":                                   resourceAlibabacloudStackMaxcomputeCu(),
		"alibabacloudstack_mongodb_instance":                                resourceAlibabacloudStackMongoDBInstance(),
		"alibabacloudstack_mongodb_account":                                 resourceAlibabacloudStackMongodbAccount(),
		"alibabacloudstack_mongodb_sharding_instance":                       resourceAlibabacloudStackMongoDBShardingInstance(),
		"alibabacloudstack_mongodb_shardinginstance":                        resourceAlibabacloudStackMongoDBShardingInstance(),
		"alibabacloudstack_mongodb_shardinginstance_csnode_address":         resourceAlibabacloudStackMongodbShardingInstanceCsNodeAddress(),
		"alibabacloudstack_mongodb_shardinginstance_shardnode_address":      resourceAlibabacloudStackMongodbShardingInstanceShardNodeAddress(),
		"alibabacloudstack_mongodb_shardinginstance_mongosnode_address":     resourceAlibabacloudStackMongodbShardingInstanceMongosNodeAddress(),
		"alibabacloudstack_mongodb_backup":                                  resourceAlibabacloudStackMongodbBackup(),
		"alibabacloudstack_nas_access_group":                                resourceAlibabacloudStackNasAccessGroup(),
		"alibabacloudstack_nas_accessgroup":                                 resourceAlibabacloudStackNasAccessGroup(),
		"alibabacloudstack_nas_access_rule":                                 resourceAlibabacloudStackNasAccessRule(),
		"alibabacloudstack_nas_accessrule":                                  resourceAlibabacloudStackNasAccessRule(),
		"alibabacloudstack_nas_file_system":                                 resourceAlibabacloudStackNasFileSystem(),
		"alibabacloudstack_nas_filesystem":                                  resourceAlibabacloudStackNasFileSystem(),
		"alibabacloudstack_nas_lifecycle_policy":                            resourceAlibabacloudStackNasLifecyclepolicy(),
		"alibabacloudstack_nas_mount_target":                                resourceAlibabacloudStackNasMountTarget(),
		"alibabacloudstack_nas_mounttarget":                                 resourceAlibabacloudStackNasMountTarget(),
		"alibabacloudstack_cpfs_file_system":                                resourceAlibabacloudStackCpfsFileSystem(),
		"alibabacloudstack_nat_gateway":                                     resourceAlibabacloudStackNatGateway(),
		"alibabacloudstack_natgateway_natgateway":                           resourceAlibabacloudStackNatGateway(),
		"alibabacloudstack_network_acl":                                     resourceAlibabacloudStackNetworkAcl(),
		"alibabacloudstack_vpc_networkacl":                                  resourceAlibabacloudStackNetworkAcl(),
		"alibabacloudstack_vpc_network_acl_attachment":                      resourceAlibabacloudStackNetworkAclAttachment(),
		"alibabacloudstack_vpc_vswitch_network_acl_attachment":              resourceAlibabacloudStackNetworkAclAttachment(),
		"alibabacloudstack_network_acl_attachment":                          resourceAlibabacloudStackNetworkAclAttachment(),
		"alibabacloudstack_vpc_network_acl_entries":                         resourceAlibabacloudStackNetworkAclEntries(),
		"alibabacloudstack_network_acl_entries":                             resourceAlibabacloudStackNetworkAclEntries(),
		"alibabacloudstack_network_interface":                               resourceAlibabacloudStackNetworkInterface(),
		"alibabacloudstack_ecs_networkinterface":                            resourceAlibabacloudStackNetworkInterface(),
		"alibabacloudstack_network_interface_attachment":                    resourceAlibabacloudStackNetworkInterfaceAttachment(),
		"alibabacloudstack_ecs_networkinterfaceattachment":                  resourceAlibabacloudStackNetworkInterfaceAttachment(),
		"alibabacloudstack_ons_group":                                       resourceAlibabacloudStackOnsGroup(),
		"alibabacloudstack_ons_instance":                                    resourceAlibabacloudStackOnsInstance(),
		"alibabacloudstack_ons_topic":                                       resourceAlibabacloudStackOnsTopic(),
		"alibabacloudstack_oss_bucket":                                      resourceAlibabacloudStackOssBucket(),
		"alibabacloudstack_oss_bucket_quota":                                resourceAlibabacloudStackOssBucketQuota(),
		"alibabacloudstack_oss_bucket_kms":                                  resourceAlibabacloudStackOssBucketKms(),
		"alibabacloudstack_oss_bucket_object":                               resourceAlibabacloudStackOssBucketObject(),
		"alibabacloudstack_ots_instance":                                    resourceAlibabacloudStackOtsInstance(),
		"alibabacloudstack_ots_instance_attachment":                         resourceAlibabacloudStackOtsInstanceAttachment(),
		"alibabacloudstack_ots_instanceattachment":                          resourceAlibabacloudStackOtsInstanceAttachment(),
		"alibabacloudstack_ots_table":                                       resourceAlibabacloudStackOtsTable(),
		"alibabacloudstack_quick_bi_user":                                   resourceAlibabacloudStackQuickBiUser(),
		"alibabacloudstack_quick_bi_user_group":                             resourceAlibabacloudStackQuickBiUserGroup(),
		"alibabacloudstack_quick_bi_workspace":                              resourceAlibabacloudStackQuickBiWorkspace(),
		"alibabacloudstack_ram_role_attachment":                             resourceAlibabacloudStackRamRoleAttachment(),
		"alibabacloudstack_ecs_ramroleattachment":                           resourceAlibabacloudStackRamRoleAttachment(),
		"alibabacloudstack_reserved_instance":                               resourceAlibabacloudStackReservedInstance(),
		"alibabacloudstack_ecs_reservedinstance":                            resourceAlibabacloudStackReservedInstance(),
		"alibabacloudstack_ros_stack":                                       resourceAlibabacloudStackRosStack(),
		"alibabacloudstack_ros_template":                                    resourceAlibabacloudStackRosTemplate(),
		"alibabacloudstack_route_entry":                                     resourceAlibabacloudStackRouteEntry(),
		"alibabacloudstack_route_table":                                     resourceAlibabacloudStackRouteTable(),
		"alibabacloudstack_vpc_routetable":                                  resourceAlibabacloudStackRouteTable(),
		"alibabacloudstack_route_table_attachment":                          resourceAlibabacloudStackRouteTableAttachment(),
		"alibabacloudstack_vpc_routetableattachment":                        resourceAlibabacloudStackRouteTableAttachment(),
		"alibabacloudstack_router_interface":                                resourceAlibabacloudStackRouterInterface(),
		"alibabacloudstack_expressconnect_routerinterface":                  resourceAlibabacloudStackRouterInterface(),
		"alibabacloudstack_router_interface_connection":                     resourceAlibabacloudStackRouterInterfaceConnection(),
		"alibabacloudstack_security_group":                                  resourceAlibabacloudStackSecurityGroup(),
		"alibabacloudstack_ecs_securitygroup":                               resourceAlibabacloudStackSecurityGroup(),
		"alibabacloudstack_security_group_rule":                             resourceAlibabacloudStackSecurityGroupRule(),
		"alibabacloudstack_slb":                                             resourceAlibabacloudStackSlb(),
		"alibabacloudstack_slb_loadbalancer":                                resourceAlibabacloudStackSlb(),
		"alibabacloudstack_slb_acl":                                         resourceAlibabacloudStackSlbAcl(),
		"alibabacloudstack_slb_accesscontrollist":                           resourceAlibabacloudStackSlbAcl(),
		"alibabacloudstack_slb_backend_server":                              resourceAlibabacloudStackSlbBackendServer(),
		"alibabacloudstack_slb_backendserver":                               resourceAlibabacloudStackSlbBackendServer(),
		"alibabacloudstack_slb_ca_certificate":                              resourceAlibabacloudStackSlbCACertificate(),
		"alibabacloudstack_slb_cacertificate":                               resourceAlibabacloudStackSlbCACertificate(),
		"alibabacloudstack_slb_domain_extension":                            resourceAlibabacloudStackSlbDomainExtension(),
		"alibabacloudstack_slb_domainextension":                             resourceAlibabacloudStackSlbDomainExtension(),
		"alibabacloudstack_slb_listener":                                    resourceAlibabacloudStackSlbListener(),
		"alibabacloudstack_slb_master_slave_server_group":                   resourceAlibabacloudStackSlbMasterSlaveServerGroup(),
		"alibabacloudstack_slb_masterslaveservergroup":                      resourceAlibabacloudStackSlbMasterSlaveServerGroup(),
		"alibabacloudstack_slb_rule":                                        resourceAlibabacloudStackSlbRule(),
		"alibabacloudstack_slb_server_certificate":                          resourceAlibabacloudStackSlbServerCertificate(),
		"alibabacloudstack_slb_servercertificate":                           resourceAlibabacloudStackSlbServerCertificate(),
		"alibabacloudstack_slb_server_group":                                resourceAlibabacloudStackSlbServerGroup(),
		"alibabacloudstack_slb_vservergroup":                                resourceAlibabacloudStackSlbServerGroup(),
		"alibabacloudstack_snapshot":                                        resourceAlibabacloudStackSnapshot(),
		"alibabacloudstack_ecs_snapshot":                                    resourceAlibabacloudStackSnapshot(),
		"alibabacloudstack_snapshot_policy":                                 resourceAlibabacloudStackSnapshotPolicy(),
		"alibabacloudstack_ecs_autosnapshotpolicy":                          resourceAlibabacloudStackSnapshotPolicy(),
		"alibabacloudstack_snat_entry":                                      resourceAlibabacloudStackSnatEntry(),
		"alibabacloudstack_natgateway_snatentry":                            resourceAlibabacloudStackSnatEntry(),
		"alibabacloudstack_vpc":                                             resourceAlibabacloudStackVpc(),
		"alibabacloudstack_vpc_vpc":                                         resourceAlibabacloudStackVpc(),
		"alibabacloudstack_vpc_ipv6_egress_rule":                            resourceAlibabacloudStackVpcIpv6EgressRule(),
		"alibabacloudstack_vpc_ipv6egressrule":                              resourceAlibabacloudStackVpcIpv6EgressRule(),
		"alibabacloudstack_vpc_ipv6_gateway":                                resourceAlibabacloudStackVpcIpv6Gateway(),
		"alibabacloudstack_vpc_ipv6gateway":                                 resourceAlibabacloudStackVpcIpv6Gateway(),
		"alibabacloudstack_vpc_ipv6_internet_bandwidth":                     resourceAlibabacloudStackVpcIpv6InternetBandwidth(),
		"alibabacloudstack_vpc_ipv6internetbandwidth":                       resourceAlibabacloudStackVpcIpv6InternetBandwidth(),
		"alibabacloudstack_vpn_connection":                                  resourceAlibabacloudStackVpnConnection(),
		"alibabacloudstack_vpngateway_vpnconnection":                        resourceAlibabacloudStackVpnConnection(),
		"alibabacloudstack_vpn_customer_gateway":                            resourceAlibabacloudStackVpnCustomerGateway(),
		"alibabacloudstack_vpngateway_customergateway":                      resourceAlibabacloudStackVpnCustomerGateway(),
		"alibabacloudstack_vpn_gateway":                                     resourceAlibabacloudStackVpnGateway(),
		"alibabacloudstack_vpngateway_vpngateway":                           resourceAlibabacloudStackVpnGateway(),
		"alibabacloudstack_vpn_route_entry":                                 resourceAlibabacloudStackVpnRouteEntry(),
		"alibabacloudstack_vpngateway_vpnrouteentry":                        resourceAlibabacloudStackVpnRouteEntry(),
		"alibabacloudstack_vpngateway_ssl_vpn_server":                       resourceAlibabacloudStackVpngatewaySslvpnserver(),
		"alibabacloudstack_vpngateway_ssl_vpnserver":                        resourceAlibabacloudStackVpngatewaySslvpnserver(),
		"alibabacloudstack_vpngateway_ssl_vpn_client_cert":                  resourceAlibabacloudStackVpngatewaySslvpnclientcert(),
		"alibabacloudstack_vpngateway_sslvpnclientcert":                     resourceAlibabacloudStackVpngatewaySslvpnclientcert(),
		"alibabacloudstack_vswitch":                                         resourceAlibabacloudStackSwitch(),
		"alibabacloudstack_vpc_vswitch":                                     resourceAlibabacloudStackSwitch(),
		"alibabacloudstack_data_works_folder":                               resourceAlibabacloudStackDataWorksFolder(),
		"alibabacloudstack_data_works_connection":                           resourceAlibabacloudStackDataWorksConnection(),
		"alibabacloudstack_data_works_user":                                 resourceAlibabacloudStackDataWorksUser(),
		"alibabacloudstack_data_works_project":                              resourceAlibabacloudStackDataWorksProject(),
		"alibabacloudstack_data_works_user_role_binding":                    resourceAlibabacloudStackDataWorksUserRoleBinding(),
		"alibabacloudstack_data_works_remind":                               resourceAlibabacloudStackDataWorksRemind(),
		"alibabacloudstack_elasticsearch_instance":                          resourceAlibabacloudStackElasticsearch(),
		"alibabacloudstack_dbs_backup_plan":                                 resourceAlibabacloudStackDbsBackupPlan(),
		"alibabacloudstack_dbs_backupplan":                                  resourceAlibabacloudStackDbsBackupPlan(),
		"alibabacloudstack_express_connect_physical_connection":             resourceAlibabacloudStackExpressConnectPhysicalConnection(),
		"alibabacloudstack_expressconnect_physicalconnection":               resourceAlibabacloudStackExpressConnectPhysicalConnection(),
		"alibabacloudstack_expressconnect_vbr_ha":                           resourceAlibabacloudStackExpressconnectVbrHa(),
		"alibabacloudstack_expressconnect_vbr_pconn_association":            resourceAlibabacloudStackExpressconnectVbrpconnassociation(),
		"alibabacloudstack_express_connect_virtual_border_router":           resourceAlibabacloudStackExpressConnectVirtualBorderRouter(),
		"alibabacloudstack_expressconnect_virtualborderrouter":              resourceAlibabacloudStackExpressConnectVirtualBorderRouter(),
		"alibabacloudstack_expressconnect_bgp_network":                      resourceAlibabacloudStackExpressconnectBgpnetwork(),
		"alibabacloudstack_oos_template":                                    resourceAlibabacloudStackOosTemplate(),
		"alibabacloudstack_oos_execution":                                   resourceAlibabacloudStackOosExecution(),
		"alibabacloudstack_arms_alert_contact":                              resourceAlibabacloudStackArmsAlertContact(),
		"alibabacloudstack_arms_alertcontact":                               resourceAlibabacloudStackArmsAlertContact(),
		"alibabacloudstack_arms_alert_contact_group":                        resourceAlibabacloudStackArmsAlertContactGroup(),
		"alibabacloudstack_arms_alertcontactgroup":                          resourceAlibabacloudStackArmsAlertContactGroup(),
		"alibabacloudstack_arms_dispatch_rule":                              resourceAlibabacloudStackArmsDispatchRule(),
		"alibabacloudstack_arms_dispatchrule":                               resourceAlibabacloudStackArmsDispatchRule(),
		"alibabacloudstack_arms_prometheus_alert_rule":                      resourceAlibabacloudStackArmsPrometheusAlertRule(),
		"alibabacloudstack_arms_prometheusalertrule":                        resourceAlibabacloudStackArmsPrometheusAlertRule(),
		"alibabacloudstack_elasticsearch_k8s_instance":                      resourceAlibabacloudStackElasticsearch(),
		"alibabacloudstack_cloud_firewall_control_policy":                   resourceAlibabacloudStackCloudFirewallControlPolicy(),
		"alibabacloudstack_cloudfirewall_controlpolicy":                     resourceAlibabacloudStackCloudFirewallControlPolicy(),
		"alibabacloudstack_cloud_firewall_control_policy_order":             resourceAlibabacloudStackCloudFirewallControlPolicyOrder(),
		"alibabacloudstack_cloudfirewall_controlpolicyorder":                resourceAlibabacloudStackCloudFirewallControlPolicyOrder(),
		"alibabacloudstack_csb_project":                                     resourceAlibabacloudStackCsbProject(),
		"alibabacloudstack_graph_database_db_instance":                      resourceAlibabacloudStackGraphDatabaseDbInstance(),
		"alibabacloudstack_graphdatabase_dbinstance":                        resourceAlibabacloudStackGraphDatabaseDbInstance(),
		"alibabacloudstack_polardb_account":                                 resourceAlibabacloudStackPolardbAccount(),
		"alibabacloudstack_polardb_database":                                resourceAlibabacloudStackPolardbDatabase(),
		"alibabacloudstack_polardb_backuppolicy":                            resourceAlibabacloudStackPolardbBackuppolicy(),
		"alibabacloudstack_polardb_dbconnection":                            resourceAlibabacloudStackPolardbConnection(),
		"alibabacloudstack_polardb_dbinstance":                              resourceAlibabacloudStackPolardbInstance(),
		"alibabacloudstack_polardb_account_database_binding":                resourceAlibabacloudStackPolardbAccountDatabaseBinding(),
		"alibabacloudstack_acm_configuration":                               resourceAlibabacloudStackAcmConfiguration(),
		"alibabacloudstack_polardb_backup":                                  resourceAlibabacloudStackPolardbBackup(),
		"alibabacloudstack_polardb_readonly_instance":                       resourceAlibabacloudStackPolardbReadonlyInstance(),
		"alibabacloudstack_polardb_proxy":                                   resourceAlibabacloudStackPolardbproxy(),
		"alibabacloudstack_bastionhost_instance":                            resourceAlibabacloudStackBastionhostInstance(),
		"alibabacloudstack_waf_instance":                                    resourceAlibabacloudstackWafInstance(),
		"alibabacloudstack_flink_namespace":                                 resourceAlibabacloudStackFlinkNamespace(),
		"alibabacloudstack_ebs_diskreplicapair":                             resourceAlibabacloudStackEbsDiskreplicapair(),
		"alibabacloudstack_ebs_diskreplicagroup":                            resourceAlibabacloudStackEbsDiskreplicagroup(),
		"alibabacloudstack_ecs_snapshot_group":                              resourceAlibabacloudStackEcsSnapshotgroup(),
		"alibabacloudstack_ecs_invocation":                                  resourceAlibabacloudStackEcsInvocation(),
		"alibabacloudstack_vpc_dhcp_options_set":                            resourceAlibabacloudStackVpcDhcpoptionsset(),
		"alibabacloudstack_slb_access_log":                                  resourceAlibabacloudStackSlbAccesslog(),
		"alibabacloudstack_vpc_ha_vip":                                      resourceAlibabacloudStackVpcHavip(),
		"alibabacloudstack_vpngateway_vpn_pbr_route_entry":                  resourceAlibabacloudStackVpngatewayVpnpbrrouteentry(),
		"alibabacloudstack_expressconnect_bgp_group":                        resourceAlibabacloudStackExpressconnectBgpgroup(),
		"alibabacloudstack_expressconnect_bgp_peer":                         resourceAlibabacloudStackExpressconnectBgppeer(),
		"alibabacloudstack_db_proxy":                                        resourceAlibabacloudStackRdsDbproxy(),
		"alibabacloudstack_polardbx_instance":                               resourceAlibabacloudStackPolardbxInstance(),
		"alibabacloudstack_polardbx_readonly_instance":                      resourceAlibabacloudStackPolardbxReadonlyInstance(),
		"alibabacloudstack_polardbx_account":                                resourceAlibabacloudStackPolardbxAccount(),
		"alibabacloudstack_polardbx_super_account":                          resourceAlibabacloudStackPolardbxSuperAccount(),
		"alibabacloudstack_polardbx_database":                               resourceAlibabacloudStackPolardbxDatabase(),
		"alibabacloudstack_polardbx_account_database_binding":               resourceAlibabacloudStackPolardbxAccountDatabaseBinding(),
		"alibabacloudstack_polardbx_backup":                                 resourceAlibabacloudStackPolardbxBackup(),
		"alibabacloudstack_polardbx_backup_policy":                          resourceAlibabacloudStackPolardbxBackupPolicy(),
		"alibabacloudstack_polardbx_log_engine":                             resourceAlibabacloudStackPolardbxLogEngine(),
		"alibabacloudstack_polardbx_read_write_splitting_config":            resourceAlibabacloudStackPolardbxReadWriteSplittingConfig(),
		"alibabacloudstack_cen_instance":                                    resourceAlibabacloudStackCenCeninstance(),
		"alibabacloudstack_cen_transit_router_route_table":                  resourceAlibabacloudStackCenTransitRouterRouterTable(),
		"alibabacloudstack_cen_transit_router_route_entry":                  resourceAlibabacloudStackCenTransitRouterRouteEntry(),
		"alibabacloudstack_cen_transit_router_vpc_attachment":               resourceAlibabacloudStackCenTransitRouterVpcAttachment(),
		"alibabacloudstack_cen_transit_router_vbr_attachment":               resourceAlibabacloudStackCenTransitRouterVbrAttachment(),
		"alibabacloudstack_cen_transit_router_route_table_association":      resourceAlibabacloudStackCenTransitRouterRouterTableAssociation(),
		"alibabacloudstack_cen_transit_router_route_table_propagation":      resourceAlibabacloudStackCenTransitRouterRouterTablePropagation(),
		"alibabacloudstack_cen_route_map":                                   resourceAlibabacloudStackCenRouteMap(),
		"alibabacloudstack_cen_transit_router_multicast_domain":             resourceAlibabacloudStackCenTransitMulticastDomain(),
		"alibabacloudstack_cen_transit_router_multicast_domain_association": resourceAlibabacloudStackCenTransitMulticastDomainAssociation(),
		"alibabacloudstack_cen_transit_router_multicast_domain_source":      resourceAlibabacloudStackCenTransitMulticastDomainSource(),
		"alibabacloudstack_cen_transit_router_multicast_domain_member":      resourceAlibabacloudStackCenTransitMulticastDomainMember(),
		"alibabacloudstack_polardb_cluster_instance":                        resourceAlibabacloudStackPolardbClusterInstance(),
		"alibabacloudstack_polardb_cluster_account":                         resourceAlibabacloudStackPolardbClusterAccount(),
		"alibabacloudstack_polardb_cluster_database":                        resourceAlibabacloudStackPolardbClusterDatabase(),
		"alibabacloudstack_polardb_cluster_account_database_binding":        resourceAlibabacloudStackPolardbClusterAccountDatabaseBinding(),
		"alibabacloudstack_polardb_cluster_proxy":                           resourceAlibabacloudStackPolardbClusterProxy(),
		"alibabacloudstack_polardb_cluster_backup_policy":                   resourceAlibabacloudStackPolardbClusterBackupPolicy(),
		"alibabacloudstack_cen_transit_router_connect_attachment":           resourceAlibabacloudStackCenTransitRouterConnectAttachment(),
		"alibabacloudstack_cen_vbr_health_check":                            resourceAlibabacloudStackCenVbrHealthCheck(),
		"alibabacloudstack_cen_transit_router_connect_peer":                 resourceAlibabacloudStackCenTransitRouterConnectPeer(),
		"alibabacloudstack_edas_k8s_application_scaling_rule":               resourceAlibabacloudStackEdasK8sApplicationScalingRule(),
		"alibabacloudstack_edas_swimming_lane_group":                        resourceAlibabacloudStackEdasSwimmingLaneGroup(),
		"alibabacloudstack_edas_swimming_lane":                              resourceAlibabacloudStackEdasSwimmingLane(),
		"alibabacloudstack_hologram_instance":                               resourceAlibabacloudStackHologramInstance(),
		"alibabacloudstack_hologram_instance_backup_policy":                 resourceAlibabacloudStackHologramInstanceBackupPolicy(),
		"alibabacloudstack_lindorm_instance":                                resourceAlibabacloudStackLindormInstance(),
		"alibabacloudstack_lindorm_lts_instance":                            resourceAlibabacloudStackLindormLtsInstance(),
		"alibabacloudstack_api_gateway_v2_k8s_cluster":                      resourceAlibabacloudStackAPIGatewayV2K8sCluster(),
		"alibabacloudstack_api_gateway_v2_instance":                         resourceAlibabacloudStackAPIGatewayV2Instance(),
		"alibabacloudstack_api_gateway_v2_certificate":                      resourceAlibabacloudStackAPIGatewayV2Certificate(),
		"alibabacloudstack_api_gateway_v2_domain":                           resourceAlibabacloudStackAPIGatewayV2Domain(),
		"alibabacloudstack_api_gateway_v2_signature":                        resourceAlibabacloudStackAPIGatewayV2Signature(),
		"alibabacloudstack_api_gateway_v2_consumer":                         resourceAlibabacloudStackAPIGatewayV2Consumer(),
		"alibabacloudstack_api_gateway_v2_service_source":                   resourceAlibabacloudStackAPIGatewayV2ServiceSource(),
		"alibabacloudstack_api_gateway_v2_route_group":                      resourceAlibabacloudStackApiGatewayV2RouteGroup(),
		"alibabacloudstack_api_gateway_v2_cascade_link":                     resourceAlibabacloudStackApiGatewayV2CascadeLink(),
		"alibabacloudstack_api_gateway_v2_cascade_instance":                 resourceAlibabacloudStackAPIGatewayV2CascadeInstance(),
		"alibabacloudstack_api_gateway_v2_service":                          resourceAlibabacloudStackAPIGatewayV2Service(),
		"alibabacloudstack_api_gateway_v2_route":                            resourceAlibabacloudStackApiGatewayV2Route(),
		"alibabacloudstack_api_gateway_v2_mcpserver":                        resourceAlibabacloudStackAPIGatewayV2Mcpserver(),
		"alibabacloudstack_mqtt_instance":                                   resourceAlibabacloudStackMqttInstance(),
		"alibabacloudstack_mqtt_topic":                                      resourceAlibabacloudStackMqttTopic(),
		"alibabacloudstack_mqtt_group":                                      resourceAlibabacloudStackMqttGroup(),
		"alibabacloudstack_universal_dns_domain":                            resourceAlibabacloudStackUniversalDnsDomain(),
		"alibabacloudstack_universal_dns_record":                            resourceAlibabacloudStackUniversalDnsRecord(),
		"alibabacloudstack_universal_dns_line":                              resourceAlibabacloudStackUniversalDnsLine(),
		"alibabacloudstack_dns_forward_domain":                              resourceAlibabacloudStackDnsForwardDomain(),
		"alibabacloudstack_dns_line":                                        resourceAlibabacloudStackDnsLine(),
		"alibabacloudstack_dns_recursor_acl":                                resourceAlibabacloudStackDnsRecursorAcl(),
		"alibabacloudstack_dns_private_domain":                              resourceAlibabacloudStackDnsPrivateDomain(),
		"alibabacloudstack_dns_private_record":                              resourceAlibabacloudStackDnsPrivateRecord(),
		"alibabacloudstack_dns_private_line":                                resourceAlibabacloudStackDnsPrivateLine(),
		"alibabacloudstack_dns_gtm_instance":                                resourceAlibabacloudStackDnsGtmInstance(),
		"alibabacloudstack_dns_gtm_addresspool":                             resourceAlibabacloudStackDnsGtmAddressPool(),
		"alibabacloudstack_dns_gtm_access_strategy":                         resourceAlibabacloudStackDnsGtmAccessStrategy(),
		"alibabacloudstack_prometheus_v2_instance":                          resourceAlibabacloudStackPrometheusV2Instance(),
		"alibabacloudstack_prometheus_v2_notify_group":                      resourceAlibabacloudStackPrometheusV2NotifyGroup(),
		"alibabacloudstack_prometheus_v2_contact":                           resourceAlibabacloudStackPrometheusV2Contact(),
		"alibabacloudstack_prometheus_v2_alert":                             resourceAlibabacloudStackPrometheusV2Alert(),
		"alibabacloudstack_schedulerx2_app_group":                           resourceAlibabacloudStackSchedulerx2AppGroup(),
		"alibabacloudstack_schedulerx2_job":                                 resourceAlibabacloudStackSchedulerx2Job(),
		"alibabacloudstack_schedulerx2_workflow":                            resourceAlibabacloudStackSchedulerx2Workflow(),
		"alibabacloudstack_apfs_file_system":                                resourceAlibabacloudStackApfsFileSystem(),
		"alibabacloudstack_hsm_instance":                                    resourceAlibabacloudStackHsmInstance(),
		"alibabacloudstack_hsm_cluster":                                     resourceAlibabacloudStackHsmCluster(),
		"alibabacloudstack_cspprivate_hsm_instance":                         resourceAlibabacloudStackCspprivateHsmInstance(),
		"alibabacloudstack_cspprivate_hsm_group":                            resourceAlibabacloudStackCspprivateHsmGroup(),
		"alibabacloudstack_bms_keypair":                                     resourceAlibabacloudStackBmsKeypair(),
		"alibabacloudstack_aqs_oss_scanconfig":                              resourceAlibabacloudStackAqsOssScanconfig(),
		"alibabacloudstack_aqs_anti_brute_force_rule":                       resourceAlibabacloudStackAqsAntiBruteForceRule(),
		"alibabacloudstack_aqs_web_lock":                                    resourceAlibabacloudStackAqsWebLock(),
		"alibabacloudstack_cr_ee_attestor_lifecycle_rule":                   resourceAlibabacloudStackCrEEArtifactLifecycleRule(),
		"alibabacloudstack_cloudfw_address_book":                            resourceAlibabacloudStackCloudfwAddressBook(),
		"alibabacloudstack_cloudfw_vpc_control_policy":                      resourceAlibabacloudStackCloudfwVpcControlPolicy(),
	}
	if v, err := stringToBool(os.Getenv("APSARASTACK_IN_ALIBABACLOUDSTACK")); err == nil && !v {
		return maps
	}
	new_map := map[string]*schema.Resource{}
	for key, value := range maps {
		new_map[key] = value
		if strings.HasPrefix(key, "alibabacloudstack_") {
			new_key := strings.Replace(key, "alibabacloudstack_", "apsarastack_", 1)
			new_map[new_key] = value
		}
	}
	return new_map
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var getProviderConfig = func(str string, key string) string {
		if str == "" {
			value, err := getConfigFromProfile(d, key)
			if err == nil && value != nil {
				str = value.(string)
			}
		}
		return str
	}

	accessKey := getProviderConfig(d.Get("access_key").(string), "access_key_id")
	secretKey := getProviderConfig(d.Get("secret_key").(string), "access_key_secret")
	region := getProviderConfig(d.Get("region").(string), "region_id")
	region = strings.TrimSpace(region)

	ecsRoleName := getProviderConfig(d.Get("ecs_role_name").(string), "ram_role_name")

	var eagleeye connectivity.EagleEye
	if os.Getenv("TF_EAGLEEYE_TRACEID") != "" && os.Getenv("TF_EAGLEEYE_RPCID") != "" {
		eagleeye = connectivity.EagleEye{
			TraceId: os.Getenv("TF_EAGLEEYE_TRACEID"),
			RpcId:   os.Getenv("TF_EAGLEEYE_RPCID"),
		}
	} else {
		eagleeye = connectivity.EagleEye{
			TraceId: connectivity.GenerateTraceId(),
			RpcId:   connectivity.DefaultRpcId,
		}
	}

	log.Printf("Eagleeye's trace id is: %s", eagleeye.GetTraceId())

	config := &connectivity.Config{
		AccessKey:            strings.TrimSpace(accessKey),
		SecretKey:            strings.TrimSpace(secretKey),
		EcsRoleName:          strings.TrimSpace(ecsRoleName),
		Region:               connectivity.Region(strings.TrimSpace(region)),
		RegionId:             strings.TrimSpace(region),
		ConfigurationSource:  d.Get("configuration_source").(string),
		Protocol:             d.Get("protocol").(string),
		ClientReadTimeout:    d.Get("client_read_timeout").(int),
		ClientConnectTimeout: d.Get("client_connect_timeout").(int),
		Insecure:             d.Get("insecure").(bool),
		Proxy:                d.Get("proxy").(string),
		Department:           d.Get("department").(string),
		ResourceGroup:        d.Get("resource_group").(string),
		ResourceSetName:      d.Get("resource_group_set_name").(string),
		SourceIp:             strings.TrimSpace(d.Get("source_ip").(string)),
		SecureTransport:      strings.TrimSpace(d.Get("secure_transport").(string)),
		Endpoints:            make(map[connectivity.ServiceCode]string),
		Eagleeye:             eagleeye,
		MaxRetryTimeout:      d.Get("max_retry_timeout").(int),
	}
	if v, ok := d.GetOk("security_transport"); config.SecureTransport == "" && ok && v.(string) != "" {
		config.SecureTransport = v.(string)
	}
	token := getProviderConfig(d.Get("security_token").(string), "sts_token")
	config.SecurityToken = strings.TrimSpace(token)
	config.RamRoleArn = getProviderConfig(d.Get("role_arn").(string), "ram_role_arn")
	log.Printf("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$led!!! %s", config.RamRoleArn)
	config.RamRoleSessionName = getProviderConfig("", "ram_session_name")
	if config.RamRoleSessionName == "" {
		config.RamRoleSessionName = "terraform"
	}
	expiredSeconds, err := getConfigFromProfile(d, "expired_seconds")
	if err == nil && expiredSeconds != nil {
		config.RamRoleSessionExpiration = (int)(expiredSeconds.(float64))
	}

	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	log.Printf("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$led!!! %s", assumeRoleList)
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		if assumeRole["role_arn"].(string) != "" {
			config.RamRoleArn = assumeRole["role_arn"].(string)
		}
		log.Printf("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$led!!! %s", config.RamRoleArn)
		if assumeRole["session_name"].(string) != "" {
			config.RamRoleSessionName = assumeRole["session_name"].(string)
		}
		config.RamRolePolicy = assumeRole["policy"].(string)
		if assumeRole["session_expiration"].(int) == 0 {
			if v := os.Getenv("ALIBABACLOUDSTACK_ASSUME_ROLE_SESSION_EXPIRATION"); v != "" {
				if expiredSeconds, err := strconv.Atoi(v); err == nil {
					config.RamRoleSessionExpiration = expiredSeconds
				}
			}
			if config.RamRoleSessionExpiration == 0 {
				config.RamRoleSessionExpiration = 3600
			}
		} else {
			config.RamRoleSessionExpiration = assumeRole["session_expiration"].(int)
		}

		log.Printf("[INFO] assume_role configuration set: (RamRoleArn: %q, RamRoleSessionName: %q, RamRolePolicy: %q, RamRoleSessionExpiration: %d)",
			config.RamRoleArn, config.RamRoleSessionName, config.RamRolePolicy, config.RamRoleSessionExpiration)
	}

	if err := config.MakeConfigByEcsRoleName(); err != nil {
		log.Printf("[ERROR] Assume role failed: %s", err)
		return nil, err
	}

	// domain := d.Get("domain").(string)
	// if domain != "" {
	// 	if strings.Contains(domain, "/") && d.Get("proxy").(string) != "" {
	// 		return nil, fmt.Errorf("[Error]Domain containing the character '/' is not supported for proxy configuration.")
	// 	}
	// 	// For services without generated popgw addresses, continue using asapi
	// 	var setEndpointIfEmpty = func(endpoint string, domain string) string {
	// 		if endpoint == "" {
	// 			return domain
	// 		}
	// 		return endpoint
	// 	}
	// 	for popcode := range connectivity.PopEndpoints {
	// 		if popcode == connectivity.OssDataCode {
	// 			// Oss data gateway is not configured
	// 			continue
	// 		}
	// 		if popcode == connectivity.SlSDataCode {
	// 			// SLS data gateway is not configured
	// 			continue
	// 		}
	// 		config.Endpoints[popcode] = setEndpointIfEmpty(config.Endpoints[popcode], domain)
	// 	}
	// }
	if v, ok := d.GetOk("popgw_domain"); ok && v.(string) != "" {
		popgw_domain := v.(string)
		log.Printf("Generator Popgw Endpoint: %s", popgw_domain)
		// Generate popgw addresses using the endpoint rules of each cloud product
		is_center_region := d.Get("is_center_region").(bool)
		for popcode := range connectivity.PopEndpoints {
			endpoint := connectivity.GeneratorEndpoint(popcode, region, popgw_domain, is_center_region)
			if endpoint != "" {
				config.Endpoints[popcode] = endpoint
			}
		}
	}
	if endpoints, ok := d.GetOk("endpoints"); ok {

		endpointsSet := endpoints.(*schema.Set)

		for _, endpointsSetI := range endpointsSet.List() {
			endpoints := endpointsSetI.(map[string]interface{})
			for popcode := range connectivity.PopEndpoints {
				endpoint := strings.TrimSpace(endpoints[strings.ToLower(string(popcode))].(string))
				if endpoint != "" {
					config.Endpoints[popcode] = endpoint
				}
			}
		}
	}
	StsEndpoint := d.Get("sts_endpoint").(string)
	if StsEndpoint != "" {
		config.Endpoints[connectivity.STSCode] = StsEndpoint
	}
	organizationAccessKey := d.Get("organization_accesskey").(string)
	organizationSecretKey := d.Get("organization_secretkey").(string)
	if organizationAccessKey != "" && organizationSecretKey != "" {
		config.AccessKey = organizationAccessKey
		config.SecretKey = organizationSecretKey
	}
	slsOpenAPIEndpoint := d.Get("sls_openapi_endpoint").(string)
	if slsOpenAPIEndpoint != "" {
		config.Endpoints[connectivity.SlSDataCode] = slsOpenAPIEndpoint
	}
	if kmsEndpoint, ok := d.GetOk("kms_endpoint"); ok && kmsEndpoint.(string) != "" {
		config.Endpoints[connectivity.KmsCode] = kmsEndpoint.(string)
	}

	if asapiEndpoint, ok := d.GetOk("asapi_endpoint"); ok && asapiEndpoint.(string) != "" {
		config.Endpoints[connectivity.ASAPICode] = asapiEndpoint.(string)
		config.Endpoints[connectivity.OneRouterCode] = asapiEndpoint.(string)
	}

	slsEndpoint := d.Get("sls_endpoint").(string)
	if slsEndpoint != "" {
		config.Endpoints[connectivity.SLSCode] = slsEndpoint
	}

	if strings.ToLower(config.Protocol) == "https" {
		config.Protocol = "HTTPS"
	} else {
		config.Protocol = "HTTP"
	}
	if config.RamRoleArn != "" {
		config.AccessKey, config.SecretKey, config.SecurityToken, err = getAssumeRoleAK(config)
		if err != nil {
			return nil, err
		}
	}
	config.ResourceSetName = d.Get("resource_group_set_name").(string)
	if config.Department == "" || config.ResourceGroup == "" {
		dept, resId, rgid, err := getResourceCredentials(config)
		if err != nil {
			return nil, err
		}
		config.Department = dept
		config.ResourceGroup = fmt.Sprintf("%d", rgid)
		config.ResourceGroupId = resId
	}

	if ots_instance_name, ok := d.GetOk("ots_instance_name"); ok && ots_instance_name.(string) != "" {
		config.OtsInstanceName = strings.TrimSpace(ots_instance_name.(string))
	}

	if account, ok := d.GetOk("account_id"); ok && account.(string) != "" {
		config.AccountId = strings.TrimSpace(account.(string))
	}

	if config.ConfigurationSource == "" {
		sourceName := fmt.Sprintf("Default/%s:%s", config.AccessKey, strings.Trim(uuid.New().String(), "-"))
		if len(sourceName) > 64 {
			sourceName = sourceName[:64]
		}
		config.ConfigurationSource = sourceName
	}
	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "The access key for API operations. You can retrieve this from the 'Security Management' section of the AlibabacloudStack console.",

		"secret_key": "The secret key for API operations. You can retrieve this from the 'Security Management' section of the AlibabacloudStackconsole.",

		"security_token": "security token. A security token is only required if you are using Security Token Service.",

		"insecure": "Use this to Trust self-signed certificates. It's typically used to allow insecure connections",

		"proxy": "Use this to set proxy connection",

		"popgw_domain": "Use this to override the default domain. It's typically used to connect to custom domain.",
	}
}
func endpointsSchema() *schema.Schema {
	schemas := make(map[string]*schema.Schema)
	for popcode := range connectivity.PopEndpoints {
		popcodeStr := strings.ToLower(string(popcode))
		schemas[popcodeStr] = &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: descriptions[popcodeStr+"_endpoint"],
		}
	}
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: schemas,
		},
		Set: func(v interface{}) int {
			var buf bytes.Buffer
			m := v.(map[string]interface{})
			for popcode := range connectivity.PopEndpoints {
				popcodeStr := strings.ToLower(string(popcode))
				buf.WriteString(fmt.Sprintf("%s-", m[popcodeStr].(string)))
			}
			return hashcode.String(buf.String())
		},
	}
}

func getConfigFromProfile(d *schema.ResourceData, ProfileKey string) (interface{}, error) {

	if providerConfig == nil {
		if v, ok := d.GetOk("profile"); !ok && v.(string) == "" {
			return nil, nil
		}
		current := d.Get("profile").(string)
		// Set CredsFilename, expanding home directory
		profilePath, err := homedir.Expand(d.Get("shared_credentials_file").(string))
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}
		if profilePath == "" {
			profilePath = fmt.Sprintf("%s/.alibabacloudstack/config.json", os.Getenv("HOME"))
			if runtime.GOOS == "windows" {
				profilePath = fmt.Sprintf("%s/.alibabacloudstack/config.json", os.Getenv("USERPROFILE"))
			}
		}
		providerConfig = make(map[string]interface{})
		_, err = os.Stat(profilePath)
		if !os.IsNotExist(err) {
			data, err := ioutil.ReadFile(profilePath)
			if err != nil {
				return nil, errmsgs.WrapError(err)
			}
			config := map[string]interface{}{}
			err = json.Unmarshal(data, &config)
			if err != nil {
				return nil, errmsgs.WrapError(err)
			}
			for _, v := range config["profiles"].([]interface{}) {
				if current == v.(map[string]interface{})["name"] {
					providerConfig = v.(map[string]interface{})
				}
			}
		}
	}

	mode := ""
	if v, ok := providerConfig["mode"]; ok {
		mode = v.(string)
	} else {
		return v, nil
	}
	switch ProfileKey {
	case "access_key_id", "access_key_secret":
		if mode == "EcsRamRole" {
			return "", nil
		}
	case "ram_role_name":
		if mode != "EcsRamRole" {
			return "", nil
		}
	case "sts_token":
		if mode != "StsToken" {
			return "", nil
		}
	case "ram_role_arn", "ram_session_name":
		if mode != "RamRoleArn" {
			return "", nil
		}
	case "expired_seconds":
		if mode != "RamRoleArn" {
			return float64(0), nil
		}
	}

	return providerConfig[ProfileKey], nil
}
func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Required:    true,
					Description: descriptions["assume_role_role_arn"],
					DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ASSUME_ROLE_ARN", os.Getenv("ALIBABACLOUDSTACK_ASSUME_ROLE_ARN")),
				},
				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_session_name"],
					DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ASSUME_ROLE_SESSION_NAME", ""),
				},
				"policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_policy"],
				},
				"session_expiration": {
					Type:         schema.TypeInt,
					Optional:     true,
					Description:  descriptions["assume_role_session_expiration"],
					ValidateFunc: intBetween(900, 3600),
				},
			},
		},
	}
}

func getAssumeRoleAK(config *connectivity.Config) (string, string, string, error) {
	client, err := config.Client()
	if err != nil {
		return "", "", "", err
	}
	request := sts.CreateAssumeRoleRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.Scheme = "https" // sts must be connected via https
	request.RoleArn = config.RamRoleArn
	request.RoleSessionName = config.RamRoleSessionName
	//request.DurationSeconds = requests.NewInteger(config.RamRoleSessionExpiration)
	request.Policy = config.RamRolePolicy

	conn, err := client.WithProductSDKClient(connectivity.STSCode)
	if err != nil {
		return "", "", "", err
	}
	stsClient := &sts.Client{
		Client: *conn,
	}
	response, err := stsClient.AssumeRole(request)
	addDebug(request.GetActionName(), response, request.RpcRequest, request)
	if err != nil {
		return config.AccessKey, config.SecretKey, config.SecurityToken, err
	}

	return response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, response.Credentials.SecurityToken, nil
}

func getResourceCredentials(config *connectivity.Config) (string, string, int, error) {
	endpoint := config.Endpoints[connectivity.ASCMCode]
	var client *sts.Client
	var err error
	if config.SecurityToken == "" {
		client, err = sts.NewClientWithAccessKey(config.RegionId, config.AccessKey, config.SecretKey)
	} else {
		client, err = sts.NewClientWithStsToken(config.RegionId, config.AccessKey, config.SecretKey, config.SecurityToken)
	}

	request := requests.NewCommonRequest()
	if config.Insecure {
		request.SetHTTPSInsecure(config.Insecure)
	}
	client.Domain = endpoint
	if config.Proxy != "" {
		client.SetHttpProxy(config.Proxy)
		client.SetHttpsProxy(config.Proxy)
	}
	request.RegionId = config.RegionId
	if strings.ToLower(config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Version = "2019-05-10"
	request.Method = "POST"
	request.Product = "ascm"
	request.ApiName = "ListResourceGroup"
	if !strings.HasPrefix(client.Domain, "internal.asapi.") && !strings.HasPrefix(client.Domain, "public.asapi.") {
		request.PathPattern = "/ascm/auth/resource_group/list_resource_group"
	}

	request.QueryParams = map[string]string{
		"resourceGroupName": config.ResourceSetName,
		"pageNumber":        "1",
		"pageSize":          "10",
	}
	if config.SecurityToken != "" {
		request.QueryParams["SecurityToken"] = config.SecurityToken
	}
	request.Headers["Content-Type"] = "application/json"
	request.Headers["x-ascm-product-name"] = "ascm"
	resp, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), resp, request, request.QueryParams)
	if err != nil {
		return "", "", 0, err
	}
	response := &ResourceGroup{}
	err = json.Unmarshal(resp.GetHttpContentBytes(), response)
	if err != nil {
		return "", "", 0, err
	}

	var deptId int   // Organization ID
	var resGrpId int //ID of resource set
	var resGrp string
	deptId = 0
	matched := 0
	if response.Code != "200" {
		return "", "", 0, fmt.Errorf("unable to initialize the ascm client: department or resource_group is not provided")
	}
	for _, d := range response.Data {
		if d.ResourceGroupName == config.ResourceSetName {
			matched += 1
		}
	}
	if matched == 0 {
		return "", "", 0, fmt.Errorf("resource group ID and organization not found for resource set %s", config.ResourceSetName)
	} else if matched > 1 {
		return "", "", 0, errmsgs.Error("There exists a resource group set name with the same name, Please Provider department or resource_group")
	} else {
		for _, j := range response.Data {
			if j.ResourceGroupName == config.ResourceSetName {
				deptId = j.OrganizationID
				resGrp = j.RsID
				resGrpId = j.ID
				break
			}
		}
	}

	//log.Printf("[INFO] Get Resource Group Details Succssfull for Resource set: %s : Department: %s, ResourceGroupId: %s", config.ResourceSetName, fmt.Sprint(response.Data[0].OrganizationID), fmt.Sprint(response.Data[0].ID))
	log.Printf("[INFO] Get Resource Group Details Succssfull for Resource set: %s : Department: %d, ResourceGroup: %s, ResourceGroupId: %d", config.ResourceSetName, deptId, resGrp, resGrpId)
	//return fmt.Sprint(response.Data[0].OrganizationID), fmt.Sprint(response.Data[0].ID), err
	return fmt.Sprint(deptId), resGrp, resGrpId, err

}

func waitSecondsIfWithTest(second int) {
	// Sleep for one second in test mode to prevent data caching from causing secondary plan failures
	if os.Getenv("TF_ACC") == "1" {
		time.Sleep(time.Duration(second) * time.Second)
	}
}
