---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_accesscontrollists"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-accesscontrollists"
description: |- 
  查询负载均衡(SLB)访问控制列表
---

# alibabacloudstack_slb_accesscontrollists
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_acls`

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)访问控制列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccSlbAclDataSourceBisic-10390"
}
variable "ip_version" {
  default = "ipv4"
}

resource "alibabacloudstack_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
    entry = "10.10.10.0/24"
    comment = "first"
  }
  entry_list {
    entry = "168.10.10.0/24"
    comment = "second"
  }
}

data "alibabacloudstack_slb_acls" "default" {
  ids        = ["${alibabacloudstack_slb_acl.default.id}"]
  name_regex = "${alibabacloudstack_slb_acl.default.name}"
  tags       = {
    Environment = "Test"
  }
}

output "acl_id" {
  value = "${data.alibabacloudstack_slb_acls.default.acls.0.id}"
}

output "acl_name" {
  value = "${data.alibabacloudstack_slb_acls.default.acls.0.name}"
}

output "entry_list" {
  value = "${data.alibabacloudstack_slb_acls.default.acls.0.entry_list}"
}
```

## 参数说明

以下参数是支持的：

* `ids` - (可选) ACL ID列表，用于过滤结果。如果您知道要检索的特定ACL的ID，这将非常有用。
* `name_regex` - (可选，变更时重建) 一个正则表达式字符串，用于通过ACL名称过滤结果。这允许您基于模式匹配ACL名称。
* `tags` - (可选) 标签映射，用于通过标签过滤ACL。只有具有匹配标签的ACL才会被返回。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - SLB ACL ID列表。
* `names` - SLB ACL名称列表。
* `acls` - SLB ACL列表。每个元素包含以下属性：
  * `id` - ACL的唯一ID。
  * `name` - ACL的名称。
  * `ip_version` - 访问控制列表的IP版本，它决定了其条目的类型（IPv4或IPv6）。可能的值为`ipv4`或`ipv6`。
  * `entry_list` - 与ACL关联的条目（IP地址或CIDR块）列表。每个条目包含：
    * `entry` - IP地址或CIDR块。
    * `comment` - 与此条目关联的注释。
  * `related_listeners` - 附加到此ACL的监听器列表。每个监听器包含：
    * `load_balancer_id` - 监听器所属的负载均衡器实例的ID。
    * `frontend_port` - 监听器的端口号。
    * `protocol` - 监听器使用的协议（例如，TCP、UDP、HTTP、HTTPS等）。
    * `acl_type` - 应用于监听器的ACL类型（例如，白名单/黑名单）。
  * `tags` - 分配给ACL的标签映射。