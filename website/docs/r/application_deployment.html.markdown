---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_applicationpackageattachment"
sidebar_current: "docs-Alibabacloudstack-edas-applicationpackageattachment"
description: |-
  使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Applicationpackageattachment resource.
---

# alibabacloudstack_edas_applicationpackageattachment
-> **NOTE:** Alias name has: `alibabacloudstack_application_deployment`

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas） Applicationpackageattachment resource.

## Example Usage
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

## Argument Reference

The following arguments are supported:
  * `app_id` - (Required, ForceNew) - The ID of the application. You can query the ListApplication interface. For more information, see [ListApplication](~~ 149390 ~~).
  * `group_id` - (Required, ForceNew) - The ID of the deployment Group. You can query the ListDeployGroup operation. For more information, see [ListDeployGroup](~~ 62077 ~~).<note> If you want to deploy to all groups, the parameter is' all '. </note>
  * `package_version` - (Optional, ForceNew) - The version of the deployment package. The maximum length is 64 characters. We recommend that you use a timestamp.
  * `war_url` - (Required, ForceNew) - The URL of the application deployment package (WAR or JAR). **Deployytype** is required when it is 'url'. We recommend that you use the OSS application deployment package path.
  * `last_package_version` - (ForceNew) - The version of the last deployed package. This can be useful for tracking and verifying the version history of deployments.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `last_package_version` - The version of the last deployed package. This can be useful for tracking and verifying the version history of deployments.
