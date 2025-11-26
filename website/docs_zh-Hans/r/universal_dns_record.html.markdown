---
subcategory: "Universal DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_universal_dns_record"
sidebar_current: "docs-Alibabacloudstack-universal-dns-record"
description: |-
  跨云解析域名解析
---

# alibabacloudstack_universal_dns_record

跨云解析域名解析资源，用于管理跨云共享域名的解析配置记录。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testacc58082"
}

resource "alibabacloudstack_universal_dns_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS record test"
}



resource "alibabacloudstack_universal_dns_record" "default" {
  name         = var.name
  type         = "A"
  ttl          = 300
  lba_strategy = "ALL_RR"
  line_ids = [
    "default"
  ]
  rdatas {
    value = "192.168.1.1"
  }
  rdatas {
    value = "127.0.0.1"
  }

  zone_id = alibabacloudstack_universal_dns_domain.default.id
}
```

## 参数说明

支持以下参数：

* `line_ids` - (必填) 解析线路列表。例如：`["default"]`。
* `lba_strategy` - (必填) 负载均衡策略。取值：
  * `ALL_RR`：返回全部地址（非权重）。
  * `RATIO`：按权重返回地址（权重）。
* `name` - (必填) 记录名称。例如：`"test"`。
* `rdatas` - (必填) 记录值集合。每个记录值包含以下属性：
  * `value` - (必填) 值。例如：`"192.168.1.1"`。
  * `lba_weight` - (可选) 权重。仅当`lba_strategy`为`RATIO`时有效。例如：`100`。
* `ttl` - (必填) 记录的缓存时间（秒）。例如：`300`。
* `type` - (必填) 记录类型。取值：
  * `A`：将域名指向一个IPv4地址。
  * `AAAA`：将域名指向一个IPv6地址。
  * `CNAME`：将域名指向另一个域名。
  * `MX`：将域名指向邮件服务器地址。
  * `TXT`：文本长度限制255，通常做SPF记录（反垃圾邮件）。
  * `PTR`：记录IP地址反向解析的域名。
  * `SRV`：记录提供特定的服务的服务器。
  * `NAPTR`：记录域名权威指针，通常用于ENUM。
  * `CAA`：为域名设置证书分发机构授权信息。
  * `NS`：将域名指向指定的DNS服务器解析。
* `zone_id` - (必填, 变更时重建) 域名ID。例如：`"f27ffcc8-02a4-4ca1-9ca4-6bffc4e82034"`。
* `remark` - (可选) 备注信息。

## 属性说明

以下属性会从API响应中导出：

* `id` - 资源ID，格式为{ZoneId}:{Id}。
* `create_timestamp` - 创建时间戳（秒）。
* `update_timestamp` - 修改时间戳（秒）。