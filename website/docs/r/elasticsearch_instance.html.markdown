---
subcategory: "Elasticsearch"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_elasticsearch_instance"
sidebar_current: "docs-Alibabacloudstack-elasticsearch-instance"
description: |-
  Provides a Elasticsearch on K8s instance resource.
---

# alibabacloudstack_elasticsearch_instance

Provides a Elasticsearch instance resource.


The `alibabacloudstack_elasticsearch_instance` resource allows you to manage Elasticsearch instances on Alibaba Cloud.

## Example Usage

```hcl
variable "name" {
  default = "tf-testAccES"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "192.168.0.0/16"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_elasticsearch_instance" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  description = "tf-testAccES16373"
  data_node_spec = "1C 2Gi"
  client_node_amount = "2"
  monitor_password = "Pz*Pu3NIXuk("
  cpu_type = "Intel"
  data_node_disk_type = "yoda-lvm"
  master_node_disk_size = "100"
  master_node_disk_type = "yoda-lvm"
  data_node_disk_size = "500"
  password = "MSPFu4fXosn^"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
  master_node_spec = "1C 2Gi"
  master_node_amount = "3"
  client_node_spec = "1C 2Gi"
  data_node_amount = "3"
  kibana_node_spec = "1C 2Gi"
  kibana_password = "W7mI##At(q*M"
  version = "7.10.0_ali1.6.0"
  scene = "normal"
}
```

## Argument Reference
The following arguments are supported -

* `zone_id` - (Required, ForceNew) The ID of the zone to which the elasticsearch instance belong.
* `cpu_type` - (Required, ForceNew) The CPU type of the resource. Valid values: `intel`.
* `version` - (Required) The version of Elasticsearch to deploy.
* `description` - (Optional) The description of the Elasticsearch instance. It must be 0 to 30 characters in length and can contain numbers, letters, underscores, and hyphens. It must start with a letter, a number, or a Chinese character.
* `scense` - (Required, ForceNew) Application Scenarios. Valid values are "high", "normal", "log".
* `data_node_amount` - (Required) The number of data nodes in the Elasticsearch cluster. Valid range is from 2 to 50.
* `data_node_spec` - (Required) The specification of the data nodes.
* `data_node_disk_size` - (Required, ForceNew) The disk size of the data nodes.
* `data_node_disk_type` - (Required, ForceNew) The disk type of the data nodes. Valid values are "yoda-lvm", "fast-disks", "fast-disks-ssd".
* `data_node_affinity` - (Optional) Whether the data node disk is encrypted. Default value is false.
* `kibana_node_spec` - (Optional) The specification of the kibana nodes.
* `kibana_node_password` - (Optional) The password of the kibana nodes.
* `master_node_amount` - (Required) The number of master nodes in the Elasticsearch cluster. Valid range is from 3 to 50.
* `master_node_spec` - (Optional) The specification of the master nodes.
* `master_node_disk_size` - (Optional, ForceNew) The disk size of the master nodes.
* `master_node_disk_type` - (Optional, ForceNew) The disk type of the master nodes. Valid values are "yoda-lvm", "fast-disks", "fast-disks-ssd".
* `client_node_amount` - (Optional) The number of client nodes in the Elasticsearch cluster. Valid range is from 2 to 25.
* `client_node_spec` - (Optional) The specification of the client nodes.
* `vswitch_id` - (Optional) The ID of the VSwitch in which to launch the Elasticsearch instance.
* `password` - (Required) The password for the Elasticsearch instance.
* `monitor_password` - (Required) The password for the Elasticsearch monitor service.
* `protocol` - (Optional) The protocol used by the Elasticsearch instance. Valid values are HTTP and HTTPS. Default value is HTTP.
* `setting_config` - (Optional) - A map of settings to configure the Elasticsearch instance.

## Attributes Reference
The following attributes are exported -

* `slb_address` - The Actual IP Address of the Elasticsearch instance.
* `domain` - The domain name of the Elasticsearch instance.
* `port` - The port number of the Elasticsearch instance.
* `status` - The status of the Elasticsearch instance.
* `kibana_slb_address` - The Actual IP Address of the Kibana instance.
* `kibana_domain` - The domain name of the Kibana instance.
* `kibana_port` - The port number of the Kibana instance.

## Import
Elasticsearch instances can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_elasticsearch_instance.example i-1234567890abcdef0
```