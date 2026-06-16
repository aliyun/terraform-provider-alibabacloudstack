---
name: apsarastack-terraform-code-generation
description: |
  为阿里云专有云（ApsaraStack）基础设施生成经过验证的 Terraform HCL 代码。
  使用 Terraform MCP 作为唯一文档来源。
  强制执行强制性的弃用检查、数据源解析（禁止硬编码 ID）以及专有云特定配置（popgw_domain、protocol）。
  涵盖 VPC、ECS、RDS、OSS、SLB 以及所有 alibabacloudstack_* 资源，支持自动模式匹配。
  触发词："write terraform for alibabacloudstack"、"生成专有云 Terraform"、"专有云 HCL"、"用 Terraform 部署阿里云专有云"、"alibabacloudstack provider"、"apsarastack terraform"、"飞天企业版 terraform"。
---

# 阿里云专有云 Terraform 代码生成

将自然语言描述的阿里云专有云基础设施需求转化为针对当前 `aliyun/alibabacloudstack` Provider 的已验证 Terraform 代码。资源知识在生成时从 Provider 自身的文档中获取。

## 硬性规则（严禁违反）

### 1. 凭证安全——永不泄露，永不需要

在任何地方（HCL、注释、环境变量声明、Shell 输出、日志）都**绝不**读取、打印、询问或写入 AK/SK 值。alibabacloudstack Provider 需要在 provider 块中显式配置凭证或通过环境变量（`ALIBABACLOUDSTACK_ACCESS_KEY`、`ALIBABACLOUDSTACK_SECRET_KEY`）提供。所有凭证都由 Provider 自身读取，本 Skill 绝不接触。

### 2. Terraform 执行限制

本 Skill **主要生成 Terraform HCL 代码**，但在文件生成完成后，可以执行 `terraform validate` 或 `tofu validate` 进行语法验证和修复（步骤 7）。

本 Skill **绝不**执行以下 Terraform 命令：
- `terraform init`（初始化）
- `terraform plan`（计划）
- `terraform apply`（应用）
- 其他会修改基础设施或需要网络连接的命令

本 Skill 的范围包括：
- 生成有效的 Terraform HCL 代码
- 创建项目结构和文件
- 提供文档和使用说明
- 执行 `terraform validate` 或 `tofu validate` 进行本地语法验证（步骤 7）

### 4. 专有云特定配置

**注意：本 Skill 不生成 Provider 配置或相关变量。**

与公有云 Provider 相比，alibabacloudstack Provider 需要额外的配置（例如 `popgw_domain`、`protocol`、`insecure`），但这些由用户在自己的 provider 块中管理。

本 Skill 专注于生成资源定义和数据源。

## 环境要求（软性建议）

- **Terraform ≥ 0.13**（Provider 要求）。
- **网络连接**是必需的——步骤 4.1 使用 Terraform MCP 获取 Provider 文档。
- **Go 1.13+**（如果需要从源码构建 Provider，很少需要）。

## 前置准备：安装 Terraform MCP Server

在使用本 Skill 之前，必须确保 Terraform MCP Server 已正确安装和配置。MCP Server 是本 Skill 获取 Provider 文档的主要来源。

### 安装步骤

**方法 1：通过 npm 安装（推荐）**

```bash
npm install -g @modelcontextprotocol/server-terraform
```

**方法 2：通过 pip 安装**

```bash
pip install mcp-server-terraform
```

**方法 3：从源码构建**

```bash
git clone https://github.com/modelcontextprotocol/servers.git
cd servers/src/terraform
npm install
npm run build
```

### 配置 MCP Server

在 Aone Copilot 的 MCP 配置文件中添加 Terraform Server：

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

或者如果使用 pip 安装的版本：

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

### 验证安装

重启 Aone Copilot 后，验证 MCP Server 是否正常工作：

1. 在对话中尝试调用 Terraform MCP 工具
2. 如果成功返回 Provider 信息，说明安装成功
3. 如果失败，检查：
   - MCP Server 是否正确安装在 PATH 中
   - Aone Copilot 的 MCP 配置是否正确
   - 网络连接是否正常

