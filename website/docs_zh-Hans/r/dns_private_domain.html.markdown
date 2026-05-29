---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_private_domain"
sidebar_current: "docs-Alibabacloudstack-resource-dns-private-domain"
description: |-
  云解析私有域名
---

# alibabacloudstack_dns_private_domain

云解析私有域名资源，用于创建和管理阿里云私有域名。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tfacc77356.test."
}

resource "alibabacloudstack_vpc_vpc" "default0" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}_vpc0"
}

resource "alibabacloudstack_vpc_vpc" "default1" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = "${var.name}_vpc1"
}


resource "alibabacloudstack_dns_private_domain" "default" {
  name = var.name
  vpc_ids = [
    alibabacloudstack_vpc_vpc.default0.id,
    alibabacloudstack_vpc_vpc.default1.id,
  ]
}
```

## 参数说明

支持以下参数：

* `name` - (必填, 变更时重建) 私有域名名称。必须以"."结尾，例如"example.com."。

* `remark` - (可选) 域名备注信息，用于描述该私有域名的用途。

* `vpc_ids` - (可选) 关联的VPC ID列表。将私有域名绑定到指定的VPC网络，使VPC内的ECS实例可以通过内网访问该域名。

## 属性说明

以下属性会从API响应中导出：

* `id` - 私有域名ID。

* `caller_uid` - 系统参数，表示创建该私有域名的用户ID。

* `create_timestamp` - 创建时间戳（秒）。

* `record_count` - 解析记录集总数。

* `update_timestamp` - 修改时间戳（秒）。

## Import

DNS Private Domain 可以使用 ID 导入，例如：

```
$ terraform import alibabacloudstack_dns_private_domain.example <id>
```