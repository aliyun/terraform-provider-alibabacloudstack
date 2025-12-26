# 3.18.22

## 新增

1. schedulerx应用组的编排能力（`alibabacloudstack_schedulerx2_app_group`）
2. schedulerx任务的编排能力（`alibabacloudstack_schedulerx2_job`）
3. schedulerx工作流的编排能力（`alibabacloudstack_schedulerx2_workflow`）
4. APFS可用区的查询能力（`alibabacloudstack_apfs_zones`）
5. APFS文件系统的编排能力（`alibabacloudstack_apfs_file_system`）
6. 密码机厂商及其产品的查询能力（`alibabacloudstack_hsm_vendors`）
7. 密码机HSM实例列表的查询能力（`alibabacloudstack_hsms`）
8. 密码机实例的编排能力（`alibabacloudstack_hsm_instances`）
9. 密码机集群的编排能力（`alibabacloudstack_hsm_clusters`）
10. 裸机管理密钥对的编排能力（`alibabacloudstack_bms_keypairs`）
11. 安骑士OSS扫描配置的编排能力（`alibabacloudstack_aqs_oss_scanconfigs`）
12. 安骑士主机防暴力破解防护规则的编排能力（`alibabacloudstack_aqs_anti_brute_force_rules`）
13. 安骑士Web防篡改配置的编排能力（`alibabacloudstack_aqs_web_locks`）
14. 密码服务密码机厂商及其产品的查询能力（`alibabacloudstack_cspprivate_hsm_vendors`）
15. 密码服务密码机HSM实例列表的查询能力（`alibabacloudstack_cspprivate_hsms`）
16. 密码服务密码机实例的编排能力（`alibabacloudstack_cspprivate_hsm_instance`）
17. 密码服务密码机集群的编排能力（`alibabacloudstack_cspprivate_hsm_group`）

## 下线

1. 环境下服务列表的查询能力(`environment_services_by_product`)

---

# 3.18.21

## 新增

1. 云解析服务转发域名的编排能力（`alibabacloudstack_dns_forward_domain`）
2. 云解析服务线路的编排能力（`alibabacloudstack_dns_line`）
3. 云解析私有域名的编排能力（`alibabacloudstack_dns_private_domain`）
4. 云解析私有域名解析配置记录的编排能力（`alibabacloudstack_dns_private_record`）
5. 云解析私有域名递归ACL策略的编排能力（`alibabacloudstack_dns_recursor_acl`）
6. 云解析私有域名线路的编排能力（`alibabacloudstack_dns_private_line`）
7. 云解析全局流量管理(GTM)实例的编排能力（`alibabacloudstack_dns_gtm_instance`）
8. 云解析全局调度地址池的编排能力（`alibabacloudstack_dns_gtm_addresspool`）
9. 云解析全局调度实例访问策略的编排能力（`alibabacloudstack_dns_gtm_access_strategy`）
10. CPFS文件系统资源的编排能力（`alibabacloudstack_cpfs_file_system`）
11. Prometheus v2监控实例的编排能力（`alibabacloudstack_prometheus_v2_instance`）
12. Prometheus v2配置告警通知接收人信息的编排能力（`alibabacloudstack_prometheus_v2_contact`）
13. Prometheus v2的通知组的编排能力（`alibabacloudstack_prometheus_v2_notify_group`）
14. Prometheus v2警告规则的编排能力（`alibabacloudstack_prometheus_v2_alert`）

## 下线

1. 下线ECS实例的计量数据的查询能力(`alibabacloudstack_ascm_metering_query_ecs`)

---

# 3.18.20

## 新增

