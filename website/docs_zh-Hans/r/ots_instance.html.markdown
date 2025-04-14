---
subcategory: "OTS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_instance"
sidebar_current: "docs-Alibabacloudstack-ots-instance"
description: |- 
  编排表格存储服务(OTS）实例
---

# alibabacloudstack_ots_instance

使用Provider配置的凭证在指定的资源集编排表格存储服务(OTS）实例。

## 示例用法

```hcl
# 创建一个 OTS 实例
resource "alibabacloudstack_ots_instance" "foo" {
  name          = "my-ots-instance"
  description   = "This is a test OTS instance"
  accessed_by   = "Vpc" # 可选值: Any, Vpc, ConsoleOrVpc. 默认值: Any
  instance_type = "Capacity" # 可选值: Capacity, HighPerformance. 默认值: HighPerformance
  tags = {
    Created = "TF"
    For     = "Building table"
  }
}
```

## 参数说明

支持以下参数：

* `name` - (必填，变更时重建) OTS 实例的名称。更改此参数将强制创建新资源。
* `accessed_by` - (可选) 访问 OTS 实例的网络限制。有效值：
  * `Any` - 允许所有网络访问实例。
  * `Vpc` - 仅允许从附加的 VPC 访问。
  * `ConsoleOrVpc` - 允许从 Web 控制台或附加的 VPC 访问。
  
  默认值：`Any`。

* `instance_type` - (可选，变更时重建) OTS 实例的类型。有效值：
  * `Capacity` - 适用于大数据量和高吞吐量需求的场景。
  * `HighPerformance` - 适用于需要低延迟和高性能的场景。
  
  默认值：`HighPerformance`。更改此参数将强制创建新资源。

* `description` - (可选，变更时重建) OTS 实例的简要描述。此字段在创建后无法修改。更改此参数将强制创建新资源。
* `tags` - (可选) 分配给 OTS 实例的标签映射。
* `propreties` - (可选) OTS 实例的其他属性配置。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - OTS 实例的 ID。其值与 `name` 相同。
* `name` - OTS 实例的名称。
* `description` - OTS 实例的描述。
* `accessed_by` - 访问 OTS 实例的网络限制。
* `instance_type` - OTS 实例的类型。
* `tags` - 分配给 OTS 实例的标签映射。
* `propreties` - OTS 实例的计算属性。

## 导入

OTS 实例可以使用实例 ID 或名称导入，例如：

```bash
$ terraform import alibabacloudstack_ots_instance.foo "my-ots-instance"
``` 