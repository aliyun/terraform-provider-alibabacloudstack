# 3.18.31

## 新增

1. `alibabacloudstack_edas_k8s_application` 支持 `cr_instance_id` 字段，镜像部署时支持指定容器镜像服务企业版实例
2. BMCP 裸机节点查询能力新增 `instance_alias` 和 `instance_name` 返回字段（`alibabacloudstack_bmcp_nodes`）
3. `alibabacloudstack_image_import` 的 `platform` 字段支持更多操作系统类型（RedHat、Fedora、Anolis、AlmaLinux、Rocky Linux、Kylin、UOS、FreeBSD、Windows Server 2016/2019/2022 等）

## 修复

1. 修复 `alibabacloudstack_edas_k8s_application` 未按 `package_type` 正确回写 `image_url` / `package_url` 的缺陷
2. 修复 `alibabacloudstack_edas_k8s_application` 创建、部署、删除时遇到限流错误直接失败的问题，增加自动重试
3. 修复 `alibabacloudstack_edas_k8s_application` 删除时应用不存在或已删除导致报错的问题，现可幂等完成
4. 修复 `alibabacloudstack_cr_ee_repo` 更新时误读 `repo_type` 导致更新失败的缺陷
5. 修复 `alibabacloudstack_bmcp_nodes` 数据源 `region_id` 无法正确获取的缺陷
6. 修复 VPC 资源标签操作（绑定/解绑/查询标签）请求参数传递错误导致操作失败的缺陷
7. 修复 CMS 服务端点模板配置错误
8. 移除硬编码的 OSS 服务端点模板，改为按实际环境动态解析

---

# 3.18.30

## 修复

1. 修复 `alibabacloudstack_log_store` 创建后分片未就绪导致状态刷新异常的缺陷
2. 修复 `alibabacloudstack_oss_bucket` 在 `ReBindResourceGroup` 失败时的资源泄漏问题
3. 修复 `alibabacloudstack_log_project` 对额外错误码（BadRequest、SCMG 调用失败）的处理逻辑

## 变更

1. `alibabacloudstack_ascm_ram_roles` 数据源 ID 匹配逻辑修正为 `RoleName:ID` 格式
2. 内部 Token 生成逻辑优化，限制 ClientToken 长度不超过 64 字符

---

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

# 3.18.24

## 新增

1. ots集群的查询能力（`alibabacloudstack_ots_clusters`）
2. polardb读写分析的配置能力（`alibabacloudstack_polardb_readwrite_splitting_connection`）
3. datawrok流程的编排能力（`alibabacloudstack_data_works_business`）
3. ecs宿主机类型的查询能力（`alibabacloudstack_ecs_dedicated_host_types`）
4. polardb的参数组的配置能力（`alibabacloudstack_polardb_parameter_group`）
5. quickbi用户组用户的绑定能力（`alibabacloudstack_quick_bi_user_group_user`）
6. csb服务的创建能力（`alibabacloudstack_csb_service`）
7. datawork文件类型的查询能力（`alibabacloudstack_dataworks_file_types`)
8. datawork文件的编排能力（`alibabacloudstack_data_works_file`）
9. datawork智能基线的编排能力（`alibabacloudstack_data_works_baseline`）
10. alikafka实力增加自定义端口能力
11. 增加ecs快照策略绑定磁盘Id的能力

## 修复
1. Ots相关资源编排能力的重构
2. Datawork资源不可编排的问题


## 废弃

1. 废弃arms的告警规则的编排能力（`alibabacloudstack_arms_prometheus_alert_rule`)
2. 废弃arms的派遣规则的编排能力（`alibabacloudstack_arms_dispatch_rule`）
3. 废弃cms站点监控的编排能力（`alibabacloudstack_cms_site_monitor`）
4. 废弃cms告警联系人的编排能力（`alibabacloudstack_cms_alarm_contact`）
5. 废弃cms告警联系人组的编排能力（`alibabacloudstack_cms_alarm_contact_group`）
6. 废弃mongodb备份计划的编排能力（`alibabacloudstack_dbs_backup_plan`）
7. 废弃quickbi的工作空间的编排能力（`alibabacloudstack_quick_bi_workspace`）
8. 废弃ots服务的查询能力（`alibabacloudstack_ots_service`）

---

# 3.18.23

## 新增

1. ACR企业版的镜像生命周期规则的编排能力（`alibabacloudstack_cr_ee_attestor_lifecycle_rule`)
2. 资源集和用户关系绑定能力的编排能力（`alibabacloudstack_ascm_resource_group_user_attachment`）
3. 云防火墙中的IP或端口地址簿的编排能力（`alibabacloudstack_cloudfw_address_book`)
4. VPC网络流量的访问控制规则的编排能力（`alibabacloudstack_cloudfw_vpc_control_policy`)
5. OSS单隧道的编排能力（`alibabacloudstack_oss_single_tunnel`)
6. redis参数模板的编排能力（`alibabacloudstack_kvstore_parameter_group`)
7. gpdb备份规则的编排能力（`alibabacloudstack_gpdb_backup_policy`)
8. nas目录配额的编排能力（`alibabacloudstack_nas_dir_quota`)
9. nas命名空间的编排能力（`alibabacloudstack_nas_namespace`)
10. nas统一命名空间文件存储映射的编排能力（`alibabacloudstack_nas_namespace_filesystem_attachment`)
11. nas命名空间文件系统挂载的编排能力（`alibabacloudstack_nas_namespace_mount_target`)
12. NAS跨域挂载编排的编排能力（`alibabacloudstack_nas_namespace_group`)
13. gpdb实例规格的查询能力（`alibabacloudstack_gpdb_instance_types`)
14. adb集群规格的查询能力（`alibabacloudstack_adb_cluster_types`)
15. edas k8s集群的查询能力（`alibabacloudstack_edas_k8s_clusters`)

## 修复
1. 批量修复datasource资源不能使用`ids`正确过滤的问题
2. 批量移除datasource资源不能正确过滤的`tags`字段，移除reousrce中不支持的`tags`字段
3. 批量修复多除代码缺陷，并修正测试用例


## 变更
1. `alibabacloudstack_expressconnect_physicalconnection`的`bandwidth`字段类型从字符串改为整型
2. `alibabacloudstack_dns_domain`删除`dns_servers`、`group_id`、`lang`等原先不支持的字段
3. `alibabacloudstack_ecs_deployment_set`删除不支持的`on_unable_to_redeploy_failed_instance`字段

## 废弃

1. dns域名组的编排能力（`alibabacloudstack_alidns_domaingroup`)
2. apigateway服务的查询能力`alibabacloudstack_api_gateway_service`)
3. ESS的生命周期管理能力（`alibabacloudstack_ess_lifecycle_hook`)
4. ESS的虚拟服务组的编排能力（`alibabacloudstack_ess_scalinggroup_vserver_groups`）
5. Kms密文能力（`alibabacloudstack_kms_secret`）

---

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