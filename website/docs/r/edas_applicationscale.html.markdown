---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_application_scale"
sidebar_current: "docs-Alibabacloudstack-edas-application-scale"
description: |-
  Provides a edas Application scale resource.
---

# alibabacloudstack_edas_application_scale


Provides a edas Application scale resource.

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
	//image_id = "centos_7_7_x64_20G_alibase_20200426.vhd"
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
}
```

## Argument Reference

The following arguments are supported:
  * `ecc_info` - (ForceNew) - A list of ECU (Elastic Compute Unit) information for the application instances. This parameter specifies the ECUs that will be used for scaling the application.
  * `app_id` - (Required, ForceNew) - The ID of the application. You can query the ListApplication interface. For more information, see [ListApplication](~~ 149390 ~~).
  * `deploy_group` - (Required, ForceNew) - The group of application instances that need to be expanded. For more information about how to obtain an application group, see [QueryApplicationStatus](~~ 149394 ~~).
  * `force_status` - (Optional, ForceNew) - A flag to force the status update. This parameter can be used to enforce a specific status change during scaling operations.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `ecc_info` - A list of ECU (Elastic Compute Unit) information for the application instances. This parameter specifies the ECUs that will be used for scaling the application.
