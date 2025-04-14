---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_vservergroup"
sidebar_current: "docs-Alibabacloudstack-slb-vservergroup"
description: |- 
  编排负载均衡(SLB)虚拟服务器组
---

# alibabacloudstack_slb_vservergroup
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_server_group`

使用Provider配置的凭证在指定的资源集编排负载均衡(SLB)虚拟服务器组。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccslbv_server_group93962"
}

data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/16"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "instance" {
  image_id                   = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type              = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  instance_name              = "${var.name}"
  count                      = "2"
  security_groups            = ["${alibabacloudstack_security_group.default.id}"]
  internet_max_bandwidth_out = "10"
  availability_zone          = "${data.alibabacloudstack_zones.default.zones.0.id}"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb" "default" {
  name               = "${var.name}"
  address_type       = "internet"
  specification      = "slb.s2.small"
  vswitch_id         = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_vservergroup" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  name             = "${var.name}"

  servers {
    server_ids = ["${alibabacloudstack_instance.instance[0].id}", "${alibabacloudstack_instance.instance[1].id}"]
    port       = 100
    weight     = 10
    type       = "ecs"
  }

  servers {
    server_ids = ["${alibabacloudstack_instance.instance.*.id}"]
    port       = 80
    weight     = 100
    type       = "eni"
  }
}
```

## 参数说明

支持以下参数：

* `load_balancer_id` - (必填, 变更时重建) 负载均衡实例的ID。
* `name` - (选填) 虚拟服务器组的名称。如果未指定，则默认为资源名称。
* `vserver_group_name` - (选填) VServer组的名称。如果未指定，将使用`name`字段的值。
* `servers` - (选填) 要添加的ECS实例列表。一个资源中最多可以支持20个ECS实例。它包含以下几个子字段：
  * `server_ids` - (必填) 后端服务器ID(ECS实例ID)列表。
  * `port` - (必填) 后端服务器使用的端口。有效值范围：[1-65535]。
  * `weight` - (选填) 后端服务器的权重。有效值范围：[0-100]。默认值为100。
  * `type` - (选填) 后端服务器类型。有效值：`ecs`、`eni`。默认值为`ecs`。
* `delete_protection_validation` - (选填) 在删除之前检查SLB实例的DeleteProtection。如果为true，则当其SLB实例启用了DeleteProtection时，此资源将不会被删除。默认值为false。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 虚拟服务器组的ID。
* `load_balancer_id` - 用于启动新的虚拟服务器组的负载均衡器ID。
* `name` - 虚拟服务器组的名称。
* `vserver_group_name` - VServer组的名称。
* `servers` - 已经添加的ECS实例列表，包含每个实例的`server_ids`、`port`、`weight`和`type`信息。