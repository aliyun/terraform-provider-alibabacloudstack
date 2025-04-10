---
subcategory: "CloudMonitorService"  
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

### 基础用法：  

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

* `ids` - （可选，变更时重建）报警联系人组 ID 列表。通过此参数可以限定返回结果中包含的报警联系人组范围。  
* `name_regex` - （可选，变更时重建）用于通过正则表达式匹配报警联系人组名称的字符串。通过此参数可以筛选出符合条件的报警联系人组名称。  

## 属性说明  

除了上述参数外，还导出以下属性：  

* `ids` - 报警联系人组 ID 列表。  
* `names` - 报警联系人组名称列表。  
* `groups` - CloudMonitorService 报警联系人组列表。每个元素包含以下属性：  
  * `id` - 报警联系人组的唯一标识符。  
  * `alarm_contact_group_name` - 报警联系人组的名称。  
  * `contacts` - 与此组关联的报警联系人列表。每个联系人包含以下属性：  
    * `contact_id` - 联系人的唯一标识符。  
    * `contact_name` - 联系人的名称。  
  * `describe` - 报警联系人组的描述信息，提供关于该组的详细说明。  
  * `enable_subscribed` - 指示报警组是否订阅每周报告的功能。有效值为 `true` 或 `false`。如果设置为 `true`，报警组将接收每周生成的摘要报告。  