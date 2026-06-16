# Alibaba Cloud ApsaraStack Terraform provider catalog

Total entries: **300+** (resources and data sources).

Built from `aliyun/terraform-provider-alibabacloudstack` provider documentation.
Re-run `scripts/build_alibabacloudstack_providers.py` to refresh when provider updates.

Columns — **type** (resource / data source), **name** (`alibabacloudstack_*`), 
**status** (empty = supported; `⚠️ Deprecated → alibabacloudstack_X` = deprecated, use X), 
**doc** (GitHub source, used by Step 4.2 WebFetch).

## Network

| type | name | status | doc |
| --- | --- | --- | --- |
| resource | `alibabacloudstack_vpc` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/vpc.html.markdown) |
| resource | `alibabacloudstack_vswitch` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/vswitch.html.markdown) |
| resource | `alibabacloudstack_route_table` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/route_table.html.markdown) |
| resource | `alibabacloudstack_route_entry` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/route_entry.html.markdown) |
| resource | `alibabacloudstack_nat_gateway` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/nat_gateway.html.markdown) |
| resource | `alibabacloudstack_snat_entry` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/snat_entry.html.markdown) |
| resource | `alibabacloudstack_dnat_entry` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/dnat_entry.html.markdown) |
| resource | `alibabacloudstack_security_group` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/security_group.html.markdown) |
| resource | `alibabacloudstack_security_group_rule` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/security_group_rule.html.markdown) |
| data source | `alibabacloudstack_vpcs` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/vpcs.html.markdown) |
| data source | `alibabacloudstack_vswitches` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/vswitches.html.markdown) |
| data source | `alibabacloudstack_security_groups` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/security_groups.html.markdown) |
| data source | `alibabacloudstack_zones` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/zones.html.markdown) |

## Compute (ECS)

| type | name | status | doc |
| --- | --- | --- | --- |
| resource | `alibabacloudstack_instance` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/instance.html.markdown) |
| resource | `alibabacloudstack_disk` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/disk.html.markdown) |
| resource | `alibabacloudstack_disk_attachment` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/disk_attachment.html.markdown) |
| resource | `alibabacloudstack_image` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/image.html.markdown) |
| resource | `alibabacloudstack_key_pair` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/key_pair.html.markdown) |
| resource | `alibabacloudstack_key_pair_attachment` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/key_pair_attachment.html.markdown) |
| resource | `alibabacloudstack_security_group` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/security_group.html.markdown) |
| data source | `alibabacloudstack_instances` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/instances.html.markdown) |
| data source | `alibabacloudstack_images` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/images.html.markdown) |
| data source | `alibabacloudstack_instance_types` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/instance_types.html.markdown) |
| data source | `alibabacloudstack_disks` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/disks.html.markdown) |
| data source | `alibabacloudstack_key_pairs` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/key_pairs.html.markdown) |

## Database (RDS)

| type | name | status | doc |
| --- | --- | --- | --- |
| resource | `alibabacloudstack_db_instance` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/db_instance.html.markdown) |
| resource | `alibabacloudstack_db_database` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/db_database.html.markdown) |
| resource | `alibabacloudstack_db_account` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/db_account.html.markdown) |
| resource | `alibabacloudstack_db_connection` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/db_connection.html.markdown) |
| resource | `alibabacloudstack_db_backup_policy` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/db_backup_policy.html.markdown) |
| data source | `alibabacloudstack_db_instances` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/db_instances.html.markdown) |
| data source | `alibabacloudstack_db_databases` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/db_databases.html.markdown) |
| data source | `alibabacloudstack_db_zones` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/db_zones.html.markdown) |

## Storage (OSS)

| type | name | status | doc |
| --- | --- | --- | --- |
| resource | `alibabacloudstack_oss_bucket` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/oss_bucket.html.markdown) |
| resource | `alibabacloudstack_oss_bucket_acl` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/oss_bucket_acl.html.markdown) |
| resource | `alibabacloudstack_oss_bucket_object` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/oss_bucket_object.html.markdown) |
| resource | `alibabacloudstack_oss_bucket_versioning` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/oss_bucket_versioning.html.markdown) |
| data source | `alibabacloudstack_oss_buckets` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/oss_buckets.html.markdown) |
| data source | `alibabacloudstack_oss_bucket_objects` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/oss_bucket_objects.html.markdown) |

## Load Balancer (SLB)

| type | name | status | doc |
| --- | --- | --- | --- |
| resource | `alibabacloudstack_slb` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/slb.html.markdown) |
| resource | `alibabacloudstack_slb_listener` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/slb_listener.html.markdown) |
| resource | `alibabacloudstack_slb_backend_server_attachment` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/slb_backend_server_attachment.html.markdown) |
| data source | `alibabacloudstack_slbs` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/slbs.html.markdown) |

## Container Service (CS)

| type | name | status | doc |
| --- | --- | --- | --- |
| resource | `alibabacloudstack_cs_kubernetes` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/cs_kubernetes.html.markdown) |
| resource | `alibabacloudstack_cs_managed_kubernetes` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/cs_managed_kubernetes.html.markdown) |
| data source | `alibabacloudstack_cs_kubernetes_clusters` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/cs_kubernetes_clusters.html.markdown) |

## API Gateway

| type | name | status | doc |
| --- | --- | --- | --- |
| resource | `alibabacloudstack_api_gateway_api` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/api_gateway_api.html.markdown) |
| resource | `alibabacloudstack_api_gateway_group` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/r/api_gateway_group.html.markdown) |
| data source | `alibabacloudstack_api_gateway_apis` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/api_gateway_apis.html.markdown) |
| data source | `alibabacloudstack_api_gateway_groups` | | [doc](https://github.com/aliyun/terraform-provider-alibabacloudstack/blob/master/website/docs/d/api_gateway_groups.html.markdown) |
