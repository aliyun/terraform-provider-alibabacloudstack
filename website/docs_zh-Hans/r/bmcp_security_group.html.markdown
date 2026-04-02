---
subcategory: "Bare Metal Compute Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_security_group"
sidebar_current: "docs-Alibabacloudstack-bmcp-security-group"
description: |- 
  编排BMCP安全组资源
---

# alibabacloudstack_bmcp_security_group

编排BMCP（Bare Metal Compute Platform，裸金属算力平台）安全组资源。

## 示例用法

```hcl
resource "alibabacloudstack_bmcp_security_group" "default" {
  vpc_id      = "vpc-7kkaxv72063xkz7exgv3x"
  name        = "my-security-group"
  description = "我的BMCP安全组"
}
```

## 参数说明

支持以下参数：

* `vpc_id` - (必填) 安全组所在的 VPC 的 ID。
* `name` - (必填) 安全组的名称。长度必须为2到128个字符。
* `description` - (选填) 安全组的描述。长度必须为0到256个字符。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 安全组的ID。
* `sg_id` - 安全组的ID。
* `name` - 安全组的名称。
* `description` - 安全组的描述。
* `vpc_id` - 安全组所在的 VPC 的 ID。
* `resource_group` - 安全组的资源组ID。
* `resource_group_name` - 安全组的资源组名称。
* `department` - 安全组的部门ID。
* `department_name` - 安全组的部门名称。
* `region_id` - 安全组的地域ID。
* `ascm_create_user` - 创建安全组的用户。
* `create_time` - 安全组的创建时间。
* `update_time` - 安全组的最后更新时间。