**注意**：MCP Server 是本 Skill 获取 Provider 文档的唯一来源，确保文档的准确性和时效性。如果 MCP Server 无法使用，任务将停止并需要人工介入处理。

## 工作流程

### 步骤 1. 解析需求

提取：

- `region`——默认值 `cn-hangzhou-env01-d01`（专有云区域格式）。
- `resources[]`——`{ alibabacloudstack_type, quantity, attributes }`。
- 非功能性需求：多可用区、加密、备份、高可用、IOPS。

如果存在歧义（例如"设置数据库"），**最多**提出一个澄清问题。

### 步骤 2. 确定目标目录

从用户请求中提取 `<target-dir>`（明确的路径如 `myshop-infra/`，或未指定时使用当前工作目录）。所有后续的 `fmt` / `init` / `validate` 命令都在此目录中运行。

在编写任何 `.tf` 文件之前，**必须**创建目录：

```bash
mkdir -p <target-dir>
```

所有文件写入的路径必须以 `<target-dir>/` 为前缀——绝不要直接写入当前工作目录，也不要写入通用的 `outputs/` 父目录。生成完成后，验证结构：

```bash
ls -R <target-dir>
```

### 步骤 3. 绘制架构草图

在编写任何 HCL 之前，绘制依赖关系表——每个资源一行：

| resource | depends on | AZ / placement |
| --- | --- | --- |

- 扩展 `resources[]`，添加隐含的基础设施（VPC → VSwitch → SecurityGroup → 工作负载）；用户解析时经常会跳过这些。
- 扩展后的列表是步骤 4 门禁的输入。

### 步骤 4. HCL 前置门禁（强制）

对于步骤 3 中的每个不同的 `alibabacloudstack_*` 类型（资源**和**数据源），执行以下步骤。每种类型的调用是独立的——**跨类型并行执行**。

#### 4.1 通过 Terraform MCP 获取 Provider 文档（唯一方式）

**步骤 4.1：使用 Terraform MCP**

使用 Terraform MCP 工具获取 Provider 文档：

1. 调用 `terraform::tool::search_providers`，参数：
   - `provider_name`: `alibabacloudstack`
   - `provider_namespace`: `aliyun`
   - `service_slug`: 资源名称（例如 `vpc`、`vswitch`）
   - `provider_document_type`: `resources`（数据源则为 `data-sources`）

2. 从结果中识别匹配的资源并获取其 `providerDocID`。

3. 使用 `provider_doc_id` 调用 `terraform::tool::get_provider_details` 获取完整文档。

**成功标准**：如果 MCP 返回包含必需/可选参数和示例用法的有效文档，**直接进入步骤 4.3（复述）**。

**失败处理**：如果 Terraform MCP 工具失败或返回无用的内容，**立即停止任务并请求人工处理**。

停止时应向用户报告：
- 哪个资源的文档获取失败
- MCP 返回的错误信息或空结果
- 建议用户检查：
  - Terraform MCP Server 是否正确安装和运行
  - 网络连接是否正常
  - 资源名称是否正确

**注意**：本 Skill **不再提供本地文档回退方案**。MCP Server 是获取 Provider 文档的唯一来源，确保文档的准确性和时效性。

#### 4.2 模式查找（条件性）

如果用户需求匹配 `references/resource-patterns-apsarastack.md` 中列出的产品特定惯用法（例如 RDS 跨可用区高可用、OSS 生命周期），读取相关部分。

当找到匹配的模式部分时，该部分"必需属性"表中列出的**所有属性都必须出现在生成的 HCL 中**。

```bash
grep -in "<keyword>" references/resource-patterns-apsarastack.md
```

#### 4.3 复述（阅读证明）

在编写任何 HCL 之前，发出并验证完整的每资源简报：

- **必需**参数（来自 MCP 文档的逐字列表）
- **2–5 个关键可选**参数（与用户需求相关）
- 来自文档"示例用法"的最小 HCL 片段

