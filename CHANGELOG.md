# 3.18.25

## Added

1. Added query capability for MQTT clusters (alibabacloudstack_mqtt_clusters).
2. Added query capability for ONS clusters (alibabacloudstack_ons_clusters).
3. alibabacloudstack_ascm_organization now supports retrieving PK (Primary Key) information.

## Fixed

1. Fixed an issue where Dataworks Remain was not functioning correctly.
2. Fixed an issue where Dataworks Connection was not functioning correctly.



---

# 3.18.24

## Added

1. Query capability for OTS clusters (`alibabacloudstack_ots_clusters`)
2. Configuration capability for PolarDB read/write splitting connections (`alibabacloudstack_polardb_readwrite_splitting_connection`)
3. Orchestration capability for DataWorks workflows (`alibabacloudstack_data_works_business`)
4. Query capability for ECS dedicated host types (`alibabacloudstack_ecs_dedicated_host_types`)
5. Configuration capability for PolarDB parameter groups (`alibabacloudstack_polardb_parameter_group`)
6. Binding capability for Quick BI user group members (`alibabacloudstack_quick_bi_user_group_user`)
7. Creation capability for CSB services (`alibabacloudstack_csb_service`)
8. Query capability for DataWorks file types (`alibabacloudstack_dataworks_file_types`)
9. Orchestration capability for DataWorks files (`alibabacloudstack_data_works_file`)
10. Orchestration capability for DataWorks intelligent baselines (`alibabacloudstack_data_works_baseline`)
11. Custom port support added for ALIKafka instances
12. Capability to bind ECS snapshot policies to disk IDs

## Fixed

1. Refactored orchestration capabilities for OTS-related resources  
2. Resolved issue where DataWorks resources were not orchestratable

## Deprecated

1. Deprecate orchestration capability for ARMS Prometheus alert rules (`alibabacloudstack_arms_prometheus_alert_rule`)
2. Deprecate orchestration capability for ARMS dispatch rules (`alibabacloudstack_arms_dispatch_rule`)
3. Deprecate orchestration capability for CMS site monitoring (`alibabacloudstack_cms_site_monitor`)
4. Deprecate orchestration capability for CMS alarm contacts (`alibabacloudstack_cms_alarm_contact`)
5. Deprecate orchestration capability for CMS alarm contact groups (`alibabacloudstack_cms_alarm_contact_group`)
6. Deprecate orchestration capability for MongoDB backup plans (`alibabacloudstack_dbs_backup_plan`)
7. Deprecate orchestration capability for Quick BI workspaces (`alibabacloudstack_quick_bi_workspace`)
8. Deprecate query capability for OTS services (`alibabacloudstack_ots_service`)

---

# 3.18.23

## New

1. Orchestration support for ACR Enterprise Edition image lifecycle rules (`alibabacloudstack_cr_ee_attestor_lifecycle_rule`)
2. Orchestration support for binding relationships between resource sets and users (`alibabacloudstack_ascm_resource_group_user_attachment`)
3. Orchestration support for IP or port address books in Cloud Firewall (`alibabacloudstack_cloudfw_address_book`)
4. Orchestration support for VPC network traffic access control policies (`alibabacloudstack_cloudfw_vpc_control_policy`)
5. Orchestration support for OSS single tunnel (`alibabacloudstack_oss_single_tunnel`)
6. Orchestration support for Redis parameter templates (`alibabacloudstack_kvstore_parameter_group`)
7. Orchestration support for GPDB backup policies (`alibabacloudstack_gpdb_backup_policy`)
8. Orchestration support for NAS directory quotas (`alibabacloudstack_nas_dir_quota`)
9. Orchestration support for NAS namespaces (`alibabacloudstack_nas_namespace`)
10. Orchestration support for NAS unified namespace filesystem mappings (`alibabacloudstack_nas_namespace_filesystem_attachment`)
11. Orchestration support for NAS namespace filesystem mount targets (`alibabacloudstack_nas_namespace_mount_target`)
12. Orchestration support for NAS cross-region mount configurations (`alibabacloudstack_nas_namespace_group`)
13. Data source for querying GPDB instance types (`alibabacloudstack_gpdb_instance_types`)
14. Data source for querying ADB cluster types (`alibabacloudstack_adb_cluster_types`)
15. Data source for querying EDAS Kubernetes clusters (`alibabacloudstack_edas_k8s_clusters`)

## Fixes

