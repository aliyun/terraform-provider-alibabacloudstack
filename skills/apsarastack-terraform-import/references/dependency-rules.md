# Resource Dependency Rules

Resource dependency layer diagram and import order constraints.

---

## Dependency Layer Diagram

```
Layer 0 (no dependencies, can import in parallel)
├── alibabacloudstack_vpc_vpc
├── alibabacloudstack_oss_bucket
├── alibabacloudstack_dns_domain
├── alibabacloudstack_ascm_user
├── alibabacloudstack_ascm_ram_role
├── alibabacloudstack_ascm_ram_policy
├── alibabacloudstack_kms_key
├── alibabacloudstack_key_pair
└── alibabacloudstack_ecs_image

Layer 1 (depends on Layer 0)
├── alibabacloudstack_vpc_vswitch              -> depends on alibabacloudstack_vpc_vpc
├── alibabacloudstack_security_group       -> depends on alibabacloudstack_vpc_vpc
├── alibabacloudstack_vpc_routetable          -> depends on alibabacloudstack_vpc_vpc
├── alibabacloudstack_eip_address          -> no dependency (can parallel with Layer 0)
├── alibabacloudstack_oss_bucket_acl       -> depends on alibabacloudstack_oss_bucket
├── alibabacloudstack_oss_bucket_policy    -> depends on alibabacloudstack_oss_bucket
├── alibabacloudstack_dns_record           -> depends on alibabacloudstack_dns_domain
└── alibabacloudstack_ascm_user_group_role_binding -> depends on alibabacloudstack_ascm_user + alibabacloudstack_ascm_ram_policy

Layer 2 (depends on Layer 1)
├── alibabacloudstack_nat_gateway          -> depends on alibabacloudstack_vpc_vpc + alibabacloudstack_vpc_vswitch
├── alibabacloudstack_vpn_gateway          -> depends on alibabacloudstack_vpc_vpc
├── alibabacloudstack_slb_loadbalancer    -> depends on alibabacloudstack_vpc_vswitch (intranet)
├── alibabacloudstack_alb_load_balancer    -> depends on alibabacloudstack_vpc_vpc + alibabacloudstack_vpc_vswitch
├── alibabacloudstack_alb_server_group     -> depends on alibabacloudstack_vpc_vpc
├── alibabacloudstack_security_group_rule  -> depends on alibabacloudstack_security_group
└── alibabacloudstack_route_entry          -> depends on alibabacloudstack_vpc_routetable

Layer 3 (depends on Layer 2)
├── alibabacloudstack_instance             -> depends on alibabacloudstack_vpc_vswitch + alibabacloudstack_security_group
├── alibabacloudstack_db_instance          -> depends on alibabacloudstack_vpc_vswitch
├── alibabacloudstack_kvstore_instance     -> depends on alibabacloudstack_vpc_vswitch
├── alibabacloudstack_mongodb_instance     -> depends on alibabacloudstack_vpc_vswitch
├── alibabacloudstack_polardb_cluster      -> depends on alibabacloudstack_vpc_vswitch
├── alibabacloudstack_nas_file_system      -> no dependency (can parallel with Layer 0)
├── alibabacloudstack_snat_entry           -> depends on alibabacloudstack_nat_gateway
├── alibabacloudstack_forward_entry        -> depends on alibabacloudstack_nat_gateway
├── alibabacloudstack_eip_association      -> depends on alibabacloudstack_eip_address + instance
├── alibabacloudstack_slb_listener         -> depends on alibabacloudstack_slb_loadbalancer
└── alibabacloudstack_alb_listener         -> depends on alibabacloudstack_alb_load_balancer

Layer 4 (depends on Layer 3)
├── alibabacloudstack_disk                 -> no dependency (can parallel with Layer 0)
├── alibabacloudstack_disk_attachment      -> depends on alibabacloudstack_disk + alibabacloudstack_instance
├── alibabacloudstack_db_database          -> depends on alibabacloudstack_db_instance
├── alibabacloudstack_db_account           -> depends on alibabacloudstack_db_instance
├── alibabacloudstack_slb_backend_server       -> depends on alibabacloudstack_slb_loadbalancer + alibabacloudstack_instance
├── alibabacloudstack_alb_server_group     -> depends on alibabacloudstack_alb_load_balancer (server binding)
├── alibabacloudstack_nas_mount_target     -> depends on alibabacloudstack_nas_file_system + alibabacloudstack_vpc_vswitch
└── alibabacloudstack_polardb_database     -> depends on alibabacloudstack_polardb_cluster

Layer 5 (depends on Layer 4)
├── alibabacloudstack_ess_scaling_group    -> depends on alibabacloudstack_vpc_vswitch
├── alibabacloudstack_ack_cluster -> depends on alibabacloudstack_vpc_vpc + alibabacloudstack_vpc_vswitch
└── alibabacloudstack_fc_service           -> no dependency (can parallel with Layer 0)

Layer 6 (depends on Layer 5)
├── alibabacloudstack_ess_scaling_configuration -> depends on alibabacloudstack_ess_scaling_group
├── alibabacloudstack_ess_scaling_rule          -> depends on alibabacloudstack_ess_scaling_group
├── alibabacloudstack_cs_kubernetes_node_pool   -> depends on alibabacloudstack_ack_cluster
└── alibabacloudstack_fc_function               -> depends on alibabacloudstack_fc_service
```

