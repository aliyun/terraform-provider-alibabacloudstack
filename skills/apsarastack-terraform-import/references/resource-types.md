# Special Import ID Format Quick Reference

Only covers resources with complex import ID formats that are error-prone. For all other resources, check the provider docs directly:
https://registry.terraform.io/providers/aliyun/alibabacloudstack/latest/docs

---

## Composite ID Formats (separated by `:`)

| Resource | Import ID Format | Notes |
|------|---------------|------|
| `alibabacloudstack_route_entry` | `<route-table-id>:<cidr>:<nexthop-type>:<nexthop-id>` | 4 segments |
| `alibabacloudstack_security_group_rule` | `<sg-id>:<direction>:<ip-protocol>:<port-range>:<priority>:<cidr-ip>:<policy>` | 7 segments, order matters |
| `alibabacloudstack_eip_association` | `<eip-id>:<instance-id>` | |
| `alibabacloudstack_snat_entry` | `<snat-table-id>:<snat-entry-id>` | snat-table-id obtained from NAT Gateway attributes |
| `alibabacloudstack_forward_entry` | `<forward-table-id>:<forward-entry-id>` | forward-table-id obtained from NAT Gateway attributes |
| `alibabacloudstack_disk_attachment` | `<disk-id>:<instance-id>` | |
| `alibabacloudstack_nas_mount_target` | `<fs-id>:<mount-target-domain>` | domain is the full domain name, not an ID |
| `alibabacloudstack_db_database` | `<db-instance-id>:<db-name>` | |
| `alibabacloudstack_db_account` | `<db-instance-id>:<account-name>` | |
| `alibabacloudstack_db_connection` | `<db-instance-id>:<connection-prefix>` | |
| `alibabacloudstack_polardb_database` | `<cluster-id>:<db-name>` | |
| `alibabacloudstack_ascm_ram_policy` | `<policy-name>:<policy-type>` | policy-type: `Custom` or `System` |
| `alibabacloudstack_ascm_user_group_role_binding` | `<user-name>:<policy-name>:<policy-type>` | |
| `alibabacloudstack_slb_listener` | `<slb-id>:<protocol>_<port>` | Note underscore, e.g., `lb-xxx:http_80` |
| `alibabacloudstack_cs_kubernetes_node_pool` | `<cluster-id>:<nodepool-id>` | |
| `alibabacloudstack_fc_function` | `<service-name>:<function-name>` | |
| `alibabacloudstack_kms_alias` | `alias/<alias-name>` | Must include `alias/` prefix |

---

## Non-ID Fields as Import Identifiers

| Resource | Import Identifier | Notes |
|------|------------|------|
| `alibabacloudstack_oss_bucket` | `<oss-cluster>:<bucket>` | ApsaraStack requires cluster ID prefix |
| `alibabacloudstack_oss_bucket_acl` | `<oss-cluster>:<bucket>` | Same as above |
| `alibabacloudstack_oss_bucket_policy` | `<oss-cluster>:<bucket>` | Same as above |
| `alibabacloudstack_oss_bucket_cors` | `<oss-cluster>:<bucket>` | Same as above |
| `alibabacloudstack_dns_domain` | Domain string | e.g., `example.com` |
| `alibabacloudstack_ascm_user` | Username | Not UID |
| `alibabacloudstack_ascm_ram_role` | Role name | Not ARN |
| `alibabacloudstack_key_pair` | Key pair name | Not ID |
| `alibabacloudstack_fc_service` | Service name | Not ID |

---

## IDs That Must Be Obtained from Parent Resource Attributes

| Resource | How to Get Required Import ID |
|------|------------------------|
| `alibabacloudstack_snat_entry` | `snat_table_id` obtained from `alibabacloudstack_nat_gateway.<name>.snat_table_ids` (note: comma-separated string, take the first one) |
| `alibabacloudstack_forward_entry` | `forward_table_id` obtained from `alibabacloudstack_nat_gateway.<name>.forward_table_ids` |

```bash
# Get NAT Gateway's snat_table_id
apsara vpc describe-nat-gateways --biz-region-id cn-hangzhou-env01-d01 --nat-gateway-id ngw-xxx \
  | python3 -c "import json,sys; d=json.load(sys.stdin); print(d['NatGateways']['NatGateway'][0]['SnatTableIds']['SnatTableId'][0])"
```

---

## Password/Secret Field Handling

After import, the following fields cannot be read from the API and must be handled in HCL:

```hcl
# Option 1: Use variables
resource "alibabacloudstack_db_instance" "xxx" {
  # ...
}
variable "rds_password" { sensitive = true }

# Option 2: ignore_changes (recommended, avoids plan diff)
lifecycle {
  ignore_changes = [password, account_password]
}
```

Affected resources: RDS (`password`), Redis (`password`), MongoDB (`account_password`), ECS (`password`)
