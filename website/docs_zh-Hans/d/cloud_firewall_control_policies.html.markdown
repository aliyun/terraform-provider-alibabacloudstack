---
subcategory: "CloudFirewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloud_firewall_control_policies"
sidebar_current: "docs-Alibabacloudstack-datasource-cloud-firewall-control-policies"
description: |- 
  查询防火墙控制策略
---

# alibabacloudstack_cloud_firewall_control_policies

根据指定过滤条件列出当前凭证权限可以访问的防火墙控制策略列表。

## 示例用法

以下是一个完整的示例，展示如何使用`alibabacloudstack_cloud_firewall_control_policies`数据源来获取符合特定条件的防火墙控制策略列表：

```terraform
data "alibabacloudstack_cloud_firewall_control_policies" "example" {
  direction   = "in"
  acl_action  = "accept"
  source      = "192.168.0.0/16"
  destination = "10.0.0.0/8"
  proto       = "TCP"
  ip_version  = "IPv4"
  lang        = "zh"
  description = "示例访问控制策略"

  # 可选：保存结果到文件
  output_file = "control_policies_output.txt"
}
```

## 参数参考

以下参数是支持的：

* `acl_action` - (选填, 变更时重建) - 安全访问控制策略中设置的流量通过云防火墙的方式。有效值：
  * **accept**：放行
  * **drop**：拒绝
  * **log**：观察

* `acl_uuid` - (选填, 变更时重建) - 安全访问控制策略的唯一标识ID。

* `description` - (选填, 变更时重建) - 安全访问控制策略的描述信息。

* `destination` - (选填, 变更时重建) - 安全访问控制策略中的目的地址。具体值取决于`destination_type`：
  * 当`destination_type`为`net`时，`destination`为目的CIDR(例如：`1.2.3.4/24`)。
  * 当`destination_type`为`group`时，`destination`为目的地址簿名称(例如：`db_group`)。
  * 当`destination_type`为`domain`时，`destination`为目的域名(例如：`*.aliyuncs.com`)。
  * 当`destination_type`为`location`时，`destination`为目的区域(例如：`["BJ11", "ZB"]`)。

* `direction` - (必填, 变更时重建) - 安全访问控制策略的流量方向。有效值：
  * **in**：外对内流量访问控制
  * **out**：内对外流量访问控制

* `ip_version` - (选填, 变更时重建) - 访问控制策略中的地址IP版本。有效值：
  * **IPv4**
  * **IPv6**

* `lang` - (选填, 变更时重建) - 请求和接收消息的语言类型。有效值：
  * **en**：英文
  * **zh**：中文

* `proto` - (选填, 变更时重建) - 安全访问控制策略中流量访问的安全协议类型。有效值：
  * **ANY**
  * **TCP**
  * **UDP**
  * **ICMP**

* `source` - (选填, 变更时重建) - 安全访问控制策略中的源地址。具体值取决于`source_type`：
  * 当`source_type`为`net`时，`source`为源CIDR(例如：`1.2.3.0/24`)。
  * 当`source_type`为`group`时，`source`为源地址簿名称(例如：`db_group`)。
  * 当`source_type`为`location`时，`source`为源区域(例如：`["BJ11", "ZB"]`)。


## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 控制策略ID列表。

