# 3.18.11

## Fixes
1. Fixed the issue where `alibabacloudstack_log_alert` did not automatically create charts

## Additions
1. `alibabacloudstack_log_alert` now supports webhook configuration
2. Implemented `alibabacloudstack_acm_configuration`

---

# 3.18.10

## Additions
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

## Additions
1. Added `BandwidthPackage` resource for NatGateway
2. Added `AccessLog` resource for SLB
3. Added `vpc_ipv6_isps` resource (IPv6 CIDR block support) with multi-resource binding via `ipv6_cidr_blocks`
4. Added `VSwitchNetworkAclAttachment` resource for VPC
5. Added `HaVip` resource (High-Availability Virtual IP) for VPC

## Deprecations
1. Marked SLB Listener's `logs_download_attributes` as deprecated. Use `AccessLog` resource instead

---

# 3.18.8

## Additions
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

## Additions
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

## Additions
1. Added SLB binding update/delete and batch binding for `edas_k8s_application`
2. Added `host_aliases` support for `edas_k8s_application`

## Fixes
1. Fixed memory/CPU modification capability for `edas_k8s_application`
2. Fixed `ots_instance` creation failure
3. Fixed Elasticsearch instance creation/deletion capabilities

---

# 3.18.5

## Additions
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

## Additions
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

## Additions
1. Added PVC/local mount and configuration support for `edas_k8s_app`
2. New resources:
   - `edas_k8s_service`
   - `edas_namespace`

---

# 3.18.2

## Additions
1. Added management support for:
   - `alibabacloudstack_bastionhost_instance` (Bastion Host)
   - `alibabacloudstack_waf_instance` (Web Application Firewall)
   - `alibabacloudstack_vpn_gateway` (VPN Gateway)

## Fixes
1. Fixed `vpn_gateway` creation failure

---

# 3.18.1

## Additions
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