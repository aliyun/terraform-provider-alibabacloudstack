---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_slbattachments"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-slbattachments"
description: |-
  Provides a list of edas slbattachments owned by an alibabacloudstack account.
---

# alibabacloudstack_edas_slbattachments
-> **NOTE:** Alias name has: `alibabacloudstack_edas_applications`

This data source provides a list of edas slbattachments in an alibabacloudstack account according to the specified filters.

## Example Usage
```
		variable "name" {
		  default = "%v"
		}
    variable "password" {
    }
		data "alibabacloudstack_zones" "default" {
			available_resource_creation= "VSwitch"
		}

		resource "alibabacloudstack_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = "${var.name}"
		}
		
		resource "alibabacloudstack_vswitch" "default" {
		  vpc_id            = "${alibabacloudstack_vpc.default.id}"
		  cidr_block        = "172.16.0.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}"
		}

		resource "alibabacloudstack_slb" "default" {
		  name          = "${var.name}"
		  vswitch_id    = "${alibabacloudstack_vswitch.default.id}"
      address_type  = "internet"
		  specification = "slb.s1.small"
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
		}
		
		resource "alibabacloudstack_edas_instance_cluster_attachment" "default" {
      cluster_id = "${alibabacloudstack_edas_cluster.default.id}"
      instance_ids = ["${alibabacloudstack_instance.default.id}"]
      pass_word = var.password
		}
		
		resource "alibabacloudstack_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = "${alibabacloudstack_edas_cluster.default.id}"
		  package_type = "JAR"
		  //ecu_info = ["${alibabacloudstack_edas_instance_cluster_attachment.default.ecu_map[alibabacloudstack_instance.default.id]}"]
		  ecu_info = ["${alibabacloudstack_edas_instance_cluster_attachment.default.ecu_map[alibabacloudstack_instance.default.id]}"]
		}

    resource "alibabacloudstack_edas_slb_attachment" "default" {
      app_id =        "${alibabacloudstack_edas_application.default.id}"
      slb_id =        "${alibabacloudstack_slb.default.id}"
      slb_ip =        "${alibabacloudstack_slb.default.address}"
      type   =        "${alibabacloudstack_slb.default.address_type}"
      listener_port = "22",
    }
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew) A list of Namespace IDs. Used to filter the results by specific namespace IDs.
* `names` - (Optional, ForceNew) A list of names of the SLB attachments.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Namespace name. This can be useful when you want to find namespaces that match a specific naming pattern.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `names` - A list of names of the SLB attachments.
  * `applications` -  A list of applications associated with the SLB attachments. Each element in the list is a map with the following keys:
    - `app_id` - The ID of the application.
    - `slb_id` - The ID of the SLB.
    - `slb_ip` - The IP address of the SLB.
    - `type` - The type of the SLB (e.g., "internet" or "intranet").
    - `listener_port` - The listener port of the SLB.
    - `vserver_group_id` - The ID of the VServer group associated with the SLB.
    - `slb_status` - The status of the SLB.
    - `vswitch_id` - The ID of the VSwitch associated with the SLB.