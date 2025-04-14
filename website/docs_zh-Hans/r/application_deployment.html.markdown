---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_applicationpackageattachment"
sidebar_current: "docs-Alibabacloudstack-edas-applicationpackageattachment"
description: |-
  编排绑定企业级分布式应用服务（Edas）应用包
---

# alibabacloudstack_edas_applicationpackageattachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_application_deployment`

使用Provider配置的凭证在指定的资源集下编排绑定企业级分布式应用服务（Edas）应用包。

## 示例用法
```
variable "name" {
	default = "tf-testacc-edasdeploymentbasic3295"
}
variable "password" {
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
	// availability_zone = "cn-neimeng-env30-amtest30001-a"
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
	image_id = "centos_7_7_x64_20G_alibase_20211028.vhd"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	system_disk_category = "cloud_ssd"
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
	//region_id    = "cn-neimeng-env30-d01"
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
		

resource "alibabacloudstack_application_deployment" "default" {
  app_id = "${alibabacloudstack_edas_application.default.id}"
  group_id = "all"
  war_url = "http://edas-sz.oss-cn-shenzhen.aliyuncs.com/prod/demo/SPRING_CLOUD_CONSUMER.jar"
}
```

## 参数说明

支持以下参数：
  * `app_id` - (必填，变更时重建) - 应用程序的ID。您可以查询ListApplication接口以获取应用程序列表。有关更多信息，请参见[ListApplication](~~ 149390 ~~)。
  * `group_id` - (必填，变更时重建) - 部署组的ID。您可以查询ListDeployGroup操作以获取部署组列表。有关更多信息，请参见[ListDeployGroup](~~ 62077 ~~)。<note> 如果要部署到所有组，参数为' all '。</note>
  * `package_version` - (可选，变更时重建) - 部署包的版本号。最大长度为64个字符。建议使用时间戳作为版本号。
  * `war_url` - (必填，变更时重建) - 应用程序部署包(WAR或JAR)的URL。当部署类型为'url'时，此参数为必填项。建议使用OSS路径作为部署包地址。
  * `last_package_version` - (变更时重建) - 上一次部署包的版本号。此参数可用于跟踪和验证部署的历史版本。

## 属性说明

除了上述参数外，还导出以下属性：
  * `last_package_version` - 最后一次成功部署的应用程序包版本号。此属性可用于跟踪和验证部署的历史版本信息。