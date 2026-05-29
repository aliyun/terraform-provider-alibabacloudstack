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

```hcl
data "alibabacloudstack_elasticsearch_instances" "default" {
  description_regex = "my-es-instance"
  version           = "6.7.0_with_X-Pack"
  vpc_id            = "vpc-xxxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `description_regex` - (Optional) A regex string to apply to the instance description list for filtering.
* `ids` - (Optional) A list of Elasticsearch instance IDs.
* `version` - (Optional) Elasticsearch version. If not specified, instances of all versions are returned.
* `vpc_id` - (Optional) The ID of the VPC to which the Elasticsearch instances belong.
* `output_file` - (Optional, Deprecated) The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the `local_file` provider instead.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Elasticsearch instance IDs.
* `descriptions` - A list of Elasticsearch instance descriptions.
* `instances` - A list of Elasticsearch instances. Each element contains the following attributes:
  * `id` - The ID of the Elasticsearch instance.
  * `version` - The version of Elasticsearch.
  * `description` - The description of the Elasticsearch instance.
  * `data_node_amount` - The number of data nodes in the Elasticsearch cluster.
  * `data_node_spec` - The specification of the data nodes.
  * `data_node_disk_size` - The disk size of the data nodes.
  * `data_node_disk_type` - The disk type of the data nodes.
  * `kibana_node_spec` - The specification of the Kibana nodes (only available when Kibana is enabled).
  * `kibana_slb_address` - The SLB address of Kibana (only available when Kibana is enabled).
  * `kibana_domain` - The domain name of Kibana (only available when Kibana is enabled).
  * `kibana_protocol` - The protocol used by Kibana (only available when Kibana is enabled).
  * `kibana_port` - The port number of Kibana (only available when Kibana is enabled).
  * `master_node_amount` - The number of dedicated master nodes in the Elasticsearch cluster (only available when dedicated master is enabled).
  * `master_node_spec` - The specification of the dedicated master nodes.
  * `master_node_disk_size` - The disk size of the dedicated master nodes.
  * `master_node_disk_type` - The disk type of the dedicated master nodes.
  * `client_node_amount` - The number of client nodes in the Elasticsearch cluster (only available when client nodes are enabled).
  * `client_node_spec` - The specification of the client nodes.
  * `vswitch_id` - The ID of the VSwitch in which the Elasticsearch instance is launched.
  * `status` - The status of the Elasticsearch instance.