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

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ACCESS_KEY", os.Getenv("ALIBABACLOUDSTACK_ACCESS_KEY")),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SECRET_KEY", os.Getenv("ALIBABACLOUDSTACK_SECRET_KEY")),
				Description: descriptions["secret_key"],
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_REGION", os.Getenv("ALIBABACLOUDSTACK_REGION")),
				Description: descriptions["region"],
			},
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_REGION", os.Getenv("ALIBABACLOUDSTACK_REGION")),
				Description: descriptions["region_id"],
				Deprecated:  "Use parameter region replace it.",
			},
			"role_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["assume_role_role_arn"],
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ASSUME_ROLE_ARN", os.Getenv("ALIBABACLOUDSTACK_ASSUME_ROLE_ARN")),
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SECURITY_TOKEN", os.Getenv("SECURITY_TOKEN")),
				Description: descriptions["security_token"],
			},
			"ecs_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ECS_ROLE_NAME", os.Getenv("ALIBABACLOUDSTACK_ECS_ROLE_NAME")),
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
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_PROFILE", ""),
			},
			"endpoints": endpointsSchema(),
			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_credentials_file"],
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SHARED_CREDENTIALS_FILE", ""),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_INSECURE", false),
				Description: descriptions["insecure"],
			},
			"assume_role": assumeRoleSchema(),
			"fc": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'fc' has been deprecated from provider version 1.28.0. New field 'fc' which in nested endpoints instead.",
			},
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
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SOURCE_IP", os.Getenv("ALIBABACLOUDSTACK_SOURCE_IP")),
				Description: descriptions["source_ip"],
			},
			"security_transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SECURITY_TRANSPORT", os.Getenv("ALIBABACLOUDSTACK_SECURITY_TRANSPORT")),
				//Deprecated:  "It has been deprecated from version 1.136.0 and using new field secure_transport instead.",
			},
			"secure_transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SECURE_TRANSPORT", os.Getenv("ALIBABACLOUDSTACK_SECURE_TRANSPORT")),
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
			"force_use_asapi": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["force_use_asapi"],
			},
			"is_center_region": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: descriptions["is_center_region"],
			},
			"popgw_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_POPGW_DOMAIN", nil),
				Description: descriptions["popgw_domain"],
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_DOMAIN", nil),
				Description: descriptions["domain"],
				Deprecated:  "ASAPI will no longer provide external services by default on apsarastack v3.18.1",
			},
			"ossservice_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_OSSSERVICE_DOMAIN", nil),
				Description: descriptions["ossservice_domain"],
				Deprecated:  "Use schema endpoints replace ossservice_domain.",
			},
			"kafkaopenapi_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_KAFKAOPENAPI_DOMAIN", nil),
				Description: descriptions["kafkaopenapi_domain"],
			},
			"organization_accesskey": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ORGANIZATION_ACCESSKEY", os.Getenv("ALIBABACLOUDSTACK_ORGANIZATION_ACCESSKEY")),
				Description: descriptions["organization_accesskey"],
			},
			"organization_secretkey": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ORGANIZATION_SECRETKEY", os.Getenv("ALIBABACLOUDSTACK_ORGANIZATION_SECRETKEY")),
				Description: descriptions["organization_secretkey"],
			},
			"sls_openapi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SLS_OPENAPI_ENDPOINT", nil),
				Description: descriptions["sls_openapi_endpoint"],
				Deprecated:  "Use schema endpoints replace sls_openapi_endpoint.",
			},
			"ascm_openapi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ASCM_OPENAPI_ENDPOINT", nil),
				Description: descriptions["ascm_openapi_endpoint"],
				Deprecated:  "Use schema endpoints replace ascm_openapi_endpoint.",
			},
			"sts_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_STS_ENDPOINT", nil),
				Description: descriptions["sts_endpoint"],
				Deprecated:  "Use schema endpoints replace sts_endpoint.",
			},
			"quickbi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_QUICKBI_ENDPOINT", nil),
				Description: descriptions["quickbi_endpoint"],
				Deprecated:  "Use schema endpoints replace quickbi_endpoint.",
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
			"dataworkspublic": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_DATAWORKS_PUBLIC_ENDPOINT", nil),
				Description: descriptions["dataworkspublic_endpoint"],
				Deprecated:  "Use schema endpoints replace dataworkspublic.",
			},
			"dbs_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_DBS_ENDPOINT", nil),
				Description: descriptions["dbs_endpoint"],
				Deprecated:  "Use schema endpoints replace dbs_endpoint.",
			},
		},
		DataSourcesMap: getDataSourcesMap(),
		ResourcesMap:   getResourcesMap(),
		ConfigureFunc:  providerConfigure,
	}
}

var providerConfig map[string]interface{}

