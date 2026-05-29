---
subcategory: "Prometheus 监控服务"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-prometheus-v2-instances"
description: |-
  查询阿里云Prometheus V2实例列表
---

# alibabacloudstack_prometheus_v2_instances

查询阿里云Prometheus V2实例列表。该数据源用于检索已创建的Prometheus V2监控实例，支持通过ID列表和名称正则表达式进行过滤。

## 示例用法

```hcl

variable "name" {
  default = "tfacc_prometheus9102323273752963708"
}

resource "alibabacloudstack_prometheus_v2_instance" "default" {
  cluster_name = var.name
  tags         = ["test1", "test2"]
}

data "alibabacloudstack_prometheus_v2_instances" "default" {
  name_regex = alibabacloudstack_prometheus_v2_instance.default.cluster_name
}

```

## 参数说明
以下参数支持过滤查询结果：

* `ids` (列表, 可选)：实例ID列表，用于精确过滤需要查询的Prometheus V2实例。

* `name_regex` (字符串, 可选)：实例名称的正则表达式，用于模糊匹配需要查询的Prometheus V2实例名称。

## 属性说明
以下属性被导出：

* `id` (字符串)：Prometheus V2实例的唯一标识符，等同于cluster_id。

* `cluster_id` (字符串)：Prometheus V2实例的ID。

* `cluster_name` (字符串)：Prometheus V2实例的名称。

* `cluster_type` (字符串)：Prometheus V2实例的类型。

* `http_api` (字符串)：Prometheus V2实例的HTTP API端点地址。

* `objid` (整数)：Prometheus V2实例的内部对象ID。

* `push_gateway_url` (字符串)：Prometheus V2实例的推送网关URL地址。

* `remote_write_url` (字符串)：Prometheus V2实例的远程写入URL地址。

* `security_level_tag` (字符串)：实例的安全级别标签。

* `status` (字符串)：实例的当前状态。

* `tag_set` (列表)：分配给实例的标签列表。