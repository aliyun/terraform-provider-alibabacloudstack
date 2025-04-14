---
subcategory: "ASCM"  
layout: "alibabacloudstack"  
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ecs_instance_families"  
sidebar_current: "docs-alibabacloudstack-datasource-ascm-ecs-instance-families"  
description: |-
    查询ECS实例族  
---

# alibabacloudstack_ascm_ecs_instance_families  

根据指定过滤条件列出当前凭证权限可以访问的ECS实例族列表。  

## 示例用法  

```
data "alibabacloudstack_ascm_ecs_instance_families" "default" {  
  status = "Available"  
  output_file = "ecs_instance"  
}  
output "ecs_instance" {  
  value = data.alibabacloudstack_ascm_ecs_instance_families.default.*  
}  
```  

## 参数说明  

以下参数被支持：  

* `ids` - (可选) ECS实例族ID列表，用于指定需要查询的ECS实例族。  
* `status` - (必填) 指定ECS实例族的状态来过滤结果。例如，可以设置为 `"Available"` 来获取可用的ECS实例族。  

## 属性说明  

除了上述参数外，还导出以下属性：  

* `families` - ECS实例族列表。每个元素包含以下属性：
    * `instance_type_family_id` - ECS实例族的ID。
    * `generation` - ECS实例族的代数。