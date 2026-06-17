---
name: apsarastack-terraform-code-generation
description: |
  Generate validated Terraform HCL code for Alibaba Cloud ApsaraStack infrastructure.
  Uses Terraform MCP as the sole documentation source.
  Enforces mandatory deprecation checks, data source resolution (no hardcoded IDs), and ApsaraStack-specific configuration (popgw_domain, protocol).
  Covers VPC, ECS, RDS, OSS, SLB, and all alibabacloudstack_* resources with automatic pattern matching.
  Trigger phrases: "write terraform for alibabacloudstack", "generate ApsaraStack Terraform", "ApsaraStack HCL", "deploy Alibaba Cloud ApsaraStack with Terraform", "alibabacloudstack provider", "apsarastack terraform", "ApsaraStack Enterprise Edition terraform".
---

# Alibaba Cloud ApsaraStack Terraform Code Generation

Transforms natural language descriptions of Alibaba Cloud ApsaraStack infrastructure requirements into validated Terraform code for the current `aliyun/alibabacloudstack` Provider. Resource knowledge is fetched from the Provider's own documentation at generation time.

## Hard Rules (Strictly Enforced)

### 1. Credential Security — Never Leak, Never Request

Nowhere (HCL, comments, environment variable declarations, Shell output, logs) should you **ever** read, print, ask for, or write AK/SK values. The alibabacloudstack Provider requires credentials to be explicitly configured in the provider block or provided via environment variables (`ALIBABACLOUDSTACK_ACCESS_KEY`, `ALIBABACLOUDSTACK_SECRET_KEY`). All credentials are read by the Provider itself; this Skill never touches them.

### 2. Terraform Execution Restrictions

This Skill **primarily generates Terraform HCL code**, but after file generation, can execute `terraform validate` or `tofu validate` for syntax validation and fixes (Step 7).

This Skill **never** executes the following Terraform commands:
- `terraform init` (initialization)
- `terraform plan` (planning)
- `terraform apply` (application)
- Other commands that modify infrastructure or require network connections

This Skill's scope includes:
- Generating valid Terraform HCL code
- Creating project structure and files
- Providing documentation and usage instructions
- Executing `terraform validate` or `tofu validate` for local syntax validation (Step 7)

### 4. ApsaraStack-Specific Configuration

**Note: This Skill does not generate Provider configuration or related variables.**

Compared to the public cloud Provider, the alibabacloudstack Provider requires additional configuration (e.g., `popgw_domain`, `protocol`, `insecure`), but these are managed by users in their own provider blocks.

This Skill focuses on generating resource definitions and data sources.

## Environment Requirements (Soft Recommendations)

- **Terraform ≥ 0.13** (Provider requirement).
- **Network connection** is required — Step 4.1 uses Terraform MCP to fetch Provider documentation.
- **Go 1.13+** (if building Provider from source, rarely needed).

## Prerequisites: Installing Terraform MCP Server

Before using this Skill, you must ensure the Terraform MCP Server is correctly installed and configured. The MCP Server is the primary source for this Skill to fetch Provider documentation.

### Installation Steps

**Method 1: Install via npm (Recommended)**

```bash
npm install -g @modelcontextprotocol/server-terraform
```

**Method 2: Install via pip**

```bash
pip install mcp-server-terraform
```

**Method 3: Build from source**

```bash
git clone https://github.com/modelcontextprotocol/servers.git
cd servers/src/terraform
npm install
npm run build
```

### Configure MCP Server

Add the Terraform Server to Aone Copilot's MCP configuration file:

```json
{
  "mcpServers": {
    "terraform": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-terraform"]
    }
  }
}
```

Or if using the pip-installed version:

```json
{
  "mcpServers": {
    "terraform": {
      "command": "python",
      "args": ["-m", "mcp_server_terraform"]
    }
  }
}
```

### Verify Installation

After restarting Aone Copilot, verify the MCP Server is working correctly:

1. Try calling Terraform MCP tools in the conversation
2. If Provider information is returned successfully, installation is successful
3. If failed, check:
   - Whether MCP Server is correctly installed in PATH
   - Whether Aone Copilot's MCP configuration is correct
   - Whether network connection is normal

