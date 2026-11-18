# ResourceCenter API Guide

> **ApsaraStack Note**: ResourceCenter and IaCService APIs are **public cloud services** and may not be available in ApsaraStack environments. This reference is retained for compatibility; in ApsaraStack, always fall back to native product APIs in `references/api-commands.md` for resource discovery.

Four APIs form the complete resource discovery and relationship building pipeline, with native product APIs as fallback.

API documentation entry: `https://api.aliyun.com/api/ResourceCenter/2022-12-01/`

---

## Prerequisite API: IaCService ListResourceTypes

**Called once during Phase 3 initialization to build a `terraformResourceType <-> ResourceCode` mapping table, cached to `.import/resource-type-mapping.json` for use throughout the session.**

```bash
apsara iacservice list-resource-types \
  --endpoint iacservice.aliyuncs.com \
  --status Available \
  --maxResults 200
```

If results exceed 200 entries, use `nextToken` to paginate and retrieve all data.

**Key return fields:**

```json
{
  "resourceTypes": [
    {
      "terraformResourceType": "alibabacloudstack_eip_address",
      "resourceType": "ALIYUN::EIP::Address",
      "product": "EIP",
      "supportTerraformer": true,
      "status": "Available"
    }
  ]
}
```

**Mapping rules:**

Replace the `ALIYUN` prefix in the `resourceType` field with `ACS` to get the ResourceCode required by ResourceCenter APIs:

```
ALIYUN::EIP::Address  ->  ACS::EIP::Address
ALIYUN::VPC::VPC      ->  ACS::VPC::VPC
ALIYUN::ECS::Instance ->  ACS::ECS::Instance
```

**Mapping table structure (`.import/resource-type-mapping.json`):**

```json
{
  "byTerraformType": {
    "alibabacloudstack_eip_address": {
      "resourceCode": "ACS::EIP::Address",
      "product": "EIP",
      "supportTerraformer": true
    },
    "alibabacloudstack_vpc_vpc": {
      "resourceCode": "ACS::VPC::VPC",
      "product": "VPC",
      "supportTerraformer": true
    }
  },
  "byResourceCode": {
    "ACS::EIP::Address": "alibabacloudstack_eip_address",
    "ACS::VPC::VPC": "alibabacloudstack_vpc_vpc"
  }
}
```

**Usage:**
- Phase 4: Known `terraformResourceType` -> Look up `byTerraformType` -> Get `resourceCode` -> Call search-resources
- Phase 6: ResourceCenter returns `ResourceType` -> Look up `byResourceCode` -> Get `terraformResourceType` -> Generate correct HCL resource type

**Fallback conditions:**
- Call fails or insufficient permissions -> Fall back to the static mapping table at the end of this file
- Mapping table has no corresponding entry -> Skip API A/B, use native product API directly

---

## API One: list-resource-types

Query resource types supported by ResourceCenter, used to determine the support scope of API A (search-resources) and API B (list-resource-relationships).

**Called once before Phase 4 begins, results cached for use throughout the session.**

```bash
apsara resourcecenter list-resource-types --endpoint resourcecenter.aliyuncs.com
```

**Key return fields:**

```json
{
  "ResourceTypes": [
    {
      "ResourceType": "ACS::ECS::Instance",
      "SupportedFeatures": ["ListResources", "GetResourceRelationships"]
    }
  ]
}
```

- `SupportedFeatures` contains `ListResources` -> This resource type supports API A (search-resources)
- `SupportedFeatures` contains `GetResourceRelationships` -> This resource type supports API B (list-resource-relationships)

**Fallback conditions:**
- Call returns permission error (`AccessDenied`) -> Fall back for the entire session, skip API A and API B
- Call times out or returns empty list -> Same as above

---

## API A: search-resources

Search resources accessible by the current account. **Returns basic resource metadata, not business attributes** (such as VPC CidrBlock, ECS InstanceType, etc.). Phase 6 HCL generation still requires native product API calls for complete attributes; search-resources is only used during Phase 4 resource discovery.

