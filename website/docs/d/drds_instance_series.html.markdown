---
subcategory: "Distributed Relational Database Service(DRDS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_instance_series"
description: |-
  Provides a list of Drds Instance Series to be used by the alibabacloudstack_drds_instance resource.
---

# alibabacloudstack_drds_instance_series

This data source provides the Drds Instance Series families of AlibabacloudStack.

## Example Usage

```
data "alibabacloudstack_drds_instance_series" "default" {
  
}


```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) Specifies the ID range of instance series. If not specified, returns instance series across all availability zones.
* `names` - (Optional, ForceNew) Specifies the name range of instance series.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:


* `ids` - A list of instance series IDs matching all conditions.
* `names` - A list of instance series names matching all conditions.
* `series` - A list of detailed instance series specifications. Each element contains the following attributes: 
  * `id` - Unique identifier of the instance series.
  * `name` - Name of the instance series.
