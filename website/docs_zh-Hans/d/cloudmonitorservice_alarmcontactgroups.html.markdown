---
subcategory: "Cloud Monitor Service (CMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_alarmcontactgroups"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudmonitorservice-alarmcontactgroups"
description: |- 
  查询云监控报警联系人组
---

# alibabacloudstack_cloudmonitorservice_alarmcontactgroups
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cms_alarm_contact_groups`

根据指定过滤条件列出当前凭证权限可以访问的云监控报警联系人组列表。

## 示例用法

### 基础用法

```hcl
data "alibabacloudstack_cloudmonitorservice_alarmcontactgroups" "example" {
  name_regex = "tf-testacc"
}

output "alarm_contact_groups" {
  value = data.alibabacloudstack_cloudmonitorservice_alarmcontactgroups.example.groups
}
```

## 参数说明

以下参数是支持的：

* `name_regex` - （可选，变更时重建）用于通过报警联系人组名称筛选结果的正则表达式字符串。这可以帮助您精确地定位符合条件的报警联系人组。
* `ids` - （可选，变更时重建）报警联系人组 ID 列表。您可以使用此参数来限定返回的报警联系人组范围。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 报警联系人组 ID 列表。
* `names` - 报警联系人组名称列表。
* `groups` - CloudMonitorService 报警联系人组列表。每个元素包含以下属性：
  * `id` - 报警联系人组的 ID。
  * `alarm_contact_group_name` - 报警联系人组的名称。
  * `contacts` - 与此组关联的报警联系人列表。每个联系人包含以下属性：
    * `contact_id` - 联系人的 ID。
    * `contact_name` - 联系人的名称。
  * `describe` - 报警联系人组的描述信息。
  * `enable_subscribed` - 指示报警组是否订阅每周报告。有效值为 `true` 或 `false`。如果设置为 `true`，报警组将接收每周摘要报告。