1. API网关V2版本实例类型的查询能力（`alibabacloudstack_api_gateway_v2_instance_types`）
2. API网关V2版本实例的编排能力（`alibabacloudstack_api_gateway_v2_instance`）
3. API网关V2版本级联实例的编排能力（`alibabacloudstack_api_gateway_v2_cascade_instance`）
4. API网关V2版本级联链路的编排能力（`alibabacloudstack_api_gateway_v2_cascade_link`）
5. API网关V2版本K8s集群的导入能力（`alibabacloudstack_api_gateway_v2_k8s_cluster`）
6. API网关V2版本证书的编排能力（`alibabacloudstack_api_gateway_v2_certificate`）
7. API网关V2版本消费者的编排能力（`alibabacloudstack_api_gateway_v2_consumer`）
8. API网关V2版本域名的编排能力（`alibabacloudstack_api_gateway_v2_domain`）
9. API网关V2版本MCP服务的编排能力（`alibabacloudstack_api_gateway_v2_mcpservice`）
10. API网关V2版本路由组的编排能力（`alibabacloudstack_api_gateway_v2_route_group`）
11. API网关V2版本路由的编排能力（`alibabacloudstack_api_gateway_v2_route`）
12. API网关V2版本服务的编排能力（`alibabacloudstack_api_gateway_v2_service`）
13. API网关V2版本服务来源的编排能力（`alibabacloudstack_api_gateway_v2_service_source`）
14. API网关V2版本签名的编排能力（`alibabacloudstack_api_gateway_v2_signature`）
15. 数据总线的消息组的编排能力(`alibabacloudstack_datahub_kafka_group`)
16. 云原生多模数据库 Lindorm数据库实例的编排能力(`alibabacloudstack_lindorm_instance`)
17. 云原生多模数据库 Lindorm数据库实例规格的查询能力(`alibabacloudstack_lindorm_instance_types`)
18. 云原生多模数据库 Lindorm数据库LTS实例的编排能力(`alibabacloudstack_lindorm_lts_instance`)
19. 微消息队列 MQTT 版实例的编排能力(`alibabacloudstack_mqtt_instance`)
20. 微消息队列 MQTT 版Topic的编排能力(`alibabacloudstack_mqtt_topic`)
21. 微消息队列 MQTT 版Group的编排能力(`alibabacloudstack_mqtt_group`)
22. 跨云解析DNS域名的编排能力(`alibabacloudstack_universal_dns_domain`)
23. 跨云解析DNS解析配置的编排能力(`alibabacloudstack_universal_dns_record`)
24. 跨云解析DNS解析线路的编排能力(`alibabacloudstack_universal_dns_line`)

---

# 3.18.19

## 新增

1. ack 模板的编排能力（`alibabacloudstack_ack_template`）
2. Edas游道的编排能力（`alibabacloudstack_edas_swimming_lane`）
3. Edas游道组的编排能力（`alibabacloudstack_edas_swimming_lane_group`）
4. Hologram集群的查询能力（`alibabacloudstack_hologram_clusters`）
5. Hologram实例的编排能力（`alibabacloudstack_hologram_instance`）
6. Hologram实例备份规则的编排能力（`alibabacloudstack_hologram_instance_backup_policy`）

---

# 3.18.18

## 新增

1. oss集群的查询能力
2. edas k8s应用的伸缩规则配置管理能力

## 修复

1. 修复polardb在创建Tde版本得PGSQL时报错得问题。
2. 修复security_group_rule的网卡类型缺少默认值可能会导致失败的问题。
3. 修复edas_k8s_application_scaling_rule在错误参数配置下异常且没有正确报错的问题。
4. 修复slb_vservergroup的services的ids顺序导致无法终态问题
5. 修复oss_bucket容灾模式无法开启kms加密问题
6. 修复slb_listener 配置日志监听时变量的类型不正确问题

---

# 3.18.17

## 新增

1. Cen边界路由的网络实例连接的编排能力（alibabacloudstack_cen_transit_router_vbr_attachment）
2. Cen Connect的网络实例连接的编排能力（alibabacloudstack_cen_transit_router_connect_attachment）
3. Cen健康检查能力（alibabacloudstack_cen_vbr_health_check）
4. Cen Connect的网络实例连接支持配置对端（alibabacloudstack_cen_connect_peer）
5. Cen 组播成员的编排能力（alibabacloudstack_cen_transit_router_multicast_domain_member）
6. Cen 组播源的增加Connect的支持（alibabacloudstack_cen_transit_router_multicast_domain_source）

---

# 3.18.16

## 新增

1. Polardb 共享盘实例的编排能力（`alibabacloudstack_polardb_cluster_instance`）
2. Polardb 共享盘帐号编排能力（`alibabacloudstack_polardb_cluster_account`）
3. Polardb 共享盘数据库编排能力（`alibabacloudstack_polardb_cluster_database`）
4. Polardb 共享盘帐号权限配置能力（`alibabacloudstack_polardb_cluster_account_database_binding`）
5. Polardb 共享盘数据库代理配置能力（`alibabacloudstack_polardb_cluster_proxy`）
6. Polardb 共享盘备份策略配置能力（`alibabacloudstack_polardb_cluster_backup_policy`）
7. Polardb 共享盘实例规格查询能力（`alibabacloudstack_polardb_cluster_instance_types`）
8. Polardb 共享盘代理规格查询能力（`alibabacloudstack_polardb_cluster_proxy_types`）

---

# 3.18.15

## 变更

1. alibabacloudstack_slb_listener的`tls_cipher_policy`支持在环境中动态获取可用范围

## 修复

1. 修复多个单元测试用例
2. 修复资源集下oss bucket为空时创建失败的问题
3. 修复alibabacloudstack_security_group_rules无法根据`nic_type`进行过滤的问题
4. 修复alibabacloudstack_slb_loadbalancer为生效`AddressType`的问题
5. 修复alibabacloudstack_ons_instances因为接口返回数据类型的变化导致的查询失败的问题
6. 修复kvstore instance 应为clientToken导致签名验证失败的问题
7. 修复alibabacloudstack_dns_record不能终态的问题
8. 修复alibabacloudstack_cr_repos查询数据失败的问题
9. 修复alibabacloudstack_nas_file_system创建成功了`cluster_id`会发生变化导致无法终态的问题
10. 修复alibabacloudstack_ascm_user_group_role_binding因为asapi不可能导致的问题
11. **（不兼容）** 修复alibabacloudstack_ascm_user_role_binding因为asapi不可能导致的问题。注意`role_ids`类型修改为Set[INT]
12. 修复polardb MySql版本设置TDE加密后，无法获取到加密的KmsKey的问题

