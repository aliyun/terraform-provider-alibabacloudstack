package apsarastack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/helper/hashcode"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ACCESS_KEY", os.Getenv("APSARASTACK_ACCESS_KEY")),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_SECRET_KEY", os.Getenv("APSARASTACK_SECRET_KEY")),
				Description: descriptions["secret_key"],
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_REGION", os.Getenv("APSARASTACK_REGION")),
				Description: descriptions["region"],
			},
			"role_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["assume_role_role_arn"],
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ASSUME_ROLE_ARN", os.Getenv("APSARASTACK_ASSUME_ROLE_ARN")),
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_SECURITY_TOKEN", os.Getenv("SECURITY_TOKEN")),
				Description: descriptions["security_token"],
			},
			"ecs_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ECS_ROLE_NAME", os.Getenv("APSARASTACK_ECS_ROLE_NAME")),
				Description: descriptions["ecs_role_name"],
			},
			"skip_region_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: descriptions["skip_region_validation"],
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_PROFILE", ""),
			},
			"endpoints": endpointsSchema(),
			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_credentials_file"],
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_SHARED_CREDENTIALS_FILE", ""),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_INSECURE", nil),
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
				Default:      "HTTP",
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
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_SOURCE_IP", os.Getenv("APSARASTACK_SOURCE_IP")),
				Description: descriptions["source_ip"],
			},
			"security_transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_SECURITY_TRANSPORT", os.Getenv("APSARASTACK_SECURITY_TRANSPORT")),
				//Deprecated:  "It has been deprecated from version 1.136.0 and using new field secure_transport instead.",
			},
			"secure_transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_SECURE_TRANSPORT", os.Getenv("APSARASTACK_SECURE_TRANSPORT")),
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
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_PROXY", nil),
				Description: descriptions["proxy"],
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_DOMAIN", nil),
				Description: descriptions["domain"],
			},
			"ossservice_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_OSSSERVICE_DOMAIN", nil),
				Description: descriptions["ossservice_domain"],
			},
			"kafkaopenapi_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_KAFKAOPENAPI_DOMAIN", nil),
				Description: descriptions["kafkaopenapi_domain"],
			},
			"organization_accesskey": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ORGANIZATION_ACCESSKEY", nil),
				Description: descriptions["organization_accesskey"],
			},
			"organization_secretkey": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ORGANIZATION_SECRETKEY", nil),
				Description: descriptions["organization_secretkey"],
			},
			"sls_openapi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_SLS_OPENAPI_ENDPOINT", nil),
				Description: descriptions["sls_openapi_endpoint"],
			},
			"sts_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_STS_ENDPOINT", os.Getenv("APSARASTACK_STS_ENDPOINT")),
				Description: descriptions["sts_endpoint"],
			},
			"quickbi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_QUICKBI_ENDPOINT", nil),
				Description: descriptions["quickbi_endpoint"],
			},
			"department": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_DEPARTMENT", nil),
				Description: descriptions["department"],
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_RESOURCE_GROUP", nil),
				Description: descriptions["resource_group"],
			},
			"resource_group_set_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_RESOURCE_GROUP_SET", nil),
				Description: descriptions["resource_group_set_name"],
			},
			"dataworkspublic": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_DATAWORKS_PUBLIC_ENDPOINT", nil),
				Description: descriptions["dataworkspublic_endpoint"],
			},
			"dbs_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_DBS_ENDPOINT", nil),
				Description: descriptions["dbs_endpoint"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"apsarastack_account":                              dataSourceApsaraStackAccount(),
			"apsarastack_adb_clusters":                         dataSourceApsaraStackAdbDbClusters(),
			"apsarastack_adb_zones":                            dataSourceApsaraStackAdbZones(),
			"apsarastack_adb_db_clusters":                      dataSourceApsaraStackAdbDbClusters(),
			"apsarastack_api_gateway_apis":                     dataSourceApsaraStackApiGatewayApis(),
			"apsarastack_api_gateway_apps":                     dataSourceApsaraStackApiGatewayApps(),
			"apsarastack_api_gateway_groups":                   dataSourceApsaraStackApiGatewayGroups(),
			"apsarastack_api_gateway_service":                  dataSourceApsaraStackApiGatewayService(),
			"apsarastack_ascm_resource_groups":                 dataSourceApsaraStackAscmResourceGroups(),
			"apsarastack_ascm_users":                           dataSourceApsaraStackAscmUsers(),
			"apsarastack_ascm_user_groups":                     dataSourceApsaraStackAscmUserGroups(),
			"apsarastack_ascm_logon_policies":                  dataSourceApsaraStackAscmLogonPolicies(),
			"apsarastack_ascm_ram_service_roles":               dataSourceApsaraStackAscmRamServiceRoles(),
			"apsarastack_ascm_organizations":                   dataSourceApsaraStackAscmOrganizations(),
			"apsarastack_ascm_instance_families":               dataSourceApsaraStackInstanceFamilies(),
			"apsarastack_ascm_regions_by_product":              dataSourceApsaraStackRegionsByProduct(),
			"apsarastack_ascm_service_cluster_by_product":      dataSourceApsaraStackServiceClusterByProduct(),
			"apsarastack_ascm_ecs_instance_families":           dataSourceApsaraStackEcsInstanceFamilies(),
			"apsarastack_ascm_specific_fields":                 dataSourceApsaraStackSpecificFields(),
			"apsarastack_ascm_environment_services_by_product": dataSourceApsaraStackAscmEnvironmentServicesByProduct(),
			"apsarastack_ascm_password_policies":               dataSourceApsaraStackAscmPasswordPolicies(),
			"apsarastack_ascm_quotas":                          dataSourceApsaraStackQuotas(),
			"apsarastack_ascm_metering_query_ecs":              dataSourceApsarastackAscmMeteringQueryEcs(),
			"apsarastack_ascm_roles":                           dataSourceApsaraStackAscmRoles(),
			"apsarastack_ascm_ram_policies":                    dataSourceApsaraStackAscmRamPolicies(),
			"apsarastack_ascm_ram_policies_for_user":           dataSourceApsaraStackAscmRamPoliciesForUser(),
			"apsarastack_common_bandwidth_packages":            dataSourceApsaraStackCommonBandwidthPackages(),
			"apsarastack_cr_ee_instances":                      dataSourceApsaraStackCrEEInstances(),
			"apsarastack_cr_ee_namespaces":                     dataSourceApsaraStackCrEENamespaces(),
			"apsarastack_cr_ee_repos":                          dataSourceApsaraStackCrEERepos(),
			"apsarastack_cr_ee_sync_rules":                     dataSourceApsaraStackCrEESyncRules(),
			"apsarastack_cr_namespaces":                        dataSourceApsaraStackCRNamespaces(),
			"apsarastack_cr_repos":                             dataSourceApsaraStackCRRepos(),
			"apsarastack_cs_kubernetes_clusters":               dataSourceApsaraStackCSKubernetesClusters(),
			"apsarastack_cms_alarm_contacts":                   dataSourceApsarastackCmsAlarmContacts(),
			"apsarastack_cms_alarm_contact_groups":             dataSourceApsarastackCmsAlarmContactGroups(),
			"apsarastack_cms_project_meta":                     dataSourceApsarastackCmsProjectMeta(),
			"apsarastack_cms_metric_metalist":                  dataSourceApsarastackCmsMetricMetalist(),
			"apsarastack_cms_alarms":                           dataSourceApsarastackCmsAlarms(),
			"apsarastack_datahub_service":                      dataSourceApsaraStackDatahubService(),
			"apsarastack_db_instances":                         dataSourceApsaraStackDBInstances(),
			"apsarastack_db_zones":                             dataSourceApsaraStackDBZones(),
			"apsarastack_disks":                                dataSourceApsaraStackDisks(),
			"apsarastack_dns_records":                          dataSourceApsaraStackDnsRecords(),
			"apsarastack_dns_groups":                           dataSourceApsaraStackDnsGroups(),
			"apsarastack_dns_domains":                          dataSourceApsaraStackDnsDomains(),
			"apsarastack_drds_instances":                       dataSourceApsaraStackDRDSInstances(),
			"apsarastack_dms_enterprise_instances":             dataSourceApsaraStackDmsEnterpriseInstances(),
			"apsarastack_dms_enterprise_users":                 dataSourceApsaraStackDmsEnterpriseUsers(),
			"apsarastack_ecs_commands":                         dataSourceApsaraStackEcsCommands(),
			"apsarastack_ecs_deployment_sets":                  dataSourceApsaraStackEcsDeploymentSets(),
			"apsarastack_ecs_hpc_clusters":                     dataSourceApsaraStackEcsHpcClusters(),
			"apsarastack_ecs_dedicated_hosts":                  dataSourceApsaraStackEcsDedicatedHosts(),
			"apsarastack_edas_deploy_groups":                   dataSourceApsaraStackEdasDeployGroups(),
			"apsarastack_edas_clusters":                        dataSourceApsaraStackEdasClusters(),
			"apsarastack_edas_applications":                    dataSourceApsaraStackEdasApplications(),
			"apsarastack_eips":                                 dataSourceApsaraStackEips(),
			"apsarastack_ess_scaling_configurations":           dataSourceApsaraStackEssScalingConfigurations(),
			"apsarastack_ess_scaling_groups":                   dataSourceApsaraStackEssScalingGroups(),
			"apsarastack_ess_lifecycle_hooks":                  dataSourceApsaraStackEssLifecycleHooks(),
			"apsarastack_ess_notifications":                    dataSourceApsaraStackEssNotifications(),
			"apsarastack_ess_scaling_rules":                    dataSourceApsaraStackEssScalingRules(),
			"apsarastack_ess_scheduled_tasks":                  dataSourceApsaraStackEssScheduledTasks(),
			"apsarastack_forward_entries":                      dataSourceApsaraStackForwardEntries(),
			"apsarastack_gpdb_accounts":                        dataSourceApsaraStackGpdbAccounts(),
			"apsarastack_gpdb_instances":                       dataSourceApsaraStackGpdbInstances(),
			"apsarastack_hbase_instances":                      dataSourceApsaraStackHBaseInstances(),
			"apsarastack_instances":                            dataSourceApsaraStackInstances(),
			"apsarastack_instance_type_families":               dataSourceApsaraStackInstanceTypeFamilies(),
			"apsarastack_instance_types":                       dataSourceApsaraStackInstanceTypes(),
			"apsarastack_images":                               dataSourceApsaraStackImages(),
			"apsarastack_key_pairs":                            dataSourceApsaraStackKeyPairs(),
			"apsarastack_kms_aliases":                          dataSourceApsaraStackKmsAliases(),
			"apsarastack_kms_ciphertext":                       dataSourceApsaraStackKmsCiphertext(),
			"apsarastack_kms_keys":                             dataSourceApsaraStackKmsKeys(),
			"apsarastack_kms_secrets":                          dataSourceApsaraStackKmsSecrets(),
			"apsarastack_kvstore_instances":                    dataSourceApsaraStackKVStoreInstances(),
			"apsarastack_kvstore_zones":                        dataSourceApsaraStackKVStoreZones(),
			"apsarastack_kvstore_instance_classes":             dataSourceApsaraStackKVStoreInstanceClasses(),
			"apsarastack_kvstore_instance_engines":             dataSourceApsaraStackKVStoreInstanceEngines(),
			"apsarastack_mongodb_instances":                    dataSourceApsaraStackMongoDBInstances(),
			"apsarastack_mongodb_zones":                        dataSourceApsaraStackMongoDBZones(),
			"apsarastack_maxcompute_cus":                       dataSourceApsaraStackMaxcomputeCus(),
			"apsarastack_maxcompute_users":                     dataSourceApsaraStackMaxcomputeUsers(),
			"apsarastack_maxcompute_clusters":                  dataSourceApsaraStackMaxcomputeClusters(),
			"apsarastack_maxcompute_cluster_qutaos":            dataSourceApsaraStackMaxcomputeClusterQutaos(),
			"apsarastack_maxcompute_projects":                  dataSourceApsaraStackMaxcomputeProjects(),
			"apsarastack_nas_zones":                            dataSourceApsaraStackNasZones(),
			"apsarastack_nas_protocols":                        dataSourceApsaraStackNasProtocols(),
			"apsarastack_nas_file_systems":                     dataSourceApsaraStackFileSystems(),
			"apsarastack_nas_mount_targets":                    dataSourceApsaraStackNasMountTargets(),
			"apsarastack_nas_access_rules":                     dataSourceApsaraStackAccessRules(),
			"apsarastack_nat_gateways":                         dataSourceApsaraStackNatGateways(),
			"apsarastack_network_acls":                         dataSourceApsaraStackNetworkAcls(),
			"apsarastack_network_interfaces":                   dataSourceApsaraStackNetworkInterfaces(),
			"apsarastack_oss_buckets":                          dataSourceApsaraStackOssBuckets(),
			"apsarastack_oss_bucket_objects":                   dataSourceApsaraStackOssBucketObjects(),
			"apsarastack_ons_instances":                        dataSourceApsaraStackOnsInstances(),
			"apsarastack_ons_topics":                           dataSourceApsaraStackOnsTopics(),
			"apsarastack_ons_groups":                           dataSourceApsaraStackOnsGroups(),
			"apsarastack_ots_tables":                           dataSourceApsaraStackOtsTables(),
			"apsarastack_ots_instances":                        dataSourceApsaraStackOtsInstances(),
			"apsarastack_ots_instances_attachment":             dataSourceApsaraStackOtsInstanceAttachments(),
			"apsarastack_ots_service":                          dataSourceApsaraStackOtsService(),
			"apsarastack_quick_bi_users":                       dataSourceApsaraStackQuickBiUsers(),
			"apsarastack_router_interfaces":                    dataSourceApsaraStackRouterInterfaces(),
			"apsarastack_ram_service_role_products":            dataSourceApsarastackRamServiceRoleProducts(),
			"apsarastack_route_tables":                         dataSourceApsaraStackRouteTables(),
			"apsarastack_route_entries":                        dataSourceApsaraStackRouteEntries(),
			"apsarastack_ros_stacks":                           dataSourceApsaraStackRosStacks(),
			"apsarastack_ros_templates":                        dataSourceApsaraStackRosTemplates(),
			"apsarastack_security_groups":                      dataSourceApsaraStackSecurityGroups(),
			"apsarastack_security_group_rules":                 dataSourceApsaraStackSecurityGroupRules(),
			"apsarastack_snapshots":                            dataSourceApsaraStackSnapshots(),
			"apsarastack_slb_listeners":                        dataSourceApsaraStackSlbListeners(),
			"apsarastack_slb_server_groups":                    dataSourceApsaraStackSlbServerGroups(),
			"apsarastack_slb_acls":                             dataSourceApsaraStackSlbAcls(),
			"apsarastack_slb_domain_extensions":                dataSourceApsaraStackSlbDomainExtensions(),
			"apsarastack_slb_rules":                            dataSourceApsaraStackSlbRules(),
			"apsarastack_slb_master_slave_server_groups":       dataSourceApsaraStackSlbMasterSlaveServerGroups(),
			"apsarastack_slbs":                                 dataSourceApsaraStackSlbs(),
			"apsarastack_slb_zones":                            dataSourceApsaraStackSlbZones(),
			"apsarastack_snat_entries":                         dataSourceApsaraStackSnatEntries(),
			"apsarastack_slb_server_certificates":              dataSourceApsaraStackSlbServerCertificates(),
			"apsarastack_slb_ca_certificates":                  dataSourceApsaraStackSlbCACertificates(),
			"apsarastack_slb_backend_servers":                  dataSourceApsaraStackSlbBackendServers(),
			"apsarastack_tsdb_zones":                           dataSourceApsaraStackTsdbZones(),
			"apsarastack_vpn_gateways":                         dataSourceApsaraStackVpnGateways(),
			"apsarastack_vpn_customer_gateways":                dataSourceApsaraStackVpnCustomerGateways(),
			"apsarastack_vpn_connections":                      dataSourceApsaraStackVpnConnections(),
			"apsarastack_vpc_ipv6_gateways":                    dataSourceApsaraStackVpcIpv6Gateways(),
			"apsarastack_vpc_ipv6_egress_rules":                dataSourceApsaraStackVpcIpv6EgressRules(),
			"apsarastack_vpc_ipv6_addresses":                   dataSourceApsaraStackVpcIpv6Addresses(),
			"apsarastack_vpc_ipv6_internet_bandwidths":         dataSourceApsaraStackVpcIpv6InternetBandwidths(),
			"apsarastack_vswitches":                            dataSourceApsaraStackVSwitches(),
			"apsarastack_vpcs":                                 dataSourceApsaraStackVpcs(),
			"apsarastack_zones":                                dataSourceApsaraStackZones(),
			"apsarastack_elasticsearch_instances":              dataSourceApsaraStackElasticsearch(),
			"apsarastack_elasticsearch_zones":                  dataSourceApsaraStackElaticsearchZones(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"apsarastack_ess_scaling_configuration":            resourceApsaraStackEssScalingConfiguration(),
			"apsarastack_adb_account":                          resourceApsaraStackAdbAccount(),
			"apsarastack_adb_backup_policy":                    resourceApsaraStackAdbBackupPolicy(),
			"apsarastack_adb_cluster":                          resourceApsaraStackAdbDbCluster(),
			"apsarastack_adb_connection":                       resourceApsaraStackAdbConnection(),
			"apsarastack_adb_db_cluster":                       resourceApsaraStackAdbDbCluster(),
			"apsarastack_alikafka_sasl_acl":                    resourceApsaraStackAlikafkaSaslAcl(),
			"apsarastack_alikafka_sasl_user":                   resourceApsaraStackAlikafkaSaslUser(),
			"apsarastack_alikafka_topic":                       resourceApsaraStackAlikafkaTopic(),
			"apsarastack_api_gateway_api":                      resourceApsaraStackApigatewayApi(),
			"apsarastack_api_gateway_app":                      resourceApsaraStackApigatewayApp(),
			"apsarastack_api_gateway_app_attachment":           resourceAliyunApigatewayAppAttachment(),
			"apsarastack_api_gateway_group":                    resourceApsaraStackApigatewayGroup(),
			"apsarastack_api_gateway_vpc_access":               resourceApsaraStackApigatewayVpc(),
			"apsarastack_application_deployment":               resourceApsaraStackEdasApplicationPackageAttachment(),
			"apsarastack_ascm_custom_role":                     resourceApsaraStackAscmRole(),
			"apsarastack_ascm_logon_policy":                    resourceApsaraStackLogonPolicy(),
			"apsarastack_ascm_organization":                    resourceApsaraStackAscmOrganization(),
			"apsarastack_ascm_password_policy":                 resourceApsaraStackAscmPasswordPolicy(),
			"apsarastack_ascm_quota":                           resourceApsaraStackAscmQuota(),
			"apsarastack_ascm_ram_policy":                      resourceApsaraStackAscmRamPolicy(),
			"apsarastack_ascm_ram_policy_for_role":             resourceApsaraStackAscmRamPolicyForRole(),
			"apsarastack_ascm_ram_role":                        resourceApsaraStackAscmRamRole(),
			"apsarastack_ascm_resource_group":                  resourceApsaraStackAscmResourceGroup(),
			"apsarastack_ascm_user":                            resourceApsaraStackAscmUser(),
			"apsarastack_ascm_user_group":                      resourceApsaraStackAscmUserGroup(),
			"apsarastack_ascm_user_group_resource_set_binding": resourceApsaraStackAscmUserGroupResourceSetBinding(),
			"apsarastack_ascm_user_group_role_binding":         resourceApsaraStackAscmUserGroupRoleBinding(),
			"apsarastack_ascm_user_role_binding":               resourceApsaraStackAscmUserRoleBinding(),
			"apsarastack_ascm_usergroup_user":                  resourceApsaraStackAscmUserGroupUser(),
			"apsarastack_cms_alarm":                            resourceApsaraStackCmsAlarm(),
			"apsarastack_cms_alarm_contact":                    resourceApsarastackCmsAlarmContact(),
			"apsarastack_cms_alarm_contact_group":              resourceApsarastackCmsAlarmContactGroup(),
			"apsarastack_cms_site_monitor":                     resourceApsaraStackCmsSiteMonitor(),
			"apsarastack_common_bandwidth_package":             resourceApsaraStackCommonBandwidthPackage(),
			"apsarastack_common_bandwidth_package_attachment":  resourceApsaraStackCommonBandwidthPackageAttachment(),
			"apsarastack_cr_ee_namespace":                      resourceApsaraStackCrEENamespace(),
			"apsarastack_cr_ee_repo":                           resourceApsaraStackCrEERepo(),
			"apsarastack_cr_ee_sync_rule":                      resourceApsaraStackCrEESyncRule(),
			"apsarastack_cr_namespace":                         resourceApsaraStackCRNamespace(),
			"apsarastack_cr_repo":                              resourceApsaraStackCRRepo(),
			"apsarastack_cs_kubernetes":                        resourceApsaraStackCSKubernetes(),
			"apsarastack_cs_kubernetes_node_pool":              resourceApsaraStackCSKubernetesNodePool(),
			"apsarastack_datahub_project":                      resourceApsaraStackDatahubProject(),
			"apsarastack_datahub_subscription":                 resourceApsaraStackDatahubSubscription(),
			"apsarastack_datahub_topic":                        resourceApsaraStackDatahubTopic(),
			"apsarastack_db_account":                           resourceApsaraStackDBAccount(),
			"apsarastack_db_account_privilege":                 resourceApsaraStackDBAccountPrivilege(),
			"apsarastack_db_backup_policy":                     resourceApsaraStackDBBackupPolicy(),
			"apsarastack_db_connection":                        resourceApsaraStackDBConnection(),
			"apsarastack_db_database":                          resourceApsaraStackDBDatabase(),
			"apsarastack_db_instance":                          resourceApsaraStackDBInstance(),
			"apsarastack_db_read_write_splitting_connection":   resourceApsaraStackDBReadWriteSplittingConnection(),
			"apsarastack_db_readonly_instance":                 resourceApsaraStackDBReadonlyInstance(),
			"apsarastack_disk":                                 resourceApsaraStackDisk(),
			"apsarastack_disk_attachment":                      resourceApsaraStackDiskAttachment(),
			"apsarastack_dms_enterprise_instance":              resourceApsaraStackDmsEnterpriseInstance(),
			"apsarastack_dms_enterprise_user":                  resourceApsaraStackDmsEnterpriseUser(),
			"apsarastack_dns_domain":                           resourceApsaraStackDnsDomain(),
			"apsarastack_dns_domain_attachment":                resourceApsaraStackDnsDomainAttachment(),
			"apsarastack_dns_group":                            resourceApsaraStackDnsGroup(),
			"apsarastack_dns_record":                           resourceApsaraStackDnsRecord(),
			"apsarastack_drds_instance":                        resourceApsaraStackDRDSInstance(),
			"apsarastack_dts_subscription_job":                 resourceApsaraStackDtsSubscriptionJob(),
			"apsarastack_dts_synchronization_instance":         resourceApsaraStackDtsSynchronizationInstance(),
			"apsarastack_dts_synchronization_job":              resourceApsaraStackDtsSynchronizationJob(),
			"apsarastack_ecs_command":                          resourceApsaraStackEcsCommand(),
			"apsarastack_ecs_dedicated_host":                   resourceApsaraStackEcsDedicatedHost(),
			"apsarastack_ecs_deployment_set":                   resourceApsaraStackEcsDeploymentSet(),
			"apsarastack_ecs_hpc_cluster":                      resourceApsaraStackEcsHpcCluster(),
			"apsarastack_ecs_ebs_storage_set":                  resourceApsaraStackEcsEbsStorageSets(),
			"apsarastack_edas_application":                     resourceApsaraStackEdasApplication(),
			"apsarastack_edas_application_scale":               resourceApsaraStackEdasInstanceApplicationAttachment(),
			"apsarastack_edas_cluster":                         resourceApsaraStackEdasCluster(),
			"apsarastack_edas_deploy_group":                    resourceApsaraStackEdasDeployGroup(),
			"apsarastack_edas_instance_cluster_attachment":     resourceApsaraStackEdasInstanceClusterAttachment(),
			"apsarastack_edas_k8s_application":                 resourceApsaraStackEdasK8sApplication(),
			"apsarastack_edas_k8s_cluster":                     resourceApsaraStackEdasK8sCluster(),
			"apsarastack_edas_slb_attachment":                  resourceApsaraStackEdasSlbAttachment(),
			"apsarastack_ehpc_job_template":                    resourceApsaraStackEhpcJobTemplate(),
			"apsarastack_eip":                                  resourceApsaraStackEip(),
			"apsarastack_eip_association":                      resourceApsaraStackEipAssociation(),
			"apsarastack_ess_alarm":                            resourceApsaraStackEssAlarm(),
			"apsarastack_ess_attachment":                       resourceApsarastackEssAttachment(),
			"apsarastack_ess_lifecycle_hook":                   resourceApsaraStackEssLifecycleHook(),
			"apsarastack_ess_notification":                     resourceApsaraStackEssNotification(),
			"apsarastack_ess_scaling_group":                    resourceApsaraStackEssScalingGroup(),
			"apsarastack_ess_scaling_rule":                     resourceApsaraStackEssScalingRule(),
			"apsarastack_ess_scalinggroup_vserver_groups":      resourceApsaraStackEssScalingGroupVserverGroups(),
			"apsarastack_ess_scheduled_task":                   resourceApsaraStackEssScheduledTask(),
			"apsarastack_forward_entry":                        resourceApsaraStackForwardEntry(),
			"apsarastack_gpdb_account":                         resourceApsaraStackGpdbAccount(),
			"apsarastack_gpdb_connection":                      resourceApsaraStackGpdbConnection(),
			"apsarastack_gpdb_instance":                        resourceApsaraStackGpdbInstance(),
			"apsarastack_hbase_instance":                       resourceApsaraStackHBaseInstance(),
			"apsarastack_image":                                resourceApsaraStackImage(),
			"apsarastack_image_copy":                           resourceApsaraStackImageCopy(),
			"apsarastack_image_export":                         resourceApsaraStackImageExport(),
			"apsarastack_image_import":                         resourceApsaraStackImageImport(),
			"apsarastack_image_share_permission":               resourceApsaraStackImageSharePermission(),
			"apsarastack_instance":                             resourceApsaraStackInstance(),
			"apsarastack_key_pair":                             resourceApsaraStackKeyPair(),
			"apsarastack_key_pair_attachment":                  resourceApsaraStackKeyPairAttachment(),
			"apsarastack_kms_alias":                            resourceApsaraStackKmsAlias(),
			"apsarastack_kms_ciphertext":                       resourceApsaraStackKmsCiphertext(),
			"apsarastack_kms_key":                              resourceApsaraStackKmsKey(),
			"apsarastack_kms_secret":                           resourceApsaraStackKmsSecret(),
			"apsarastack_kvstore_account":                      resourceApsaraStackKVstoreAccount(),
			"apsarastack_kvstore_backup_policy":                resourceApsaraStackKVStoreBackupPolicy(),
			"apsarastack_kvstore_connection":                   resourceApsaraStackKvstoreConnection(),
			"apsarastack_kvstore_instance":                     resourceApsaraStackKVStoreInstance(),
			"apsarastack_launch_template":                      resourceApsaraStackLaunchTemplate(),
			"apsarastack_log_machine_group":                    resourceApsaraStackLogMachineGroup(),
			"apsarastack_log_project":                          resourceApsaraStackLogProject(),
			"apsarastack_log_store":                            resourceApsaraStackLogStore(),
			"apsarastack_log_store_index":                      resourceApsaraStackLogStoreIndex(),
			"apsarastack_logtail_attachment":                   resourceApsaraStackLogtailAttachment(),
			"apsarastack_logtail_config":                       resourceApsaraStackLogtailConfig(),
			"apsarastack_maxcompute_cu":                        resourceApsaraStackMaxcomputeCu(),
			"apsarastack_maxcompute_project":                   resourceApsaraStackMaxcomputeProject(),
			"apsarastack_maxcompute_user":                      resourceApsaraStackMaxcomputeUser(),
			"apsarastack_mongodb_instance":                     resourceApsaraStackMongoDBInstance(),
			"apsarastack_mongodb_sharding_instance":            resourceApsaraStackMongoDBShardingInstance(),
			"apsarastack_nas_access_group":                     resourceApsaraStackNasAccessGroup(),
			"apsarastack_nas_access_rule":                      resourceApsaraStackNasAccessRule(),
			"apsarastack_nas_file_system":                      resourceApsaraStackNasFileSystem(),
			"apsarastack_nas_mount_target":                     resourceApsaraStackNasMountTarget(),
			"apsarastack_nat_gateway":                          resourceApsaraStackNatGateway(),
			"apsarastack_network_acl":                          resourceApsaraStackNetworkAcl(),
			"apsarastack_network_acl_attachment":               resourceApsaraStackNetworkAclAttachment(),
			"apsarastack_network_acl_entries":                  resourceApsaraStackNetworkAclEntries(),
			"apsarastack_network_interface":                    resourceApsaraStackNetworkInterface(),
			"apsarastack_network_interface_attachment":         resourceNetworkInterfaceAttachment(),
			"apsarastack_ons_group":                            resourceApsaraStackOnsGroup(),
			"apsarastack_ons_instance":                         resourceApsaraStackOnsInstance(),
			"apsarastack_ons_topic":                            resourceApsaraStackOnsTopic(),
			"apsarastack_oss_bucket":                           resourceApsaraStackOssBucket(),
			"apsarastack_oss_bucket_kms":                       resourceApsaraStackOssBucketKms(),
			"apsarastack_oss_bucket_object":                    resourceApsaraStackOssBucketObject(),
			"apsarastack_ots_instance":                         resourceApsaraStackOtsInstance(),
			"apsarastack_ots_instance_attachment":              resourceApsaraStackOtsInstanceAttachment(),
			"apsarastack_ots_table":                            resourceApsaraStackOtsTable(),
			"apsarastack_quick_bi_user":                        resourceApsaraStackQuickBiUser(),
			"apsarastack_quick_bi_user_group":                  resourceApsaraStackQuickBiUserGroup(),
			"apsarastack_quick_bi_workspace":                   resourceApsaraStackQuickBiWorkspace(),
			"apsarastack_ram_role_attachment":                  resourceApsaraStackRamRoleAttachment(),
			"apsarastack_reserved_instance":                    resourceApsaraStackReservedInstance(),
			"apsarastack_ros_stack":                            resourceApsaraStackRosStack(),
			"apsarastack_ros_template":                         resourceApsaraStackRosTemplate(),
			"apsarastack_route_entry":                          resourceApsaraStackRouteEntry(),
			"apsarastack_route_table":                          resourceApsaraStackRouteTable(),
			"apsarastack_route_table_attachment":               resourceApsaraStackRouteTableAttachment(),
			"apsarastack_router_interface":                     resourceApsaraStackRouterInterface(),
			"apsarastack_router_interface_connection":          resourceApsaraStackRouterInterfaceConnection(),
			"apsarastack_security_group":                       resourceApsaraStackSecurityGroup(),
			"apsarastack_security_group_rule":                  resourceApsaraStackSecurityGroupRule(),
			"apsarastack_slb":                                  resourceApsaraStackSlb(),
			"apsarastack_slb_acl":                              resourceApsaraStackSlbAcl(),
			"apsarastack_slb_backend_server":                   resourceApsaraStackSlbBackendServer(),
			"apsarastack_slb_ca_certificate":                   resourceApsaraStackSlbCACertificate(),
			"apsarastack_slb_domain_extension":                 resourceApsaraStackSlbDomainExtension(),
			"apsarastack_slb_listener":                         resourceApsaraStackSlbListener(),
			"apsarastack_slb_master_slave_server_group":        resourceApsaraStackSlbMasterSlaveServerGroup(),
			"apsarastack_slb_rule":                             resourceApsaraStackSlbRule(),
			"apsarastack_slb_server_certificate":               resourceApsaraStackSlbServerCertificate(),
			"apsarastack_slb_server_group":                     resourceApsaraStackSlbServerGroup(),
			"apsarastack_snapshot":                             resourceApsaraStackSnapshot(),
			"apsarastack_snapshot_policy":                      resourceApsaraStackSnapshotPolicy(),
			"apsarastack_snat_entry":                           resourceApsaraStackSnatEntry(),
			"apsarastack_vpc":                                  resourceApsaraStackVpc(),
			"apsarastack_vpc_ipv6_egress_rule":                 resourceApsaraStackVpcIpv6EgressRule(),
			"apsarastack_vpc_ipv6_gateway":                     resourceApsaraStackVpcIpv6Gateway(),
			"apsarastack_vpc_ipv6_internet_bandwidth":          resourceApsaraStackVpcIpv6InternetBandwidth(),
			"apsarastack_vpn_connection":                       resourceApsaraStackVpnConnection(),
			"apsarastack_vpn_customer_gateway":                 resourceApsaraStackVpnCustomerGateway(),
			"apsarastack_vpn_gateway":                          resourceApsaraStackVpnGateway(),
			"apsarastack_vpn_route_entry":                      resourceApsaraStackVpnRouteEntry(),
			"apsarastack_vswitch":                              resourceApsaraStackSwitch(),
			"apsarastack_data_works_folder":                    resourceApsaraStackDataWorksFolder(),
			"apsarastack_data_works_connection":                resourceApsaraStackDataWorksConnection(),
			"apsarastack_data_works_user":                      resourceApsaraStackDataWorksUser(),
			"apsarastack_data_works_project":                   resourceApsaraStackDataWorksProject(),
			"apsarastack_data_works_user_role_binding":         resourceApsaraStackDataWorksUserRoleBinding(),
			"apsarastack_data_works_remind":                    resourceApsaraStackDataWorksRemind(),
			"apsarastack_elasticsearch_instance":               resourceApsaraStackElasticsearch(),
			"apsarastack_dbs_backup_plan":                      resourceApsaraStackDbsBackupPlan(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var providerConfig map[string]interface{}

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
	if region == "" {
		region = DEFAULT_REGION
	}

	ecsRoleName := getProviderConfig(d.Get("ecs_role_name").(string), "ram_role_name")

	config := &connectivity.Config{
		AccessKey:            strings.TrimSpace(accessKey),
		SecretKey:            strings.TrimSpace(secretKey),
		EcsRoleName:          strings.TrimSpace(ecsRoleName),
		Region:               connectivity.Region(strings.TrimSpace(region)),
		RegionId:             strings.TrimSpace(region),
		SkipRegionValidation: d.Get("skip_region_validation").(bool),
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
	}
	if v, ok := d.GetOk("security_transport"); config.SecureTransport == "" && ok && v.(string) != "" {
		config.SecureTransport = v.(string)
	}
	token := getProviderConfig(d.Get("security_token").(string), "sts_token")
	config.SecurityToken = strings.TrimSpace(token)
	config.RamRoleArn = getProviderConfig(d.Get("role_arn").(string), "ram_role_arn")
	log.Printf("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$led!!! %s", config.RamRoleArn)
	config.RamRoleSessionName = getProviderConfig("", "ram_session_name")
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
		if config.RamRoleSessionName == "" {
			config.RamRoleSessionName = "terraform"
		}
		config.RamRolePolicy = assumeRole["policy"].(string)
		if assumeRole["session_expiration"].(int) == 0 {
			if v := os.Getenv("APSARASTACK_ASSUME_ROLE_SESSION_EXPIRATION"); v != "" {
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
		return nil, err
	}
	ossServicedomain := d.Get("ossservice_domain").(string)
	if ossServicedomain != "" {
		config.OssServerEndpoint = ossServicedomain
	}
	domain := d.Get("domain").(string)
	if domain != "" {
		config.EcsEndpoint = domain
		config.VpcEndpoint = domain
		config.SlbEndpoint = domain
		config.OssEndpoint = domain
		config.AscmEndpoint = domain
		config.RdsEndpoint = domain
		config.OnsEndpoint = domain
		config.KmsEndpoint = domain
		config.LogEndpoint = domain
		config.CrEndpoint = domain
		config.EssEndpoint = domain
		config.DnsEndpoint = domain
		config.KVStoreEndpoint = domain
		config.GpdbEndpoint = domain
		config.DdsEndpoint = domain
		config.CsEndpoint = domain
		config.CmsEndpoint = domain
		config.HitsdbEndpoint = domain
		config.MaxComputeEndpoint = domain
		config.OtsEndpoint = domain
		config.DatahubEndpoint = domain
		config.EdasEndpoint = domain
		config.AdbEndpoint = domain
		config.RosEndpoint = domain
		config.DtsEndpoint = domain
		config.AlikafkaEndpoint = domain
		config.NasEndpoint = domain
		config.ApigatewayEndpoint = domain
		config.DmsEnterpriseEndpoint = domain
		config.HBaseEndpoint = domain
		config.DrdsEndpoint = domain
		config.QuickbiEndpoint = domain
		config.ElasticsearchEndpoint = domain
		config.DataworkspublicEndpoint = domain
		config.DbsEndpoint = domain
	} else {

		endpointsSet := d.Get("endpoints").(*schema.Set)

		for _, endpointsSetI := range endpointsSet.List() {
			endpoints := endpointsSetI.(map[string]interface{})
			config.EcsEndpoint = strings.TrimSpace(endpoints["ecs"].(string))
			config.VpcEndpoint = strings.TrimSpace(endpoints["vpc"].(string))
			config.AscmEndpoint = strings.TrimSpace(endpoints["ascm"].(string))
			config.RdsEndpoint = strings.TrimSpace(endpoints["rds"].(string))
			config.OssEndpoint = strings.TrimSpace(endpoints["oss"].(string))
			config.OnsEndpoint = strings.TrimSpace(endpoints["ons"].(string))
			config.KmsEndpoint = strings.TrimSpace(endpoints["kms"].(string))
			config.LogEndpoint = strings.TrimSpace(endpoints["log"].(string))
			config.SlbEndpoint = strings.TrimSpace(endpoints["slb"].(string))
			config.CrEndpoint = strings.TrimSpace(endpoints["cr"].(string))
			config.EssEndpoint = strings.TrimSpace(endpoints["ess"].(string))
			config.DnsEndpoint = strings.TrimSpace(endpoints["dns"].(string))
			config.KVStoreEndpoint = strings.TrimSpace(endpoints["kvstore"].(string))
			config.GpdbEndpoint = strings.TrimSpace(endpoints["gpdb"].(string))
			config.DdsEndpoint = strings.TrimSpace(endpoints["dds"].(string))
			config.CsEndpoint = strings.TrimSpace(endpoints["cs"].(string))
			config.CmsEndpoint = strings.TrimSpace(endpoints["cms"].(string))
			config.OtsEndpoint = strings.TrimSpace(endpoints["ots"].(string))
			config.DatahubEndpoint = strings.TrimSpace(endpoints["datahub"].(string))
			config.AdbEndpoint = strings.TrimSpace(endpoints["adb"].(string))
			config.StsEndpoint = strings.TrimSpace(endpoints["sts"].(string))
			config.RosEndpoint = strings.TrimSpace(endpoints["ros"].(string))
			config.DtsEndpoint = strings.TrimSpace(endpoints["dts"].(string))
			config.AlikafkaEndpoint = strings.TrimSpace(endpoints["alikafka"].(string))
			config.NasEndpoint = strings.TrimSpace(endpoints["nas"].(string))
			config.ApigatewayEndpoint = strings.TrimSpace(endpoints["apigateway"].(string))
			config.DmsEnterpriseEndpoint = strings.TrimSpace(endpoints["dms_enterprise"].(string))
			config.HBaseEndpoint = strings.TrimSpace(endpoints["hbase"].(string))
			config.DrdsEndpoint = strings.TrimSpace(endpoints["drds"].(string))
			config.QuickbiEndpoint = strings.TrimSpace(endpoints["quickbi"].(string))
			config.DataworkspublicEndpoint = strings.TrimSpace(endpoints["dataworkspublic"].(string))
			config.DbsEndpoint = strings.TrimSpace(endpoints["dbs"].(string))
		}
	}
	DbsEndpoint := d.Get("dbs_endpoint").(string)
	if DbsEndpoint != "" {
		config.DbsEndpoint = DbsEndpoint
	}
	DataworkspublicEndpoint := d.Get("dataworkspublic").(string)
	if DataworkspublicEndpoint != "" {
		config.DataworkspublicEndpoint = DataworkspublicEndpoint
	}
	QuickbiEndpoint := d.Get("quickbi_endpoint").(string)
	if QuickbiEndpoint != "" {
		config.QuickbiEndpoint = QuickbiEndpoint
	}
	kafkaOpenApidomain := d.Get("kafkaopenapi_domain").(string)
	if kafkaOpenApidomain != "" {
		config.AlikafkaOpenAPIEndpoint = kafkaOpenApidomain
	}
	StsEndpoint := d.Get("sts_endpoint").(string)
	if StsEndpoint != "" {
		config.StsEndpoint = StsEndpoint
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
		config.SLSOpenAPIEndpoint = slsOpenAPIEndpoint
	}
	if strings.ToLower(config.Protocol) == "https" {
		config.Protocol = "HTTPS"
	} else {
		config.Protocol = "HTTP"
	}

	config.ResourceSetName = d.Get("resource_group_set_name").(string)
	if config.Department == "" || config.ResourceGroup == "" {
		dept, rg, err := getResourceCredentials(config)
		if err != nil {
			return nil, err
		}
		config.Department = dept
		config.ResourceGroup = rg
	}

	if config.RamRoleArn != "" {
		config.AccessKey, config.SecretKey, config.SecurityToken, err = getAssumeRoleAK(config)
		if err != nil {
			return nil, err
		}
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
		"access_key": "The access key for API operations. You can retrieve this from the 'Security Management' section of the ApsaraStack console.",

		"secret_key": "The secret key for API operations. You can retrieve this from the 'Security Management' section of the ApsaraStackconsole.",

		"security_token": "security token. A security token is only required if you are using Security Token Service.",

		"insecure": "Use this to Trust self-signed certificates. It's typically used to allow insecure connections",

		"proxy": "Use this to set proxy connection",

		"domain": "Use this to override the default domain. It's typically used to connect to custom domain.",
	}
}
func endpointsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cbn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cbn_endpoint"],
				},

				"ecs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ecs_endpoint"],
				},
				"sts": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sts_endpoint"],
				},
				"ascm": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ascm_endpoint"],
				},
				"rds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["rds_endpoint"],
				},
				"slb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["slb_endpoint"],
				},
				"vpc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["vpc_endpoint"],
				},
				"cen": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cen_endpoint"],
				},
				"ess": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ess_endpoint"],
				},
				"oss": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["oss_endpoint"],
				},
				"ons": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ons_endpoint"],
				},
				"alikafka": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["alikafka_endpoint"],
				},
				"dns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dns_endpoint"],
				},
				"ram": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ram_endpoint"],
				},
				"cs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cs_endpoint"],
				},
				"cr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cr_endpoint"],
				},
				"cdn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cdn_endpoint"],
				},

				"kms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["kms_endpoint"],
				},

				"ots": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ots_endpoint"],
				},

				"cms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cms_endpoint"],
				},

				"pvtz": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["pvtz_endpoint"],
				},
				// log service is sls service
				"log": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["log_endpoint"],
				},
				"drds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["drds_endpoint"],
				},
				"dds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dds_endpoint"],
				},
				"polardb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["polardb_endpoint"],
				},
				"gpdb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["gpdb_endpoint"],
				},
				"kvstore": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["kvstore_endpoint"],
				},
				"fc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["fc_endpoint"],
				},
				"apigateway": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["apigateway_endpoint"],
				},
				"datahub": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["datahub_endpoint"],
				},
				"mns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mns_endpoint"],
				},
				"location": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["location_endpoint"],
				},
				"elasticsearch": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["elasticsearch_endpoint"],
				},
				"nas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["nas_endpoint"],
				},
				"actiontrail": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["actiontrail_endpoint"],
				},
				"cas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cas_endpoint"],
				},
				"bssopenapi": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["bssopenapi_endpoint"],
				},
				"ddoscoo": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ddoscoo_endpoint"],
				},
				"ddosbgp": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ddosbgp_endpoint"],
				},
				"emr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["emr_endpoint"],
				},
				"market": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["market_endpoint"],
				},
				"adb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["adb_endpoint"],
				},
				"maxcompute": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["maxcompute_endpoint"],
				},
				"dms_enterprise": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dms_enterprise_endpoint"],
				},
				"quickbi": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["quickbi_endpoint"],
				},
				"dataworkspublic": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_DATAWORKS_PUBLIC_ENDPOINT", nil),
					Description: descriptions["dataworkspublic_endpoint"],
				},
				"dbs_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_DBS_ENDPOINT", nil),
					Description: descriptions["dbs_endpoint"],
				},
			},
		},
		Set: endpointsToHash,
	}
}
func endpointsToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["ascm"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ecs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["rds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["slb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["vpc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cen"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ess"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["oss"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ons"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["alikafka"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ram"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cdn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["kms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ots"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["pvtz"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ascm"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["log"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["drds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["gpdb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["kvstore"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["polardb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["fc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["apigateway"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["datahub"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["mns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["location"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["elasticsearch"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["nas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["actiontrail"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["bssopenapi"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ddoscoo"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ddosbgp"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["emr"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["market"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["adb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cbn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["maxcompute"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dms_enterprise"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["quickbi"].(string)))

	return hashcode.String(buf.String())
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
			return nil, WrapError(err)
		}
		if profilePath == "" {
			profilePath = fmt.Sprintf("%s/.apsarastack/config.json", os.Getenv("HOME"))
			if runtime.GOOS == "windows" {
				profilePath = fmt.Sprintf("%s/.apsarastack/config.json", os.Getenv("USERPROFILE"))
			}
		}
		providerConfig = make(map[string]interface{})
		_, err = os.Stat(profilePath)
		if !os.IsNotExist(err) {
			data, err := ioutil.ReadFile(profilePath)
			if err != nil {
				return nil, WrapError(err)
			}
			config := map[string]interface{}{}
			err = json.Unmarshal(data, &config)
			if err != nil {
				return nil, WrapError(err)
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
					DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ASSUME_ROLE_ARN", os.Getenv("APSARASTACK_ASSUME_ROLE_ARN")),
				},
				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_session_name"],
					DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ASSUME_ROLE_SESSION_NAME", ""),
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
	request := sts.CreateAssumeRoleRequest()
	request.RoleArn = config.RamRoleArn
	if config.RamRoleSessionName == "" {
		config.RamRoleSessionName = "terraform"
	}
	request.RoleSessionName = config.RamRoleSessionName
	//request.DurationSeconds = requests.NewInteger(config.RamRoleSessionExpiration)
	request.Policy = config.RamRolePolicy
	request.Scheme = "https"
	request.SetHTTPSInsecure(true)
	request.Domain = config.StsEndpoint
	request.Headers["x-ascm-product-name"] = "sts"
	request.Headers["x-acs-organizationId"] = config.Department

	var client *sts.Client
	var err error
	if config.SecurityToken == "" {
		client, err = sts.NewClientWithAccessKey(config.RegionId, config.AccessKey, config.SecretKey)
	} else {
		client, err = sts.NewClientWithStsToken(config.RegionId, config.AccessKey, config.SecretKey, config.SecurityToken)
	}

	if err != nil {
		return "", "", "", err
	}

	client.Domain = config.StsEndpoint
	client.AppendUserAgent(connectivity.Terraform, connectivity.TerraformVersion)
	client.AppendUserAgent(connectivity.Provider, connectivity.ProviderVersion)
	client.AppendUserAgent(connectivity.Module, config.ConfigurationSource)
	client.SetHTTPSInsecure(config.Insecure)
	if config.Proxy != "" {
		client.SetHttpProxy(config.Proxy)
	}
	response, err := client.AssumeRole(request)
	if err != nil {
		return config.AccessKey, config.SecretKey, config.SecurityToken, err
	}

	return response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, response.Credentials.SecurityToken, nil
}

func getResourceCredentials(config *connectivity.Config) (string, string, error) {
	endpoint := config.AscmEndpoint
	if endpoint == "" {
		return "", "", fmt.Errorf("unable to initialize the ascm client: endpoint or domain is not provided for ascm service")
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(config.RegionId, string(connectivity.ASCMCode), endpoint)
	}
	ascmClient, err := sdk.NewClientWithAccessKey(config.RegionId, config.AccessKey, config.SecretKey)
	if err != nil {
		return "", "", fmt.Errorf("unable to initialize the ascm client: %#v", err)
	}

	ascmClient.AppendUserAgent(connectivity.Terraform, connectivity.TerraformVersion)
	ascmClient.AppendUserAgent(connectivity.Provider, connectivity.ProviderVersion)
	ascmClient.AppendUserAgent(connectivity.Module, config.ConfigurationSource)
	ascmClient.SetHTTPSInsecure(config.Insecure)
	ascmClient.Domain = endpoint
	if config.Proxy != "" {
		ascmClient.SetHttpProxy(config.Proxy)
	}
	if config.ResourceSetName == "" {
		return "", "", fmt.Errorf("errror while fetching resource group details, resource group set name can not be empty")
	}
	request := requests.NewCommonRequest()
	if config.Insecure {
		request.SetHTTPSInsecure(config.Insecure)
	}
	request.RegionId = config.RegionId
	request.Method = "GET"         // Set request method
	request.Product = "ascm"       // Specify product
	request.Domain = endpoint      // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2019-05-10" // Specify product version
	// Set request scheme. Default: http
	if strings.ToLower(config.Protocol) == "https" {
		log.Printf("PROTOCOL SET TO HTTPS")
		request.Scheme = "https"
	} else {
		log.Printf("PROTOCOL SET TO HTTP")
		request.Scheme = "http"
	}
	if config.Insecure {
		ascmClient.SetHTTPSInsecure(config.Insecure)
	}
	request.ApiName = "ListResourceGroup"
	request.QueryParams = map[string]string{
		"AccessKeySecret":   config.SecretKey,
		"Product":           "ascm",
		"Department":        config.Department,
		"ResourceGroup":     config.ResourceGroup,
		"RegionId":          config.RegionId,
		"Action":            "ListResourceGroup",
		"Version":           "2019-05-10",
		"SignatureVersion":  "1.0",
		"resourceGroupName": config.ResourceSetName,
	}
	resp := responses.BaseResponse{}
	if config.Insecure {
		request.SetHTTPSInsecure(config.Insecure)
	}
	request.TransToAcsRequest()
	err = ascmClient.DoAction(request, &resp)
	if err != nil {
		return "", "", err
	}
	response := &ResourceGroup{}
	err = json.Unmarshal(resp.GetHttpContentBytes(), response)
	if err != nil {
		return "", "", err
	}
	var deptId int   // Organization ID
	var resGrpId int //ID of resource set
	deptId = 0
	resGrpId = 0
	if len(response.Data) == 0 || response.Code != "200" {
		if len(response.Data) == 0 {
			return "", "", fmt.Errorf("resource group ID and organization not found for resource set %s", config.ResourceSetName)
		}
		return "", "", fmt.Errorf("unable to initialize the ascm client: department or resource_group is not provided")
	} else {
		for _, j := range response.Data {
			if j.ResourceGroupName == config.ResourceSetName {
				deptId = j.OrganizationID
				resGrpId = j.ID
				break
			}
		}
	}

	//log.Printf("[INFO] Get Resource Group Details Succssfull for Resource set: %s : Department: %s, ResourceGroupId: %s", config.ResourceSetName, fmt.Sprint(response.Data[0].OrganizationID), fmt.Sprint(response.Data[0].ID))
	log.Printf("[INFO] Get Resource Group Details Succssfull for Resource set: %s : Department: %s, ResourceGroupId: %s", config.ResourceSetName, fmt.Sprint(deptId), fmt.Sprint(resGrpId))
	//return fmt.Sprint(response.Data[0].OrganizationID), fmt.Sprint(response.Data[0].ID), err
	return fmt.Sprint(deptId), fmt.Sprint(resGrpId), err

}

func wiatSecondsIfWithTest(second int) {
	// 测试模式下休眠一秒，防止数据缓存导致二次plan失败
	if os.Getenv("TF_ACC") == "1" {
		time.Sleep(time.Duration(second) * time.Second)
	}
}
