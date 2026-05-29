---
subcategory: "Distributed Relational Database Service(DRDS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_instance_specifications"
sidebar_current: "docs-alibabacloudstack-drds-instance-specifications"
description: |-
  Provides a list of Drds Instance Series to be used by the alibabacloudstack_drds_instance resource.
---

# alibabacloudstack_drds_instance_specifications

This data source provides the Drds Instance specifications of AlibabacloudStack.

## Example Usage

```
data "alibabacloudstack_drds_instance_series" "default" {
  
}

data "alibabacloudstack_drds_instance_specifications" "default {
  series = data.alibabacloudstack_drds_instance_series.default
  cpu = 8
  memory = 32
}


```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) Specifies the ID range of instance specifications. If not specified, returns instance specifications across all availability zones.
* `names` - (Optional, ForceNew) Specifies the name range of instance specifications.
* `series` - (Optional)Filter the results to a specific drds series.
* `cpu` - (Optional) Filter the results to a specific number of cpu cores.
* `memory` - (Optional) Filter the results to a specific memory size in GB.
* `sorted_by` - (Optional, ForceNew) Sort mode, valid values: `CPU`, `Memory`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:


* `ids` - A list of instance series IDs matching all conditions.
* `names` - A list of instance series names matching all conditions.
* `specifications` - A list of detailed instance series specifications. Each element contains the following attributes: 
  * `id` - Unique identifier of the instance series.
  * `name` - Name of the instance series.
  * `cpu` - Number of cpu cores.
  * `memory` - Memory size in GB.
  * `series` - Drds series id.
