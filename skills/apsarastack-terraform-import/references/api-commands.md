# apsara-cli Command Quick Reference

Resource discovery commands organized by product, including pagination parameters and output formatting.

**Notes:**
- Commands marked with a checkmark have been verified with a test account and can be used directly
- Unmarked commands are generated based on training knowledge; verify parameters with `apsara <product> help` before executing
- For resource types not in this file, Claude will look up apsara-cli or provider documentation on the fly during the discovery phase
- **All commands use plugin mode (lowercase-hyphenated format)**; parameter names do not use CamelCase

---

## General Parameter Notes

### ApsaraStack Connection Parameters (Required)

```bash
# Environment variables
export ALIBABACLOUDSTACK_ACCESS_KEY="your-ak"
export ALIBABACLOUDSTACK_SECRET_KEY="your-sk"
export ALIBABACLOUDSTACK_REGION="<REGION>"
export ALIBABACLOUDSTACK_POPGW_DOMAIN="<POPGW_DOMAIN>"
export HTTP_PROXY="http://<PROXY_HOST>:<PROXY_PORT>"   # If behind proxy

# Full command template for ApsaraStack
apsara <product> <Action> \
  --access-key-id "$ALIBABACLOUDSTACK_ACCESS_KEY" \
  --access-key-secret "$ALIBABACLOUDSTACK_SECRET_KEY" \
  --region "$ALIBABACLOUDSTACK_REGION" \
  --endpoint "<product>.$ALIBABACLOUDSTACK_POPGW_DOMAIN" \
  --skip-secure-verify

# Example: VPC discovery on 17E
apsara vpc describe-vpcs \
  --access-key-id "$ALIBABACLOUDSTACK_ACCESS_KEY" \
  --access-key-secret "$ALIBABACLOUDSTACK_SECRET_KEY" \
  --region "$ALIBABACLOUDSTACK_REGION" \
  --endpoint "vpc.$ALIBABACLOUDSTACK_POPGW_DOMAIN" \
  --skip-secure-verify
```

| Parameter | Description | Required |
|-----------|-------------|----------|
| `--endpoint` | POP Gateway endpoint: `<product>.<popgw_domain>` | Yes |
| `--skip-secure-verify` | Skip TLS certificate verification (self-signed certs) | Yes |
| `--access-key-id` | Access Key ID | Yes |
| `--access-key-secret` | Access Key Secret | Yes |
| `--region` | ApsaraStack region ID | Yes |
| `HTTP_PROXY` | HTTP proxy for intranet access | Env-dependent |

### Common Parameters

```bash
# Set Region (common to all commands)
REGION="cn-hangzhou-env01-d01"

# Region parameters (two types, different meanings):
# --biz-region-id  Business parameter, corresponds to the API's RegionId, passed in request body. Use when API requires RegionId
# --region         Framework parameter, controls which region's endpoint to connect to, not passed to API request body
# Rule: Prefer --biz-region-id (covers all scenarios).
#       --region can substitute only when the API's RegionId is optional (e.g., describe-vswitches, where setting endpoint region causes API to default query in that region)

# Pagination parameters
--page-size 50       # Items per page (max 100)
--page-number 1      # Page number (starts from 1)

# Output formatting (table format)
--output cols=Field1,Field2 rows=Path.To.Array

# JSON format output (for parsing details)
# Default output is JSON, no additional parameters needed
```

---

## Networking Products (VPC)