如果缺少必需或可选参数，返回步骤 4.1。跳过或使用部分复述是严重失败。

### 步骤 5. 生成

#### 5.1 根据复述编写 HCL，而非凭记忆

**仅**使用步骤 4.3 中建立的参数。如果需要未在复述简报中的参数，通过更深入的读取重新获取步骤 4.2；不要猜测。

在编写字段之前，在 `references/deprecated-fields-apsarastack.md` 中查找资源：

```bash
grep '`alibabacloudstack_<resource>`' references/deprecated-fields-apsarastack.md
```

如果用户需求涉及具有特定使用模式的产品，还需查阅 `references/resource-patterns-apsarastack.md`。

#### 5.2 数据源强制（强制——禁止硬编码 ID）

通过 `data` 块解析，绝不使用字面量：

- `zone_id` → `data "alibabacloudstack_zones"`。
- `image_id` → `data "alibabacloudstack_images"`（通过 `name_regex`、`owners = "system"`、`most_recent = true` 过滤）。
- `instance_type` → `data "alibabacloudstack_instance_types"`。

**数据源 `ids` 参数防空字符串**：

当 data source 使用 `ids` 参数按 ID 列表查询时，**严禁传入空字符串**。如果变量可能为空，必须使用 `coalesce()` 兜底：

```hcl
# ❌ 错误：var.vpc_id 可能为空字符串，导致 ids = [""]
data "alibabacloudstack_vpc_vpcs" "existing" {
  ids = [var.vpc_id]
}

# ✅ 正确：使用 coalesce 兜底，空值时传入一个不存在的 ID
data "alibabacloudstack_vpc_vpcs" "existing" {
  ids = [coalesce(var.vpc_id, "vpc-nonexistent-placeholder")]
}
```

这样当变量为空时，data source 返回空列表（`vpcs` 长度为 0），不会报错。

**条件创建漂移问题（重要）**：

当实现"先查询再创建"逻辑时，**严禁使用 data source 查询结果来控制 `count`**。因为：
- 第一次 apply：data source 查不到 → 创建资源
- 第二次 apply：data source 查到了 → `count` 变为 0 → **删除资源**

**正确做法**：用变量是否为空来控制创建逻辑，data source 仅用于获取已存在资源的属性：

```hcl
# ❌ 错误：用 data source 结果控制 count，会导致第二次 apply 删除资源
locals {
  vpc_exists = length(data.alibabacloudstack_vpc_vpcs.existing.vpcs) > 0
}
resource "alibabacloudstack_vpc_vpc" "new" {
  count = local.vpc_exists ? 0 : 1  # 第二次 apply 会变成 0，删除资源！
}

# ✅ 正确：用变量是否为空控制 count，data source 仅用于获取属性
locals {
  vpc_exists = var.vpc_id != ""
  vpc_id     = local.vpc_exists ? var.vpc_id : alibabacloudstack_vpc_vpc.new[0].id
}
resource "alibabacloudstack_vpc_vpc" "new" {
  count = local.vpc_exists ? 0 : 1  # 只要 var.vpc_id 不变，count 就不会变
}
```

**关键原则**：`count` 的条件必须基于**用户输入的变量**，而非**数据源查询结果**。

#### 5.3 Provider 块（内容契约）

**重要：本 Skill 不生成 `provider "alibabacloudstack" {}` 块。** 
Provider 配置是用户的责任，应单独管理。

本 Skill 仅生成 `terraform { required_providers {} }` 块来声明 Provider 要求。

### 文件组织（强制）

生成的代码**必须**组织为恰好四个文件：

1. **`version.tf`**——Terraform 和 Provider 版本约束
   
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

2. **`variables.tf`**——资源配置的变量声明
   
   - 仅包含资源所需的变量（不包含 Provider 凭证）
   - 在适当的地方使用合理的默认值
   - 使用 `sensitive = true` 标记敏感变量

3. **`main.tf`**——主逻辑，包括：
   
   - 数据源（可用区、镜像、实例类型等）
   - 资源定义（VPC、ECS、RDS 等）
   - 所有基础设施组件

