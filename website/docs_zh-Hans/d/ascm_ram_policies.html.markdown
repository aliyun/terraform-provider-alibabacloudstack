---
subcategory: "ASCM"  
layout: "alibabacloudstack"  
page_title: "AlibabacloudStack: alibabacloudstack_ascm_ram_policies"  
sidebar_current: "docs-alibabacloudstack-datasource-ascm-ram-policies"  
description: |-
    查询RAM策略  

---

# alibabacloudstack_ascm_ram_policies  

根据指定过滤条件列出当前凭证权限可以访问的所有RAM策略列表。  

## 示例用法  

```
resource "alibabacloudstack_ascm_ram_policy" "default" {  
  name = "TestPolicy"  
  description = "Testing"  
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"  

}  
data "alibabacloudstack_ascm_ram_policies" "default" {  
  name_regex = alibabacloudstack_ascm_ram_policy.default.name  
}  
output "ram_policies" {  
  value = data.alibabacloudstack_ascm_ram_policies.default.*  
}  
```  

## 参数说明  

支持以下参数：  

* `ids` - (可选) RAM策略ID列表。用于通过策略ID筛选结果。  
* `name_regex` - (可选) 用于通过RAM策略名称过滤结果的正则表达式字符串。可以通过该参数匹配特定命名规则的策略。  
* `region` - (可选) 策略所属的区域名称。如果指定了区域，则只返回该区域下的策略。  

## 属性说明  

除了上述参数外，还导出以下属性：  

* `policies` - 策略列表。每个元素包含以下属性：  
    * `id` - 策略的ID，唯一标识一个RAM策略。  
    * `name` - 策略名称，用于标识策略的名称。  
    * `description` - 策略的描述信息，提供关于策略用途或功能的详细说明。  
    * `ctime` - RAM策略的创建时间，格式为标准时间戳。  
    * `cuser_id` - 策略创建者的ID，标识创建该策略的用户。  
    * `region` - 策略所属的区域名称，表示该策略适用的区域。  
    * `policy_document` - 策略文档，定义了策略的具体权限和规则。  
    * `output_file` - （可选）保存数据源结果的文件名（在运行`terraform plan`之后）。可以通过该属性将查询结果保存到指定文件中，便于后续使用或记录。  