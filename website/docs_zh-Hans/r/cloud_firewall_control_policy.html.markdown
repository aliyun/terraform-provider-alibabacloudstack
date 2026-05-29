---
subcategory: "云防火墙"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloud_firewall_control_policy"
sidebar_current: "docs-Alibabacloudstack-resource-cloud-firewall-control-policy"
description: |-
  管理云防火墙访问控制策略。
---

# alibabacloudstack_cloud_firewall_control_policy

管理云防火墙访问控制策略，用于定义通过云防火墙流量的访问控制规则。

-> **注意：** 该资源也可以使用以下别名引用：
-> - `alibabacloudstack_cloudfw_controlpolicy`

## 示例代码

### 基本用法

```hcl
resource "alibabacloudstack_cloud_firewall_control_policy" "default" {
  acl_action       = "accept"
  application_name = "HTTP"
  description      = "tf-testacc-control-policy"
  destination      = "192.168.0.0/24"
  destination_type = "net"
  direction        = "out"
  proto            = "TCP"
  source           = "10.0.0.0/8"
  source_type      = "net"
}
```

## 参数说明

支持以下参数：

* `acl_action` - （必填）访问控制策略对匹配流量执行的动作。取值：`accept`（放行）、`drop`（拒绝）、`log`（观察）。
* `application_name` - （必填）访问控制策略支持的应用类型。取值：`ANY`、`HTTP`、`HTTPS`、`MQTT`、`Memcache`、`MongoDB`、`MySQL`、`RDP`、`Redis`、`SMTP`、`SMTPS`、`SSH`、`SSL`、`VNC`。
* `description` - （必填）访问控制策略的描述信息。
* `destination` - （必填）访问控制策略中的目的地址。取值取决于 `destination_type`：
  - 当 `destination_type` 为 `net` 时，填写 CIDR 地址段（如 `192.168.0.0/24`）。
  - 当 `destination_type` 为 `group` 时，填写目的地址簿名称（如 `db_group`）。
  - 当 `destination_type` 为 `domain` 时，填写目的域名（如 `*.example.com`）。
  - 当 `destination_type` 为 `location` 时，填写区域代码（如 `BJ11`、`ZB`）。
* `destination_type` - （必填）目的地址类型。取值：`group`（地址簿）、`location`（区域）、`net`（地址段）、`domain`（域名）。
* `direction` - （必填）访问控制策略的流量方向。取值：`in`（内对外访问）、`out`（外对内访问）。修改此参数会强制重新创建资源。
* `proto` - （必填）访问控制策略中流量访问的协议类型。取值：`ANY`、`TCP`、`UDP`、`ICMP`。
* `source` - （必填）访问控制策略中的源地址。取值取决于 `source_type`：
  - 当 `source_type` 为 `net` 时，填写 CIDR 地址段（如 `192.168.0.0/24`）。
  - 当 `source_type` 为 `group` 时，填写源地址簿名称（如 `db_group`）。
  - 当 `source_type` 为 `location` 时，填写区域代码（如 `BJ11`、`ZB`）。
* `source_type` - （必填）源地址类型。取值：`group`（地址簿）、`location`（区域）、`net`（地址段）。
* `dest_port` - （可选）访问控制策略中流量访问的目的端口。当 `dest_port_type` 为 `port` 时需要设置。
* `dest_port_group` - （可选）访问控制策略中流量访问的目的端口地址簿名称。当 `dest_port_type` 为 `group` 时需要设置。
* `dest_port_type` - （可选）目的端口类型。取值：`group`（端口组）、`port`（端口）。
* `lang` - （可选）请求语言。取值：`en`（英文）、`zh`（中文）。
* `release` - （可选）访问控制策略的状态。策略创建后默认启用。取值：`true`（启用）、`false`（停用）。
* `source_ip` - （可选）请求的源 IP 地址。

-> **注意：** `dest_port` 和 `dest_port_group` 参数根据 `dest_port_type` 的值互斥。当 `dest_port_type` 为 `port` 时，设置 `dest_port`；当 `dest_port_type` 为 `group` 时，设置 `dest_port_group`。

## 属性说明

除上述参数外，该资源还导出以下属性：

* `id` - 资源 ID。格式为 `<acl_uuid>:<direction>`。
* `acl_uuid` - 访问控制策略的唯一标识符。

## 导入

云防火墙访问控制策略可以使用 `acl_uuid` 和 `direction` 导入，例如：

```
$ terraform import alibabacloudstack_cloud_firewall_control_policy.example acl-12345678:out
```
