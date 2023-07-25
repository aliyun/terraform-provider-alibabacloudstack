package alibabacloudstack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
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
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_DOMAIN", nil),
				Description: descriptions["domain"],
			},
			"ossservice_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_OSSSERVICE_DOMAIN", nil),
				Description: descriptions["ossservice_domain"],
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
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ORGANIZATION_ACCESSKEY", nil),
				Description: descriptions["organization_accesskey"],
			},
			"organization_secretkey": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_ORGANIZATION_SECRETKEY", nil),
				Description: descriptions["organization_secretkey"],
			},
			"sls_openapi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_SLS_OPENAPI_ENDPOINT", nil),
				Description: descriptions["sls_openapi_endpoint"],
			},
			"sts_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_STS_ENDPOINT", os.Getenv("ALIBABACLOUDSTACK_STS_ENDPOINT")),
				Description: descriptions["sts_endpoint"],
			},
			"quickbi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_QUICKBI_ENDPOINT", nil),
				Description: descriptions["quickbi_endpoint"],
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
			},
			"dbs_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_DBS_ENDPOINT", nil),
				Description: descriptions["dbs_endpoint"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alibabacloudstack_account":                                dataSourceAlibabacloudStackAccount(),
			"alibabacloudstack_adb_clusters":                           dataSourceAlibabacloudStackAdbDbClusters(),
			"alibabacloudstack_adb_zones":                              dataSourceAlibabacloudStackAdbZones(),
			"alibabacloudstack_adb_db_clusters":                        dataSourceAlibabacloudStackAdbDbClusters(),
			"alibabacloudstack_api_gateway_apis":                       dataSourceAlibabacloudStackApiGatewayApis(),
			"alibabacloudstack_api_gateway_apps":                       dataSourceAlibabacloudStackApiGatewayApps(),
			"alibabacloudstack_api_gateway_groups":                     dataSourceAlibabacloudStackApiGatewayGroups(),
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
			"alibabacloudstack_cr_ee_instances":                        dataSourceAlibabacloudStackCrEEInstances(),
			"alibabacloudstack_cr_ee_namespaces":                       dataSourceAlibabacloudStackCrEENamespaces(),
			"alibabacloudstack_cr_ee_repos":                            dataSourceAlibabacloudStackCrEERepos(),
			"alibabacloudstack_cr_ee_sync_rules":                       dataSourceAlibabacloudStackCrEESyncRules(),
			"alibabacloudstack_cr_namespaces":                          dataSourceAlibabacloudStackCRNamespaces(),
			"alibabacloudstack_cr_repos":                               dataSourceAlibabacloudStackCRRepos(),
			"alibabacloudstack_cs_kubernetes_clusters":                 dataSourceAlibabacloudStackCSKubernetesClusters(),
			"alibabacloudstack_cms_alarm_contacts":                     dataSourceAlibabacloudstackCmsAlarmContacts(),
			"alibabacloudstack_cms_alarm_contact_groups":               dataSourceAlibabacloudstackCmsAlarmContactGroups(),
			"alibabacloudstack_cms_project_meta":                       dataSourceAlibabacloudstackCmsProjectMeta(),
			"alibabacloudstack_cms_metric_metalist":                    dataSourceAlibabacloudstackCmsMetricMetalist(),
			"alibabacloudstack_cms_alarms":                             dataSourceAlibabacloudstackCmsAlarms(),
			"alibabacloudstack_datahub_service":                        dataSourceAlibabacloudStackDatahubService(),
			"alibabacloudstack_db_instances":                           dataSourceAlibabacloudStackDBInstances(),
			"alibabacloudstack_db_zones":                               dataSourceAlibabacloudStackDBZones(),
			"alibabacloudstack_disks":                                  dataSourceAlibabacloudStackDisks(),
			"alibabacloudstack_dns_records":                            dataSourceAlibabacloudStackDnsRecords(),
			"alibabacloudstack_dns_groups":                             dataSourceAlibabacloudStackDnsGroups(),
			"alibabacloudstack_dns_domains":                            dataSourceAlibabacloudStackDnsDomains(),
			"alibabacloudstack_drds_instances":                         dataSourceAlibabacloudStackDRDSInstances(),
			"alibabacloudstack_dms_enterprise_instances":               dataSourceAlibabacloudStackDmsEnterpriseInstances(),
			"alibabacloudstack_dms_enterprise_users":                   dataSourceAlibabacloudStackDmsEnterpriseUsers(),
			"alibabacloudstack_ecs_commands":                           dataSourceAlibabacloudStackEcsCommands(),
			"alibabacloudstack_ecs_deployment_sets":                    dataSourceAlibabacloudStackEcsDeploymentSets(),
			"alibabacloudstack_ecs_hpc_clusters":                       dataSourceAlibabacloudStackEcsHpcClusters(),
			"alibabacloudstack_ecs_dedicated_hosts":                    dataSourceAlibabacloudStackEcsDedicatedHosts(),
			"alibabacloudstack_edas_deploy_groups":                     dataSourceAlibabacloudStackEdasDeployGroups(),
			"alibabacloudstack_edas_clusters":                          dataSourceAlibabacloudStackEdasClusters(),
			"alibabacloudstack_edas_applications":                      dataSourceAlibabacloudStackEdasApplications(),
			"alibabacloudstack_eips":                                   dataSourceAlibabacloudStackEips(),
			"alibabacloudstack_ess_scaling_configurations":             dataSourceAlibabacloudStackEssScalingConfigurations(),
			"alibabacloudstack_ess_scaling_groups":                     dataSourceAlibabacloudStackEssScalingGroups(),
			"alibabacloudstack_ess_lifecycle_hooks":                    dataSourceAlibabacloudStackEssLifecycleHooks(),
			"alibabacloudstack_ess_notifications":                      dataSourceAlibabacloudStackEssNotifications(),
			"alibabacloudstack_ess_scaling_rules":                      dataSourceAlibabacloudStackEssScalingRules(),
			"alibabacloudstack_ess_scheduled_tasks":                    dataSourceAlibabacloudStackEssScheduledTasks(),
			"alibabacloudstack_forward_entries":                        dataSourceAlibabacloudStackForwardEntries(),
			"alibabacloudstack_gpdb_accounts":                          dataSourceAlibabacloudStackGpdbAccounts(),
			"alibabacloudstack_gpdb_instances":                         dataSourceAlibabacloudStackGpdbInstances(),
			"alibabacloudstack_hbase_instances":                        dataSourceAlibabacloudStackHBaseInstances(),
			"alibabacloudstack_instances":                              dataSourceAlibabacloudStackInstances(),
			"alibabacloudstack_instance_type_families":                 dataSourceAlibabacloudStackInstanceTypeFamilies(),
			"alibabacloudstack_instance_types":                         dataSourceAlibabacloudStackInstanceTypes(),
			"alibabacloudstack_images":                                 dataSourceAlibabacloudStackImages(),
			"alibabacloudstack_key_pairs":                              dataSourceAlibabacloudStackKeyPairs(),
			"alibabacloudstack_kms_aliases":                            dataSourceAlibabacloudStackKmsAliases(),
			"alibabacloudstack_kms_ciphertext":                         dataSourceAlibabacloudStackKmsCiphertext(),
			"alibabacloudstack_kms_keys":                               dataSourceAlibabacloudStackKmsKeys(),
			"alibabacloudstack_kms_secrets":                            dataSourceAlibabacloudStackKmsSecrets(),
			"alibabacloudstack_kvstore_instances":                      dataSourceAlibabacloudStackKVStoreInstances(),
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
			"alibabacloudstack_nas_mount_targets":                      dataSourceAlibabacloudStackNasMountTargets(),
			"alibabacloudstack_nas_access_rules":                       dataSourceAlibabacloudStackAccessRules(),
			"alibabacloudstack_nat_gateways":                           dataSourceAlibabacloudStackNatGateways(),
			"alibabacloudstack_network_acls":                           dataSourceAlibabacloudStackNetworkAcls(),
			"alibabacloudstack_network_interfaces":                     dataSourceAlibabacloudStackNetworkInterfaces(),
			"alibabacloudstack_oss_buckets":                            dataSourceAlibabacloudStackOssBuckets(),
			"alibabacloudstack_oss_bucket_objects":                     dataSourceAlibabacloudStackOssBucketObjects(),
			"alibabacloudstack_ons_instances":                          dataSourceAlibabacloudStackOnsInstances(),
			"alibabacloudstack_ons_topics":                             dataSourceAlibabacloudStackOnsTopics(),
			"alibabacloudstack_ons_groups":                             dataSourceAlibabacloudStackOnsGroups(),
			"alibabacloudstack_ots_tables":                             dataSourceAlibabacloudStackOtsTables(),
			"alibabacloudstack_ots_instances":                          dataSourceAlibabacloudStackOtsInstances(),
			"alibabacloudstack_ots_instances_attachment":               dataSourceAlibabacloudStackOtsInstanceAttachments(),
			"alibabacloudstack_ots_service":                            dataSourceAlibabacloudStackOtsService(),
			"alibabacloudstack_quick_bi_users":                         dataSourceAlibabacloudStackQuickBiUsers(),
			"alibabacloudstack_router_interfaces":                      dataSourceAlibabacloudStackRouterInterfaces(),
			"alibabacloudstack_ram_service_role_products":              dataSourceAlibabacloudstackRamServiceRoleProducts(),
			"alibabacloudstack_route_tables":                           dataSourceAlibabacloudStackRouteTables(),
			"alibabacloudstack_route_entries":                          dataSourceAlibabacloudStackRouteEntries(),
			"alibabacloudstack_ros_stacks":                             dataSourceAlibabacloudStackRosStacks(),
			"alibabacloudstack_ros_templates":                          dataSourceAlibabacloudStackRosTemplates(),
			"alibabacloudstack_security_groups":                        dataSourceAlibabacloudStackSecurityGroups(),
			"alibabacloudstack_security_group_rules":                   dataSourceAlibabacloudStackSecurityGroupRules(),
			"alibabacloudstack_snapshots":                              dataSourceAlibabacloudStackSnapshots(),
			"alibabacloudstack_slb_listeners":                          dataSourceAlibabacloudStackSlbListeners(),
			"alibabacloudstack_slb_server_groups":                      dataSourceAlibabacloudStackSlbServerGroups(),
			"alibabacloudstack_slb_acls":                               dataSourceAlibabacloudStackSlbAcls(),
			"alibabacloudstack_slb_domain_extensions":                  dataSourceAlibabacloudStackSlbDomainExtensions(),
			"alibabacloudstack_slb_rules":                              dataSourceAlibabacloudStackSlbRules(),
			"alibabacloudstack_slb_master_slave_server_groups":         dataSourceAlibabacloudStackSlbMasterSlaveServerGroups(),
			"alibabacloudstack_slbs":                                   dataSourceAlibabacloudStackSlbs(),
			"alibabacloudstack_slb_zones":                              dataSourceAlibabacloudStackSlbZones(),
			"alibabacloudstack_snat_entries":                           dataSourceAlibabacloudStackSnatEntries(),
			"alibabacloudstack_slb_server_certificates":                dataSourceAlibabacloudStackSlbServerCertificates(),
			"alibabacloudstack_slb_ca_certificates":                    dataSourceAlibabacloudStackSlbCACertificates(),
			"alibabacloudstack_slb_backend_servers":                    dataSourceAlibabacloudStackSlbBackendServers(),
			"alibabacloudstack_tsdb_zones":                             dataSourceAlibabacloudStackTsdbZones(),
			"alibabacloudstack_vpn_gateways":                           dataSourceAlibabacloudStackVpnGateways(),
			"alibabacloudstack_vpn_customer_gateways":                  dataSourceAlibabacloudStackVpnCustomerGateways(),
			"alibabacloudstack_vpn_connections":                        dataSourceAlibabacloudStackVpnConnections(),
			"alibabacloudstack_vpc_ipv6_gateways":                      dataSourceAlibabacloudStackVpcIpv6Gateways(),
			"alibabacloudstack_vpc_ipv6_egress_rules":                  dataSourceAlibabacloudStackVpcIpv6EgressRules(),
			"alibabacloudstack_vpc_ipv6_addresses":                     dataSourceAlibabacloudStackVpcIpv6Addresses(),
			"alibabacloudstack_vpc_ipv6_internet_bandwidths":           dataSourceAlibabacloudStackVpcIpv6InternetBandwidths(),
			"alibabacloudstack_vswitches":                              dataSourceAlibabacloudStackVSwitches(),
			"alibabacloudstack_vpcs":                                   dataSourceAlibabacloudStackVpcs(),
			"alibabacloudstack_zones":                                  dataSourceAlibabacloudStackZones(),
			"alibabacloudstack_elasticsearch_instances":                dataSourceAlibabacloudStackElasticsearch(),
			"alibabacloudstack_elasticsearch_zones":                    dataSourceAlibabacloudStackElaticsearchZones(),
			"alibabacloudstack_ehpc_job_templates":                     dataSourceAlibabacloudStackEhpcJobTemplates(),
			"alibabacloudstack_oos_executions":                         dataSourceAlibabacloudStackOosExecutions(),
			"alibabacloudstack_oos_templates":                          dataSourceAlibabacloudStackOosTemplates(),
			"alibabacloudstack_express_connect_physical_connections":   dataSourceAlibabacloudStackExpressConnectPhysicalConnections(),
			"alibabacloudstack_express_connect_access_points":          dataSourceAlibabacloudStackExpressConnectAccessPoints(),
			"alibabacloudstack_express_connect_virtual_border_routers": dataSourceAlibabacloudStackExpressConnectVirtualBorderRouters(),
			"alibabacloudStack_cloud_firewall_control_policies":        dataSourceAlibabacloudStackCloudFirewallControlPolicies(),
			"alibabacloudstack_ecs_ebs_storage_sets":                   dataSourceAlibabacloudStackEcsEbsStorageSets(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"alibabacloudstack_ess_scaling_configuration":             resourceAlibabacloudStackEssScalingConfiguration(),
			"alibabacloudstack_adb_account":                           resourceAlibabacloudStackAdbAccount(),
			"alibabacloudstack_adb_backup_policy":                     resourceAlibabacloudStackAdbBackupPolicy(),
			"alibabacloudstack_adb_cluster":                           resourceAlibabacloudStackAdbDbCluster(),
			"alibabacloudstack_adb_connection":                        resourceAlibabacloudStackAdbConnection(),
			"alibabacloudstack_adb_db_cluster":                        resourceAlibabacloudStackAdbDbCluster(),
			"alibabacloudstack_alikafka_sasl_acl":                     resourceAlibabacloudStackAlikafkaSaslAcl(),
			"alibabacloudstack_alikafka_sasl_user":                    resourceAlibabacloudStackAlikafkaSaslUser(),
			"alibabacloudstack_alikafka_topic":                        resourceAlibabacloudStackAlikafkaTopic(),
			"alibabacloudstack_api_gateway_api":                       resourceAlibabacloudStackApigatewayApi(),
			"alibabacloudstack_api_gateway_app":                       resourceAlibabacloudStackApigatewayApp(),
			"alibabacloudstack_api_gateway_app_attachment":            resourceAliyunApigatewayAppAttachment(),
			"alibabacloudstack_api_gateway_group":                     resourceAlibabacloudStackApigatewayGroup(),
			"alibabacloudstack_api_gateway_vpc_access":                resourceAlibabacloudStackApigatewayVpc(),
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
			"alibabacloudstack_cms_alarm_contact":                     resourceAlibabacloudstackCmsAlarmContact(),
			"alibabacloudstack_cms_alarm_contact_group":               resourceAlibabacloudstackCmsAlarmContactGroup(),
			"alibabacloudstack_cms_site_monitor":                      resourceAlibabacloudStackCmsSiteMonitor(),
			"alibabacloudstack_common_bandwidth_package":              resourceAlibabacloudStackCommonBandwidthPackage(),
			"alibabacloudstack_common_bandwidth_package_attachment":   resourceAlibabacloudStackCommonBandwidthPackageAttachment(),
			"alibabacloudstack_cr_ee_namespace":                       resourceAlibabacloudStackCrEENamespace(),
			"alibabacloudstack_cr_ee_repo":                            resourceAlibabacloudStackCrEERepo(),
			"alibabacloudstack_cr_ee_sync_rule":                       resourceAlibabacloudStackCrEESyncRule(),
			"alibabacloudstack_cr_namespace":                          resourceAlibabacloudStackCRNamespace(),
			"alibabacloudstack_cr_repo":                               resourceAlibabacloudStackCRRepo(),
			"alibabacloudstack_cs_kubernetes":                         resourceAlibabacloudStackCSKubernetes(),
			"alibabacloudstack_cs_kubernetes_node_pool":               resourceAlibabacloudStackCSKubernetesNodePool(),
			"alibabacloudstack_datahub_project":                       resourceAlibabacloudStackDatahubProject(),
			"alibabacloudstack_datahub_subscription":                  resourceAlibabacloudStackDatahubSubscription(),
			"alibabacloudstack_datahub_topic":                         resourceAlibabacloudStackDatahubTopic(),
			"alibabacloudstack_db_account":                            resourceAlibabacloudStackDBAccount(),
			"alibabacloudstack_db_account_privilege":                  resourceAlibabacloudStackDBAccountPrivilege(),
			"alibabacloudstack_db_backup_policy":                      resourceAlibabacloudStackDBBackupPolicy(),
			"alibabacloudstack_db_connection":                         resourceAlibabacloudStackDBConnection(),
			"alibabacloudstack_db_database":                           resourceAlibabacloudStackDBDatabase(),
			"alibabacloudstack_db_instance":                           resourceAlibabacloudStackDBInstance(),
			"alibabacloudstack_db_read_write_splitting_connection":    resourceAlibabacloudStackDBReadWriteSplittingConnection(),
			"alibabacloudstack_db_readonly_instance":                  resourceAlibabacloudStackDBReadonlyInstance(),
			"alibabacloudstack_disk":                                  resourceAlibabacloudStackDisk(),
			"alibabacloudstack_disk_attachment":                       resourceAlibabacloudStackDiskAttachment(),
			"alibabacloudstack_dms_enterprise_instance":               resourceAlibabacloudStackDmsEnterpriseInstance(),
			"alibabacloudstack_dms_enterprise_user":                   resourceAlibabacloudStackDmsEnterpriseUser(),
			"alibabacloudstack_dns_domain":                            resourceAlibabacloudStackDnsDomain(),
			"alibabacloudstack_dns_domain_attachment":                 resourceAlibabacloudStackDnsDomainAttachment(),
			"alibabacloudstack_dns_group":                             resourceAlibabacloudStackDnsGroup(),
			"alibabacloudstack_dns_record":                            resourceAlibabacloudStackDnsRecord(),
			"alibabacloudstack_drds_instance":                         resourceAlibabacloudStackDRDSInstance(),
			"alibabacloudstack_dts_subscription_job":                  resourceAlibabacloudStackDtsSubscriptionJob(),
			"alibabacloudstack_dts_synchronization_instance":          resourceAlibabacloudStackDtsSynchronizationInstance(),
			"alibabacloudstack_dts_synchronization_job":               resourceAlibabacloudStackDtsSynchronizationJob(),
			"alibabacloudstack_ecs_command":                           resourceAlibabacloudStackEcsCommand(),
			"alibabacloudstack_ecs_dedicated_host":                    resourceAlibabacloudStackEcsDedicatedHost(),
			"alibabacloudstack_ecs_deployment_set":                    resourceAlibabacloudStackEcsDeploymentSet(),
			"alibabacloudstack_ecs_hpc_cluster":                       resourceAlibabacloudStackEcsHpcCluster(),
			"alibabacloudstack_ecs_ebs_storage_set":                   resourceAlibabacloudStackEcsEbsStorageSets(),
			"alibabacloudstack_edas_application":                      resourceAlibabacloudStackEdasApplication(),
			"alibabacloudstack_edas_application_scale":                resourceAlibabacloudStackEdasInstanceApplicationAttachment(),
			"alibabacloudstack_edas_cluster":                          resourceAlibabacloudStackEdasCluster(),
			"alibabacloudstack_edas_deploy_group":                     resourceAlibabacloudStackEdasDeployGroup(),
			"alibabacloudstack_edas_instance_cluster_attachment":      resourceAlibabacloudStackEdasInstanceClusterAttachment(),
			"alibabacloudstack_edas_k8s_application":                  resourceAlibabacloudStackEdasK8sApplication(),
			"alibabacloudstack_edas_k8s_cluster":                      resourceAlibabacloudStackEdasK8sCluster(),
			"alibabacloudstack_edas_slb_attachment":                   resourceAlibabacloudStackEdasSlbAttachment(),
			"alibabacloudstack_ehpc_job_template":                     resourceAlibabacloudStackEhpcJobTemplate(),
			"alibabacloudstack_eip":                                   resourceAlibabacloudStackEip(),
			"alibabacloudstack_eip_association":                       resourceAlibabacloudStackEipAssociation(),
			"alibabacloudstack_ess_alarm":                             resourceAlibabacloudStackEssAlarm(),
			"alibabacloudstack_ess_attachment":                        resourceAlibabacloudstackEssAttachment(),
			"alibabacloudstack_ess_lifecycle_hook":                    resourceAlibabacloudStackEssLifecycleHook(),
			"alibabacloudstack_ess_notification":                      resourceAlibabacloudStackEssNotification(),
			"alibabacloudstack_ess_scaling_group":                     resourceAlibabacloudStackEssScalingGroup(),
			"alibabacloudstack_ess_scaling_rule":                      resourceAlibabacloudStackEssScalingRule(),
			"alibabacloudstack_ess_scalinggroup_vserver_groups":       resourceAlibabacloudStackEssScalingGroupVserverGroups(),
			"alibabacloudstack_ess_scheduled_task":                    resourceAlibabacloudStackEssScheduledTask(),
			"alibabacloudstack_forward_entry":                         resourceAlibabacloudStackForwardEntry(),
			"alibabacloudstack_gpdb_account":                          resourceAlibabacloudStackGpdbAccount(),
			"alibabacloudstack_gpdb_connection":                       resourceAlibabacloudStackGpdbConnection(),
			"alibabacloudstack_gpdb_instance":                         resourceAlibabacloudStackGpdbInstance(),
			"alibabacloudstack_hbase_instance":                        resourceAlibabacloudStackHBaseInstance(),
			"alibabacloudstack_image":                                 resourceAlibabacloudStackImage(),
			"alibabacloudstack_image_copy":                            resourceAlibabacloudStackImageCopy(),
			"alibabacloudstack_image_export":                          resourceAlibabacloudStackImageExport(),
			"alibabacloudstack_image_import":                          resourceAlibabacloudStackImageImport(),
			"alibabacloudstack_image_share_permission":                resourceAlibabacloudStackImageSharePermission(),
			"alibabacloudstack_instance":                              resourceAlibabacloudStackInstance(),
			"alibabacloudstack_key_pair":                              resourceAlibabacloudStackKeyPair(),
			"alibabacloudstack_key_pair_attachment":                   resourceAlibabacloudStackKeyPairAttachment(),
			"alibabacloudstack_kms_alias":                             resourceAlibabacloudStackKmsAlias(),
			"alibabacloudstack_kms_ciphertext":                        resourceAlibabacloudStackKmsCiphertext(),
			"alibabacloudstack_kms_key":                               resourceAlibabacloudStackKmsKey(),
			"alibabacloudstack_kms_secret":                            resourceAlibabacloudStackKmsSecret(),
			"alibabacloudstack_kvstore_account":                       resourceAlibabacloudStackKVstoreAccount(),
			"alibabacloudstack_kvstore_backup_policy":                 resourceAlibabacloudStackKVStoreBackupPolicy(),
			"alibabacloudstack_kvstore_connection":                    resourceAlibabacloudStackKvstoreConnection(),
			"alibabacloudstack_kvstore_instance":                      resourceAlibabacloudStackKVStoreInstance(),
			"alibabacloudstack_launch_template":                       resourceAlibabacloudStackLaunchTemplate(),
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
			"alibabacloudstack_nas_access_group":                      resourceAlibabacloudStackNasAccessGroup(),
			"alibabacloudstack_nas_access_rule":                       resourceAlibabacloudStackNasAccessRule(),
			"alibabacloudstack_nas_file_system":                       resourceAlibabacloudStackNasFileSystem(),
			"alibabacloudstack_nas_mount_target":                      resourceAlibabacloudStackNasMountTarget(),
			"alibabacloudstack_nat_gateway":                           resourceAlibabacloudStackNatGateway(),
			"alibabacloudstack_network_acl":                           resourceAlibabacloudStackNetworkAcl(),
			"alibabacloudstack_network_acl_attachment":                resourceAlibabacloudStackNetworkAclAttachment(),
			"alibabacloudstack_network_acl_entries":                   resourceAlibabacloudStackNetworkAclEntries(),
			"alibabacloudstack_network_interface":                     resourceAlibabacloudStackNetworkInterface(),
			"alibabacloudstack_network_interface_attachment":          resourceNetworkInterfaceAttachment(),
			"alibabacloudstack_ons_group":                             resourceAlibabacloudStackOnsGroup(),
			"alibabacloudstack_ons_instance":                          resourceAlibabacloudStackOnsInstance(),
			"alibabacloudstack_ons_topic":                             resourceAlibabacloudStackOnsTopic(),
			"alibabacloudstack_oss_bucket":                            resourceAlibabacloudStackOssBucket(),
			"alibabacloudstack_oss_bucket_quota":                      resourceAlibabacloudStackOssBucketQuota(),
			"alibabacloudstack_oss_bucket_kms":                        resourceAlibabacloudStackOssBucketKms(),
			"alibabacloudstack_oss_bucket_object":                     resourceAlibabacloudStackOssBucketObject(),
			"alibabacloudstack_ots_instance":                          resourceAlibabacloudStackOtsInstance(),
			"alibabacloudstack_ots_instance_attachment":               resourceAlibabacloudStackOtsInstanceAttachment(),
			"alibabacloudstack_ots_table":                             resourceAlibabacloudStackOtsTable(),
			"alibabacloudstack_quick_bi_user":                         resourceAlibabacloudStackQuickBiUser(),
			"alibabacloudstack_quick_bi_user_group":                   resourceAlibabacloudStackQuickBiUserGroup(),
			"alibabacloudstack_quick_bi_workspace":                    resourceAlibabacloudStackQuickBiWorkspace(),
			"alibabacloudstack_ram_role_attachment":                   resourceAlibabacloudStackRamRoleAttachment(),
			"alibabacloudstack_reserved_instance":                     resourceAlibabacloudStackReservedInstance(),
			"alibabacloudstack_ros_stack":                             resourceAlibabacloudStackRosStack(),
			"alibabacloudstack_ros_template":                          resourceAlibabacloudStackRosTemplate(),
			"alibabacloudstack_route_entry":                           resourceAlibabacloudStackRouteEntry(),
			"alibabacloudstack_route_table":                           resourceAlibabacloudStackRouteTable(),
			"alibabacloudstack_route_table_attachment":                resourceAlibabacloudStackRouteTableAttachment(),
			"alibabacloudstack_router_interface":                      resourceAlibabacloudStackRouterInterface(),
			"alibabacloudstack_router_interface_connection":           resourceAlibabacloudStackRouterInterfaceConnection(),
			"alibabacloudstack_security_group":                        resourceAlibabacloudStackSecurityGroup(),
			"alibabacloudstack_security_group_rule":                   resourceAlibabacloudStackSecurityGroupRule(),
			"alibabacloudstack_slb":                                   resourceAlibabacloudStackSlb(),
			"alibabacloudstack_slb_acl":                               resourceAlibabacloudStackSlbAcl(),
			"alibabacloudstack_slb_backend_server":                    resourceAlibabacloudStackSlbBackendServer(),
			"alibabacloudstack_slb_ca_certificate":                    resourceAlibabacloudStackSlbCACertificate(),
			"alibabacloudstack_slb_domain_extension":                  resourceAlibabacloudStackSlbDomainExtension(),
			"alibabacloudstack_slb_listener":                          resourceAlibabacloudStackSlbListener(),
			"alibabacloudstack_slb_master_slave_server_group":         resourceAlibabacloudStackSlbMasterSlaveServerGroup(),
			"alibabacloudstack_slb_rule":                              resourceAlibabacloudStackSlbRule(),
			"alibabacloudstack_slb_server_certificate":                resourceAlibabacloudStackSlbServerCertificate(),
			"alibabacloudstack_slb_server_group":                      resourceAlibabacloudStackSlbServerGroup(),
			"alibabacloudstack_snapshot":                              resourceAlibabacloudStackSnapshot(),
			"alibabacloudstack_snapshot_policy":                       resourceAlibabacloudStackSnapshotPolicy(),
			"alibabacloudstack_snat_entry":                            resourceAlibabacloudStackSnatEntry(),
			"alibabacloudstack_vpc":                                   resourceAlibabacloudStackVpc(),
			"alibabacloudstack_vpc_ipv6_egress_rule":                  resourceAlibabacloudStackVpcIpv6EgressRule(),
			"alibabacloudstack_vpc_ipv6_gateway":                      resourceAlibabacloudStackVpcIpv6Gateway(),
			"alibabacloudstack_vpc_ipv6_internet_bandwidth":           resourceAlibabacloudStackVpcIpv6InternetBandwidth(),
			"alibabacloudstack_vpn_connection":                        resourceAlibabacloudStackVpnConnection(),
			"alibabacloudstack_vpn_customer_gateway":                  resourceAlibabacloudStackVpnCustomerGateway(),
			"alibabacloudstack_vpn_gateway":                           resourceAlibabacloudStackVpnGateway(),
			"alibabacloudstack_vpn_route_entry":                       resourceAlibabacloudStackVpnRouteEntry(),
			"alibabacloudstack_vswitch":                               resourceAlibabacloudStackSwitch(),
			"alibabacloudstack_data_works_folder":                     resourceAlibabacloudStackDataWorksFolder(),
			"alibabacloudstack_data_works_connection":                 resourceAlibabacloudStackDataWorksConnection(),
			"alibabacloudstack_data_works_user":                       resourceAlibabacloudStackDataWorksUser(),
			"alibabacloudstack_data_works_project":                    resourceAlibabacloudStackDataWorksProject(),
			"alibabacloudstack_data_works_user_role_binding":          resourceAlibabacloudStackDataWorksUserRoleBinding(),
			"alibabacloudstack_data_works_remind":                     resourceAlibabacloudStackDataWorksRemind(),
			"alibabacloudstack_elasticsearch_instance":                resourceAlibabacloudStackElasticsearch(),
			"alibabacloudstack_dbs_backup_plan":                       resourceAlibabacloudStackDbsBackupPlan(),
			"alibabacloudstack_express_connect_physical_connection":   resourceAlibabacloudStackExpressConnectPhysicalConnection(),
			"alibabacloudstack_express_connect_virtual_border_router": resourceAlibabacloudStackExpressConnectVirtualBorderRouter(),
			"alibabacloudstack_oos_template":                          resourceAlibabacloudStackOosTemplate(),
			"alibabacloudstack_oos_execution":                         resourceAlibabacloudStackOosExecution(),
			"alibabacloudstack_arms_alert_contact":                    resourceAlibabacloudStackArmsAlertContact(),
			"alibabacloudstack_arms_alert_contact_group":              resourceAlibabacloudStackArmsAlertContactGroup(),
			"alibabacloudstack_arms_dispatch_rule":                    resourceAlibabacloudStackArmsDispatchRule(),
			"alibabacloudstack_arms_prometheus_alert_rule":            resourceAlibabacloudStackArmsPrometheusAlertRule(),
			"alibabacloudstack_elasticsearch_k8s_instance":            resourceAlibabacloudStackElasticsearchOnk8s(),
			"alibabacloudstack_cloud_firewall_control_policy":         resourceAlibabacloudStackCloudFirewallControlPolicy(),
			"alibabacloudstack_cloud_firewall_control_policy_order":   resourceAlibabacloudStackCloudFirewallControlPolicyOrder(),
			"alibabacloudstack_csb_project":                           resourceAlibabacloudStackCsbProject(),
			"alibabacloudstack_graph_database_db_instance":            resourceAlibabacloudStackGraphDatabaseDbInstance(),
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
		config.OosEndpoint = domain
		config.ArmsEndpoint = domain
		config.CloudfwEndpoint = domain
		config.CsbEndpoint = domain
		config.GdbEndpoint = domain
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
			config.CsbEndpoint = strings.TrimSpace(endpoints["csb"].(string))
			config.GdbEndpoint = strings.TrimSpace(endpoints["gdb"].(string))
			config.DataworkspublicEndpoint = strings.TrimSpace(endpoints["dataworkspublic"].(string))
			config.DbsEndpoint = strings.TrimSpace(endpoints["dbs"].(string))
			config.OosEndpoint = strings.TrimSpace(endpoints["oos"].(string))
			config.ArmsEndpoint = strings.TrimSpace(endpoints["arms"].(string))
			config.CloudfwEndpoint = strings.TrimSpace(endpoints["cloudfw"].(string))
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
		"access_key": "The access key for API operations. You can retrieve this from the 'Security Management' section of the AlibabacloudStack console.",

		"secret_key": "The secret key for API operations. You can retrieve this from the 'Security Management' section of the AlibabacloudStackconsole.",

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
					DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_DATAWORKS_PUBLIC_ENDPOINT", nil),
					Description: descriptions["dataworkspublic_endpoint"],
				},
				"dbs_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_DBS_ENDPOINT", nil),
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
	request := sts.CreateAssumeRoleRequest()
	request.RoleArn = config.RamRoleArn
	if config.RamRoleSessionName == "" {
		config.RamRoleSessionName = "terraform"
	}
	request.RoleSessionName = config.RamRoleSessionName
	//request.DurationSeconds = requests.NewInteger(config.RamRoleSessionExpiration)
	request.Policy = config.RamRolePolicy
	request.Scheme = "https"
	request.SetHTTPSInsecure(config.Insecure)

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
