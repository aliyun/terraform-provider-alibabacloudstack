---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_record"
sidebar_current: "docs-alibabacloudstack-resource-dns-record"
description: |-
  编排DNS域名记录
---

# alibabacloudstack_dns_record

使用Provider配置的凭证在指定的资源集下编排DNS域名记录。

## 示例用法

```
resource "alibabacloudstack_dns_domain" "default" {
  domain_name = "domaintest."
  remark = "testing Domain"
}

# 创建一个新的域名记录
resource "alibabacloudstack_dns_record" "default" {
  zone_id   = alibabacloudstack_dns_domain.default.id
  name = "testing_record"
  type        = "A"
  remark = "testing Record"
  ttl         = 300
  lba_strategy = "ALL_RR"
  rr_set      = ["192.168.2.4","192.168.2.7","10.0.0.4"]
}

output "record" {
  value = alibabacloudstack_dns_record.default.*
}
```

## 参数说明

以下是支持的参数：

* `zone_id` - (必填) 此记录所属的 DNS 域的 ID。
* `name` - (必填) 域名记录的名称。该主机记录最多可以有 253 个字符，每个用“.”分隔的部分最多可以有 63 个字符，并且必须仅包含字母数字字符或连字符，例如“-”、“.”、“*”、“@”，并且不能以“-”开头或结尾。
* `type` - (必填) 域名记录的类型。有效值为 `A`, `NS`, `MX`, `TXT`, `CNAME`, `SRV`, `AAAA`, `CAA`, `REDIRECT_URL` 和 `FORWORD_URL`。
* `lba_strategy` - (必填) 域名记录的负载均衡策略。有效值为 `ALL_RR`（轮询）和 `RATIO`（加权）。
* `rr_set` - (可选) 域名记录的值列表。当 `type` 为 `MX`, `NS`, `CNAME`, `SRV` 时，服务器会将 `value` 视为完全限定域名，因此不需要在末尾添加“.”。
* `ttl` - (可选) 域名记录的有效时间（TTL）。其范围取决于云解析的版本。免费版是 `[600, 86400]`，基础版是 `[120, 86400]`，标准版是 `[60, 86400]`，高级版是 `[10, 86400]`，独享版是 `[1, 86400]`。默认值为 `300`。
* `remark` - (可选) 域名记录的备注信息，用于描述该记录的具体用途。

## 属性说明

以下属性会被导出：

* `record_id` - DNS 记录的唯一标识符（ID）。
* `type` - 当前记录的类型，例如 `A`, `CNAME`, `MX` 等。
* `name` - 当前记录的名称（即主机记录）。
* `rr_set` - 当前记录的值列表，具体值取决于记录类型。
* `ttl` - 当前记录的有效时间（TTL），单位为秒。
* `zone_id` - 此记录所属的 DNS 域的 ID。
* `lba_strategy` - 当前记录的负载均衡策略，例如 `ALL_RR` 或 `RATIO`。
* `line_ids` - 与当前 DNS 记录关联的线路 ID 列表。
* `remark` - 当前记录的备注信息，用于描述记录的具体用途。