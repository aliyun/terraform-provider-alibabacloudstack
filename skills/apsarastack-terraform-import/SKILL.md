---
name: apsarastack-terraform-import
description: This skill should be used when the user asks to "导入专有云资源到 Terraform", "terraform import ApsaraStack", "terraform import 专有云", "将现有专有云资源纳入 Terraform 管理", "阿里云专有云资源迁移 Terraform", "生成 terraform state", "import alibabacloudstack resources", "阿里云专有云 IaC 迁移", "阿里云专有云 Terraform 导入", or needs to manage existing Alibaba Cloud ApsaraStack resources with Terraform. Guides users step-by-step through environment check, authentication, resource discovery, HCL generation, state import, validation, and dependency graph. Supports both one-time migration and incremental sync.
version: 0.1.0
---

# Alibaba Cloud ApsaraStack Resource Terraform Import Wizard

Step-by-step guide to import existing resources under an Alibaba Cloud ApsaraStack account into local Terraform management. Automatically performs environment checks, resource discovery, HCL generation, state import, etc. Users only need to confirm at key decision points.

Two modes are supported:
- **One-time migration**: Complete Phase 1-9 end to end
- **Incremental sync**: When a Terraform working directory already exists, jump directly to incremental sync mode

## Resource Support Scope

- **Theoretical scope**: Any resource in the alibabacloudstack provider that has a Read method (i.e., can read state from the cloud). Most resources can be imported directly via the `terraform import` command or `import` blocks (Terraform 1.5+); a few resources without an Importer can be imported by manually constructing tfstate, but this requires exact matching of the provider's state structure and carries some risk. Check support status at `https://registry.terraform.io/providers/aliyun/alibabacloudstack/latest/docs`.
- **Pre-built resources**: ~50 common resource types (networking, compute, storage, databases, load balancing, DNS, RAM, ACK, FC) with verified API discovery commands, HCL templates, and import ID formats in `references/`
- **Beyond pre-built scope**: Will look up provider docs on the fly to provide commands and templates; accuracy depends on documentation content; users should verify before executing

---

## Command Format Requirements

**All apsara-cli commands must use plugin mode (lowercase-hyphenated format). CamelCase command names and parameter names are prohibited.** Violating this rule will cause parameter errors or call failures.

**Example (correct plugin-mode format):**
```bash
apsara vpc describe-vpcs --biz-region-id cn-shanghai --page-size 50
```