---

## Recommended Import Order

Execute imports in the following order to avoid dependency errors:

```
Batch 1: VPC, OSS Bucket, DNS Domain, RAM resources, KMS keys
Batch 2: VSwitch, Security Group, Route Table, EIP
Batch 3: NAT Gateway, SLB, ALB, Security Group Rules
Batch 4: ECS Instance, RDS, Redis, MongoDB, NAS
Batch 5: Disk Attachment, EIP Association, SNAT/DNAT, SLB Listener, DB Accounts
Batch 6: Auto Scaling Groups, ACK Clusters
Batch 7: Scaling Configurations, Node Pools, FC Functions
```

---

## Key Dependency Attributes Per Resource

### alibabacloudstack_vpc_vswitch
```hcl
vpc_id    = alibabacloudstack_vpc_vpc.<name>.id      # Must reference, never hardcode
zone_id   = "cn-hangzhou-env01-d01"             # Fixed value, does not reference other resources
```

### alibabacloudstack_security_group
```hcl
vpc_id = alibabacloudstack_vpc_vpc.<name>.id
```

### alibabacloudstack_instance
```hcl
vswitch_id      = alibabacloudstack_vpc_vswitch.<name>.id
security_groups = [alibabacloudstack_security_group.<name>.id]
# Note: security_groups is a list, can reference multiple security groups
```

### alibabacloudstack_nat_gateway
```hcl
vpc_id     = alibabacloudstack_vpc_vpc.<name>.id
vswitch_id = alibabacloudstack_vpc_vswitch.<name>.id  # Required for Enhanced type
```

### alibabacloudstack_snat_entry
```hcl
snat_table_id     = split(",", alibabacloudstack_nat_gateway.<name>.snat_table_ids)[0]
source_vswitch_id = alibabacloudstack_vpc_vswitch.<name>.id
snat_ip           = alibabacloudstack_eip_address.<name>.ip_address
```

### alibabacloudstack_eip_association
```hcl
allocation_id = alibabacloudstack_eip_address.<name>.id
instance_id   = alibabacloudstack_instance.<name>.id
instance_type = "EcsInstance"  # Or "Nat", "SlbInstance", etc.
```

### alibabacloudstack_db_instance
```hcl
vswitch_id = alibabacloudstack_vpc_vswitch.<name>.id
# Note: RDS does not directly reference security groups; access is controlled via security_ips whitelist
```

### alibabacloudstack_slb_loadbalancer (intranet)
```hcl
vswitch_id   = alibabacloudstack_vpc_vswitch.<name>.id
address_type = "intranet"
```

### alibabacloudstack_slb_backend_server
```hcl
load_balancer_id = alibabacloudstack_slb_loadbalancer.<name>.id
instance_ids     = [alibabacloudstack_instance.<name>.id]
# Note: instance_ids is a list
```

---

## Circular Dependency Detection

The following combinations are prone to circular dependencies and require special attention:

1. **ECS + Security Group**: ECS references security group; security group rules reference ECS private IP
   - Solution: Use CIDR in security group rules instead of instance IP, or use security group ID cross-referencing

2. **SLB + ECS**: SLB backend references ECS; ECS security group rules reference SLB
   - Solution: Use SLB IP range in security group rules, don't directly reference SLB resource

3. **NAT + EIP**: NAT Gateway needs EIP; EIP association needs NAT Gateway ID
   - Solution: Import NAT Gateway first, then import EIP association

---

## Cross-VPC Resource References

When resources span VPCs (e.g., VPC Peering), additional attention is needed:

```hcl
# VPC Peering connection
resource "alibabacloudstack_vpc_peer_connection" "peer_<name>" {
  peer_connection_name = "<name>"
  vpc_id               = alibabacloudstack_vpc_vpc.vpc_a.id
  accepting_ali_uid    = "<peer-account-id>"
  accepting_region_id  = "cn-beijing"
  accepting_vpc_id     = "<peer-vpc-id>"
}
```

---

## Multi-Region Resources

When managing resources across multiple Regions, use provider alias:

```hcl
provider "alibabacloudstack" {
  alias      = "beijing"
  region     = "cn-beijing"
  access_key = var.access_key
  secret_key = var.secret_key
}

resource "alibabacloudstack_vpc_vpc" "vpc_beijing" {
  provider   = alibabacloudstack.beijing
  vpc_name   = "prod-beijing"
  cidr_block = "172.16.0.0/12"
}
```
