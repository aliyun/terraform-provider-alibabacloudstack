---
subcategory: "Elasticsearch"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_elasticsearch_instance"
sidebar_current: "docs-Alibabacloudstack-elasticsearch-instance"
description: |-
  提供一个在阿里云上的Elasticsearch实例资源。
---

# alibabacloudstack_elasticsearch_instance

提供一个在阿里云上的Elasticsearch实例资源。

## 概述

`alibabacloudstack_elasticsearch_instance` 资源允许您在阿里云上管理Elasticsearch实例。

## 示例用法
```hcl
resource "alibabacloudstack_elasticsearch_instance" "example" {
  description            = "example_elasticsearch"
  vswitch_id             = "vsw-1234567890abcdef0"
  version                = "6.7_with_X-Pack"
  instance_charge_type   = "PostPaid"
  data_node_amount       = 2
  data_node_spec         = "elasticsearch.sn2ne.large"
  data_node_disk_size    = 20
  data_node_disk_type    = "cloud_efficiency"
  enable_public          = true
  public_whitelist       = ["0.0.0.0/0"]
  enable_kibana_public_network = true
  kibana_whitelist       = ["0.0.0.0/0"]
  zone_count             = 2
}
```

## 参数说明
以下参数是支持的：

* `description` - (可选) - Elasticsearch实例的描述。长度必须在0到30个字符之间，可以包含数字、字母、下划线和连字符。必须以字母、数字或中文字符开头。
* `vswitch_id` - (必选) - 启动Elasticsearch实例的VSwitch的ID。
* `password` - (可选) - Elasticsearch实例的密码。
* `kms_encrypted_password` - (可选) - 用于Elasticsearch实例的KMS加密密码。
* `kms_encryption_context` - (可选) - 用于为KMS加密提供额外上下文的键值对映射。
* `version` - (必选) - 要部署的Elasticsearch版本。
* `tags` - (可选) - 要分配给资源的标签映射。
* `instance_charge_type` - (可选) - 实例的计费类型。有效值为 PrePaid 和 PostPaid。默认值为 PostPaid。
* `period` - (可选) - 实例的订阅周期。有效值为 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36。默认值为 1。
* `data_node_amount` - (必选) - Elasticsearch集群中的数据节点数量。有效范围是从 2 到 50。
* `data_node_spec` - (必选) - 数据节点的规格。
* `data_node_disk_size` - (必选) - 数据节点的磁盘大小。
* `data_node_disk_type` - (必选) - 数据节点的磁盘类型。
* `data_node_disk_encrypted` - (可选) - 数据节点磁盘是否加密。默认值为 false。
* `private_whitelist` - (可选) - Elasticsearch实例的私有IP白名单。
* `enable_public` - (可选) - 是否启用对Elasticsearch实例的公共访问。默认值为 false。
* `public_whitelist` - (可选) - Elasticsearch实例的公共IP白名单。
* `master_node_spec` - (可选) - 主节点的规格。
* `client_node_amount` - (可选) - Elasticsearch集群中的客户端节点数量。有效范围是从 2 到 25。
* `client_node_spec` - (可选) - 客户端节点的规格。
* `protocol` - (可选) - Elasticsearch实例使用的协议。有效值为 HTTP 和 HTTPS。默认值为 HTTP。
* `zone_count` - (可选) - 部署Elasticsearch实例的可用区数量。有效范围是从 1 到 3。默认值为 1。
* `resource_group_id` - (可选) - Elasticsearch实例所属的资源组ID。
* `setting_config` - (可选) - 用于配置Elasticsearch实例的设置映射。

## 属性说明
以下属性会被导出：

* `domain` - Elasticsearch实例的域名。
* `port` - Elasticsearch实例的端口号。
* `status` - Elasticsearch实例的状态。
* `kibana_domain` - Kibana实例的域名。
* `kibana_port` - Kibana实例的端口号。
* `enable_kibana_public_network` - 是否启用对Kibana实例的公共访问。
* `enable_kibana_private_network` - 是否启用对Kibana实例的私有访问。
* `kibana_whitelist` - Kibana实例的公共IP白名单。
* `kibana_private_whitelist` - Kibana实例的私有IP白名单。

## 导入
Elasticsearch实例可以使用ID进行导入，例如：

```bash
$ terraform import alibabacloudstack_elasticsearch_instance.example i-1234567890abcdef0
```