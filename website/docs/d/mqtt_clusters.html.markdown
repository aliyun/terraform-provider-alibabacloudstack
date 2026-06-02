---
subcategory: "Message Queuing Telemetry Transport"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_clusters"
description: |-
  Provides a list of MQTT Clusters.
---

# alibabacloudstack_mqtt_clusters

This data source provides the MQTT Clusters available in ApsaraStack.

## Example Usage

```hcl
data "alibabacloudstack_mqtt_clusters" "example" {
  name_regex = "^cluster-.*"
}

output "mqtt_clusters" {
  value = data.alibabacloudstack_mqtt_clusters.example.clusters
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of cluster IDs to filter results.
* `name_regex` - (Optional, ForceNew) A regex string to filter clusters by name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of cluster IDs.
* `clusters` - A list of MQTT Clusters. Each element contains the following attributes:
  * `id` - The ID of the cluster.
  * `name` - The name of the cluster.