4. **`outputs.tf`**——输出定义
   
   - 导出重要的资源 ID、IP、名称等
   - 为每个输出包含描述性说明

### 版本约束规则

- Provider 版本约束：**必须使用 `< 3.19.0`** 作为专有云兼容性的默认上限。
- **不要**生成任何包含凭证、popgw_domain、protocol 或其他配置的 `provider "alibabacloudstack" {}` 块。
- 用户必须在自己的 `.tf` 文件中自行配置 provider 块。

**生成后验证**：

```bash
# 验证所有四个必需文件都存在
for file in version.tf variables.tf main.tf outputs.tf; do
  test -f <target-dir>/$file && echo "OK: $file exists" || echo "MISSING: $file"
done

# 验证 required_providers 包含 aliyun/alibabacloudstack 且具有正确的版本约束
grep -Rq 'alibabacloudstack.*source.*=.*"aliyun/alibabacloudstack"' \
  <target-dir>/version.tf && echo OK_SOURCE || echo BAD_SOURCE

grep -Rq 'version.*=.*"< 3.19.0"' <target-dir>/version.tf \
  && echo OK_VERSION_CONSTRAINT || echo BAD_VERSION_CONSTRAINT

# 确保本 Skill 没有生成任何 provider 块
! grep -Rq 'provider "alibabacloudstack"' <target-dir>/*.tf \
  && echo OK_NO_PROVIDER_BLOCK || echo UNEXPECTED_PROVIDER_BLOCK
```

所有检查必须通过。

#### 5.4 样式基线

- 2 空格缩进；块内 `=` 对齐；snake_case 语义化资源标签。
- 每个支持标签的资源都应携带非空的 `tags` 块。

#### 5.5 弃用字段审计——静态 grep 检查（强制）

在需要 `terraform` 之前运行——这是对刚编写的 HCL 的纯 grep 检查。对于本次生成中的每个资源，对照 `references/deprecated-fields-apsarastack.md` grep 项目并处理每种行类型：

- **rename** 行 → 用新字段名替换旧字段名。
- **split / soft-split** 行 → **不要**在父资源上写入内联字段；仅在需要时声明替换的子资源。
- **deprecated-no-replacement** 行 → 停止使用该字段，无替代品。

**审计后验证（bash grep——必须全部返回 OK）**：

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

**硬性门禁：必须在步骤 6 之前通过**——如果出现任何 `DEPRECATED:` 行，修复 HCL 并重新运行直到所有行返回 `OK:`。

### 步骤 6. 覆盖率检查 + 总结

**强制——无论生成结果如何都运行。**

**覆盖率检查。** 枚举生成的 HCL 中的资源块并与步骤 3 的草图比较。如果缺少任何草图行，返回步骤 5 并添加。

**总结模板**——以用户的语言打印：

```
已写入文件：
<path/to/file1>
<path/to/file2>
...

验证：待执行（进入步骤 7）

弃用路由：<如果重新路由：`<original_name>` → `<new_name>`；否则：无>

<可选：架构说明、设计决策、部署提示>
```

### 步骤 7. Terraform/OpenTofu 验证与修复（条件性执行）

**触发条件**：在文件生成完成后，检查用户本地是否安装了 `terraform` 或 `tofu`（OpenTofu）。

**检测工具可用性**：

```bash
# 检测 terraform
which terraform 2>/dev/null || command -v terraform 2>/dev/null

# 检测 tofu (OpenTofu)
which tofu 2>/dev/null || command -v tofu 2>/dev/null
```

**执行逻辑**：

1. **如果检测到 `terraform`**：
   - 在 `<target-dir>` 中执行 `terraform init`
   - 如果 `init` 成功，继续执行 `terraform validate`
     - 如果验证通过，记录成功状态
     - 如果验证失败，分析错误信息并自动修复 HCL 代码，然后重新验证（最多重试 3 次）
   - 如果 `init` 失败，输出提示信息并跳过验证步骤

