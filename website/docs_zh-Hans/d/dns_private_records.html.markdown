---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_private_records"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-private-records"
description: |-
  查询阿里云云解析私有域名解析记录
---

# alibabacloudstack_dns_private_records

> 阿里云云解析私有域名解析记录查询数据源

## 示例用法

```hcl
  
variable "name" {
  default = "tf-testacc56057"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS record test"
}

resource "alibabacloudstack_dns_private_record" "default" {
  zone_id      = alibabacloudstack_dns_private_domain.default.id
  name         = var.name
  type         = "A"
  ttl          = 300
  lba_strategy = "ALL_RR"
  line_ids     = ["default"]
  rdatas {
    value = "192.168.1.1"
  }
  rdatas {
    value = "127.0.0.1"
  }
}

data "alibabacloudstack_dns_private_records" "default" {
  zone_id    = alibabacloudstack_dns_private_record.default.zone_id
  name_regex = "tf-testacc[0-9]+"
}
```

## 参数说明
以下参数支持过滤查询结果：

* `zone_id` (必填)：域名ID，用于指定要查询的域名区域。

* `name_regex` (可选)：解析记录名称的正则表达式，用于按名称过滤结果。

* `ids` (可选)：要查询的解析记录ID列表，格式为"zone_id:id"。

## 属性说明
以下属性被导出：

* `id`：数据源的唯一标识符，格式为"zone_id1:id1,zone_id2:id2,..."

* `records`：查询到的解析记录列表。每个记录包含以下属性：
  * `id`：解析记录ID。
  * `create_timestamp`：创建时间戳（秒）。
  * `line_ids`：解析线路列表。
  * `lba_strategy`：负载均衡策略。取值说明：
    - ALL_RR：返回全部地址（非权重）。
    - RATIO：按权重返回地址（权重）。
  * `name`：记录名称。
  * `rdatas`：记录值列表。每个值包含以下属性：
    * `lba_weight`：权重。
    * `value`：值。
  * `remark`：备注。
  * `ttl`：记录的缓存时间。
  * `type`：记录类型。取值说明：
    - A：将域名指向一个IPV4地址。
    - AAAA：将域名指向一个IPV6地址。
    - CNAME：将域名指向另一个域名。
    - MX：将域名指向邮件服务器地址。
    - TXT：文本长度限制255，通常做SPF记录（反垃圾邮件）。
    - PTR：记录IP地址反向解析的域名。
    - SRV：记录提供特定的服务的服务器。
    - NAPTR：记录域名权威指针，通常用于ENUM。
    - CAA：为域名设置证书分发机构授权信息。
    - NS：将域名指向指定的DNS服务器解析。
  * `update_timestamp`：修改时间戳（秒）。
  * `zone_id`：域名ID。