---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "阿里云栈: alibabacloudstack_cen_transitroutervpcattachments"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transitroutervpcattachments"
description: |-
  提供阿里云栈账户拥有的cen transitroutervpcattachments列表。
---

# alibabacloudstack\_cen\_transitroutervpcattachments

该数据源根据指定的过滤条件提供阿里云栈账户中的cen transitroutervpcattachments列表。

## 使用示例
```
variable "name" {
  default = "tf-testAccRouterVpcAttachmentsDatasource18484"
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



resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_vpc_attachment" "default" {
	transit_router_attachment_name = "${var.name}"
	transit_router_attachment_description = "${var.name}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	vpc_id = "${alibabacloudstack_vpc_vswitch.default.vpc_id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
	zone_mappings {
			vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
			 zone_id = "${alibabacloudstack_vpc_vswitch.default.zone_id}"
		}
}

data "alibabacloudstack_cen_transit_router_vpc_attachments" "default" {
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	vpc_id = "${alibabacloudstack_vpc_vswitch.default.vpc_id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
	name_regex = "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}"
}
```


## 参数引用

支持以下参数：
  * `ids` - (可选) - vpc附件的ID列表
  * `cen_id` - (必填) - cen实例ID
  * `vpc_id` - (必填) - vpc ID
  * `tags` - (可选) - 资源的标签
    
    * `tag_key` - (可选) - 标签键
    
    * `tag_value` - (可选) - 标签值
  * `transit_router_id` - (必填) - 路由器ID
  * `name_regex` - (可选) - vpc附件的名称正则表达式
  * `description_regex` - (可选) - vpc附件的描述正则表达式

## 属性引用

除了上述列出的参数外，还导出以下属性：
  * `transitrouterattachments` - vpc附件列表