2. **如果检测到 `tofu`（但未检测到 `terraform`）**：
   - 在 `<target-dir>` 中执行 `tofu init`
   - 如果 `init` 成功，继续执行 `tofu validate`
     - 如果验证通过，记录成功状态
     - 如果验证失败，分析错误信息并自动修复 HCL 代码，然后重新验证（最多重试 3 次）
   - 如果 `init` 失败，输出提示信息并跳过验证步骤

3. **如果两者都未检测到**：
   - 跳过验证步骤
   - 在总结中注明"验证：已跳过（未检测到 terraform 或 tofu）"
   - 提醒用户手动安装 Terraform 或 OpenTofu 以进行本地验证

**验证执行示例**：

```bash
cd <target-dir>

# 使用 terraform 验证
if command -v terraform &>/dev/null; then
  echo "检测到 Terraform，执行初始化..."
  terraform init
  
  if [ $? -eq 0 ]; then
    echo "初始化成功，执行验证..."
    terraform validate
    
    if [ $? -ne 0 ]; then
      echo "验证失败，分析错误并修复..."
      # 根据错误信息修复 HCL 文件
      # 重新验证（最多 3 次）
    fi
  else
    echo "初始化失败，跳过验证步骤"
    echo "提示：请检查网络连接、Provider 配置或手动执行 'terraform init'"
  fi
  
# 或使用 tofu 验证
elif command -v tofu &>/dev/null; then
  echo "检测到 OpenTofu，执行初始化..."
  tofu init
  
  if [ $? -eq 0 ]; then
    echo "初始化成功，执行验证..."
    tofu validate
    
    if [ $? -ne 0 ]; then
      echo "验证失败，分析错误并修复..."
      # 根据错误信息修复 HCL 文件
      # 重新验证（最多 3 次）
    fi
  else
    echo "初始化失败，跳过验证步骤"
    echo "提示：请检查网络连接、Provider 配置或手动执行 'tofu init'"
  fi
  
else
  echo "未检测到 Terraform 或 OpenTofu，跳过验证"
fi
```

**常见错误及修复策略**：

- **语法错误**（括号不匹配、缺少逗号等）→ 修正 HCL 语法
- **未定义的变量** → 检查 `variables.tf` 中是否声明
- **未定义的资源引用** → 检查资源名称是否正确
- **类型不匹配** → 修正变量类型或添加类型转换
- **Provider 配置缺失** → 提醒用户自行配置 provider 块（本 Skill 不生成）
- **init 失败** → 可能原因：网络问题、Provider 版本不存在、专有云环境配置问题；建议用户手动执行 init 并查看详细错误

**验证后更新总结**：

在步骤 6 的总结模板中更新验证状态：

```
验证：<成功/失败/已跳过>
- 工具：<terraform/tofu/无>
- 初始化：<成功/失败/跳过>
- 结果：<验证通过/验证失败，已修复 X 处错误/无法自动修复/init 失败跳过>
- 错误详情：<如有错误，列出关键错误信息>
```

**注意**：
- 本步骤会执行 `terraform init` 或 `tofu init` 来初始化 Provider 插件，但不会执行 `plan`、`apply` 等会修改基础设施的命令
- 如果 `init` 失败（通常由于网络问题或 Provider 配置问题），将跳过验证步骤并向用户报告
- 如果验证失败且无法自动修复，向用户报告具体错误并建议手动修复
- 最多重试 3 次验证-修复循环，超过后停止并报告剩余错误

## 参考资料

| 来源 | 何时读取 |
| --- | --- |
| Terraform MCP（`terraform::tool::search_providers`、`terraform::tool::get_provider_details`） | 步骤 4.1——每个资源的权威必需/可选参数（唯一来源） |
| `references/deprecated-fields-apsarastack.md`（本地） | 步骤 5.1——已知的字段级重命名 |
| `references/resource-patterns-apsarastack.md`（本地） | 步骤 5.1——产品特定惯用法 |
| `references/auth-and-network-apsarastack.md`（本地） | 用户参考凭证配置 |