1. Batch-fixed issues where data sources could not correctly filter resources using the `ids` attribute
2. Batch-removed unsupported tags fields from data sources that couldn't properly filter by `tags`, and removed unsupported `tags` fields from resources
3. Batch-fixed multiple code defects and corrected test cases

## Changes

1. Changed the type of the bandwidth attribute in `alibabacloudstack_expressconnect_physicalconnection` from string to integer
2. Removed previously unsupported fields (`dns_servers`, `group_id`, `lang`, etc.) from `alibabacloudstack_dns_domain`
3. Removed unsupported field `on_unable_to_redeploy_failed_instance` from `alibabacloudstack_ecs_deployment_set`

## Deprecations

1. DNS domain group orchestration capability (alibabacloudstack_alidns_domaingroup)
2. API Gateway service data source (alibabacloudstack_api_gateway_service)
3. ESS lifecycle management capability (alibabacloudstack_ess_lifecycle_hook)
4. ESS scaling group virtual server group orchestration capability (alibabacloudstack_ess_scalinggroup_vserver_groups)
5. KMS secret capability (alibabacloudstack_kms_secret)

---

# 3.18.22

## New

1. Orchestration support for SchedulerX application groups (`alibabacloudstack_schedulerx2_app_group`)
2. Orchestration support for SchedulerX jobs (`alibabacloudstack_schedulerx2_job`)
3. Orchestration support for SchedulerX workflows (`alibabacloudstack_schedulerx2_workflow`)
4. Data source for querying APFS zones (`alibabacloudstack_apfs_zones`)
5. Orchestration support for APFS file systems (`alibabacloudstack_apfs_file_system`)
6. Data source for querying HSM vendors and their products (`alibabacloudstack_hsm_vendors`)
7. Data source for listing HSM instances (`alibabacloudstack_hsms`)
8. Orchestration support for HSM instances (`alibabacloudstack_hsm_instances`)
9. Orchestration support for HSM clusters (`alibabacloudstack_hsm_clusters`)
10. Orchestration support for bare metal server (BMS) key pairs (`alibabacloudstack_bms_keypairs`)
11. Orchestration support for AQS OSS scan configurations (`alibabacloudstack_aqs_oss_scanconfigs`)
12. Orchestration support for AQS anti-brute-force protection rules for hosts (`alibabacloudstack_aqs_anti_brute_force_rules`)
13. Orchestration support for AQS web tamper protection configurations (`alibabacloudstack_aqs_web_locks`)
14. Data source for querying CSP Private HSM vendors and their products (`alibabacloudstack_cspprivate_hsm_vendors`)
15. Data source for listing CSP Private HSM instances (`alibabacloudstack_cspprivate_hsms`)
16. Orchestration support for CSP Private HSM instances (`alibabacloudstack_cspprivate_hsm_instance`)
17. Orchestration support for CSP Private HSM clusters (`alibabacloudstack_cspprivate_hsm_group`)

## Removals

1. Data source for querying service lists under an environment (`environment_services_by_product`)

---

# 3.18.21

## New

1. Orchestration capability for Alibaba Cloud DNS forward domains (`alibabacloudstack_dns_forward_domain`)
2. Orchestration capability for Alibaba Cloud DNS line configurations (`alibabacloudstack_dns_line`)
3. Orchestration capability for Alibaba Cloud private DNS domains (`alibabacloudstack_dns_private_domain`)
4. Orchestration capability for DNS record configurations of private domains (`alibabacloudstack_dns_private_record`)
5. Orchestration capability for recursive ACL policies of private DNS domains (`alibabacloudstack_dns_recursor_acl`)
6. Orchestration capability for line configurations of private DNS domains (`alibabacloudstack_dns_private_line`)
7. Orchestration capability for Global Traffic Manager (GTM) instances (`alibabacloudstack_dns_gtm_instance`)
8. Orchestration capability for GTM address pools (`alibabacloudstack_dns_gtm_addresspool`)
9. Orchestration capability for GTM instance access strategies (`alibabacloudstack_dns_gtm_access_strategy`)
10. Orchestration capability for CPFS file system resources (`alibabacloudstack_cpfs_file_system`)
11. Orchestration capability for Prometheus v2 monitoring instances (`alibabacloudstack_prometheus_v2_instance`)
12. Orchestration capability for configuring alert notification recipients in Prometheus v2 (`alibabacloudstack_prometheus_v2_contact`)
13. Orchestration capability for notification groups in Prometheus v2 (`alibabacloudstack_prometheus_v2_notify_group`)
14. Orchestration capability for alert rules in Prometheus v2 (`alibabacloudstack_prometheus_v2_alert`)

