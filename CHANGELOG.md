# 3.18.17

## New Features

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