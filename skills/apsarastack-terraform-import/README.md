# 阿里云专有云资源 Terraform 导入 (apsarastack-terraform-import)

将阿里云专有云（ApsaraStack）账号下的现有资源导入到 Terraform 管理。自动执行环境检查、资源发现、HCL 生成、state 导入等操作。

> **注意**：本 Skill 面向阿里云专有云环境（ApsaraStack），使用 `aliyun/alibabacloudstack` Provider。公有云环境请使用对应的公有云导入 Skill。

## 功能概述

- **环境检查**：自动检测 apsara-cli、Terraform、Python3 等工具安装状态
- **认证配置**：支持环境变量、CLI 配置文件、ECS RAM 角色等多种认证方式
- **资源发现**：通过 ResourceCenter API 或原生产品 API 批量查询资源
- **依赖分析**：构建资源依赖图，按正确顺序导入
- **HCL 生成**：自动生成 Terraform 资源配置代码
- **State 导入**：执行 terraform import 将资源纳入管理
- **增量同步**：支持已有工作目录的增量资源同步

## 支持模式

- **一次性迁移**：完整执行环境检查到验证的全流程
- **增量同步**：在已有 Terraform 工作目录基础上，同步新增/删除/变更的资源

## 目录结构

```
.
├── .qoder-plugin/
│   └── plugin.json          # 插件清单
├── skills/
│   └── apsarastack-terraform-import/
│       ├── SKILL.md          # Skill 主定义文件
│       ├── examples/
│       │   └── provider.tf   # Provider 配置示例
│       └── references/
│           ├── api-commands.md      # apsara-cli 命令速查
│           ├── dependency-rules.md  # 资源依赖规则
│           ├── incremental-sync.md  # 增量同步详细流程
│           ├── ram-policies.md      # RAM 权限参考
│           ├── resource-hub-api.md  # ResourceCenter API 指南
│           ├── resource-types.md    # 资源类型映射表
│           ├── state-management.md  # terraform state 命令参考
│           └── terraform-patterns.md # HCL 代码模板与常见修复
└── README.md
```

## 来源

- **原始仓库**：https://github.com/aliyun/alibabacloud-aiops-skills
- **Skill 路径**：skills/migrationom/solutions/apsarastack-terraform-import
- **Agent Explorer**：https://agentexplorer.aliyuncs.com

## 触发条件

当用户提到以下关键词时，该 skill 将被激活：

- 导入专有云资源到 Terraform
- terraform import 专有云
- terraform import ApsaraStack
- 将现有专有云资源纳入 Terraform 管理
- 阿里云专有云资源迁移 Terraform
- 生成 terraform state
- import alibabacloudstack resources
- 阿里云专有云 IaC 迁移
- 阿里云专有云 Terraform 导入

## 支持范围

理论上支持 alibabacloudstack provider 中任何具有 Read 方法的资源。
预置了约 50 种常见资源类型的验证命令、HCL 模板和导入 ID 格式：
VPC、ECS、SLB、OSS、RDS、Redis、MongoDB、PolarDB、NAT Gateway、EIP、Security Group、VSwitch、DNS、RAM、ACK、FC 等。
