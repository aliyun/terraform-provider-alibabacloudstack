---
subcategory: "Elasticsearch"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_elasticsearch_instances"
sidebar_current: "docs-alibabacloudstack-datasource-elasticsearch-instances"
description: |-
  Provides a collection of Elasticsearch instances according to the specified filters.
---

# alibabacloudstack_elasticsearch_instances

The `alibabacloudstack_elasticsearch_instances` data source provides a collection of Elasticsearch instances available in alibabacloudstack account.
Filters support description regex, searches by tags, and other filters which are listed below.

## Example Usage

```
variable "name" {
  default = "tf-testacc-alikafkainstance18734"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_alikafka_instance" "default" {
  name = "${var.name}"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  sasl =      true
  plaintext = true
  spec =      "Broker4C16G"
}

data "alibabacloudstack_alikafka_instances" "default" {
  enable_details = "true"
  name_regex = "${alibabacloudstack_alikafka_instance.default.name}"
}
```

## Argument Reference

The following arguments are supported:

* `description_regex` - (Optional) A regex string to apply to the instance description.
* `ids` - (Optional, Available 1.52.1+) A list of Elasticsearch instance IDs.
* `version` - (Optional) Elasticsearch version. Options are `5.5.3_with_X-Pack`, `6.3.2_with_X-Pack` and `6.7.0_with_X-Pack`. If no value is specified, all versions are returned.
* `tags` - (Optional, Available 1.74.0+) A map of tags assigned to instances.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Elasticsearch instance IDs.
* `descriptions` - A list of Elasticsearch instance descriptions.
* `instances` - A list of Elasticsearch instances. Its every element contains the following attributes:
  * `id` - The ID of the Elasticsearch instance.
  * `zone_id` - The ID of the zone to which the elasticsearch instance belong.
  * `cpu_type` - The CPU type of the resource. 
  * `version` - The version of Elasticsearch to deploy.
  * `description` - The description of the Elasticsearch instance.
  * `scense` - Application Scenarios. 
  * `data_node_amount` - The number of data nodes in the Elasticsearch cluster.
  * `data_node_spec` - The specification of the data nodes.
  * `data_node_disk_size` - The disk size of the data nodes.
  * `data_node_disk_type` - The disk type of the data nodes.
  * `data_node_affinity` - Whether the data node disk is encrypted.
  * `kibana_node_spec` - The specification of the kibana nodes.
  * `kibana_node_password` - The password of the kibana nodes.
  * `master_node_amount` - The number of master nodes in the Elasticsearch cluster.
  * `master_node_spec` - The specification of the master nodes.
  * `master_node_disk_size` - The disk size of the master nodes
  * `master_node_disk_type` - The disk type of the master nodes.
  * `client_node_amount` - The number of client nodes in the Elasticsearch cluster.
  * `client_node_spec` - The specification of the client nodes.
  * `vswitch_id` - The ID of the VSwitch in which to launch the Elasticsearch instance.