---
subcategory: "Managed Service for Prometheus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_instances"
description: |-
  Queries the list of Prometheus V2 instances in Alibaba Cloud.
---

# alibabacloudstack_prometheus_v2_instances

Queries the list of Prometheus V2 instances in Alibaba Cloud. This data source retrieves created Prometheus V2 monitoring instances and supports filtering by ID list and name regular expression.

## Example Usage

```hcl
variable "name" {
  default = "tfacc_prometheus9102323273752963708"
}

resource "alibabacloudstack_prometheus_v2_instance" "default" {
  cluster_name = var.name
  tags         = ["test1", "test2"]
}

data "alibabacloudstack_prometheus_v2_instances" "default" {
  name_regex = alibabacloudstack_prometheus_v2_instance.default.cluster_name
}
```

## Argument Reference

The following arguments support filtering query results:

* `ids` (List, Optional): A list of instance IDs used to precisely filter the Prometheus V2 instances to query.

* `name_regex` (String, Optional): A regular expression for the instance name, used to fuzzy match the names of Prometheus V2 instances to query.

## Attributes Reference

The following attributes are exported:

* `id` (String): The unique identifier of the Prometheus V2 instance, equivalent to cluster_id.

* `cluster_id` (String): The ID of the Prometheus V2 instance.

* `cluster_name` (String): The name of the Prometheus V2 instance.

* `cluster_type` (String): The type of the Prometheus V2 instance.

* `http_api` (String): The HTTP API endpoint address of the Prometheus V2 instance.

* `objid` (Integer): The internal object ID of the Prometheus V2 instance.

* `push_gateway_url` (String): The push gateway URL address of the Prometheus V2 instance.

* `remote_write_url` (String): The remote write URL address of the Prometheus V2 instance.

* `security_level_tag` (String): The security level tag of the instance.

* `status` (String): The current status of the instance.

* `tag_set` (List): A list of tags assigned to the instance.