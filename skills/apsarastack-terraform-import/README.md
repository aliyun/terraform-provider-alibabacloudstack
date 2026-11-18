# Alibaba Cloud ApsaraStack Resource Terraform Import (apsarastack-terraform-import)

Import existing resources under an Alibaba Cloud ApsaraStack account into Terraform management. Automatically performs environment checks, resource discovery, HCL generation, state import, and more.

> **Note**: This skill is designed for the Alibaba Cloud ApsaraStack environment, using the `aliyun/alibabacloudstack` Provider. For public cloud environments, please use the corresponding public cloud import skill.

## Feature Overview

- **Environment Check**: Automatically detects the installation status of tools such as apsara-cli, Terraform, and Python3
- **Authentication Configuration**: Supports multiple authentication methods including environment variables, CLI configuration files, and ECS RAM roles
- **Resource Discovery**: Batch query resources via ResourceCenter API or native product APIs
- **Dependency Analysis**: Builds resource dependency graphs and imports in the correct order
- **HCL Generation**: Automatically generates Terraform resource configuration code
- **State Import**: Executes `terraform import` to bring resources under management
- **Incremental Sync**: Supports incremental resource sync for existing working directories

## Supported Modes

- **One-time Migration**: Complete end-to-end process from environment check to validation
- **Incremental Sync**: Sync newly added / deleted / changed resources based on an existing Terraform working directory

## Directory Structure

```
.
├── .qoder-plugin/
│   └── plugin.json          # Plugin manifest
├── skills/
│   └── apsarastack-terraform-import/
│       ├── SKILL.md          # Main skill definition file
│       ├── examples/
│       │   └── provider.tf   # Provider configuration example
│       └── references/
│           ├── api-commands.md      # apsara-cli command quick reference
│           ├── dependency-rules.md  # Resource dependency rules
│           ├── incremental-sync.md  # Detailed incremental sync workflow
│           ├── ram-policies.md      # RAM permission references
│           ├── resource-hub-api.md  # ResourceCenter API guide
│           ├── resource-types.md    # Resource type mapping table
│           ├── state-management.md  # terraform state command reference
│           └── terraform-patterns.md # HCL code templates and common fixes
└── README.md
```

## Source

- **Original Repository**: https://github.com/aliyun/alibabacloud-aiops-skills
- **Skill Path**: skills/migrationom/solutions/apsarastack-terraform-import
- **Agent Explorer**: https://agentexplorer.aliyuncs.com

## Trigger Conditions

This skill will be activated when the user mentions the following keywords:

- Import ApsaraStack resources to Terraform
- terraform import ApsaraStack
- terraform import ApsaraStack
- Bring existing ApsaraStack resources into Terraform management
- Alibaba Cloud ApsaraStack resource migration to Terraform
- Generate terraform state
- import alibabacloudstack resources
- Alibaba Cloud ApsaraStack IaC migration
- Alibaba Cloud ApsaraStack Terraform import

## Supported Scope

Theoretically supports any resource in the alibabacloudstack provider that has a Read method.
Pre-configured with verified API discovery commands, HCL templates, and import ID formats for approximately 50 common resource types:
VPC, ECS, SLB, OSS, RDS, Redis, MongoDB, PolarDB, NAT Gateway, EIP, Security Group, VSwitch, DNS, RAM, ACK, FC, etc.
