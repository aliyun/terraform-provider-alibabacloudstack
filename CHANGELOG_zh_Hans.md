# 3.16.24

## 新增

1. MongoDB 实例类型的查询能力（`alibabacloudstack_mongodb_instance_types`）
2. 多个数据源支持查询从其他组织共享的资源能力，新增 `shared` 参数

## 修复

1. 修复 `alibabacloudstack_ascm_organization` 资源创建/更新/删除逻辑，提升稳定性
2. 修复 `alibabacloudstack_rds_dbconnection` 支持 `schema_version` 参数配置
3. 修复 `alibabacloudstack_cr_namespace` 创建参数不正确导致创建失败的问题
4. 修复 `alibabacloudstack_ascm_ram_policy` 描述字段处理逻辑
5. 修复 `alibabacloudstack_ascm_organizations` 数据源支持 `ids` 过滤和新增 `primary_key` 字段

---

# 3.16.23

## 新增

1. 创建Polardb时支持配置参数模板。

---

# 3.16.22

## 修复

1. 修复logtail_config 资源的 plugin 输入类型。

---

# 3.16.21

## 变更

1.  oss object调度不再依赖用户配置oss cluster的endpoint地址，terraform会自动获取

## 修复

1. oss bucket编排时accountinfo错误导致编排失败的问题
2. cms_alarm 资源创建时没有正确传输webhook属性的问题

---

# 3.16.20

## 修复

1. 修复polardb在创建Tde版本得PGSQL时报错得问题。
2. 修复security_group_rule的网卡类型缺少默认值可能会导致失败的问题。
3. 修复edas_k8s_application_scaling_rule在错误参数配置下异常且没有正确报错的问题。

---

# 3.16.19

## 新增

1. oss集群的查询能力
2. edas k8s应用的伸缩规则配置管理能力

## 修复

1. 修复slb_vservergroup的services的ids顺序导致无法终态问题
2. 修复oss_bucket容灾模式无法开启kms加密问题
3. 修复slb_listener 配置日志监听时变量的类型不正确问题

---

# 3.16.18

## 新增

1. slb_listener 的`tls_cipher_policy`类型支持动态获取。 

## 修复

1. polardb MySql版本设置TDE加密后，无法获取到加密的KmsKey
2. mongodb 支持开启和关闭Sql审计

---

# 3.16.17

## 修复

1. polardb PG数据库创建后，芯片类型为空
2. polardb database 绑定用户后无法查询到，也不能正常使用
3. mongodb_instance  ssl 开启状态检查修复

## 新增

1. polardb 新增支持ACL功能配置
2. 新增服务角色的授权资源 

---

# 3.16.16

## 修复

1. 修复ecs description长度限制是256字符而不是256个字符的问题
2. 修复edas_k8s_application package url属性因为回读导致不终态的问题

---

# 3.16.15

## 修复：

1. slb loadbalancer增加network_type属性带来的不兼容
2. 修复alibabacloudstack_log_alert没有自动创建图表的问题

## 增加

1. alibabacloudstack_log_alert支持配置webhook
2. 实现alibabacloudstack_acm_configuration

---

# 3.16.14

## 修复

1. 修复polardb修改内核参数的能力

## 增加

1. 增加polardb instance配置tag的能力
2. 允许edas k8sapp 绑定存量slb实例

---

# 3.16.13

## 修复

1. 修复elasticsearch instance的创建和删除能力

---

# 3.16.12

## 新增

1. edas_k8s_application slb绑定的update，delete功能， 批量绑定功能
2. edas_k8s_application 支持设置 host_aliases

## 修复

1. edas_k8s_application 变更 内存和cpu的能力
2. ots_instance创建失败的问题

---

# 3.16.11

## 新增

1. redis实现开启和关闭ssl功能
2. slb支持指定address的经典网络
3. oss_bucket 设置tags

## 修复

1. 修复ploardb account在创建时因为资源可能没有就绪导致的失败
2. 修复poladb instance在tde没有开启时，资源因为encrypt_algorithm无法终态的问题
3. 修复kvstore_instnaceclass无法正常查询redis实例规格的问题
4. 修复cs_cluster集群的nodepool无法缩扩节点的问题

---

# 3.16.10

## 新增

1. 拉起所有资源通过id import的能力，73项关键资源的Import功能进行修复验证修复

## 变更

