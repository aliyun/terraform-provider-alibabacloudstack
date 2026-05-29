---
subcategory: "裸金属算力平台 BMCP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bcmp_security_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-bcmp-security-groups"
description: |- 
  获取BMCP安全组列表
---

# alibabacloudstack_bcmp_security_groups

该数据源根据指定的过滤器获取AlibabacloudStack账户中的BMCP（Bare Metal Compute Platform，裸金属算力平台）安全组列表。

## 示例用法

```hcl
# 声明数据源
data "alibabacloudstack_bcmp_security_groups" "example" {
  vpc_id = "vpc-7kkaxv72063xkz7exgv3x"
  name   = "my-security-group"
}

output "security_group_ids" {
  value = data.alibabacloudstack_bcmp_security_groups.example.ids
}

output "security_group_names" {
  value = data.alibabacloudstack_bcmp_security_groups.example.names
}
```

## 参数说明

支持以下参数：

* `vpc_id` - (选填) 用于过滤结果的 VPC ID。
* `name` - (选填) 用于过滤结果的安全组名称。
* `name_regex` - (选填) 用于按安全组名称过滤结果的正则表达式字符串。
* `ids` - (选填) 用于过滤结果的安全组 ID 列表。如果未指定，将考虑所有安全组。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `names` - 匹配的安全组名称列表。
* `ids` - 匹配的安全组 ID 列表。
* `security_groups` - 匹配的安全组列表。每个元素包含以下属性：
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