**Note**: The MCP Server is the only source for this Skill to fetch Provider documentation, ensuring accuracy and timeliness. If the MCP Server is unavailable, the task will stop and require human intervention.

## Workflow

### Step 1. Parse Requirements

Extract:

- `region` — default `cn-hangzhou-env01-d01` (ApsaraStack region format).
- `resources[]` — `{ alibabacloudstack_type, quantity, attributes }`.
- Non-functional requirements: multi-AZ, encryption, backup, high availability, IOPS.

If ambiguity exists (e.g., "set up database"), ask **at most** one clarifying question.

### Step 2. Determine Target Directory

Extract `<target-dir>` from user request (explicit path like `myshop-infra/`, or current working directory if unspecified). All subsequent `fmt` / `init` / `validate` commands run in this directory.

Before writing any `.tf` files, **must** create directory:

```bash
mkdir -p <target-dir>
```

All file write paths must be prefixed with `<target-dir>/` — never write directly to current working directory or generic `outputs/` parent directory. After generation, verify structure:

```bash
ls -R <target-dir>
```

### Step 3. Draw Architecture Sketch

Before writing any HCL, draw a dependency table — one row per resource:

| resource | depends on | AZ / placement |
| --- | --- | --- |

- Extend `resources[]`, adding implicit infrastructure (VPC → VSwitch → SecurityGroup → workloads); users often skip these when parsing.
- The extended list is input for Step 4 gating.

### Step 4. HCL Pre-flight Gating (Mandatory)

For each distinct `alibabacloudstack_*` type (resources **and** data sources) in Step 3, execute the following steps. Calls per type are independent — **execute in parallel across types**.

#### 4.1 Fetch Provider Documentation via Terraform MCP (Sole Source)

**Step 4.1: Use Terraform MCP**

Use Terraform MCP tools to fetch Provider documentation:

1. Call `terraform::tool::search_providers` with parameters:
   - `provider_name`: `alibabacloudstack`
   - `provider_namespace`: `aliyun`
   - `service_slug`: resource name (e.g., `vpc`, `vswitch`)
   - `provider_document_type`: `resources` (`data-sources` for data sources)

2. Identify matching resource from results and obtain its `providerDocID`.

3. Call `terraform::tool::get_provider_details` with `provider_doc_id` to fetch complete documentation.

**Success criteria**: If MCP returns valid documentation with required/optional parameters and example usage, **proceed directly to Step 4.3 (Recitation)**.

**Failure handling**: If Terraform MCP tool fails or returns useless content, **immediately stop task and request human intervention**.

When stopping, report to user:
- Which resource's documentation fetch failed
- Error message or empty result returned by MCP
- Advise user to check:
  - Whether Terraform MCP Server is correctly installed and running
  - Whether network connection is normal
  - Whether resource name is correct

**Note**: This Skill **no longer provides local documentation fallback**. MCP Server is the sole source for fetching Provider documentation, ensuring accuracy and timeliness.

#### 4.2 Pattern Lookup (Conditional)

If user requirements match product-specific idioms listed in `references/resource-patterns-apsarastack.md` (e.g., RDS cross-AZ HA, OSS lifecycle), read relevant sections.

When a matching pattern section is found, **all attributes listed in that section's "Required Attributes" table must appear in the generated HCL**.

```bash
grep -in "<keyword>" references/resource-patterns-apsarastack.md
```

#### 4.3 Recitation (Reading Proof)

Before writing any HCL, issue and verify a complete per-resource briefing:

- **Required** parameters (verbatim list from MCP documentation)
- **2-5 key optional** parameters (relevant to user requirements)
- Minimal HCL snippet from documentation's "Example Usage"

If required or optional parameters are missing, return to Step 4.1. Skipping or using partial recitation is a critical failure.

### Step 5. Generation

#### 5.1 Write HCL from Recitation, Not Memory

