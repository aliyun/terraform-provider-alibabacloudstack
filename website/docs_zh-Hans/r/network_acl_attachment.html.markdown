---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_network_acl_attachment"
sidebar_current: "docs-alibabacloudstack-resource-network-acl-attachment"
description: |-
  编排绑定云网络访问控制列表(Network ACL)与交换机(vswitch)
---

# alibabacloudstack_network_acl_attachment

使用Provider配置的凭证在指定的资源集编排绑定云网络访问控制列表(Network ACL)与交换机(vswitch)。

-> **已弃用:** 此资源已被弃用。请使用 [alibabacloudstack_network_acl](https://www.terraform.io/docs/providers/alibabacloudstack/r/network_acl.html) 资源中的 `resources` 属性替代。
请注意，由于此资源与 `alibabacloudstack_network_acl` 的 `resources` 属性存在冲突，因此无法同时使用。

-> **注意:** 目前，该资源仅在以下区域可用：香港(cn-hongkong)、印度(ap-south-1)和印尼(ap-southeast-1)。

## 示例用法

### 基础用法

```
variable "name" {
  default = "NatGatewayConfigSpec"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_network_acl" "default" {
  vpc_id           = alibabacloudstack_vpc.default.id
  network_acl_name = var.name
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alibabacloudstack_network_acl_attachment" "default" {
  network_acl_id = alibabacloudstack_network_acl.default.id
  resources {
    resource_id   = alibabacloudstack_vswitch.default.id
    resource_type = "VSwitch"
  }
}
```

## 参数说明

支持以下参数：

* `network_acl_id` - (必填，变更时重建) 网络ACL的ID，该字段不可更改。
* `resources` - (必填) 与网络ACL关联的资源列表。详细信息见以下描述：
  * `resource_id` - (必填) 要与网络ACL关联的资源ID。
  * `resource_type` - (必填) 要与网络ACL关联的资源类型。目前仅支持 `VSwitch`。

## 属性说明

导出以下属性：

* `id` - 网络ACL关联的唯一标识符。格式为 `<network_acl_id>:<一个唯一ID>`。
* `resources` - 与网络ACL关联的资源列表。该列表包含所有成功绑定到网络ACL的资源信息。