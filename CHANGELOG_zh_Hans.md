# 3.18.29

## 新增

1. 日志服务集群的查询能力（`alibabacloudstack_log_clusters`）
2. BMCP 裸机密钥对的查询能力新增 BMCP 前缀别名（`alibabacloudstack_bmcp_keypairs`）
3. BMCP 裸机安全组的查询能力新增 BMCP 前缀别名（`alibabacloudstack_bmcp_security_groups`）
4. BMCP 裸机安全组规则的查询能力新增 BMCP 前缀别名（`alibabacloudstack_bmcp_security_group_rules`）

## 修复

1. 修复 `alibabacloudstack_mongodb_instance` 资源创建和读取的缺陷
2. 修复 Client `ConnectTimeout` 连接超时配置错误

## 变更

1. `ecs_dedicated_host_cluster` 数据源重命名为 `ecs_dedicated_host_clusters`
2. `alibabacloudstack_ack_cluster` 增强节点池相关字段处理能力，支持无默认节点池场景
3. `alibabacloudstack_cs_kubernetes_node_pool` 修复节点更新逻辑
4. 日志服务资源全量增强：`alibabacloudstack_log_project`、`alibabacloudstack_log_store`、`alibabacloudstack_log_store_index`、`alibabacloudstack_log_alert`、`alibabacloudstack_log_machine_group`、`alibabacloudstack_logtail_attachment`、`alibabacloudstack_logtail_config`
5. Provider 废弃参数描述更新：`sls_openapi_endpoint`、`asapi_endpoint`、`sls_endpoint`
6. 清理 `alibabacloudstack_datahub_*`、`alibabacloudstack_db_*`、`alibabacloudstack_dms_enterprise_*`、`alibabacloudstack_dns_*`、`alibabacloudstack_drds_*`、`alibabacloudstack_dts_*`、`alibabacloudstack_ecs_*` 等废弃别名资源

## 下线

1. 下线 ASCM 访问密钥资源的编排能力（`alibabacloudstack_ascm_access_key`）

---

# 3.18.28

## 修复

1. 修复 `alibabacloudstack_ack_cluster` 更新节点池 vswitch 配置时的引用错误逻辑
2. 修复 `alibabacloudstack_cs_kubernetes_node_pool` 多节点池场景下 refresh 操作失败的缺陷

## 变更

1. OSS 资源全量改造，`alibabacloudstack_oss_bucket`、`alibabacloudstack_oss_bucket_object`、`alibabacloudstack_oss_bucket_kms`、`alibabacloudstack_oss_bucket_quota` 摆脱对 ASAPI 的依赖，改用云产品自建网关获取 Endpoint
2. `alibabacloudstack_cs_kubernetes_node_pool` 支持通过已有 K8s 集群创建节点池，新增 `platform = "Custom"` 支持
3. `alibabacloudstack_ack_cluster` 简化集群更新逻辑，修复状态迁移处理
4. 修复 Client `ReadTimeout` 和 `ConnectTimeout` 超时单位错误（从 Hour 修正为 Millisecond）

---

# 3.18.27

## 新增

1. BMCP裸机计算节点的编排能力（`alibabacloudstack_bmcp_cluster`）
2. BMCP裸机计算节点的查询能力（`alibabacloudstack_bmcp_clusters`）
3. BMS裸机计算节点规格的查询能力（`alibabacloudstack_bms_machinetypes`）
4. ASCM RAM服务角色的查询能力（`alibabacloudstack_ascm_ram_service_roles`）
5. 增强虚拟专有网络（EVPC）实例的查询能力（`alibabacloudstack_evpc_vpcs`）

## 下线

1. 下线 MaxCompute 用户的编排能力（`alibabacloudstack_maxcompute_user`）

## 修复

1. 修复 `alibabacloudstack_oss_bucket` 生命周期、版本控制等属性的读写逻辑缺陷
2. 修复 `alibabacloudstack_oss_bucket` CORS 规则更新时缺少校验导致的错误
3. 修复 `alibabacloudstack_image_export` 导出镜像时未等待任务完成的 bug
4. 修复 `alibabacloudstack_vpc_vswitch` SDK 兼容性问题

## 变更

1. `alibabacloudstack_oss_bucket` 大规模重构，重构资源管理内部实现
2. `alibabacloudstack_oss_bucket_object` 简化对象操作逻辑
3. `alibabacloudstack_oss_bucket_kms` 重构 KMS 配置逻辑

---

# 3.18.26

## 新增

1. BMCP裸机计算节点镜像的查询能力（`alibabacloudstack_bmcp_images`）
2. BMCP裸机计算节点规格的查询能力（`alibabacloudstack_bmcp_machinetypes`）
3. BMCP裸机计算节点的查询能力（`alibabacloudstack_bmcp_nodes`）
4. BCMP裸机计算节点密钥对的编排能力（`alibabacloudstack_bcmp_keypair`）
5. BCMP裸机计算节点密钥对的查询能力（`alibabacloudstack_bcmp_keypairs`）
6. BCMP裸机计算网络安全组的编排能力（`alibabacloudstack_bcmp_security_group`）
7. BCMP裸机计算网络安全组的查询能力（`alibabacloudstack_bcmp_security_groups`）
8. BCMP裸机计算网络安全组规则的编排能力（`alibabacloudstack_bcmp_security_group_rule`）
9. BCMP裸机计算网络安全组规则的查询能力（`alibabacloudstack_bcmp_security_group_rules`）
10. 增强虚拟专有网络（EVPC）实例的编排能力（`alibabacloudstack_evpc_evpc`）
11. 增强虚拟专有网络（EVPC）实例的查询能力（`alibabacloudstack_evpc_evpcs`）

## 修复

1. 修复 `alibabacloudstack_ascm_organization` 创建组织时因同名称已存在导致的重复创建和超时问题
2. 修复 `alibabacloudstack_ascm_organization` 更新组织名称时发送冗余请求的问题
3. 修复 `alibabacloudstack_ascm_organization` Read 方法中返回结构体字段访问错误
4. 修复 `alibabacloudstack_cr_namespace` 创建时缺少必要参数导致失败的问题
5. 修复 `alibabacloudstack_cr_namespace` 更新时缺少 NamespaceName 参数的问题
6. 修复 `alibabacloudstack_rds_dbconnection` 数据库连接地址创建和更新逻辑缺陷
7. 修复 `alibabacloudstack_kvstore_instance` 测试用例并优化实例编排逻辑
8. 修复 `alibabacloudstack_kvstore_connection` SSL 连接字段处理问题
9. 修复 `alibabacloudstack_oos_execution` 和 `alibabacloudstack_oos_template` 执行与模板相关逻辑缺陷
10. 修复 `alibabacloudstack_ascm_organization` 删除逻辑中重试条件错误

## 变更

1. `alibabacloudstack_kvstore_instance` 恢复 `cpu_type` 字段可用性（移除 Deprecated 标记）
2. `alibabacloudstack_kvstore_instance` 优化 `node_type` 字段的 diff 处理逻辑
3. 批量修复十余个 DataSource 增加 `shared` 参数支持查询共享资源（ECS、OSS、SLB、CR、NAT、VPC 等）

---

# 3.18.25

## 新增

1. mqtt集群的查询能力（`alibabacloudstack_mqtt_clusters`）
2. ons集群的查询能力（`alibabacloudstack_ons_clusters`）
3. ascm组织`alibabacloudstack_ascm_organization`支持获取PK信息

## 修复
1. Datawork Remain无法正常使用的问题
2. Datawork Connection无法正常使用的问题

---