## Removals

1. Deprecation of the query capability for ECS instance metering data (`alibabacloudstack_ascm_metering_query_ecs`)

---

# 3.18.20

## New

1.Query capability for API Gateway V2 instance types (`alibabacloudstack_api_gateway_v2_instance_types`)
2. Orchestration capability for API Gateway V2 instances (`alibabacloudstack_api_gateway_v2_instance`)
3. Orchestration capability for API Gateway V2 cascade instances (`alibabacloudstack_api_gateway_v2_cascade_instance`)
4. Orchestration capability for API Gateway V2 cascade links (`alibabacloudstack_api_gateway_v2_cascade_link`)
5. Import capability for API Gateway V2 Kubernetes clusters (`alibabacloudstack_api_gateway_v2_k8s_cluster`)
6. Orchestration capability for API Gateway V2 certificates (`alibabacloudstack_api_gateway_v2_certificate`)
7. Orchestration capability for API Gateway V2 consumers (`alibabacloudstack_api_gateway_v2_consumer`)
8. Orchestration capability for API Gateway V2 domains (`alibabacloudstack_api_gateway_v2_domain`)
9. Orchestration capability for API Gateway V2 MCP services (`alibabacloudstack_api_gateway_v2_mcpservice`)
10. Orchestration capability for API Gateway V2 route groups (`alibabacloudstack_api_gateway_v2_route_group`)
11. Orchestration capability for API Gateway V2 routes (`alibabacloudstack_api_gateway_v2_route`)
12. Orchestration capability for API Gateway V2 services (`alibabacloudstack_api_gateway_v2_service`)
13. Orchestration capability for API Gateway V2 service sources (`alibabacloudstack_api_gateway_v2_service_source`)
14. Orchestration capability for API Gateway V2 signatures (`alibabacloudstack_api_gateway_v2_signature`)
15. Orchestration capability for DataHub Kafka consumer groups (`alibabacloudstack_datahub_kafka_group`)
16. Orchestration capability for ApsaraDB for Lindorm instances (`alibabacloudstack_lindorm_instance`)
17. Query capability for ApsaraDB for Lindorm instance types (`alibabacloudstack_lindorm_instance_types`)
18. Orchestration capability for ApsaraDB for Lindorm LTS instances (`alibabacloudstack_lindorm_lts_instance`)
19. Orchestration capability for Message Queue for MQTT instances (`alibabacloudstack_mqtt_instance`)
20. Orchestration capability for Message Queue for MQTT topics (`alibabacloudstack_mqtt_topic`)
21. Orchestration capability for Message Queue for MQTT groups (`alibabacloudstack_mqtt_group`)
22. Orchestration capability for Universal DNS domains (`alibabacloudstack_universal_dns_domain`)
23. Orchestration capability for Universal DNS record configurations (`alibabacloudstack_universal_dns_record`)
24. Orchestration capability for Universal DNS resolution lines (`alibabacloudstack_universal_dns_line`)

---

# 3.18.19

## New

1. ACK template orchestration capability (`alibabacloudstack_ack_template`)
2. EDAS swimming lane orchestration capability (`alibabacloudstack_edas_swimming_lane`)
3. EDAS swimming lane group orchestration capability (`alibabacloudstack_edas_swimming_lane_group`)
4. Hologram cluster data source (`alibabacloudstack_hologram_clusters`)
5. Hologram instance orchestration capability (`alibabacloudstack_hologram_instance`)
6. Hologram instance backup policy orchestration capability (`alibabacloudstack_hologram_instance_backup_policy`)

---

# 3.18.18

## New

1. Query capability for OSS clusters
2. Configuration management capability for auto-scaling rules of EDAS K8s applications

## Fixes

1. Fixed the issue where creating TDE-enabled PolarDB for PostgreSQL versions caused errors.
2. Added default value handling for NIC type in security group rules (`security_group_rule`) to prevent potential failures.
3. Resolved improper exception handling in EDAS Kubernetes application scaling rules (`edas_k8s_application_scaling_rule`) with invalid parameter configurations.
4. Fixed the issue where the order of service IDs in SLB VServerGroup caused failure to reach the final state
5. Fixed the issue where KMS encryption could not be enabled for OSS buckets in disaster recovery mode
6. Fixed the issue of incorrect variable types when configuring log monitoring for SLB listeners