1. alibabacloudstack_ascm_user_group_resource_set_binding的主键从resourceSetId修改为resourceSetId:userGroupId:ascmRoleId。apply时资源会删除后重建，业务影响范围可控，但不影响模板兼容性
2. alibabacloudstack_ascm_user_group_role_binding过时，用alibabacloudstack_ascm_user_group中的role_ids包含了相关的功能
3. alibabacloudstack_ascm_user_role_binding过时，用alibabacloudstack_ascm_user中的role_ids包含了相关的功能

## 修复

1. alibabacloudstack_polardb_dbinstance 在非mysql引擎时，Tde开启错误问题
2. 通过环境变量设置OSS bucket时失败的原因

---

# 3.16.9

## 新增
1. cs_kubernetes 创建时允许设置 storage_set
## 修复
1. polardb使用pgsql引擎时，tde启动失败问题
2. 补充部分资源文档信息。

---

# 3.16.8

## 修复

1. 补充文档
2. 创建edas_k8s_app时，增加一个cree_repo_id传参
3. 修复polardb_instance 同时开启Tde和ssl时报错的问题
4. 修复edas_k8s_cluster 二次apply时出现ForceNew的问题
5. 修复edas_k8s_service, edas_k8s_app read缺陷导致资源无法正常终态的问题

---

# 3.16.7

## 修复

1. 修复cr-ee的namespace和repository在popgw模式下调度失败的问题
2. 修复cr-ee的instance, namespace, repo查询失败的问题

---

# 3.16.6

## 修复

1. edas_k8s_service: 更换read接口，补充部分只读属性。完善相关表文档
2. cr_repo将name长度调整到64与页面能力保持一致
3. 修复部分文档页面结果错误问题

---

# 3.16.5

## 新增

1. 组织资源集过滤修改精确匹配。
2. 修复vpngateway资源创建。
3. 汇丰资源平滑迁移测试
4. 新增新资源 Polardb_instance
5. 新增新资源 Polardb_database
6. 新增新资源 Polardb_account
7. 新增新资源 Polardb_backup_policy
8. 新增新资源 Polardb_dbconnection
9. 新增新资源 Polardb_Zone

---

# 3.16.4

## 新增

1. Edas_k8s_app创建时，支持设置PVC挂载，本地挂载，配置设置
2. 新增新资源 Edas_k8s_service
3. 新增新资源 Edas_namespace

---

# 3.16.3

## 修复

1. 修复alibabacloudstack_datahub_project在popgw模式下的问题，comment不在允许修改
2. 修复alibabacloudstack_datahub_topic在popgw模式下的问题
3. 修复alibabacloudstack_datahub_subscription在popgw模式下的问题，comment不在允许修改

---

# 3.16.2

## 修复

1. 修复DNS不支持HTTPS请求协议，所有AliDNS资源强制使用http
2. OSS_Bucket的创建兼容318x的数据格式
3. 修复在popgw模式下，OSS使用ECS endpoint的问题
4. db_readonly_instance逻辑漏洞修复
5. kvstore_instance_class问题series可选值大小写规范
6. 修复部分测试用例无法导出测试模板的问题

---

# 3.16.1

## 新增

1. 支持apsarastack用户平滑迁移
2. 支持sls的log_alert资源的增删改查
3. redis_instance(kvstore_instance)开启TDE的能力
4. ess_scalinggroup的multi_az_policy属性支持
5. ecs instance的data_disk_tags
6. ecs_disk的auto_snapshot_policy_id和enable_automated_snapshot_policy支持
7. oss_bucket增加bucket_sync（同步），storage_capacity（容量限制），sse_algorithm（加密方式)
8. slb_listener增加logs_download_attributes
9. image_copy的加密功能
10. cs_k8s指定security_group_id的功能

## 修复

1. 修复ecs在更新image时，system_disk的tags未能正常继承的问题

## 优化

1. 优化slb_listener acl相关的参数的依赖互斥检查逻辑
2. 对github.com/aliyun/aliyun-datahub-sdk-go的依赖

---

# 3.16.0 

## 新增

1. 支持ASAPI和POPAPI的调度， 支持单独制定云产品的endpoint
2. 支持ASAPI的权限处理机制，所有的请求会通过ASAPI的进行鉴权和消息错误封装
3. 支持鹰眼ID，优化报错信息和消息提示
4. 通过自动化生成补充和提升了测试覆盖率，修复了测试case和Provider逻辑，提升了测试通过率