**Only** use parameters established in Step 4.3. If parameters not in the recitation briefing are needed, re-fetch Step 4.2 through deeper reading; do not guess.

Before writing fields, look up the resource in `references/deprecated-fields-apsarastack.md`:

```bash
grep '`alibabacloudstack_<resource>`' references/deprecated-fields-apsarastack.md
```

If user requirements involve products with specific usage patterns, also consult `references/resource-patterns-apsarastack.md`.

#### 5.2 Data Source Mandate (Mandatory — No Hardcoded IDs)

Resolve via `data` blocks, never use literals:

- `zone_id` → `data "alibabacloudstack_zones"`。
- `image_id` → `data "alibabacloudstack_images"` (filter via `name_regex`, `owners = "system"`, `most_recent = true`).
- `instance_type` → `data "alibabacloudstack_instance_types"`。

**Data Source `ids` Parameter — Prevent Empty Strings:**

When a data source uses the `ids` parameter to query by ID list, **passing empty strings is strictly prohibited**. If a variable may be empty, must use `coalesce()` as fallback:

```hcl
# ❌ Wrong: var.vpc_id may be empty string, resulting in ids = [""]
data "alibabacloudstack_vpc_vpcs" "existing" {
  ids = [var.vpc_id]
}

# ✅ Correct: Use coalesce fallback, pass a non-existent ID when empty
data "alibabacloudstack_vpc_vpcs" "existing" {
  ids = [coalesce(var.vpc_id, "vpc-nonexistent-placeholder")]
}
```

This way, when the variable is empty, the data source returns an empty list (`vpcs` length is 0), without error.

**Conditional Creation Drift Issue (Important):**

When implementing "query-then-create" logic, **using data source query results to control `count` is strictly prohibited**. Because:
- First apply: data source finds nothing → create resource
- Second apply: data source finds it → `count` becomes 0 → **delete resource**

**Correct approach**: Use whether variable is empty to control creation logic; data source is only for fetching attributes of existing resources:

```hcl
# ❌ Wrong: Using data source results to control count will cause resource deletion on second apply
locals {
  vpc_exists = length(data.alibabacloudstack_vpc_vpcs.existing.vpcs) > 0
}
resource "alibabacloudstack_vpc_vpc" "new" {
  count = local.vpc_exists ? 0 : 1  # Second apply will become 0, deleting resource!
}

# ✅ Correct: Use variable emptiness to control count; data source only for fetching attributes
locals {
  vpc_exists = var.vpc_id != ""
  vpc_id     = local.vpc_exists ? var.vpc_id : alibabacloudstack_vpc_vpc.new[0].id
}
resource "alibabacloudstack_vpc_vpc" "new" {
  count = local.vpc_exists ? 0 : 1  # As long as var.vpc_id does not change, count will not change
}
```

**Key principle**: The condition for `count` must be based on **user-input variables**, not **data source query results**.

#### 5.3 Provider Block (Content Contract)

**Important: This Skill does not generate `provider "alibabacloudstack" {}` blocks.** 
Provider configuration is the user's responsibility and should be managed separately.

This Skill only generates `terraform { required_providers {} }` blocks to declare Provider requirements.

### File Organization (Mandatory)

Generated code **must** be organized into exactly four files:

1. **`version.tf`** — Terraform and Provider version constraints
   
   ```hcl
   terraform {
     required_version = ">= 0.13"
     required_providers {
       alibabacloudstack = {
         source  = "aliyun/alibabacloudstack"
         version = "< 3.19.0"
       }
     }
   }
   ```

2. **`variables.tf`** — Variable declarations for resource configuration
   
   - Contains only variables needed for resources (no Provider credentials)
   - Use sensible defaults where appropriate
   - Mark sensitive variables with `sensitive = true`

3. **`main.tf`** — Main logic, including:
   
   - Data sources (zones, images, instance types, etc.)
   - Resource definitions (VPC, ECS, RDS, etc.)
   - All infrastructure components

4. **`outputs.tf`** — Output definitions
   
   - Export important resource IDs, IPs, names, etc.
   - Include descriptive descriptions for each output

### Version Constraint Rules

