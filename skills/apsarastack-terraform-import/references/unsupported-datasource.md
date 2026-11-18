# 不支持 Datasource 覆盖的资源清单

> 基于 `terraform-provider-apsarastack` v3.18xR 分支的 `provider.go` 中注册的所有 datasource 与 `references/api-commands.md` 中的 CLI 命令做完整对比。
> 以下资源在发现阶段 **无法使用 Terraform datasource** 完成，必须回退到 `apsara` CLI 或 ASAPI 直接调用。

---

## 一、完全没有 Datasource 的产品（CLI 必填）

| 产品 | CLI 命令示例 | 对应 Terraform 资源 | 备注 |
|------|-------------|-------------------|------|
| **FC (Function Compute)** | `apsara fc GET /services` | `alibabacloudstack_fc_service` | 无 service/function datasource |
| | `apsara fc GET /services/<svc>/functions` | `alibabacloudstack_fc_function` | |
| **ALB (应用型负载均衡)** | `apsara alb list-load-balancers` | `alibabacloudstack_alb_loadbalancer` | 无 ALB 任何 datasource |
| | `apsara alb list-listeners` | `alibabacloudstack_alb_listener` | |
| | `apsara alb list-server-groups` | `alibabacloudstack_alb_server_group` | |
| **DNS Domain Group** | `apsara alidns describe-domain-groups` | — | 仅分组信息，非核心资源 |

---

## 二、有列表 Datasource 但缺少子资源/详情 Datasource

### 2.1 OSS Bucket（列表有，属性详情无）

| 属性 | CLI 命令 | 说明 |
|------|---------|------|
| Bucket stat | `apsara ossutil stat oss://<bucket>` | 存储量、对象数等 |
| ACL | `apsara ossutil bucket-acl oss://<bucket>` | 访问权限 |
| Versioning | `apsara ossutil bucket-versioning oss://<bucket>` | 版本控制状态 |
| Lifecycle | `apsara ossutil lifecycle --method get oss://<bucket>` | 生命周期规则 |
| CORS | `apsara ossutil cors --method get oss://<bucket>` | 跨域配置 |
| Encryption | `apsara ossutil bucket-encryption --method get oss://<bucket>` | 服务端加密 |

**已有 Datasource**: `alibabacloudstack_oss_buckets`（仅返回 bucket 列表和基础属性）
**缺失 Datasource**: 单个 bucket 的完整属性查询

### 2.2 RDS 子资源（实例列表有，子资源无）

| 子资源 | CLI 命令 | 对应 Terraform 资源 | 说明 |
|--------|---------|-------------------|------|
| Database 列表 | `apsara rds describe-databases` | `alibabacloudstack_db_database` | 无 db_databases datasource |
| Account 列表 | `apsara rds describe-accounts` | `alibabacloudstack_db_account` | 无 db_accounts datasource |
| IP 白名单 | `apsara rds describe-db-instance-ip-array-list` | `alibabacloudstack_db_instance` 的 `security_ips` | 字段在实例详情中，但独立查询无 datasource |
| 连接地址 | `apsara rds describe-db-instance-net-info` | `alibabacloudstack_db_connection` | 无 db_connections datasource |

**已有 Datasource**: `alibabacloudstack_db_instances` / `alibabacloudstack_rds_dbinstances`
**缺失 Datasource**: db_databases, db_accounts, db_connections 等子资源

> 注意：PolarDB 有 `alibabacloudstack_polardb_databases` 和 `alibabacloudstack_polardb_accounts` datasource，但 RDS MySQL 没有。

### 2.3 Redis (KVStore) 子资源

| 子资源 | CLI 命令 | 对应 Terraform 资源 | 说明 |
|--------|---------|-------------------|------|
| IP 白名单 | `apsara r-kvstore describe-security-ips` | `alibabacloudstack_kvstore_instance` 的 `security_ips` | 字段在实例属性中，无独立 datasource |

**已有 Datasource**: `alibabacloudstack_kvstore_instances`

### 2.4 MongoDB (DDS) 子资源

| 子资源 | CLI 命令 | 对应 Terraform 资源 | 说明 |
|--------|---------|-------------------|------|
| IP 白名单 | `apsara dds describe-security-ips` | `alibabacloudstack_mongodb_instance` 的 `security_ip_list` | 字段在实例属性中，无独立 datasource |

