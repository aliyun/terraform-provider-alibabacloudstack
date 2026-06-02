---
subcategory: "云服务总线 CSB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_csb_service"
sidebar_current: "docs-Alibabacloudstack-resource-csb-service"
description: |-
  提供管理 CSB 服务的 Alibabacloudstack 资源。
---

# alibabacloudstack_csb_service

该资源可帮助您管理 CSB 服务。

有关 CSB 服务及其使用方法的更多信息，请参阅 [创建服务](https://help.aliyun.com/apsara/enterprise/v_3_18_0_30393230/csb/apsarastack-developer-guide/createservice.html)

## 示例

基本用法

```hcl
resource "alibabacloudstack_csb_service" "service" {
  csb_id          = "your-csb-id"
  project_id      = "your-project-id"
  service_name    = "example-service"
  service_version = "1.0.0"
  provide_type    = "RESTful"
  skip_auth       = false
  all_visiable    = true
  route_conf_json = jsonencode({
    importConf = {
      accessEndpointJSON = jsonencode({
        # endpoint configuration
      })
    }
  })
}
```

## 参数说明

支持以下参数：

* `csb_id` - （必填，变更时强制重建，Int）CSB 实例的唯一 ID。修改此参数将强制创建新资源。
* `project_id` - （必填，String）服务所属项目的 ID。
* `service_name` - （必填，变更时强制重建，String）服务名称。修改此参数将强制创建新资源。
* `service_version` - （必填，变更时强制重建，String）服务版本（例如 '1.0.0'）。修改此参数将强制创建新资源。
* `route_conf_json` - （必填，String）JSON 格式的路由配置。该字段包含导入配置和访问端点设置。
* `alias` - （可选，String）服务的别名。
* `model_version` - （可选，String）服务的模型版本。默认值：`2.0`。
* `skip_auth` - （可选，Bool）是否跳过此服务的身份验证。默认值：`false`。
* `all_visiable` - （可选，Bool）服务是否对所有用户可见。默认值：`true`。
* `scope` - （可选，String）服务的范围（例如，'0' 表示私有，'1' 表示公开）。默认值：`0`。
* `consume_types` - （可选，String 列表）消费类型列表。有效值：`Restful`、`WebService`。
* `provide_type` - （可选，String）服务的提供类型。有效值：`RESTful`、`SpringCloud`、`HSF`、`WebService`、`DUBBO`、`JDBC`。默认值：`Restful`。
* `cas_serv_targets` - （可选，String 列表）CAS 服务目标列表。
* `qps` - （可选，Int）每秒查询数（只读监控指标）。默认值：`0`。
* `interface_name` - （可选，String）接口名称（对于 RESTful 服务通常为空）。
* `ip_white_str` - （可选，String）IP 白名单字符串。
* `ip_black_str` - （可选，String）IP 黑名单字符串。
* `err_def_json` - （可选，String）JSON 格式的错误定义。
* `access_params_json` - （可选，String）JSON 格式的访问参数。

## 属性说明

导出以下属性：

* `id` - 资源的唯一标识符，格式为 `csb_id:service_id`。
* `service_id` - CSB 服务的内部 ID。
* `csb_id` - CSB 实例的唯一 ID。
* `project_id` - 服务所属项目的 ID。
* `service_name` - 服务名称。
* `service_version` - 服务版本。
* `alias` - 服务别名。
* `model_version` - 服务的模型版本。
* `skip_auth` - 是否跳过身份验证。
* `all_visiable` - 服务是否对所有用户可见。
* `scope` - 服务的范围。
* `provide_type` - 服务的提供类型。
* `consume_types` - 消费类型列表。
* `cas_serv_targets` - CAS 服务目标列表。
* `qps` - 每秒查询数。
* `status` - 服务状态（例如，0=未激活，1=已激活）。
* `modified_time` - 最后修改时间戳（毫秒）。
* `principal_name` - 与服务关联的主体名称。
* `owner_id` - 服务的所有者 ID。
* `user_id` - 服务创建者的用户 ID。
* `interface_name` - 接口名称。
* `ssl` - 是否启用了 SSL。
* `ott_flag` - OTF 标志状态。
* `policy_handler` - 策略处理器（例如，'accept'）。
* `ip_white_str` - IP 白名单字符串。
* `ip_black_str` - IP 黑名单字符串。
* `err_def_json` - 错误定义 JSON。
* `access_params_json` - 访问参数 JSON。
* `route_conf_json` - 路由配置 JSON。

## 导入

CSB 服务可以通过 `csb_id` 和 `service_id` 的组合导入，以冒号分隔，例如：

```
$ terraform import alibabacloudstack_csb_service.example <csb_id>:<service_id>
```