---

# 3.18.17

## New

1. Orchestration capability for CEN border router network instance connections (`alibabacloudstack_cen_transit_router_vbr_attachment`)
2. Orchestration capability for CEN Connect network instance connections (`alibabacloudstack_cen_transit_router_connect_attachment`)
3. CEN VBR health check capability (`alibabacloudstack_cen_vbr_health_check`)
4. CEN Connect network instance connections now support peer configuration (`alibabacloudstack_cen_connect_peer`)
5. Orchestration capability for CEN multicast domain members (`alibabacloudstack_cen_transit_router_multicast_domain_member`)
6. Added Connect support for CEN multicast domain sources (`alibabacloudstack_cen_transit_router_multicast_domain_source`)

---

# 3.18.16

## New

1. Polardb cluster instance orchestration capability (`alibabacloudstack_polardb_cluster_instance`)
2. Polardb cluster account orchestration capability (`alibabacloudstack_polardb_cluster_account`)
3. Polardb cluster database orchestration capability (`alibabacloudstack_polardb_cluster_database`)
4. Polardb cluster account permission configuration capability (`alibabacloudstack_polardb_cluster_account_database_binding`)
5. Polardb cluster proxy configuration capability (`alibabacloudstack_polardb_cluster_proxy`)
6. Polardb cluster backup policy configuration capability (`alibabacloudstack_polardb_cluster_backup_policy`)
7. Polardb cluster instance type query capability (`alibabacloudstack_polardb_cluster_instance_types`)
8. Polardb cluster proxy specification query capability (`alibabacloudstack_polardb_cluster_proxy_types`)

---

# 3.18.15

## Changes
1. Added support for dynamically retrieving available TLS cipher policy options in the environment for `alibabacloudstack_slb_listener`

## Fixes
1. Fixed multiple unit test cases
2. Fixed creation failure when OSS bucket under resource group is empty
3. Fixed inability to filter `alibabacloudstack_security_group_rules` by `nic_type`
4. Fixed ineffective AddressType configuration for `alibabacloudstack_slb_loadbalancer`
5. Fixed query failures for `alibabacloudstack_ons_instances` caused by API response data type changes
6. Fixed signature verification failures in KVStore instances caused by `clientToken` parameter issues
7. Fixed non-terminal state handling for `alibabacloudstack_dns_record` resources
8. Fixed data retrieval failures in `alibabacloudstack_cr_repos`
9. Fixed unstable `cluster_id` values causing state mismatches when creating `alibabacloudstack_nas_file_system`
10. Fixed `alibabacloudstack_ascm_user_group_role_binding` issues caused by ASAPI unavailability
11. **(Incompatible)** Fixed `alibabacloudstack_ascm_user_role_binding` issues caused by ASAPI unavailability.  Modified `role_ids` type to `Set[INT]`
12. Fixed PolarDB for MySQL TDE encryption configuration failure to retrieve KMS keys

---

# 3.18.14

## New

1. Added support for ACL functionality configuration in PolarDB  
2. Added authorization resources for service roles  

## Fixes

1. Fixed empty chip type display after PolarDB PG database creation  
2. Fixed issue where bound users couldn't be queried or used normally after PolarDB database user binding  
3. Fixed SSL enablement status check for MongoDB instances  
4. MaxCompute CU resource adaptation for version 3.18  
5. MaxCompute User resource adaptation for version 3.18  
6. MaxCompute Project resource adaptation for version 3.18  

---

# 3.18.13

## New