**Basic usage:**
```bash
apsara resourcecenter search-resources \
  --filter '[{"Key":"ResourceType","MatchType":"Equals","Value":["ACS::ECS::Instance"]}]' \
  --max-results 50
```

**Supported filter conditions:**

| Filter Parameter | Supported Match Types |
|---------|-------------|
| `ResourceType` | Equals |
| `RegionId` | Equals |
| `ResourceId` | Equals, Prefix |
| `ResourceName` | Equals, Contains |
| `Tag` | Contains, NotContains, NotExists (JSON format: `{"key":"k","value":"v"}`) |
| `VpcId` | Equals |
| `VSwitchId` | Equals |
| `IpAddress` | Equals, Contains |
| `ResourceGroupId` | Equals, Exists, NotExists |

Multiple filter conditions have AND relationship; multiple values within the same condition have OR relationship.

**Filter by Region + resource type:**
```bash
apsara resourcecenter search-resources \
  --filter '[{"Key":"ResourceType","MatchType":"Equals","Value":["ACS::ECS::Instance"]},{"Key":"RegionId","MatchType":"Equals","Value":["cn-hangzhou-env01-d01"]}]' \
  --max-results 50
```

**Query by resource ID (to confirm resource type):**
```bash
apsara resourcecenter search-resources \
  --filter '[{"Key":"ResourceId","MatchType":"Equals","Value":["vpc-bp1xxx"]}]' \
  --max-results 1
```

**Pagination:**
```bash
apsara resourcecenter search-resources \
  --filter '[...]' \
  --max-results 50 \
  --next-token <NextToken from previous page>
```
No `NextToken` in results indicates no more data. MaxResults maximum is 500.

**Return fields:**

```json
{
  "Resources": [
    {
      "AccountId": "151266687691****",
      "ResourceId": "vtb-bp1xxx",
      "ResourceType": "ACS::VPC::RouteTable",
      "ResourceName": "group1",
      "RegionId": "cn-hangzhou-env01-d01",
      "ZoneId": "cn-hangzhou-env01-d01",
      "ResourceGroupId": "rg-acfmzawhxxc****",
      "CreateTime": "2021-06-30T09:20:08Z",
      "ExpireTime": "2021-07-30T09:20:08Z",
      "IpAddresses": ["192.168.1.2"],
      "IpAddressAttributes": [
        {"IpAddress": "192.168.1.2", "NetworkType": "Public", "Version": "Ipv4"}
      ],
      "Tags": [{"Key": "test_key", "Value": "test_value"}],
      "Deleted": false
    }
  ],
  "NextToken": "eyJzZWFyY2hBZnRlcnMiOlsiMTAwMTU2Nzk4MTU1OSJd****"
}
```

Note: Whether fields like `ZoneId`, `CreateTime`, `IpAddresses` are returned depends on each cloud service; not guaranteed for all resource types.

**Resource type format (`ACS::<Product>::<ResourceType>`):**

| Cloud Product | ResourceCenter Type |
|--------|-------------------|
| VPC | `ACS::VPC::VPC` |
| VSwitch | `ACS::VPC::VSwitch` |
| Route Table | `ACS::VPC::RouteTable` |
| EIP | `ACS::VPC::EipAddress` |

Complete list obtained dynamically via `iacservice list-resource-types`. Static fallback mapping table below:

| Terraform Resource Type | ResourceCode |
|-------------------|-------------|
| `alibabacloudstack_vpc_vpc` | `ACS::VPC::VPC` |
| `alibabacloudstack_vpc_vswitch` | `ACS::VPC::VSwitch` |
| `alibabacloudstack_vpc_routetable` | `ACS::VPC::RouteTable` |
| `alibabacloudstack_eip_address` | `ACS::EIP::Address` |
| `alibabacloudstack_nat_gateway` | `ACS::VPC::NatGateway` |
| `alibabacloudstack_instance` | `ACS::ECS::Instance` |
| `alibabacloudstack_security_group` | `ACS::ECS::SecurityGroup` |
| `alibabacloudstack_db_instance` | `ACS::RDS::DBInstance` |
| `alibabacloudstack_kvstore_instance` | `ACS::KVStore::Instance` |
| `alibabacloudstack_mongodb_instance` | `ACS::DDS::DBInstance` |
| `alibabacloudstack_oss_bucket` | `ACS::OSS::Bucket` |
| `alibabacloudstack_slb_loadbalancer` | `ACS::SLB::LoadBalancer` |
| `alibabacloudstack_ack_cluster` | `ACS::CS::KubernetesCluster` |