* `policies` - 云防火墙控制策略列表。每个元素包含以下属性：
  * `id` - 控制策略的ID。格式为`<acl_uuid>:<direction>`。
  * `acl_uuid` - 安全访问控制策略的唯一标识ID。
  * `acl_action` - 安全访问控制策略中设置的流量通过云防火墙的方式。有效值：
    * **accept**：放行
    * **drop**：拒绝
    * **log**：观察
  * `application_id` - 安全访问控制策略中设置访问的应用ID。
  * `application_name` - 安全访问控制策略支持的应用类型有以下几种：
    * **ANY**
    * **HTTP**
    * **HTTPS**
    * **MySQL**
    * **SMTP**
    * **SMTPS**
    * **RDP**
    * **VNC**
    * **SSH**
    * **Redis**
    * **MQTT**
    * **MongoDB**
    * **Memcache**
    * **SSL**
  * `description` - 安全访问控制策略的描述信息。
  * `dest_port` - 安全访问控制策略中流量访问的目的端口。
  * `dest_port_group` - 安全访问控制策略中流量访问的目的端口地址簿名称。
  * `dest_port_group_ports` - 目的端口地址簿中包含的端口列表。
  * `dest_port_type` - 安全访问控制策略中流量访问的目的端口类型。有效值：
    * **port**：端口
    * **group**：端口地址簿
  * `destination` - 安全访问控制策略中的目的地址。具体值取决于`destination_type`：
    * 当`destination_type`为`net`时，`destination`为目的CIDR(例如：`1.2.3.4/24`)。
    * 当`destination_type`为`group`时，`destination`为目的地址簿名称(例如：`db_group`)。
    * 当`destination_type`为`domain`时，`destination`为目的域名(例如：`*.aliyuncs.com`)。
    * 当`destination_type`为`location`时，`destination`为目的区域(例如：`["BJ11", "ZB"]`)。
  * `destination_group_cidrs` - 安全访问控制策略中的目的地址簿中的网段列表。
  * `destination_group_type` - 安全访问控制策略中的目的地址簿类型。有效值：
    * **ip**：IP地址簿，包含一个或多个IP地址段。
    * **tag**：ECS标签地址簿，包含一个或多个ECS标签下的IP地址。
    * **domain**：域名地址簿，包含一个或多个域名地址。
    * **threat**：威胁地址簿，包含一个或多个恶意IP或域名地址。
    * **backsrc**：回源地址簿，包含一个或多个DDoS防护实例或WAF实例的回源地址。
  * `destination_type` - 安全访问控制策略中的目的地址类型。有效值：
    * **net**：目的网段(CIDR)
    * **group**：目的地址簿
    * **domain**：目的域名
    * **location**：目的区域
  * `direction` - 安全访问控制策略的流量方向。有效值：
    * **in**：外对内流量访问控制
    * **out**：内对外流量访问控制
  * `dns_result` - DNS解析结果。
  * `dns_result_time` - DNS解析时间。
  * `hit_times` - 安全访问控制策略命中次数统计。
  * `order` - 安全访问控制策略生效的优先级。优先级数字从1开始顺序递增，优先级数字越小，优先级越高。特殊值：
    * **-1**：表示优先级最低。
  * `proto` - 安全访问控制策略中流量访问的安全协议类型。有效值：
    * **ANY**
    * **TCP**
    * **UDP**
    * **ICMP**
  * `release` - 安全访问控制策略是否生效。
  * `source` - 安全访问控制策略中的源地址。具体值取决于`source_type`：
    * 当`source_type`为`net`时，`source`为源CIDR(例如：`1.2.3.0/24`)。
    * 当`source_type`为`group`时，`source`为源地址簿名称(例如：`db_group`)。
    * 当`source_type`为`location`时，`source`为源区域(例如：`["BJ11", "ZB"]`)。
  * `source_group_cidrs` - 安全访问控制策略中的源地址簿中的网段列表。
  * `source_group_type` - 安全访问控制策略中的源地址簿类型。有效值：
    * **ip**：IP地址簿，包含一个或多个IP地址段。
    * **tag**：ECS标签地址簿，包含一个或多个ECS标签下的IP地址。
    * **domain**：域名地址簿，包含一个或多个域名地址。
    * **threat**：威胁地址簿，包含一个或多个恶意IP或域名地址。
    * **backsrc**：回源地址簿，包含一个或多个DDoS防护实例或WAF实例的回源地址。
  * `source_type` - 安全访问控制策略中的源地址类型。有效值：
    * **net**：源网段(CIDR)
    * **group**：源地址簿
    * **location**：源区域
