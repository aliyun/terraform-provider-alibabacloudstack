---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_instance_engines"
sidebar_current: "docs-alibabacloudstack-datasource-kvstore-instance-engines"
description: |-
    Provides a list of KVStore instacne engines info.
---

# alibabacloudstack_kvstore_instance_engines

This data source provides the KVStore instance engines resource available info of Apsara Stack Cloud.

## Example Usage

```tf
data "alibabacloudstack_zones" "resources" {
  available_resource_creation = "KVStore"
}

data "alibabacloudstack_kvstore_instance_engines" "resources" {
  zone_id              = "${data.alibabacloudstack_zones.resources.zones.0.id}"
  instance_charge_type = "PrePaid"
  engine               = "Redis"
  engine_version       = "5.0"
  output_file          = "./engines.txt"
}

output "first_kvstore_instance_class" {
  value = "${data.alibabacloudstack_kvstore_instance_engines.resources.instance_engines.0.engine}"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required) The Zone to launch the KVStore instance.
* `instance_charge_type` - (Optional) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PrePaid`.
* `engine` - (Optional) Database type. Options are `Redis`, `Memcache`. Default to `Redis`.
* `engine_version` - (Optional) Database version required by the user. Value options of Redis can refer to the latest docs [detail info](https://www.alibabacloud.com/help/doc-detail/60873.htm) `EngineVersion`. Value of Memcache should be empty.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instance_engines` - A list of KVStore available instance engines. Each element contains the following attributes:
    * `zone_id` - The Zone to launch the KVStore instance.
    * `engine` - Database type.
    * `engine_version` - KVStore Instance version.