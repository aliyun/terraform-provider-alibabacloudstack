---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "阿里云栈: alibabacloudstack_cen_transit_router_vpc_attachment"
sidebar_current: "docs-Alibabacloudstack-resource-cen-transit-router-vpc-attachment"
description: |-
  提供一个cen Transitroutervpcattachment资源。
---

# alibabacloudstack\_cen\_transitroutervpcattachment

提供一个cen Transitroutervpcattachment资源。

## 使用示例
```
variable "name" {
	default = "tf-testaccrouter_vpc_attachment33915"
}


data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}




resource "alibabacloudstack_vpc_vswitch" "vswitchv2" {
	name = "${var.name}v2"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	cidr_block = "172.16.0.0/24"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  }

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}



resource "alibabacloudstack_cen_transit_router_vpc_attachment" "default" {
  zone_mappings {
    vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
    zone_id = "${alibabacloudstack_vpc_vswitch.default.zone_id}"
  }
  
  auto_create_vpc_route = "true"
  route_table_association_enabled = "true"
  route_table_propagation_enabled = "true"
  cen_id = "${alibabacloudstack_cen_instance.default.id}"
  vpc_id = "${alibabacloudstack_vpc_vswitch.default.vpc_id}"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}
```


## 参数引用

支持以下参数：
  * `cen_id` - (必填) - cen实例ID
  * `tags` - (可选) - 资源的标签
    
    * `tag_key` - (可选) - 标签键
    
    * `tag_value` - (可选) - 标签值
  * `transit_router_attachment_description` - (可选) - vpc附件的描述
  * `transit_router_id` - (必填) - 路由器ID
  * `transit_router_attachment_name` - (可选) - vpc附件的名称
  * `auto_create_vpc_route` - (可选) - 自动创建vpc路由条目
  * `route_table_propagation_enabled` - (可选) - 自动创建路由表传播
  * `route_table_association_enabled` - (可选) - 自动创建路由表关联
  * `vpc_id` - (必填) - vpc ID
  * `zone_mappings` - (必填) - 区域映射
    
    * `vswitch_id` - (必填) - 交换机ID
    
    * `zone_id` - (必填) - 区域ID
    

## 属性引用

除了上述列出的参数外，还导出以下属性：
  * `auto_publish_route_enabled` - 自动发布路由是否启用
  * `charge_type` - 计费类型
  * `creation_time` - 创建时间
  * `resource_type` - 资源类型
  * `status` - vpc附件实例的状态
  * `transit_router_attachment_id` - vpc附件实例的ID
  * `vpc_owner_id` - vpc所有者ID