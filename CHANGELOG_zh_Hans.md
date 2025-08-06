# 3.18.11

## 修复

1. 修复 alibabacloudstack_log_alert 没有自动创建图表的问题

## 新增

1. alibabacloudstack_log_alert 支持配置 `webhook`
2. 实现 alibabacloudstack_acm_configuration

---

# 3.18.10

## 新增

1. 新增 VPN 资源 vpn_pbr_route_entry（策略路由表）
2. 新增 VPN 资源 ssl_vpnserver（SSL 服务端）
3. 新增 VPN 资源 ssl_vpn_client_cert（SSL 客户端）
4. 新增 ExpressConnect 资源 bgp_group（BGP 组）
5. 新增 ExpressConnect 资源 bgp_peer（BGP 邻居）
6. 新增 ExpressConnect 资源 vbr_ha（快速倒置组）
7. 新增 ExpressConnect 资源 vbr_pconn_association（多物理专线）
8. 新增 ExpressConnect 资源 bgp_network（BGP 网络）
9. 新增 NAS 资源 nas_lifecycle_policy（生命周期规则）

## 修复

1. **（不兼容）** NAS 资源可用区查询，返回结构发生较大变化

## 下线

1. **（不兼容）** NAS 资源下线 protocols 查询，可用新版本 Zone 查询替代

---

# 3.18.9

## 新增

1. 新增 NatGateway 资源 BandwidthPackage（带宽包）
2. 新增 SLB 资源 AccessLog（访问日志开启）
3. 新增 VPC 资源 vpc_ipv6_isps（可用 IPv6 网段）的查询能力，vpc 资源支持 `ipv6_cidr_blocks` 绑定多个资源
4. 新增 VPC 资源 VSwitchNetworkAclAttachment（vswitch acl 资源绑定）能力
5. 新增 VPC 资源 HaVip（高可用浮动 IP）

## 废弃

1. SLB Listener `logs_download_attributes` 属性标记废弃，请用 AccessLog 资源代替

---

# 3.18.8

## 新增

1. 新增 EBS 资源 DiskReplicaPair（云盘异步复制）
2. 新增 EBS 资源 DiskReplicaGroup（一致性复制组）
3. 新增 ECS 资源 SnapshotGroup（快照一致性组）
4. 新增 ECS 资源 DedicatedHostCluster（宿主机组）
5. 新增 ECS 资源 Invocation（命令执行）
6. 新增 ECS 资源 DhcpOptionsSet（DHCP 选项集）
7. 新增 ESS 资源 ScheduledTask 对属性 `scaling_group_id` 的支持

## 下线

1. 下线 LaunchTemplate 的能力，相关功能 ASCM 页面未开放

---

# 3.18.7

## 新增

1. flink_namespace 创建和查询能力
2. 新增 PolarDB Instance 配置 `tag` 的能力
3. 新增 CS K8s 集群的磁盘加密能力
4. 允许 EDAS K8sApp 绑定存量 SLB 实例

## 修复

1. 修复 PolarDB 修改内核参数的能力

## 废弃

1. 所有 DataSource 的 `output_file` 属性已标记废弃，计划 3.19.0 下线。可用 `local_file` Provider 替代。

---

# 3.18.6

## 新增

1. edas_k8s_application SLB 绑定的 update、delete 功能，批量绑定功能
2. edas_k8s_application 支持设置 `host_aliases`

## 修复

1. edas_k8s_application 变更内存和 CPU 的能力
2. ots_instance 创建失败的问题
3. 修复 Elasticsearch Instance 的创建和删除能力

---

# 3.18.5

## 新增

1. 新增 alikafka_instance 的创建和查询能力
2. Redis 实现开启和关闭 SSL 功能
3. SLB 支持经典网络指定 `address`
4. oss_bucket 支持设置 `tags`

## 变更