- Provider version constraint: **must use `< 3.19.0`** as the default upper bound for ApsaraStack compatibility.
- **Do not** generate any `provider "alibabacloudstack" {}` blocks containing credentials, popgw_domain, protocol, or other configuration.
- Users must configure provider blocks themselves in their own `.tf` files.

**Post-generation validation**:

```bash
# Verify all four required files exist
for file in version.tf variables.tf main.tf outputs.tf; do
  test -f <target-dir>/$file && echo "OK: $file exists" || echo "MISSING: $file"
done

# Verify required_providers contains aliyun/alibabacloudstack with correct version constraint
grep -Rq 'alibabacloudstack.*source.*=.*"aliyun/alibabacloudstack"' \
  <target-dir>/version.tf && echo OK_SOURCE || echo BAD_SOURCE

grep -Rq 'version.*=.*"< 3.19.0"' <target-dir>/version.tf \
  && echo OK_VERSION_CONSTRAINT || echo BAD_VERSION_CONSTRAINT

# Ensure this Skill did not generate any provider blocks
! grep -Rq 'provider "alibabacloudstack"' <target-dir>/*.tf \
  && echo OK_NO_PROVIDER_BLOCK || echo UNEXPECTED_PROVIDER_BLOCK
```

All checks must pass.

#### 5.4 Style Baseline

- 2-space indentation; `=` aligned within blocks; snake_case for semantic resource labels.
- Every resource that supports tags should carry a non-empty `tags` block.

#### 5.5 Deprecated Field Audit — Static grep Check (Mandatory)

Run before requiring `terraform` — this is a pure grep check on the just-written HCL. For each resource in this generation, grep items against `references/deprecated-fields-apsarastack.md` and handle each line type:

- **rename** lines → replace old field name with new field name.
- **split / soft-split** lines → **do not** write inline fields on parent resource; only declare replacement sub-resources when needed.
- **deprecated-no-replacement** lines → stop using the field, no replacement.

**Post-audit validation (bash grep — must all return OK)**:

```bash
grep '| `alibabacloudstack_' references/deprecated-fields-apsarastack.md | while IFS='|' read _ resource field kind _; do
  resource=$(echo "$resource" | tr -d ' `')
  field=$(echo "$field" | tr -d ' ')
  kind=$(echo "$kind" | tr -d ' ')
  if grep -Rq "resource \"$resource\"" <target-dir>/*.tf; then
    case "$kind" in
      rename|deprecated-no-replacement)
        grep -q "\b$field\b\s*=" <target-dir>/*.tf \
          && echo "DEPRECATED: $resource.$field" || echo "OK: $resource.$field"
        ;;
      split|soft-split)
        grep -q "\b$field\b\s*=" <target-dir>/*.tf \
          && echo "DEPRECATED: $resource.$field (inline)" \
          || echo "OK: $resource.$field (not inline)"
        ;;
    esac
  fi
done
```

**Hard gating: must pass before Step 6** — if any `DEPRECATED:` lines appear, fix HCL and re-run until all lines return `OK:`.

### Step 6. Coverage Check + Summary

**Mandatory — run regardless of generation results.**

**Coverage check.** Enumerate resource blocks in generated HCL and compare with Step 3 sketch. If any sketch lines are missing, return to Step 5 and add them.

**Summary template** — print in user's language:

```
Files written:
<path/to/file1>
<path/to/file2>
...

Validation: pending (proceed to Step 7)

Deprecation routing: <if rerouted: `<original_name>` → `<new_name>`; otherwise: none>

<optional: architecture notes, design decisions, deployment tips>
```

### Step 7. Terraform/OpenTofu Validation and Fixes (Conditional Execution)

**Trigger condition**: After file generation, check whether user has `terraform` or `tofu` (OpenTofu) installed locally.

**Detect tool availability**:

```bash
# Detect terraform
which terraform 2>/dev/null || command -v terraform 2>/dev/null

# Detect tofu (OpenTofu)
which tofu 2>/dev/null || command -v tofu 2>/dev/null
```

**Execution logic**:

1. **If `terraform` is detected**:
   - Execute `terraform init` in `<target-dir>`
   - If `init` succeeds, continue with `terraform validate`
     - If validation passes, record success status
     - If validation fails, analyze error messages and automatically fix HCL code, then re-validate (retry up to 3 times)
   - If `init` fails, output prompt message and skip validation step

2. **If `tofu` is detected (but `terraform` is not detected)**:
   - Execute `tofu init` in `<target-dir>`
   - If `init` succeeds, continue with `tofu validate`
     - If validation passes, record success status
     - If validation fails, analyze error messages and automatically fix HCL code, then re-validate (retry up to 3 times)
   - If `init` fails, output prompt message and skip validation step

3. **If neither is detected**:
   - Skip validation step
   - Note in summary "Validation: skipped (terraform or tofu not detected)"
   - Remind user to manually install Terraform or OpenTofu for local validation

**Validation execution example**:

```bash
cd <target-dir>

# Validate using terraform
if command -v terraform &>/dev/null; then
  echo "Terraform detected, executing initialization..."
  terraform init
  
  if [ $? -eq 0 ]; then
    echo "Initialization successful, executing validation..."
    terraform validate
    
    if [ $? -ne 0 ]; then
      echo "Validation failed, analyzing errors and fixing..."
      # Fix HCL files based on error messages
      # Re-validate (up to 3 times)
    fi
  else
    echo "Initialization failed, skipping validation step"
    echo "Tip: Please check network connection, Provider configuration, or manually execute 'terraform init'"
  fi
  
# Or validate using tofu
elif command -v tofu &>/dev/null; then
  echo "OpenTofu detected, executing initialization..."
  tofu init
  
  if [ $? -eq 0 ]; then
    echo "Initialization successful, executing validation..."
    tofu validate
    
    if [ $? -ne 0 ]; then
      echo "Validation failed, analyzing errors and fixing..."
      # Fix HCL files based on error messages
      # Re-validate (up to 3 times)
    fi
  else
    echo "Initialization failed, skipping validation step"
    echo "Tip: Please check network connection, Provider configuration, or manually execute 'tofu init'"
  fi
  
else
  echo "Terraform or OpenTofu not detected, skipping validation"
fi
```

**Common errors and fix strategies**:

- **Syntax errors** (mismatched brackets, missing commas, etc.) → Fix HCL syntax
- **Undefined variables** → Check if declared in `variables.tf`
- **Undefined resource references** → Check if resource names are correct
- **Type mismatches** → Fix variable types or add type conversions
- **Missing Provider configuration** → Remind user to configure provider block themselves (this Skill does not generate)
- **init failure** → Possible causes: network issues, Provider version does not exist, ApsaraStack environment configuration issues; advise user to manually execute init and check detailed errors

**Update summary after validation**:

Update validation status in Step 6 summary template:

```
Validation: <success/failure/skipped>
- Tool: <terraform/tofu/none>
- Initialization: <success/failure/skipped>
- Result: <validation passed/validation failed, X errors fixed/unable to auto-fix/init failed skipped>
- Error details: <if errors, list key error messages>
```

**Note**:
- This step will execute `terraform init` or `tofu init` to initialize Provider plugins, but will not execute `plan`, `apply`, or other commands that modify infrastructure
- If `init` fails (usually due to network issues or Provider configuration issues), will skip validation step and report to user
- If validation fails and cannot be automatically fixed, report specific errors to user and suggest manual fixes
- Retry validation-fix cycle up to 3 times, stop after and report remaining errors

## References

| Source | When to read |
| --- | --- |
| Terraform MCP (`terraform::tool::search_providers`, `terraform::tool::get_provider_details`) | Step 4.1 — authoritative required/optional parameters for each resource (sole source) |
| `references/deprecated-fields-apsarastack.md` (local) | Step 5.1 — known field-level renames |
| `references/resource-patterns-apsarastack.md` (local) | Step 5.1 — product-specific idioms |
| `references/auth-and-network-apsarastack.md` (local) | User reference for credential configuration |