**已有 Datasource**: `alibabacloudstack_mongodb_instances`

### 2.5 ACK 子资源

| 子资源 | CLI 命令 | 对应 Terraform 资源 | 说明 |
|--------|---------|-------------------|------|
| Node Pool 列表 | `apsara cs GET /clusters/<id>/nodepools` | `alibabacloudstack_cs_kubernetes_node_pool` | 无 ack_nodepools datasource |

**已有 Datasource**: `alibabacloudstack_ack_clusters`, `alibabacloudstack_ack_templates`

---

## 三、有 Datasource 但功能不完整（需 CLI 补充）

| 资源 | Datasource | 缺失能力 | CLI 补充命令 |
|------|-----------|---------|-------------|
| Security Group Rules | `alibabacloudstack_security_group_rules` | 可能无法按 ingress/egress 分别过滤 | `apsara ecs describe-security-group-attribute --direction ingress` |
| ECS Images | `alibabacloudstack_ecs_images` | 可能需要 owner 过滤 | `apsara ecs describe-images --image-owner-alias self` |
| ESS Scaling Configuration | `alibabacloudstack_ess_scaling_configurations` | 需指定 scaling_group_id | `apsara ess describe-scaling-configurations --scaling-group-id` |

---

## 四、发现阶段策略建议（与 SKILL.md Phase 4 对齐）

```
Phase 4 资源发现流程（Datasource First, CLI Fallback）:

Step 0: 分类资源类型
  ├── Group A (有 datasource): VPC, VSwitch, ECS, SecurityGroup, EIP, NAT, SLB,
  │                            OSS(list), RDS(instance), Redis, MongoDB, DNS, ASCM
  └── Group B (无 datasource): FC, ALB, ACK NodePool, RDS DB/Account, OSS properties

Step 1: Group A 用 Terraform Datasource 发现
  ├── 创建 .import/discovery/discovery.tf（参考 examples/discovery.tf）
  ├── terraform init && terraform refresh
  ├── terraform output -json > discovery-output.json
  └── python3 tests/parse-discovery.py discovery-output.json

Step 2: Group B 非空检查 + CLI 依赖处理
  ├── Group B 为空 → 跳过 Step 3/4，直接用 Group A 报告进入 Phase 5
  └── Group B 非空 → 检查 apsara-cli

Step 3: apsara-cli 安装检查
  ├── 已安装 → 继续 Step 4
  ├── 未安装 → 尝试自动安装（brew / curl）
  │   ├── 安装成功 → 继续 Step 4
  │   └── 安装失败 → 向用户展示 Group B 清单，提示必须安装 CLI，
  │                 提供手动安装指导，停止所有后续任务
  └── 停止后等待用户手动安装 CLI，再重新触发

Step 4: Group B 用 apsara-cli 发现
  └── 执行 references/api-commands.md 中的对应命令

Step 5: 合并报告
  └── 合并 Group A (datasource) + Group B (CLI)，生成统一发现报告
```

---

## 五、核心导入资源覆盖率统计

| 类别 | 核心资源数 | Datasource 覆盖 | CLI 回退 | 覆盖率 |
|------|-----------|----------------|----------|--------|
| 网络 (VPC/SLB/EIP/NAT/VPN) | 8 | 8 | 0 | 100% |
| 计算 (ECS/Disk/Image/KeyPair) | 5 | 5 | 0 | 100% |
| 存储 (OSS Bucket) | 1 | 0.5 | 0.5 | 50% |
| 数据库 (RDS/Redis/MongoDB) | 4 | 2 | 2 | 50% |
| 负载均衡 (SLB) | 1 | 1 | 0 | 100% |
| DNS | 2 | 2 | 0 | 100% |
| 身份 (ASCM) | 3 | 3 | 0 | 100% |
| 容器 (ACK) | 1 | 0.5 | 0.5 | 50% |
| 函数计算 (FC) | 2 | 0 | 2 | 0% |
| **合计** | **27** | **22** | **5** | **81%** |

> 注："0.5" 表示列表可用 datasource，但子资源/属性需 CLI 补充。

---

*本清单基于 provider v3.18xR 分支代码，后续版本可能新增 datasource。*
