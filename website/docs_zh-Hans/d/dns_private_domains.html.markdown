---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_private_domains"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-private-domains"
description: |-
  查询阿里云DNS私有域名列表
---

# alibabacloudstack_dns_private_domains

查询阿里云DNS私有域名列表。该数据源用于检索已创建的私有域名信息，包括域名ID、名称、关联VPC等详细信息。

## 示例用法

```hcl

variable "name" {
  default = "tf-testacc4077889550832453245."
}

resource "alibabacloudstack_vpc_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}_vpc"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name    = var.name
  remark  = var.name
  vpc_ids = [alibabacloudstack_vpc_vpc.default.id]
}

data "alibabacloudstack_dns_private_domains" "default" {
  name_regex = alibabacloudstack_dns_private_domain.default.name
}

```

## 参数说明
以下参数用于过滤查询结果：

- **id** (可选)：域名ID，用于精确匹配单个域名。
- **ids** (可选)：域名ID列表，用于过滤多个特定ID的域名。
- **name** (可选)：域名名称，用于精确匹配域名。
- **name_regex** (可选)：域名名称正则表达式，用于模糊匹配域名名称。
- **vpc_id** (可选)：VPC ID，用于查询关联了特定VPC的域名。

## 属性说明
以下属性被导出：

- **id** (字符串)：数据源ID，由查询结果的哈希值生成。
- **domains** (列表)：查询到的私有域名列表，每个元素包含以下属性：
  - **id** (字符串)：域名ID。
  - **name** (字符串)：域名名称。
  - **caller_uid** (字符串)：调用者UID。
  - **create_timestamp** (整数)：域名创建时间戳（秒）。
  - **gtm_instance_count** (整数)：关联的全局流量管理实例总数。
  - **record_count** (整数)：域名下解析记录集总数。
  - **region_and_vpcs** (列表)：区域和VPC关联信息列表，每个元素包含：
    - **region_id** (字符串)：区域ID。
    - **vpcs** (列表)：关联的VPC列表，每个元素包含：
      - **id** (字符串)：VPC ID。
      - **name** (字符串)：VPC名称。
  - **remark** (字符串)：域名备注信息。
  - **update_timestamp** (整数)：域名最后修改时间戳（秒）。
- **ids** (列表)：查询到的所有域名ID列表。