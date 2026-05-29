---
subcategory: "实时计算 Flink 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack:alibabacloudstack_flink_namespace"
sidebar_current: "docs-Alibabacloudstack-resource-flink-namespace"
description: |-
  编排实时计算命名空间
---

# alibabacloudstack_flink_namespace

使用Provider配置的凭证在指定的资源集下编排实时计算命名空间。


## 示例用法

### 基础用法

```
variable name{
 default = "<Your NameSpace Name>"
}

resource "alibabacloudstack_ascm_user" "default" {
  display_name = var.name
  mobile_nation_code = "86"
  login_name = var.name
  login_policy_id = "1"
  role_ids = [
               "8",
               "9"
             ]
  cellphone_number = "13612345678"
  email = "${var.name}@gmail.com"
}


resource "alibabacloudstack_flink_namespace" "default" {
  owner_uid = "${alibabacloudstack_ascm_user.default.user_uid}"
  name      = var.name
  cu        = 1
  cpu_type  = "Intel"
}
```

## 参数说明

支持以下参数：

* `name` - (必填，变更时重建) 仓库命名空间的名称。
* `cu` -  (必填) 仓库命名空间保留的资源数量。
* `cpu_type` - (必填，变更时重建) 资源CPU类型。有效值：`Intel`。
* `owner_uid` -  (可选，变更时重建) 仓库命名空间管理员Id。

## 属性说明

导出以下属性：

* `id` - Flink仓库命名空间的唯一标识符。其值与 `name` 参数相同。


## 导入
Flink的命名空间可以使用ID进行导入，例如：

```
$ terraform import alibabacloudstack_flink_namespace.default namespace_name
```