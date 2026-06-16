# ApsaraStack Terraform deprecated fields reference

Deprecated field mappings for `alibabacloudstack_*` resources.
Consult this file during Step 5.1 before writing any field.

| Resource | Field | Kind | Action | Since |
| --- | --- | --- | --- | --- |

| `alibabacloudstack_db_instance` | `name` | rename | Use `db_instance_name` instead; value semantics unchanged. | v3.16.0 |
| `alibabacloudstack_db_database` | `name` | rename | Use `db_name` instead; value semantics unchanged. | v3.16.0 |
| `alibabacloudstack_security_group` | `name` | rename | Use `security_group_name` instead; value semantics unchanged. | v3.16.0 |
| `alibabacloudstack_vpc` | `name` | rename | Use `vpc_name` instead; value semantics unchanged. | v3.16.0 |
| `alibabacloudstack_vswitch` | `name` | rename | Use `vswitch_name` instead; value semantics unchanged. | v3.16.0 |
| `alibabacloudstack_instance` | `instance_name` | rename | Use `name` instead; value semantics unchanged. | v3.16.4 |
| `alibabacloudstack_oss_bucket` | `acl` | split | Inline ACL deprecated — declare separate `alibabacloudstack_oss_bucket_acl` resource alongside the parent and remove inline `acl = …`. | v3.16.0 |
| `alibabacloudstack_oss_bucket` | `cors_rule` | split | Inline CORS deprecated — declare separate `alibabacloudstack_oss_bucket_cors` resource. | v3.16.0 |
| `alibabacloudstack_oss_bucket` | `logging` | split | Inline logging deprecated — declare separate `alibabacloudstack_oss_bucket_logging` resource. | v3.16.0 |
| `alibabacloudstack_oss_bucket` | `website` | split | Inline website configuration deprecated — declare separate `alibabacloudstack_oss_bucket_website` resource. | v3.16.0 |
| `alibabacloudstack_slb` | `name` | rename | Use `load_balancer_name` instead; value semantics unchanged. | v3.16.0 |
| `alibabacloudstack_slb_listener` | `bandwidth` | deprecated-no-replacement | Deprecated without a documented replacement — use `alibabacloudstack_slb_listener_bandwidth` sub-resource. | v3.16.0 |