**Fallback conditions:**
- Resource type is not in the list returned by `list-resource-types`
- Call returns `NoPermission` or other errors
- In these cases, fall back to native product APIs in `references/api-commands.md`

---

## API B: list-resource-relationships

Query relationships between a specified resource and its related resources, used to dynamically build the resource dependency graph. **`RegionId`, `ResourceType`, and `ResourceId` are all required parameters.**

Results are a flat list where each record represents an association between "current resource" and "a related resource". **Direction is not distinguished** (no parent-child or lateral distinction); the caller determines dependency direction based on business semantics.

**Basic usage:**
```bash
apsara resourcecenter list-resource-relationships \
  --region-id cn-hangzhou \
  --resource-type ACS::VPC::VPC \
  --resource-id vpc-bp1xxx
```

**Filter by related resource type (only VSwitch):**
```bash
apsara resourcecenter list-resource-relationships \
  --region-id cn-hangzhou \
  --resource-type ACS::VPC::VPC \
  --resource-id vpc-bp1xxx \
  --related-resource-filter '[{"Key":"RelatedResourceType","MatchType":"Equals","Value":["ACS::VPC::VSwitch"]}]'
```

**Pagination:**
```bash
apsara resourcecenter list-resource-relationships \
  --region-id cn-hangzhou \
  --resource-type ACS::VPC::VPC \
  --resource-id vpc-bp1xxx \
  --max-results 50 \
  --next-token <NextToken from previous page>
```
MaxResults range is 1-500, default 20.

**Supported RelatedResourceFilter parameters:**

| Parameter | Description | Supported Match Types |
|------|------|-------------|
| `RelatedResourceRegionId` | Related resource region ID | Equals |
| `RelatedResourceType` | Related resource type | Equals |
| `RelatedResourceId` | Related resource ID | Equals |

**Return fields:**

```json
{
  "ResourceRelationships": [
    {
      "RegionId": "cn-hangzhou-env01-d01",
      "ResourceType": "ACS::ACK::Cluster",
      "ResourceId": "m-eb3hji****",
      "RelatedResourceRegionId": "cn-shanghai",
      "RelatedResourceType": "ACS::VPC::VPC",
      "RelatedResourceId": "vpc-uf6m5okksddm6c9lh7***"
    }
  ],
  "NextToken": "eyJzZWFyY2hBZnRlcnMiOlsiMTAwMTU2Nzk4MTU1OSJd****"
}
```

Note: Related resources may be in different Regions (`RelatedResourceRegionId` differs from the requested `RegionId`); be aware of cross-Region associations.

**Phase 5 usage:**

1. From Phase 4 discovered resources, identify root resources (VPC, OSS Bucket, DNS domains, RAM users, etc. with no parent dependencies)
2. For each root resource, call API B to get all related resources
3. Recursively call API B for related resources until no new resources appear, building the complete relationship graph
4. Dependency direction is determined by business semantics: VPC -> VSwitch (VSwitch depends on VPC), VSwitch -> ECS (ECS depends on VSwitch), etc. Refer to `references/dependency-rules.md`
5. Topologically sort the relationship graph to generate import order

**Phase 6 usage:**

For related resources discovered via API B, their dependency relationships are already determined. Combined with the reference field whitelist extracted in Phase 3, find the corresponding `_id` fields in the current resource and directly replace with terraform resource address, without guessing.

**Fallback conditions:**
- Resource type is not in the list returned by `list-resource-types`
- Call returns `NoPermission` or other errors
- In these cases, fall back to static rules in `references/dependency-rules.md`
