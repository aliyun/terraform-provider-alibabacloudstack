---
subcategory: "企业控制台(ASCM)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_policy_for_role"
sidebar_current: "docs-alibabacloudstack-resource-ascm-ram-policy-for-role"
description: |-
    编排绑定ASCM的RAM策略和RAM角色
---

# alibabacloudstack_ascm_ram_policy_for_role

使用Provider配置的凭证在ASCM中绑定RAM策略到RAM角色。

-> **注意:** 此资源也可以使用以下别名引用：`apsarastack_ascm_ram_policy_for_role`。

## 示例用法

```
resource "alibabacloudstack_ascm_ram_policy" "default" {
  name = "Testpolicy"
  description = "Testing Complete"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"

}

resource "alibabacloudstack_ascm_ram_role" "default" {
  role_name = "TestRole"
  description = "TestingRole"
  organization_visibility = "organizationVisibility.global"
}

resource "alibabacloudstack_ascm_ram_policy_for_role" "default" {
  ram_policy_id = alibabacloudstack_ascm_ram_policy.default.ram_id
  role_id = alibabacloudstack_ascm_ram_role.default.role_id
}
```

## 参数说明

支持以下参数：

* `ram_policy_id` - (必填，变更时重建) 要绑定到角色的RAM策略ID。更改此参数会强制创建新资源。
* `role_id` - (必填，变更时重建) 要绑定策略的RAM角色ID。更改此参数会强制创建新资源。

## 属性说明

此资源目前未定义任何输出属性。

## Import

ASCM RAM策略角色绑定可以使用 ram_policy_id 和 role_id 导入，以冒号分隔，例如：

```
$ terraform import alibabacloudstack_ascm_ram_policy_for_role.example <ram_policy_id>:<role_id>
```