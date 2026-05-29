---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_read_write_splitting_config"
sidebar_current: "docs-Alibabacloudstack-resource-polardbx-read-write-splitting-config"
description: |-
  Provides a PolarDBX read-write splitting configuration resource.
---
# Alibaba Cloud PolarDBX Read Write Splitting Configuration

Provides a PolarDBX read-write splitting configuration resource.

## Example Usage

```hcl
resource "alibabacloudstack_polardbx_read_write_splitting_config" "example" {
  db_instance_id                  = "example-polardbx-instance-id"
  attend_htap_list                = ["example-polardbx-instance-id1", "example-polardbx-instance-id2"]
  auto_attend_htap                = true
  delay_execution_strategy        = 3
  enable_consistent_replica_read  = true
  enable_htap                     = true
  master_read_weight              = 50
  storage_delay_threshold         = 1000
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the PolarDBX instance.
* `attend_htap_list` - (Optional) List of HTAP instances to attend.
* `auto_attend_htap` - (Optional) Whether to automatically attend HTAP instances.
* `delay_execution_strategy` - (Optional) Delay execution strategy. Valid values are 0 and 1.
* `enable_consistent_replica_read` - (Optional) Whether to enable consistent replica read.
* `enable_htap` - (Optional) Whether to enable HTAP.
* `master_read_weight` - (Optional) The read weight of the master instance. Valid values are 1 to 100.
* `storage_delay_threshold` - (Optional) The storage delay threshold.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, same as `db_instance_id`.

## Import

PolarDBX read-write splitting configuration can be imported using the DB instance ID, e.g.

```bash
$ terraform import alibabacloudstack_polardbx_read_write_splitting_config.example <db_instance_id>
```