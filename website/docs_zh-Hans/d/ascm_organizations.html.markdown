---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_organizations"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-organizations"
description: |-
    查询ascm组织
---

# alibabacloudstack_ascm_organizations  

根据指定过滤条件列出当前凭证权限可以访问的组织列表。  

## 示例用法  

```
resource "alibabacloudstack_ascm_organization" "default" {  
  name = "Test_org"  
}  
output "orgres" {  
  value = alibabacloudstack_ascm_organization.default.*  
}  
data "alibabacloudstack_ascm_organizations" "default" {  
    name_regex = alibabacloudstack_ascm_organization.default.name  
    parent_id = alibabacloudstack_ascm_organization.default.parent_id  
}  
output "orgs" {  
  value = data.alibabacloudstack_ascm_organizations.default.*  
}  
```  

## 参数说明  

以下是支持的参数：  

* `ids` - (可选) 组织ID列表。可以通过此参数指定需要查询的组织ID集合。  
* `name_regex` - (可选) 用于按组织名称过滤结果的正则表达式字符串。通过此参数，可以根据组织名称进行模糊匹配。  
* `parent_id` - (可选) 通过指定的组织父级ID过滤结果。此参数用于筛选特定父级下的子组织。  

## 属性说明  

除了上述列出的参数外，还导出以下属性：  

* `organizations` - 组织列表。每个元素包含以下属性：  
  * `id` - 组织的唯一标识符。  
  * `name` - 组织的名称。  
  * `cuser_id` - 创建该组织的用户的ID。  
  * `muser_id` - 最后修改该组织的用户的ID。  
  * `alias` - 组织的别名，通常用于简化或替代正式名称。  
  * `parent_id` - 该组织所属的父级组织ID。如果为顶级组织，则此值为空。  
  * `internal` - 表示该组织是否为内部组织。可能的值为 `true` 或 `false`。  

请根据实际需求使用这些参数和属性来查询和管理组织信息。