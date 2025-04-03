---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_network_acl_entries"
sidebar_current: "docs-alibabacloudstack-resource-network-acl-entries"
description: |-
  编排网络ACL条目资源
---

# alibabacloudstack_network_acl_entries

使用Provider配置的凭证在指定的资源集编排网络ACL条目资源以创建入站和出站条目。

-> **注意：** 目前，该资源仅在香港(cn-hongkong)、印度(ap-south-1)和印度尼西亚(ap-southeast-1)地区可用。

-> **注意：** 它不支持并发操作，且入站和出站条目的顺序决定了优先级。

-> **注意：** 使用此资源需要开启白名单。

-> **已弃用：** 此资源已被弃用。请使用 `ingress_acl_entries` 和 `egress_acl_entries` 替代，并与 [alibabacloudstack_network_acl](https://www.terraform.io/docs/providers/alibabacloudstack/r/network_acl.html) 资源一起使用。

## 示例用法

### 基础用法

```
variable "name" {
  default = "NetworkAclEntries"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_network_acl" "default" {
  vpc_id = alibabacloudstack_vpc.default.id
  name   = var.name
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
  name              = var.name
}

resource "alibabacloudstack_network_acl_attachment" "default" {
  network_acl_id = alibabacloudstack_network_acl.default.id
  resources {
    resource_id   = alibabacloudstack_vswitch.default.id
    resource_type = "VSwitch"
  }
}

resource "alibabacloudstack_network_acl_entries" "default" {
  network_acl_id = alibabacloudstack_network_acl.default.id
  ingress {
    protocol       = "all"
    port           = "-1/-1"
    source_cidr_ip = "0.0.0.0/32"
    name           = var.name
    entry_type     = "custom"
    policy         = "accept"
    description    = var.name
  }
  egress {
    protocol            = "all"
    port                = "-1/-1"
    destination_cidr_ip = "0.0.0.0/32"
    name                = var.name
    entry_type          = "custom"
    policy              = "accept"
    description         = var.name
  }
}
```

## 参数参考

支持以下参数：

* `network_acl_id` - (必填，变更时重建) 网络 ACL 的 ID，此字段不能更改。
* `ingress` - (可选) 网络 ACL 的入站条目列表。入站条目的顺序决定了优先级。详细信息请参阅 Block Ingress。资源映射支持以下内容：
  * `description` - (可选) 入站条目的描述。
  * `source_cidr_ip` - (可选) 入站条目的源 IP。
  * `entry_type` - (可选) 入站条目的条目类型。必须为 `custom` 或 `system`。默认值为 `custom`。
  * `name` - (可选) 入站条目的名称。
  * `policy` - (可选) 入站条目的策略。必须为 `accept` 或 `drop`。
  * `port` - (可选) 入站条目的端口。
  * `protocol` - (可选) 入站条目的协议。
* `egress` - (可选) 网络 ACL 的出站条目列表。出站条目的顺序决定了优先级。详细信息请参阅 Block Egress。资源映射支持以下内容：
  * `description` - (可选) 出站条目的描述。
  * `destination_cidr_ip` - (可选) 出站条目的目标 IP。
  * `entry_type` - (可选) 出站条目的条目类型。必须为 `custom` 或 `system`。默认值为 `custom`。
  * `name` - (可选) 出站条目的名称。
  * `policy` - (可选) 出站条目的策略。必须为 `accept` 或 `drop`。
  * `port` - (可选) 出站条目的端口。
  * `protocol` - (可选) 出站条目的协议。
* `network_acl_id` - (必填，变更时重建) 网络 ACL 的 ID，此字段不能更改。

## 属性参考

导出以下属性：

* `id` - 网络 ACL 条目的 ID。格式为 `<network_acl_id>:<a unique id>`。
* `description` - 入站或出站条目的描述。
* `source_cidr_ip` - 入站条目的源 IP。
* `entry_type` - 入站或出站条目的条目类型。
* `name` - 入站或出站条目的名称。
* `policy` - 入站或出站条目的策略。
* `port` - 入站或出站条目的端口。
* `protocol` - 入站或出站条目的协议。
* `destination_cidr_ip` - 出站条目的目标 IP。