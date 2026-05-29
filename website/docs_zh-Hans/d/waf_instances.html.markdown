---
subcategory: "Web应用防火墙" 
layout: "alibabacloudstack" 
page_title: "Alibabacloudstack: alibabacloudstack_waf_instances" 
sidebar_current: "docs-alibabacloudstack-resource_waf_instances" 
description: |- 
提供一个基于过滤条件的 Alibaba Cloud Stack WAF 实例列表。
---

# alibabacloudstack_waf_instances
-> **NOTE:** 该资源等效别名有:：alibabacloudstack_waf_instances

根据提供的过滤条件提供 Alibaba Cloud Stack 中的 WAF 实例列表。

## 示例用法
hcl
data "alibabacloudstack_waf_instances" "example" {
  ids = ["waf-instance-1", "waf-instance-2"]
  output_file = "output.json"
}

output "instances" {
  value = data.alibabacloudstack_waf_instances.example.ids
}
## 参数说明
支持以下参数：

* `ids` - (可选，强制更新) 要筛选结果的 WAF 实例 ID 列表。
* `name` - (输出属性) WAF 实例的名称。
* `instance_status` - (输出属性) WAF 实例的当前状态。
* `instance_make_status` - (输出属性) WAF 实例的创建或供应状态。
* `output_file` - (可选) 将结果保存为 JSON 格式的文件名。
* `Deprecated`: 此字段已弃用，并将在版本 3.19.0 中删除。请改用 local_file 提供程序。
* `vpc_vswitch` - (输出属性) 包含 WAF 实例的 VPC 和 vSwitch 配置的列表，其结构如下：
* `vswitch_name` - vSwitch 的名称。
* `vswitch` - vSwitch 的 ID。
* `cidr_block` - vSwitch 的 CIDR 网段。
* `available_zone` - vSwitch 所在的可用区。
* `vpc` - 关联 VPC 的 ID。
* `vpc_name` - 关联 VPC 的名称。
* `detector_specs` - (输出属性) 检测引擎的规格。
* `detector_version` - (输出属性) 检测引擎的版本或等级。
* `detector_nodenum` - (输出属性) 单个可用区内部署的检测引擎节点数。
属性输出
除以上所有参数外，还导出以下属性：

* `ids` - 符合条件的 WAF 实例 ID 列表。
* `instances` - 具有详细属性的 WAF 实例列表。每个条目包含：
* `name` - (输出属性) WAF 实例的名称。
* `instance_status` - (输出属性) WAF 实例的当前状态。
* `instance_make_status` - (输出属性) WAF 实例的创建或供应状态。
* `vpc_vswitch` - (输出属性) 包含 WAF 实例的 VPC 和 vSwitch 配置的列表。
* `detector_specs` - (输出属性) 检测引擎的规格。
* `detector_version` - (输出属性) 检测引擎的版本或等级。
* `detector_nodenum` - (输出属性) 单个可用区内部署的检测引擎节点数。