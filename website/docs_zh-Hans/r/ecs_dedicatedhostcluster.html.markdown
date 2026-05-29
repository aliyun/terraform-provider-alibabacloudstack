---
subcategory: "云服务器 ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_dedicatedhostcluster"
sidebar_current: "docs-Alibabacloudstack-ecs-dedicatedhostcluster"
description: |-
  Provides a ecs Dedicatedhostcluster resource.
---

# alibabacloudstack\_ecs\_dedicatedhostcluster

Provides a ecs Dedicatedhostcluster resource.

## 示例用法
```
variable "name" {
    default = "tf-testaccecsdedicated_hostcluster20348"
}



data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}





resource "alibabacloudstack_ecs_dedicated_host_cluster" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  dedicated_host_cluster_name = "${var.name}"
}
```

## 参数参考

支持以下参数：
  * `dedicated_host_cluster_name` - (选填) - 专有宿主机集群名
  * `description` - (选填) - 描述
  * `zone_id` - (选填) - 可用区ID

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `dedicated_host_cluster_id` - 专有宿主机集群ID
