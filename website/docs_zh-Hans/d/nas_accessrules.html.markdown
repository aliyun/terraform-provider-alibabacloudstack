---
subcategory: "NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_accessrules"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-accessrules"
description: |- 
  查询文件存储（NAS）访问规则列表
---

# alibabacloudstack_nas_accessrules
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_nas_access_rules`

根据指定过滤条件列出当前凭证权限可以访问的文件存储（NAS）访问规则列表

## 示例用法

```hcl
variable "name" {
    default = "tf-testAccAccessGroupsdatasource-866785"
}

resource "alibabacloudstack_nas_access_group" "default" {
    access_group_name = var.name
    access_group_type = "Vpc"
    description       = "tf-testAccAccessGroupsdatasource"
}

resource "alibabacloudstack_nas_access_rule" "default" {
    access_group_name = alibabacloudstack_nas_access_group.default.access_group_name
    source_cidr_ip    = "168.1.1.0/16"
    rw_access_type    = "RDWR"
    user_access_type  = "no_squash"
    priority          = 2
}

data "alibabacloudstack_nas_access_rules" "default" {
    access_group_name = alibabacloudstack_nas_access_group.default.access_group_name
    source_cidr_ip    = alibabacloudstack_nas_access_rule.default.source_cidr_ip
}

output "access_rule_id" {
    value = data.alibabacloudstack_nas_access_rules.default.rules[0].access_rule_id
}
```

## 参数参考

以下参数是支持的：

* `access_group_name` - (必填，变更时重建) 权限组名称。这是一个必填字段，创建后无法修改。
* `source_cidr_ip` - (选填) 地址或地址段，用于指定访问规则的IP范围的CIDR块。
* `rw_access` - (选填) 授权对象对文件系统的读写权限。有效值包括：
  * `RDONLY`: 只读访问。
  * `RDWR`: 读写访问。
* `user_access` - (选填) 授权对象的系统用户对文件系统的访问权限。有效值包括：
  * `no_squash`: 不进行根压缩。
  * `root_squash`: 根压缩。
  * `all_squash`: 全部压缩。
* `ids` - (选填) 用于过滤结果的NAS访问规则ID列表。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - NAS访问规则ID列表。
* `rules` - NAS访问规则列表。每个元素包含以下属性：
  * `source_cidr_ip` - 地址或地址段，用于指定访问规则的IP范围的CIDR块。
  * `priority` - 访问规则的优先级。
  * `access_rule_id` - 访问规则的ID。
  * `user_access` - 授权对象的系统用户对文件系统的访问权限。
  * `rw_access` - 授权对象对文件系统的读写权限。