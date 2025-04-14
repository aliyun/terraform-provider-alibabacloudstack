---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_accesscontrollist"
sidebar_current: "docs-Alibabacloudstack-slb-accesscontrollist"
description: |- 
  编排负载均衡(SLB)访问控制
---

# alibabacloudstack_slb_accesscontrollist
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_acl`

使用Provider配置的凭证在指定的资源集编排负载均衡(SLB)访问控制。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccslbaccess_control_list69490"
}

resource "alibabacloudstack_slb_accesscontrollist" "default" {
  acl_name            = "Rdk_test_name01"
  address_ip_version  = "ipv4"

  entry_list {
    entry   = "10.10.10.0/24"
    comment = "first"
  }

  entry_list {
    entry   = "168.10.10.0/24"
    comment = "second"
  }

  tags = {
    CreatedBy = "Terraform"
    Purpose   = "Testing"
  }
}
```

## 参数说明

支持以下参数：
  * `acl_name` - (必填) 访问控制策略组名称。
  * `address_ip_version` - (可选，变更时重建) 关联的实例的IP类型。有效值为 `ipv4` 和 `ipv6`。默认值为 `ipv4`。
  * `entry_list` - (可选) 要添加的条目(IP地址或CIDR块)列表。一个资源中最多可以支持50个条目。每个条目包含：
    * `entry` - (必填) 一个IP地址或CIDR块。
    * `comment` - (可选) 该条目的注释。
  * `tags` - (可选) 分配给资源的标签映射。
  * `resource_group_id` - (可选，变更时重建) 资源组的ID。
  * `name` - (可选) 已废弃的名称字段。字段 `name` 已被废弃，并将在未来的版本中移除。请改用新的字段 `acl_name`。
  * `ip_version` - (可选，变更时重建) 已废弃的IP版本字段。字段 `ip_version` 已被废弃，并将在未来的版本中移除。请改用新的字段 `address_ip_version`。

## 属性说明

除了上述所有参数外，还导出了以下属性：
  * `id` - 访问控制列表的ID。
  * `acl_name` - 访问控制策略组名称。
  * `address_ip_version` - 关联的实例的IP类型。有效值为 `ipv4` 和 `ipv6`。
  * `name` - 已废弃的名称字段。
  * `ip_version` - 已废弃的IP版本字段。