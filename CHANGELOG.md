# 3.16.19

## New

1. Query capability for OSS clusters
2. Configuration management capability for auto-scaling rules of EDAS K8s applications

## Fixes

1. Fixed the issue where the order of service IDs in SLB VServerGroup caused failure to reach the final state
2. Fixed the issue where KMS encryption could not be enabled for OSS buckets in disaster recovery mode
3. Fixed the issue of incorrect variable types when configuring log monitoring for SLB listeners

---

# 3.16.18

## New

1. Added dynamic retrieval support for the `tls_cipher_policy` type in SLB listeners.

## Fixes

1. Fixed PolarDB for MySQL TDE encryption configuration failure to retrieve KMS encryption keys
2. Added enable/disable control for MongoDB SQL audit capability

---

# 3.16.17

## Fixes

1. chip type is empty after creating polardb PG database
2. polardb database cannot be queried or used normally after binding user
3. mongodb_instance ssl status check fix

## New

1. polardb adds support for ACL function configuration
2. added authorization resources for service roles

---

# 3.16.16

## Fixes

1. Fixed the issue that the length limit of ecs description is 256 characters instead of 256 characters
2. Fixed the issue that the edas_k8s_application package url attribute does not reach the final state due to readback

---

# 3.16.15

## Fixes:

1. slb loadbalancer adds network_type attribute to bring incompatibility
2. Fixed the issue that alibabacloudstack_log_alert did not automatically create charts

## New

1. alibabacloudstack_log_alert supports configuring webhook
2. Implemented alibabacloudstack_acm_configuration

---

# 3.16.14

## Fixes

1. Fixed polardb's ability to modify kernel parameters

## New

1. Added polardb instance configuration tag capability
2. Allow edas k8sapp to bind existing slb instances

---

# 3.16.13

## Fixes

1. Fixed elasticsearch instance creation and deletion capabilities

---

# 3.16.12

## New

1. edas_k8s_application slb binding update, delete function, batch binding function
2. edas_k8s_application supports setting host_aliases

## Fixes

1. edas_k8s_application change memory and cpu capabilities
2. Fixed ots_instance creation failure issue

---

# 3.16.11

## New

1. redis implements ssl enable and disable functions
2. slb supports specifying address for classic network
3. oss_bucket setting tags

## Fixes

1. Fixed the failure of ploardb account during creation due to possible unready resources
2. Fixed the issue that poladb instance cannot reach final state when tde is not enabled due to encrypt_algorithm
3. Fixed the issue that kvstore_instnaceclass cannot normally query redis instance specifications
4. Fixed the issue that cs_cluster cluster nodepool cannot scale nodes

---

# 3.16.10

## New

1. Pull up the import capability of all resources through id, and fix and verify the import function of 73 key resources

## Changes

1. The primary key of alibabacloudstack_ascm_user_group_resource_set_binding is changed from resourceSetId to resourceSetId:userGroupId:ascmRoleId. Resources will be deleted and rebuilt during apply, but the business impact is controllable and does not affect template compatibility
2. alibabacloudstack_ascm_user_group_role_binding is deprecated, and the role_ids in alibabacloudstack_ascm_user_group includes related functions
3. alibabacloudstack_ascm_user_role_binding is deprecated, and the role_ids in alibabacloudstack_ascm_user includes related functions

## Fixes

1. Fixed the issue of incorrect Tde activation when alibabacloudstack_polardb_dbinstance is not mysql engine
2. Fixed the issue of failing to set OSS bucket via environment variables

---

# 3.16.9

## New
1. cs_kubernetes allows setting storage_set when creating
## Fixes
1. Fixed the issue of tde startup failure when polardb uses pgsql engine
2. Supplement some resource documentation information.

---

# 3.16.8

## Fixes

1. Supplement documentation
2. When creating edas_k8s_app, add a cree_repo_id parameter
3. Fixed the issue that polardb_instance reports an error when enabling both Tde and ssl at the same time
4. Fixed the ForceNew issue of edas_k8s_cluster during secondary apply
5. Fixed the read defect of edas_k8s_service and edas_k8s_app that caused resources to fail to reach normal final state

---

# 3.16.7

## Fixes

1. Fixed the issue that cr-ee's namespace and repository failed to schedule in popgw mode
2. Fixed the issue that cr-ee's instance, namespace, and repo query failed

---

# 3.16.6

## Fixes

1. edas_k8s_service: Replace read interface, supplement some read-only attributes. Improve related table documentation
2. cr_repo adjusts name length to 64 to be consistent with page capabilities
3. Fixed some documentation page result error issues

---

# 3.16.5

## New

1. Organization resource set filtering modification exact match.
2. Fixed vpngateway resource creation.
3. HSBC resource smooth migration test
4. Added new resource Polardb_instance
5. Added new resource Polardb_database
6. Added new resource Polardb_account
7. Added new resource Polardb_backup_policy
8. Added new resource Polardb_dbconnection
9. Added new resource Polardb_Zone

---

# 3.16.4

## New

1. When creating Edas_k8s_app, support setting PVC mounting, local mounting, and configuration settings
2. Added new resource Edas_k8s_service
3. Added new resource Edas_namespace

---

# 3.16.3

## Fixes

1. Fixed the issue of alibabacloudstack_datahub_project in popgw mode, comment is no longer allowed to be modified
2. Fixed the issue of alibabacloudstack_datahub_topic in popgw mode
3. Fixed the issue of alibabacloudstack_datahub_subscription in popgw mode, comment is no longer allowed to be modified

---

# 3.16.2

## Fixes

1. Fixed DNS does not support HTTPS request protocol, all AliDNS resources force the use of http
2. OSS_Bucket creation is compatible with 318x data format
3. Fixed the issue of OSS using ECS endpoint in popgw mode
4. db_readonly_instance logic vulnerability fix
5. kvstore_instance_class issue series optional value case specification
6. Fixed the issue that some test cases cannot export test templates

---

# 3.16.1

## New

1. Support for apsarastack user smooth migration
2. Support for sls log_alert resource CRUD operations
3. TDE enable capability for redis_instance (kvstore_instance)
4. multi_az_policy attribute support for ess_scalinggroup
5. data_disk_tags for ecs instance
6. auto_snapshot_policy_id and enable_automated_snapshot_policy support for ecs_disk
7. oss_bucket adds bucket_sync (synchronization), storage_capacity (capacity limit), sse_algorithm (encryption method)
8. slb_listener adds logs_download_attributes
9. Encryption function for image_copy
10. cs_k8s function to specify security_group_id

## Fixes

1. Fixed the issue that when ecs updates image, system_disk tags could not be inherited normally

## Optimizations

1. Optimized slb_listener acl-related parameter dependency mutex check logic
2. Dependency on github.com/aliyun/aliyun-datahub-sdk-go

---

# 3.16.0

## New

1. Support ASAPI and POPAPI scheduling, support separate designation of cloud product endpoints
2. Support ASAPI permission processing mechanism, all requests will be authenticated and message error encapsulated through ASAPI
3. Support Eagle Eye ID, optimize error messages and message prompts
4. Supplement and improve test coverage through automated generation, fix test cases and Provider logic, and improve test pass rate