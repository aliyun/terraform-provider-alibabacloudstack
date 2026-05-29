---
subcategory: "企业控制台(ASCM)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_organization"
sidebar_current: "docs-Alibabacloudstack-resource-ascm-organization"
description: |-
  编排 ASCM 组织资源。
---

# alibabacloudstack_ascm_organization

使用 Provider 配置的凭证编排 ASCM 组织资源。

-> **Note:** 该资源也可以使用以下别名引用：
- `apsarastack_ascm_organization`

## 示例用法

```hcl
resource "alibabacloudstack_ascm_organization" "default" {
  name      = "apsara_Organization"
  parent_id = "1"
}
```

## 参数说明

以下是支持的参数：

* `name` - (必填) 组织的名称。该名称长度必须为 2 到 128 个字符。
* `parent_id` - (可选) 父组织的 ID。默认值为 `"1"`。
* `person_num` - (可选) 保留参数，目前暂无实际用途。
* `resource_group_num` - (可选) 保留参数，目前暂无实际用途。

## 属性说明

以下属性被导出：

* `id` - 组织的 ID。
* `org_id` - 组织的 UUID。
* `primary_key` - 组织关联的主键。
* `aliyunid` - 组织关联的 Aliyun ID。

## Import

ASCM 组织可以通过组织 ID 导入，例如：

```
$ terraform import alibabacloudstack_ascm_organization.example 12345
```