1. alibabacloudstack_ascm_user_group_resource_set_binding 的主键从 `resourceSetId` 修改为 `resourceSetId:userGroupId:ascmRoleId`。apply 时资源会删除后重建，业务影响范围可控，但不影响模板兼容性
2. 标记 alibabacloudstack_ascm_user_group_role_binding 为废弃，其功能已由 alibabacloudstack_ascm_user_group 的 `role_ids` 包含
3. 标记 alibabacloudstack_ascm_user_role_binding 为废弃，其功能已由 alibabacloudstack_ascm_user 的 `role_ids` 包含

## 修复

1. PolarDB 使用 PostgreSQL 引擎时，TDE 启动失败问题
2. 通过环境变量设置 OSS Bucket 时失败的问题
3. alibabacloudstack_polardb_dbinstance 在非 MySQL 引擎时，TDE 开启错误问题
4. 修复 PolarDB Account 在创建时因资源未就绪导致的失败
5. 修复 PolarDB Instance 在 TDE 未开启时，因 `encrypt_algorithm` 无法终态的问题
6. 修复 CS Cluster 集群的 NodePool 无法缩扩节点的问题

---

# 3.18.4

## 新增

1. 组织资源集过滤修改为精确匹配
2. 修复 vpngateway 资源创建
3. 部分资源平滑迁移测试
4. 新增资源 PolarDB_instance
5. 新增资源 PolarDB_database
6. 新增资源 PolarDB_account
7. 新增资源 PolarDB_backup_policy
8. 新增资源 PolarDB_dbconnection
9. 新增资源 PolarDB_Zone

## 修复

1. edas_k8s_service：更换 read 接口，补充部分只读属性，完善相关文档
2. cr_repo 将 name 长度调整至 64，与页面能力保持一致
3. 修复部分文档页面错误问题
4. 修复 CR-EE 的 namespace 和 repository 在 popgw 模式下调度失败的问题
5. 修复 CR-EE 的 instance、namespace、repo 查询失败的问题
6. 创建 edas_k8s_app 时，新增 `cree_repo_id` 传参
7. 修复 PolarDB_instance 同时开启 TDE 和 SSL 时报错的问题
8. 修复 edas_k8s_cluster 二次 apply 时出现 ForceNew 的问题
9. 修复 edas_k8s_service、edas_k8s_app read 缺陷导致资源无法终态的问题

---

# 3.18.3

## 新增

1. edas_k8s_app 创建时支持设置 PVC 挂载、本地挂载、配置设置
2. 新增资源 edas_k8s_service
3. 新增资源 edas_namespace

---

# 3.18.2

## 新增

1. Terraform 支持创建/修改/删除以下资源：
   - alibabacloudstack_bastionhost_instance（堡垒机实例）
   - alibabacloudstack_waf_instance（WAF 实例）
   - alibabacloudstack_vpn_gateway（VPN 网关）

## 修复

1. 修复 vpn_gateway 创建失败的问题

---

# 3.18.1

## 新增

1. 云盘加密方式支持

## 修复

1. ECS 更换镜像时系统盘 Tag 丢失的问题
2. 修复 Disk 加密时默认加密方式不传参导致无 `tag` 的问题
3. 修复使用加密快照创建 Disk 时失败的问题
4. 修复 SLS Alarm 未设置 `mute_util` 时编排失败的问题
5. 修复 alibabacloudstack_datahub_project 在 popgw 模式下 `comment` 不可修改的问题
6. 修复 alibabacloudstack_datahub_topic 在 popgw 模式下的问题
7. 修复 alibabacloudstack_datahub_subscription 在 popgw 模式下 `comment` 不可修改的问题

---

# 3.18.0

> 基于 3.16.1 版本分支

## 修复

1. 修改 alikafka 的预置 popgw 域名格式，适配 318x（与专有云 3.16.2 版本不兼容）

## 移除

1. 移除 asapi domain 的配置方式 `domain` 和 `force_use_asapi`