```bash
# Verified: VPC list (uses --biz-region-id)
apsara vpc describe-vpcs --biz-region-id $REGION --page-size 50 \
  --output cols=VpcId,VpcName,CidrBlock,Status rows=Vpcs.Vpc

# Verified: VPC details (JSON)
apsara vpc describe-vpcs --biz-region-id $REGION --vpc-id vpc-bp1xxx

# Verified: VSwitch list (uses --region)
apsara vpc describe-vswitches --region $REGION --page-size 50 \
  --output cols=VSwitchId,VSwitchName,CidrBlock,ZoneId,VpcId rows=VSwitches.VSwitch

# Verified: VSwitch details
apsara vpc describe-vswitches --region $REGION --vswitch-id vsw-bp1xxx

# Route table list
apsara vpc describe-route-tables --biz-region-id $REGION --page-size 50 \
  --output cols=RouteTableId,RouteTableName,VpcId,RouteTableType rows=RouteTables.RouteTable

# EIP list (uses --biz-region-id)
apsara vpc describe-eip-addresses --biz-region-id $REGION --page-size 50 \
  --output cols=AllocationId,IpAddress,Status,InstanceId,InstanceType rows=EipAddresses.EipAddress

# EIP details
apsara vpc describe-eip-addresses --biz-region-id $REGION --allocation-id eip-bp1xxx

# Verified: NAT Gateway list (uses --biz-region-id)
apsara vpc describe-nat-gateways --biz-region-id $REGION --page-size 50 \
  --output cols=NatGatewayId,Name,VpcId,Status,NatType rows=NatGateways.NatGateway

# NAT Gateway details
apsara vpc describe-nat-gateways --biz-region-id $REGION --nat-gateway-id ngw-bp1xxx

# SNAT entry list
apsara vpc describe-snat-table-entries --biz-region-id $REGION \
  --snat-table-id stb-bp1xxx --page-size 50 \
  --output cols=SnatEntryId,SnatIp,SourceCidr,Status rows=SnatTableEntries.SnatTableEntry

# VPN Gateway list
apsara vpc describe-vpn-gateways --biz-region-id $REGION --page-size 50 \
  --output cols=VpnGatewayId,Name,VpcId,Status rows=VpnGateways.VpnGateway
```

---

## Compute Products (ECS)

```bash
# Verified: ECS instance list (uses --biz-region-id)
apsara ecs describe-instances --biz-region-id $REGION --page-size 50 \
  --output cols=InstanceId,InstanceName,Status,InstanceType,ZoneId,VpcAttributes.VpcId rows=Instances.Instance

# Verified: ECS instance details (JSON, includes all attributes)
apsara ecs describe-instances --biz-region-id $REGION \
  --instance-ids '["i-bp1xxx"]'

# Verified: Security group list (uses --biz-region-id)
apsara ecs describe-security-groups --biz-region-id $REGION --page-size 50 \
  --output cols=SecurityGroupId,SecurityGroupName,VpcId,SecurityGroupType rows=SecurityGroups.SecurityGroup

# Security group details (with rules)
apsara ecs describe-security-group-attribute --biz-region-id $REGION \
  --security-group-id sg-bp1xxx --direction ingress
apsara ecs describe-security-group-attribute --biz-region-id $REGION \
  --security-group-id sg-bp1xxx --direction egress

# Verified: Disk list (filter by instance)
apsara ecs describe-disks --biz-region-id $REGION --instance-id i-bp1xxx --page-size 50 \
  --output cols=DiskId,DiskName,Status,Size,Type,Category rows=Disks.Disk

# Disk details
apsara ecs describe-disks --biz-region-id $REGION --disk-ids '["d-bp1xxx"]'

# Image list (custom images)
apsara ecs describe-images --biz-region-id $REGION --image-owner-alias self --page-size 50 \
  --output cols=ImageId,ImageName,Status,OSType rows=Images.Image

# Key pair list
apsara ecs describe-key-pairs --biz-region-id $REGION --page-size 50 \
  --output cols=KeyPairName,KeyPairFingerPrint rows=KeyPairs.KeyPair

# Auto Scaling group list
apsara ess describe-scaling-groups --biz-region-id $REGION --page-size 50 \
  --output cols=ScalingGroupId,ScalingGroupName,LifecycleState rows=ScalingGroups.ScalingGroup

# Scaling configuration list
apsara ess describe-scaling-configurations --biz-region-id $REGION \
  --scaling-group-id asg-bp1xxx --page-size 50 \
  --output cols=ScalingConfigurationId,ScalingConfigurationName,LifecycleState rows=ScalingConfigurations.ScalingConfiguration
```

