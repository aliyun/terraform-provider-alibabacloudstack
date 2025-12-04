---
subcategory: "Prometheus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_instance"
sidebar_current: "docs-Alibabacloudstack-prometheus-prometheus_v2_instance"
description: |-
  管理Prometheus V2实例资源
---

# alibabacloudstack_prometheus_v2_instance

管理阿里云Prometheus V2实例资源，用于创建、读取、更新和删除Prometheus监控实例。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tfacc_prometheus15993"
}


resource "alibabacloudstack_prometheus_v2_instance" "default" {
  cluster_name = var.name
  tags = [
    "test1",
    "test2"
  ]
}
```

## 参数说明

支持以下参数：

* `cluster_name` - (必填, 变更时重建) Prometheus实例的名称。名称长度需符合API规范，创建后无法直接修改，修改时将触发资源重建。
* `tags` - (可选) 实例的标签集合。支持通过标签管理资源，更新时仅影响标签配置。

## 属性说明

以下属性导出为资源状态：

* `id` - 实例的唯一标识符（ClusterId）。
* `cluster_id` - Prometheus实例的集群ID，与id属性值相同。
* `http_api` - Prometheus HTTP API的访问地址，用于查询和管理监控数据。
* `objid` - 实例的内部对象ID（整数类型），用于API操作标识。
* `push_gateway_url` - Push Gateway的访问地址，用于接收监控指标推送。
* `remote_write_url` - Remote Write的访问地址，用于远程写入监控数据。