**Parameter conversion rules:**
- **Region**: Prefer `--biz-region-id` (business parameter, corresponds to the API's RegionId); `--region` (framework parameter, controls endpoint region) should only be used when the API's RegionId is optional
- **Pagination**: `--page-size` (not PageSize)
- **Others**: Convert CamelCase parameter names to lowercase-hyphenated (VpcId -> `--vpc-id`, InstanceIds -> `--instance-ids`). Use `apsara <product> <command> --help` to confirm
- **OSS**: Use `apsara ossutil ls oss://` exclusively

---

## Phase 1: Environment Check

Run the following commands to detect the environment:

```bash
echo "=== apsara-cli ===" && apsara --version 2>&1 || echo "NOT_INSTALLED"
echo "=== terraform ===" && terraform --version 2>&1 || echo "NOT_INSTALLED"
echo "=== python3 ===" && python3 --version 2>&1 || echo "NOT_INSTALLED"
echo "=== OS ===" && uname -s 2>&1
```

**Tool requirements by tier:**

| Tool | Required? | When needed |
|------|-----------|-------------|
| terraform | **Mandatory** | Always (init, refresh, import, plan) |
| python3 | **Mandatory** | JSON parsing, report generation |
| apsara-cli | *Conditional* | Only when discovering resources without datasource coverage (see `references/unsupported-datasource.md`) |

Automatically handle missing tools based on output:
- **terraform missing** (macOS): Run `brew tap hashicorp/tap && brew install hashicorp/tap/terraform`
- **terraform missing** (Linux): Direct user to https://developer.hashicorp.com/terraform/install for manual installation (requires root)
- **apsara-cli missing**: Note it in the environment report. Do **not** block here — defer installation to Phase 4 when (and only when) a Group B resource is encountered. This avoids unnecessary installation if all target resources have datasource coverage.

Once terraform and python3 are ready, proceed to Phase 2.

**ApsaraStack Environment Notes:**

- **Region format**: ApsaraStack regions typically follow the format `cn-hangzhou-env01-d01` (not `cn-hangzhou` like public cloud)
- **Endpoint discovery**: If `apsara` commands fail with endpoint errors, the user may need to configure endpoints via `apsara configure --profile <name>` or environment variables
- **ASAPI endpoint**: Some ApsaraStack environments require ASAPI endpoint configuration for certain product APIs

---

## Phase 2: Authentication Configuration

This Skill relies on the ApsaraStack default credential chain for authentication, supporting the following methods (by priority):
1. Environment variables (ALIBABACLOUDSTACK_ACCESS_KEY / ALIBABACLOUDSTACK_SECRET_KEY) — **highest priority**
2. apsara-cli configuration file (~/.aliyun/config.json)
3. ECS RAM role (only when running on ECS)

**ApsaraStack-specific credential variables:**
- `ALIBABACLOUDSTACK_ACCESS_KEY` / `ALIBABACLOUDSTACK_SECRET_KEY` — Required
- `ALIBABACLOUDSTACK_REGION` — Required (e.g., `cn-hangzhou-env01-d01`)
- `ALIBABACLOUDSTACK_POPGW_DOMAIN` — POP gateway domain, often required for ApsaraStack
- `ALIBABACLOUDSTACK_PROTOCOL` — Optional, defaults to `HTTP`
- `ALIBABACLOUDSTACK_INSECURE` — Optional, set to `true` to skip TLS verification
- `ALIBABACLOUDSTACK_DEPARTMENT` — Optional, organization department
- `ALIBABACLOUDSTACK_RESOURCE_GROUP` — Optional, resource group name

**Important**: The Terraform alibabacloudstack provider and apsara-cli each read credentials independently, both prioritizing environment variables. When environment variables and CLI profile point to different accounts, environment variables take precedence.

**Verify authentication:**

Verify identity via the default credential chain (without explicitly passing AK/SK):

```bash
apsara sts get-caller-identity 2>&1
```

Analyze the output:
- Confirm AccountId and UserId, record Region
- **If the credential source is unexpected** (e.g., environment variables and CLI profile point to different accounts), inform the user of the credential conflict and suggest unifying credential configuration before retrying

If an error occurs, prompt the user to:
- Check whether credentials are configured (via `apsara configure` or environment variables)
- Check whether `popgw_domain` is correctly set (ApsaraStack environments often require this)
- Check permissions (requires ReadOnlyAccess or equivalent read-only policy; policy name may vary in ApsaraStack)
- Refer to `references/ram-policies.md` for required permissions
---

## Phase 3: Working Directory Initialization

Ask the user for the following information:
- Target working directory path (default `~/apsarastack-terraform`)
- Target Region (confirm if already inferred from Phase 2, otherwise ask)
- Whether a Terraform working directory already exists (skip this Phase if so)

After user confirmation, automatically complete the following:

1. Create working directory and `.import/` subdirectory (for storing intermediate artifacts)
2. Generate `provider.tf` in the target directory (refer to `examples/provider.tf`), replacing the region default with the user-confirmed Region. The provider block must include ApsaraStack-specific configuration (`popgw_domain`, `protocol`, `insecure` as needed)
3. Run `terraform init` in the target directory
4. Analyze init output, confirm provider download succeeded
5. **(Optional enhancement)** Run `terraform providers schema -json`, extract fields ending in `_id` or `_ids` from all resources to build a **reference field whitelist** for Phase 6. If schema output is too large causing parse timeout, skip this step; Phase 6 will fall back to static dependency rules for reference resolution
6. **(Optional enhancement)** Call `iacservice list-resource-types` (paginate to get all), replace the `ALIYUN` prefix in the `resourceType` field with `ACS` to build a `terraformResourceType <-> ResourceCode` bidirectional mapping table, written to `.import/resource-type-mapping.json` for use in Phase 4/5/6. If the call fails or insufficient permissions, fall back to the static mapping table in `references/resource-hub-api.md`. **Skipping this step does not affect the core import flow**, only affects ResourceCenter API resource type conversion
7. Inform user that initialization is complete, proceed to Phase 4

---

## Phase 4: Resource Discovery

**Discovery Strategy: Datasource First, CLI Fallback**

All resource discovery follows a **two-tier strategy**:
1. **Tier 1 — Terraform Datasource** (preferred): Reuses provider auth config, returns Terraform-native JSON, no extra tools needed
2. **Tier 2 — apsara-cli** (fallback): Only for resources without datasource coverage. See `references/unsupported-datasource.md` for the full gap list.

**Step 0 — Classify resource types**

For each target resource type the user wants to discover, determine whether a datasource exists:
- Check `references/unsupported-datasource.md`
- Or check the provider docs for a corresponding `data source` page

Split into two groups:
- **Group A (datasource available)**: VPC, VSwitch, ECS, Security Group, EIP, NAT Gateway, SLB, OSS Bucket (list), RDS Instance, Redis, MongoDB, DNS, ASCM Users/Roles, etc.
- **Group B (no datasource)**: FC Service/Function, ALB, ACK Node Pool, RDS Database/Account, OSS Bucket properties (ACL/lifecycle/etc.)

**Step 1 — Discover Group A via Terraform Datasource**

1. Create `.import/discovery/discovery.tf` based on `examples/discovery.tf`, uncommenting only the datasource blocks for Group A resource types
2. Run the following commands:
   ```bash
   cd .import/discovery
   terraform init
   terraform refresh
   terraform output -json > ../discovery-output.json
   cd ../..
   ```
3. Parse `discovery-output.json` with Python to extract resource lists:
   ```python
   import json
   with open('.import/discovery-output.json') as f:
       data = json.load(f)
   # Example: extract VPC list
   vpcs = data.get('vpcs', {}).get('value', {}).get('vpcs', [])
   ```
4. Generate a partial discovery report for Group A resources

**Step 2 — Check if Group B is non-empty and handle CLI dependency**

If Group B is empty:
- Skip Step 3 and Step 4
- Use the Group A report directly, proceed to Phase 5

If Group B is non-empty:
- **Check apsara-cli installation**:
  ```bash
  which apsara >/dev/null 2>&1 && echo "INSTALLED" || echo "NOT_INSTALLED"
  ```

- **If NOT installed — prompt manual installation**:
  - Present the Group B resource list to the user
  - Explain: "The following resource types have no Terraform datasource and require apsara-cli for discovery, but it is not installed: [list]"
  - Prompt: "Please manually install apsara-cli first. Contact your ApsaraStack platform team or administrator to obtain the installation package."
  - **Stop all subsequent tasks** and wait for the user to install apsara-cli manually before continuing

**Step 3 — Discover Group B via apsara-cli**

For each resource type in Group B, execute the corresponding CLI commands from `references/api-commands.md`.

Parse JSON output, extract resource IDs and basic attributes, append to the discovery report.

**Step 4 — Merge and generate discovery report**

Merge Group A (datasource) and Group B (CLI) results into a single report:

```
Resource Discovery Report
=========================
Total resource types: N
Total resources: M

By type:
  alibabacloudstack_vpc_vpc        : 3  (datasource)
  alibabacloudstack_vpc_vswitch    : 6  (datasource)
  alibabacloudstack_instance       : 12 (datasource)
  alibabacloudstack_oss_bucket     : 2  (datasource + CLI for properties)
  alibabacloudstack_fc_service     : 1  (CLI)
  ...
```

**Empty discovery results handling:**

If a resource type (e.g., NAT Gateway) returns empty results but the user expects them to exist, proactively investigate:
- Confirm whether the resource is in a different Region
- Confirm whether the current credentials have Describe permission for that resource
- Confirm whether the resource has been released (can check audit logs)
- Inform user that discovery results don't match expectations, ask for confirmation before continuing

---

## Phase 5: Select Import Scope and Depth

Present the resource discovery report and ask the user:
- Import all or select specific resource types?
- Filter by VPC scope?
- Filter by Tag?
- **Import depth strategy**:
  1. Import only discovered resources (no related resource queries, depth=0)
  2. Import one level of related resources (depth=1)
  3. Import complete dependency tree (depth=unlimited, recommended)

**Build Resource Dependency Graph (based on depth strategy)**

Based on the user's chosen depth strategy:

**Depth=0 (import only discovered resources):**
- Skip API B calls
- Directly use static dependency rules from `references/dependency-rules.md` to sort Phase 4 discovered resources
- Generate import order

**Depth=1 (import one level of related resources):**
- For each resource discovered in Phase 4, if its type is in API B's support list -> Call API B (list-resource-relationships) to get directly related resources
- **Note**: API B is often unavailable in ApsaraStack environments; if it fails, immediately fall back to static dependency rules
- Do not recursively call API B for returned related resources
- Merge discovered resources and one-level related resources, topologically sort the resource graph, generate import order

**Depth=unlimited (import complete dependency tree, recommended):**
- From Phase 4 discovered resources, identify root resources with no parent dependencies (VPC, OSS Bucket, DNS domains, ASCM users, etc.)
- For each root resource, if its type is in API B's support list -> Call API B to get related resources
- **Note**: API B is often unavailable in ApsaraStack environments; if it fails, immediately fall back to static dependency rules
- Recursively call API B for related resources until no new resources appear (new resource set == empty set)
- For resource types not supported by API B -> Fall back to static dependency rules in `references/dependency-rules.md`
- Merge all relationship trees, topologically sort the resource graph, generate import order

**Fallback static import order (when API B is unavailable):**
```
Layer 0: VPC, OSS Bucket (no dependencies)
Layer 1: VSwitch, Security Group, EIP (depend on VPC)
Layer 2: NAT Gateway, SLB (depend on VPC + VSwitch)
Layer 3: ECS Instance (depends on VSwitch + Security Group)
Layer 4: Disk (depends on ECS)
Layer 5: RDS, Redis, MongoDB (depend on VSwitch)
Layer 6: DNS Record (depends on Domain)
```

For detailed dependency rules, refer to `references/dependency-rules.md`.

---

## Phase 6: Generate Terraform HCL Code

Process resources in batches by type.

**Data source priority (reuse Phase 4 results when possible):**

- **Group A resources (datasource-discovered in Phase 4)**: Reuse `.import/discovery-output.json` directly. The datasource output already contains full resource attributes in Terraform-native format. No additional API calls needed.
- **Group B resources (CLI-discovered in Phase 4)**: Execute detail retrieval commands via apsara-cli (JSON format), parse JSON, then generate HCL.
- **Hybrid resources (e.g., OSS Bucket)**: List-level data from datasource, property-level data from CLI if needed.

**Generation steps:**

1. Load existing data (discovery-output.json for Group A, or execute CLI for Group B)
2. Parse JSON, generate corresponding .tf files and write to working directory

One file per resource type, filename format: `<product>_<resource_type>.tf`:

```
vpc_vpcs.tf
vpc_vswitches.tf
vpc_security_groups.tf
vpc_eips.tf
vpc_nat_gateways.tf
ecs_instances.tf
ecs_disks.tf
oss_buckets.tf
rds_instances.tf
kvstore_instances.tf
dds_instances.tf
slb_loadbalancers.tf
alidns_domains.tf
alidns_records.tf
```

Refer to `references/terraform-patterns.md` for HCL code templates and `examples/` for examples.

**Cross-Resource Reference Resolution (three-layer strategy)**

When generating HCL for each resource, determine for each field value whether to replace with a terraform resource address, using the following priority:

1. **API B relationship first**: If the resource was discovered via API B (list-resource-relationships), its parent or related resources are already determined. Combined with the reference field whitelist extracted in Phase 3, find type-matching `_id` fields in the current resource and directly replace with the corresponding terraform resource address. The relationship is definitive, no guessing needed.

2. **Whitelist + mapping table dual-match fallback**: For resources not covered by API B, if the field name is in the Phase 3 reference field whitelist and the field value has a corresponding entry in the resource ID mapping table, replace with terraform resource address.

3. **dependency-rules.md static rules fallback**: For edge cases not covered by the above two layers (composite IDs, reference fields not ending in `_id`, etc.), refer to static rules in `references/dependency-rules.md`.

**Resource ID Mapping Table**

After generating each resource block, immediately append a record to the mapping table:

```
<cloud resource ID>  ->  <terraform resource address>
# Example:
vpc-bp1xxx   ->  alibabacloudstack_vpc_vpc.vpc_prod_main
vsw-bp1aaa   ->  alibabacloudstack_vpc_vswitch.vsw_prod_a
sg-bp1xxx    ->  alibabacloudstack_security_group.sg_web
```

Write the mapping table to `.import/id-mapping.json` for use in subsequent phases.

**Other key principles:**
- Resource names use `<type>_<name>` format (e.g., `alibabacloudstack_vpc_vpc.vpc_prod_main`)
- Password/secret fields use variables (`var.db_password`)
- Preserve original tags

---

## Phase 7: Import State

Based on the topological sort generated in Phase 5, generate the complete list of `terraform import` commands, present to user for confirmation, then execute in batches.

**Batch 1 — Network foundation:**
```bash
terraform import alibabacloudstack_vpc_vpc.vpc_<name> <vpc-id>
terraform import alibabacloudstack_vpc_vswitch.vsw_<name> <vsw-id>
terraform import alibabacloudstack_security_group.sg_<name> <sg-id>
```

**Batch 2 — Network ancillary:**
```bash
terraform import alibabacloudstack_nat_gateway.nat_<name> <nat-id>
terraform import alibabacloudstack_eip_address.eip_<name> <eip-id>
```

**Batch 3 — Compute:**
```bash
terraform import alibabacloudstack_instance.ecs_<name> <instance-id>
terraform import alibabacloudstack_disk.disk_<name> <disk-id>
```

**Batch 4 — Storage:**
```bash
terraform import alibabacloudstack_oss_bucket.bucket_<name> <oss-cluster>:<bucket-name>
```

**Note**: ApsaraStack OSS bucket import requires `<oss_cluster>:<bucket>` format (unlike public cloud which only needs `<bucket>`). The `oss_cluster` is the ID of the OSS cluster.

**Batch 5 — Database:**
```bash
terraform import alibabacloudstack_db_instance.rds_<name> <db-id>
terraform import alibabacloudstack_kvstore_instance.redis_<name> <redis-id>
terraform import alibabacloudstack_mongodb_instance.mongo_<name> <mongo-id>
```

**Batch 6 — Load Balancing:**
```bash
terraform import alibabacloudstack_slb_loadbalancer.slb_<name> <slb-id>
```

**Batch 7 — DNS:**
```bash
terraform import alibabacloudstack_dns_domain.domain_<name> <domain-name>
```

After each batch, analyze output and handle common errors:
- `ResourceNotFound`: Resource ID is incorrect, re-verify
- `PermissionDenied`: Need Describe permission for the corresponding resource
- `already managed`: Resource already in state, skip
- ID format error: Refer to import format in `references/resource-types.md`

---

## Phase 8: Validation

Execute the following commands to validate import results:

```bash
# List all resources in state
terraform state list 2>&1

# Run plan, target is No changes
terraform plan -refresh=true 2>&1
```

Analyze plan output:
- `No changes` -> Validation passed, proceed to Phase 9
- Has diff -> Analyze the cause, refer to common diff fix patterns in `references/terraform-patterns.md`, provide fix recommendations and inform user, execute fix after user confirmation

Common diff causes:
- Tag format mismatch (need to adjust tags syntax in HCL)
- Password fields cannot be read from API (use `ignore_changes` to ignore)
- Computed field differences (add `lifecycle { ignore_changes = [...] }`)
- Default value differences (explicitly set values matching the cloud)

---

## Phase 9: Resource Relationship Graph

Execute the following commands to generate a resource relationship graph:

```bash
# Export state JSON
terraform show -json > .import/tf_state.json 2>&1

# Extract dependencies
python3 -c "
import json, sys
with open('.import/tf_state.json') as f:
    state = json.load(f)
resources = state.get('values', {}).get('root_module', {}).get('resources', [])
print(f'Total resources: {len(resources)}')
for r in resources:
    deps = r.get('depends_on', [])
    if deps:
        for d in deps:
            print(f'  {d} --> {r[\"address\"]}')
" 2>&1
```

Based on output, generate a text-format resource relationship graph (tree structure) showing:
- Network layer -> Compute layer -> Data layer dependency chains
- Security group cross-cutting references
- Load balancer backend binding relationships

Refer to `examples/dependency-graph.txt` for example format.

---

## Incremental Sync Mode

Use this mode when a Terraform working directory already exists. Automatically detects three types of changes:

**1. Detect new resources (exist in cloud, not in state):**

Execute the following commands to compare state and cloud ID lists:
```bash
# ECS example
terraform state list | grep alibabacloudstack_instance
apsara ecs describe-instances --biz-region-id $REGION --page-size 100 \
  --output cols=InstanceId rows=Instances.Instance 2>&1
```

Compare the two lists, find the difference set, generate HCL + import commands for new resources, inform user for confirmation before executing.

**2. Detect deleted resources (exist in state, not in cloud):**

Run existence checks for each resource in state; if API returns empty, execute:
```bash
terraform state rm alibabacloudstack_instance.<name>
```

**3. Detect configuration drift:**

Execute the following command to detect drift:
```bash
terraform plan -refresh=true 2>&1
```

For detailed incremental sync procedures, refer to `references/incremental-sync.md`.

---

## References

- `references/resource-types.md` — Complete resource type mapping table (20+ resource types)
- `references/terraform-patterns.md` — HCL code templates + common diff fixes
- `references/state-management.md` — terraform state command reference
- `references/api-commands.md` — apsara-cli command quick reference (native product API fallback path)
- `references/dependency-rules.md` — Resource dependency layers and import order (API B fallback path)
- `references/resource-hub-api.md` — ResourceCenter API guide (list-resource-types / search-resources / list-resource-relationships)
- `references/incremental-sync.md` — Detailed incremental sync procedures
