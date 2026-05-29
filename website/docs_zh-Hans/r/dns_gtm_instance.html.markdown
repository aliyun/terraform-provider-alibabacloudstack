---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_instance"
sidebar_current: "docs-Alibabacloudstack-resource-dns-gtm-instance"
description: |-
  创建和管理云解析全局流量管理(GTM)实例
---

# alibabacloudstack_dns_gtm_instance

使用Provider配置的凭证在指定的资源集创建和管理云解析全局流量管理(GTM)实例。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tfacc24074"
}
resource "alibabacloudstack_vpc_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}_vpc0"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS record test"
  vpc_ids = [
    "${alibabacloudstack_vpc_vpc.default.id}"
  ]
}




resource "alibabacloudstack_dns_gtm_instance" "default" {
  name    = var.name
  prefix  = var.name
  zone_id = alibabacloudstack_dns_private_domain.default.id
  ttl     = "300"
}
```

## 参数说明

支持以下参数：

* `name` - (必填) 调度实例名称。名称长度为1到128个字符，不能以`http://`或`https://`开头。
* `prefix` - (必填) 调度域名。用于构建完整的调度域名，例如当prefix为"www"且zone_name为"example.com."时，完整调度域名为"www.example.com."。
* `ttl` - (必填) 全局TTL（缓存时间），单位为秒。表示DNS解析结果在客户端缓存的时间长度。
* `zone_id` - (必填) 归属调度域的域名ID。该ID对应于已创建的云解析域名。

## 属性说明

以下属性会从API响应中导出：

* `id` - 调度实例ID。
* `create_timestamp` - 创建时间戳（秒）。
* `update_timestamp` - 修改时间戳（秒）。
* `zone_name` - 域名名称。表示该调度实例所属的域名。