---

## Storage Products (OSS)

```bash
# OSS Bucket list (global, no Region needed)
apsara ossutil ls oss://

# OSS Bucket details
apsara ossutil stat oss://bucket-name

# OSS Bucket ACL
apsara ossutil bucket-acl oss://bucket-name

# OSS Bucket versioning
apsara ossutil bucket-versioning oss://bucket-name

# OSS Bucket lifecycle
apsara ossutil lifecycle --method get oss://bucket-name

# OSS Bucket CORS
apsara ossutil cors --method get oss://bucket-name

# OSS Bucket encryption
apsara ossutil bucket-encryption --method get oss://bucket-name

# NAS file system list
apsara nas describe-file-systems --biz-region-id $REGION --page-size 50 \
  --output cols=FileSystemId,FileSystemType,Status,StorageType rows=FileSystems.FileSystem

# NAS mount target list
apsara nas describe-mount-targets --biz-region-id $REGION \
  --file-system-id fs-xxx --page-size 50 \
  --output cols=MountTargetDomain,Status,VpcId,VSwitchId rows=MountTargets.MountTarget
```

---

## Database Products (RDS)

```bash
# RDS instance list (uses --biz-region-id)
apsara rds describe-db-instances --biz-region-id $REGION --page-size 50 \
  --output cols=DBInstanceId,DBInstanceDescription,Engine,EngineVersion,DBInstanceStatus,DBInstanceClass rows=Items.DBInstance

# RDS instance details
apsara rds describe-db-instance-attribute --biz-region-id $REGION \
  --db-instance-id rm-bp1xxx

# RDS database list
apsara rds describe-databases --biz-region-id $REGION \
  --db-instance-id rm-bp1xxx \
  --output cols=DBName,DBStatus,CharacterSetName rows=Databases.Database

# RDS account list
apsara rds describe-accounts --biz-region-id $REGION \
  --db-instance-id rm-bp1xxx \
  --output cols=AccountName,AccountStatus,AccountType rows=Accounts.DBInstanceAccount

# RDS IP whitelist
apsara rds describe-db-instance-ip-array-list --biz-region-id $REGION \
  --db-instance-id rm-bp1xxx

# RDS connection addresses
apsara rds describe-db-instance-net-info --biz-region-id $REGION \
  --db-instance-id rm-bp1xxx
```

---

## Cache Products (Redis/KVStore)

```bash
# Redis instance list (uses --biz-region-id)
apsara r-kvstore describe-instances --biz-region-id $REGION --page-size 50 \
  --output cols=InstanceId,InstanceName,InstanceStatus,InstanceType,EngineVersion rows=Instances.KVStoreInstance

# Redis instance details
apsara r-kvstore describe-instance-attribute --biz-region-id $REGION \
  --instance-id r-bp1xxx

# Redis IP whitelist
apsara r-kvstore describe-security-ips --biz-region-id $REGION \
  --instance-id r-bp1xxx
```

---

## Document Database (MongoDB/DDS)

```bash
# MongoDB instance list (uses --biz-region-id)
apsara dds describe-db-instances --biz-region-id $REGION --page-size 50 \
  --output cols=DBInstanceId,DBInstanceDescription,DBInstanceStatus,DBInstanceType rows=DBInstances.DBInstance

# MongoDB instance details
apsara dds describe-db-instance-attribute --biz-region-id $REGION \
  --db-instance-id dds-bp1xxx

# MongoDB IP whitelist
apsara dds describe-security-ips --biz-region-id $REGION \
  --db-instance-id dds-bp1xxx
```

---

## Load Balancing (SLB)

