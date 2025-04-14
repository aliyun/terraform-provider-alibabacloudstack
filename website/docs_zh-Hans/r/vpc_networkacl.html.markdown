---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_networkacl"
sidebar_current: "docs-Alibabacloudstack-vpc-networkacl"
description: |- 
  编排VPC的网络访问控制列表(ACL）
---

# alibabacloudstack_vpc_networkacl
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_network_acl`

使用Provider配置的凭证在指定的资源集编排VPC的网络访问控制列表(ACL）。

## 示例用法

```terraform
variable "name" {
	default = "tf-testaccnetworkacl40200"
}
variable "name_change" {
	default = "tf-testaccnetworkacl40200_change"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  name       = var.name
}

resource "alibabacloudstack_vswitch" "default0" {
  vpc_id            = alibabacloudstack_vpc.default.id
  name              = var.name
  cidr_block        = cidrsubnet(alibabacloudstack_vpc.default.cidr_block, 4, 4)
  availability_zone = data.alibabacloudstack_zones.default.ids.0
}

resource "alibabacloudstack_vswitch" "default1" {
  vpc_id            = alibabacloudstack_vpc.default.id
  name              = var.name_change
  cidr_block        = cidrsubnet(alibabacloudstack_vpc.default.cidr_block, 4, 5)
  availability_zone = data.alibabacloudstack_zones.default.ids.0
}

resource "alibabacloudstack_network_acl" "default" {
  vpc_id           = alibabacloudstack_vpc.default.id
  network_acl_name = "tf-testaccnetworkacl40200"
  description      = "This is a test network ACL"

  ingress_acl_entries {
    description            = "Allow SSH and HTTP traffic"
    network_acl_entry_name = "tcp22-80-ingress"
    source_cidr_ip         = "0.0.0.0/0"
    policy                 = "accept"
    port                   = "22/80"
    protocol               = "tcp"
  }

  egress_acl_entries {
    description            = "Allow all outbound traffic"
    network_acl_entry_name = "all-egress"
    destination_cidr_ip    = "0.0.0.0/0"
    policy                 = "accept"
    port                   = "-1/-1"
    protocol               = "all"
  }

  resources {
    resource_id   = alibabacloudstack_vswitch.default0.id
    resource_type = "VSwitch"
  }

  resources {
    resource_id   = alibabacloudstack_vswitch.default1.id
    resource_type = "VSwitch"
  }
}
```

## 参数说明

支持以下参数：

* `vpc_id` - (必填, 变更时重建) - 关联的VPC的ID。此字段在创建后无法更改。
* `network_acl_name` - (选填) - 网络ACL的名称。名称长度为1～128个字符，不能以http://或https://开头。
* `description` - (选填) - 网络ACL的描述信息。描述长度为1～256个字符，不能以http://或https://开头。

### 入站规则

* `ingress_acl_entries` - (选填) - 网络ACL的入站规则列表。规则的顺序决定了优先级。每个条目支持以下内容：

  * `description` - (选填) - 入站规则的描述。
  * `network_acl_entry_name` - (选填) - 入站规则的名称。
  * `policy` - (选填) - 入站规则的策略。有效值：`accept`、`drop`。
  * `port` - (选填) - 入站规则的端口范围。
  * `protocol` - (选填) - 入站规则的协议。有效值：`icmp`、`gre`、`tcp`、`udp`、`all`。
  * `source_cidr_ip` - (选填) - 入站规则的源CIDR IP。

### 出站规则

* `egress_acl_entries` - (选填) - 网络ACL的出站规则列表。规则的顺序决定了优先级。每个条目支持以下内容：

  * `description` - (选填) - 出站规则的描述。
  * `network_acl_entry_name` - (选填) - 出站规则的名称。
  * `policy` - (选填) - 出站规则的策略。有效值：`accept`、`drop`。
  * `port` - (选填) - 出站规则的端口范围。
  * `protocol` - (选填) - 出站规则的协议。有效值：`icmp`、`gre`、`tcp`、`udp`、`all`。
  * `destination_cidr_ip` - (选填) - 出站规则的目标CIDR IP。

### 关联资源

* `resources` - (选填) - 与此网络ACL相关的资源列表。每个资源支持以下内容：

  * `resource_id` - (选填) - 相关资源的ID。
  * `resource_type` - (选填) - 相关资源的类型。有效值：`VSwitch`。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 网络ACL实例的ID。
* `status` - 网络ACL的状态。
* `egress_acl_entries` - 出方向规则信息。包含出站规则的具体配置和优先级。
* `ingress_acl_entries` - 入方向规则信息。包含入站规则的具体配置和优先级。
* `network_acl_name` - 网络ACL的名称。
* `name` - 已废弃字段，建议使用`network_acl_name`代替。
* `vpc_id` - (计算得出) - 关联的VPC的ID。