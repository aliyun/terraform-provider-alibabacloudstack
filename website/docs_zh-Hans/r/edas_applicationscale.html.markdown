---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_application_scale"
sidebar_current: "docs-Alibabacloudstack-edas-application-scale"
description: |- 
  编排企业级分布式应用服务（Edas） 应用伸缩信息
---

# alibabacloudstack_edas_application_scale

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） 应用伸缩信息。

## 示例用法

```hcl
variable "password" {
}

variable "name" {
	default = "tf-testacc-edasiaattachment1441"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}	

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
	name = "${var.name}"
	description = "New security group"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "default" {
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	image_id="centos_7_7_x64_20G_alibase_20211028.vhd"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	system_disk_category = "cloud_efficiency"
	system_disk_size ="60"
	instance_type = "ecs.xn4.small"

	security_groups = ["${alibabacloudstack_security_group.default.id}"]
	instance_name = "${var.name}"
	tags = {
	Name = "TerraformTest-instance"
	}
}

resource "alibabacloudstack_edas_cluster" "default" {
	cluster_name = "${var.name}"
	cluster_type = 2
	network_mode = 2
	vpc_id       = "${alibabacloudstack_vpc.default.id}"
	region_id    = "cn-neimeng-env30-d01"
}

resource "alibabacloudstack_edas_instance_cluster_attachment" "default" {
	cluster_id = "${alibabacloudstack_edas_cluster.default.id}"
	instance_ids = ["${alibabacloudstack_instance.default.id}"]
	pass_word = "${var.password}"
}

resource "alibabacloudstack_edas_application" "default" {
	application_name = "${var.name}"
	cluster_id = "${alibabacloudstack_edas_cluster.default.id}"
	package_type = "JAR"
	ecu_info = ["${alibabacloudstack_edas_instance_cluster_attachment.default.ecu_map[alibabacloudstack_instance.default.id]}"]
	build_pack_id = "15"
}

data "alibabacloudstack_edas_deploy_groups" "default" {
	app_id = "${alibabacloudstack_edas_application.default.id}"
}

resource "alibabacloudstack_edas_application_scale" "default" {
  deploy_group = "${data.alibabacloudstack_edas_deploy_groups.default.groups.0.group_id}"
  ecu_info = [
               "${alibabacloudstack_edas_instance_cluster_attachment.default.ecu_map[alibabacloudstack_instance.default.id]}"
             ]
  app_id = "${alibabacloudstack_edas_application.default.id}"
  force_status = "true"
}
```

## 参数参考

支持以下参数：
  * `ecu_info` - (强制更新) - 应用程序实例的ECU(弹性计算单元)信息列表。此参数指定将用于扩展应用程序的ECUs。
  * `app_id` - (必填，强制更新) - 应用程序的ID。可以通过查询ListApplication接口获取更多信息，参见 [ListApplication](~~149390~~)。
  * `deploy_group` - (必填，强制更新) - 需要扩容的应用实例分组。获取应用实例分组，请参见 [QueryApplicationStatus](~~149394~~)。
  * `force_status` - (选填，强制更新) - 强制状态更新标志。此参数可用于在扩展操作期间强制执行特定的状态更改。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `ecu_info` - 应用程序实例的ECU(弹性计算单元)信息列表。此参数指定将用于扩展应用程序的ECUs。