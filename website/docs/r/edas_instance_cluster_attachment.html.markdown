---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_instance_cluster_attachment"
sidebar_current: "docs-alibabacloudstack-resource-edas-instance-cluster-attachment"
description: |-
  Provides an EDAS instance cluster attachment resource.
---

# alibabacloudstack\_edas\_instance\_cluster\_attachment

Provides an EDAS instance cluster attachment resource.




## Example Usage

Basic Usage

```
resource "alibabacloudstack_edas_instance_cluster_attachment" "default" {
  cluster_id   = var.cluster_id
  instance_ids = var.instance_ids
}

```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the cluster that you want to create the application.
* `instance_ids` - (Required, ForceNew) The ID of instance. Type: list.
* `pass_word` - (Required, ForceNew) The login password for the ECS instance in the cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<cluster_id>:<instance_id1,instance_id2>`.
* `status_map` - The status map of the resource supplied above. The key is instance_id and the values are 1(running) 0(converting) -1(failed) and -2(offline).
* `ecu_map` - The ecu map of the resource supplied above. The key is instance_id and the value is ecu_id.
* `cluster_member_ids` - The cluster members map of the resource supplied above. The key is instance_id and the value is cluster_member_id.


