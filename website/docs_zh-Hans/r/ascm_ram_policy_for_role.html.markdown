---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_policy_for_role"
sidebar_current: "docs-alibabacloudstack-resource-ascm-ram-policy-for-role"
description: |-
    编排绑定ASCM的RAM策略和RAM角色
---

# alibabacloudstack_ascm_ram_policy_for_role

使用Provider配置的凭证在指定的资源集下编排绑定ASCM的RAM策略和RAM角色。

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
output "ramrolebinder" {
  value = alibabacloudstack_ascm_ram_policy_for_role.default.*
}

```

## 参数说明

支持以下参数：

* `ram_policy_id` - (必填) 要绑定的RAM策略的ID。
* `role_id` - (必填，变更时重建) 要绑定的角色ID。

## 属性说明

此资源目前未定义任何输出属性。