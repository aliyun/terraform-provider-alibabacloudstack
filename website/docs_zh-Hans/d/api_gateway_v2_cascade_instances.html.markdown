---
subcategory: "API 网关（API Gateway）V2 版"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_cascade_instances"
description: |-
  提供API网关V2级联实例列表
---

# alibabacloudstack_api_gateway_v2_cascade_instances

该数据源提供当前用户在阿里云上的API网关V2级联实例列表。

## 示例用法

基础用法

```hcl
data "alibabacloudstack_api_gateway_v2_cascade_instances" "example" {
  
}

output "first_api_gateway_v2_cascade_instance_id" {
  value = "${data.alibabacloudstack_api_gateway_v2_cascade_instances.example.instances.0.id}"
}
```

## 参数参考

以下参数被支持：

* `ids` - (可选) 级联实例ID列表。
* `name_regex` - (可选) 用于根据级联实例名称过滤结果的正则表达式。

## 属性参考

以下属性会被导出：

* `ids` - 级联实例ID列表。
* `names` - 级联实例名称列表。
* `instances` - 级联实例列表。每个元素包含以下属性：
  * `id` - 级联实例的ID。
  * `instance_type` - 级联实例的类型。
  * `instance_name` - 级联实例的名称。
  * `cascade_instance_id` - 级联实例ID。
  * `create_time` - 级联实例的创建时间。
  * `status` - 级联实例的状态。