func stringToBool(value string) (bool, error) {
	// 将字符串转换为小写以便于比较
	value = strings.ToLower(value)

	// 检查常见的布尔值表示
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
			"alibabacloudstack_account":                                dataSourceAlibabacloudStackAccount(),
			"alibabacloudstack_adb_clusters":                           dataSourceAlibabacloudStackAdbDbClusters(),
			"alibabacloudstack_adb_zones":                              dataSourceAlibabacloudStackAdbZones(),
			"alibabacloudstack_adb_db_clusters":                        dataSourceAlibabacloudStackAdbDbClusters(),
			"alibabacloudstack_adb_dbclusters":                         dataSourceAlibabacloudStackAdbDbClusters(),
			"alibabacloudstack_api_gateway_apis":                       dataSourceAlibabacloudStackApiGatewayApis(),
			"alibabacloudstack_apigateway_apis":                        dataSourceAlibabacloudStackApiGatewayApis(),
			"alibabacloudstack_api_gateway_apps":                       dataSourceAlibabacloudStackApiGatewayApps(),
			"alibabacloudstack_api_gateway_groups":                     dataSourceAlibabacloudStackApiGatewayGroups(),
			"alibabacloudstack_apigateway_apigroups":                   dataSourceAlibabacloudStackApiGatewayGroups(),
			"alibabacloudstack_api_gateway_service":                    dataSourceAlibabacloudStackApiGatewayService(),
			"alibabacloudstack_ascm_resource_groups":                   dataSourceAlibabacloudStackAscmResourceGroups(),
			"alibabacloudstack_ascm_users":                             dataSourceAlibabacloudStackAscmUsers(),
			"alibabacloudstack_ascm_user_groups":                       dataSourceAlibabacloudStackAscmUserGroups(),
			"alibabacloudstack_ascm_logon_policies":                    dataSourceAlibabacloudStackAscmLogonPolicies(),
			"alibabacloudstack_ascm_ram_service_roles":                 dataSourceAlibabacloudStackAscmRamServiceRoles(),
			"alibabacloudstack_ascm_organizations":                     dataSourceAlibabacloudStackAscmOrganizations(),
			"alibabacloudstack_ascm_instance_families":                 dataSourceAlibabacloudStackInstanceFamilies(),
			"alibabacloudstack_ascm_regions_by_product":                dataSourceAlibabacloudStackRegionsByProduct(),
			"alibabacloudstack_ascm_service_cluster_by_product":        dataSourceAlibabacloudStackServiceClusterByProduct(),
			"alibabacloudstack_ascm_ecs_instance_families":             dataSourceAlibabacloudStackEcsInstanceFamilies(),
			"alibabacloudstack_ascm_specific_fields":                   dataSourceAlibabacloudStackSpecificFields(),
			"alibabacloudstack_ascm_environment_services_by_product":   dataSourceAlibabacloudStackAscmEnvironmentServicesByProduct(),
			"alibabacloudstack_ascm_password_policies":                 dataSourceAlibabacloudStackAscmPasswordPolicies(),
			"alibabacloudstack_ascm_quotas":                            dataSourceAlibabacloudStackQuotas(),
			"alibabacloudstack_ascm_metering_query_ecs":                dataSourceAlibabacloudstackAscmMeteringQueryEcs(),
			"alibabacloudstack_ascm_roles":                             dataSourceAlibabacloudStackAscmRoles(),
			"alibabacloudstack_ascm_ram_policies":                      dataSourceAlibabacloudStackAscmRamPolicies(),
			"alibabacloudstack_ascm_ram_policies_for_user":             dataSourceAlibabacloudStackAscmRamPoliciesForUser(),
			"alibabacloudstack_common_bandwidth_packages":              dataSourceAlibabacloudStackCommonBandwidthPackages(),
			"alibabacloudstack_cbwp_commonbandwidthpackages":           dataSourceAlibabacloudStackCommonBandwidthPackages(),
			"alibabacloudstack_cr_ee_instances":                        dataSourceAlibabacloudStackCrEEInstances(),
			"alibabacloudstack_cr_ee_namespaces":                       dataSourceAlibabacloudStackCrEENamespaces(),
			"alibabacloudstack_cr_ee_repos":                            dataSourceAlibabacloudStackCrEERepos(),
			"alibabacloudstack_cr_repositories":                        dataSourceAlibabacloudStackCrEERepos(),
			"alibabacloudstack_cr_ee_sync_rules":                       dataSourceAlibabacloudStackCrEESyncRules(),
			"alibabacloudstack_cr_namespaces":                          dataSourceAlibabacloudStackCRNamespaces(),
			"alibabacloudstack_cr_repos":                               dataSourceAlibabacloudStackCRRepos(),
			"alibabacloudstack_cs_kubernetes_clusters":                 dataSourceAlibabacloudStackCSKubernetesClusters(),
			"alibabacloudstack_cs_kubernetes_clusters_kubeconfig":      dataSourceAlibabacloudStackCSKubernetesClustersKubeConfig(),
			"alibabacloudstack_cms_metric_rule_templates":              dataSourceAlibabacloudStackCmsMetricRuleTemplates(),
			"alibabacloudstack_cms_alarm_contacts":                     dataSourceAlibabacloudstackCmsAlarmContacts(),
			"alibabacloudstack_cloudmonitorservice_alarmcontacts":      dataSourceAlibabacloudstackCmsAlarmContacts(),
			"alibabacloudstack_cms_alarm_contact_groups":               dataSourceAlibabacloudstackCmsAlarmContactGroups(),
			"alibabacloudstack_cloudmonitorservice_alarmcontactgroups": dataSourceAlibabacloudstackCmsAlarmContactGroups(),
			"alibabacloudstack_cms_project_meta":                       dataSourceAlibabacloudstackCmsProjectMeta(),
			"alibabacloudstack_cms_metric_metalist":                    dataSourceAlibabacloudstackCmsMetricMetalist(),
			"alibabacloudstack_cms_alarms":                             dataSourceAlibabacloudstackCmsAlarms(),
			"alibabacloudstack_cloudmonitorservice_metricalarmrules":   dataSourceAlibabacloudstackCmsAlarms(),
			"alibabacloudstack_datahub_service":                        dataSourceAlibabacloudStackDatahubService(),
			"alibabacloudstack_db_instances":                           dataSourceAlibabacloudStackDBInstances(),
			"alibabacloudstack_db_zones":                               dataSourceAlibabacloudStackDBZones(),
			"alibabacloudstack_disks":                                  dataSourceAlibabacloudStackDisks(),
			"alibabacloudstack_ecs_disks":                              dataSourceAlibabacloudStackDisks(),
			"alibabacloudstack_dns_records":                            dataSourceAlibabacloudStackDnsRecords(),
			"alibabacloudstack_dns_groups":                             dataSourceAlibabacloudStackDnsGroups(),
			"alibabacloudstack_alidns_domaingroups":                    dataSourceAlibabacloudStackDnsGroups(),
			"alibabacloudstack_dns_domains":                            dataSourceAlibabacloudStackDnsDomains(),
			"alibabacloudstack_drds_instances":                         dataSourceAlibabacloudStackDRDSInstances(),
			"alibabacloudstack_dms_enterprise_instances":               dataSourceAlibabacloudStackDmsEnterpriseInstances(),
			"alibabacloudstack_dmsenterprise_instances":                dataSourceAlibabacloudStackDmsEnterpriseInstances(),
			"alibabacloudstack_dms_enterprise_users":                   dataSourceAlibabacloudStackDmsEnterpriseUsers(),
			"alibabacloudstack_dmsenterprise_users":                    dataSourceAlibabacloudStackDmsEnterpriseUsers(),
			"alibabacloudstack_ecs_commands":                           dataSourceAlibabacloudStackEcsCommands(),
			"alibabacloudstack_ecs_deployment_sets":                    dataSourceAlibabacloudStackEcsDeploymentSets(),
			"alibabacloudstack_ecs_deploymentsets":                     dataSourceAlibabacloudStackEcsDeploymentSets(),
			"alibabacloudstack_ecs_hpc_clusters":                       dataSourceAlibabacloudStackEcsHpcClusters(),
			"alibabacloudstack_ecs_hpcclusters":                        dataSourceAlibabacloudStackEcsHpcClusters(),
			"alibabacloudstack_ecs_dedicated_hosts":                    dataSourceAlibabacloudStackEcsDedicatedHosts(),
			"alibabacloudstack_ecs_dedicatedhosts":                     dataSourceAlibabacloudStackEcsDedicatedHosts(),
			"alibabacloudstack_edas_deploy_groups":                     dataSourceAlibabacloudStackEdasDeployGroups(),
			"alibabacloudstack_edas_deploygroups":                      dataSourceAlibabacloudStackEdasDeployGroups(),
			"alibabacloudstack_edas_clusters":                          dataSourceAlibabacloudStackEdasClusters(),
			"alibabacloudstack_edas_applications":                      dataSourceAlibabacloudStackEdasApplications(),
			"alibabacloudstack_edas_slbattachments":                    dataSourceAlibabacloudStackEdasApplications(),
			"alibabacloudstack_eips":                                   dataSourceAlibabacloudStackEips(),
			"alibabacloudstack_eip_addresses":                          dataSourceAlibabacloudStackEips(),
			"alibabacloudstack_ess_scaling_configurations":             dataSourceAlibabacloudStackEssScalingConfigurations(),
			"alibabacloudstack_ess_scaling_groups":                     dataSourceAlibabacloudStackEssScalingGroups(),
			"alibabacloudstack_ess_lifecycle_hooks":                    dataSourceAlibabacloudStackEssLifecycleHooks(),
			"alibabacloudstack_ess_notifications":                      dataSourceAlibabacloudStackEssNotifications(),
			"alibabacloudstack_autoscaling_notifications":              dataSourceAlibabacloudStackEssNotifications(),
			"alibabacloudstack_ess_scaling_rules":                      dataSourceAlibabacloudStackEssScalingRules(),
			"alibabacloudstack_ess_scheduled_tasks":                    dataSourceAlibabacloudStackEssScheduledTasks(),
			"alibabacloudstack_autoscaling_scheduledtasks":             dataSourceAlibabacloudStackEssScheduledTasks(),
			"alibabacloudstack_forward_entries":                        dataSourceAlibabacloudStackForwardEntries(),
			"alibabacloudstack_natgateway_forwardentries":              dataSourceAlibabacloudStackForwardEntries(),
			"alibabacloudstack_gpdb_accounts":                          dataSourceAlibabacloudStackGpdbAccounts(),
			"alibabacloudstack_gpdb_instances":                         dataSourceAlibabacloudStackGpdbInstances(),
			"alibabacloudstack_hbase_instances":                        dataSourceAlibabacloudStackHBaseInstances(),
			"alibabacloudstack_hbase_clusters":                         dataSourceAlibabacloudStackHBaseInstances(),
			"alibabacloudstack_instances":                              dataSourceAlibabacloudStackInstances(),
			"alibabacloudstack_ecs_instances":                          dataSourceAlibabacloudStackInstances(),
			"alibabacloudstack_instance_type_families":                 dataSourceAlibabacloudStackInstanceTypeFamilies(),
			"alibabacloudstack_instance_types":                         dataSourceAlibabacloudStackInstanceTypes(),
			"alibabacloudstack_images":                                 dataSourceAlibabacloudStackImages(),
			"alibabacloudstack_ecs_images":                             dataSourceAlibabacloudStackImages(),
			"alibabacloudstack_key_pairs":                              dataSourceAlibabacloudStackKeyPairs(),
			"alibabacloudstack_ecs_keypairs":                           dataSourceAlibabacloudStackKeyPairs(),
			"alibabacloudstack_kms_aliases":                            dataSourceAlibabacloudStackKmsAliases(),
			"alibabacloudstack_kms_ciphertext":                         dataSourceAlibabacloudStackKmsCiphertext(),
			"alibabacloudstack_kms_keys":                               dataSourceAlibabacloudStackKmsKeys(),
			"alibabacloudstack_kms_secrets":                            dataSourceAlibabacloudStackKmsSecrets(),
			"alibabacloudstack_kvstore_instances":                      dataSourceAlibabacloudStackKVStoreInstances(),
			"alibabacloudstack_redis_tairinstances":                    dataSourceAlibabacloudStackKVStoreInstances(),
			"alibabacloudstack_kvstore_zones":                          dataSourceAlibabacloudStackKVStoreZones(),
			"alibabacloudstack_kvstore_instance_classes":               dataSourceAlibabacloudStackKVStoreInstanceClasses(),
			"alibabacloudstack_kvstore_instance_engines":               dataSourceAlibabacloudStackKVStoreInstanceEngines(),
			"alibabacloudstack_mongodb_instances":                      dataSourceAlibabacloudStackMongoDBInstances(),
			"alibabacloudstack_mongodb_zones":                          dataSourceAlibabacloudStackMongoDBZones(),
			"alibabacloudstack_maxcompute_cus":                         dataSourceAlibabacloudStackMaxcomputeCus(),
			"alibabacloudstack_maxcompute_users":                       dataSourceAlibabacloudStackMaxcomputeUsers(),
			"alibabacloudstack_maxcompute_clusters":                    dataSourceAlibabacloudStackMaxcomputeClusters(),
			"alibabacloudstack_maxcompute_cluster_qutaos":              dataSourceAlibabacloudStackMaxcomputeClusterQutaos(),
			"alibabacloudstack_maxcompute_projects":                    dataSourceAlibabacloudStackMaxcomputeProjects(),
			"alibabacloudstack_nas_zones":                              dataSourceAlibabacloudStackNasZones(),
			"alibabacloudstack_nas_protocols":                          dataSourceAlibabacloudStackNasProtocols(),
			"alibabacloudstack_nas_file_systems":                       dataSourceAlibabacloudStackFileSystems(),
			"alibabacloudstack_nas_filesystems":                        dataSourceAlibabacloudStackFileSystems(),
			"alibabacloudstack_nas_mount_targets":                      dataSourceAlibabacloudStackNasMountTargets(),
			"alibabacloudstack_nas_mounttargets":                       dataSourceAlibabacloudStackNasMountTargets(),
			"alibabacloudstack_nas_access_rules":                       dataSourceAlibabacloudStackAccessRules(),
			"alibabacloudstack_nas_accessrules":                        dataSourceAlibabacloudStackAccessRules(),
			"alibabacloudstack_nat_gateways":                           dataSourceAlibabacloudStackNatGateways(),
			"alibabacloudstack_natgateway_natgateways":                 dataSourceAlibabacloudStackNatGateways(),
			"alibabacloudstack_network_acls":                           dataSourceAlibabacloudStackNetworkAcls(),
			"alibabacloudstack_vpc_networkacls":                        dataSourceAlibabacloudStackNetworkAcls(),
			"alibabacloudstack_network_interfaces":                     dataSourceAlibabacloudStackNetworkInterfaces(),
			"alibabacloudstack_ecs_networkinterfaces":                  dataSourceAlibabacloudStackNetworkInterfaces(),
			"alibabacloudstack_oss_buckets":                            dataSourceAlibabacloudStackOssBuckets(),
			"alibabacloudstack_oss_bucket_objects":                     dataSourceAlibabacloudStackOssBucketObjects(),
			"alibabacloudstack_ons_instances":                          dataSourceAlibabacloudStackOnsInstances(),
			"alibabacloudstack_ons_topics":                             dataSourceAlibabacloudStackOnsTopics(),
			"alibabacloudstack_ons_groups":                             dataSourceAlibabacloudStackOnsGroups(),
			"alibabacloudstack_ots_tables":                             dataSourceAlibabacloudStackOtsTables(),
			"alibabacloudstack_ots_instances":                          dataSourceAlibabacloudStackOtsInstances(),
			"alibabacloudstack_ots_instance_attachments":               dataSourceAlibabacloudStackOtsInstanceAttachments(),
			"alibabacloudstack_ots_instance_attachment":                dataSourceAlibabacloudStackOtsInstanceAttachments(),
			"alibabacloudstack_ots_service":                            dataSourceAlibabacloudStackOtsService(),
			"alibabacloudstack_quick_bi_users":                         dataSourceAlibabacloudStackQuickBiUsers(),
			"alibabacloudstack_router_interfaces":                      dataSourceAlibabacloudStackRouterInterfaces(),
			"alibabacloudstack_expressconnect_routerinterfaces":        dataSourceAlibabacloudStackRouterInterfaces(),
			"alibabacloudstack_ram_service_role_products":              dataSourceAlibabacloudstackRamServiceRoleProducts(),
			"alibabacloudstack_route_tables":                           dataSourceAlibabacloudStackRouteTables(),
			"alibabacloudstack_vpc_routetables":                        dataSourceAlibabacloudStackRouteTables(),
			"alibabacloudstack_route_entries":                          dataSourceAlibabacloudStackRouteEntries(),
			"alibabacloudstack_ros_stacks":                             dataSourceAlibabacloudStackRosStacks(),
			"alibabacloudstack_ros_templates":                          dataSourceAlibabacloudStackRosTemplates(),
			"alibabacloudstack_security_groups":                        dataSourceAlibabacloudStackSecurityGroups(),
			"alibabacloudstack_ecs_securitygroups":                     dataSourceAlibabacloudStackSecurityGroups(),
			"alibabacloudstack_security_group_rules":                   dataSourceAlibabacloudStackSecurityGroupRules(),
			"alibabacloudstack_snapshots":                              dataSourceAlibabacloudStackSnapshots(),
			"alibabacloudstack_ecs_snapshots":                          dataSourceAlibabacloudStackSnapshots(),
			"alibabacloudstack_slb_listeners":                          dataSourceAlibabacloudStackSlbListeners(),
			"alibabacloudstack_slb_server_groups":                      dataSourceAlibabacloudStackSlbServerGroups(),
			"alibabacloudstack_slb_vservergroups":                      dataSourceAlibabacloudStackSlbServerGroups(),
			"alibabacloudstack_slb_acls":                               dataSourceAlibabacloudStackSlbAcls(),
			"alibabacloudstack_slb_accesscontrollists":                 dataSourceAlibabacloudStackSlbAcls(),
			"alibabacloudstack_slb_domain_extensions":                  dataSourceAlibabacloudStackSlbDomainExtensions(),
			"alibabacloudstack_slb_domainextensions":                   dataSourceAlibabacloudStackSlbDomainExtensions(),
			"alibabacloudstack_slb_rules":                              dataSourceAlibabacloudStackSlbRules(),
			"alibabacloudstack_slb_master_slave_server_groups":         dataSourceAlibabacloudStackSlbMasterSlaveServerGroups(),
			"alibabacloudstack_slb_masterslaveservergroups":            dataSourceAlibabacloudStackSlbMasterSlaveServerGroups(),
			"alibabacloudstack_slbs":                                   dataSourceAlibabacloudStackSlbs(),
			"alibabacloudstack_slb_loadbalancers":                      dataSourceAlibabacloudStackSlbs(),
			"alibabacloudstack_slb_zones":                              dataSourceAlibabacloudStackSlbZones(),
			"alibabacloudstack_snat_entries":                           dataSourceAlibabacloudStackSnatEntries(),
			"alibabacloudstack_natgateway_snatentries":                 dataSourceAlibabacloudStackSnatEntries(),
			"alibabacloudstack_slb_server_certificates":                dataSourceAlibabacloudStackSlbServerCertificates(),
			"alibabacloudstack_slb_servercertificates":                 dataSourceAlibabacloudStackSlbServerCertificates(),
			"alibabacloudstack_slb_ca_certificates":                    dataSourceAlibabacloudStackSlbCACertificates(),
			"alibabacloudstack_slb_cacertificates":                     dataSourceAlibabacloudStackSlbCACertificates(),
			"alibabacloudstack_slb_backend_servers":                    dataSourceAlibabacloudStackSlbBackendServers(),
			"alibabacloudstack_slb_backendservers":                     dataSourceAlibabacloudStackSlbBackendServers(),
			"alibabacloudstack_tsdb_zones":                             dataSourceAlibabacloudStackTsdbZones(),
			"alibabacloudstack_vpn_gateways":                           dataSourceAlibabacloudStackVpnGateways(),
			"alibabacloudstack_vpngateway_vpngateways":                 dataSourceAlibabacloudStackVpnGateways(),
			"alibabacloudstack_vpn_customer_gateways":                  dataSourceAlibabacloudStackVpnCustomerGateways(),
			"alibabacloudstack_vpngateway_customergateways":            dataSourceAlibabacloudStackVpnCustomerGateways(),
			"alibabacloudstack_vpn_connections":                        dataSourceAlibabacloudStackVpnConnections(),
			"alibabacloudstack_vpngateway_vpnconnections":              dataSourceAlibabacloudStackVpnConnections(),
			"alibabacloudstack_vpc_ipv6_gateways":                      dataSourceAlibabacloudStackVpcIpv6Gateways(),
			"alibabacloudstack_vpc_ipv6_egress_rules":                  dataSourceAlibabacloudStackVpcIpv6EgressRules(),
			"alibabacloudstack_vpc_ipv6_egressrules":                   dataSourceAlibabacloudStackVpcIpv6EgressRules(),
			"alibabacloudstack_vpc_ipv6_addresses":                     dataSourceAlibabacloudStackVpcIpv6Addresses(),
			"alibabacloudstack_vpc_ipv6_internet_bandwidths":           dataSourceAlibabacloudStackVpcIpv6InternetBandwidths(),
			"alibabacloudstack_vpc_ipv6_internetbandwidths":            dataSourceAlibabacloudStackVpcIpv6InternetBandwidths(),
			"alibabacloudstack_vswitches":                              dataSourceAlibabacloudStackVSwitches(),
			"alibabacloudstack_vpc_vswitches":                          dataSourceAlibabacloudStackVSwitches(),
			"alibabacloudstack_vpcs":                                   dataSourceAlibabacloudStackVpcs(),
			"alibabacloudstack_vpc_vpcs":                               dataSourceAlibabacloudStackVpcs(),
			"alibabacloudstack_zones":                                  dataSourceAlibabacloudStackZones(),
			"alibabacloudstack_elasticsearch_instances":                dataSourceAlibabacloudStackElasticsearch(),
			"alibabacloudstack_elasticsearch_zones":                    dataSourceAlibabacloudStackElaticsearchZones(),
			"alibabacloudstack_ehpc_job_templates":                     dataSourceAlibabacloudStackEhpcJobTemplates(),
			"alibabacloudstack_oos_executions":                         dataSourceAlibabacloudStackOosExecutions(),
			"alibabacloudstack_oos_templates":                          dataSourceAlibabacloudStackOosTemplates(),
			"alibabacloudstack_express_connect_physical_connections":   dataSourceAlibabacloudStackExpressConnectPhysicalConnections(),
			"alibabacloudstack_expressconnect_physicalconnections":     dataSourceAlibabacloudStackExpressConnectPhysicalConnections(),
			"alibabacloudstack_express_connect_access_points":          dataSourceAlibabacloudStackExpressConnectAccessPoints(),
			"alibabacloudstack_express_connect_virtual_border_routers": dataSourceAlibabacloudStackExpressConnectVirtualBorderRouters(),
			"alibabacloudstack_expressconnect_virtualborderrouters":    dataSourceAlibabacloudStackExpressConnectVirtualBorderRouters(),
			"alibabacloudStack_cloud_firewall_control_policies":        dataSourceAlibabacloudStackCloudFirewallControlPolicies(),
			"alibabacloudstack_ecs_ebs_storage_sets":                   dataSourceAlibabacloudStackEcsEbsStorageSets(),
		}
	if v, err := stringToBool(os.Getenv("APSARASTACK_IN_ALIBABACLOUDSTACK")); err != nil && !v {
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
			"alibabacloudstack_ess_scaling_configuration":             resourceAlibabacloudStackEssScalingConfiguration(),
			"alibabacloudstack_adb_account":                           resourceAlibabacloudStackAdbAccount(),
			"alibabacloudstack_adb_backup_policy":                     resourceAlibabacloudStackAdbBackupPolicy(),
			"alibabacloudstack_adb_backuppolicy":                      resourceAlibabacloudStackAdbBackupPolicy(),
			"alibabacloudstack_adb_cluster":                           resourceAlibabacloudStackAdbDbCluster(),
			"alibabacloudstack_adb_connection":                        resourceAlibabacloudStackAdbConnection(),
			"alibabacloudstack_adb_db_cluster":                        resourceAlibabacloudStackAdbDbCluster(),
			"alibabacloudstack_adb_dbcluster":                         resourceAlibabacloudStackAdbDbCluster(),
			"alibabacloudstack_alikafka_sasl_acl":                     resourceAlibabacloudStackAlikafkaSaslAcl(),
			"alibabacloudstack_alikafka_sasl_user":                    resourceAlibabacloudStackAlikafkaSaslUser(),
			"alibabacloudstack_alikafka_topic":                        resourceAlibabacloudStackAlikafkaTopic(),
			"alibabacloudstack_api_gateway_api":                       resourceAlibabacloudStackApigatewayApi(),
			"alibabacloudstack_apigateway_api":                        resourceAlibabacloudStackApigatewayApi(),
			"alibabacloudstack_api_gateway_app":                       resourceAlibabacloudStackApigatewayApp(),
			"alibabacloudstack_api_gateway_app_attachment":            resourceAliyunApigatewayAppAttachment(),
			"alibabacloudstack_apigateway_app":                        resourceAliyunApigatewayAppAttachment(),
			"alibabacloudstack_api_gateway_group":                     resourceAlibabacloudStackApigatewayGroup(),
			"alibabacloudstack_apigateway_apigroup":                   resourceAlibabacloudStackApigatewayGroup(),
			"alibabacloudstack_api_gateway_vpc_access":                resourceAlibabacloudStackApigatewayVpc(),
			"alibabacloudstack_apigateway_vpc":                        resourceAlibabacloudStackApigatewayVpc(),
			"alibabacloudstack_application_deployment":                resourceAlibabacloudStackEdasApplicationPackageAttachment(),
			"alibabacloudstack_ascm_custom_role":                      resourceAlibabacloudStackAscmRole(),
			"alibabacloudstack_ascm_logon_policy":                     resourceAlibabacloudStackLogonPolicy(),
			"alibabacloudstack_ascm_organization":                     resourceAlibabacloudStackAscmOrganization(),
			"alibabacloudstack_ascm_password_policy":                  resourceAlibabacloudStackAscmPasswordPolicy(),
			"alibabacloudstack_ascm_quota":                            resourceAlibabacloudStackAscmQuota(),
			"alibabacloudstack_ascm_ram_policy":                       resourceAlibabacloudStackAscmRamPolicy(),
			"alibabacloudstack_ascm_ram_policy_for_role":              resourceAlibabacloudStackAscmRamPolicyForRole(),
			"alibabacloudstack_ascm_ram_role":                         resourceAlibabacloudStackAscmRamRole(),
			"alibabacloudstack_ascm_resource_group":                   resourceAlibabacloudStackAscmResourceGroup(),
			"alibabacloudstack_ascm_user":                             resourceAlibabacloudStackAscmUser(),
			"alibabacloudstack_ascm_user_group":                       resourceAlibabacloudStackAscmUserGroup(),
			"alibabacloudstack_ascm_user_group_resource_set_binding":  resourceAlibabacloudStackAscmUserGroupResourceSetBinding(),
			"alibabacloudstack_ascm_user_group_role_binding":          resourceAlibabacloudStackAscmUserGroupRoleBinding(),
			"alibabacloudstack_ascm_user_role_binding":                resourceAlibabacloudStackAscmUserRoleBinding(),
			"alibabacloudstack_ascm_usergroup_user":                   resourceAlibabacloudStackAscmUserGroupUser(),
			"alibabacloudstack_cms_alarm":                             resourceAlibabacloudStackCmsAlarm(),
			"alibabacloudstack_cloudmonitorservice_metricalarmrule":   resourceAlibabacloudStackCmsAlarm(),
			"alibabacloudstack_cms_alarm_contact":                     resourceAlibabacloudstackCmsAlarmContact(),
			"alibabacloudstack_cloudmonitorservice_alarmcontact":      resourceAlibabacloudstackCmsAlarmContact(),
			"alibabacloudstack_cms_alarm_contact_group":               resourceAlibabacloudstackCmsAlarmContactGroup(),
			"alibabacloudstack_cloudmonitorservice_alarmcontactgroup": resourceAlibabacloudstackCmsAlarmContactGroup(),
			"alibabacloudstack_cms_metric_rule_template":              resourceAlibabacloudCmsMetricRuleTemplate(),
			"alibabacloudstack_cms_site_monitor":                      resourceAlibabacloudStackCmsSiteMonitor(),
			"alibabacloudstack_cloudmonitorservice_sitemonitor":       resourceAlibabacloudStackCmsSiteMonitor(),
			"alibabacloudstack_common_bandwidth_package":              resourceAlibabacloudStackCommonBandwidthPackage(),
			"alibabacloudstack_cbwp_commonbandwidthpackage":           resourceAlibabacloudStackCommonBandwidthPackage(),
			"alibabacloudstack_common_bandwidth_package_attachment":   resourceAlibabacloudStackCommonBandwidthPackageAttachment(),
			"alibabacloudstack_cbwp_commonbandwidthpackageattachment": resourceAlibabacloudStackCommonBandwidthPackageAttachment(),
			"alibabacloudstack_cr_ee_namespace":                       resourceAlibabacloudStackCrEENamespace(),
			"alibabacloudstack_cr_ee_repo":                            resourceAlibabacloudStackCrEERepo(),
			"alibabacloudstack_cr_repository":                         resourceAlibabacloudStackCrEERepo(),
			"alibabacloudstack_cr_ee_sync_rule":                       resourceAlibabacloudStackCrEESyncRule(),
			"alibabacloudstack_cr_namespace":                          resourceAlibabacloudStackCRNamespace(),
			"alibabacloudstack_cr_repo":                               resourceAlibabacloudStackCRRepo(),
			"alibabacloudstack_cs_kubernetes":                         resourceAlibabacloudStackCSKubernetes(),
			"alibabacloudstack_ack_cluster":                           resourceAlibabacloudStackCSKubernetes(),
			"alibabacloudstack_cs_kubernetes_node_pool":               resourceAlibabacloudStackCSKubernetesNodePool(),
			"alibabacloudstack_datahub_project":                       resourceAlibabacloudStackDatahubProject(),
			"alibabacloudstack_datahub_subscription":                  resourceAlibabacloudStackDatahubSubscription(),
			"alibabacloudstack_datahub_topic":                         resourceAlibabacloudStackDatahubTopic(),
			"alibabacloudstack_db_account":                            resourceAlibabacloudStackDBAccount(),
			"alibabacloudstack_rds_account":                           resourceAlibabacloudStackDBAccount(),
			"alibabacloudstack_db_account_privilege":                  resourceAlibabacloudStackDBAccountPrivilege(),
			"alibabacloudstack_db_backup_policy":                      resourceAlibabacloudStackDBBackupPolicy(),
			"alibabacloudstack_rds_backuppolicy":                      resourceAlibabacloudStackDBBackupPolicy(),
			"alibabacloudstack_db_connection":                         resourceAlibabacloudStackDBConnection(),
			"alibabacloudstack_rds_dbinstance":                        resourceAlibabacloudStackDBInstance(),
			"alibabacloudstack_db_database":                           resourceAlibabacloudStackDBDatabase(),
			"alibabacloudstack_rds_database":                          resourceAlibabacloudStackDBDatabase(),
			"alibabacloudstack_db_instance":                           resourceAlibabacloudStackDBInstance(),
			"alibabacloudstack_db_read_write_splitting_connection":    resourceAlibabacloudStackDBReadWriteSplittingConnection(),
			"alibabacloudstack_db_readonly_instance":                  resourceAlibabacloudStackDBReadonlyInstance(),
			"alibabacloudstack_disk":                                  resourceAlibabacloudStackDisk(),
			"alibabacloudstack_ecs_disk":                              resourceAlibabacloudStackDisk(),
			"alibabacloudstack_disk_attachment":                       resourceAlibabacloudStackDiskAttachment(),
			"alibabacloudstack_ecs_diskattachment":                    resourceAlibabacloudStackDiskAttachment(),
			"alibabacloudstack_dms_enterprise_instance":               resourceAlibabacloudStackDmsEnterpriseInstance(),
			"alibabacloudstack_dmsenterprise_instance":                resourceAlibabacloudStackDmsEnterpriseInstance(),
			"alibabacloudstack_dms_enterprise_user":                   resourceAlibabacloudStackDmsEnterpriseUser(),
			"alibabacloudstack_dmsenterprise_user":                    resourceAlibabacloudStackDmsEnterpriseUser(),
			"alibabacloudstack_dns_domain":                            resourceAlibabacloudStackDnsDomain(),
			"alibabacloudstack_dns_domain_attachment":                 resourceAlibabacloudStackDnsDomainAttachment(),
			"alibabacloudstack_alidns_domainattachment":               resourceAlibabacloudStackDnsDomainAttachment(),
			"alibabacloudstack_dns_group":                             resourceAlibabacloudStackDnsGroup(),
			"alibabacloudstack_alidns_domaingroup":                    resourceAlibabacloudStackDnsGroup(),
			"alibabacloudstack_dns_record":                            resourceAlibabacloudStackDnsRecord(),
			"alibabacloudstack_drds_instance":                         resourceAlibabacloudStackDRDSInstance(),
			"alibabacloudstack_dts_subscription_job":                  resourceAlibabacloudStackDtsSubscriptionJob(),
			"alibabacloudstack_dts_subscriptionjob":                   resourceAlibabacloudStackDtsSubscriptionJob(),
			"alibabacloudstack_dts_synchronization_instance":          resourceAlibabacloudStackDtsSynchronizationInstance(),
			"alibabacloudstack_dts_synchronizationinstance":           resourceAlibabacloudStackDtsSynchronizationInstance(),
			"alibabacloudstack_dts_synchronization_job":               resourceAlibabacloudStackDtsSynchronizationJob(),
			"alibabacloudstack_ecs_command":                           resourceAlibabacloudStackEcsCommand(),
			"alibabacloudstack_ecs_dedicated_host":                    resourceAlibabacloudStackEcsDedicatedHost(),
			"alibabacloudstack_ecs_dedicatedhost":                     resourceAlibabacloudStackEcsDedicatedHost(),
			"alibabacloudstack_ecs_deployment_set":                    resourceAlibabacloudStackEcsDeploymentSet(),
			"alibabacloudstack_ecs_deploymentset":                     resourceAlibabacloudStackEcsDeploymentSet(),
			"alibabacloudstack_ecs_hpc_cluster":                       resourceAlibabacloudStackEcsHpcCluster(),
			"alibabacloudstack_ecs_hpccluster":                        resourceAlibabacloudStackEcsHpcCluster(),
			"alibabacloudstack_ecs_ebs_storage_set":                   resourceAlibabacloudStackEcsEbsStorageSets(),
			"alibabacloudstack_edas_application":                      resourceAlibabacloudStackEdasApplication(),
			"alibabacloudstack_edas_slbattachment":                    resourceAlibabacloudStackEdasApplication(),
			"alibabacloudstack_edas_application_scale":                resourceAlibabacloudStackEdasInstanceApplicationAttachment(),
			"alibabacloudstack_edas_cluster":                          resourceAlibabacloudStackEdasCluster(),
			"alibabacloudstack_edas_deploy_group":                     resourceAlibabacloudStackEdasDeployGroup(),
			"alibabacloudstack_edas_deploygroup":                      resourceAlibabacloudStackEdasDeployGroup(),
			"alibabacloudstack_edas_instance_cluster_attachment":      resourceAlibabacloudStackEdasInstanceClusterAttachment(),
			"alibabacloudstack_edas_instanceclusterattachment":        resourceAlibabacloudStackEdasInstanceClusterAttachment(),
			"alibabacloudstack_edas_k8s_application":                  resourceAlibabacloudStackEdasK8sApplication(),
			"alibabacloudstack_edas_k8s_cluster":                      resourceAlibabacloudStackEdasK8sCluster(),
			"alibabacloudstack_edas_slb_attachment":                   resourceAlibabacloudStackEdasSlbAttachment(),
			"alibabacloudstack_ehpc_job_template":                     resourceAlibabacloudStackEhpcJobTemplate(),
			"alibabacloudstack_eip":                                   resourceAlibabacloudStackEip(),
			"alibabacloudstack_eip_address":                           resourceAlibabacloudStackEip(),
			"alibabacloudstack_eip_association":                       resourceAlibabacloudStackEipAssociation(),
			"alibabacloudstack_ess_alarm":                             resourceAlibabacloudStackEssAlarm(),
			"alibabacloudstack_autoscaling_alarmtask":                 resourceAlibabacloudStackEssAlarm(),
			"alibabacloudstack_ess_attachment":                        resourceAlibabacloudstackEssAttachment(),
			"alibabacloudstack_ess_lifecycle_hook":                    resourceAlibabacloudStackEssLifecycleHook(),
			"alibabacloudstack_ess_notification":                      resourceAlibabacloudStackEssNotification(),
			"alibabacloudstack_autoscaling_notification":              resourceAlibabacloudStackEssNotification(),
			"alibabacloudstack_ess_scaling_group":                     resourceAlibabacloudStackEssScalingGroup(),
			"alibabacloudstack_ess_scaling_rule":                      resourceAlibabacloudStackEssScalingRule(),
			"alibabacloudstack_ess_scalinggroup_vserver_groups":       resourceAlibabacloudStackEssScalingGroupVserverGroups(),
			"alibabacloudstack_ess_scheduled_task":                    resourceAlibabacloudStackEssScheduledTask(),
			"alibabacloudstack_autoscaling_scheduledtask":             resourceAlibabacloudStackEssScheduledTask(),
			"alibabacloudstack_forward_entry":                         resourceAlibabacloudStackForwardEntry(),
			"alibabacloudstack_natgateway_forwardentry":               resourceAlibabacloudStackForwardEntry(),
			"alibabacloudstack_gpdb_account":                          resourceAlibabacloudStackGpdbAccount(),
			"alibabacloudstack_gpdb_connection":                       resourceAlibabacloudStackGpdbConnection(),
			"alibabacloudstack_gpdb_instance":                         resourceAlibabacloudStackGpdbInstance(),
			"alibabacloudstack_hbase_instance":                        resourceAlibabacloudStackHBaseInstance(),
			"alibabacloudstack_hbase_cluster":                         resourceAlibabacloudStackHBaseInstance(),
			"alibabacloudstack_image":                                 resourceAlibabacloudStackImage(),
			"alibabacloudstack_ecs_image":                             resourceAlibabacloudStackImage(),
			"alibabacloudstack_image_copy":                            resourceAlibabacloudStackImageCopy(),
			"alibabacloudstack_image_export":                          resourceAlibabacloudStackImageExport(),
			"alibabacloudstack_image_import":                          resourceAlibabacloudStackImageImport(),
			"alibabacloudstack_image_share_permission":                resourceAlibabacloudStackImageSharePermission(),
			"alibabacloudstack_instance":                              resourceAlibabacloudStackInstance(),
			"alibabacloudstack_ecs_instance":                          resourceAlibabacloudStackInstance(),
			"alibabacloudstack_key_pair":                              resourceAlibabacloudStackKeyPair(),
			"alibabacloudstack_ecs_keypair":                           resourceAlibabacloudStackKeyPair(),
			"alibabacloudstack_key_pair_attachment":                   resourceAlibabacloudStackKeyPairAttachment(),
			"alibabacloudstack_ecs_keypairattachment":                 resourceAlibabacloudStackKeyPairAttachment(),
			"alibabacloudstack_kms_alias":                             resourceAlibabacloudStackKmsAlias(),
			"alibabacloudstack_kms_ciphertext":                        resourceAlibabacloudStackKmsCiphertext(),
			"alibabacloudstack_kms_key":                               resourceAlibabacloudStackKmsKey(),
			"alibabacloudstack_kms_secret":                            resourceAlibabacloudStackKmsSecret(),
			"alibabacloudstack_kvstore_account":                       resourceAlibabacloudStackKVstoreAccount(),
			"alibabacloudstack_redis_account":                         resourceAlibabacloudStackKVstoreAccount(),
			"alibabacloudstack_kvstore_backup_policy":                 resourceAlibabacloudStackKVStoreBackupPolicy(),
			"alibabacloudstack_kvstore_connection":                    resourceAlibabacloudStackKvstoreConnection(),
			"alibabacloudstack_redis_connection":                      resourceAlibabacloudStackKvstoreConnection(),
			"alibabacloudstack_kvstore_instance":                      resourceAlibabacloudStackKVStoreInstance(),
			"alibabacloudstack_redis_tairinstance":                    resourceAlibabacloudStackKVStoreInstance(),
			"alibabacloudstack_launch_template":                       resourceAlibabacloudStackLaunchTemplate(),
			"alibabacloudstack_ecs_launchtemplate":                    resourceAlibabacloudStackLaunchTemplate(),
			"alibabacloudstack_log_alert":                             resourceAlibabacloudStackLogAlert(),
			"alibabacloudstack_log_machine_group":                     resourceAlibabacloudStackLogMachineGroup(),
			"alibabacloudstack_log_project":                           resourceAlibabacloudStackLogProject(),
			"alibabacloudstack_log_store":                             resourceAlibabacloudStackLogStore(),
			"alibabacloudstack_log_store_index":                       resourceAlibabacloudStackLogStoreIndex(),
			"alibabacloudstack_logtail_attachment":                    resourceAlibabacloudStackLogtailAttachment(),
			"alibabacloudstack_logtail_config":                        resourceAlibabacloudStackLogtailConfig(),
			"alibabacloudstack_maxcompute_project":                    resourceAlibabacloudStackMaxcomputeProject(),
			"alibabacloudstack_maxcompute_user":                       resourceAlibabacloudStackMaxcomputeUser(),
			"alibabacloudstack_maxcompute_cu":                         resourceAlibabacloudStackMaxcomputeCu(),
			"alibabacloudstack_mongodb_instance":                      resourceAlibabacloudStackMongoDBInstance(),
			"alibabacloudstack_mongodb_sharding_instance":             resourceAlibabacloudStackMongoDBShardingInstance(),
			"alibabacloudstack_mongodb_shardinginstance":              resourceAlibabacloudStackMongoDBShardingInstance(),
			"alibabacloudstack_nas_access_group":                      resourceAlibabacloudStackNasAccessGroup(),
			"alibabacloudstack_nas_accessgroup":                       resourceAlibabacloudStackNasAccessGroup(),
			"alibabacloudstack_nas_access_rule":                       resourceAlibabacloudStackNasAccessRule(),
			"alibabacloudstack_nas_accessrule":                        resourceAlibabacloudStackNasAccessRule(),
			"alibabacloudstack_nas_file_system":                       resourceAlibabacloudStackNasFileSystem(),
			"alibabacloudstack_nas_filesystem":                        resourceAlibabacloudStackNasFileSystem(),
			"alibabacloudstack_nas_mount_target":                      resourceAlibabacloudStackNasMountTarget(),
			"alibabacloudstack_nas_mounttarget":                       resourceAlibabacloudStackNasMountTarget(),
			"alibabacloudstack_nat_gateway":                           resourceAlibabacloudStackNatGateway(),
			"alibabacloudstack_natgateway_natgateway":                 resourceAlibabacloudStackNatGateway(),
			"alibabacloudstack_network_acl":                           resourceAlibabacloudStackNetworkAcl(),
			"alibabacloudstack_vpc_networkacl":                        resourceAlibabacloudStackNetworkAcl(),
			"alibabacloudstack_network_acl_attachment":                resourceAlibabacloudStackNetworkAclAttachment(),
			"alibabacloudstack_network_acl_entries":                   resourceAlibabacloudStackNetworkAclEntries(),
			"alibabacloudstack_network_interface":                     resourceAlibabacloudStackNetworkInterface(),
			"alibabacloudstack_ecs_networkinterface":                  resourceAlibabacloudStackNetworkInterface(),
			"alibabacloudstack_network_interface_attachment":          resourceNetworkInterfaceAttachment(),
			"alibabacloudstack_ecs_networkinterfaceattachment":        resourceNetworkInterfaceAttachment(),
			"alibabacloudstack_ons_group":                             resourceAlibabacloudStackOnsGroup(),
			"alibabacloudstack_ons_instance":                          resourceAlibabacloudStackOnsInstance(),
			"alibabacloudstack_ons_topic":                             resourceAlibabacloudStackOnsTopic(),
			"alibabacloudstack_oss_bucket":                            resourceAlibabacloudStackOssBucket(),
			"alibabacloudstack_oss_bucket_quota":                      resourceAlibabacloudStackOssBucketQuota(),
			"alibabacloudstack_oss_bucket_kms":                        resourceAlibabacloudStackOssBucketKms(),
			"alibabacloudstack_oss_bucket_object":                     resourceAlibabacloudStackOssBucketObject(),
			"alibabacloudstack_ots_instance":                          resourceAlibabacloudStackOtsInstance(),
			"alibabacloudstack_ots_instance_attachment":               resourceAlibabacloudStackOtsInstanceAttachment(),
			"alibabacloudstack_ots_instanceattachment":                resourceAlibabacloudStackOtsInstanceAttachment(),
			"alibabacloudstack_ots_table":                             resourceAlibabacloudStackOtsTable(),
			"alibabacloudstack_quick_bi_user":                         resourceAlibabacloudStackQuickBiUser(),
			"alibabacloudstack_quick_bi_user_group":                   resourceAlibabacloudStackQuickBiUserGroup(),
			"alibabacloudstack_quick_bi_workspace":                    resourceAlibabacloudStackQuickBiWorkspace(),
			"alibabacloudstack_ram_role_attachment":                   resourceAlibabacloudStackRamRoleAttachment(),
			"alibabacloudstack_ecs_ramroleattachment":                 resourceAlibabacloudStackRamRoleAttachment(),
			"alibabacloudstack_reserved_instance":                     resourceAlibabacloudStackReservedInstance(),
			"alibabacloudstack_ecs_reservedinstance":                  resourceAlibabacloudStackReservedInstance(),
			"alibabacloudstack_ros_stack":                             resourceAlibabacloudStackRosStack(),
			"alibabacloudstack_ros_template":                          resourceAlibabacloudStackRosTemplate(),
			"alibabacloudstack_route_entry":                           resourceAlibabacloudStackRouteEntry(),
			"alibabacloudstack_route_table":                           resourceAlibabacloudStackRouteTable(),
			"alibabacloudstack_vpc_routetable":                        resourceAlibabacloudStackRouteTable(),
			"alibabacloudstack_route_table_attachment":                resourceAlibabacloudStackRouteTableAttachment(),
			"alibabacloudstack_vpc_routetableattachment":              resourceAlibabacloudStackRouteTableAttachment(),
			"alibabacloudstack_router_interface":                      resourceAlibabacloudStackRouterInterface(),
			"alibabacloudstack_expressconnect_routerinterface":        resourceAlibabacloudStackRouterInterface(),
			"alibabacloudstack_router_interface_connection":           resourceAlibabacloudStackRouterInterfaceConnection(),
			"alibabacloudstack_security_group":                        resourceAlibabacloudStackSecurityGroup(),
			"alibabacloudstack_ecs_securitygroup":                     resourceAlibabacloudStackSecurityGroup(),
			"alibabacloudstack_security_group_rule":                   resourceAlibabacloudStackSecurityGroupRule(),
			"alibabacloudstack_slb":                                   resourceAlibabacloudStackSlb(),
			"alibabacloudstack_slb_loadbalancer":                      resourceAlibabacloudStackSlb(),
			"alibabacloudstack_slb_acl":                               resourceAlibabacloudStackSlbAcl(),
			"alibabacloudstack_slb_accesscontrollist":                 resourceAlibabacloudStackSlbAcl(),
			"alibabacloudstack_slb_backend_server":                    resourceAlibabacloudStackSlbBackendServer(),
			"alibabacloudstack_slb_backendserver":                     resourceAlibabacloudStackSlbBackendServer(),
			"alibabacloudstack_slb_ca_certificate":                    resourceAlibabacloudStackSlbCACertificate(),
			"alibabacloudstack_slb_cacertificate":                     resourceAlibabacloudStackSlbCACertificate(),
			"alibabacloudstack_slb_domain_extension":                  resourceAlibabacloudStackSlbDomainExtension(),
			"alibabacloudstack_slb_domainextension":                   resourceAlibabacloudStackSlbDomainExtension(),
			"alibabacloudstack_slb_listener":                          resourceAlibabacloudStackSlbListener(),
			"alibabacloudstack_slb_master_slave_server_group":         resourceAlibabacloudStackSlbMasterSlaveServerGroup(),
			"alibabacloudstack_slb_masterslaveservergroup":            resourceAlibabacloudStackSlbMasterSlaveServerGroup(),
			"alibabacloudstack_slb_rule":                              resourceAlibabacloudStackSlbRule(),
			"alibabacloudstack_slb_server_certificate":                resourceAlibabacloudStackSlbServerCertificate(),
			"alibabacloudstack_slb_servercertificate":                 resourceAlibabacloudStackSlbServerCertificate(),
			"alibabacloudstack_slb_server_group":                      resourceAlibabacloudStackSlbServerGroup(),
			"alibabacloudstack_slb_vservergroup":                      resourceAlibabacloudStackSlbServerGroup(),
			"alibabacloudstack_snapshot":                              resourceAlibabacloudStackSnapshot(),
			"alibabacloudstack_ecs_snapshot":                          resourceAlibabacloudStackSnapshot(),
			"alibabacloudstack_snapshot_policy":                       resourceAlibabacloudStackSnapshotPolicy(),
			"alibabacloudstack_ecs_autosnapshotpolicy":                resourceAlibabacloudStackSnapshotPolicy(),
			"alibabacloudstack_snat_entry":                            resourceAlibabacloudStackSnatEntry(),
			"alibabacloudstack_natgateway_snatentry":                  resourceAlibabacloudStackSnatEntry(),
			"alibabacloudstack_vpc":                                   resourceAlibabacloudStackVpc(),
			"alibabacloudstack_vpc_vpc":                               resourceAlibabacloudStackVpc(),
			"alibabacloudstack_vpc_ipv6_egress_rule":                  resourceAlibabacloudStackVpcIpv6EgressRule(),
			"alibabacloudstack_vpc_ipv6egressrule":                    resourceAlibabacloudStackVpcIpv6EgressRule(),
			"alibabacloudstack_vpc_ipv6_gateway":                      resourceAlibabacloudStackVpcIpv6Gateway(),
			"alibabacloudstack_vpc_ipv6gateway":                       resourceAlibabacloudStackVpcIpv6Gateway(),
			"alibabacloudstack_vpc_ipv6_internet_bandwidth":           resourceAlibabacloudStackVpcIpv6InternetBandwidth(),
			"alibabacloudstack_vpc_ipv6internetbandwidth":             resourceAlibabacloudStackVpcIpv6InternetBandwidth(),
			"alibabacloudstack_vpn_connection":                        resourceAlibabacloudStackVpnConnection(),
			"alibabacloudstack_vpngateway_vpnconnection":              resourceAlibabacloudStackVpnConnection(),
			"alibabacloudstack_vpn_customer_gateway":                  resourceAlibabacloudStackVpnCustomerGateway(),
			"alibabacloudstack_vpngateway_customergateway":            resourceAlibabacloudStackVpnCustomerGateway(),
			"alibabacloudstack_vpn_gateway":                           resourceAlibabacloudStackVpnGateway(),
			"alibabacloudstack_vpngateway_vpngateway":                 resourceAlibabacloudStackVpnGateway(),
			"alibabacloudstack_vpn_route_entry":                       resourceAlibabacloudStackVpnRouteEntry(),
			"alibabacloudstack_vpngateway_vpnrouteentry":              resourceAlibabacloudStackVpnRouteEntry(),
			"alibabacloudstack_vswitch":                               resourceAlibabacloudStackSwitch(),
			"alibabacloudstack_vpc_vswitch":                           resourceAlibabacloudStackSwitch(),
			"alibabacloudstack_data_works_folder":                     resourceAlibabacloudStackDataWorksFolder(),
			"alibabacloudstack_data_works_connection":                 resourceAlibabacloudStackDataWorksConnection(),
			"alibabacloudstack_data_works_user":                       resourceAlibabacloudStackDataWorksUser(),
			"alibabacloudstack_data_works_project":                    resourceAlibabacloudStackDataWorksProject(),
			"alibabacloudstack_data_works_user_role_binding":          resourceAlibabacloudStackDataWorksUserRoleBinding(),
			"alibabacloudstack_data_works_remind":                     resourceAlibabacloudStackDataWorksRemind(),
			"alibabacloudstack_elasticsearch_instance":                resourceAlibabacloudStackElasticsearch(),
			"alibabacloudstack_dbs_backup_plan":                       resourceAlibabacloudStackDbsBackupPlan(),
			"alibabacloudstack_dbs_backupplan":                        resourceAlibabacloudStackDbsBackupPlan(),
			"alibabacloudstack_express_connect_physical_connection":   resourceAlibabacloudStackExpressConnectPhysicalConnection(),
			"alibabacloudstack_expressconnect_physicalconnection":     resourceAlibabacloudStackExpressConnectPhysicalConnection(),
			"alibabacloudstack_express_connect_virtual_border_router": resourceAlibabacloudStackExpressConnectVirtualBorderRouter(),
			"alibabacloudstack_expressconnect_virtualborderrouter":    resourceAlibabacloudStackExpressConnectVirtualBorderRouter(),
			"alibabacloudstack_oos_template":                          resourceAlibabacloudStackOosTemplate(),
			"alibabacloudstack_oos_execution":                         resourceAlibabacloudStackOosExecution(),
			"alibabacloudstack_arms_alert_contact":                    resourceAlibabacloudStackArmsAlertContact(),
			"alibabacloudstack_arms_alertcontact":                     resourceAlibabacloudStackArmsAlertContact(),
			"alibabacloudstack_arms_alert_contact_group":              resourceAlibabacloudStackArmsAlertContactGroup(),
			"alibabacloudstack_arms_alertcontactgroup":                resourceAlibabacloudStackArmsAlertContactGroup(),
			"alibabacloudstack_arms_dispatch_rule":                    resourceAlibabacloudStackArmsDispatchRule(),
			"alibabacloudstack_arms_dispatchrule":                     resourceAlibabacloudStackArmsDispatchRule(),
			"alibabacloudstack_arms_prometheus_alert_rule":            resourceAlibabacloudStackArmsPrometheusAlertRule(),
			"alibabacloudstack_arms_prometheusalertrule":              resourceAlibabacloudStackArmsPrometheusAlertRule(),
			"alibabacloudstack_elasticsearch_k8s_instance":            resourceAlibabacloudStackElasticsearchOnk8s(),
			"alibabacloudstack_cloud_firewall_control_policy":         resourceAlibabacloudStackCloudFirewallControlPolicy(),
			"alibabacloudstack_cloudfirewall_controlpolicy":           resourceAlibabacloudStackCloudFirewallControlPolicy(),
			"alibabacloudstack_cloud_firewall_control_policy_order":   resourceAlibabacloudStackCloudFirewallControlPolicyOrder(),
			"alibabacloudstack_cloudfirewall_controlpolicyorder":      resourceAlibabacloudStackCloudFirewallControlPolicyOrder(),
			"alibabacloudstack_csb_project":                           resourceAlibabacloudStackCsbProject(),
			"alibabacloudstack_graph_database_db_instance":            resourceAlibabacloudStackGraphDatabaseDbInstance(),
			"alibabacloudstack_graphdatabase_dbinstance":              resourceAlibabacloudStackGraphDatabaseDbInstance(),
		}
	if v, err := stringToBool(os.Getenv("APSARASTACK_IN_ALIBABACLOUDSTACK")); err != nil && !v {
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
	if os.Getenv("TF_EAGLEEYE_TRACEID") != "" && os.Getenv("TF_EAGLEEYE_TRACEID") != "" {
		eagleeye = connectivity.EagleEye{
			TraceId: os.Getenv("TF_EAGLEEYE_TRACEID"),
			RpcId:   os.Getenv("TF_EAGLEEYE_TRACEID"),
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
	ossServicedomain := d.Get("ossservice_domain").(string)
	if ossServicedomain != "" {
		config.Endpoints[connectivity.OssDataCode] = ossServicedomain
	}

	domain := d.Get("domain").(string)
	if domain != "" {
		if strings.Contains(domain, "/") && d.Get("proxy").(string) != "" {
			return nil, fmt.Errorf("[Error]Domain containing the character '/' is not supported for proxy configuration.")
		}
		// 没有生成popgw地址的，继续使用asapi
		var setEndpointIfEmpty = func(endpoint string, domain string) string {
			if endpoint == "" {
				return domain
			}
			return endpoint
		}
		for popcode := range connectivity.PopEndpoints {
			if popcode == connectivity.OssDataCode {
				// oss的数据网关不做配置
				continue
			}
			if popcode == connectivity.SlSDataCode {
				// SLS的数据网关不做配置
				continue
			}
			config.Endpoints[popcode] = setEndpointIfEmpty(config.Endpoints[popcode], domain)
		}
	}
	if v, ok := d.GetOk("popgw_domain"); !d.Get("force_use_asapi").(bool) && ok && v.(string) != "" {
		popgw_domain := v.(string)
		log.Printf("Generator Popgw Endpoint: %s", popgw_domain)
		// 使用各云产品的endpoint的规则生成popgw地址
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
	DbsEndpoint := d.Get("dbs_endpoint").(string)
	if DbsEndpoint != "" {
		config.Endpoints[connectivity.DDSCode] = DbsEndpoint
	}
	DataworkspublicEndpoint := d.Get("dataworkspublic").(string)
	if DataworkspublicEndpoint != "" {
		config.Endpoints[connectivity.DataworkspublicCode] = DataworkspublicEndpoint
	}
	QuickbiEndpoint := d.Get("quickbi_endpoint").(string)
	if QuickbiEndpoint != "" {
		config.Endpoints[connectivity.QuickbiCode] = QuickbiEndpoint
	}
	kafkaOpenApidomain := d.Get("kafkaopenapi_domain").(string)
	if kafkaOpenApidomain != "" {
		config.Endpoints[connectivity.ALIKAFKACode] = kafkaOpenApidomain
	}
	StsEndpoint := d.Get("sts_endpoint").(string)
	if StsEndpoint != "" {
		config.Endpoints[connectivity.STSCode] = StsEndpoint
	}
	organizationAccessKey := d.Get("organization_accesskey").(string)
	if organizationAccessKey != "" {
		config.OrganizationAccessKey = organizationAccessKey
	}
	organizationSecretKey := d.Get("organization_secretkey").(string)
	if organizationSecretKey != "" {
		config.OrganizationSecretKey = organizationSecretKey
	}
	slsOpenAPIEndpoint := d.Get("sls_openapi_endpoint").(string)
	if slsOpenAPIEndpoint != "" {
		config.Endpoints[connectivity.SlSDataCode] = slsOpenAPIEndpoint
	}
	ascmOpenAPIEndpoint := d.Get("ascm_openapi_endpoint").(string)
	if ascmOpenAPIEndpoint != "" {
		config.Endpoints[connectivity.ASCMCode] = ascmOpenAPIEndpoint
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
		dept, _, rgid, err := getResourceCredentials(config)
		if err != nil {
			return nil, err
		}
		config.Department = dept
		config.ResourceGroup = fmt.Sprintf("%d", rgid)
		config.ResourceGroupId = rgid
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

		"domain": "Use this to override the default domain. It's typically used to connect to custom domain.",
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
	request.Scheme = "https" // sts必须是https连接
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
	if len(response.Data) == 0 || response.Code != "200" {
		if len(response.Data) == 0 {
			return "", "", 0, fmt.Errorf("resource group ID and organization not found for resource set %s", config.ResourceSetName)
		}
		return "", "", 0, fmt.Errorf("unable to initialize the ascm client: department or resource_group is not provided")
	} else if len(response.Data) > 1{
		return "", "", 0, fmt.Errorf("There exists a resource group set name with the same name, Please Provider department or resource_group")
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
	// 测试模式下休眠一秒，防止数据缓存导致二次plan失败
	if os.Getenv("TF_ACC") == "1" {
		time.Sleep(time.Duration(second) * time.Second)
	}
}
