---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_accessrule"
sidebar_current: "docs-Alibabacloudstack-nas-accessrule"
description: |- 
  编排文件存储（NAS）访问规则
---

# alibabacloudstack_nas_accessrule
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_nas_access_rule`

使用Provider配置的凭证在指定的资源集编排文件存储（NAS）访问规则。

## 示例用法

### 基础用法

```terraform
variable "name" {
    default = "tf-testaccnasaccess_rule99652"
}

resource "alibabacloudstack_nas_access_group" "example" {
  access_group_name = "tf-NasConfigName"
  access_group_type = "Vpc"
  description       = "tf-testAccNasConfig"
}

resource "alibabacloudstack_nas_accessrule" "default" {
  access_group_name = alibabacloudstack_nas_access_group.example.access_group_name
  source_cidr_ip    = "1.1.1.1/0"
  rw_access_type    = "RDWR"
  user_access_type  = "no_squash"
  priority          = 1
}
```

## 参数参考

支持以下参数：

* `access_group_name` - (必填, 变更时重建) - 权限组名称。此参数是必填的，并且一旦设置后无法修改。
* `source_cidr_ip` - (必填) - 您希望允许访问 NAS 文件系统的地址或地址段。例如，`1.1.1.1/0`。
* `rw_access_type` - (选填) - 规则的读写权限类型。有效值为：
  * `RDWR`: 读写访问(默认)。
  * `RDONLY`: 只读访问。
* `user_access_type` - (选填) - 规则的用户权限类型。有效值为：
  * `no_squash`: 对 root 用户无限制(默认)。
  * `root_squash`: 限制 root 用户拥有完全权限。
  * `all_squash`: 限制所有用户拥有完全权限。
* `priority` - (选填) - 规则的优先级。有效范围为 1-100。数字越小表示优先级越高。默认值为 `1`。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 此资源的 ID。格式为 `<access_group_name>:<access_rule_id>`。
* `access_rule_id` - NAS 访问规则的唯一标识符。