```bash
# SLB instance list (uses --biz-region-id)
apsara slb describe-load-balancers --biz-region-id $REGION --page-size 50 \
  --output cols=LoadBalancerId,LoadBalancerName,LoadBalancerStatus,AddressType,VpcId rows=LoadBalancers.LoadBalancer

# SLB instance details
apsara slb describe-load-balancer-attribute --biz-region-id $REGION \
  --load-balancer-id lb-bp1xxx

# SLB listener list
apsara slb describe-load-balancer-listeners --biz-region-id $REGION \
  --load-balancer-id lb-bp1xxx \
  --output cols=ListenerPort,ListenerProtocol,Status rows=Listeners

# SLB backend servers
apsara slb describe-health-status --biz-region-id $REGION \
  --load-balancer-id lb-bp1xxx

# ALB instance list (uses --region)
apsara alb list-load-balancers --region $REGION --max-results 50

# ALB listener list
apsara alb list-listeners --region $REGION --max-results 50

# ALB server group list
apsara alb list-server-groups --region $REGION --max-results 50
```

---

## DNS (Alidns)

```bash
# Domain list (global, no Region needed)
apsara alidns describe-domains --page-size 50 \
  --output cols=DomainId,DomainName,RecordCount,GroupName rows=Domains.Domain

# Domain DNS record list
apsara alidns describe-domain-records --domain-name example.com --page-size 500 \
  --output cols=RecordId,RR,Type,Value,TTL,Status rows=DomainRecords.Record

# DNS group list
apsara alidns describe-domain-groups --page-size 50 \
  --output cols=GroupId,GroupName rows=DomainGroups.DomainGroup
```

---

## Security and Identity (RAM)

```bash
# RAM user list
apsara ram list-users --output cols=UserName,DisplayName,CreateDate rows=Users.User

# RAM role list
apsara ram list-roles --output cols=RoleName,RoleId,CreateDate rows=Roles.Role

# RAM policy list (custom)
apsara ram list-policies --policy-type Custom \
  --output cols=PolicyName,PolicyType,CreateDate rows=Policies.Policy

# RAM user attached policies
apsara ram list-policies-for-user --user-name username \
  --output cols=PolicyName,PolicyType rows=Policies.Policy

# KMS key list
apsara kms list-keys --page-size 100 \
  --output cols=KeyId,KeyArn rows=Keys.Key

# KMS key details
apsara kms describe-key --key-id key-id
```

---

## Container Service (ACK)

```bash
# ACK cluster list
apsara cs GET /clusters 2>&1

# ACK node pool list
apsara cs GET /clusters/<cluster-id>/nodepools 2>&1
```

---

## Function Compute (FC)

```bash
# FC service list
apsara fc GET /services 2>&1

# FC function list
apsara fc GET /services/<service-name>/functions 2>&1
```

---

## FAQ

**Q: Command returns empty results**
- Check if Region is correct
- Check if account has Describe permission for the corresponding resource
- Some resources (e.g., OSS) are global and do not need a Region parameter

**Q: Region parameter error**
- `--biz-region-id` is a business parameter (corresponds to the API's RegionId), `--region` is a framework parameter (controls endpoint region)
- Prefer `--biz-region-id`; `--region` can substitute only when the API's RegionId is optional
- If a command reports `--biz-region-id is required`, the API's RegionId is mandatory and `--region` alone is not sufficient
- Global resources (OSS, DNS, RAM) do not need a Region parameter

**Q: Pagination issues (more than 50 resources)**
```bash
# Use page number pagination
apsara vpc describe-vpcs --biz-region-id $REGION --page-size 50 --page-number 2

# Or use --pager parameter (apsara-cli supports automatic page merging)
apsara ecs describe-instances --biz-region-id $REGION --page-size 100 --pager
```

**Q: Output formatting not working**
- Confirm the `rows=` path is correct (try without `--output` first to see raw JSON structure)
- Use `.` for nested fields (e.g., `VpcAttributes.VpcId`)

**Q: Unsure about parameter names**
- Use `apsara <product> <command> --help` to see all supported parameters
- Plugin mode parameters are always lowercase-hyphenated (e.g., `--instance-id`, `--vpc-id`)