1. Added support for managing instances of Enterprise Cloud Network CEN (cen_instance)
2. Added support for CEN transit router route tables (cen_transit_router_route_table)
3. Added support for CEN transit router route entries (cen_transit_router_route_entry)
4. Added support for managing VPC-type transit router attachments in CEN (cen_transit_router_vpc_attachment)
5. Added support for associating CEN transit router route tables (cen_transit_router_route_table_association)
6. Added support for propagating routes in CEN transit router route tables (cen_transit_router_route_table_propagation)
7. Added support for CEN route maps (alibabacloudstack_cen_route_map)
8. Added support for CEN multicast domains (cen_transit_router_multicast_domain)
9. Added support for binding multicast switches in CEN (cen_transit_router_multicast_domain_association)
10. Added support for CEN multicast domain sources (cen_transit_router_multicast_domain_source)
11. Added support for querying Polardbx 2.0 instance specifications (instance_type)
12. Added support for managing Polardbx 2.0 instances (polardbx_instance)
13. Added support for managing Polardbx 2.0 read-only instances (polardbx_readonly_instance)
14. Added support for managing Polardbx 2.0 user accounts (polardbx_account)
15. Added support for configuring Polardbx 2.0 super accounts (including three-tier separation) (polardbx_super_account)
16. Added support for managing Polardbx 2.0 database tables (polardbx_database)
17. Added support for managing Polardbx 2.0 backups (polardbx_backup)
18. Added support for configuring Polardbx 2.0 backup policies (polardbx_backup_policy)
19. Added support for binding accounts and databases in Polardbx 2.0 (polardbx_account_database_binding)
20. Added support for configuring Polardbx 2.0 instance log engines (polardbx_log_engine)
21. Added support for configuring read/write splitting in Polardbx 2.0 (polardbx_read_write_splitting_config)

---

# 3.18.12

## New

1. Added MongoDB account management capability (`mongodb_account`)
2. Added MongoDB backup management capability (`mongodb_backup`)
3. Added network connection management capability for CS nodes in MongoDB sharded instances (`mongodb_shardinginstance_csnode_address`)
4. Added network connection management capability for Shard nodes in MongoDB sharded instances (`mongodb_shardinginstance_shardnode_address`)
5. Added network connection management capability for Mongos nodes in MongoDB sharded instances (`mongodb_shardinginstance_mongonode_address`)
6. Added DRDS private RDS instance management capability (`drds_rds_instance`)
7. Added DRDS private database management capability (`drds_rds_instance`)
8. Added DRDS account management capability (`drds_account`)
9. Added DRDS instance series query capability (`drds_instance_series`)
10. Added DRDS instance specifications query capability (`drds_instance_specifications`)
11. Added DRDS private read-only instance management capability (`readonly_instance`)
12. Added PolarDB backup management capability (`polardb_backup`)
13. Added PolarDB instance type specifications query capability (`polardb_instance_types`)
14. Added PolarDB read-only instance management capability (`polardb_readonly_instance`)
15. Added PolarDB database proxy (including read/write splitting) management capability (`polardb_proxy`)
16. Added RDS instance type specifications query capability (`rds_instance_types`)
17. Added RDS backup management capability (`rds_backup`)
18. Added RDS proxy management capability (`rds_backup`)
19. Added database connection address management capability for MongoDB replica instances
20. Added management capability for CS nodes in MongoDB sharded instances

## Changes

1. **(Incompatible)**: Added support for setting descriptions on nodes in MongoDB sharded instances. This field is now **mandatory in Terraform (TF)** and must be **unique within the same instance**.

---

# 3.18.11

## Fixes

1. Fixed the issue where `alibabacloudstack_log_alert` did not automatically create charts

## New

1. `alibabacloudstack_log_alert` now supports webhook configuration
2. **(Incompatible)** `alibabacloudstack_acm_configuration`

---

# 3.18.10

## New

1. Added VPN resources:
   - `vpn_pbr_route_entry` (Policy-Based Routing Table)
   - `ssl_vpnserver` (SSL Server)
   - `ssl_vpn_client_cert` (SSL Client Certificate)
2. Added ExpressConnect resources:
   - `bgp_group` (BGP Group)
   - `bgp_peer` (BGP Peer)
   - `vbr_ha` (Fast Failover Group)
   - `vbr_pconn_association` (Multi-Physical Connection)
   - `bgp_network` (BGP Network)
3. Added NAS resource `nas_lifecycle_policy` (Lifecycle Policy)

## Fixes

1. **(Incompatible)** Major structural changes in NAS resource availability zone query responses

## Deprecations

1. **(Incompatible)** Deprecated `protocols` query in NAS resources. Use the new `Zone` query instead

---

# 3.18.9

## New

1. Added `BandwidthPackage` resource for NatGateway
2. Added `AccessLog` resource for SLB
3. Added `vpc_ipv6_isps` resource (IPv6 CIDR block support) with multi-resource binding via `ipv6_cidr_blocks`
4. Added `VSwitchNetworkAclAttachment` resource for VPC
5. Added `HaVip` resource (High-Availability Virtual IP) for VPC

## Deprecations

1. Marked SLB Listener's `logs_download_attributes` as deprecated. Use `AccessLog` resource instead

---

# 3.18.8

