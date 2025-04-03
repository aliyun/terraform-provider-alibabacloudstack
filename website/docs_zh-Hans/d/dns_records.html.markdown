---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_records"
sidebar_current: "docs-alibabacloudstack-datasource-dns-records"
description: |-
    查询DNS域名解析记录
---

# alibabacloudstack_dns_records

根据指定过滤条件列出当前凭证权限可以访问的DNS域名解析记录列表。

## 示例用法

```
resource "alibabacloudstack_dns_domain" "default" {
  domain_name = "domaintest."
  remark = "testing Domain"
}

# 创建一个新的域名解析记录
resource "alibabacloudstack_dns_record" "default" {
  zone_id   = alibabacloudstack_dns_domain.default.domain_id
  name = "testing_record"
  type        = "A"
  remark = "testing Record"
  ttl         = 300
  lba_strategy = "ALL_RR"
  rr_set      = ["192.168.2.4","192.168.2.7","10.0.0.4"]
}

data "alibabacloudstack_dns_records" "default"{
 zone_id         = alibabacloudstack_dns_record.default.zone_id
 name = alibabacloudstack_dns_record.default.name
}
output "records" {
  value = data.alibabacloudstack_dns_records.default.*
}
```

## 参数参考

支持以下参数：

* `zone_id` - (必填) 与解析记录关联的域名 ID。
* `host_record_regex` - (可选，强制更新)主机记录正则表达式。
* `type` - (可选) 记录类型。有效选项包括 `A`, `NS`, `MX`, `TXT`, `CNAME`, `SRV`, `AAAA`, `REDIRECT_URL`, `FORWORD_URL`。
* `ids` - (可选) 记录 ID 列表。
* `name` - (可选) DNS 记录名称。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 记录 ID 列表。
* `records` - 记录列表。每个元素包含以下属性：
  * `record_id` - 记录的 ID。
  * `zone_id` - 记录所属域名的 ID。
  * `name` - 域名的主机记录。
  * `type` - 记录类型。
  * `ttl` - 记录的 TTL。
  * `remark` - 记录的描述。
  * `rr_set` - 记录的 RrSet。