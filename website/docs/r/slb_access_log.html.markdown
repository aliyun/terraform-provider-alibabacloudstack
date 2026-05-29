---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_access_log"
sidebar_current: "docs-Alibabacloudstack-slb-access-log"
description: |-
  Provides a slb Access logs resource.
---

# alibabacloudstack\_slb\_access_log

Provides a slb Access logs resource.

## Example Usage
```
variable "name" {
  default = "tf-test-log88673"
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


resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  	cidr_ip = "192.168.0.0/16"
}


data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  //name_regex  = "arm_centos_7_6_20G_20211110.raw"
  //name_regex  = "^arm_centos_7"
  most_recent = true
  owners      = "system"
}


data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
	default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

resource "alibabacloudstack_slb" "default" {
  name          = "${var.name}_slb"
  vswitch_id    = "${alibabacloudstack_vpc_vswitch.default.id}"
}

resource "alibabacloudstack_slb_server_certificate" "servercertificate" {
  name               = "slbservercertificate"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq************bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIIDRjCCAq************bJJyOm5LqoiA=\n-----END RSA PRIVATE KEY-----"
}

resource "alibabacloudstack_slb_listener" "default" {
    load_balancer_id            = alibabacloudstack_slb.default.id
    server_certificate_id       = alibabacloudstack_slb_server_certificate.servercertificate.id
    sticky_session              = "off"
    sticky_session_type         = "insert"
    cookie_timeout              = 1000
    frontend_port               = 120
    backend_port                = 120
    enable_http2                = "on"
    acl_status                  = "off"
    acl_type                    = "white"
    protocol                    = "https"
    bandwidth                   = -1
    gzip                        = true
    x_forwarded_for {
        retrive_slb_ip = false
        retrive_slb_id = false
        retrive_slb_proto = false
    }
    tls_cipher_policy           = "tls_cipher_policy_1_2"
    health_check                = "on"
    health_check_type           = "http"
    health_check_uri            = "/"
    health_check_connect_port   = 20
    health_check_method         = "head"
    healthy_threshold           = "3"
    unhealthy_threshold         = "3"
    health_check_timeout        = "5"
    health_check_interval       = "2"
    health_check_http_code      = "http_2xx,http_3xx"
    description                 = "testslblistener"
}

resource "alibabacloudstack_log_project" "default" {
		name = "${var.name}_project"
		description = "test"
	}
	resource "alibabacloudstack_log_store" "default" {
		name = "${var.name}_store"
		project = "${alibabacloudstack_log_project.default.name}"	
		retention_period      = "30"
		shard_count           = "2"
		enable_web_tracking   = false
		auto_split            = true
		max_split_shard_count = "64"
		append_meta           = true
	}



resource "alibabacloudstack_slb_access_log" "default" {
  log_project = "${alibabacloudstack_log_project.default.name}"
  log_store = "${alibabacloudstack_log_store.default.name}"
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  depends_on = [
                 "alibabacloudstack_slb_listener.default"
               ]
}
```

## Argument Reference

The following arguments are supported:
  * `load_balancer_id` - (Optional, ForceNew) - Slb instance id
  * `log_project` - (Optional, ForceNew) - LogProject of user SLS
  * `log_store` - (Optional, ForceNew) - LogStore of user SLS
  * `log_type` - (Optional) - Access Log. The default is layer7
  * `role_name` - (Optional) - Your role name in sls: Default is aliyunlogwriteonlyrole

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