---

# 3.18.14

## 新增

1. polardb 新增支持ACL功能配置
2. 新增服务角色的授权资源 


## 修复

1. polardb PG数据库创建后，芯片类型为空
2. polardb database 绑定用户后无法查询到，也不能正常使用
3. mongodb_instance  ssl 开启状态检查修复
4. maxcompute_cu 资源适配3.18
5. maxcompute_user 资源适配3.18
6. maxcompute_project 资源适配3.18

---

# 3.18.13

## 新增

1. 新增支持企业云网络Cen的实例管理能力(cen_instance)
2. 新增支持企业云网络Cen的转发路由表(cen_transit_router_route_table)
3. 新增支持企业云网络Cen的转发路由表条目(cen_transit_router_route_entry)
4. 新增支持企业云网络Cen的转发路由链接管理VPC类型(cen_transit_router_vpc_attachment)
5. 新增支持企业云网络Cen的转发路由关联转发(cen_transit_router_route_table_association)
6. 新增支持企业云网络Cen的转发路由路由学习(cen_transit_router_route_table_propagation)
7. 新增支持企业云网络Cen的转发路由路由策略(alibabacloudstack_cen_route_map)
8. 新增支持企业云网络Cen的组播(cen_transit_router_multicast_domain)
9. 新增支持企业云网络Cen的组播交换机绑定(cen_transit_router_multicast_domain_association)
10. 新增支持企业云网络Cen的组播源(cen_transit_router_multicast_domain_source)
11. 新增支持Polardbx 2.0的实例规格查询能力(instance_type)
12. 新增支持Polardbx 2.0的实例管理能力(polardbx_instance)
13. 新增支持Polardbx 2.0的只读实例管理能力(polardbx_readonly_instance)
14. 新增支持Polardbx 2.0的用户帐号管理能力(polardbx_account)
15. 新增支持Polardbx 2.0的超级帐号（含三权分立）配置能力(polardbx_super_account)
16. 新增支持Polardbx 2.0的数据库表管理能力(polardbx_database)
17. 新增支持Polardbx 2.0的备份管理能力(polardbx_backup)
18. 新增支持Polardbx 2.0的备份规则配置能力(backup_policy)
19. 新增支持Polardbx 2.0的帐号和数据库绑定能力(polardbx_account_database_binding)
20. 新增支持Polardbx 2.0的实例日志引擎配置能力(polardbx_log_engine)
21. 新增支持Polardbx 2.0的读写分离配置能力(polardbx_read_write_splitting_config)

---

# 3.18.12

## 新增

1. 新增支持Mongodb帐号的管理能力（mongodb_account）
2. 新增支持Mongodb备份的管理能力(mongodb_backup)
3. 新增支持Mongodb分片实例cs节点网络连接的管理能力(mongodb_shardinginstance_csnode_address)
4. 新增支持Mongodb分片实例shard节点网络连接的管理能力(mongodb_shardinginstance_shardnode_address)
5. 新增支持Mongodb分片实例mongs节点网络连接的管理能力(mongodb_shardinginstance_mongonode_address)
6. 新增支持DRDS私有RDS实例的管理能力（drds_rds_instance）
7. 新增支持DRDS私有数据库的管理能力（drds_rds_instance）
8. 新增支持DRDS私有帐号的管理能力（drds_account）
9. 新增支持DRDS规格簇的查询能力（drds_instance_series）
10. 新增支持DRDS规格的查询能力（drds_instance_specifications）
11. 新增支持DRDS私有只读实例的管理能力（readonly_instance）
12. 新增支持POLARDB备份的管理能力（polardb_backup）
13. 新增支持POLARDB规格的查询能力（polardb_instance_types）
14. 新增支持POLARDB只读实例的管理能力（polardb_readonly_instance）
15. 新增支持POLARDB数据库代理（含读写分离）的管理能力（polardb_proxy）
16. 新增支持RDS规格的查询能力（rds_instance_types）
17. 新增支持RDS备份的管理能力（rds_backup）
18. 新增支持RDS代理的管理能力（rds_backup）
19. 新增加Mongodb副本实例的数据库连接地址的管理能力
20. 新增加Mongodb分片实例csnode的管理能力

## 变更

1. **（不兼容）** 新增加Mongodb分片实例各node支持设置描述，该字段在TF下为必填，且同一实例下需要唯一

---

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