---
subcategory: "堡垒机"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bastionhost_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-bastionhost-instances"
description: |-
 根据提供的过滤条件提供一组 Alibaba Cloud Stack Bastionhost 实例。
---


示例用法
以下示例通过描述正则表达式检索 Bastionhost 实例，并将结果写入文件：

```hcl
data "alibabacloudstack_bastionhost_instances" "example" {
  description_regex = "^example"
  output_file       = "output.json"
}

output "instances" {
  value = data.alibabacloudstack_bastionhost_instances.example.ids
}
```
## 参数说明
支持以下参数：

* `description_regex` - (可选) 用于根据实例描述进行过滤的正则表达式。
* `output_file` - (可选) 将结果保存为 JSON 格式的文件名。
* `Deprecated`: 此字段已弃用，将在版本 3.19.0 中移除。请改用 local_file 提供程序。
* `ids` - (可选, ForceNew) 要过滤的一组 Bastionhost 实例 ID。
* `tags` - (可选) 用于按标签过滤 Bastionhost 实例的键值对映射。
## 属性参考
除了上述所有参数外，还导出以下属性：

* `ids` - 匹配条件的 Bastionhost 实例 ID 列表。
* `descriptions` - 对应匹配实例的描述列表。
* `instances` - 包含详细属性的 Bastionhost 实例列表。每个条目包含：
* `id` - Bastionhost 实例的 ID。
* `description` - 实例的描述。
* `user_vswitch_id` - 实例关联的 VSwitch ID。
* `private_domain` - 实例的私有域名或内网端点。
* `public_domain` - 实例的公共域名或外网端点（如果可用）。
* `instance_status` - 实例的当前状态。
* `license_code` - 实例的许可证代码。
* `public_network_access` - 布尔值，表示是否启用公网访问。
* `security_group_ids` - 实例关联的安全组 ID 列表。
* `tags` - 分配给实例的标签。