## New

1. Added EBS resources:
   - `DiskReplicaPair` (Disk Async Replication)
   - `DiskReplicaGroup` (Consistency Replication Group)
2. Added ECS resources:
   - `SnapshotGroup` (Snapshot Consistency Group)
   - `DedicatedHostCluster` (Dedicated Host Cluster)
   - `Invocation` (Command Execution)
   - `DhcpOptionsSet` (DHCP Options Set)
3. Added `scaling_group_id` support for ESS `ScheduledTask`

## Removals

1. Removed `LaunchTemplate` functionality (unavailable in ASCM console)

---

# 3.18.7

## New

1. Added `flink_namespace` creation/query capabilities
2. Added `tag` configuration support for PolarDB instances
3. Added disk encryption support for CS Kubernetes clusters
4. Enabled EDAS K8sApp to bind existing SLB instances

## Fixes

1. Fixed PolarDB kernel parameter modification capability

## Deprecations
1. Marked all DataSource `output_file` attributes as deprecated (scheduled for removal in 3.19.0). Use `local_file` provider instead

---

# 3.18.6

## New

1. Added SLB binding update/delete and batch binding for `edas_k8s_application`
2. Added `host_aliases` support for `edas_k8s_application`

## Fixes
1. Fixed memory/CPU modification capability for `edas_k8s_application`
2. Fixed `ots_instance` creation failure
3. Fixed Elasticsearch instance creation/deletion capabilities

---

# 3.18.5

## New

1. Added `alikafka_instance` creation/query capabilities
2. Implemented SSL enable/disable for Redis
3. Added classic network `address` specification for SLB
4. Added `tags` support for `oss_bucket`

## Changes

1. Changed primary key of `alibabacloudstack_ascm_user_group_resource_set_binding` from `resourceSetId` to `resourceSetId:userGroupId:ascmRoleId` (recreation required on apply)
2. Deprecated `alibabacloudstack_ascm_user_group_role_binding` (functionality merged into `alibabacloudstack_ascm_user_group.role_ids`)
3. Deprecated `alibabacloudstack_ascm_user_role_binding` (functionality merged into `alibabacloudstack_ascm_user.role_ids`)

## Fixes

1. Fixed TDE activation failure for PolarDB PostgreSQL engines
2. Fixed OSS bucket configuration via environment variables
3. Fixed TDE activation errors for non-MySQL PolarDB engines
4. Fixed PolarDB account creation readiness issues
5. Fixed PolarDB instance state issues when TDE is disabled
6. Fixed node scaling in CS cluster node pools

---

# 3.18.4

## New

1. Implemented exact matching for organizational resource set filtering
2. Added PolarDB resources:
   - `polardb_instance`
   - `polardb_database`
   - `polardb_account`
   - `polardb_backup_policy`
   - `polardb_dbconnection`
   - `polardb_zone`

## Fixes

1. Enhanced `edas_k8s_service` with new read API and documentation
2. Adjusted `cr_repo` name length to 64 characters (UI alignment)
3. Fixed CR-EE resource queries in popgw mode
4. Added `cree_repo_id` parameter for `edas_k8s_app`
5. Fixed PolarDB TDE+SSL coexistence issue
6. Resolved `edas_k8s_cluster` ForceNew on reapply
7. Fixed terminal state issues for EDAS resources

---

# 3.18.3

## New

1. Added PVC/local mount and configuration support for `edas_k8s_app`
2. New resources:
   - `edas_k8s_service`
   - `edas_namespace`

---

# 3.18.2

## New

1. Added management support for:
   - `alibabacloudstack_bastionhost_instance` (Bastion Host)
   - `alibabacloudstack_waf_instance` (Web Application Firewall)
   - `alibabacloudstack_vpn_gateway` (VPN Gateway)

## Fixes

1. Fixed `vpn_gateway` creation failure

---

# 3.18.1

## New

1. Added disk encryption methods

## Fixes
1. Fixed system disk tag loss during ECS image replacement
2. Fixed default disk encryption tag issues
3. Fixed encrypted snapshot disk creation
4. Resolved SLS Alarm `mute_util` validation
5. Fixed popgw mode issues for DataHub resources

---

# 3.18.0

> Branch: 3.16.2

## Fixes

1. Modified alikafka popgw domain format for 318x compatibility **(apsarastack 3.16.2 incompatible)**

## Removals

1. Removed `domain` and `force_use_asapi` ASAPI configurations