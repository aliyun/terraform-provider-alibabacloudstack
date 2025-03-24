---
subcategory: "CloudFirewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfirewall_controlpolicy"
sidebar_current: "docs-Alibabacloudstack-cloudfirewall-controlpolicy"
description: |- 
  编排云防火墙控制策略
---

# alibabacloudstack_cloudfirewall_controlpolicy
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cloud_firewall_control_policy`

使用Provider配置的凭证在指定的资源集下编排云防火墙控制策略。

## 示例用法

### 基础用法

以下是一个基本的示例，展示如何创建一个简单的云防火墙控制策略：

```terraform
variable "name" {
    default = "tf-testacccloud_firewallcontrol_policy46819"
}

resource "alibabacloudstack_cloudfirewall_controlpolicy" "default" {
  source           = "0.0.0.0/0"
  proto            = "ANY"
  destination      = "0.0.0.0/0"
  application_name = "ANY"
  acl_action       = "accept"
  dest_port_type   = "port"
  release          = "true"
  description      = "test"
  direction        = "in"
  source_type      = "net"
  dest_port        = "80"
  destination_type = "net"
}
```

### 高级用法

以下是一个更复杂的示例，包含 `acl_uuid` 和 `ip_version` 参数：

```terraform
resource "alibabacloudstack_cloudfirewall_controlpolicy" "example" {
  acl_uuid         = "example-acl-uuid"
  ip_version       = "ipv4"
  source           = "192.168.1.0/24"
  proto            = "TCP"
  destination      = "10.0.0.0/24"
  application_name = "HTTP"
  acl_action       = "drop"
  dest_port_type   = "port"
  release          = "false"
  description      = "Advanced example"
  direction        = "out"
  source_type      = "net"
  dest_port        = "8080"
  destination_type = "net"
}
```

## 参数参考

支持以下参数：

* `acl_action` - (必填) - 安全访问控制策略中设置的流量通过云防火墙的方式。
  - **accept**：放行
  - **drop**：拒绝
  - **log**：观察

* `acl_uuid` - (强制新建, 可选) - 安全访问控制策略的唯一标识ID。如果未指定，Terraform将自动生成一个。

* `application_name` - (必填) - 安全访问控制策略支持的应用类型有以下几种:
  - **ANY**
  - **HTTP**
  - **HTTPS**
  - **MySQL**
  - **SMTP**
  - **SMTPS**
  - **RDP**
  - **VNC**
  - **SSH**
  - **Redis**
  - **MQTT**
  - **MongoDB**
  - **Memcache**
  - **SSL**
  > **ANY** 表示策略应用在所有类型的应用中。

* `description` - (必填) - 安全访问控制策略的描述信息。

* `dest_port` - (选填) - 安全访问控制策略中流量访问的目的端口。如果 `dest_port_type` 设置为 `port`，则需要此参数。

* `dest_port_group` - (选填) - 安全访问控制策略中流量访问的目的端口地址簿名称。如果 `dest_port_type` 设置为 `group`，则需要此参数。

* `dest_port_type` - (选填) - 安全访问控制策略中流量访问的目的端口类型。
  - **port**：端口
  - **group**：端口地址簿

* `destination` - (必填) - 安全访问控制策略中的目的地址。
  - 当 `destination_type` 为 `net` 时，`destination` 为目的CIDR。例如：`1.2.3.4/24`
  - 当 `destination_type` 为 `group` 时，`destination` 为目的地址簿名称。例如：`db_group`
  - 当 `destination_type` 为 `domain` 时，`destination` 为目的域名。例如：`*.aliyuncs.com`
  - 当 `destination_type` 为 `location` 时，`destination` 为目的区域(具体区域位置编码见后文)。例如：`["BJ11", "ZB"]`

* `destination_type` - (必填) - 安全访问控制策略中的目的地址类型。
  - **net**：目的网段(CIDR)
  - **group**：目的地址簿
  - **domain**：目的域名
  - **location**：目的区域

* `direction` - (必填, 变更时重建) - 安全访问控制策略的流量方向。
  - **in**：外对内流量访问控制
  - **out**：内对外流量访问控制

* `ip_version` - (选填) - IP版本。有效值：`ipv4`, `ipv6`。

* `lang` - (选填) - 请求和接收消息的语言类型。有效值：`en`, `zh`。

* `proto` - (必填) - 安全访问控制策略中流量访问的安全协议类型。
  - **ANY**
  - **TCP**
  - **UDP**
  - **ICMP**

* `release` - (选填) - 指定安全访问控制策略是否生效。默认情况下，创建后访问控制策略是启用的。有效值：`true`, `false`。

* `source` - (必填) - 安全访问控制策略中的源地址。
  - 当 `source_type` 为 `net` 时，`source` 为源CIDR。例如：`1.2.3.0/24`
  - 当 `source_type` 为 `group` 时，`source` 为源地址簿名称。例如：`db_group`
  - 当 `source_type` 为 `location` 时，`source` 为源区域(具体区域位置编码见后文)。例如：`["BJ11", "ZB"]`

* `source_ip` - (选填) - 源IP地址。

* `source_type` - (必填) - 安全访问控制策略中的源地址类型。
  - **net**：源网段(CIDR)
  - **group**：源地址簿
  - **location**：源区域

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 云防火墙控制策略的唯一标识符。格式为 `<acl_uuid>:<direction>`。
* `acl_uuid` - 安全访问控制策略的唯一标识ID。
* `dest_port` - 安全访问控制策略中流量访问的目的端口。
* `dest_port_group` - 安全访问控制策略中流量访问的目的端口地址簿名称。
* `dest_port_type` - 安全访问控制策略中流量访问的目的端口类型。
  - **port**：端口
  - **group**：端口地址簿
* `release` - 指定安全访问控制策略是否生效。
* `source_ip` - 源IP地址。